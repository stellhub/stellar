package grpcgoadapter

import (
	"context"
	"testing"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/discovery"
	"github.com/stellhub/stellar/interceptor"
	"github.com/stellhub/stellar/loadbalancer"
	grpcbalancer "google.golang.org/grpc/balancer"
)

func TestNewNamedClientConnFromConfig(t *testing.T) {
	cfg := &config.GRPCClientConfig{
		Timeout:  "3s",
		Insecure: boolPtr(true),
		Clients: map[string]config.GRPCNamedClientConfig{
			"user-service": {
				Target:  "dns:///localhost:19091",
				Timeout: "2s",
			},
		},
	}

	conn, target, err := NewNamedClientConnFromConfig(context.Background(), cfg, "user-service", nil)
	if err != nil {
		t.Fatalf("new grpc client conn: %v", err)
	}
	defer conn.Close()

	if target != "dns:///localhost:19091" {
		t.Fatalf("unexpected target %q", target)
	}
}

func TestNewNamedClientConnFromConfigRejectsDisabledClient(t *testing.T) {
	cfg := &config.GRPCClientConfig{
		Enabled: boolPtr(false),
		Clients: map[string]config.GRPCNamedClientConfig{
			"user-service": {Target: "dns:///localhost:19091"},
		},
	}

	if _, _, err := NewNamedClientConnFromConfig(context.Background(), cfg, "user-service", nil); err == nil {
		t.Fatalf("expected disabled grpc client error")
	}
}

func TestNewNamedClientConnFromConfigRequiresNamedClient(t *testing.T) {
	cfg := &config.GRPCClientConfig{
		Clients: map[string]config.GRPCNamedClientConfig{},
	}

	if _, _, err := NewNamedClientConnFromConfig(context.Background(), cfg, "user-service", nil); err == nil {
		t.Fatalf("expected missing grpc client error")
	}
}

func TestGRPCResolverAddressCarriesDiscoveryEndpoint(t *testing.T) {
	endpoint := discovery.Endpoint{
		Name:       "grpc",
		Protocol:   "grpc",
		Host:       "127.0.0.1",
		Port:       19091,
		InstanceID: "user-service-1",
	}
	address := grpcResolverAddress(endpoint, endpoint.Address())

	extracted, ok := endpointFromResolverAddress(address)
	if !ok {
		t.Fatalf("expected endpoint on resolver address")
	}
	if extracted.InstanceID != endpoint.InstanceID || extracted.Address() != endpoint.Address() {
		t.Fatalf("unexpected endpoint %#v", extracted)
	}
}

func TestLoadBalancerRequestUsesDiscoveryServiceForRouting(t *testing.T) {
	ctx := loadbalancer.ContextWithRequest(context.Background(), loadbalancer.Request{
		Service: "example.v1.UserAPI",
	})
	request := loadBalancerRequestFromPickInfo(ctx, grpcbalancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/example.v1.UserAPI/Get",
	}, "user-service")

	if request.Service != "user-service" {
		t.Fatalf("expected discovery service for routing, got %q", request.Service)
	}
	if request.Attributes["grpc.service"] != "example.v1.UserAPI" {
		t.Fatalf("expected grpc service attribute, got %#v", request.Attributes)
	}
}

func TestGRPCClientInvocationUsesNamedService(t *testing.T) {
	inv := grpcClientInvocation(context.Background(), "/example.v1.UserAPI/Get", nil, "user-service")

	if inv.Service != "user-service" {
		t.Fatalf("expected named client service, got %q", inv.Service)
	}
	if inv.Attributes["grpc.service"] != "example.v1.UserAPI" {
		t.Fatalf("expected protobuf service attribute, got %#v", inv.Attributes)
	}
}

func TestGRPCServerInvocationUsesConfiguredServiceName(t *testing.T) {
	inv := grpcInvocation(interceptor.KindGRPCServer, "/example.v1.UserAPI/Get", nil, nil, "order-service")

	if inv.Service != "order-service" {
		t.Fatalf("expected configured service name, got %q", inv.Service)
	}
	if inv.Attributes["grpc.service"] != "example.v1.UserAPI" {
		t.Fatalf("expected protobuf service attribute, got %#v", inv.Attributes)
	}
}

func boolPtr(value bool) *bool {
	return &value
}
