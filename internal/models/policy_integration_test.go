package models

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadPolicyStrategyCount verifies that the production policy ships with the
// full strategy catalog. This guards against packaging the app with a truncated
// `policy.v1.json`, which would cause the ticker entry dropdown to show only a
// subset of strategies.
func TestLoadPolicyStrategyCount(t *testing.T) {
	policyPath := filepath.Join("data", "policy.v1.json")

	if _, err := os.Stat(policyPath); err != nil {
		t.Skipf("policy file not present at %s: %v", policyPath, err)
	}

	policy, err := LoadPolicy(policyPath)
	if err != nil {
		t.Fatalf("failed to load policy: %v", err)
	}

	const expectedStrategies = 12
	if len(policy.Strategies) < expectedStrategies {
		t.Fatalf("policy only defines %d strategies (expected at least %d)", len(policy.Strategies), expectedStrategies)
	}

	// Spot-check that a few core strategies are present.
	core := []string{"Alt10", "Alt26", "Alt22", "Alt47", "Baseline"}
	for _, id := range core {
		if _, ok := policy.Strategies[id]; !ok {
			t.Fatalf("strategy %s missing from policy", id)
		}
	}
}
