package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// FeatureFlag represents a single feature flag configuration
type FeatureFlag struct {
	Enabled      bool   `json:"enabled"`
	Description  string `json:"description"`
	Phase        int    `json:"phase"`
	SinceVersion string `json:"since_version"`
}

// FeatureFlags represents the complete feature flags configuration
type FeatureFlags struct {
	Version string                 `json:"version"`
	Flags   map[string]FeatureFlag `json:"flags"`
}

// LoadFeatureFlags loads feature flags from the specified path
func LoadFeatureFlags(path string) (*FeatureFlags, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read feature flags file: %w", err)
	}

	var flags FeatureFlags
	if err := json.Unmarshal(data, &flags); err != nil {
		return nil, fmt.Errorf("failed to parse feature flags JSON: %w", err)
	}

	return &flags, nil
}

// IsEnabled checks if a specific feature flag is enabled
// Returns false if the flag doesn't exist (fail-safe default)
func (ff *FeatureFlags) IsEnabled(flagName string) bool {
	if flag, exists := ff.Flags[flagName]; exists {
		return flag.Enabled
	}
	return false // Default to disabled for safety
}

// GetFlag returns the complete FeatureFlag configuration for a given flag name
// Returns nil if the flag doesn't exist
func (ff *FeatureFlags) GetFlag(flagName string) *FeatureFlag {
	if flag, exists := ff.Flags[flagName]; exists {
		return &flag
	}
	return nil
}

// ListEnabledFlags returns a list of all enabled flag names
func (ff *FeatureFlags) ListEnabledFlags() []string {
	enabled := []string{}
	for name, flag := range ff.Flags {
		if flag.Enabled {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

// ListPhase2Flags returns a list of all Phase 2 feature flag names
func (ff *FeatureFlags) ListPhase2Flags() []string {
	phase2 := []string{}
	for name, flag := range ff.Flags {
		if flag.Phase == 2 {
			phase2 = append(phase2, name)
		}
	}
	return phase2
}
