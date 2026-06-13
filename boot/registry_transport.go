package boot

import (
	"context"
	"errors"

	serviceregistry "github.com/stellhub/stellar/registry"
)

type registryTransport struct {
	registry   *serviceregistry.Registry
	instance   serviceregistry.Instance
	registered bool
}

func (t *registryTransport) Name() string {
	if t == nil || t.registry == nil || t.registry.AdapterName() == "" {
		return "registry"
	}
	return "registry-" + t.registry.AdapterName()
}

func (t *registryTransport) Start(ctx context.Context) error {
	if t == nil || t.registry == nil {
		return nil
	}
	if err := t.registry.Register(ctx, t.instance); err != nil {
		return err
	}
	t.registered = true
	return nil
}

func (t *registryTransport) Stop(ctx context.Context) error {
	if t == nil || t.registry == nil {
		return nil
	}
	var joined error
	if t.registered {
		joined = errors.Join(joined, t.registry.Deregister(ctx, t.instance))
		t.registered = false
	}
	joined = errors.Join(joined, t.registry.Close(ctx))
	return joined
}
