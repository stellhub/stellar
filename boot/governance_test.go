package boot

import (
	"context"
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

func containsStarter(starters []Starter, name string) bool {
	for _, starter := range starters {
		if starter != nil && starter.Name() == name {
			return true
		}
	}
	return false
}
