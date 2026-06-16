package observability

import (
	"context"
	"testing"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func TestMessagingMetricsAreRecorded(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	provider := New(WithMeterProvider(meterProvider))

	ctx, finishSend := provider.StartMessagingProducer(context.Background(), MessagingClientRequest{
		System:          "stellflow",
		OperationName:   "send",
		OperationType:   MessagingOperationSend,
		DestinationName: "orders.created",
		ClientID:        "orders-service",
		ServerAddress:   "localhost",
		ServerPort:      9092,
	})
	finishSend(MessagingClientResult{Messages: 1, PartitionID: "0"})

	ctx, finishPoll := provider.StartMessagingConsumer(ctx, MessagingClientRequest{
		System:            "stellflow",
		OperationName:     "poll",
		OperationType:     MessagingOperationReceive,
		ConsumerGroupName: "orders-worker",
		ClientID:          "orders-service",
	})
	finishPoll(MessagingClientResult{Messages: 2})

	metrics := collectMetrics(t, reader)
	assertMetricExists(t, metrics, "messaging.client.operation.duration")
	assertMetricExists(t, metrics, "messaging.client.sent.messages")
	assertMetricExists(t, metrics, "messaging.client.consumed.messages")
}
