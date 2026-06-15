package loadbalancer

import (
	"context"
	"strings"
	"time"

	"github.com/stellhub/stellar/discovery"
)

const (
	DefaultPolicy        = PolicyP2C
	PolicyP2C            = "p2c"
	PolicyRoundRobin     = "round_robin"
	PolicyRandom         = "random"
	PolicyWeightedRound  = "weighted_round_robin"
	PolicyLeastRequest   = "least_request"
	PolicyConsistentHash = "consistent_hash"
)

type Request struct {
	Protocol   string
	Service    string
	Operation  string
	Method     string
	Path       string
	Target     string
	HashKey    string
	Headers    map[string][]string
	Attributes map[string]any
}

type Result struct {
	StatusCode int
	Duration   time.Duration
	Err        error
}

type Pick struct {
	Endpoint discovery.Endpoint
	Done     func(Result)
}

type Balancer interface {
	Pick(context.Context, Request, []discovery.Endpoint) (Pick, error)
}

type Router interface {
	Route(context.Context, Request, []discovery.Endpoint) ([]discovery.Endpoint, error)
}

type RouterFunc func(context.Context, Request, []discovery.Endpoint) ([]discovery.Endpoint, error)

func (f RouterFunc) Route(ctx context.Context, request Request, endpoints []discovery.Endpoint) ([]discovery.Endpoint, error) {
	if f == nil {
		return CloneEndpoints(endpoints), nil
	}
	return f(ctx, request, endpoints)
}

func NormalizePolicy(policy string) string {
	switch strings.ToLower(strings.TrimSpace(policy)) {
	case PolicyRoundRobin:
		return PolicyRoundRobin
	case PolicyRandom:
		return PolicyRandom
	case PolicyWeightedRound:
		return PolicyWeightedRound
	case PolicyLeastRequest:
		return PolicyLeastRequest
	case PolicyConsistentHash, "ring_hash", "hash":
		return PolicyConsistentHash
	case PolicyP2C, "power_of_two_choices", "least_request_p2c", "":
		return PolicyP2C
	default:
		return PolicyP2C
	}
}

func EndpointKey(endpoint discovery.Endpoint) string {
	endpoint = discovery.NormalizeEndpoint(endpoint)
	if endpoint.InstanceID != "" {
		return endpoint.InstanceID + "|" + endpoint.Name + "|" + endpoint.Protocol + "|" + endpoint.Address()
	}
	return endpoint.Name + "|" + endpoint.Protocol + "|" + endpoint.Address()
}

func CloneEndpoints(values []discovery.Endpoint) []discovery.Endpoint {
	copied := make([]discovery.Endpoint, 0, len(values))
	for _, value := range values {
		copied = append(copied, discovery.NormalizeEndpoint(value))
	}
	return copied
}

func HeaderValue(headers map[string][]string, key string) string {
	for current, values := range headers {
		if strings.EqualFold(current, key) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}
