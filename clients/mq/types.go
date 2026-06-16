package mq

import (
	"context"
	"time"
)

const (
	DefaultName      = "mq"
	AdapterStellFlow = "stellflow"
	AdapterKafka     = "kafka"
)

type Header struct {
	Key   string
	Value []byte
}

type Message struct {
	Topic        string
	Partition    int32
	PartitionSet bool
	Offset       int64
	Key          []byte
	Value        []byte
	Headers      []Header
	Timestamp    time.Time
}

type Metadata struct {
	Topic     string
	Partition int32
	Offset    int64
}

type Producer interface {
	Send(context.Context, Message) (Metadata, error)
}

type Consumer interface {
	Subscribe(context.Context, []string) error
	Poll(context.Context) ([]Message, error)
	Commit(context.Context) error
}

type Adapter interface {
	Name() string
	Producer() Producer
	Consumer() Consumer
	Close(context.Context) error
}
