package internal

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/stellhub/stellar"
	stellarhttp "github.com/stellhub/stellar/transport/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type metricsStarter struct {
	metrics *orderMetrics
	seq     atomic.Int64
}

func NewMetricsStarter() *metricsStarter {
	return &metricsStarter{}
}

func (s *metricsStarter) Name() string {
	return "opentelemetry-custom-metrics-example"
}

func (s *metricsStarter) Condition(stellar.StarterContext) bool {
	return true
}

func (s *metricsStarter) Init(_ context.Context, app *stellar.App) error {
	metrics, err := newOrderMetrics(app.Observability())
	if err != nil {
		return err
	}
	s.metrics = metrics

	api := app.HTTP().Group("/api/v1")
	api.GET("/otel/status", s.handleStatus)
	api.GET("/otel/simulate", s.handleSimulate)
	stellarhttp.Handle(
		api,
		http.MethodPost,
		"/otel/orders",
		stellarhttp.JSONBinder[createOrderRequest](),
		s.handleCreateOrder,
		stellarhttp.JSONEncoder[createOrderResponse],
	)
	return nil
}

func (s *metricsStarter) Start(context.Context) error {
	return nil
}

func (s *metricsStarter) Stop(context.Context) error {
	return nil
}

func (s *metricsStarter) handleStatus(_ context.Context, _ *stellarhttp.Request) (*stellarhttp.Response, error) {
	return stellarhttp.JSON(http.StatusOK, map[string]any{
		"metrics_endpoint": "/metrics",
		"custom_metrics": []string{
			"example.orders.created",
			"example.orders.in_flight",
			"example.order.amount",
			"example.order.processing.duration",
		},
	}), nil
}

func (s *metricsStarter) handleSimulate(ctx context.Context, req *stellarhttp.Request) (*stellarhttp.Response, error) {
	count := 1
	if raw := req.Query.Get("count"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		if parsed > 0 {
			count = parsed
		}
	}
	if count > 100 {
		count = 100
	}

	for i := 0; i < count; i++ {
		_, err := s.recordOrder(ctx, createOrderRequest{
			CustomerID: "demo-customer",
			Amount:     float64(50 + i),
			Channel:    "simulate",
		})
		if err != nil {
			return nil, err
		}
	}
	return stellarhttp.JSON(http.StatusOK, map[string]any{
		"created": count,
		"message": "custom metrics recorded; visit /metrics to inspect unified output",
	}), nil
}

func (s *metricsStarter) handleCreateOrder(ctx context.Context, req *createOrderRequest) (*createOrderResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	return s.recordOrder(ctx, *req)
}

func (s *metricsStarter) recordOrder(ctx context.Context, req createOrderRequest) (*createOrderResponse, error) {
	if s.metrics == nil {
		return nil, errors.New("custom metrics are not initialized")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	channel := req.Channel
	if channel == "" {
		channel = "api"
	}

	start := time.Now()
	attrs := []attribute.KeyValue{
		attribute.String("example.channel", channel),
	}
	s.metrics.inFlight.Add(ctx, 1, metric.WithAttributes(attrs...))
	defer s.metrics.inFlight.Add(ctx, -1, metric.WithAttributes(attrs...))

	time.Sleep(10 * time.Millisecond)
	orderID := "order-" + strconv.FormatInt(s.seq.Add(1), 10)
	s.metrics.created.Add(ctx, 1, metric.WithAttributes(attrs...))
	s.metrics.amount.Record(ctx, req.Amount, metric.WithAttributes(attrs...))
	s.metrics.duration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(attrs...))

	return &createOrderResponse{
		OrderID: orderID,
		Status:  "created",
		Amount:  req.Amount,
		Channel: channel,
	}, nil
}
