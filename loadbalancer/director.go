package loadbalancer

import (
	"context"

	"github.com/stellhub/stellar/discovery"
)

type DirectorOption func(*Director)

type Director struct {
	resolver discovery.Resolver
	target   discovery.Target
	router   Router
	balancer Balancer
}

func NewDirector(resolver discovery.Resolver, target discovery.Target, options ...DirectorOption) *Director {
	director := &Director{
		resolver: resolver,
		target:   discovery.NormalizeTarget(target),
		balancer: New(DefaultPolicy),
	}
	for _, option := range options {
		if option != nil {
			option(director)
		}
	}
	if director.balancer == nil {
		director.balancer = New(DefaultPolicy)
	}
	return director
}

func WithRouter(router Router) DirectorOption {
	return func(director *Director) {
		director.router = router
	}
}

func WithBalancer(balancer Balancer) DirectorOption {
	return func(director *Director) {
		if balancer != nil {
			director.balancer = balancer
		}
	}
}

func WithPolicy(policy string) DirectorOption {
	return WithBalancer(New(policy))
}

func (d *Director) Pick(ctx context.Context, request Request) (Pick, error) {
	if d == nil || d.resolver == nil {
		return Pick{}, discovery.ErrNoAvailableEndpoint
	}
	endpoints, err := d.resolver.Resolve(ctx, d.target)
	if err != nil {
		return Pick{}, err
	}
	if d.router != nil {
		endpoints, err = d.router.Route(ctx, request, endpoints)
		if err != nil {
			return Pick{}, err
		}
	}
	if len(endpoints) == 0 {
		return Pick{}, discovery.ErrNoAvailableEndpoint
	}
	return d.balancer.Pick(ctx, request, endpoints)
}

func (d *Director) Close(ctx context.Context) error {
	if d == nil || d.resolver == nil {
		return nil
	}
	return d.resolver.Close(ctx)
}
