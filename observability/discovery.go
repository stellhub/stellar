package observability

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	DiscoveryResultOK       = resultOK
	DiscoveryResultError    = resultError
	DiscoveryResultCacheHit = "cache_hit"
	DiscoveryResultRefresh  = "refresh"
	DiscoveryResultStale    = "stale"
)

type DiscoveryRequest struct {
	Resolver    string
	Operation   string
	Namespace   string
	Service     string
	Protocol    string
	LoadBalance string
}

type DiscoveryResult struct {
	Result    string
	Endpoints int
	CacheAge  time.Duration
	Stale     bool
	Err       error
}

func (p *Provider) StartDiscovery(ctx context.Context, request DiscoveryRequest) (context.Context, func(DiscoveryResult)) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	request = normalizeDiscoveryRequest(request)
	attrs := p.discoveryAttrs(request)
	start := time.Now()

	return ctx, func(result DiscoveryResult) {
		if !p.discoveryMetrics {
			return
		}
		resultName := normalizeDiscoveryResult(result)
		resultAttrs := append(attrs,
			attribute.String("discovery.result", resultName),
			attribute.Bool("discovery.cache.stale", result.Stale),
		)
		if p.discoveryOperations != nil {
			p.discoveryOperations.Add(ctx, 1, metric.WithAttributes(resultAttrs...))
		}
		if p.discoveryDuration != nil {
			p.discoveryDuration.Record(ctx, durationSeconds(start), metric.WithAttributes(resultAttrs...))
		}
		p.RecordDiscoveryState(ctx, request, result.Endpoints, result.CacheAge, result.Stale)
	}
}

func (p *Provider) RecordDiscoveryWatchEvent(ctx context.Context, request DiscoveryRequest, eventType string, endpoints int) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !p.discoveryMetrics {
		return
	}
	request = normalizeDiscoveryRequest(request)
	attrs := append(p.discoveryAttrs(request),
		attribute.String("discovery.event.type", valueOrUnknown(eventType)),
	)
	if p.discoveryWatchEvents != nil {
		p.discoveryWatchEvents.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	p.RecordDiscoveryState(ctx, request, endpoints, 0, false)
}

func (p *Provider) RecordDiscoveryState(ctx context.Context, request DiscoveryRequest, endpoints int, cacheAge time.Duration, stale bool) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !p.discoveryMetrics {
		return
	}
	attrs := append(p.discoveryStateAttrs(normalizeDiscoveryRequest(request)),
		attribute.Bool("discovery.cache.stale", stale),
	)
	if p.discoveryEndpoints != nil && endpoints >= 0 {
		p.discoveryEndpoints.Record(ctx, int64(endpoints), metric.WithAttributes(attrs...))
	}
	if p.discoveryCacheAge != nil && cacheAge >= 0 {
		p.discoveryCacheAge.Record(ctx, cacheAge.Seconds(), metric.WithAttributes(attrs...))
	}
}

func (p *Provider) initDiscoveryMetrics() {
	if p == nil || p.meter == nil {
		return
	}
	operations, err := p.meter.Int64Counter(
		"discovery.client.operation.count",
		metric.WithDescription("Number of client-side discovery operations."),
	)
	if err == nil {
		p.discoveryOperations = operations
	}
	duration, err := p.meter.Float64Histogram(
		"discovery.client.operation.duration",
		metric.WithDescription("Duration of client-side discovery operations."),
		metric.WithUnit("s"),
	)
	if err == nil {
		p.discoveryDuration = duration
	}
	endpoints, err := p.meter.Int64Gauge(
		"discovery.client.endpoints",
		metric.WithDescription("Number of endpoints currently visible to a discovery target."),
		metric.WithUnit("{endpoint}"),
	)
	if err == nil {
		p.discoveryEndpoints = endpoints
	}
	cacheAge, err := p.meter.Float64Gauge(
		"discovery.client.cache.age",
		metric.WithDescription("Age of the local discovery cache."),
		metric.WithUnit("s"),
	)
	if err == nil {
		p.discoveryCacheAge = cacheAge
	}
	events, err := p.meter.Int64Counter(
		"discovery.client.watch.event.count",
		metric.WithDescription("Number of client-side discovery watch events."),
	)
	if err == nil {
		p.discoveryWatchEvents = events
	}
}

func (p *Provider) discoveryAttrs(request DiscoveryRequest) []attribute.KeyValue {
	attrs := append(p.commonAttrs(),
		attribute.String("discovery.resolver.name", request.Resolver),
		attribute.String("discovery.operation.name", request.Operation),
		attribute.String("discovery.namespace", request.Namespace),
		attribute.String("discovery.service.name", request.Service),
		attribute.String("network.protocol.name", request.Protocol),
	)
	if request.LoadBalance != "" && request.LoadBalance != "unknown" {
		attrs = append(attrs, attribute.String("discovery.load_balance.policy", request.LoadBalance))
	}
	return attrs
}

func (p *Provider) discoveryStateAttrs(request DiscoveryRequest) []attribute.KeyValue {
	attrs := append(p.commonAttrs(),
		attribute.String("discovery.resolver.name", request.Resolver),
		attribute.String("discovery.namespace", request.Namespace),
		attribute.String("discovery.service.name", request.Service),
		attribute.String("network.protocol.name", request.Protocol),
	)
	if request.LoadBalance != "" && request.LoadBalance != "unknown" {
		attrs = append(attrs, attribute.String("discovery.load_balance.policy", request.LoadBalance))
	}
	return attrs
}

func normalizeDiscoveryRequest(request DiscoveryRequest) DiscoveryRequest {
	request.Resolver = valueOrUnknown(request.Resolver)
	request.Operation = valueOrUnknown(request.Operation)
	request.Namespace = valueOrUnknown(request.Namespace)
	request.Service = valueOrUnknown(request.Service)
	request.Protocol = strings.ToLower(valueOrUnknown(request.Protocol))
	request.LoadBalance = strings.ToLower(strings.TrimSpace(request.LoadBalance))
	return request
}

func normalizeDiscoveryResult(result DiscoveryResult) string {
	if result.Err != nil {
		return DiscoveryResultError
	}
	switch result.Result {
	case DiscoveryResultCacheHit, DiscoveryResultRefresh, DiscoveryResultStale:
		return result.Result
	default:
		return DiscoveryResultOK
	}
}
