package mq

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
	sfconsumer "github.com/stellhub/stellflow-go-sdk/consumer"
	sfsdkobservability "github.com/stellhub/stellflow-go-sdk/observability"
	sfproducer "github.com/stellhub/stellflow-go-sdk/producer"
	sfmessage "github.com/stellhub/stellflow-go-sdk/protocol/message"
	"github.com/stellhub/stellflow-go-sdk/stellflow"
)

type stellFlowAdapter struct {
	factory  *stellflow.ClientFactory
	producer *stellFlowProducer
	consumer *stellFlowConsumer
}

type stellFlowProducer struct {
	client *sfproducer.Client
}

type stellFlowConsumer struct {
	client      *sfconsumer.Client
	pollTimeout string
}

func newStellFlowAdapter(ctx context.Context, cfg *config.MQConfig, provider *observability.Provider) (Adapter, error) {
	options, err := stellFlowOptionsFromConfig(cfg, provider)
	if err != nil {
		return nil, err
	}
	factory, err := stellflow.NewClientFactory(options)
	if err != nil {
		return nil, err
	}
	adapter := &stellFlowAdapter{factory: factory}
	if producerEnabled(cfg.Producer) {
		adapter.producer = &stellFlowProducer{client: factory.NewProducer()}
	}
	if consumerEnabled(cfg.Consumer) {
		consumerClient := factory.NewConsumer()
		adapter.consumer = &stellFlowConsumer{
			client:      consumerClient,
			pollTimeout: cfg.Consumer.PollTimeout,
		}
		if len(cfg.Consumer.Topics) > 0 {
			if err := adapter.consumer.Subscribe(ctx, cfg.Consumer.Topics); err != nil {
				_ = adapter.Close(ctx)
				return nil, err
			}
		}
	}
	return adapter, nil
}

func (a *stellFlowAdapter) Name() string {
	return AdapterStellFlow
}

func (a *stellFlowAdapter) Producer() Producer {
	if a == nil {
		return nil
	}
	return a.producer
}

func (a *stellFlowAdapter) Consumer() Consumer {
	if a == nil {
		return nil
	}
	return a.consumer
}

func (a *stellFlowAdapter) Close(ctx context.Context) error {
	if a == nil {
		return nil
	}
	var joined error
	if a.producer != nil && a.producer.client != nil {
		joined = stderrors.Join(joined, a.producer.client.Close(contextOrBackground(ctx)))
	}
	if a.consumer != nil && a.consumer.client != nil {
		joined = stderrors.Join(joined, a.consumer.client.Close())
	}
	if a.factory != nil {
		joined = stderrors.Join(joined, a.factory.Close())
	}
	return joined
}

func (p *stellFlowProducer) Send(ctx context.Context, message Message) (Metadata, error) {
	if p == nil || p.client == nil {
		return Metadata{}, fmt.Errorf("stellar: stellflow producer is not configured")
	}
	partition := sfproducer.NoPartition
	if message.PartitionSet {
		partition = message.Partition
	}
	metadata, err := p.client.Send(contextOrBackground(ctx), sfproducer.Record{
		Topic:     message.Topic,
		Partition: partition,
		Key:       append([]byte(nil), message.Key...),
		Value:     append([]byte(nil), message.Value...),
		Headers:   stellFlowHeaders(message.Headers),
	})
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{
		Topic:     metadata.Topic,
		Partition: metadata.Partition,
		Offset:    metadata.Offset,
	}, nil
}

func (c *stellFlowConsumer) Subscribe(ctx context.Context, topics []string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("stellar: stellflow consumer is not configured")
	}
	return c.client.Subscribe(contextOrBackground(ctx), topics)
}

func (c *stellFlowConsumer) Poll(ctx context.Context) ([]Message, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("stellar: stellflow consumer is not configured")
	}
	pollCtx := contextOrBackground(ctx)
	timeout, err := durationFromConfig(c.pollTimeout, "consumer.poll_timeout", 0)
	if err != nil {
		return nil, err
	}
	if timeout > 0 {
		if _, ok := pollCtx.Deadline(); !ok {
			var cancel context.CancelFunc
			pollCtx, cancel = context.WithTimeout(pollCtx, timeout)
			defer cancel()
		}
	}
	records, err := c.client.Poll(pollCtx)
	if err != nil {
		return nil, err
	}
	messages := make([]Message, 0, len(records))
	for _, record := range records {
		messages = append(messages, Message{
			Topic:        record.Topic,
			Partition:    record.Partition,
			PartitionSet: true,
			Offset:       record.Offset,
			Key:          append([]byte(nil), record.Key...),
			Value:        append([]byte(nil), record.Value...),
			Headers:      headersFromStellFlow(record.Headers),
		})
	}
	return messages, nil
}

func (c *stellFlowConsumer) Commit(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("stellar: stellflow consumer is not configured")
	}
	return c.client.Commit(contextOrBackground(ctx))
}

func stellFlowOptionsFromConfig(cfg *config.MQConfig, provider *observability.Provider) (stellflow.Options, error) {
	requestTimeout, err := durationFromConfig(cfg.Producer.RequestTimeout, "producer.request_timeout", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	if requestTimeout == 0 {
		requestTimeout, err = durationFromConfig(cfg.Consumer.PollTimeout, "consumer.poll_timeout", 0)
		if err != nil {
			return stellflow.Options{}, err
		}
	}
	retryBackoff, err := durationFromConfig(cfg.Producer.RetryBackoff, "producer.retry_backoff", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	deliveryTimeout, err := durationFromConfig(cfg.Producer.DeliveryTimeout, "producer.delivery_timeout", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	linger, err := durationFromConfig(cfg.Producer.Linger, "producer.linger", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	sessionTimeout, err := durationFromConfig(cfg.Consumer.SessionTimeout, "consumer.session_timeout", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	heartbeatInterval, err := durationFromConfig(cfg.Consumer.HeartbeatInterval, "consumer.heartbeat_interval", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	maxPollInterval, err := durationFromConfig(cfg.Consumer.MaxPollInterval, "consumer.max_poll_interval", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	autoCommitInterval, err := durationFromConfig(cfg.Consumer.AutoCommitInterval, "consumer.auto_commit_interval", 0)
	if err != nil {
		return stellflow.Options{}, err
	}
	return stellflow.Options{
		BootstrapServers: brokersFromConfig(cfg),
		ClientID:         cfg.ClientID,
		RequestTimeout:   requestTimeout,
		Observability:    stellFlowObservability(provider),
		Producer: sfproducer.Options{
			Acks:                          cfg.Producer.Acks,
			TimeoutMs:                     cfg.Producer.TimeoutMs,
			DeliveryTimeout:               deliveryTimeout,
			RequestTimeout:                requestTimeout,
			RetryMaxAttempts:              cfg.Producer.RetryMaxAttempts,
			RetryBackoff:                  retryBackoff,
			MaxInFlight:                   cfg.Producer.MaxInFlight,
			Ordering:                      stellFlowOrdering(cfg.Producer.Ordering),
			BatchSize:                     cfg.Producer.BatchSize,
			BatchBytes:                    cfg.Producer.BatchBytes,
			Linger:                        linger,
			QueueSize:                     cfg.Producer.QueueSize,
			DisableAutoCreateTopics:       cfg.Producer.DisableAutoCreateTopics,
			AutoCreateTopicPartitionCount: cfg.Producer.AutoCreateTopicPartitionCount,
			Idempotent:                    cfg.Producer.Idempotent,
			TransactionalID:               cfg.Producer.TransactionalID,
		},
		Consumer: sfconsumer.Options{
			GroupID:            cfg.Consumer.GroupID,
			MemberID:           cfg.Consumer.MemberID,
			SessionTimeout:     sessionTimeout,
			HeartbeatInterval:  heartbeatInterval,
			MaxPollInterval:    maxPollInterval,
			FetchMaxBytes:      cfg.Consumer.FetchMaxBytes,
			PartitionMaxBytes:  cfg.Consumer.PartitionMaxBytes,
			CommitMetadata:     cfg.Consumer.CommitMetadata,
			OffsetReset:        stellFlowOffsetReset(cfg.Consumer.AutoOffsetReset),
			EnableAutoCommit:   cfg.Consumer.EnableAutoCommit,
			AutoCommitInterval: autoCommitInterval,
		},
	}, nil
}

func stellFlowObservability(provider *observability.Provider) sfsdkobservability.Options {
	if provider == nil {
		return sfsdkobservability.Options{}
	}
	if !provider.MessageProducerTraceEnabled() && !provider.MessageConsumerTraceEnabled() {
		return sfsdkobservability.Options{}
	}
	return sfsdkobservability.Options{TracerProvider: provider.TracerProvider()}
}

func stellFlowOrdering(value string) sfproducer.OrderingStrategy {
	if strings.EqualFold(strings.TrimSpace(value), "none") {
		return sfproducer.OrderingNone
	}
	return sfproducer.OrderingPerPartition
}

func stellFlowOffsetReset(value string) sfconsumer.OffsetResetStrategy {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "earliest":
		return sfconsumer.OffsetResetEarliest
	case "none":
		return sfconsumer.OffsetResetNone
	default:
		return sfconsumer.OffsetResetLatest
	}
}

func stellFlowHeaders(headers []Header) []sfmessage.RecordHeader {
	if len(headers) == 0 {
		return nil
	}
	values := make([]sfmessage.RecordHeader, 0, len(headers))
	for _, header := range headers {
		key := strings.TrimSpace(header.Key)
		if key == "" {
			continue
		}
		values = append(values, sfmessage.RecordHeader{
			Key:   &key,
			Value: append([]byte(nil), header.Value...),
		})
	}
	return values
}

func headersFromStellFlow(headers []sfmessage.RecordHeader) []Header {
	if len(headers) == 0 {
		return nil
	}
	values := make([]Header, 0, len(headers))
	for _, header := range headers {
		if header.Key == nil || strings.TrimSpace(*header.Key) == "" {
			continue
		}
		values = append(values, Header{
			Key:   *header.Key,
			Value: append([]byte(nil), header.Value...),
		})
	}
	return values
}
