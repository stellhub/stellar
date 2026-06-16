package internal

import (
	"errors"

	"github.com/stellhub/stellar"
)

type sendMessageRequest struct {
	Topic   string            `json:"topic"`
	Key     string            `json:"key"`
	Value   string            `json:"value"`
	Headers map[string]string `json:"headers"`
}

type sendMessageResponse struct {
	Adapter   string `json:"adapter"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
}

type receiveMessagesRequest struct {
	Topic  string   `json:"topic"`
	Topics []string `json:"topics"`
	Commit bool     `json:"commit"`
}

type receiveMessagesResponse struct {
	Adapter string          `json:"adapter"`
	Count   int             `json:"count"`
	Records []messageRecord `json:"records"`
}

type messageRecord struct {
	Topic     string            `json:"topic"`
	Partition int32             `json:"partition"`
	Offset    int64             `json:"offset"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	Headers   map[string]string `json:"headers,omitempty"`
}

func errMQNotConfigured() error {
	return errors.New("mq is not configured")
}

func headersFromRequest(headers map[string]string) []stellar.MessageQueueHeader {
	if len(headers) == 0 {
		return nil
	}
	values := make([]stellar.MessageQueueHeader, 0, len(headers))
	for key, value := range headers {
		if key == "" {
			continue
		}
		values = append(values, stellar.MessageQueueHeader{
			Key:   key,
			Value: []byte(value),
		})
	}
	return values
}

func responseRecords(records []stellar.MessageQueueMessage) []messageRecord {
	if len(records) == 0 {
		return nil
	}
	values := make([]messageRecord, 0, len(records))
	for _, record := range records {
		values = append(values, messageRecord{
			Topic:     record.Topic,
			Partition: record.Partition,
			Offset:    record.Offset,
			Key:       string(record.Key),
			Value:     string(record.Value),
			Headers:   headersToMap(record.Headers),
		})
	}
	return values
}

func headersToMap(headers []stellar.MessageQueueHeader) map[string]string {
	if len(headers) == 0 {
		return nil
	}
	values := make(map[string]string, len(headers))
	for _, header := range headers {
		if header.Key == "" {
			continue
		}
		values[header.Key] = string(header.Value)
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

func producerEnabled(cfg *stellar.MQConfig) bool {
	return cfg != nil && (cfg.Producer.Enabled == nil || *cfg.Producer.Enabled)
}

func consumerEnabled(cfg *stellar.MQConfig) bool {
	if cfg == nil {
		return false
	}
	if cfg.Consumer.Enabled != nil {
		return *cfg.Consumer.Enabled
	}
	return cfg.Consumer.GroupID != "" || len(cfg.Consumer.Topics) > 0
}

func consumerTopics(cfg *stellar.MQConfig) []string {
	if cfg == nil {
		return nil
	}
	if len(cfg.Consumer.Topics) > 0 {
		return append([]string(nil), cfg.Consumer.Topics...)
	}
	if cfg.Producer.DefaultTopic != "" {
		return []string{cfg.Producer.DefaultTopic}
	}
	return nil
}
