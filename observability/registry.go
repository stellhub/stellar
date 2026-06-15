package observability

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	resultOK    = "ok"
	resultError = "error"
)

type RegistryRequest struct {
	Adapter    string
	Operation  string
	Namespace  string
	Service    string
	InstanceID string
}

type RegistryResult struct {
	Instances int
	Endpoints int
	Err       error
}

func (p *Provider) StartRegistry(ctx context.Context, request RegistryRequest) (context.Context, func(RegistryResult)) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	request = normalizeRegistryRequest(request)
	attrs := p.registryAttrs(request)
	start := time.Now()

	return ctx, func(result RegistryResult) {
		if !p.registryMetrics {
			return
		}
		resultAttrs := append(attrs, attribute.String("registry.result", resultName(result.Err)))
		if p.registryOperations != nil {
			p.registryOperations.Add(ctx, 1, metric.WithAttributes(resultAttrs...))
		}
		if p.registryDuration != nil {
			p.registryDuration.Record(ctx, durationSeconds(start), metric.WithAttributes(resultAttrs...))
		}
		p.RecordRegistryState(ctx, request, result.Instances, result.Endpoints)
	}
}

func (p *Provider) RecordRegistryWatchEvent(ctx context.Context, request RegistryRequest, eventType string, instances int, endpoints int) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !p.registryMetrics {
		return
	}
	request = normalizeRegistryRequest(request)
	attrs := append(p.registryAttrs(request),
		attribute.String("registry.event.type", valueOrUnknown(eventType)),
	)
	if p.registryWatchEvents != nil {
		p.registryWatchEvents.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	p.RecordRegistryState(ctx, request, instances, endpoints)
}

func (p *Provider) RecordRegistryState(ctx context.Context, request RegistryRequest, instances int, endpoints int) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !p.registryMetrics {
		return
	}
	attrs := p.registryStateAttrs(normalizeRegistryRequest(request))
	if p.registryInstances != nil && instances >= 0 {
		p.registryInstances.Record(ctx, int64(instances), metric.WithAttributes(attrs...))
	}
	if p.registryEndpoints != nil && endpoints >= 0 {
		p.registryEndpoints.Record(ctx, int64(endpoints), metric.WithAttributes(attrs...))
	}
}

func (p *Provider) initRegistryMetrics() {
	if p == nil || p.meter == nil {
		return
	}
	operations, err := p.meter.Int64Counter(
		"registry.client.operation.count",
		metric.WithDescription("Number of registry client operations."),
	)
	if err == nil {
		p.registryOperations = operations
	}
	duration, err := p.meter.Float64Histogram(
		"registry.client.operation.duration",
		metric.WithDescription("Duration of registry client operations."),
		metric.WithUnit("s"),
	)
	if err == nil {
		p.registryDuration = duration
	}
	instances, err := p.meter.Int64Gauge(
		"registry.client.instances",
		metric.WithDescription("Number of service instances observed by registry operations."),
		metric.WithUnit("{instance}"),
	)
	if err == nil {
		p.registryInstances = instances
	}
	endpoints, err := p.meter.Int64Gauge(
		"registry.client.endpoints",
		metric.WithDescription("Number of service endpoints observed by registry operations."),
		metric.WithUnit("{endpoint}"),
	)
	if err == nil {
		p.registryEndpoints = endpoints
	}
	events, err := p.meter.Int64Counter(
		"registry.client.watch.event.count",
		metric.WithDescription("Number of registry watch events."),
	)
	if err == nil {
		p.registryWatchEvents = events
	}
}

func (p *Provider) registryAttrs(request RegistryRequest) []attribute.KeyValue {
	attrs := append(p.commonAttrs(),
		attribute.String("registry.system.name", request.Adapter),
		attribute.String("registry.operation.name", request.Operation),
		attribute.String("registry.namespace", request.Namespace),
		attribute.String("registry.service.name", request.Service),
	)
	if request.InstanceID != "" {
		attrs = append(attrs, attribute.String("service.instance.id", request.InstanceID))
	}
	return attrs
}

func (p *Provider) registryStateAttrs(request RegistryRequest) []attribute.KeyValue {
	attrs := append(p.commonAttrs(),
		attribute.String("registry.system.name", request.Adapter),
		attribute.String("registry.namespace", request.Namespace),
		attribute.String("registry.service.name", request.Service),
	)
	return attrs
}

func normalizeRegistryRequest(request RegistryRequest) RegistryRequest {
	request.Adapter = valueOrUnknown(request.Adapter)
	request.Operation = valueOrUnknown(request.Operation)
	request.Namespace = valueOrUnknown(request.Namespace)
	request.Service = valueOrUnknown(request.Service)
	request.InstanceID = strings.TrimSpace(request.InstanceID)
	return request
}

func resultName(err error) string {
	if err != nil {
		return resultError
	}
	return resultOK
}

func valueOrUnknown(value string) string {
	if strings.TrimSpace(value) == "" {
		return "unknown"
	}
	return strings.TrimSpace(value)
}
