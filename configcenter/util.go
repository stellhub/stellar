package configcenter

import (
	"context"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stellhub/stellar/config"
)

const defaultConfigDataID = "application.yaml"

type sourceSpec struct {
	dataID    string
	configKey string
	group     string
	format    string
	required  bool
}

func contextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func sourceSpecs(cfg *config.ConfigCenterConfig) []sourceSpec {
	if cfg == nil {
		return []sourceSpec{{dataID: defaultConfigDataID, configKey: defaultConfigDataID, required: true}}
	}
	if len(cfg.Sources) == 0 {
		key := firstNonBlank(cfg.ConfigKey, cfg.DataID, defaultConfigDataID)
		dataID := firstNonBlank(cfg.DataID, cfg.ConfigKey, defaultConfigDataID)
		return []sourceSpec{{
			dataID:    dataID,
			configKey: key,
			group:     cfg.Group,
			format:    cfg.Format,
			required:  true,
		}}
	}
	specs := make([]sourceSpec, 0, len(cfg.Sources))
	for _, item := range cfg.Sources {
		required := true
		if item.Required != nil {
			required = *item.Required
		}
		key := firstNonBlank(item.ConfigKey, item.DataID, cfg.ConfigKey, cfg.DataID, defaultConfigDataID)
		dataID := firstNonBlank(item.DataID, item.ConfigKey, cfg.DataID, cfg.ConfigKey, defaultConfigDataID)
		specs = append(specs, sourceSpec{
			dataID:    dataID,
			configKey: key,
			group:     firstNonBlank(item.Group, cfg.Group),
			format:    firstNonBlank(item.Format, cfg.Format),
			required:  required,
		})
	}
	return specs
}

func configCenterEndpoints(cfg *config.ConfigCenterConfig, fallback string) []string {
	if cfg == nil {
		return []string{fallback}
	}
	values := append([]string(nil), cfg.Endpoints...)
	if len(values) == 0 && strings.TrimSpace(cfg.Endpoint) != "" {
		values = append(values, cfg.Endpoint)
	}
	if len(values) == 0 && fallback != "" {
		values = append(values, fallback)
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func firstEndpoint(cfg *config.ConfigCenterConfig, fallback string) string {
	values := configCenterEndpoints(cfg, fallback)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func sourceName(source Source) string {
	return firstNonBlank(source.Key, source.DataID, defaultConfigDataID)
}

func sourceFileName(source Source) string {
	name := sourceName(source)
	if ext := strings.ToLower(filepath.Ext(name)); ext == ".yaml" || ext == ".yml" {
		return name
	}
	switch strings.ToLower(strings.TrimSpace(source.Format)) {
	case "yaml", "yml", "application/yaml", "application/x-yaml":
		return name + ".yaml"
	default:
		return defaultConfigDataID
	}
}

func sourceSupportsYAML(source Source) bool {
	ext := strings.ToLower(filepath.Ext(sourceName(source)))
	format := strings.ToLower(strings.TrimSpace(source.Format))
	switch format {
	case "yaml", "yml", "application/yaml", "application/x-yaml":
		return true
	case "":
		return ext == "" || ext == ".yaml" || ext == ".yml"
	default:
		return false
	}
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func parseDuration(value string) (time.Duration, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	return time.ParseDuration(value)
}

func hostPortFromEndpoint(endpoint string, defaultPort int) (string, int, string) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "127.0.0.1", defaultPort, ""
	}
	if parsed, err := url.Parse(endpoint); err == nil && parsed.Host != "" {
		host := parsed.Hostname()
		port := defaultPort
		if parsed.Port() != "" {
			if parsedPort, err := strconv.Atoi(parsed.Port()); err == nil {
				port = parsedPort
			}
		}
		return host, port, parsed.Scheme
	}
	host, portValue, err := net.SplitHostPort(endpoint)
	if err == nil {
		port, _ := strconv.Atoi(portValue)
		return host, port, ""
	}
	if strings.Contains(endpoint, ":") {
		parts := strings.Split(endpoint, ":")
		port, _ := strconv.Atoi(parts[len(parts)-1])
		return strings.Join(parts[:len(parts)-1], ":"), port, ""
	}
	return endpoint, defaultPort, ""
}
