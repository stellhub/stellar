package boot

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stellhub/stellar/config"
)

func TestNewConfiguredGovernanceRequiresConfigCenterEndpoint(t *testing.T) {
	enabled := true
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			Enabled: &enabled,
		},
	}.Normalize()

	_, err := NewConfigured(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "requires config_center endpoint") {
		t.Fatalf("expected governance config center endpoint error, got %v", err)
	}
}

func TestNewConfiguredGovernanceRouteStarterIsOptional(t *testing.T) {
	enabled := true
	disabled := false
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			Enabled: &enabled,
			ConfigCenter: config.GovernanceConfigCenterConfig{
				Endpoint: "http://localhost:8060",
			},
			Route: config.GovernanceFeatureConfig{
				Enabled: &disabled,
			},
		},
	}.Normalize()

	app, err := NewConfigured(context.Background(), cfg)
	if err != nil {
		t.Fatalf("new configured app: %v", err)
	}
	if app.governanceRouteEnabled {
		t.Fatalf("expected governance route to be disabled")
	}
	if containsStarter(app.starters, "governance-route") {
		t.Fatalf("did not expect route starter to be registered")
	}
	if !containsStarter(app.starters, "governance-stellorbit") {
		t.Fatalf("expected stellorbit starter to be registered")
	}
}

func TestNewConfiguredGovernanceRegistersRouteStarter(t *testing.T) {
	enabled := true
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			Enabled: &enabled,
			ConfigCenter: config.GovernanceConfigCenterConfig{
				Endpoint: "http://localhost:8060",
			},
			Route: config.GovernanceFeatureConfig{
				Enabled: &enabled,
			},
		},
	}.Normalize()

	app, err := NewConfigured(context.Background(), cfg)
	if err != nil {
		t.Fatalf("new configured app: %v", err)
	}
	if !app.governanceRouteEnabled {
		t.Fatalf("expected governance route to be enabled")
	}
	if !containsStarter(app.starters, "governance-route") {
		t.Fatalf("expected route starter to be registered")
	}
}

func TestConfiguredHeaderRateLimitRulesBuildServerScopes(t *testing.T) {
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			RateLimit: config.GovernanceRateLimitConfig{
				Headers: []config.GovernanceHeaderRateLimitConfig{
					{Header: "X-Tenant-ID", Rate: 3, Burst: 4},
					{Transport: "grpc", Header: "x-api-key", CoordinationMode: "global_quota"},
				},
			},
		},
	}.Normalize()

	rules, err := configuredHeaderRateLimitRules(cfg.Governance)
	if err != nil {
		t.Fatalf("build header rate limit rules: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected two rules, got %#v", rules)
	}
	if rules[0].Scope.Transport != "http.server" || rules[1].Scope.Transport != "grpc.server" {
		t.Fatalf("unexpected scopes %#v %#v", rules[0].Scope, rules[1].Scope)
	}
	if rules[0].Spec["limitMode"] != "HEADER" || rules[0].Spec["rate"] != int64(3) || rules[0].Spec["burst"] != int64(4) {
		t.Fatalf("unexpected http header rule spec %#v", rules[0].Spec)
	}
	if rules[1].Spec["coordinationMode"] != "GLOBAL_QUOTA" || rules[1].Spec["mode"] != "distributed" {
		t.Fatalf("unexpected grpc header coordination %#v", rules[1].Spec)
	}
	extractor, ok := rules[1].Spec["keyExtractor"].(map[string]any)
	if !ok {
		t.Fatalf("expected key extractor map, got %#v", rules[1].Spec["keyExtractor"])
	}
	keys, ok := extractor["keys"].([]any)
	if !ok || len(keys) != 1 {
		t.Fatalf("expected one key extractor entry, got %#v", extractor)
	}
}

func TestDistributedHeaderRateLimitRequiresPulsarAddress(t *testing.T) {
	disabled := false
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			RateLimit: config.GovernanceRateLimitConfig{
				Distributed: config.GovernanceDistributedRateLimitConfig{
					Enabled: &disabled,
				},
				Headers: []config.GovernanceHeaderRateLimitConfig{
					{Transport: "grpc", Header: "x-api-key", CoordinationMode: "global_quota"},
				},
			},
		},
	}.Normalize()

	_, err := newGovernanceRateLimitStarter(context.Background(), cfg.Governance, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "distributed rate limit requires") {
		t.Fatalf("expected distributed rate limit address error, got %v", err)
	}
}

func TestDistributedRateLimitEnabledCanDisableUnusedAddress(t *testing.T) {
	disabled := false
	cfg := config.GovernanceRateLimitConfig{
		Distributed: config.GovernanceDistributedRateLimitConfig{
			Enabled: &disabled,
			Address: "127.0.0.1:19091",
		},
	}
	if distributedRateLimitClientEnabled(cfg) {
		t.Fatalf("expected explicit disabled distributed config without distributed rules to stay disabled")
	}
}

func TestConfiguredHeaderRateLimitRejectsInvalidCoordinationMode(t *testing.T) {
	cfg := config.Config{
		AppName:     "order-service",
		Environment: config.EnvDev,
		Governance: &config.GovernanceConfig{
			RateLimit: config.GovernanceRateLimitConfig{
				Headers: []config.GovernanceHeaderRateLimitConfig{
					{Header: "x-api-key", CoordinationMode: "global_quato"},
				},
			},
		},
	}.Normalize()

	_, err := configuredHeaderRateLimitRules(cfg.Governance)
	if err == nil || !strings.Contains(err.Error(), "coordination_mode") {
		t.Fatalf("expected invalid coordination mode error, got %v", err)
	}
}

func TestInitialSyncErrorHonorsFailFastOnBootstrap(t *testing.T) {
	syncErr := errors.New("sync failed")
	if err := (&stellorbitGovernanceStarter{}).initialSyncError(syncErr); err != nil {
		t.Fatalf("expected non fail-fast starter to ignore initial sync error, got %v", err)
	}
	starter := &stellorbitGovernanceStarter{
		cfg: config.GovernanceConfig{
			FailFastOnBootstrap: true,
		},
	}
	if err := starter.initialSyncError(syncErr); !errors.Is(err, syncErr) {
		t.Fatalf("expected fail-fast starter to return sync error, got %v", err)
	}
}

func containsStarter(starters []Starter, name string) bool {
	for _, starter := range starters {
		if starter != nil && starter.Name() == name {
			return true
		}
	}
	return false
}
