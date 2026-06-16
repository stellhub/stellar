package mq

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

type Client struct {
	adapter       Adapter
	provider      *observability.Provider
	clientID      string
	defaultTopic  string
	consumerGroup string
	serverAddress string
	serverPort    int
}

func NewFromConfig(ctx context.Context, cfg *config.MQConfig, provider *observability.Provider) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: mq config is required")
	}
	if provider == nil {
		provider = observability.New()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	adapter, err := newAdapter(ctx, cfg, provider)
	if err != nil {
		return nil, err
	}
	serverAddress, serverPort := firstBrokerAddress(cfg.Brokers)
	return &Client{
		adapter:       adapter,
		provider:      provider,
		clientID:      cfg.ClientID,
		defaultTopic:  cfg.Producer.DefaultTopic,
		consumerGroup: cfg.Consumer.GroupID,
		serverAddress: serverAddress,
		serverPort:    serverPort,
	}, nil
}

func (c *Client) AdapterName() string {
	if c == nil || c.adapter == nil {
		return ""
	}
	return c.adapter.Name()
}

func (c *Client) Send(ctx context.Context, message Message) (Metadata, error) {
	if strings.TrimSpace(message.Topic) == "" && c != nil {
		message.Topic = strings.TrimSpace(c.defaultTopic)
	}
	if strings.TrimSpace(message.Topic) == "" {
		return Metadata{}, fmt.Errorf("stellar: mq message topic is required")
	}
	if err := c.validateProducer(); err != nil {
		return Metadata{}, err
	}
	ctx = contextOrBackground(ctx)
	request := c.messagingRequest(
		"send",
		observability.MessagingOperationSend,
		message.Topic,
		partitionFromMessage(message),
	)
	ctx, finish := c.provider.StartMessagingProducer(ctx, request)
	metadata, err := c.adapter.Producer().Send(ctx, message)
	result := observability.MessagingClientResult{
		Messages:    1,
		PartitionID: observabilityPartitionID(metadata.Partition),
		Err:         err,
	}
	if err != nil {
		result.PartitionID = partitionFromMessage(message)
	}
	finish(result)
	return metadata, err
}

func (c *Client) Subscribe(ctx context.Context, topics []string) error {
	if len(topics) == 0 {
		return fmt.Errorf("stellar: mq subscribe topics are required")
	}
	if err := c.validateConsumer(); err != nil {
		return err
	}
	ctx = contextOrBackground(ctx)
	request := c.messagingRequest("subscribe", observability.MessagingOperationReceive, strings.Join(topics, ","), "")
	ctx, finish := c.provider.StartMessagingConsumer(ctx, request)
	err := c.adapter.Consumer().Subscribe(ctx, topics)
	finish(observability.MessagingClientResult{Err: err})
	return err
}

func (c *Client) Poll(ctx context.Context) ([]Message, error) {
	if err := c.validateConsumer(); err != nil {
		return nil, err
	}
	ctx = contextOrBackground(ctx)
	request := c.messagingRequest("poll", observability.MessagingOperationReceive, "", "")
	ctx, finish := c.provider.StartMessagingConsumer(ctx, request)
	messages, err := c.adapter.Consumer().Poll(ctx)
	finish(observability.MessagingClientResult{Messages: len(messages), Err: err})
	return messages, err
}

func (c *Client) Commit(ctx context.Context) error {
	if err := c.validateConsumer(); err != nil {
		return err
	}
	ctx = contextOrBackground(ctx)
	request := c.messagingRequest("commit", observability.MessagingOperationSettle, "", "")
	ctx, finish := c.provider.StartMessagingConsumer(ctx, request)
	err := c.adapter.Consumer().Commit(ctx)
	finish(observability.MessagingClientResult{Err: err})
	return err
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.adapter == nil {
		return nil
	}
	return c.adapter.Close(contextOrBackground(ctx))
}

func (c *Client) validateProducer() error {
	if c == nil || c.adapter == nil || c.adapter.Producer() == nil {
		return fmt.Errorf("stellar: mq producer is not configured")
	}
	if c.provider == nil {
		c.provider = observability.New()
	}
	return nil
}

func (c *Client) validateConsumer() error {
	if c == nil || c.adapter == nil || c.adapter.Consumer() == nil {
		return fmt.Errorf("stellar: mq consumer is not configured")
	}
	if c.provider == nil {
		c.provider = observability.New()
	}
	return nil
}

func (c *Client) messagingRequest(operationName string, operationType string, destination string, partitionID string) observability.MessagingClientRequest {
	system := ""
	if c != nil && c.adapter != nil {
		system = c.adapter.Name()
	}
	return observability.MessagingClientRequest{
		System:            system,
		OperationName:     operationName,
		OperationType:     operationType,
		DestinationName:   destination,
		PartitionID:       partitionID,
		ConsumerGroupName: c.consumerGroup,
		ClientID:          c.clientID,
		ServerAddress:     c.serverAddress,
		ServerPort:        c.serverPort,
	}
}

func newAdapter(ctx context.Context, cfg *config.MQConfig, provider *observability.Provider) (Adapter, error) {
	switch normalizeAdapter(cfg.Adapter) {
	case AdapterKafka:
		return newKafkaAdapter(ctx, cfg, provider)
	case AdapterStellFlow:
		return newStellFlowAdapter(ctx, cfg, provider)
	default:
		return nil, fmt.Errorf("stellar: unsupported mq adapter %q", cfg.Adapter)
	}
}

func normalizeAdapter(adapter string) string {
	if strings.TrimSpace(adapter) == "" {
		return AdapterStellFlow
	}
	switch strings.ToLower(strings.TrimSpace(adapter)) {
	case "stell-flow", "stell_flow":
		return AdapterStellFlow
	default:
		return strings.ToLower(strings.TrimSpace(adapter))
	}
}

func producerEnabled(cfg config.MQProducerConfig) bool {
	return cfg.Enabled == nil || *cfg.Enabled
}

func consumerEnabled(cfg config.MQConsumerConfig) bool {
	if cfg.Enabled != nil {
		return *cfg.Enabled
	}
	return strings.TrimSpace(cfg.GroupID) != "" || len(cfg.Topics) > 0
}

func contextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func partitionFromMessage(message Message) string {
	if !message.PartitionSet || message.Partition < 0 {
		return ""
	}
	return observabilityPartitionID(message.Partition)
}

func observabilityPartitionID(partition int32) string {
	if partition < 0 {
		return ""
	}
	return strconv.FormatInt(int64(partition), 10)
}

func firstBrokerAddress(brokers []string) (string, int) {
	if len(brokers) == 0 {
		return "", 0
	}
	value := strings.TrimSpace(brokers[0])
	if value == "" {
		return "", 0
	}
	if strings.Contains(value, "://") {
		if parsed, err := url.Parse(value); err == nil {
			value = parsed.Host
		}
	}
	host, portValue, err := net.SplitHostPort(value)
	if err != nil {
		return value, 0
	}
	port, _ := strconv.Atoi(portValue)
	return host, port
}
