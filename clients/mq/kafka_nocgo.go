//go:build !cgo

package mq

import (
	"context"
	"fmt"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

func newKafkaAdapter(context.Context, *config.MQConfig, *observability.Provider) (Adapter, error) {
	return nil, fmt.Errorf("stellar: kafka mq adapter requires CGO because confluent-kafka-go depends on librdkafka")
}
