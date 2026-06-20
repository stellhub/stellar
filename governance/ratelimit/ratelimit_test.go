package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
)

func TestLocalLimiterRejectsWhenBucketIsEmpty(t *testing.T) {
	limiter := NewLocalLimiter()
	request := LimitRequest{
		Rule:     governance.Rule{ID: "local-limit"},
		Resource: "orders",
		QuotaKey: "tenant-a",
		Rate:     1,
		Burst:    1,
		Cost:     1,
	}

	first, err := limiter.Allow(context.Background(), request)
	if err != nil {
		t.Fatalf("allow first request: %v", err)
	}
	if !first.Allowed {
		t.Fatalf("expected first request to be allowed")
	}
	second, err := limiter.Allow(context.Background(), request)
	if err != nil {
		t.Fatalf("allow second request: %v", err)
	}
	if second.Allowed || !second.Limited || second.RetryAfter <= 0 {
		t.Fatalf("expected second request to be limited, got %#v", second)
	}
}

func TestLocalLimiterWaitHonorsContext(t *testing.T) {
	limiter := NewLocalLimiter()
	request := LimitRequest{
		Rule:     governance.Rule{ID: "blocking-limit"},
		Resource: "orders",
		QuotaKey: "tenant-a",
		Rate:     1,
		Burst:    1,
		Cost:     1,
	}
	_, _ = limiter.Allow(context.Background(), request)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	decision, err := limiter.Wait(ctx, request)
	if err != nil {
		t.Fatalf("wait request: %v", err)
	}
	if decision.Allowed || !decision.Limited {
		t.Fatalf("expected context-limited decision, got %#v", decision)
	}
}

func TestPolicyRejectsMatchingRateLimitRule(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{{
			ID:      "server-limit",
			Kind:    governance.RuleKindRateLimit,
			Enabled: true,
			Scope: governance.Scope{
				Transport: "http.server",
				Service:   "GET /orders",
			},
			Spec: map[string]any{
				"rate":  1,
				"burst": 1,
			},
		}},
	})
	policy := NewPolicy(PolicyOptions{Store: store})
	inv := &interceptor.Invocation{
		Kind:     interceptor.KindHTTPServer,
		Protocol: "http",
		Service:  "GET /orders",
		Method:   "GET",
		Path:     "/orders",
		Headers:  interceptor.Header{"x-tenant-id": []string{"tenant-a"}},
	}
	next := func(context.Context, *interceptor.Invocation, any) (any, error) {
		return "ok", nil
	}

	if _, err := policy.Intercept(context.Background(), inv, nil, next); err != nil {
		t.Fatalf("first intercept should be allowed: %v", err)
	}
	if _, err := policy.Intercept(context.Background(), inv, nil, next); err == nil {
		t.Fatalf("expected second intercept to be rate limited")
	}
}

func TestPolicyUsesDistributedFallbackWhenClientIsMissing(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{{
			ID:      "distributed-limit",
			Kind:    governance.RuleKindRateLimit,
			Enabled: true,
			Scope: governance.Scope{
				Transport: "http.client",
				Service:   "order-service",
			},
			Spec: map[string]any{
				"mode": "distributed",
			},
		}},
	})
	policy := NewPolicy(PolicyOptions{Store: store, DistributedFallback: "fail_closed"})
	inv := &interceptor.Invocation{
		Kind:     interceptor.KindHTTPClient,
		Protocol: "http",
		Service:  "order-service",
		Method:   "GET",
		Path:     "/orders",
	}
	nextCalled := false
	next := func(context.Context, *interceptor.Invocation, any) (any, error) {
		nextCalled = true
		return "ok", nil
	}

	if _, err := policy.Intercept(context.Background(), inv, nil, next); err == nil {
		t.Fatalf("expected missing distributed client to use fail-closed fallback")
	}
	if nextCalled {
		t.Fatalf("did not expect next handler to be called")
	}
}
