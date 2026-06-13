package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	stellmapregistry "github.com/stellhub/stellmap-go-sdk/registry"

	"github.com/stellhub/stellar/config"
)

const defaultStellMapEndpoint = "http://localhost:8080"

type StellMapAdapter struct {
	client            *stellmapregistry.Client
	heartbeatInterval time.Duration
	mu                sync.Mutex
	registrars        map[string]*stellmapregistry.Registrar
}

func NewStellMapAdapter(cfg *config.RegistryConfig) (*StellMapAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: registry config is required")
	}
	timeout, err := parseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid stellmap timeout %q: %w", cfg.Timeout, err)
	}
	options := []stellmapregistry.ClientOption{
		stellmapregistry.WithUserAgent("stellar"),
	}
	if timeout > 0 {
		options = append(options, stellmapregistry.WithHTTPClient(&http.Client{Timeout: timeout}))
	}
	if cfg.Token != "" {
		token := strings.TrimSpace(cfg.Token)
		if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = "Bearer " + token
		}
		options = append(options, stellmapregistry.WithDefaultHeader("Authorization", token))
	}
	client, err := stellmapregistry.NewClient(valueOrDefault(firstEndpoint(cfg), defaultStellMapEndpoint), options...)
	if err != nil {
		return nil, err
	}
	heartbeatInterval, err := parseDuration(cfg.HeartbeatInterval)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid stellmap heartbeat_interval %q: %w", cfg.HeartbeatInterval, err)
	}
	return &StellMapAdapter{
		client:            client,
		heartbeatInterval: heartbeatInterval,
		registrars:        map[string]*stellmapregistry.Registrar{},
	}, nil
}

func (a *StellMapAdapter) Name() string {
	return AdapterStellMap
}

func (a *StellMapAdapter) Register(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	request := stellMapRegisterRequest(instance)
	options := []stellmapregistry.RegistrarOption{
		stellmapregistry.WithDeregisterOnStop(true),
	}
	if a.heartbeatInterval > 0 {
		options = append(options, stellmapregistry.WithHeartbeatInterval(a.heartbeatInterval))
	}
	registrar, err := stellmapregistry.NewRegistrar(a.client, request, options...)
	if err != nil {
		return err
	}
	key := registryKey(instance)
	a.stopRegistrar(contextOrBackground(ctx), key)
	if err := registrar.Start(contextOrBackground(ctx)); err != nil {
		return err
	}
	a.mu.Lock()
	a.registrars[key] = registrar
	a.mu.Unlock()
	return nil
}

func (a *StellMapAdapter) Deregister(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	key := registryKey(instance)
	if a.stopRegistrar(contextOrBackground(ctx), key) {
		return nil
	}
	return a.client.Deregister(contextOrBackground(ctx), stellmapregistry.DeregisterRequest{
		Namespace:  instance.Namespace,
		Service:    instance.Service,
		InstanceID: instance.InstanceID,
	})
}

func (a *StellMapAdapter) Discover(ctx context.Context, query Query) ([]Instance, error) {
	query = normalizeQuery(query)
	values, err := a.client.QueryInstances(contextOrBackground(ctx), stellMapQueryOptions(query))
	if err != nil {
		return nil, err
	}
	instances := make([]Instance, 0, len(values))
	for _, value := range values {
		instances = append(instances, instanceFromStellMap(value))
	}
	return filterInstances(instances, query), nil
}

func (a *StellMapAdapter) Watch(ctx context.Context, query Query) (Watcher, error) {
	query = normalizeQuery(query)
	watcher, watchCtx := newChannelWatcher(ctx)
	includeSnapshot := true
	sdkWatcher := a.client.WatchInstances(watchCtx, stellmapregistry.WatchOptions{
		QueryOptions:    stellMapQueryOptions(query),
		IncludeSnapshot: &includeSnapshot,
		AutoReconnect:   true,
	})
	go func() {
		defer watcher.closeEvents()
		defer func() { _ = sdkWatcher.Close() }()
		for {
			select {
			case event, ok := <-sdkWatcher.Events():
				if !ok {
					return
				}
				converted := eventFromStellMap(event)
				if converted.Type == "" {
					continue
				}
				if !watcher.publish(watchCtx, converted) {
					return
				}
			case <-sdkWatcher.Done():
				return
			case <-watchCtx.Done():
				return
			}
		}
	}()
	return watcher, nil
}

func (a *StellMapAdapter) Close(ctx context.Context) error {
	a.mu.Lock()
	registrars := make([]*stellmapregistry.Registrar, 0, len(a.registrars))
	for key, registrar := range a.registrars {
		registrars = append(registrars, registrar)
		delete(a.registrars, key)
	}
	a.mu.Unlock()
	var joined error
	for _, registrar := range registrars {
		if err := registrar.Stop(contextOrBackground(ctx)); err != nil {
			joined = joinErrors(joined, err)
		}
	}
	return joined
}

func (a *StellMapAdapter) stopRegistrar(ctx context.Context, key string) bool {
	a.mu.Lock()
	registrar := a.registrars[key]
	delete(a.registrars, key)
	a.mu.Unlock()
	if registrar == nil {
		return false
	}
	_ = registrar.Stop(ctx)
	return true
}

func stellMapRegisterRequest(instance Instance) stellmapregistry.RegisterRequest {
	endpoints := make([]stellmapregistry.Endpoint, 0, len(instance.Endpoints))
	for _, endpoint := range instance.Endpoints {
		endpoints = append(endpoints, stellmapregistry.Endpoint{
			Name:     endpoint.Name,
			Protocol: endpoint.Protocol,
			Host:     endpoint.Host,
			Port:     int32(endpoint.Port),
			Path:     endpoint.Path,
			Weight:   int32(endpointWeight(endpoint)),
		})
	}
	return stellmapregistry.RegisterRequest{
		Namespace:       instance.Namespace,
		Service:         instance.Service,
		InstanceID:      instance.InstanceID,
		Zone:            instance.Zone,
		Labels:          cloneStringMap(instance.Labels),
		Metadata:        cloneStringMap(instance.Metadata),
		Endpoints:       endpoints,
		LeaseTTLSeconds: instance.TTLSeconds,
	}
}

func stellMapQueryOptions(query Query) stellmapregistry.QueryOptions {
	return stellmapregistry.QueryOptions{
		Namespace: query.Namespace,
		Service:   query.Service,
		Zone:      query.Zone,
		Labels:    append([]string(nil), query.Labels...),
	}
}

func instanceFromStellMap(value stellmapregistry.Instance) Instance {
	endpoints := make([]Endpoint, 0, len(value.Endpoints))
	for _, endpoint := range value.Endpoints {
		endpoints = append(endpoints, Endpoint{
			Name:     endpoint.Name,
			Protocol: endpoint.Protocol,
			Host:     endpoint.Host,
			Port:     int(endpoint.Port),
			Path:     endpoint.Path,
			Weight:   int(endpoint.Weight),
		})
	}
	return normalizeInstance(Instance{
		Namespace:  value.Namespace,
		Service:    value.Service,
		InstanceID: value.InstanceID,
		Zone:       value.Zone,
		Labels:     cloneStringMap(value.Labels),
		Metadata:   cloneStringMap(value.Metadata),
		Endpoints:  endpoints,
		TTLSeconds: value.LeaseTTLSeconds,
	})
}

func eventFromStellMap(value stellmapregistry.WatchEvent) Event {
	event := Event{Type: string(value.Type)}
	if value.Instance != nil {
		instance := instanceFromStellMap(*value.Instance)
		event.Instance = &instance
	}
	if len(value.Instances) > 0 {
		event.Instances = make([]Instance, 0, len(value.Instances))
		for _, instance := range value.Instances {
			event.Instances = append(event.Instances, instanceFromStellMap(instance))
		}
	}
	return event
}

func registryKey(instance Instance) string {
	instance = normalizeInstance(instance)
	return instance.Namespace + "/" + instance.Service + "/" + instance.InstanceID
}

func joinErrors(left error, right error) error {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	return fmt.Errorf("%w; %w", left, right)
}
