package ratelimit

import (
	"context"
	"strings"
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

func TestModeFromSpecUsesCoordinationMode(t *testing.T) {
	if got := modeFromSpec(map[string]any{"limitMode": "QPS", "coordinationMode": "GLOBAL_QUOTA"}, "local"); got != "distributed" {
		t.Fatalf("expected distributed mode, got %q", got)
	}
	if got := modeFromSpec(map[string]any{"limitMode": "HEADER", "coordinationMode": "LOCAL_ONLY"}, "distributed"); got != "local" {
		t.Fatalf("expected local mode, got %q", got)
	}
}

func TestPolicyAppliesStaticHeaderKeyExtractorRule(t *testing.T) {
	policy := NewPolicy(PolicyOptions{
		Rules: []governance.Rule{{
			ID:      "header-limit",
			Kind:    governance.RuleKindRateLimit,
			Enabled: true,
			Scope: governance.Scope{
				Transport: "http.server",
			},
			Spec: map[string]any{
				"limitMode":        "HEADER",
				"coordinationMode": "LOCAL_ONLY",
				"rate":             1,
				"burst":            1,
				"keyExtractor": map[string]any{
					"keys": []any{map[string]any{
						"name":      "tenant",
						"source":    "HEADER",
						"key":       "x-tenant-id",
						"normalize": "lowercase",
					}},
				},
			},
		}},
	})
	inv := &interceptor.Invocation{
		Kind:     interceptor.KindHTTPServer,
		Protocol: "http",
		Service:  "order-service",
		Method:   "GET",
		Path:     "/orders",
		Headers:  interceptor.Header{"X-Tenant-ID": []string{"Tenant-A"}},
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

func TestDistributedRequestDefersKeyExtractorQuotaKey(t *testing.T) {
	policy := NewPolicy(PolicyOptions{})
	rule := governance.Rule{
		ID:      "distributed-header-limit",
		Kind:    governance.RuleKindRateLimit,
		Enabled: true,
		Spec: map[string]any{
			"limitMode":        "HEADER",
			"coordinationMode": "GLOBAL_QUOTA",
			"keyExtractor": map[string]any{
				"keys": []any{map[string]any{
					"name":   "tenant",
					"source": "GRPC_METADATA",
					"key":    "x-tenant-id",
				}},
			},
		},
	}
	inv := &interceptor.Invocation{
		Kind:     interceptor.KindGRPCServer,
		Protocol: "grpc",
		Service:  "order-service",
		Method:   "GetOrder",
		Path:     "/orders.OrderService/GetOrder",
		Headers:  interceptor.Header{"x-tenant-id": []string{"tenant-a"}},
	}

	request := policy.limitRequest(inv, rule)
	if request.QuotaKey != "" {
		t.Fatalf("expected distributed key extractor quota key to be deferred, got %q", request.QuotaKey)
	}
	if request.Attributes["x-tenant-id"] != "tenant-a" || request.Attributes["metadata.x-tenant-id"] != "tenant-a" {
		t.Fatalf("expected grpc metadata in attributes, got %#v", request.Attributes)
	}
}

func TestUnsupportedLocalLimitModesHonorFailClosed(t *testing.T) {
	for _, limitMode := range []string{"CONCURRENCY", "CONNECTION", "CUSTOM"} {
		t.Run(limitMode, func(t *testing.T) {
			policy := NewPolicy(PolicyOptions{
				Rules: []governance.Rule{{
					ID:      "local-" + strings.ToLower(limitMode),
					Kind:    governance.RuleKindRateLimit,
					Enabled: true,
					Scope: governance.Scope{
						Transport: "http.server",
					},
					Spec: map[string]any{
						"limitMode":        limitMode,
						"coordinationMode": "LOCAL_ONLY",
						"fallbackPolicy": map[string]any{
							"failPolicy": "FAIL_CLOSED",
						},
					},
				}},
			})
			inv := &interceptor.Invocation{
				Kind:     interceptor.KindHTTPServer,
				Protocol: "http",
				Service:  "order-service",
				Method:   "GET",
				Path:     "/orders",
			}
			next := func(context.Context, *interceptor.Invocation, any) (any, error) {
				return "ok", nil
			}

			if _, err := policy.Intercept(context.Background(), inv, nil, next); err == nil {
				t.Fatalf("expected unsupported local %s mode to fail closed", limitMode)
			}
		})
	}
}

func TestLocalTokenBucketLimitModesAreSupported(t *testing.T) {
	for _, limitMode := range []string{"QPS", "HEADER", "HOT_KEY", "QUOTA", "BANDWIDTH", "MODEL"} {
		t.Run(limitMode, func(t *testing.T) {
			policy := NewPolicy(PolicyOptions{
				Rules: []governance.Rule{{
					ID:      "local-" + strings.ToLower(limitMode),
					Kind:    governance.RuleKindRateLimit,
					Enabled: true,
					Scope: governance.Scope{
						Transport: "http.server",
					},
					Spec: map[string]any{
						"limitMode":        limitMode,
						"coordinationMode": "LOCAL_ONLY",
						"rate":             1,
						"burst":            1,
					},
				}},
			})
			inv := &interceptor.Invocation{
				Kind:     interceptor.KindHTTPServer,
				Protocol: "http",
				Service:  "order-service",
				Method:   "GET",
				Path:     "/orders",
				Headers:  interceptor.Header{"x-tenant-id": []string{strings.ToLower(limitMode)}},
			}
			next := func(context.Context, *interceptor.Invocation, any) (any, error) {
				return "ok", nil
			}

			if _, err := policy.Intercept(context.Background(), inv, nil, next); err != nil {
				t.Fatalf("expected first local %s request to be allowed: %v", limitMode, err)
			}
			if _, err := policy.Intercept(context.Background(), inv, nil, next); err == nil {
				t.Fatalf("expected second local %s request to be rate limited", limitMode)
			}
		})
	}
}
