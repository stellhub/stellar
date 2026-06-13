package registry

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/stellhub/stellar/config"
)

const defaultEtcdEndpoint = "localhost:2379"

type EtcdAdapter struct {
	client  *clientv3.Client
	prefix  string
	mu      sync.Mutex
	cancels map[string]context.CancelFunc
	leases  map[string]clientv3.LeaseID
}

func NewEtcdAdapter(ctx context.Context, cfg *config.RegistryConfig) (*EtcdAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: registry config is required")
	}
	timeout, err := parseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid etcd timeout %q: %w", cfg.Timeout, err)
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   registryEndpoints(cfg, defaultEtcdEndpoint),
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: timeout,
		Context:     contextOrBackground(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &EtcdAdapter{
		client:  client,
		prefix:  valueOrDefault(cfg.Prefix, "/stellar/registry"),
		cancels: map[string]context.CancelFunc{},
		leases:  map[string]clientv3.LeaseID{},
	}, nil
}

func (a *EtcdAdapter) Name() string {
	return AdapterEtcd
}

func (a *EtcdAdapter) Register(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	key := instanceKey(a.prefix, instance)
	value, err := encodeInstance(instance)
	if err != nil {
		return err
	}

	a.stopKeepAlive(key)
	if instance.TTLSeconds <= 0 {
		_, err := a.client.Put(contextOrBackground(ctx), key, string(value))
		return err
	}

	grant, err := a.client.Grant(contextOrBackground(ctx), instance.TTLSeconds)
	if err != nil {
		return err
	}
	if _, err := a.client.Put(contextOrBackground(ctx), key, string(value), clientv3.WithLease(grant.ID)); err != nil {
		_, _ = a.client.Revoke(contextOrBackground(ctx), grant.ID)
		return err
	}

	keepAliveCtx, cancel := context.WithCancel(context.Background())
	responses, err := a.client.KeepAlive(keepAliveCtx, grant.ID)
	if err != nil {
		cancel()
		_, _ = a.client.Revoke(contextOrBackground(ctx), grant.ID)
		return err
	}
	a.mu.Lock()
	a.cancels[key] = cancel
	a.leases[key] = grant.ID
	a.mu.Unlock()
	go func() {
		for range responses {
		}
	}()
	return nil
}

func (a *EtcdAdapter) Deregister(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	key := instanceKey(a.prefix, instance)
	a.stopKeepAlive(key)
	_, err := a.client.Delete(contextOrBackground(ctx), key)
	return err
}

func (a *EtcdAdapter) Discover(ctx context.Context, query Query) ([]Instance, error) {
	query = normalizeQuery(query)
	response, err := a.client.Get(contextOrBackground(ctx), servicePrefix(a.prefix, query), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	instances := make([]Instance, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		instance, err := decodeInstance(kv.Value)
		if err != nil {
			continue
		}
		instances = append(instances, instance)
	}
	return filterInstances(instances, query), nil
}

func (a *EtcdAdapter) Watch(ctx context.Context, query Query) (Watcher, error) {
	query = normalizeQuery(query)
	watcher, watchCtx := newChannelWatcher(ctx)
	go func() {
		defer watcher.closeEvents()

		if instances, err := a.Discover(watchCtx, query); err == nil {
			if !watcher.publish(watchCtx, Event{Type: EventSnapshot, Instances: instances}) {
				return
			}
		}

		responses := a.client.Watch(watchCtx, servicePrefix(a.prefix, query), clientv3.WithPrefix(), clientv3.WithPrevKV())
		for response := range responses {
			if response.Canceled {
				return
			}
			for _, event := range response.Events {
				switch event.Type {
				case mvccpb.PUT:
					instance, err := decodeInstance(event.Kv.Value)
					if err != nil || !instanceMatchesQuery(instance, query) {
						continue
					}
					if !watcher.publish(watchCtx, Event{Type: EventUpsert, Instance: &instance}) {
						return
					}
				case mvccpb.DELETE:
					instance := Instance{
						Namespace:  query.Namespace,
						Service:    query.Service,
						InstanceID: lastPathSegment(string(event.Kv.Key)),
					}
					if event.PrevKv != nil {
						if previous, err := decodeInstance(event.PrevKv.Value); err == nil {
							instance = previous
						}
					}
					if !watcher.publish(watchCtx, Event{Type: EventDelete, Instance: &instance}) {
						return
					}
				}
			}
		}
	}()
	return watcher, nil
}

func (a *EtcdAdapter) Close(context.Context) error {
	a.mu.Lock()
	for key, cancel := range a.cancels {
		cancel()
		delete(a.cancels, key)
	}
	a.leases = map[string]clientv3.LeaseID{}
	a.mu.Unlock()
	return a.client.Close()
}

func (a *EtcdAdapter) stopKeepAlive(key string) {
	a.mu.Lock()
	cancel := a.cancels[key]
	lease := a.leases[key]
	delete(a.cancels, key)
	delete(a.leases, key)
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if lease != 0 {
		_, _ = a.client.Revoke(context.Background(), lease)
	}
}

func instanceMatchesQuery(instance Instance, query Query) bool {
	return len(filterInstances([]Instance{instance}, query)) == 1
}

func lastPathSegment(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return ""
	}
	index := strings.LastIndex(path, "/")
	if index < 0 {
		return path
	}
	return path[index+1:]
}
