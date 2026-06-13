package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stellhub/stellar/config"
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
	adapter Adapter
}

func New(adapter Adapter) (*Registry, error) {
	if adapter == nil {
		return nil, fmt.Errorf("stellar: registry adapter is required")
	}
	return &Registry{adapter: adapter}, nil
}

func NewFromConfig(ctx context.Context, cfg *config.RegistryConfig) (*Registry, error) {
	adapter, err := NewAdapterFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return New(adapter)
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
	if err := validateInstance(instance); err != nil {
		return err
	}
	if err := r.validate(); err != nil {
		return err
	}
	return r.adapter.Register(contextOrBackground(ctx), normalizeInstance(instance))
}

func (r *Registry) Deregister(ctx context.Context, instance Instance) error {
	if strings.TrimSpace(instance.Namespace) == "" || strings.TrimSpace(instance.Service) == "" || strings.TrimSpace(instance.InstanceID) == "" {
		return fmt.Errorf("stellar: namespace, service and instance_id are required")
	}
	if err := r.validate(); err != nil {
		return err
	}
	return r.adapter.Deregister(contextOrBackground(ctx), normalizeInstance(instance))
}

func (r *Registry) Discover(ctx context.Context, query Query) ([]Instance, error) {
	if strings.TrimSpace(query.Namespace) == "" || strings.TrimSpace(query.Service) == "" {
		return nil, fmt.Errorf("stellar: namespace and service are required")
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return r.adapter.Discover(contextOrBackground(ctx), normalizeQuery(query))
}

func (r *Registry) Watch(ctx context.Context, query Query) (Watcher, error) {
	if strings.TrimSpace(query.Namespace) == "" || strings.TrimSpace(query.Service) == "" {
		return nil, fmt.Errorf("stellar: namespace and service are required")
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return r.adapter.Watch(contextOrBackground(ctx), normalizeQuery(query))
}

func (r *Registry) Close(ctx context.Context) error {
	if r == nil || r.adapter == nil {
		return nil
	}
	return r.adapter.Close(contextOrBackground(ctx))
}

func (r *Registry) validate() error {
	if r == nil || r.adapter == nil {
		return fmt.Errorf("stellar: registry is not configured")
	}
	return nil
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
