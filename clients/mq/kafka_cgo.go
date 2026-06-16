//go:build cgo

package mq

import (
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

type kafkaAdapter struct {
	producer *kafkaProducer
	consumer *kafkaConsumer
}

type kafkaProducer struct {
	client       *kafka.Producer
	flushTimeout time.Duration
}

type kafkaConsumer struct {
	client      *kafka.Consumer
	pollTimeout time.Duration
}

func newKafkaAdapter(ctx context.Context, cfg *config.MQConfig, _ *observability.Provider) (Adapter, error) {
	adapter := &kafkaAdapter{}
	if producerEnabled(cfg.Producer) {
		producer, err := newKafkaProducer(cfg)
		if err != nil {
			return nil, err
		}
		adapter.producer = producer
	}
	if consumerEnabled(cfg.Consumer) {
		consumer, err := newKafkaConsumer(cfg)
		if err != nil {
			_ = adapter.Close(ctx)
			return nil, err
		}
		adapter.consumer = consumer
		if len(cfg.Consumer.Topics) > 0 {
			if err := adapter.consumer.Subscribe(ctx, cfg.Consumer.Topics); err != nil {
				_ = adapter.Close(ctx)
				return nil, err
			}
		}
	}
	return adapter, nil
}

func (a *kafkaAdapter) Name() string {
	return AdapterKafka
}

func (a *kafkaAdapter) Producer() Producer {
	if a == nil {
		return nil
	}
	return a.producer
}

func (a *kafkaAdapter) Consumer() Consumer {
	if a == nil {
		return nil
	}
	return a.consumer
}

func (a *kafkaAdapter) Close(context.Context) error {
	if a == nil {
		return nil
	}
	var joined error
	if a.producer != nil {
		joined = stderrors.Join(joined, a.producer.Close())
	}
	if a.consumer != nil {
		joined = stderrors.Join(joined, a.consumer.Close())
	}
	return joined
}

func newKafkaProducer(cfg *config.MQConfig) (*kafkaProducer, error) {
	conf := kafkaConfig(cfg, cfg.Producer.Properties)
	if cfg.Producer.Acks != 0 {
		_ = conf.SetKey("acks", strconv.FormatInt(int64(cfg.Producer.Acks), 10))
	}
	if cfg.Producer.DeliveryTimeout != "" {
		duration, err := durationFromConfig(cfg.Producer.DeliveryTimeout, "producer.delivery_timeout", 0)
		if err != nil {
			return nil, err
		}
		_ = conf.SetKey("message.timeout.ms", strconv.FormatInt(duration.Milliseconds(), 10))
	}
	if cfg.Producer.RequestTimeout != "" {
		duration, err := durationFromConfig(cfg.Producer.RequestTimeout, "producer.request_timeout", 0)
		if err != nil {
			return nil, err
		}
		_ = conf.SetKey("request.timeout.ms", strconv.FormatInt(duration.Milliseconds(), 10))
	}
	if cfg.Producer.Idempotent {
		_ = conf.SetKey("enable.idempotence", true)
	}
	if strings.TrimSpace(cfg.Producer.TransactionalID) != "" {
		_ = conf.SetKey("transactional.id", strings.TrimSpace(cfg.Producer.TransactionalID))
	}
	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	flushTimeout, err := durationFromConfig(cfg.Producer.FlushTimeout, "producer.flush_timeout", 10*time.Second)
	if err != nil {
		producer.Close()
		return nil, err
	}
	return &kafkaProducer{client: producer, flushTimeout: flushTimeout}, nil
}

func newKafkaConsumer(cfg *config.MQConfig) (*kafkaConsumer, error) {
	if strings.TrimSpace(cfg.Consumer.GroupID) == "" {
		return nil, fmt.Errorf("stellar: kafka consumer group_id is required")
	}
	conf := kafkaConfig(cfg, cfg.Consumer.Properties)
	_ = conf.SetKey("group.id", strings.TrimSpace(cfg.Consumer.GroupID))
	_ = conf.SetKey("enable.auto.commit", cfg.Consumer.EnableAutoCommit)
	if strings.TrimSpace(cfg.Consumer.AutoOffsetReset) != "" {
		_ = conf.SetKey("auto.offset.reset", strings.ToLower(strings.TrimSpace(cfg.Consumer.AutoOffsetReset)))
	}
	if cfg.Consumer.SessionTimeout != "" {
		duration, err := durationFromConfig(cfg.Consumer.SessionTimeout, "consumer.session_timeout", 0)
		if err != nil {
			return nil, err
		}
		_ = conf.SetKey("session.timeout.ms", strconv.FormatInt(duration.Milliseconds(), 10))
	}
	if cfg.Consumer.AutoCommitInterval != "" {
		duration, err := durationFromConfig(cfg.Consumer.AutoCommitInterval, "consumer.auto_commit_interval", 0)
		if err != nil {
			return nil, err
		}
		_ = conf.SetKey("auto.commit.interval.ms", strconv.FormatInt(duration.Milliseconds(), 10))
	}
	consumer, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, err
	}
	pollTimeout, err := durationFromConfig(cfg.Consumer.PollTimeout, "consumer.poll_timeout", time.Second)
	if err != nil {
		_ = consumer.Close()
		return nil, err
	}
	return &kafkaConsumer{client: consumer, pollTimeout: pollTimeout}, nil
}

func (p *kafkaProducer) Send(ctx context.Context, message Message) (Metadata, error) {
	if p == nil || p.client == nil {
		return Metadata{}, fmt.Errorf("stellar: kafka producer is not configured")
	}
	partition := kafka.PartitionAny
	if message.PartitionSet {
		partition = message.Partition
	}
	delivery := make(chan kafka.Event, 1)
	err := p.client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &message.Topic,
			Partition: partition,
		},
		Key:       append([]byte(nil), message.Key...),
		Value:     append([]byte(nil), message.Value...),
		Headers:   kafkaHeaders(message.Headers),
		Timestamp: message.Timestamp,
	}, delivery)
	if err != nil {
		return Metadata{}, err
	}
	select {
	case event := <-delivery:
		delivered, ok := event.(*kafka.Message)
		if !ok {
			return Metadata{}, fmt.Errorf("stellar: unexpected kafka delivery event %T", event)
		}
		if delivered.TopicPartition.Error != nil {
			return Metadata{}, delivered.TopicPartition.Error
		}
		topic := message.Topic
		if delivered.TopicPartition.Topic != nil {
			topic = *delivered.TopicPartition.Topic
		}
		return Metadata{
			Topic:     topic,
			Partition: delivered.TopicPartition.Partition,
			Offset:    int64(delivered.TopicPartition.Offset),
		}, nil
	case <-contextOrBackground(ctx).Done():
		return Metadata{}, contextOrBackground(ctx).Err()
	}
}

func (p *kafkaProducer) Close() error {
	if p == nil || p.client == nil {
		return nil
	}
	timeoutMs := int(p.flushTimeout.Milliseconds())
	if timeoutMs <= 0 {
		timeoutMs = int((10 * time.Second).Milliseconds())
	}
	p.client.Flush(timeoutMs)
	p.client.Close()
	return nil
}

func (c *kafkaConsumer) Subscribe(_ context.Context, topics []string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("stellar: kafka consumer is not configured")
	}
	return c.client.SubscribeTopics(topics, nil)
}

func (c *kafkaConsumer) Poll(ctx context.Context) ([]Message, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("stellar: kafka consumer is not configured")
	}
	pollTimeout := timeoutFromContext(contextOrBackground(ctx), c.pollTimeout)
	event := c.client.Poll(int(pollTimeout.Milliseconds()))
	if event == nil {
		if err := contextOrBackground(ctx).Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	switch value := event.(type) {
	case *kafka.Message:
		if value.TopicPartition.Error != nil {
			return nil, value.TopicPartition.Error
		}
		return []Message{kafkaMessage(value)}, nil
	case kafka.Error:
		if value.IsTimeout() {
			return nil, nil
		}
		return nil, value
	default:
		return nil, nil
	}
}

func (c *kafkaConsumer) Commit(context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("stellar: kafka consumer is not configured")
	}
	_, err := c.client.Commit()
	return err
}

func (c *kafkaConsumer) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func kafkaConfig(cfg *config.MQConfig, properties map[string]string) *kafka.ConfigMap {
	conf := &kafka.ConfigMap{}
	for key, value := range mergeProperties(cfg.Properties, properties) {
		_ = conf.SetKey(key, value)
	}
	_ = conf.SetKey("bootstrap.servers", strings.Join(brokersFromConfig(cfg), ","))
	if strings.TrimSpace(cfg.ClientID) != "" {
		_ = conf.SetKey("client.id", strings.TrimSpace(cfg.ClientID))
	}
	return conf
}

func kafkaHeaders(headers []Header) []kafka.Header {
	if len(headers) == 0 {
		return nil
	}
	values := make([]kafka.Header, 0, len(headers))
	for _, header := range headers {
		if strings.TrimSpace(header.Key) == "" {
			continue
		}
		values = append(values, kafka.Header{
			Key:   header.Key,
			Value: append([]byte(nil), header.Value...),
		})
	}
	return values
}

func kafkaMessage(message *kafka.Message) Message {
	topic := ""
	if message.TopicPartition.Topic != nil {
		topic = *message.TopicPartition.Topic
	}
	return Message{
		Topic:        topic,
		Partition:    message.TopicPartition.Partition,
		PartitionSet: true,
		Offset:       int64(message.TopicPartition.Offset),
		Key:          append([]byte(nil), message.Key...),
		Value:        append([]byte(nil), message.Value...),
		Headers:      headersFromKafka(message.Headers),
		Timestamp:    message.Timestamp,
	}
}

func headersFromKafka(headers []kafka.Header) []Header {
	if len(headers) == 0 {
		return nil
	}
	values := make([]Header, 0, len(headers))
	for _, header := range headers {
		if strings.TrimSpace(header.Key) == "" {
			continue
		}
		values = append(values, Header{
			Key:   header.Key,
			Value: append([]byte(nil), header.Value...),
		})
	}
	return values
}

func timeoutFromContext(ctx context.Context, fallback time.Duration) time.Duration {
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			return 0
		}
		if fallback <= 0 || timeout < fallback {
			return timeout
		}
	}
	if fallback <= 0 {
		return time.Second
	}
	return fallback
}
