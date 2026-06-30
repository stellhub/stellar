package stellorbit

import (
	"context"
	"errors"
	"testing"

	"github.com/stellhub/stellar/governance"
	orbitgovernance "github.com/stellhub/stellorbit-go-sdk/governance"
)

func TestConvertRouteRuleFlattensFilterSpec(t *testing.T) {
	rule := orbitgovernance.Rule{
		RuleID:        "route-v1",
		RuleName:      "route-v1",
		ConfigKey:     "governance/order-service/route",
		RuleType:      orbitgovernance.RuleTypeRoute,
		TargetService: "order-service",
		Status:        orbitgovernance.RuleStatusActive,
		Priority:      10,
		Revision:      7,
		Checksum:      "abc",
		Content: map[string]any{
			"scope": map[string]any{
				"transport": "http.client",
				"method":    "GET",
				"path":      "/orders",
			},
			"filter": map[string]any{
				"endpoint_name": "http-v1",
				"zone":          "zone-a",
			},
		},
	}

	converted, ok := ConvertRule(rule)
	if !ok {
		t.Fatalf("expected rule to be converted")
	}
	if converted.Kind != governance.RuleKindRoute || !converted.Enabled {
		t.Fatalf("unexpected converted rule %#v", converted)
	}
	if converted.Scope.Service != "order-service" || converted.Scope.Transport != "http.client" || converted.Scope.Method != "GET" {
		t.Fatalf("unexpected converted scope %#v", converted.Scope)
	}
	if converted.Spec["endpoint_name"] != "http-v1" || converted.Spec["zone"] != "zone-a" {
		t.Fatalf("expected filter fields to be promoted, got %#v", converted.Spec)
	}
	if converted.Metadata["rule_checksum"] != "abc" || converted.Version != "abc" {
		t.Fatalf("unexpected metadata/version %#v version=%q", converted.Metadata, converted.Version)
	}
}

func TestConvertRateLimitRulePromotesEnterpriseSpec(t *testing.T) {
	rule := orbitgovernance.Rule{
		RuleID:        "limit-v2",
		RuleName:      "limit-v2",
		ConfigKey:     "governance/order-service/limit",
		RuleType:      orbitgovernance.RuleTypeRateLimit,
		TargetService: "order-service",
		Status:        orbitgovernance.RuleStatusActive,
		Content: map[string]any{
			"scope": map[string]any{
				"transport": "http.server",
			},
			"limit": map[string]any{
				"limitMode":         "HEADER",
				"limitType":         "HEADER",
				"limitAlgorithm":    "TOKEN_BUCKET",
				"trafficProtocol":   "HTTP",
				"executionLocation": "APPLICATION",
				"coordinationMode":  "GLOBAL_QUOTA",
				"quotaConfig": map[string]any{
					"limit": 10,
				},
				"burstConfig": map[string]any{
					"capacity": 20,
				},
				"keyExtractor": map[string]any{
					"keys": []any{map[string]any{
						"name":      "tenant",
						"source":    "HEADER",
						"key":       "x-tenant-id",
						"normalize": "LOWERCASE",
					}},
				},
				"fallbackPolicy": map[string]any{
					"failPolicy": "FAIL_CLOSED",
				},
			},
		},
	}

	converted, ok := ConvertRule(rule)
	if !ok {
		t.Fatalf("expected rule to be converted")
	}
	if converted.Kind != governance.RuleKindRateLimit || converted.Scope.Transport != "http.server" {
		t.Fatalf("unexpected converted rule %#v", converted)
	}
	if converted.Spec["limitMode"] != "HEADER" || converted.Spec["coordinationMode"] != "GLOBAL_QUOTA" || converted.Spec["mode"] != "distributed" {
		t.Fatalf("expected enterprise rate limit fields, got %#v", converted.Spec)
	}
	extractor, ok := converted.Spec["keyExtractor"].(map[string]any)
	if !ok {
		t.Fatalf("expected keyExtractor map, got %#v", converted.Spec["keyExtractor"])
	}
	keys, ok := extractor["keys"].([]any)
	if !ok || len(keys) != 1 {
		t.Fatalf("expected keyExtractor keys, got %#v", extractor)
	}
}

func TestSyncStoreReturnsErrorForMissingTarget(t *testing.T) {
	_, err := SyncStore(context.Background(), governance.NewStore(), nil)
	if !errors.Is(err, ErrSyncTargetRequired) {
		t.Fatalf("expected sync target error, got %v", err)
	}
}

func TestConvertAuthRuleKind(t *testing.T) {
	cases := []struct {
		name string
		kind any
		want governance.RuleKind
	}{
		{name: "authn", kind: "authentication", want: governance.RuleKindAuthentication},
		{name: "signing", kind: "signing", want: governance.RuleKindSigning},
		{name: "default", kind: "authorization", want: governance.RuleKindAuthorization},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			converted, ok := ConvertRule(orbitgovernance.Rule{
				RuleID:        tt.name,
				RuleName:      tt.name,
				ConfigKey:     tt.name,
				RuleType:      orbitgovernance.RuleTypeAuth,
				TargetService: "order-service",
				Status:        orbitgovernance.RuleStatusActive,
				Content:       map[string]any{"kind": tt.kind},
			})
			if !ok {
				t.Fatalf("expected rule to be converted")
			}
			if converted.Kind != tt.want {
				t.Fatalf("expected kind %q, got %q", tt.want, converted.Kind)
			}
		})
	}
}
