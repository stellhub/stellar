package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/stellhub/stellar/config"
)

func normalizeInstance(instance Instance) Instance {
	instance.Namespace = valueOrDefault(strings.TrimSpace(instance.Namespace), "default")
	instance.Service = strings.TrimSpace(instance.Service)
	instance.InstanceID = strings.TrimSpace(instance.InstanceID)
	instance.Zone = strings.TrimSpace(instance.Zone)
	instance.Labels = cloneStringMap(instance.Labels)
	instance.Metadata = cloneStringMap(instance.Metadata)
	instance.Endpoints = append([]Endpoint(nil), instance.Endpoints...)
	for i := range instance.Endpoints {
		instance.Endpoints[i] = normalizeEndpoint(instance.Endpoints[i])
	}
	return instance
}

func normalizeQuery(query Query) Query {
	query.Namespace = valueOrDefault(strings.TrimSpace(query.Namespace), "default")
	query.Service = strings.TrimSpace(query.Service)
	query.Zone = strings.TrimSpace(query.Zone)
	query.Labels = append([]string(nil), query.Labels...)
	return query
}

func normalizeEndpoint(endpoint Endpoint) Endpoint {
	endpoint.Name = strings.TrimSpace(endpoint.Name)
	endpoint.Protocol = valueOrDefault(strings.TrimSpace(endpoint.Protocol), "http")
	endpoint.Host = strings.TrimSpace(endpoint.Host)
	endpoint.Path = strings.TrimSpace(endpoint.Path)
	if endpoint.Weight <= 0 {
		endpoint.Weight = 100
	}
	return endpoint
}

func endpointsFromConfig(values []config.RegistryServiceEndpointConfig) []Endpoint {
	endpoints := make([]Endpoint, 0, len(values))
	for _, value := range values {
		endpoints = append(endpoints, normalizeEndpoint(Endpoint{
			Name:     value.Name,
			Protocol: value.Protocol,
			Host:     value.Host,
			Port:     value.Port,
			Path:     value.Path,
			Weight:   value.Weight,
		}))
	}
	return endpoints
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	copied := make(map[string]string, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func mergeStringMap(base map[string]string, extra map[string]string) map[string]string {
	merged := cloneStringMap(base)
	if merged == nil {
		merged = map[string]string{}
	}
	for key, value := range extra {
		if strings.TrimSpace(key) == "" {
			continue
		}
		merged[key] = value
	}
	return merged
}

func labelSelectorsMatch(labels map[string]string, selectors []string) bool {
	for _, selector := range selectors {
		selector = strings.TrimSpace(selector)
		if selector == "" {
			continue
		}
		key, value, hasValue := strings.Cut(selector, "=")
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		actual, ok := labels[key]
		if !ok {
			return false
		}
		if hasValue && actual != strings.TrimSpace(value) {
			return false
		}
	}
	return true
}

func filterInstances(instances []Instance, query Query) []Instance {
	query = normalizeQuery(query)
	filtered := make([]Instance, 0, len(instances))
	for _, instance := range instances {
		instance = normalizeInstance(instance)
		if query.Zone != "" && instance.Zone != query.Zone {
			continue
		}
		if len(query.Labels) > 0 && !labelSelectorsMatch(instance.Labels, query.Labels) {
			continue
		}
		filtered = append(filtered, instance)
	}
	return filtered
}

func endpointWeight(endpoint Endpoint) int {
	if endpoint.Weight <= 0 {
		return 100
	}
	return endpoint.Weight
}

func contextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func valueOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func durationSeconds(value string) int64 {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return int64(duration.Seconds())
}

func parseDuration(value string) (time.Duration, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func firstEndpoint(cfg *config.RegistryConfig) string {
	if cfg == nil {
		return ""
	}
	if strings.TrimSpace(cfg.Endpoint) != "" {
		return strings.TrimSpace(cfg.Endpoint)
	}
	if len(cfg.Endpoints) > 0 {
		return strings.TrimSpace(cfg.Endpoints[0])
	}
	return ""
}

func registryEndpoints(cfg *config.RegistryConfig, fallback string) []string {
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

func encodeInstance(instance Instance) ([]byte, error) {
	return json.Marshal(normalizeInstance(instance))
}

func decodeInstance(data []byte) (Instance, error) {
	var instance Instance
	if err := json.Unmarshal(data, &instance); err != nil {
		return Instance{}, err
	}
	return normalizeInstance(instance), nil
}

func instanceKey(prefix string, instance Instance) string {
	instance = normalizeInstance(instance)
	prefix = "/" + strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "/" {
		prefix = "/stellar/registry"
	}
	return fmt.Sprintf("%s/%s/%s/%s", prefix, instance.Namespace, instance.Service, instance.InstanceID)
}

func servicePrefix(prefix string, query Query) string {
	query = normalizeQuery(query)
	prefix = "/" + strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "/" {
		prefix = "/stellar/registry"
	}
	return fmt.Sprintf("%s/%s/%s/", prefix, query.Namespace, query.Service)
}
