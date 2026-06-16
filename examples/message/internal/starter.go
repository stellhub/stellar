package internal

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/stellhub/stellar"
	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type messageStarter struct {
	name          string
	cfg           stellar.Config
	mu            sync.Mutex
	subscribedKey string
}

func NewMessageStarter(name string) *messageStarter {
	if strings.TrimSpace(name) == "" {
		name = "message-quickstart"
	}
	return &messageStarter{name: name}
}

func (s *messageStarter) Name() string {
	return s.name
}

func (s *messageStarter) Condition(ctx stellar.StarterContext) bool {
	s.cfg = ctx.Config()
	return true
}

func (s *messageStarter) Init(_ context.Context, app *stellar.App) error {
	api := app.HTTP().Group("/api/v1")
	api.GET("/messages/status", s.handleStatus(app))
	stellarhttp.Handle(
		api,
		http.MethodPost,
		"/messages/send",
		stellarhttp.JSONBinder[sendMessageRequest](),
		s.handleSend(app),
		stellarhttp.JSONEncoder[sendMessageResponse],
	)
	stellarhttp.Handle(
		api,
		http.MethodPost,
		"/messages/receive",
		stellarhttp.JSONBinder[receiveMessagesRequest](),
		s.handleReceive(app),
		stellarhttp.JSONEncoder[receiveMessagesResponse],
	)
	return nil
}

func (s *messageStarter) Start(context.Context) error {
	return nil
}

func (s *messageStarter) Stop(context.Context) error {
	return nil
}

func (s *messageStarter) handleStatus(app *stellar.App) stellarhttp.Endpoint {
	return func(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
		mq, ok := app.MessageQueue()
		if !ok {
			return stellarhttp.JSON(http.StatusNotFound, map[string]any{
				"configured": false,
				"message":    "mq is not configured",
			}), nil
		}
		return stellarhttp.JSON(http.StatusOK, map[string]any{
			"configured": true,
			"adapter":    mq.AdapterName(),
			"producer":   producerEnabled(s.cfg.MQ),
			"consumer":   consumerEnabled(s.cfg.MQ),
			"topics":     consumerTopics(s.cfg.MQ),
		}), nil
	}
}

func (s *messageStarter) handleSend(app *stellar.App) func(context.Context, *sendMessageRequest) (*sendMessageResponse, error) {
	return func(ctx context.Context, req *sendMessageRequest) (*sendMessageResponse, error) {
		mq, ok := app.MessageQueue()
		if !ok {
			return nil, errMQNotConfigured()
		}
		topic := strings.TrimSpace(req.Topic)
		if topic == "" && s.cfg.MQ != nil {
			topic = strings.TrimSpace(s.cfg.MQ.Producer.DefaultTopic)
		}
		metadata, err := mq.Send(ctx, stellar.MessageQueueMessage{
			Topic:   topic,
			Key:     []byte(req.Key),
			Value:   []byte(req.Value),
			Headers: headersFromRequest(req.Headers),
		})
		if err != nil {
			return nil, err
		}
		return &sendMessageResponse{
			Adapter:   mq.AdapterName(),
			Topic:     metadata.Topic,
			Partition: metadata.Partition,
			Offset:    metadata.Offset,
		}, nil
	}
}

func (s *messageStarter) handleReceive(app *stellar.App) func(context.Context, *receiveMessagesRequest) (*receiveMessagesResponse, error) {
	return func(ctx context.Context, req *receiveMessagesRequest) (*receiveMessagesResponse, error) {
		mq, ok := app.MessageQueue()
		if !ok {
			return nil, errMQNotConfigured()
		}
		if err := s.ensureSubscribed(ctx, mq, req); err != nil {
			return nil, err
		}
		records, err := mq.Poll(ctx)
		if err != nil {
			return nil, err
		}
		if req.Commit && len(records) > 0 {
			if err := mq.Commit(ctx); err != nil {
				return nil, err
			}
		}
		return &receiveMessagesResponse{
			Adapter: mq.AdapterName(),
			Count:   len(records),
			Records: responseRecords(records),
		}, nil
	}
}

func (s *messageStarter) ensureSubscribed(ctx context.Context, mq *stellar.MessageQueue, req *receiveMessagesRequest) error {
	topics := s.receiveTopics(req)
	key := strings.Join(topics, ",")

	s.mu.Lock()
	if s.subscribedKey == key {
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()

	if err := mq.Subscribe(ctx, topics); err != nil {
		return err
	}

	s.mu.Lock()
	s.subscribedKey = key
	s.mu.Unlock()
	return nil
}

func (s *messageStarter) receiveTopics(req *receiveMessagesRequest) []string {
	topics := make([]string, 0, len(req.Topics)+1)
	if strings.TrimSpace(req.Topic) != "" {
		topics = append(topics, strings.TrimSpace(req.Topic))
	}
	for _, topic := range req.Topics {
		if strings.TrimSpace(topic) == "" {
			continue
		}
		topics = append(topics, strings.TrimSpace(topic))
	}
	if len(topics) > 0 {
		return topics
	}
	if s.cfg.MQ != nil && len(s.cfg.MQ.Consumer.Topics) > 0 {
		return append([]string(nil), s.cfg.MQ.Consumer.Topics...)
	}
	if s.cfg.MQ != nil && strings.TrimSpace(s.cfg.MQ.Producer.DefaultTopic) != "" {
		return []string{strings.TrimSpace(s.cfg.MQ.Producer.DefaultTopic)}
	}
	return []string{"stellar.quickstart.messages"}
}
