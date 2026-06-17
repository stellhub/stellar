package internal

import (
	"strings"

	"github.com/stellhub/stellar"
)

type sourceSummary struct {
	Key          string `json:"key"`
	DataID       string `json:"data_id"`
	Group        string `json:"group"`
	Format       string `json:"format"`
	Required     bool   `json:"required"`
	ContentBytes int    `json:"content_bytes"`
}

func bootstrapSummary(cfg *stellar.ConfigCenterConfig) map[string]any {
	if cfg == nil {
		return nil
	}
	return map[string]any{
		"adapter":    cfg.Adapter,
		"endpoint":   firstNonBlank(cfg.Endpoint, firstString(cfg.Endpoints)),
		"namespace":  cfg.Namespace,
		"group":      cfg.Group,
		"app_id":     cfg.AppID,
		"client_id":  cfg.ClientID,
		"env":        cfg.Env,
		"config_key": cfg.ConfigKey,
		"data_id":    cfg.DataID,
		"labels":     cfg.Labels,
	}
}

func httpServerSummary(cfg *stellar.HTTPServerConfig) map[string]any {
	if cfg == nil {
		return nil
	}
	return map[string]any{
		"enabled": cfg.Enabled == nil || *cfg.Enabled,
		"addr":    cfg.Addr,
		"adapter": cfg.Adapter,
	}
}

func mqSummary(cfg *stellar.MQConfig) map[string]any {
	if cfg == nil {
		return nil
	}
	return map[string]any{
		"adapter":       cfg.Adapter,
		"brokers":       cfg.Brokers,
		"client_id":     cfg.ClientID,
		"default_topic": cfg.Producer.DefaultTopic,
		"group_id":      cfg.Consumer.GroupID,
		"topics":        cfg.Consumer.Topics,
	}
}

func registrySummary(cfg *stellar.RegistryConfig) map[string]any {
	if cfg == nil {
		return nil
	}
	return map[string]any{
		"adapter":     cfg.Adapter,
		"endpoints":   cfg.Endpoints,
		"namespace":   cfg.Namespace,
		"service":     cfg.Service,
		"instance_id": cfg.InstanceID,
	}
}

func sourceSummaries(sources []stellar.ConfigCenterSource) []sourceSummary {
	if len(sources) == 0 {
		return nil
	}
	values := make([]sourceSummary, 0, len(sources))
	for _, source := range sources {
		values = append(values, sourceSummary{
			Key:          source.Key,
			DataID:       source.DataID,
			Group:        source.Group,
			Format:       source.Format,
			Required:     source.Required,
			ContentBytes: len(source.Content),
		})
	}
	return values
}

func firstString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
