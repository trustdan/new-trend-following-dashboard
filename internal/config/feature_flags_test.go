package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFeatureFlags(t *testing.T) {
	// Create a temporary feature flags file
	tmpDir := t.TempDir()
	flagsPath := filepath.Join(tmpDir, "feature.flags.json")

	flagsJSON := `{
  "version": "1.0.0",
  "flags": {
    "test_feature": {
      "enabled": true,
      "description": "Test feature",
      "phase": 1,
      "since_version": "1.0.0"
    },
    "disabled_feature": {
      "enabled": false,
      "description": "Disabled feature",
      "phase": 2,
      "since_version": "2.0.0"
    }
  }
}`

	if err := os.WriteFile(flagsPath, []byte(flagsJSON), 0644); err != nil {
		t.Fatalf("Failed to create test flags file: %v", err)
	}

	// Test loading
	flags, err := LoadFeatureFlags(flagsPath)
	if err != nil {
		t.Fatalf("Failed to load feature flags: %v", err)
	}

	if flags.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", flags.Version)
	}

	if len(flags.Flags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(flags.Flags))
	}
}

func TestIsEnabled(t *testing.T) {
	flags := &FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]FeatureFlag{
			"enabled_feature": {
				Enabled:      true,
				Description:  "Enabled",
				Phase:        1,
				SinceVersion: "1.0.0",
			},
			"disabled_feature": {
				Enabled:      false,
				Description:  "Disabled",
				Phase:        2,
				SinceVersion: "2.0.0",
			},
		},
	}

	tests := []struct {
		name     string
		flagName string
		expected bool
	}{
		{"Enabled feature", "enabled_feature", true},
		{"Disabled feature", "disabled_feature", false},
		{"Non-existent feature", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flags.IsEnabled(tt.flagName)
			if result != tt.expected {
				t.Errorf("IsEnabled(%s) = %v, expected %v", tt.flagName, result, tt.expected)
			}
		})
	}
}

func TestGetFlag(t *testing.T) {
	flags := &FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]FeatureFlag{
			"test_feature": {
				Enabled:      true,
				Description:  "Test",
				Phase:        1,
				SinceVersion: "1.0.0",
			},
		},
	}

	// Test existing flag
	flag := flags.GetFlag("test_feature")
	if flag == nil {
		t.Fatal("Expected flag to exist")
	}
	if !flag.Enabled {
		t.Error("Expected flag to be enabled")
	}
	if flag.Description != "Test" {
		t.Errorf("Expected description 'Test', got '%s'", flag.Description)
	}

	// Test non-existent flag
	nilFlag := flags.GetFlag("nonexistent")
	if nilFlag != nil {
		t.Error("Expected nil for non-existent flag")
	}
}

func TestListEnabledFlags(t *testing.T) {
	flags := &FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]FeatureFlag{
			"enabled1": {Enabled: true, Phase: 1, SinceVersion: "1.0.0"},
			"disabled": {Enabled: false, Phase: 1, SinceVersion: "1.0.0"},
			"enabled2": {Enabled: true, Phase: 2, SinceVersion: "2.0.0"},
		},
	}

	enabled := flags.ListEnabledFlags()
	if len(enabled) != 2 {
		t.Errorf("Expected 2 enabled flags, got %d", len(enabled))
	}

	// Check that both enabled flags are in the list
	foundEnabled1 := false
	foundEnabled2 := false
	for _, name := range enabled {
		if name == "enabled1" {
			foundEnabled1 = true
		}
		if name == "enabled2" {
			foundEnabled2 = true
		}
	}
	if !foundEnabled1 || !foundEnabled2 {
		t.Error("Not all enabled flags were found in list")
	}
}

func TestListPhase2Flags(t *testing.T) {
	flags := &FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]FeatureFlag{
			"phase1_feature":  {Enabled: true, Phase: 1, SinceVersion: "1.0.0"},
			"phase2_feature1": {Enabled: false, Phase: 2, SinceVersion: "2.0.0"},
			"phase2_feature2": {Enabled: false, Phase: 2, SinceVersion: "2.1.0"},
		},
	}

	phase2 := flags.ListPhase2Flags()
	if len(phase2) != 2 {
		t.Errorf("Expected 2 Phase 2 flags, got %d", len(phase2))
	}

	// Verify phase 1 feature is not in the list
	for _, name := range phase2 {
		if name == "phase1_feature" {
			t.Error("Phase 1 feature should not be in Phase 2 list")
		}
	}
}

func TestLoadFeatureFlags_FileNotFound(t *testing.T) {
	_, err := LoadFeatureFlags("nonexistent.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadFeatureFlags_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	flagsPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(flagsPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadFeatureFlags(flagsPath)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}
