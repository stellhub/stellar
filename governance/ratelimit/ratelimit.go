package ratelimit

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
	stellpulsar "github.com/stellhub/stellpulsar-go-sdk"
)

const Name = "governance-rate-limit"

type LimitBehavior string

const (
	LimitBehaviorReject LimitBehavior = "reject"
	LimitBehaviorBlock  LimitBehavior = "block"
)

type Limiter interface {
	Allow(ctx context.Context, request LimitRequest) (LimitDecision, error)
	Wait(ctx context.Context, request LimitRequest) (LimitDecision, error)
}

type LimitRequest struct {
	Rule       governance.Rule
	Transport  string
	Service    string
	Method     string
	Path       string
	Target     string
	Resource   string
	QuotaKey   string
	Cost       int64
	Rate       int64
	Burst      int64
	Attributes map[string]string
}

type LimitDecision struct {
	Allowed    bool
	Limited    bool
	RuleID     string
	Reason     string
	RetryAfter time.Duration
	Fallback   bool
	RemoteCode string
}

type PolicyOptions struct {
	Store               *governance.Store
	Metrics             *governance.Metrics
	Local               Limiter
	Distributed         Limiter
	DefaultMode         string
	DefaultBehavior     LimitBehavior
	DistributedFallback string
	DefaultRate         int64
	DefaultBurst        int64
}

type Policy struct {
	store           *governance.Store
	metrics         *governance.Metrics
	local           Limiter
	distributed     Limiter
	defaultMode     string
	defaultBehavior LimitBehavior
	defaultRate     int64
	defaultBurst    int64
}

func NewPolicy(options PolicyOptions) *Policy {
	defaultMode := strings.ToLower(strings.TrimSpace(options.DefaultMode))
	if defaultMode == "" {
		defaultMode = "local"
	}
	behavior := options.DefaultBehavior
	if behavior == "" {
		behavior = LimitBehaviorReject
	}
	local := options.Local
	if local == nil {
		local = NewLocalLimiter()
	}
	distributed := options.Distributed
	if distributed == nil {
		distributed = NewDistributedLimiter(DistributedOptions{Fallback: options.DistributedFallback})
	}
	defaultRate := options.DefaultRate
	if defaultRate <= 0 {
		defaultRate = 100
	}
	defaultBurst := options.DefaultBurst
	if defaultBurst <= 0 {
		defaultBurst = defaultRate
	}
	return &Policy{
		store:           options.Store,
		metrics:         options.Metrics,
		local:           local,
		distributed:     distributed,
		defaultMode:     defaultMode,
		defaultBehavior: behavior,
		defaultRate:     defaultRate,
		defaultBurst:    defaultBurst,
	}
}

func Definitions(options PolicyOptions) []interceptor.Definition {
	policy := NewPolicy(options)
	return []interceptor.Definition{
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageAdmission, Name, policy),
		interceptor.Framework(interceptor.KindGRPCServer, interceptor.StageAdmission, Name, policy),
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
	rules := p.store.Rules(governance.RuleKindRateLimit, func(rule governance.Rule) bool {
		return rule.Enabled && invocationMatches(rule.Scope, inv)
	})
	for _, rule := range rules {
		request := p.limitRequest(inv, rule)
		behavior := behaviorFromSpec(rule.Spec, p.defaultBehavior)
		limiter := p.limiterFor(modeFromSpec(rule.Spec, p.defaultMode))
		if limiter == nil {
			continue
		}
		start := time.Now()
		decision, err := decide(ctx, limiter, behavior, request)
		if behavior == LimitBehaviorBlock {
			p.recordWait(ctx, request, behavior, time.Since(start), decision)
		}
		if err != nil {
			return nil, err
		}
		if !decision.Allowed {
			p.recordReject(ctx, request, behavior, decision)
			return nil, apperrors.New(apperrors.CodeUnavailable, "stellar: rate limit exceeded", http.StatusTooManyRequests)
		}
		p.recordAllow(ctx, request, behavior, decision)
	}
	return next(ctx, inv, req)
}

func decide(ctx context.Context, limiter Limiter, behavior LimitBehavior, request LimitRequest) (LimitDecision, error) {
	switch behavior {
	case LimitBehaviorBlock:
		return limiter.Wait(ctx, request)
	default:
		return limiter.Allow(ctx, request)
	}
}

func (p *Policy) limiterFor(mode string) Limiter {
	switch normalizeMode(mode) {
	case "distributed", "global-sync", "global-quota", "edge":
		return p.distributed
	default:
		return p.local
	}
}

func (p *Policy) limitRequest(inv *interceptor.Invocation, rule governance.Rule) LimitRequest {
	resource := resourceKey(inv, rule)
	attrs := invocationAttributes(inv)
	cost := int64Spec(rule.Spec, 1, "cost")
	rate := int64Spec(rule.Spec, p.defaultRate, "rate", "qps", "quota")
	burst := int64Spec(rule.Spec, p.defaultBurst, "burst")
	if burst <= 0 {
		burst = rate
	}
	return LimitRequest{
		Rule:       rule,
		Transport:  invocationTransport(inv),
		Service:    invocationService(inv),
		Method:     invocationMethod(inv),
		Path:       invocationPath(inv),
		Target:     invocationTarget(inv),
		Resource:   resource,
		QuotaKey:   quotaKey(inv, rule, resource),
		Cost:       cost,
		Rate:       rate,
		Burst:      burst,
		Attributes: attrs,
	}
}

func (p *Policy) recordAllow(ctx context.Context, request LimitRequest, behavior LimitBehavior, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitAllow(ctx, metricAttrs(request, behavior, decision, "allowed"))
}

func (p *Policy) recordReject(ctx context.Context, request LimitRequest, behavior LimitBehavior, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitReject(ctx, metricAttrs(request, behavior, decision, "rejected"))
}

func (p *Policy) recordWait(ctx context.Context, request LimitRequest, behavior LimitBehavior, duration time.Duration, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitWait(ctx, duration, metricAttrs(request, behavior, decision, "wait"))
}

func metricAttrs(request LimitRequest, behavior LimitBehavior, decision LimitDecision, outcome string) governance.MetricAttrs {
	if decision.Fallback {
		outcome = outcome + ".fallback"
	}
	return governance.MetricAttrs{
		Adapter:   "stellorbit",
		RuleKind:  string(request.Rule.Kind),
		RuleID:    request.Rule.ID,
		Transport: request.Transport,
		Service:   request.Service,
		Method:    request.Method,
		Resource:  request.Resource,
		Mode:      modeFromSpec(request.Rule.Spec, ""),
		Behavior:  string(behavior),
		Outcome:   outcome,
	}
}

type LocalLimiter struct {
	buckets sync.Map
}

type bucket struct {
	mu     sync.Mutex
	tokens float64
	last   time.Time
}

func NewLocalLimiter() *LocalLimiter {
	return &LocalLimiter{}
}

func (l *LocalLimiter) Allow(_ context.Context, request LimitRequest) (LimitDecision, error) {
	current := l.bucketFor(request)
	allowed, retryAfter := current.take(request.Rate, request.Burst, request.Cost)
	if allowed {
		return allowedDecision(request.Rule.ID), nil
	}
	return limitedDecision(request.Rule.ID, "local_rate_limited", retryAfter), nil
}

func (l *LocalLimiter) Wait(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		decision, err := l.Allow(ctx, request)
		if err != nil || decision.Allowed {
			return decision, err
		}
		delay := decision.RetryAfter
		if delay <= 0 {
			delay = 10 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return limitedDecision(request.Rule.ID, ctx.Err().Error(), 0), nil
		case <-timer.C:
		}
	}
}

func (l *LocalLimiter) bucketFor(request LimitRequest) *bucket {
	key := request.Rule.ID + ":" + request.Resource + ":" + request.QuotaKey
	value, _ := l.buckets.LoadOrStore(key, &bucket{tokens: float64(request.Burst), last: time.Now()})
	return value.(*bucket)
}

func (b *bucket) take(rate int64, burst int64, cost int64) (bool, time.Duration) {
	if rate <= 0 {
		return true, 0
	}
	if burst <= 0 {
		burst = rate
	}
	if cost <= 0 {
		cost = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	if b.last.IsZero() {
		b.last = now
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * float64(rate)
	if maxTokens := float64(burst); b.tokens > maxTokens {
		b.tokens = maxTokens
	}
	b.last = now
	if b.tokens >= float64(cost) {
		b.tokens -= float64(cost)
		return true, 0
	}
	missing := float64(cost) - b.tokens
	wait := time.Duration(missing / float64(rate) * float64(time.Second))
	if wait <= 0 {
		wait = time.Millisecond
	}
	return false, wait
}

type DistributedLimiter struct {
	client   stellpulsar.StellpulsarClient
	timeout  time.Duration
	fallback string
	appID    string
}

type DistributedOptions struct {
	Client   stellpulsar.StellpulsarClient
	Timeout  time.Duration
	Fallback string
	AppID    string
}

func NewDistributedLimiter(options DistributedOptions) *DistributedLimiter {
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 100 * time.Millisecond
	}
	fallback := strings.ToLower(strings.TrimSpace(options.Fallback))
	if fallback == "" {
		fallback = "fail_open"
	}
	return &DistributedLimiter{
		client:   options.Client,
		timeout:  timeout,
		fallback: fallback,
		appID:    options.AppID,
	}
}

func (l *DistributedLimiter) Allow(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if l == nil || l.client == nil {
		return l.fallbackDecision(request, "stellpulsar_client_missing"), nil
	}
	callCtx := ctx
	cancel := func() {}
	if l.timeout > 0 {
		callCtx, cancel = context.WithTimeout(ctx, l.timeout)
	}
	defer cancel()
	result, err := l.client.TryAcquire(callCtx, stellpulsar.RateLimitRequest{
		RequestID:       requestID(request),
		ApplicationCode: l.applicationCode(request),
		TargetService:   request.Service,
		Resource:        request.Resource,
		Method:          request.Method,
		TenantID:        request.Attributes["tenant_id"],
		UserID:          request.Attributes["user_id"],
		QuotaKey:        request.QuotaKey,
		Cost:            request.Cost,
		Attributes:      request.Attributes,
	})
	if err != nil {
		return l.fallbackDecision(request, err.Error()), nil
	}
	return decisionFromPulsar(request, result), nil
}

func (l *DistributedLimiter) Wait(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		decision, err := l.Allow(ctx, request)
		if err != nil || decision.Allowed {
			return decision, err
		}
		delay := decision.RetryAfter
		if delay <= 0 {
			delay = 50 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return l.fallbackDecision(request, ctx.Err().Error()), nil
		case <-timer.C:
		}
	}
}

func (l *DistributedLimiter) applicationCode(request LimitRequest) string {
	if strings.TrimSpace(l.appID) != "" {
		return strings.TrimSpace(l.appID)
	}
	return request.Service
}

func (l *DistributedLimiter) fallbackDecision(request LimitRequest, reason string) LimitDecision {
	allowed := l == nil || l.fallback != "fail_closed"
	return LimitDecision{
		Allowed:    allowed,
		Limited:    !allowed,
		RuleID:     request.Rule.ID,
		Reason:     reason,
		Fallback:   true,
		RemoteCode: "CLIENT_ERROR",
	}
}

func decisionFromPulsar(request LimitRequest, result *stellpulsar.RateLimitResult) LimitDecision {
	if result == nil {
		return LimitDecision{Allowed: true, RuleID: request.Rule.ID, Reason: "no_result"}
	}
	return LimitDecision{
		Allowed:    result.Permitted,
		Limited:    result.Limited,
		RuleID:     firstNonBlank(result.RuleID, request.Rule.ID),
		Reason:     result.Reason,
		RetryAfter: time.Duration(result.RetryAfterMS) * time.Millisecond,
		Fallback:   result.Fallback,
		RemoteCode: string(result.Decision),
	}
}

func allowedDecision(ruleID string) LimitDecision {
	return LimitDecision{Allowed: true, RuleID: ruleID, Reason: "allowed"}
}

func limitedDecision(ruleID string, reason string, retryAfter time.Duration) LimitDecision {
	return LimitDecision{Allowed: false, Limited: true, RuleID: ruleID, Reason: reason, RetryAfter: retryAfter}
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

func behaviorFromSpec(spec map[string]any, fallback LimitBehavior) LimitBehavior {
	switch strings.ToLower(strings.TrimSpace(stringSpec(spec, "behavior", "limit_behavior"))) {
	case string(LimitBehaviorBlock):
		return LimitBehaviorBlock
	case string(LimitBehaviorReject):
		return LimitBehaviorReject
	default:
		if fallback == "" {
			return LimitBehaviorReject
		}
		return fallback
	}
}

func modeFromSpec(spec map[string]any, fallback string) string {
	if value := stringSpec(spec, "mode", "backend"); value != "" {
		return normalizeMode(value)
	}
	return normalizeMode(fallback)
}

func normalizeMode(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), "_", "-"))
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

func quotaKey(inv *interceptor.Invocation, rule governance.Rule, resource string) string {
	if value := stringSpec(rule.Spec, "quota_key", "quotaKey"); value != "" {
		return value
	}
	if inv != nil {
		if value := inv.Headers.Get("x-stellar-quota-key"); value != "" {
			return value
		}
		if value := inv.Headers.Get("x-tenant-id"); value != "" {
			return value
		}
		if inv.Attributes != nil {
			for _, key := range []string{"quota_key", "quotaKey", "tenant_id", "tenantId", "user_id", "userId"} {
				if value := fmt.Sprint(inv.Attributes[key]); value != "" && value != "<nil>" {
					return strings.TrimSpace(value)
				}
			}
		}
	}
	return resource
}

func invocationAttributes(inv *interceptor.Invocation) map[string]string {
	attrs := map[string]string{}
	if inv == nil {
		return attrs
	}
	put(attrs, "transport", string(inv.Kind))
	put(attrs, "service", inv.Service)
	put(attrs, "method", inv.Method)
	put(attrs, "path", inv.Path)
	put(attrs, "target", inv.Target)
	put(attrs, "tenant_id", inv.Headers.Get("x-tenant-id"))
	put(attrs, "user_id", inv.Headers.Get("x-user-id"))
	put(attrs, "quota_key", inv.Headers.Get("x-stellar-quota-key"))
	for key, value := range inv.Attributes {
		put(attrs, key, fmt.Sprint(value))
	}
	return attrs
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

func invocationPath(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Path
}

func invocationTarget(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Target
}

func requestID(request LimitRequest) string {
	if value := request.Attributes["request_id"]; value != "" {
		return value
	}
	return request.Rule.ID + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func int64Spec(spec map[string]any, fallback int64, keys ...string) int64 {
	for _, key := range keys {
		if spec == nil {
			return fallback
		}
		value, ok := spec[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case int:
			if typed > 0 {
				return int64(typed)
			}
		case int64:
			if typed > 0 {
				return typed
			}
		case float64:
			if typed > 0 {
				return int64(typed)
			}
		case string:
			parsed, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
			if err == nil && parsed > 0 {
				return parsed
			}
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
		if text != "" && text != "<nil>" {
			return text
		}
	}
	return ""
}

func put(values map[string]string, key string, value string) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key != "" && value != "" && value != "<nil>" {
		values[key] = value
	}
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
