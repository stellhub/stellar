package loadbalancer

import (
	"context"
	"sort"
	"strings"

	"github.com/stellhub/stellar/discovery"
	"github.com/stellhub/stellar/governance"
)

type RouteRule struct {
	Name     string
	Priority int
	Match    RouteMatch
	Filter   EndpointFilter
}

type RouteMatch struct {
	Protocol string
	Service  string
	Method   string
	Path     string
}

type EndpointFilter struct {
	Protocol     string
	EndpointName string
	Zone         string
	InstanceIDs  []string
	Labels       map[string]string
	Metadata     map[string]string
}

type StaticRouter struct {
	Rules []RouteRule
}

func (r StaticRouter) Route(_ context.Context, request Request, endpoints []discovery.Endpoint) ([]discovery.Endpoint, error) {
	if len(r.Rules) == 0 {
		return CloneEndpoints(endpoints), nil
	}
	rules := append([]RouteRule(nil), r.Rules...)
	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].Priority == rules[j].Priority {
			return rules[i].Name < rules[j].Name
		}
		return rules[i].Priority < rules[j].Priority
	})
	for _, rule := range rules {
		if !rule.Match.matches(request) {
			continue
		}
		filtered := filterEndpoints(endpoints, rule.Filter)
		if len(filtered) > 0 {
			return filtered, nil
		}
	}
	return CloneEndpoints(endpoints), nil
}

type GovernanceRouter struct {
	Store    *governance.Store
	Fallback Router
}

func NewGovernanceRouter(store *governance.Store, fallback Router) *GovernanceRouter {
	return &GovernanceRouter{Store: store, Fallback: fallback}
}

func (r *GovernanceRouter) Route(ctx context.Context, request Request, endpoints []discovery.Endpoint) ([]discovery.Endpoint, error) {
	if r == nil || r.Store == nil {
		return routeWithFallback(ctx, request, endpoints, nil)
	}
	rules := r.Store.Rules(governance.RuleKindRoute, func(rule governance.Rule) bool {
		return rule.Enabled && scopeMatches(rule.Scope, request)
	})
	for _, rule := range rules {
		filter := filterFromSpec(rule.Spec)
		filtered := filterEndpoints(endpoints, filter)
		if len(filtered) > 0 {
			return filtered, nil
		}
	}
	return routeWithFallback(ctx, request, endpoints, r.Fallback)
}

func routeWithFallback(ctx context.Context, request Request, endpoints []discovery.Endpoint, fallback Router) ([]discovery.Endpoint, error) {
	if fallback == nil {
		return CloneEndpoints(endpoints), nil
	}
	return fallback.Route(ctx, request, endpoints)
}

func (m RouteMatch) matches(request Request) bool {
	if m.Protocol != "" && !strings.EqualFold(m.Protocol, request.Protocol) {
		return false
	}
	if m.Service != "" && !strings.EqualFold(m.Service, request.Service) {
		return false
	}
	if m.Method != "" && !strings.EqualFold(m.Method, request.Method) {
		return false
	}
	if m.Path != "" && !strings.HasPrefix(request.Path, m.Path) {
		return false
	}
	return true
}

func scopeMatches(scope governance.Scope, request Request) bool {
	if scope.Transport != "" && !strings.EqualFold(scope.Transport, request.Protocol+".client") {
		return false
	}
	if scope.Service != "" && !strings.EqualFold(scope.Service, request.Service) {
		return false
	}
	if scope.Method != "" && !strings.EqualFold(scope.Method, request.Method) {
		return false
	}
	if scope.Path != "" && !strings.HasPrefix(request.Path, scope.Path) {
		return false
	}
	if scope.Target != "" && !strings.EqualFold(scope.Target, request.Target) {
		return false
	}
	return true
}

func filterEndpoints(endpoints []discovery.Endpoint, filter EndpointFilter) []discovery.Endpoint {
	values := make([]discovery.Endpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		endpoint = discovery.NormalizeEndpoint(endpoint)
		if !endpointMatches(endpoint, filter) {
			continue
		}
		values = append(values, endpoint)
	}
	return values
}

func endpointMatches(endpoint discovery.Endpoint, filter EndpointFilter) bool {
	if filter.Protocol != "" && !strings.EqualFold(filter.Protocol, endpoint.Protocol) {
		return false
	}
	if filter.EndpointName != "" && !strings.EqualFold(filter.EndpointName, endpoint.Name) {
		return false
	}
	if filter.Zone != "" && !strings.EqualFold(filter.Zone, endpoint.Zone) {
		return false
	}
	if len(filter.InstanceIDs) > 0 && !stringIn(endpoint.InstanceID, filter.InstanceIDs) {
		return false
	}
	if !mapContains(endpoint.Labels, filter.Labels) {
		return false
	}
	if !mapContains(endpoint.Metadata, filter.Metadata) {
		return false
	}
	return true
}

func filterFromSpec(spec map[string]any) EndpointFilter {
	return EndpointFilter{
		Protocol:     stringSpec(spec, "protocol"),
		EndpointName: stringSpec(spec, "endpoint_name"),
		Zone:         stringSpec(spec, "zone"),
		InstanceIDs:  stringSliceSpec(spec, "instance_ids"),
		Labels:       stringMapSpec(spec, "labels"),
		Metadata:     stringMapSpec(spec, "metadata"),
	}
}

func stringSpec(spec map[string]any, key string) string {
	if spec == nil {
		return ""
	}
	value, _ := spec[key].(string)
	return strings.TrimSpace(value)
}

func stringSliceSpec(spec map[string]any, key string) []string {
	if spec == nil {
		return nil
	}
	switch value := spec[key].(type) {
	case []string:
		return append([]string(nil), value...)
	case []any:
		result := make([]string, 0, len(value))
		for _, item := range value {
			if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
				result = append(result, strings.TrimSpace(text))
			}
		}
		return result
	case string:
		if strings.TrimSpace(value) == "" {
			return nil
		}
		return []string{strings.TrimSpace(value)}
	default:
		return nil
	}
}

func stringMapSpec(spec map[string]any, key string) map[string]string {
	if spec == nil {
		return nil
	}
	switch value := spec[key].(type) {
	case map[string]string:
		copied := make(map[string]string, len(value))
		for key, item := range value {
			copied[key] = item
		}
		return copied
	case map[string]any:
		result := make(map[string]string, len(value))
		for itemKey, itemValue := range value {
			if text, ok := itemValue.(string); ok {
				result[itemKey] = text
			}
		}
		return result
	default:
		return nil
	}
}

func mapContains(values map[string]string, required map[string]string) bool {
	for key, value := range required {
		if values == nil {
			return false
		}
		current, ok := values[key]
		if !ok {
			return false
		}
		if value != "" && current != value {
			return false
		}
	}
	return true
}

func stringIn(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if strings.EqualFold(value, strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}
