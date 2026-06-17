package configcenter

import (
	"context"
	"errors"
)

const (
	DefaultName      = "configcenter"
	AdapterStellNula = "stellnula"
	AdapterNacos     = "nacos"
)

var ErrNotSupported = errors.New("stellar: config center operation is not supported")

type Source struct {
	Key      string
	DataID   string
	Group    string
	Format   string
	Content  string
	Required bool
}

type Adapter interface {
	Name() string
	Load(context.Context) ([]Source, error)
	Close(context.Context) error
}
