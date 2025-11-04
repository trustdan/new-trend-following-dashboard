package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tf-engine/internal/models"
)

const settingsFile = "data/ui/settings.json"

// SaveSettings saves user settings to disk
func SaveSettings(settings *models.Settings) error {
	// Ensure directory exists
	dir := filepath.Dir(settingsFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal settings to JSON
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(settingsFile, data, 0644)
}

// LoadSettings loads user settings from disk
func LoadSettings() (*models.Settings, error) {
	// Check if file exists
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// Return default settings if file doesn't exist
		return models.DefaultSettings(), nil
	}

	// Read file
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		return models.DefaultSettings(), err
	}

	// Unmarshal JSON
	var settings models.Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return models.DefaultSettings(), err
	}

	return &settings, nil
}
