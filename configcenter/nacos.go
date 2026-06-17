package configcenter

import (
	"context"
	"fmt"
	"strings"
	"time"

	nacosclients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/stellhub/stellar/config"
)

const (
	defaultNacosEndpoint = "127.0.0.1:8848"
	defaultNacosGroup    = "DEFAULT_GROUP"
)

type NacosAdapter struct {
	client config_client.IConfigClient
	cfg    *config.ConfigCenterConfig
}

func NewNacosAdapter(cfg *config.ConfigCenterConfig) (*NacosAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: config center config is required")
	}
	timeout, err := parseDuration(cfg.RequestTimeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid nacos request_timeout %q: %w", cfg.RequestTimeout, err)
	}
	clientConfig := &constant.ClientConfig{
		NamespaceId: nacosNamespace(cfg.Namespace),
		Username:    cfg.Username,
		Password:    cfg.Password,
		AppName:     cfg.AppID,
	}
	if timeout > 0 {
		clientConfig.TimeoutMs = uint64(timeout / time.Millisecond)
	}
	if cfg.Properties != nil {
		clientConfig.AccessKey = cfg.Properties["access_key"]
		clientConfig.SecretKey = cfg.Properties["secret_key"]
		clientConfig.LogDir = cfg.Properties["log_dir"]
		clientConfig.LogLevel = cfg.Properties["log_level"]
		clientConfig.CacheDir = cfg.Properties["cache_dir"]
	}
	client, err := nacosclients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: nacosServerConfigs(cfg),
	})
	if err != nil {
		return nil, err
	}
	return &NacosAdapter{client: client, cfg: cfg}, nil
}

func (a *NacosAdapter) Name() string {
	return AdapterNacos
}

func (a *NacosAdapter) Load(context.Context) ([]Source, error) {
	if a == nil || a.client == nil {
		return nil, fmt.Errorf("stellar: nacos config center is not configured")
	}
	sources := make([]Source, 0, len(sourceSpecs(a.cfg)))
	for _, spec := range sourceSpecs(a.cfg) {
		dataID := firstNonBlank(spec.dataID, spec.configKey, defaultConfigDataID)
		group := nacosGroup(spec.group)
		content, err := a.client.GetConfig(vo.ConfigParam{
			DataId: dataID,
			Group:  group,
		})
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(content) == "" {
			if spec.required {
				return nil, fmt.Errorf("stellar: nacos config %q in group %q is empty", dataID, group)
			}
			continue
		}
		sources = append(sources, Source{
			Key:      firstNonBlank(spec.configKey, dataID),
			DataID:   dataID,
			Group:    group,
			Format:   spec.format,
			Content:  content,
			Required: spec.required,
		})
	}
	return sources, nil
}

func (a *NacosAdapter) Close(context.Context) error {
	if a != nil && a.client != nil {
		a.client.CloseClient()
	}
	return nil
}

func nacosServerConfigs(cfg *config.ConfigCenterConfig) []constant.ServerConfig {
	endpoints := configCenterEndpoints(cfg, defaultNacosEndpoint)
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
	if namespace == "" || namespace == "default" || namespace == "public" {
		return ""
	}
	return namespace
}

func nacosGroup(group string) string {
	group = strings.TrimSpace(group)
	if group == "" || group == "default" {
		return defaultNacosGroup
	}
	return group
}
