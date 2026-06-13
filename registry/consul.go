package registry

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/stellhub/stellar/config"
)

const defaultConsulEndpoint = "127.0.0.1:8500"

type ConsulAdapter struct {
	client            *api.Client
	namespace         string
	datacenter        string
	token             string
	heartbeatInterval time.Duration
	mu                sync.Mutex
	cancels           map[string]context.CancelFunc
}

func NewConsulAdapter(ctx context.Context, cfg *config.RegistryConfig) (*ConsulAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: registry config is required")
	}
	consulCfg := api.DefaultConfig()
	endpoint := valueOrDefault(firstEndpoint(cfg), defaultConsulEndpoint)
	address, scheme := consulAddress(endpoint, cfg.Scheme)
	consulCfg.Address = address
	if scheme != "" {
		consulCfg.Scheme = scheme
	}
	consulCfg.Token = cfg.Token
	consulCfg.Datacenter = cfg.Datacenter
	timeout, err := parseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid consul timeout %q: %w", cfg.Timeout, err)
	}
	if timeout > 0 {
		consulCfg.HttpClient.Timeout = timeout
	}
	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}
	heartbeatInterval, err := parseDuration(cfg.HeartbeatInterval)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid consul heartbeat_interval %q: %w", cfg.HeartbeatInterval, err)
	}
	return &ConsulAdapter{
		client:            client,
		namespace:         cfg.Namespace,
		datacenter:        cfg.Datacenter,
		token:             cfg.Token,
		heartbeatInterval: heartbeatInterval,
		cancels:           map[string]context.CancelFunc{},
	}, nil
}

func (a *ConsulAdapter) Name() string {
	return AdapterConsul
}

func (a *ConsulAdapter) Register(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	endpoint := instance.Endpoints[0]
	metadata := consulMetadata(instance, endpoint)
	registration := &api.AgentServiceRegistration{
		ID:      instance.InstanceID,
		Name:    instance.Service,
		Address: endpoint.Host,
		Port:    endpoint.Port,
		Tags:    consulTags(instance.Labels),
		Meta:    metadata,
		Weights: &api.AgentWeights{
			Passing: endpointWeight(endpoint),
			Warning: endpointWeight(endpoint),
		},
	}
	if namespace := consulNamespace(instance.Namespace); namespace != "" {
		registration.Namespace = namespace
	}
	if instance.TTLSeconds > 0 {
		ttl := time.Duration(instance.TTLSeconds) * time.Second
		registration.Check = &api.AgentServiceCheck{
			CheckID:                        consulCheckID(instance.InstanceID),
			Name:                           instance.Service + " heartbeat",
			TTL:                            ttl.String(),
			Status:                         api.HealthPassing,
			DeregisterCriticalServiceAfter: (ttl * 3).String(),
		}
	}
	opts := api.ServiceRegisterOpts{ReplaceExistingChecks: true, Token: a.token}.WithContext(contextOrBackground(ctx))
	if err := a.client.Agent().ServiceRegisterOpts(registration, opts); err != nil {
		return err
	}
	if instance.TTLSeconds > 0 {
		a.startHeartbeat(instance.InstanceID, instance.TTLSeconds)
	}
	return nil
}

func (a *ConsulAdapter) Deregister(ctx context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	a.stopHeartbeat(instance.InstanceID)
	options := &api.QueryOptions{
		Token: a.token,
	}
	if namespace := consulNamespace(instance.Namespace); namespace != "" {
		options.Namespace = namespace
	}
	return a.client.Agent().ServiceDeregisterOpts(instance.InstanceID, options.WithContext(contextOrBackground(ctx)))
}

func (a *ConsulAdapter) Discover(ctx context.Context, query Query) ([]Instance, error) {
	query = normalizeQuery(query)
	options := (&api.QueryOptions{
		Datacenter: a.datacenter,
		Token:      a.token,
	}).WithContext(contextOrBackground(ctx))
	if namespace := consulNamespace(query.Namespace); namespace != "" {
		options.Namespace = namespace
	}
	entries, _, err := a.client.Health().Service(query.Service, "", query.PassingOnly, options)
	if err != nil {
		return nil, err
	}
	instances := make([]Instance, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || entry.Service == nil {
			continue
		}
		instances = append(instances, instanceFromConsulService(query.Namespace, entry.Service))
	}
	return filterInstances(instances, query), nil
}

func (a *ConsulAdapter) Watch(ctx context.Context, query Query) (Watcher, error) {
	query = normalizeQuery(query)
	watcher, watchCtx := newChannelWatcher(ctx)
	go func() {
		defer watcher.closeEvents()
		var waitIndex uint64
		for {
			options := (&api.QueryOptions{
				Datacenter: a.datacenter,
				Token:      a.token,
				WaitIndex:  waitIndex,
				WaitTime:   5 * time.Minute,
			}).WithContext(watchCtx)
			if namespace := consulNamespace(query.Namespace); namespace != "" {
				options.Namespace = namespace
			}
			entries, meta, err := a.client.Health().Service(query.Service, "", query.PassingOnly, options)
			if err != nil {
				select {
				case <-time.After(time.Second):
					continue
				case <-watchCtx.Done():
					return
				}
			}
			if meta != nil && meta.LastIndex != 0 {
				waitIndex = meta.LastIndex
			}
			instances := make([]Instance, 0, len(entries))
			for _, entry := range entries {
				if entry == nil || entry.Service == nil {
					continue
				}
				instances = append(instances, instanceFromConsulService(query.Namespace, entry.Service))
			}
			if !watcher.publish(watchCtx, Event{Type: EventSnapshot, Instances: filterInstances(instances, query)}) {
				return
			}
		}
	}()
	return watcher, nil
}

func (a *ConsulAdapter) Close(context.Context) error {
	a.mu.Lock()
	for id, cancel := range a.cancels {
		cancel()
		delete(a.cancels, id)
	}
	a.mu.Unlock()
	return nil
}

func (a *ConsulAdapter) startHeartbeat(instanceID string, ttlSeconds int64) {
	a.stopHeartbeat(instanceID)
	interval := a.heartbeatInterval
	if interval <= 0 {
		interval = time.Duration(ttlSeconds) * time.Second / 3
	}
	if interval <= 0 {
		interval = 10 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.mu.Lock()
	a.cancels[instanceID] = cancel
	a.mu.Unlock()
	go func() {
		checkID := consulCheckID(instanceID)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		_ = a.client.Agent().UpdateTTL(checkID, "stellar heartbeat", api.HealthPassing)
		for {
			select {
			case <-ticker.C:
				_ = a.client.Agent().UpdateTTL(checkID, "stellar heartbeat", api.HealthPassing)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *ConsulAdapter) stopHeartbeat(instanceID string) {
	a.mu.Lock()
	cancel := a.cancels[instanceID]
	delete(a.cancels, instanceID)
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func consulAddress(endpoint string, configuredScheme string) (string, string) {
	endpoint = strings.TrimSpace(endpoint)
	scheme := strings.TrimSpace(configuredScheme)
	if parsed, err := url.Parse(endpoint); err == nil && parsed.Host != "" {
		if scheme == "" {
			scheme = parsed.Scheme
		}
		return parsed.Host, scheme
	}
	host, port, _ := hostPortFromEndpoint(endpoint, 8500)
	return net.JoinHostPort(host, fmt.Sprintf("%d", port)), scheme
}

func consulMetadata(instance Instance, endpoint Endpoint) map[string]string {
	return mergeStringMap(instance.Metadata, map[string]string{
		"stellar_namespace": instance.Namespace,
		"stellar_zone":      instance.Zone,
		"stellar_protocol":  endpoint.Protocol,
		"stellar_path":      endpoint.Path,
	})
}

func consulTags(labels map[string]string) []string {
	tags := make([]string, 0, len(labels))
	for key, value := range labels {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if strings.TrimSpace(value) == "" {
			tags = append(tags, strings.TrimSpace(key))
			continue
		}
		tags = append(tags, strings.TrimSpace(key)+"="+value)
	}
	return tags
}

func labelsFromConsulTags(tags []string) map[string]string {
	labels := map[string]string{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		key, value, hasValue := strings.Cut(tag, "=")
		if hasValue {
			labels[key] = value
			continue
		}
		labels[tag] = "true"
	}
	return labels
}

func instanceFromConsulService(namespace string, service *api.AgentService) Instance {
	metadata := cloneStringMap(service.Meta)
	if namespaceFromMeta := metadata["stellar_namespace"]; namespaceFromMeta != "" {
		namespace = namespaceFromMeta
	}
	zone := metadata["stellar_zone"]
	protocol := valueOrDefault(metadata["stellar_protocol"], "http")
	path := metadata["stellar_path"]
	delete(metadata, "stellar_namespace")
	delete(metadata, "stellar_zone")
	delete(metadata, "stellar_protocol")
	delete(metadata, "stellar_path")
	return normalizeInstance(Instance{
		Namespace:  namespace,
		Service:    service.Service,
		InstanceID: service.ID,
		Zone:       zone,
		Labels:     labelsFromConsulTags(service.Tags),
		Metadata:   metadata,
		Endpoints: []Endpoint{{
			Name:     protocol,
			Protocol: protocol,
			Host:     service.Address,
			Port:     service.Port,
			Path:     path,
			Weight:   service.Weights.Passing,
		}},
	})
}

func consulCheckID(instanceID string) string {
	return "service:" + instanceID
}

func consulNamespace(namespace string) string {
	namespace = strings.TrimSpace(namespace)
	if namespace == "" || namespace == "default" {
		return ""
	}
	return namespace
}
