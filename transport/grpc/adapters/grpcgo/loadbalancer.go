package grpcgoadapter

import (
	"context"
	"time"

	"github.com/stellhub/stellar/discovery"
	"github.com/stellhub/stellar/loadbalancer"
	"google.golang.org/grpc/attributes"
	grpcbalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type endpointAttributeKey struct{}

func grpcResolverAddress(endpoint discovery.Endpoint, address string) resolver.Address {
	endpoint = discovery.NormalizeEndpoint(endpoint)
	attrs := attributes.New(endpointAttributeKey{}, endpoint)
	return resolver.Address{
		Addr:               address,
		BalancerAttributes: attrs,
		Attributes:         attrs,
	}
}

func endpointFromResolverAddress(address resolver.Address) (discovery.Endpoint, bool) {
	if address.BalancerAttributes != nil {
		if endpoint, ok := address.BalancerAttributes.Value(endpointAttributeKey{}).(discovery.Endpoint); ok {
			return discovery.NormalizeEndpoint(endpoint), true
		}
	}
	if address.Attributes != nil {
		if endpoint, ok := address.Attributes.Value(endpointAttributeKey{}).(discovery.Endpoint); ok {
			return discovery.NormalizeEndpoint(endpoint), true
		}
	}
	return discovery.Endpoint{}, false
}

func newLoadBalancerBuilder(name string, policy string, router loadbalancer.Router, service string) grpcbalancer.Builder {
	return base.NewBalancerBuilder(name, &stellarPickerBuilder{
		picker:  loadbalancer.New(policy),
		router:  router,
		service: service,
	}, base.Config{HealthCheck: true})
}

type stellarPickerBuilder struct {
	picker  *loadbalancer.Picker
	router  loadbalancer.Router
	service string
}

func (b *stellarPickerBuilder) Build(info base.PickerBuildInfo) grpcbalancer.Picker {
	entries := make([]grpcEndpointEntry, 0, len(info.ReadySCs))
	for subConn, subConnInfo := range info.ReadySCs {
		endpoint, ok := endpointFromResolverAddress(subConnInfo.Address)
		if !ok {
			endpoint = discovery.Endpoint{
				Protocol: "grpc",
				Host:     subConnInfo.Address.Addr,
			}
		}
		entries = append(entries, grpcEndpointEntry{
			endpoint: endpoint,
			subConn:  subConn,
		})
	}
	return &stellarGRPCPicker{
		entries: entries,
		picker:  b.picker,
		router:  b.router,
		service: b.service,
	}
}

type grpcEndpointEntry struct {
	endpoint discovery.Endpoint
	subConn  grpcbalancer.SubConn
}

type stellarGRPCPicker struct {
	entries []grpcEndpointEntry
	picker  *loadbalancer.Picker
	router  loadbalancer.Router
	service string
}

func (p *stellarGRPCPicker) Pick(info grpcbalancer.PickInfo) (grpcbalancer.PickResult, error) {
	if len(p.entries) == 0 {
		return grpcbalancer.PickResult{}, grpcbalancer.ErrNoSubConnAvailable
	}
	ctx := info.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	request := loadBalancerRequestFromPickInfo(ctx, info, p.service)
	endpoints := make([]discovery.Endpoint, 0, len(p.entries))
	subConns := make(map[string]grpcbalancer.SubConn, len(p.entries))
	for _, entry := range p.entries {
		endpoint := discovery.NormalizeEndpoint(entry.endpoint)
		endpoints = append(endpoints, endpoint)
		subConns[loadbalancer.EndpointKey(endpoint)] = entry.subConn
	}
	var err error
	if p.router != nil {
		endpoints, err = p.router.Route(ctx, request, endpoints)
		if err != nil {
			return grpcbalancer.PickResult{}, err
		}
	}
	picker := p.picker
	if picker == nil {
		picker = loadbalancer.New(loadbalancer.DefaultPolicy)
	}
	pick, err := picker.Pick(ctx, request, endpoints)
	if err != nil {
		return grpcbalancer.PickResult{}, err
	}
	subConn, ok := subConns[loadbalancer.EndpointKey(pick.Endpoint)]
	if !ok {
		if pick.Done != nil {
			pick.Done(loadbalancer.Result{Err: discovery.ErrNoAvailableEndpoint})
		}
		return grpcbalancer.PickResult{}, grpcbalancer.ErrNoSubConnAvailable
	}
	start := time.Now()
	return grpcbalancer.PickResult{
		SubConn: subConn,
		Done: func(done grpcbalancer.DoneInfo) {
			if pick.Done != nil {
				pick.Done(loadbalancer.Result{
					Duration: time.Since(start),
					Err:      done.Err,
				})
			}
		},
	}, nil
}

func loadBalancerRequestFromPickInfo(ctx context.Context, info grpcbalancer.PickInfo, service string) loadbalancer.Request {
	grpcService, method := splitFullMethod(info.FullMethodName)
	if service == "" {
		service = grpcService
	}
	return loadbalancer.MergeContextRequest(ctx, loadbalancer.Request{
		Protocol:   "grpc",
		Service:    service,
		Operation:  info.FullMethodName,
		Method:     method,
		Path:       info.FullMethodName,
		Headers:    map[string][]string(headersFromOutgoingContext(ctx)),
		Attributes: map[string]any{"grpc.service": grpcService},
	})
}
