package internal

import (
	"errors"

	"github.com/stellhub/stellar"
	"go.opentelemetry.io/otel/metric"
)

type orderMetrics struct {
	created  metric.Int64Counter
	inFlight metric.Int64UpDownCounter
	amount   metric.Float64Histogram
	duration metric.Float64Histogram
}

func newOrderMetrics(provider *stellar.ObservabilityProvider) (*orderMetrics, error) {
	if provider == nil {
		return nil, errors.New("observability provider is required")
	}
	meter := provider.Meter()

	created, err := meter.Int64Counter(
		"example.orders.created",
		metric.WithDescription("Number of example orders created by the custom metrics demo."),
		metric.WithUnit("{order}"),
	)
	if err != nil {
		return nil, err
	}
	inFlight, err := meter.Int64UpDownCounter(
		"example.orders.in_flight",
		metric.WithDescription("Number of example orders currently being processed."),
		metric.WithUnit("{order}"),
	)
	if err != nil {
		return nil, err
	}
	amount, err := meter.Float64Histogram(
		"example.order.amount",
		metric.WithDescription("Order amount observed by the custom metrics demo."),
		metric.WithUnit("{currency}"),
	)
	if err != nil {
		return nil, err
	}
	duration, err := meter.Float64Histogram(
		"example.order.processing.duration",
		metric.WithDescription("Duration of example order processing."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &orderMetrics{
		created:  created,
		inFlight: inFlight,
		amount:   amount,
		duration: duration,
	}, nil
}
