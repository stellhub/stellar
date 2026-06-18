# OpenTelemetry Custom Metrics Example

This example shows how application code can use the framework-managed OpenTelemetry provider to define custom business metrics and expose them through the same `/metrics` endpoint as Stellar's built-in metrics.

Run:

```bash
go run ./examples/opentelemetry/custom-metrics
```

Record one order:

```bash
curl -X POST http://localhost:18088/api/v1/otel/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"customer-1001","amount":128.5,"channel":"api"}'
```

Record several sample orders:

```bash
curl "http://localhost:18088/api/v1/otel/simulate?count=5"
```

Inspect the unified metrics output:

```bash
curl http://localhost:18088/metrics
```

Custom metric names:

```text
example.orders.created
example.orders.in_flight
example.order.amount
example.order.processing.duration
```

The important integration point is `app.Observability().Meter()`, which returns the same OpenTelemetry meter configured by `application.yml`.

