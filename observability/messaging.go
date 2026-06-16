package observability

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	MessagingOperationSend    = "send"
	MessagingOperationReceive = "receive"
	MessagingOperationSettle  = "settle"
)

type MessagingClientRequest struct {
	System            string
	OperationName     string
	OperationType     string
	DestinationName   string
	PartitionID       string
	ConsumerGroupName string
	ClientID          string
	ServerAddress     string
	ServerPort        int
}

type MessagingClientResult struct {
	Messages    int
	PartitionID string
	Err         error
}

func (p *Provider) StartMessagingProducer(ctx context.Context, request MessagingClientRequest) (context.Context, func(MessagingClientResult)) {
	return p.startMessagingClient(ctx, messagingProducer, request)
}

func (p *Provider) StartMessagingConsumer(ctx context.Context, request MessagingClientRequest) (context.Context, func(MessagingClientResult)) {
	return p.startMessagingClient(ctx, messagingConsumer, request)
}

func (p *Provider) initMessagingMetrics() {
	if p == nil || p.meter == nil {
		return
	}
	duration, err := p.meter.Float64Histogram(
		"messaging.client.operation.duration",
		metric.WithDescription("Duration of messaging client operations."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
	)
	if err == nil {
		p.messagingDuration = duration
	}
	sent, err := p.meter.Int64Counter(
		"messaging.client.sent.messages",
		metric.WithDescription("Number of messages attempted to be sent to the broker."),
		metric.WithUnit("{message}"),
	)
	if err == nil {
		p.messagingSentMessages = sent
	}
	read, err := p.meter.Int64Counter(
		"messaging.client.consumed.messages",
		metric.WithDescription("Number of messages received from the broker."),
		metric.WithUnit("{message}"),
	)
	if err == nil {
		p.messagingReadMessages = read
	}
}

type messagingClientRole string

const (
	messagingProducer messagingClientRole = "producer"
	messagingConsumer messagingClientRole = "consumer"
)

func (p *Provider) startMessagingClient(ctx context.Context, role messagingClientRole, request MessagingClientRequest) (context.Context, func(MessagingClientResult)) {
	if p == nil {
		p = New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	request = normalizeMessagingRequest(request)
	attrs := p.messagingAttrs(request)
	span := trace.SpanFromContext(ctx)
	if p.messagingTraceEnabled(role) {
		ctx, span = p.Tracer().Start(
			ctx,
			"messaging "+request.OperationName,
			trace.WithSpanKind(spanKindForMessaging(role)),
			trace.WithAttributes(attrs...),
		)
	}
	start := time.Now()

	return ctx, func(result MessagingClientResult) {
		duration := durationSeconds(start)
		resultAttrs := append([]attribute.KeyValue(nil), attrs...)
		if strings.TrimSpace(result.PartitionID) != "" {
			resultAttrs = append(resultAttrs, attribute.String("messaging.destination.partition.id", result.PartitionID))
		}
		if result.Err != nil {
			resultAttrs = append(resultAttrs, attribute.String("error.type", errorType(result.Err)))
		}

		if p.messagingTraceEnabled(role) {
			span.SetAttributes(resultAttrs...)
			if result.Err != nil {
				span.RecordError(result.Err)
				span.SetStatus(codes.Error, result.Err.Error())
			}
			span.End()
		}
		if p.messagingMetricsEnabled(role) {
			p.recordMessagingMetrics(ctx, role, request, result, resultAttrs, duration)
		}
	}
}

func (p *Provider) recordMessagingMetrics(
	ctx context.Context,
	role messagingClientRole,
	request MessagingClientRequest,
	result MessagingClientResult,
	attrs []attribute.KeyValue,
	duration float64,
) {
	if p.messagingDuration != nil {
		p.messagingDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
	}
	if role == messagingProducer && request.OperationType == MessagingOperationSend && p.messagingSentMessages != nil {
		count := result.Messages
		if count <= 0 && result.Err == nil {
			count = 1
		}
		if count > 0 {
			p.messagingSentMessages.Add(ctx, int64(count), metric.WithAttributes(attrs...))
		}
	}
	if role == messagingConsumer && request.OperationType == MessagingOperationReceive && p.messagingReadMessages != nil && result.Messages > 0 {
		p.messagingReadMessages.Add(ctx, int64(result.Messages), metric.WithAttributes(attrs...))
	}
}

func (p *Provider) messagingAttrs(request MessagingClientRequest) []attribute.KeyValue {
	attrs := append(p.commonAttrs(),
		attribute.String("messaging.system", request.System),
		attribute.String("messaging.operation.name", request.OperationName),
		attribute.String("messaging.operation.type", request.OperationType),
	)
	if request.DestinationName != "" {
		attrs = append(attrs, attribute.String("messaging.destination.name", request.DestinationName))
	}
	if request.PartitionID != "" {
		attrs = append(attrs, attribute.String("messaging.destination.partition.id", request.PartitionID))
	}
	if request.ConsumerGroupName != "" {
		attrs = append(attrs, attribute.String("messaging.consumer.group.name", request.ConsumerGroupName))
	}
	if request.ClientID != "" {
		attrs = append(attrs, attribute.String("messaging.client.id", request.ClientID))
	}
	if request.ServerAddress != "" {
		attrs = append(attrs, attribute.String("server.address", request.ServerAddress))
	}
	if request.ServerPort > 0 {
		attrs = append(attrs, attribute.Int("server.port", request.ServerPort))
	}
	return attrs
}

func (p *Provider) messagingTraceEnabled(role messagingClientRole) bool {
	if p == nil {
		return true
	}
	switch role {
	case messagingProducer:
		return p.messageProducerTrace
	case messagingConsumer:
		return p.messageConsumerTrace
	default:
		return true
	}
}

func (p *Provider) messagingMetricsEnabled(role messagingClientRole) bool {
	if p == nil {
		return true
	}
	switch role {
	case messagingProducer:
		return p.messageProducerMetrics
	case messagingConsumer:
		return p.messageConsumerMetrics
	default:
		return true
	}
}

func spanKindForMessaging(role messagingClientRole) trace.SpanKind {
	if role == messagingProducer {
		return trace.SpanKindProducer
	}
	return trace.SpanKindConsumer
}

func normalizeMessagingRequest(request MessagingClientRequest) MessagingClientRequest {
	request.System = valueOrUnknown(request.System)
	request.OperationName = valueOrUnknown(request.OperationName)
	request.OperationType = valueOrUnknown(request.OperationType)
	request.DestinationName = strings.TrimSpace(request.DestinationName)
	request.PartitionID = strings.TrimSpace(request.PartitionID)
	request.ConsumerGroupName = strings.TrimSpace(request.ConsumerGroupName)
	request.ClientID = strings.TrimSpace(request.ClientID)
	request.ServerAddress = strings.TrimSpace(request.ServerAddress)
	return request
}

func partitionID(partition int32) string {
	if partition < 0 {
		return ""
	}
	return strconv.FormatInt(int64(partition), 10)
}

func errorType(err error) string {
	if err == nil {
		return ""
	}
	errType := reflect.TypeOf(err)
	for errType.Kind() == reflect.Pointer {
		errType = errType.Elem()
	}
	if errType.Name() == "" {
		return "_OTHER"
	}
	if errType.PkgPath() == "" {
		return errType.Name()
	}
	return errType.PkgPath() + "." + errType.Name()
}
