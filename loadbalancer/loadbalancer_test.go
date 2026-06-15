package loadbalancer

import (
	"context"
	"testing"

	"github.com/stellhub/stellar/discovery"
	"github.com/stellhub/stellar/governance"
)

func TestStaticRouterFiltersEndpointsBeforeBalancing(t *testing.T) {
	router := StaticRouter{Rules: []RouteRule{{
		Name: "canary",
		Match: RouteMatch{
			Protocol: "http",
			Service:  "user-service",
		},
		Filter: EndpointFilter{
			Labels: map[string]string{"version": "v2"},
		},
	}}}

	endpoints, err := router.Route(context.Background(), Request{
		Protocol: "http",
		Service:  "user-service",
	}, []discovery.Endpoint{
		{Host: "127.0.0.1", Port: 8081, Labels: map[string]string{"version": "v1"}},
		{Host: "127.0.0.1", Port: 8082, Labels: map[string]string{"version": "v2"}},
	})
	if err != nil {
		t.Fatalf("route endpoints: %v", err)
	}
	if len(endpoints) != 1 || endpoints[0].Port != 8082 {
		t.Fatalf("unexpected routed endpoints %#v", endpoints)
	}
}

func TestGovernanceRouterUsesRouteRules(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{{
			ID:       "route-v2",
			Kind:     governance.RuleKindRoute,
			Enabled:  true,
			Priority: 10,
			Scope: governance.Scope{
				Transport: "http.client",
				Service:   "user-service",
			},
			Spec: map[string]any{
				"labels": map[string]any{"version": "v2"},
			},
		}},
	})
	router := NewGovernanceRouter(store, nil)

	endpoints, err := router.Route(context.Background(), Request{
		Protocol: "http",
		Service:  "user-service",
	}, []discovery.Endpoint{
		{Host: "127.0.0.1", Port: 8081, Labels: map[string]string{"version": "v1"}},
		{Host: "127.0.0.1", Port: 8082, Labels: map[string]string{"version": "v2"}},
	})
	if err != nil {
		t.Fatalf("route endpoints: %v", err)
	}
	if len(endpoints) != 1 || endpoints[0].Port != 8082 {
		t.Fatalf("unexpected governance routed endpoints %#v", endpoints)
	}
}

func TestLeastRequestAvoidsBusyEndpoint(t *testing.T) {
	picker := New(PolicyLeastRequest)
	endpoints := []discovery.Endpoint{
		{Host: "127.0.0.1", Port: 8081},
		{Host: "127.0.0.1", Port: 8082},
	}
	busy, err := picker.Pick(context.Background(), Request{}, endpoints)
	if err != nil {
		t.Fatalf("pick busy endpoint: %v", err)
	}

	selected, err := picker.Pick(context.Background(), Request{}, endpoints)
	if err != nil {
		t.Fatalf("pick endpoint: %v", err)
	}
	busy.Done(Result{})
	selected.Done(Result{})
	if EndpointKey(selected.Endpoint) == EndpointKey(busy.Endpoint) {
		t.Fatalf("expected least request to avoid busy endpoint")
	}
}

func TestConsistentHashKeepsSameKeyOnSameEndpoint(t *testing.T) {
	picker := New(PolicyConsistentHash)
	endpoints := []discovery.Endpoint{
		{InstanceID: "a", Host: "127.0.0.1", Port: 8081},
		{InstanceID: "b", Host: "127.0.0.1", Port: 8082},
		{InstanceID: "c", Host: "127.0.0.1", Port: 8083},
	}
	first, err := picker.Pick(context.Background(), Request{HashKey: "tenant-a"}, endpoints)
	if err != nil {
		t.Fatalf("pick endpoint: %v", err)
	}
	second, err := picker.Pick(context.Background(), Request{HashKey: "tenant-a"}, endpoints)
	if err != nil {
		t.Fatalf("pick endpoint: %v", err)
	}
	first.Done(Result{})
	second.Done(Result{})
	if EndpointKey(first.Endpoint) != EndpointKey(second.Endpoint) {
		t.Fatalf("expected same hash key to select same endpoint")
	}
}

func TestConsistentHashIgnoresEndpointOrder(t *testing.T) {
	picker := New(PolicyConsistentHash)
	endpoints := []discovery.Endpoint{
		{InstanceID: "a", Host: "127.0.0.1", Port: 8081},
		{InstanceID: "b", Host: "127.0.0.1", Port: 8082},
		{InstanceID: "c", Host: "127.0.0.1", Port: 8083},
	}
	reordered := []discovery.Endpoint{endpoints[2], endpoints[0], endpoints[1]}

	first, err := picker.Pick(context.Background(), Request{HashKey: "tenant-a"}, endpoints)
	if err != nil {
		t.Fatalf("pick endpoint: %v", err)
	}
	second, err := picker.Pick(context.Background(), Request{HashKey: "tenant-a"}, reordered)
	if err != nil {
		t.Fatalf("pick endpoint: %v", err)
	}
	first.Done(Result{})
	second.Done(Result{})
	if EndpointKey(first.Endpoint) != EndpointKey(second.Endpoint) {
		t.Fatalf("expected endpoint order not to change hash selection")
	}
}
