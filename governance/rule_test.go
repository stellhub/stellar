package governance_test

import (
	"testing"

	"github.com/stellhub/stellar/governance"
)

func TestStoreReplacesSnapshotAtomically(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Version: "v1",
		Rules: []governance.Rule{
			{ID: "rate-limit", Kind: governance.RuleKindRateLimit, Priority: 20},
		},
	})

	store.Replace(governance.Snapshot{
		Version: "v2",
		Rules: []governance.Rule{
			{ID: "auth", Kind: governance.RuleKindAuthentication, Priority: 10},
		},
	})

	snapshot := store.Snapshot()
	if snapshot.Version != "v2" {
		t.Fatalf("unexpected version %q", snapshot.Version)
	}
	if len(snapshot.Rules) != 1 || snapshot.Rules[0].ID != "auth" {
		t.Fatalf("unexpected rules %#v", snapshot.Rules)
	}
}

func TestStoreRulesReturnsSortedCopies(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{
			{ID: "b", Kind: governance.RuleKindRateLimit, Priority: 20, Spec: map[string]any{"qps": 100}},
			{ID: "a", Kind: governance.RuleKindRateLimit, Priority: 10, Spec: map[string]any{"qps": 50}},
			{ID: "auth", Kind: governance.RuleKindAuthentication, Priority: 1},
		},
	})

	rules := store.Rules(governance.RuleKindRateLimit, nil)
	if len(rules) != 2 {
		t.Fatalf("unexpected rule count %d", len(rules))
	}
	if rules[0].ID != "a" || rules[1].ID != "b" {
		t.Fatalf("unexpected rule order %#v", rules)
	}

	rules[0].Spec["qps"] = 1
	again := store.Rules(governance.RuleKindRateLimit, nil)
	if again[0].Spec["qps"] != 50 {
		t.Fatalf("expected store snapshot to be immutable from returned copies")
	}
}
