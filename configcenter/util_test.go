package configcenter

import (
	"testing"

	"github.com/stellhub/stellar/config"
)

func TestSourceSpecsDefaultToApplicationYAML(t *testing.T) {
	cfg := (&config.Config{
		AppName: "orders",
		ConfigCenter: &config.ConfigCenterConfig{
			Adapter: "stellnula",
		},
	}).Normalize()

	specs := sourceSpecs(cfg.ConfigCenter)
	if len(specs) != 1 {
		t.Fatalf("expected one default source, got %#v", specs)
	}
	if specs[0].configKey != defaultConfigDataID || specs[0].dataID != defaultConfigDataID || !specs[0].required {
		t.Fatalf("unexpected default source spec %#v", specs[0])
	}
}

func TestSourceSpecsUseExplicitSources(t *testing.T) {
	required := false
	cfg := &config.ConfigCenterConfig{
		DataID: "application.yaml",
		Group:  "DEFAULT_GROUP",
		Sources: []config.ConfigCenterSourceConfig{{
			DataID:   "orders.yaml",
			Format:   "yaml",
			Required: &required,
		}},
	}

	specs := sourceSpecs(cfg)
	if len(specs) != 1 {
		t.Fatalf("expected one source, got %#v", specs)
	}
	if specs[0].dataID != "orders.yaml" || specs[0].configKey != "orders.yaml" || specs[0].group != "DEFAULT_GROUP" || specs[0].required {
		t.Fatalf("unexpected source spec %#v", specs[0])
	}
}

func TestSourceFileNameUsesYAMLFormat(t *testing.T) {
	name := sourceFileName(Source{Key: "orders", Format: "yaml"})
	if name != "orders.yaml" {
		t.Fatalf("unexpected source file name %q", name)
	}
}

func TestSourceSupportsYAMLRejectsExplicitNonYAML(t *testing.T) {
	if sourceSupportsYAML(Source{Key: "orders.properties"}) {
		t.Fatalf("expected properties source to be rejected")
	}
	if sourceSupportsYAML(Source{Key: "orders", Format: "properties"}) {
		t.Fatalf("expected explicit properties format to be rejected")
	}
	if sourceSupportsYAML(Source{Key: "orders.yaml", Format: "properties"}) {
		t.Fatalf("expected explicit properties format to override yaml extension")
	}
	if !sourceSupportsYAML(Source{Key: "orders", Format: "yaml"}) {
		t.Fatalf("expected yaml format to be accepted")
	}
	if !sourceSupportsYAML(Source{Key: "orders.yaml"}) {
		t.Fatalf("expected yaml extension to be accepted")
	}
}
