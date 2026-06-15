package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

const (
	DefaultName     = "registry"
	AdapterEtcd     = "etcd"
	AdapterConsul   = "consul"
	AdapterNacos    = "nacos"
	AdapterStellMap = "stellmap"
	EventSnapshot   = "snapshot"
	EventUpsert     = "upsert"
	EventDelete     = "delete"
)

var ErrNotSupported = errors.New("stellar: registry operation is not supported")

type Adapter interface {
	Name() string
	Register(context.Context, Instance) error
	Deregister(context.Context, Instance) error
	Discover(context.Context, Query) ([]Instance, error)
	Watch(context.Context, Query) (Watcher, error)
	Close(context.Context) error
}

type Registry struct {
	adapter        Adapter
	observability  *observability.Provider
	metricsEnabled bool
}

type Option func(*Registry)

func WithObservability(provider *observability.Provider) Option {
	return func(registry *Registry) {
		if provider != nil {
			registry.observability = provider
		}
	}
}

func WithMetricsEnabled(enabled bool) Option {
	return func(registry *Registry) {
		registry.metricsEnabled = enabled
	}
}

func New(adapter Adapter, options ...Option) (*Registry, error) {
	if adapter == nil {
		return nil, fmt.Errorf("stellar: registry adapter is required")
	}
	registry := &Registry{
		adapter:        adapter,
		metricsEnabled: true,
	}
	for _, option := range options {
		option(registry)
	}
	if registry.observability == nil {
		registry.observability = observability.New()
	}
	return registry, nil
}

func NewFromConfig(ctx context.Context, cfg *config.RegistryConfig, options ...Option) (*Registry, error) {
	adapter, err := NewAdapterFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if cfg != nil && cfg.Observability.Metrics != nil {
		options = append(options, WithMetricsEnabled(*cfg.Observability.Metrics))
	}
	return New(adapter, options...)
}

func NewAdapterFromConfig(ctx context.Context, cfg *config.RegistryConfig) (Adapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: registry config is required")
	}
	adapter := strings.ToLower(strings.TrimSpace(cfg.Adapter))
	if adapter == "" {
		adapter = AdapterStellMap
	}
	switch adapter {
	case AdapterEtcd:
		return NewEtcdAdapter(ctx, cfg)
	case AdapterConsul:
		return NewConsulAdapter(ctx, cfg)
	case AdapterNacos:
		return NewNacosAdapter(cfg)
	case AdapterStellMap, "stell-map", "stell_map":
		return NewStellMapAdapter(cfg)
	default:
		return nil, fmt.Errorf("stellar: unsupported registry adapter %q", cfg.Adapter)
	}
}

func (r *Registry) AdapterName() string {
	if r == nil || r.adapter == nil {
		return ""
	}
	return r.adapter.Name()
}

func (r *Registry) Register(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	ctx, finish := r.startOperation(ctx, "register", instance.Namespace, instance.Service, instance.InstanceID)
	if err := validateInstance(instance); err != nil {
		finish(observability.RegistryResult{Endpoints: len(instance.Endpoints), Err: err})
		return err
	}
	if err := r.validate(); err != nil {
		finish(observability.RegistryResult{Endpoints: len(instance.Endpoints), Err: err})
		return err
	}
	err := r.adapter.Register(contextOrBackground(ctx), instance)
	finish(observability.RegistryResult{Instances: 1, Endpoints: len(instance.Endpoints), Err: err})
	return err
}

func (r *Registry) Deregister(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	ctx, finish := r.startOperation(ctx, "deregister", instance.Namespace, instance.Service, instance.InstanceID)
	if strings.TrimSpace(instance.Namespace) == "" || strings.TrimSpace(instance.Service) == "" || strings.TrimSpace(instance.InstanceID) == "" {
		err := fmt.Errorf("stellar: namespace, service and instance_id are required")
		finish(observability.RegistryResult{Endpoints: len(instance.Endpoints), Err: err})
		return err
	}
	if err := r.validate(); err != nil {
		finish(observability.RegistryResult{Endpoints: len(instance.Endpoints), Err: err})
		return err
	}
	err := r.adapter.Deregister(contextOrBackground(ctx), instance)
	finish(observability.RegistryResult{Instances: 1, Endpoints: len(instance.Endpoints), Err: err})
	return err
}

func (r *Registry) Discover(ctx context.Context, query Query) ([]Instance, error) {
	query = normalizeQuery(query)
	ctx, finish := r.startOperation(ctx, "discover", query.Namespace, query.Service, "")
	if strings.TrimSpace(query.Namespace) == "" || strings.TrimSpace(query.Service) == "" {
		err := fmt.Errorf("stellar: namespace and service are required")
		finish(observability.RegistryResult{Err: err})
		return nil, err
	}
	if err := r.validate(); err != nil {
		finish(observability.RegistryResult{Err: err})
		return nil, err
	}
	instances, err := r.adapter.Discover(contextOrBackground(ctx), query)
	finish(observability.RegistryResult{
		Instances: len(instances),
		Endpoints: countEndpoints(instances),
		Err:       err,
	})
	return instances, err
}

func (r *Registry) Watch(ctx context.Context, query Query) (Watcher, error) {
	query = normalizeQuery(query)
	ctx, finish := r.startOperation(ctx, "watch", query.Namespace, query.Service, "")
	if strings.TrimSpace(query.Namespace) == "" || strings.TrimSpace(query.Service) == "" {
		err := fmt.Errorf("stellar: namespace and service are required")
		finish(observability.RegistryResult{Err: err})
		return nil, err
	}
	if err := r.validate(); err != nil {
		finish(observability.RegistryResult{Err: err})
		return nil, err
	}
	watcher, err := r.adapter.Watch(contextOrBackground(ctx), query)
	finish(observability.RegistryResult{Err: err})
	if err != nil {
		return nil, err
	}
	if !r.metricsEnabled {
		return watcher, nil
	}
	return newObservableWatcher(ctx, watcher, r.observability, observability.RegistryRequest{
		Adapter:   r.AdapterName(),
		Operation: "watch",
		Namespace: query.Namespace,
		Service:   query.Service,
	}), nil
}

func (r *Registry) Close(ctx context.Context) error {
	if r == nil || r.adapter == nil {
		return nil
	}
	ctx, finish := r.startOperation(ctx, "close", "", "", "")
	err := r.adapter.Close(contextOrBackground(ctx))
	finish(observability.RegistryResult{Err: err})
	return err
}

func (r *Registry) validate() error {
	if r == nil || r.adapter == nil {
		return fmt.Errorf("stellar: registry is not configured")
	}
	return nil
}

func (r *Registry) startOperation(ctx context.Context, operation string, namespace string, service string, instanceID string) (context.Context, func(observability.RegistryResult)) {
	provider := (*observability.Provider)(nil)
	adapter := ""
	if r != nil {
		provider = r.observability
		adapter = r.AdapterName()
		if !r.metricsEnabled {
			if ctx == nil {
				ctx = context.Background()
			}
			return ctx, func(observability.RegistryResult) {}
		}
	}
	if provider == nil {
		provider = observability.New()
	}
	return provider.StartRegistry(ctx, observability.RegistryRequest{
		Adapter:    adapter,
		Operation:  operation,
		Namespace:  namespace,
		Service:    service,
		InstanceID: instanceID,
	})
}

type Endpoint struct {
	Name     string
	Protocol string
	Host     string
	Port     int
	Path     string
	Weight   int
}

type Instance struct {
	Namespace  string
	Service    string
	InstanceID string
	Zone       string
	Labels     map[string]string
	Metadata   map[string]string
	Endpoints  []Endpoint
	TTLSeconds int64
}

type Query struct {
	Namespace   string
	Service     string
	Zone        string
	Labels      []string
	PassingOnly bool
}

type Event struct {
	Type      string
	Instance  *Instance
	Instances []Instance
}

type Watcher interface {
	Events() <-chan Event
	Close() error
}

func InstanceFromConfig(cfg *config.RegistryConfig) (Instance, bool) {
	if cfg == nil {
		return Instance{}, false
	}
	instance := Instance{
		Namespace:  cfg.Namespace,
		Service:    cfg.Service,
		InstanceID: cfg.InstanceID,
		Zone:       cfg.Zone,
		Labels:     cloneStringMap(cfg.Labels),
		Metadata:   cloneStringMap(cfg.Metadata),
		Endpoints:  endpointsFromConfig(cfg.ServiceEndpoints),
		TTLSeconds: durationSeconds(cfg.TTL),
	}
	if strings.TrimSpace(instance.Namespace) == "" || strings.TrimSpace(instance.Service) == "" || strings.TrimSpace(instance.InstanceID) == "" || len(instance.Endpoints) == 0 {
		return Instance{}, false
	}
	return normalizeInstance(instance), true
}

func QueryFromConfig(cfg *config.RegistryConfig) (Query, bool) {
	if cfg == nil || strings.TrimSpace(cfg.Namespace) == "" || strings.TrimSpace(cfg.Service) == "" {
		return Query{}, false
	}
	return normalizeQuery(Query{
		Namespace:   cfg.Namespace,
		Service:     cfg.Service,
		Zone:        cfg.Zone,
		PassingOnly: true,
	}), true
}

func validateInstance(instance Instance) error {
	instance = normalizeInstance(instance)
	if instance.Namespace == "" || instance.Service == "" || instance.InstanceID == "" {
		return fmt.Errorf("stellar: namespace, service and instance_id are required")
	}
	if len(instance.Endpoints) == 0 {
		return fmt.Errorf("stellar: at least one service endpoint is required")
	}
	for _, endpoint := range instance.Endpoints {
		if endpoint.Host == "" || endpoint.Port <= 0 {
			return fmt.Errorf("stellar: endpoint host and port are required")
		}
	}
	return nil
}
