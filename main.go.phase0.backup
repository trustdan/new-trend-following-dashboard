package main

import (
	"fmt"
	"os"
	"tf-engine/internal/config"
)

// Phase 0 stub - verifies infrastructure is working
// Full application will be built in Phase 1+
func main() {
	fmt.Println("TF-Engine 2.0 - Phase 0 Infrastructure Test")
	fmt.Println("===========================================")

	// Test 1: Load and verify policy
	fmt.Println("\n1. Testing policy signature validation...")
	// Policy validation happens in scripts/verify_policy_hash.go
	// Just verify the file exists here
	if _, err := os.Stat("data/policy.v1.json"); err != nil {
		fmt.Println("❌ Policy file not found")
		os.Exit(1)
	}
	fmt.Println("✅ Policy file exists")

	// Test 2: Load feature flags
	fmt.Println("\n2. Testing feature flags system...")
	flags, err := config.LoadFeatureFlags("feature.flags.json")
	if err != nil {
		fmt.Printf("❌ Failed to load feature flags: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Feature flags loaded (version %s)\n", flags.Version)

	// Test 3: Verify Phase 2 features are disabled
	fmt.Println("\n3. Verifying Phase 2 features are disabled...")
	phase2Flags := flags.ListPhase2Flags()
	allDisabled := true
	for _, flagName := range phase2Flags {
		if flags.IsEnabled(flagName) {
			fmt.Printf("❌ Phase 2 feature '%s' is enabled (should be disabled)\n", flagName)
			allDisabled = false
		}
	}
	if allDisabled {
		fmt.Printf("✅ All %d Phase 2 features are disabled\n", len(phase2Flags))
	} else {
		os.Exit(1)
	}

	// Test 4: Verify folder structure
	fmt.Println("\n4. Verifying folder structure...")
	requiredDirs := []string{
		"internal/screens/_post_mvp",
		"internal/testdata",
		"internal/config",
		"data/backups",
		"logs",
		"scripts",
	}
	allExist := true
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); err != nil {
			fmt.Printf("❌ Directory missing: %s\n", dir)
			allExist = false
		}
	}
	if allExist {
		fmt.Printf("✅ All %d required directories exist\n", len(requiredDirs))
	} else {
		os.Exit(1)
	}

	fmt.Println("\n===========================================")
	fmt.Println("✅ Phase 0 Infrastructure Test PASSED")
	fmt.Println("\nNext steps:")
	fmt.Println("  - Phase 1: Build foundation layer (navigation, persistence, cooldown)")
	fmt.Println("  - See plans/roadmap.md for details")
}
