package configcenter

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
	stellnula "github.com/stellhub/stellnula-go-sdk"
)

const defaultStellNulaEndpoint = "http://localhost:8060"

type StellNulaAdapter struct {
	client *stellnula.Client
	cfg    *config.ConfigCenterConfig
}

func NewStellNulaAdapter(cfg *config.ConfigCenterConfig, provider *observability.Provider) (*StellNulaAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: config center config is required")
	}
	options, err := stellNulaOptionsFromConfig(cfg, provider)
	if err != nil {
		return nil, err
	}
	client, err := stellnula.NewClient(options)
	if err != nil {
		return nil, err
	}
	return &StellNulaAdapter{client: client, cfg: cfg}, nil
}

func (a *StellNulaAdapter) Name() string {
	return AdapterStellNula
}

func (a *StellNulaAdapter) Load(ctx context.Context) ([]Source, error) {
	if a == nil || a.client == nil {
		return nil, fmt.Errorf("stellar: stellnula config center is not configured")
	}
	snapshot, err := a.client.SyncNow(contextOrBackground(ctx))
	if err != nil {
		return nil, err
	}

	sources := make([]Source, 0, len(sourceSpecs(a.cfg)))
	for _, spec := range sourceSpecs(a.cfg) {
		content, contentType, ok := stellNulaSnapshotValue(snapshot, spec.configKey)
		if !ok && spec.dataID != spec.configKey {
			content, contentType, ok = stellNulaSnapshotValue(snapshot, spec.dataID)
		}
		if !ok {
			if spec.required {
				return nil, fmt.Errorf("stellar: stellnula config %q is missing", spec.configKey)
			}
			continue
		}
		sources = append(sources, Source{
			Key:      spec.configKey,
			DataID:   spec.dataID,
			Group:    spec.group,
			Format:   firstNonBlank(spec.format, contentType),
			Content:  content,
			Required: spec.required,
		})
	}
	return sources, nil
}

func (a *StellNulaAdapter) Close(context.Context) error {
	if a == nil || a.client == nil {
		return nil
	}
	return a.client.Close()
}

func stellNulaOptionsFromConfig(cfg *config.ConfigCenterConfig, provider *observability.Provider) (stellnula.Options, error) {
	requestTimeout, err := parseDuration(cfg.RequestTimeout)
	if err != nil {
		return stellnula.Options{}, fmt.Errorf("stellar: invalid stellnula request_timeout %q: %w", cfg.RequestTimeout, err)
	}
	watchTimeout, err := parseDuration(cfg.WatchTimeout)
	if err != nil {
		return stellnula.Options{}, fmt.Errorf("stellar: invalid stellnula watch_timeout %q: %w", cfg.WatchTimeout, err)
	}
	retryDelay, err := parseDuration(cfg.RetryDelay)
	if err != nil {
		return stellnula.Options{}, fmt.Errorf("stellar: invalid stellnula retry_delay %q: %w", cfg.RetryDelay, err)
	}

	options := stellnula.Options{
		Endpoint:                 firstEndpoint(cfg, defaultStellNulaEndpoint),
		GRPCEndpoint:             cfg.GRPCEndpoint,
		APIToken:                 firstNonBlank(cfg.APIToken, cfg.Token),
		APIVersion:               cfg.APIVersion,
		SDKVersion:               cfg.SDKVersion,
		AppID:                    cfg.AppID,
		ClientID:                 cfg.ClientID,
		Env:                      cfg.Env,
		Region:                   cfg.Region,
		Zone:                     cfg.Zone,
		Cluster:                  cfg.Cluster,
		Namespace:                cfg.Namespace,
		Group:                    cfg.Group,
		ClientIP:                 cfg.ClientIP,
		HostName:                 cfg.HostName,
		Labels:                   cfg.Labels,
		SnapshotDirectory:        cfg.SnapshotDirectory,
		FailFastOnBootstrap:      cfg.FailFastOnBootstrap,
		RequestTimeout:           requestTimeout,
		WatchTimeout:             watchTimeout,
		RetryDelay:               retryDelay,
		AcceptLargeFileReference: true,
	}
	if cfg.GRPCPlaintext != nil {
		options.GRPCPlaintext = stellnula.Bool(*cfg.GRPCPlaintext)
	}
	if cfg.WatchEnabled != nil {
		options.WatchEnabled = stellnula.Bool(*cfg.WatchEnabled)
	}
	if requestTimeout > 0 {
		options.HTTPClient = &http.Client{Timeout: requestTimeout}
	}
	if provider != nil {
		options.Observability = stellnula.Observability{
			TracerProvider: provider.TracerProvider(),
			MeterProvider:  provider.MeterProvider(),
		}
	}
	return options, nil
}

func stellNulaSnapshotValue(snapshot stellnula.Snapshot, key string) (string, string, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", "", false
	}
	for _, entry := range snapshot.Entries {
		if entry.Deleted {
			continue
		}
		if entry.ConfigKey == key || entry.ConfigID == key {
			return entry.ConfigValue(), entry.ContentType, true
		}
	}
	if value, ok := snapshot.GetValue(key); ok {
		return value, "", true
	}
	return "", "", false
}
