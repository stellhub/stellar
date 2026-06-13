package registry

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	nacosclients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/stellhub/stellar/config"
)

const (
	defaultNacosEndpoint = "127.0.0.1:8848"
	defaultNacosGroup    = "DEFAULT_GROUP"
	defaultNacosCluster  = "DEFAULT"
)

type NacosAdapter struct {
	client    naming_client.INamingClient
	group     string
	cluster   string
	ephemeral bool
}

func NewNacosAdapter(cfg *config.RegistryConfig) (*NacosAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: registry config is required")
	}
	timeout, err := parseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid nacos timeout %q: %w", cfg.Timeout, err)
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId: nacosNamespace(cfg.Namespace),
		Username:    cfg.Username,
		Password:    cfg.Password,
	}
	if timeout > 0 {
		clientConfig.TimeoutMs = uint64(timeout / time.Millisecond)
	}
	serverConfigs := nacosServerConfigs(cfg)
	client, err := nacosclients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return nil, err
	}
	ephemeral := true
	if cfg.Ephemeral != nil {
		ephemeral = *cfg.Ephemeral
	}
	return &NacosAdapter{
		client:    client,
		group:     valueOrDefault(cfg.Group, defaultNacosGroup),
		cluster:   valueOrDefault(cfg.Cluster, defaultNacosCluster),
		ephemeral: ephemeral,
	}, nil
}

func (a *NacosAdapter) Name() string {
	return AdapterNacos
}

func (a *NacosAdapter) Register(_ context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	for _, endpoint := range instance.Endpoints {
		ok, err := a.client.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          endpoint.Host,
			Port:        uint64(endpoint.Port),
			Weight:      float64(endpointWeight(endpoint)),
			Enable:      true,
			Healthy:     true,
			Metadata:    nacosMetadata(instance, endpoint),
			ClusterName: a.cluster,
			ServiceName: instance.Service,
			GroupName:   a.group,
			Ephemeral:   a.ephemeral,
		})
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("stellar: nacos register instance rejected")
		}
	}
	return nil
}

func (a *NacosAdapter) Deregister(_ context.Context, instance Instance) error {
	instance = normalizeInstance(instance)
	for _, endpoint := range instance.Endpoints {
		ok, err := a.client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          endpoint.Host,
			Port:        uint64(endpoint.Port),
			Cluster:     a.cluster,
			ServiceName: instance.Service,
			GroupName:   a.group,
			Ephemeral:   a.ephemeral,
		})
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("stellar: nacos deregister instance rejected")
		}
	}
	return nil
}

func (a *NacosAdapter) Discover(_ context.Context, query Query) ([]Instance, error) {
	query = normalizeQuery(query)
	clusters := []string{a.cluster}
	var values []model.Instance
	var err error
	if query.PassingOnly {
		values, err = a.client.SelectInstances(vo.SelectInstancesParam{
			Clusters:    clusters,
			ServiceName: query.Service,
			GroupName:   a.group,
			HealthyOnly: true,
		})
	} else {
		values, err = a.client.SelectAllInstances(vo.SelectAllInstancesParam{
			Clusters:    clusters,
			ServiceName: query.Service,
			GroupName:   a.group,
		})
	}
	if err != nil {
		return nil, err
	}
	return filterInstances(instancesFromNacos(query.Namespace, values), query), nil
}

func (a *NacosAdapter) Watch(ctx context.Context, query Query) (Watcher, error) {
	query = normalizeQuery(query)
	watcher, watchCtx := newChannelWatcher(ctx)
	if instances, err := a.Discover(watchCtx, query); err == nil {
		if !watcher.publish(watchCtx, Event{Type: EventSnapshot, Instances: instances}) {
			watcher.closeEvents()
			return watcher, nil
		}
	}
	param := &vo.SubscribeParam{
		ServiceName: query.Service,
		Clusters:    []string{a.cluster},
		GroupName:   a.group,
		SubscribeCallback: func(values []model.Instance, err error) {
			if err != nil {
				return
			}
			instances := filterInstances(instancesFromNacos(query.Namespace, values), query)
			_ = watcher.publish(watchCtx, Event{Type: EventSnapshot, Instances: instances})
		},
	}
	if err := a.client.Subscribe(param); err != nil {
		watcher.closeEvents()
		return nil, err
	}
	go func() {
		defer watcher.closeEvents()
		<-watchCtx.Done()
		_ = a.client.Unsubscribe(param)
	}()
	return watcher, nil
}

func (a *NacosAdapter) Close(context.Context) error {
	if a.client != nil {
		a.client.CloseClient()
	}
	return nil
}

func nacosServerConfigs(cfg *config.RegistryConfig) []constant.ServerConfig {
	endpoints := registryEndpoints(cfg, defaultNacosEndpoint)
	configs := make([]constant.ServerConfig, 0, len(endpoints))
	for _, endpoint := range endpoints {
		host, port, scheme := hostPortFromEndpoint(endpoint, 8848)
		if cfg.Scheme != "" {
			scheme = cfg.Scheme
		}
		options := []constant.ServerOption{}
		if scheme != "" {
			options = append(options, constant.WithScheme(scheme))
		}
		server := constant.NewServerConfig(host, uint64(port), options...)
		configs = append(configs, *server)
	}
	return configs
}

func nacosNamespace(namespace string) string {
	namespace = strings.TrimSpace(namespace)
	if namespace == "" || namespace == "default" {
		return ""
	}
	return namespace
}

func nacosMetadata(instance Instance, endpoint Endpoint) map[string]string {
	metadata := mergeStringMap(instance.Metadata, map[string]string{
		"stellar_namespace":     instance.Namespace,
		"stellar_zone":          instance.Zone,
		"stellar_instance_id":   instance.InstanceID,
		"stellar_endpoint_name": endpoint.Name,
		"stellar_protocol":      endpoint.Protocol,
		"stellar_path":          endpoint.Path,
		"stellar_weight":        strconv.Itoa(endpointWeight(endpoint)),
	})
	for key, value := range instance.Labels {
		metadata["stellar_label_"+key] = value
	}
	return metadata
}

func instancesFromNacos(namespace string, values []model.Instance) []Instance {
	grouped := map[string]Instance{}
	for _, value := range values {
		metadata := cloneStringMap(value.Metadata)
		instanceID := metadata["stellar_instance_id"]
		if instanceID == "" {
			instanceID = value.InstanceId
		}
		if instanceID == "" {
			instanceID = value.Ip + ":" + strconv.FormatUint(value.Port, 10)
		}
		instanceNamespace := valueOrDefault(metadata["stellar_namespace"], namespace)
		labels := map[string]string{}
		for key, item := range metadata {
			if strings.HasPrefix(key, "stellar_label_") {
				labels[strings.TrimPrefix(key, "stellar_label_")] = item
				delete(metadata, key)
			}
		}
		endpoint := normalizeEndpoint(Endpoint{
			Name:     metadata["stellar_endpoint_name"],
			Protocol: valueOrDefault(metadata["stellar_protocol"], "http"),
			Host:     value.Ip,
			Port:     int(value.Port),
			Path:     metadata["stellar_path"],
			Weight:   int(value.Weight),
		})
		zone := metadata["stellar_zone"]
		delete(metadata, "stellar_namespace")
		delete(metadata, "stellar_zone")
		delete(metadata, "stellar_instance_id")
		delete(metadata, "stellar_endpoint_name")
		delete(metadata, "stellar_protocol")
		delete(metadata, "stellar_path")
		delete(metadata, "stellar_weight")

		instance := grouped[instanceID]
		if instance.InstanceID == "" {
			instance = Instance{
				Namespace:  instanceNamespace,
				Service:    nacosServiceName(value.ServiceName),
				InstanceID: instanceID,
				Zone:       zone,
				Labels:     labels,
				Metadata:   metadata,
			}
		}
		instance.Endpoints = append(instance.Endpoints, endpoint)
		grouped[instanceID] = normalizeInstance(instance)
	}
	instances := make([]Instance, 0, len(grouped))
	for _, instance := range grouped {
		instances = append(instances, instance)
	}
	return instances
}

func nacosServiceName(serviceName string) string {
	serviceName = strings.TrimSpace(serviceName)
	if index := strings.LastIndex(serviceName, "@@"); index >= 0 && index < len(serviceName)-2 {
		return serviceName[index+2:]
	}
	return serviceName
}
