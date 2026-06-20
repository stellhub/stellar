package circuitbreaker

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	apperrors "github.com/stellhub/stellar/errors"
	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
)

const Name = "governance-circuit-breaker"

type Policy struct {
	store    *governance.Store
	metrics  *governance.Metrics
	breakers sync.Map
}

type Settings struct {
	FailureThreshold    int
	OpenTimeout         time.Duration
	HalfOpenMaxRequests int
}

type breakerState string

const (
	stateClosed   breakerState = "closed"
	stateOpen     breakerState = "open"
	stateHalfOpen breakerState = "half_open"
)

type breaker struct {
	mu       sync.Mutex
	state    breakerState
	failures int
	openedAt time.Time
	probes   int
}

func NewPolicy(store *governance.Store, metrics *governance.Metrics) *Policy {
	return &Policy{store: store, metrics: metrics}
}

func Definitions(store *governance.Store, metrics *governance.Metrics) []interceptor.Definition {
	policy := NewPolicy(store, metrics)
	return []interceptor.Definition{
		interceptor.Framework(interceptor.KindHTTPClient, interceptor.StageAdmission, Name, policy),
		interceptor.Framework(interceptor.KindGRPCClient, interceptor.StageAdmission, Name, policy),
	}
}

func (p *Policy) Name() string {
	return Name
}

func (p *Policy) Intercept(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
	if p == nil || p.store == nil {
		return next(ctx, inv, req)
	}
	rules := p.store.Rules(governance.RuleKindCircuitBreaker, func(rule governance.Rule) bool {
		return rule.Enabled && invocationMatches(rule.Scope, inv)
	})
	if len(rules) == 0 {
		return next(ctx, inv, req)
	}

	rule := rules[0]
	resource := resourceKey(inv, rule)
	settings := settingsFromSpec(rule.Spec)
	current := p.breakerFor(rule.ID + ":" + resource)
	if !current.allow(settings) {
		if p.metrics != nil {
			p.metrics.RecordCircuitBreakerReject(ctx, governance.MetricAttrs{
				Adapter:   "stellorbit",
				RuleKind:  string(rule.Kind),
				RuleID:    rule.ID,
				Transport: invocationTransport(inv),
				Service:   invocationService(inv),
				Method:    invocationMethod(inv),
				Resource:  resource,
				Outcome:   "rejected",
			})
		}
		return nil, apperrors.New(apperrors.CodeUnavailable, "stellar: circuit breaker is open", http.StatusServiceUnavailable)
	}

	resp, err := next(ctx, inv, req)
	changed := current.after(settings, isFailure(resp, err))
	if changed && p.metrics != nil {
		p.metrics.RecordCircuitBreakerStateChange(ctx, governance.MetricAttrs{
			Adapter:   "stellorbit",
			RuleKind:  string(rule.Kind),
			RuleID:    rule.ID,
			Transport: invocationTransport(inv),
			Service:   invocationService(inv),
			Method:    invocationMethod(inv),
			Resource:  resource,
			Outcome:   string(current.stateSnapshot()),
		})
	}
	return resp, err
}

func (p *Policy) breakerFor(key string) *breaker {
	value, _ := p.breakers.LoadOrStore(key, &breaker{state: stateClosed})
	return value.(*breaker)
}

func (b *breaker) allow(settings Settings) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	switch b.state {
	case stateOpen:
		if now.Sub(b.openedAt) < settings.OpenTimeout {
			return false
		}
		b.state = stateHalfOpen
		b.probes = 0
	case stateHalfOpen:
	default:
		b.state = stateClosed
	}
	if b.state == stateHalfOpen {
		if b.probes >= settings.HalfOpenMaxRequests {
			return false
		}
		b.probes++
	}
	return true
}

func (b *breaker) after(settings Settings, failed bool) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	old := b.state
	if failed {
		b.failures++
		if b.state == stateHalfOpen || b.failures >= settings.FailureThreshold {
			b.state = stateOpen
			b.openedAt = time.Now()
			b.probes = 0
		}
		return old != b.state
	}
	b.failures = 0
	if b.state == stateHalfOpen {
		b.state = stateClosed
		b.probes = 0
	}
	return old != b.state
}

func (b *breaker) stateSnapshot() breakerState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func settingsFromSpec(spec map[string]any) Settings {
	return Settings{
		FailureThreshold:    intSpec(spec, 5, "failure_threshold", "failureThreshold", "threshold"),
		OpenTimeout:         durationSpec(spec, 30*time.Second, "open_timeout", "openTimeout", "sleep_window", "sleepWindow"),
		HalfOpenMaxRequests: intSpec(spec, 1, "half_open_max_requests", "halfOpenMaxRequests", "half_open_requests"),
	}
}

func invocationMatches(scope governance.Scope, inv *interceptor.Invocation) bool {
	if inv == nil {
		return true
	}
	if scope.Transport != "" && !strings.EqualFold(scope.Transport, string(inv.Kind)) && !strings.EqualFold(scope.Transport, inv.Protocol+".client") && !strings.EqualFold(scope.Transport, inv.Protocol+".server") {
		return false
	}
	if scope.Service != "" && !strings.EqualFold(scope.Service, inv.Service) {
		return false
	}
	if scope.Method != "" && !strings.EqualFold(scope.Method, inv.Method) {
		return false
	}
	if scope.Path != "" && !strings.HasPrefix(inv.Path, scope.Path) {
		return false
	}
	if scope.Target != "" && !strings.EqualFold(scope.Target, inv.Target) {
		return false
	}
	return true
}

func isFailure(resp any, err error) bool {
	if err != nil {
		return true
	}
	if response, ok := resp.(*http.Response); ok && response != nil {
		return response.StatusCode >= http.StatusInternalServerError
	}
	return false
}

func resourceKey(inv *interceptor.Invocation, rule governance.Rule) string {
	if value := stringSpec(rule.Spec, "resource", "resource_key", "resourceKey"); value != "" {
		return value
	}
	if inv == nil {
		return rule.ID
	}
	return strings.Trim(strings.Join([]string{inv.Service, inv.Method, inv.Path}, ":"), ":")
}

func invocationTransport(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return string(inv.Kind)
}

func invocationService(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Service
}

func invocationMethod(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Method
}

func intSpec(spec map[string]any, fallback int, keys ...string) int {
	for _, key := range keys {
		value, ok := spec[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case int:
			if typed > 0 {
				return typed
			}
		case int64:
			if typed > 0 {
				return int(typed)
			}
		case float64:
			if typed > 0 {
				return int(typed)
			}
		case string:
			parsed, err := strconv.Atoi(strings.TrimSpace(typed))
			if err == nil && parsed > 0 {
				return parsed
			}
		}
	}
	return fallback
}

func durationSpec(spec map[string]any, fallback time.Duration, keys ...string) time.Duration {
	for _, key := range keys {
		value := stringSpec(spec, key)
		if value == "" {
			continue
		}
		duration, err := time.ParseDuration(value)
		if err == nil && duration > 0 {
			return duration
		}
	}
	return fallback
}

func stringSpec(spec map[string]any, keys ...string) string {
	for _, key := range keys {
		if spec == nil {
			return ""
		}
		value, ok := spec[key]
		if !ok {
			continue
		}
		text := strings.TrimSpace(fmt.Sprint(value))
		if text != "" {
			return text
		}
	}
	return ""
}
