package observability

import (
	"context"
	"testing"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestRegistryMetricsAreRecorded(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	provider := New(WithMeterProvider(meterProvider))

	ctx, finish := provider.StartRegistry(context.Background(), RegistryRequest{
		Adapter:    "stellmap",
		Operation:  "discover",
		Namespace:  "default",
		Service:    "user-service",
		InstanceID: "user-service-1",
	})
	finish(RegistryResult{Instances: 2, Endpoints: 3})
	provider.RecordRegistryWatchEvent(ctx, RegistryRequest{
		Adapter:   "stellmap",
		Operation: "watch",
		Namespace: "default",
		Service:   "user-service",
	}, "snapshot", 2, 3)

	metrics := collectMetrics(t, reader)
	assertMetricExists(t, metrics, "registry.client.operation.count")
	assertMetricExists(t, metrics, "registry.client.operation.duration")
	assertMetricExists(t, metrics, "registry.client.instances")
	assertMetricExists(t, metrics, "registry.client.endpoints")
	assertMetricExists(t, metrics, "registry.client.watch.event.count")
}

func TestDiscoveryMetricsAreRecorded(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	provider := New(WithMeterProvider(meterProvider))

	ctx, finish := provider.StartDiscovery(context.Background(), DiscoveryRequest{
		Resolver:    "registry",
		Operation:   "resolve",
		Namespace:   "default",
		Service:     "user-service",
		Protocol:    "http",
		LoadBalance: "round_robin",
	})
	finish(DiscoveryResult{Result: DiscoveryResultRefresh, Endpoints: 2})
	provider.RecordDiscoveryWatchEvent(ctx, DiscoveryRequest{
		Resolver:  "registry",
		Operation: "watch",
		Namespace: "default",
		Service:   "user-service",
		Protocol:  "http",
	}, "snapshot", 2)

	metrics := collectMetrics(t, reader)
	assertMetricExists(t, metrics, "discovery.client.operation.count")
	assertMetricExists(t, metrics, "discovery.client.operation.duration")
	assertMetricExists(t, metrics, "discovery.client.endpoints")
	assertMetricExists(t, metrics, "discovery.client.cache.age")
	assertMetricExists(t, metrics, "discovery.client.watch.event.count")
}

func collectMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var metrics metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &metrics); err != nil {
		t.Fatalf("collect metrics: %v", err)
	}
	return metrics
}

func assertMetricExists(t *testing.T, metrics metricdata.ResourceMetrics, name string) {
	t.Helper()
	for _, scopeMetrics := range metrics.ScopeMetrics {
		for _, metric := range scopeMetrics.Metrics {
			if metric.Name == name {
				return
			}
		}
	}
	t.Fatalf("expected metric %q", name)
}
