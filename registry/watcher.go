package registry

import (
	"context"
	"sync"
)

type channelWatcher struct {
	cancel     context.CancelFunc
	events     chan Event
	cancelOnce sync.Once
	closeOnce  sync.Once
}

func newChannelWatcher(parent context.Context) (*channelWatcher, context.Context) {
	ctx, cancel := context.WithCancel(contextOrBackground(parent))
	return &channelWatcher{
		cancel: cancel,
		events: make(chan Event, 128),
	}, ctx
}

func (w *channelWatcher) Events() <-chan Event {
	return w.events
}

func (w *channelWatcher) Close() error {
	w.cancelOnce.Do(func() {
		if w.cancel != nil {
			w.cancel()
		}
	})
	return nil
}

func (w *channelWatcher) closeEvents() {
	w.cancelOnce.Do(func() {
		if w.cancel != nil {
			w.cancel()
		}
	})
	w.closeOnce.Do(func() {
		close(w.events)
	})
}

func (w *channelWatcher) publish(ctx context.Context, event Event) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	select {
	case w.events <- event:
		return true
	case <-ctx.Done():
		return false
	}
}
