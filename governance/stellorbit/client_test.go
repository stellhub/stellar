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
