package registry

import (
	"context"
	"sync"

	"github.com/stellhub/stellar/observability"
)

type observableWatcher struct {
	ctx           context.Context
	base          Watcher
	observability *observability.Provider
	request       observability.RegistryRequest
	events        chan Event
	once          sync.Once
}

func newObservableWatcher(ctx context.Context, base Watcher, provider *observability.Provider, request observability.RegistryRequest) Watcher {
	if base == nil {
		return nil
	}
	if provider == nil {
		provider = observability.New()
	}
	return &observableWatcher{
		ctx:           contextOrBackground(ctx),
		base:          base,
		observability: provider,
		request:       request,
		events:        make(chan Event, 128),
	}
}

func (w *observableWatcher) Events() <-chan Event {
	w.once.Do(func() {
		go w.run()
	})
	return w.events
}

func (w *observableWatcher) Close() error {
	if w == nil || w.base == nil {
		return nil
	}
	return w.base.Close()
}

func (w *observableWatcher) run() {
	defer close(w.events)
	events := w.base.Events()
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			w.observability.RecordRegistryWatchEvent(w.ctx, w.request, event.Type, eventInstanceCount(event), eventEndpointCount(event))
			select {
			case w.events <- event:
			case <-w.ctx.Done():
				return
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func eventInstanceCount(event Event) int {
	if len(event.Instances) > 0 {
		return len(event.Instances)
	}
	if event.Instance != nil {
		return 1
	}
	return 0
}

func eventEndpointCount(event Event) int {
	if len(event.Instances) > 0 {
		return countEndpoints(event.Instances)
	}
	if event.Instance != nil {
		return len(event.Instance.Endpoints)
	}
	return 0
}

func countEndpoints(instances []Instance) int {
	total := 0
	for _, instance := range instances {
		total += len(instance.Endpoints)
	}
	return total
}
