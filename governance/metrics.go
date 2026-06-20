package governance

import (
	"context"
	"strings"
	"time"

	"github.com/stellhub/stellar/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	ruleSync                   metric.Int64Counter
	routeMatches               metric.Int64Counter
	circuitBreakerStateChanges metric.Int64Counter
	circuitBreakerRejects      metric.Int64Counter
	rateLimitAllows            metric.Int64Counter
	rateLimitRejects           metric.Int64Counter
	rateLimitWaitDuration      metric.Float64Histogram
	authAllows                 metric.Int64Counter
	authDenies                 metric.Int64Counter
}

type MetricAttrs struct {
	Adapter   string
	RuleKind  string
	RuleID    string
	Transport string
	Service   string
	Method    string
	Resource  string
	Mode      string
	Behavior  string
	Outcome   string
}

func NewMetrics(provider *observability.Provider) *Metrics {
	if provider == nil {
		provider = observability.New()
	}
	meter := provider.Meter()
	metrics := &Metrics{}
	if value, err := meter.Int64Counter("governance.rule.sync.count", metric.WithDescription("Number of governance rule sync operations.")); err == nil {
		metrics.ruleSync = value
	}
	if value, err := meter.Int64Counter("governance.route.match.count", metric.WithDescription("Number of governance route rule matches.")); err == nil {
		metrics.routeMatches = value
	}
	if value, err := meter.Int64Counter("governance.circuitbreaker.state.change.count", metric.WithDescription("Number of circuit breaker state changes.")); err == nil {
		metrics.circuitBreakerStateChanges = value
	}
	if value, err := meter.Int64Counter("governance.circuitbreaker.reject.count", metric.WithDescription("Number of circuit breaker rejections.")); err == nil {
		metrics.circuitBreakerRejects = value
	}
	if value, err := meter.Int64Counter("governance.ratelimit.allow.count", metric.WithDescription("Number of rate limit allowed decisions.")); err == nil {
		metrics.rateLimitAllows = value
	}
	if value, err := meter.Int64Counter("governance.ratelimit.reject.count", metric.WithDescription("Number of rate limit rejected decisions.")); err == nil {
		metrics.rateLimitRejects = value
	}
	if value, err := meter.Float64Histogram("governance.ratelimit.wait.duration", metric.WithDescription("Duration spent waiting for blocking rate limit permits."), metric.WithUnit("s")); err == nil {
		metrics.rateLimitWaitDuration = value
	}
	if value, err := meter.Int64Counter("governance.auth.allow.count", metric.WithDescription("Number of auth allowed decisions.")); err == nil {
		metrics.authAllows = value
	}
	if value, err := meter.Int64Counter("governance.auth.deny.count", metric.WithDescription("Number of auth denied decisions.")); err == nil {
		metrics.authDenies = value
	}
	return metrics
}

func (m *Metrics) RecordRuleSync(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.ruleSync == nil {
		return
	}
	m.ruleSync.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordRouteMatch(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.routeMatches == nil {
		return
	}
	m.routeMatches.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordCircuitBreakerStateChange(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.circuitBreakerStateChanges == nil {
		return
	}
	m.circuitBreakerStateChanges.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordCircuitBreakerReject(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.circuitBreakerRejects == nil {
		return
	}
	m.circuitBreakerRejects.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordRateLimitAllow(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.rateLimitAllows == nil {
		return
	}
	m.rateLimitAllows.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordRateLimitReject(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.rateLimitRejects == nil {
		return
	}
	m.rateLimitRejects.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordRateLimitWait(ctx context.Context, duration time.Duration, attrs MetricAttrs) {
	if m == nil || m.rateLimitWaitDuration == nil {
		return
	}
	m.rateLimitWaitDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordAuthAllow(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.authAllows == nil {
		return
	}
	m.authAllows.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func (m *Metrics) RecordAuthDeny(ctx context.Context, attrs MetricAttrs) {
	if m == nil || m.authDenies == nil {
		return
	}
	m.authDenies.Add(ctx, 1, metric.WithAttributes(metricAttributes(attrs)...))
}

func metricAttributes(values MetricAttrs) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("adapter", valueOrUnknown(values.Adapter)),
		attribute.String("rule_kind", valueOrUnknown(values.RuleKind)),
		attribute.String("rule_id", valueOrUnknown(values.RuleID)),
		attribute.String("transport", valueOrUnknown(values.Transport)),
		attribute.String("service", valueOrUnknown(values.Service)),
		attribute.String("method", valueOrUnknown(values.Method)),
		attribute.String("resource", valueOrUnknown(values.Resource)),
		attribute.String("mode", valueOrUnknown(values.Mode)),
		attribute.String("behavior", valueOrUnknown(values.Behavior)),
		attribute.String("outcome", valueOrUnknown(values.Outcome)),
	}
}

func valueOrUnknown(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	return value
}
