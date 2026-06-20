package circuitbreaker

import (
	"context"
	"errors"
	"testing"

	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
)

func TestPolicyOpensBreakerAfterFailureThreshold(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{{
			ID:      "client-breaker",
			Kind:    governance.RuleKindCircuitBreaker,
			Enabled: true,
			Scope: governance.Scope{
				Transport: "http.client",
				Service:   "order-service",
			},
			Spec: map[string]any{
				"failure_threshold": 1,
				"open_timeout":      "1m",
			},
		}},
	})
	policy := NewPolicy(store, nil)
	inv := &interceptor.Invocation{
		Kind:     interceptor.KindHTTPClient,
		Protocol: "http",
		Service:  "order-service",
		Method:   "GET",
		Path:     "/orders",
	}

	_, err := policy.Intercept(context.Background(), inv, nil, func(context.Context, *interceptor.Invocation, any) (any, error) {
		return nil, errors.New("downstream unavailable")
	})
	if err == nil {
		t.Fatalf("expected downstream failure")
	}

	_, err = policy.Intercept(context.Background(), inv, nil, func(context.Context, *interceptor.Invocation, any) (any, error) {
		return "ok", nil
	})
	if err == nil {
		t.Fatalf("expected open breaker rejection")
	}
}
