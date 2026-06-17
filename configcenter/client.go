package configcenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

type Client struct {
	adapter Adapter
}

func New(adapter Adapter) (*Client, error) {
	if adapter == nil {
		return nil, fmt.Errorf("stellar: config center adapter is required")
	}
	return &Client{adapter: adapter}, nil
}

func NewFromConfig(ctx context.Context, cfg *config.ConfigCenterConfig, provider *observability.Provider) (*Client, error) {
	adapter, err := NewAdapterFromConfig(ctx, cfg, provider)
	if err != nil {
		return nil, err
	}
	return New(adapter)
}

func NewAdapterFromConfig(ctx context.Context, cfg *config.ConfigCenterConfig, provider *observability.Provider) (Adapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: config center config is required")
	}
	adapter := strings.ToLower(strings.TrimSpace(cfg.Adapter))
	if adapter == "" {
		adapter = AdapterStellNula
	}
	switch adapter {
	case AdapterStellNula, "stell-nula", "stell_nula":
		return NewStellNulaAdapter(cfg, provider)
	case AdapterNacos:
		return NewNacosAdapter(cfg)
	default:
		return nil, fmt.Errorf("stellar: unsupported config center adapter %q", cfg.Adapter)
	}
}

func (c *Client) AdapterName() string {
	if c == nil || c.adapter == nil {
		return ""
	}
	return c.adapter.Name()
}

func (c *Client) Load(ctx context.Context) ([]Source, error) {
	if c == nil || c.adapter == nil {
		return nil, fmt.Errorf("stellar: config center is not configured")
	}
	return c.adapter.Load(contextOrBackground(ctx))
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.adapter == nil {
		return nil
	}
	return c.adapter.Close(contextOrBackground(ctx))
}

func Load(ctx context.Context, base config.Config, provider *observability.Provider) (config.Config, *Client, error) {
	cfg := base.Normalize()
	if cfg.ConfigCenter == nil {
		return cfg, nil, nil
	}
	if cfg.ConfigCenter.Enabled != nil && !*cfg.ConfigCenter.Enabled {
		return cfg, nil, nil
	}

	client, err := NewFromConfig(ctx, cfg.ConfigCenter, provider)
	if err != nil {
		return cfg, nil, err
	}
	sources, err := client.Load(ctx)
	if err != nil {
		_ = client.Close(ctx)
		return cfg, nil, err
	}

	merged := cfg
	for _, source := range sources {
		if strings.TrimSpace(source.Content) == "" {
			if source.Required {
				_ = client.Close(ctx)
				return cfg, nil, fmt.Errorf("stellar: config center source %q is empty", sourceName(source))
			}
			continue
		}
		if !sourceSupportsYAML(source) {
			_ = client.Close(ctx)
			return cfg, nil, fmt.Errorf("stellar: config center source %q only supports yaml content", sourceName(source))
		}
		remote, err := config.LoadBytes(sourceFileName(source), []byte(source.Content))
		if err != nil {
			_ = client.Close(ctx)
			return cfg, nil, fmt.Errorf("load config center source %q: %w", sourceName(source), err)
		}
		merged = config.Merge(merged, remote)
	}

	merged.ConfigCenter = cfg.ConfigCenter
	return merged.Normalize(), client, nil
}
