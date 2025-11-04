package screens

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// TestTickerEntry_Validate tests the Validate method
func TestTickerEntry_Validate(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*appcore.AppState)
		expected bool
	}{
		{
			name: "No trade - invalid",
			setup: func(state *appcore.AppState) {
				state.CurrentTrade = nil
			},
			expected: false,
		},
		{
			name: "Trade with no ticker - invalid",
			setup: func(state *appcore.AppState) {
				state.CurrentTrade = &models.Trade{
					Sector:   "Healthcare",
					Strategy: "Alt10",
				}
			},
			expected: false,
		},
		{
			name: "Trade with no strategy - invalid",
			setup: func(state *appcore.AppState) {
				state.CurrentTrade = &models.Trade{
					Sector: "Healthcare",
					Ticker: "UNH",
				}
			},
			expected: false,
		},
		{
			name: "Trade with ticker and strategy - valid",
			setup: func(state *appcore.AppState) {
				state.CurrentTrade = &models.Trade{
					Sector:   "Healthcare",
					Ticker:   "UNH",
					Strategy: "Alt10",
				}
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.Policy = createTestPolicy()
			tt.setup(state)

			window := test.NewWindow(nil)
			defer window.Close()

			screen := NewTickerEntry(state, window)

			// Act
			result := screen.Validate()

			// Assert
			if result != tt.expected {
				t.Errorf("Validate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTickerEntry_GetName tests the GetName method
func TestTickerEntry_GetName(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	name := screen.GetName()

	// Assert
	expected := "ticker_entry"
	if name != expected {
		t.Errorf("GetName() = %v, want %v", name, expected)
	}
}

// TestTickerEntry_SetNavCallbacks tests navigation callback registration
func TestTickerEntry_SetNavCallbacks(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	nextCalled := false
	backCalled := false
	cancelCalled := false

	// Act
	screen.SetNavCallbacks(
		func() error { nextCalled = true; return nil },
		func() error { backCalled = true; return nil },
		func() { cancelCalled = true },
	)

	// Trigger callbacks
	screen.onNext()
	screen.onBack()
	screen.onCancel()

	// Assert
	if !nextCalled {
		t.Error("onNext callback not called")
	}
	if !backCalled {
		t.Error("onBack callback not called")
	}
	if !cancelCalled {
		t.Error("onCancel callback not called")
	}
}

// TestTickerEntry_Render_NoSector tests rendering without sector selection
func TestTickerEntry_Render_NoSector(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}
}

// TestTickerEntry_Render_WithSector tests rendering with sector selected
func TestTickerEntry_Render_WithSector(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}
}

// TestTickerEntry_StrategyFiltering tests strategy filtering by sector
func TestTickerEntry_StrategyFiltering(t *testing.T) {
	tests := []struct {
		name             string
		sector           string
		expectedCount    int
		shouldContain    string
		shouldNotContain string
	}{
		{
			name:             "Healthcare sector - shows Alt10",
			sector:           "Healthcare",
			expectedCount:    5, // Healthcare has 5 strategies in test policy
			shouldContain:    "Alt10",
			shouldNotContain: "Alt26", // Alt26 is for Technology
		},
		{
			name:             "Technology sector - shows Alt26",
			sector:           "Technology",
			expectedCount:    5, // Technology has 5 strategies in test policy
			shouldContain:    "Alt26",
			shouldNotContain: "Alt43", // Alt43 is only for Healthcare
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.Policy = createTestPolicy()
			state.CurrentTrade = &models.Trade{
				Sector: tt.sector,
			}

			window := test.NewWindow(nil)
			defer window.Close()

			screen := NewTickerEntry(state, window)

			// Act
			strategies := screen.getFilteredStrategies()

			// Assert
			if len(strategies) != tt.expectedCount {
				t.Errorf("Expected %d strategies, got %d", tt.expectedCount, len(strategies))
			}

			// Check shouldContain
			found := false
			for _, s := range strategies {
				if contains(s, tt.shouldContain) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %s in strategies, but didn't", tt.shouldContain)
			}

			// Check shouldNotContain
			for _, s := range strategies {
				if contains(s, tt.shouldNotContain) {
					t.Errorf("Did not expect to find %s in strategies, but found it", tt.shouldNotContain)
				}
			}
		})
	}
}

// TestTickerEntry_StrategyFiltering_NoSector tests filtering when no sector selected
func TestTickerEntry_StrategyFiltering_NoSector(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	// No sector selected

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	strategies := screen.getFilteredStrategies()

	// Assert
	if len(strategies) != 0 {
		t.Errorf("Expected 0 strategies when no sector selected, got %d", len(strategies))
	}
}

// TestTickerEntry_TickerUppercase tests ticker auto-uppercase conversion
func TestTickerEntry_TickerUppercase(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	screen.tickerEntry.SetText("unh") // Lowercase input

	// Assert
	if screen.tickerEntry.Text != "UNH" {
		t.Errorf("Expected ticker to be uppercase 'UNH', got '%s'", screen.tickerEntry.Text)
	}

	if state.CurrentTrade.Ticker != "UNH" {
		t.Errorf("Expected trade ticker to be 'UNH', got '%s'", state.CurrentTrade.Ticker)
	}
}

// TestTickerEntry_StrategySelection tests strategy selection updates trade state
func TestTickerEntry_StrategySelection(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	screen.onStrategySelected("Alt10 - Profit Targets")

	// Assert
	if state.CurrentTrade.Strategy != "Alt10" {
		t.Errorf("Expected strategy to be 'Alt10', got '%s'", state.CurrentTrade.Strategy)
	}
}

// TestTickerEntry_ContinueButtonState tests continue button enable/disable logic
func TestTickerEntry_ContinueButtonState(t *testing.T) {
	tests := []struct {
		name            string
		ticker          string
		strategy        string
		expectedEnabled bool
	}{
		{
			name:            "No ticker, no strategy - disabled",
			ticker:          "",
			strategy:        "",
			expectedEnabled: false,
		},
		{
			name:            "Ticker only - disabled",
			ticker:          "UNH",
			strategy:        "",
			expectedEnabled: false,
		},
		{
			name:            "Strategy only - disabled",
			ticker:          "",
			strategy:        "Alt10",
			expectedEnabled: false,
		},
		{
			name:            "Both ticker and strategy - enabled",
			ticker:          "UNH",
			strategy:        "Alt10",
			expectedEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.Policy = createTestPolicy()
			state.CurrentTrade = &models.Trade{
				Sector:   "Healthcare",
				Ticker:   tt.ticker,
				Strategy: tt.strategy,
			}

			window := test.NewWindow(nil)
			defer window.Close()

			screen := NewTickerEntry(state, window)
			screen.Render() // Initialize UI components

			// Act
			screen.updateContinueButton()

			// Assert
			if tt.expectedEnabled && screen.continueBtn.Disabled() {
				t.Error("Expected continue button to be enabled, but it's disabled")
			}
			if !tt.expectedEnabled && !screen.continueBtn.Disabled() {
				t.Error("Expected continue button to be disabled, but it's enabled")
			}
		})
	}
}

// TestTickerEntry_StartCooldown tests cooldown activation
func TestTickerEntry_StartCooldown(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector:   "Healthcare",
		Ticker:   "UNH",
		Strategy: "Alt10",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	nextCalled := false
	screen.SetNavCallbacks(
		func() error { nextCalled = true; return nil },
		nil,
		nil,
	)

	// Act
	screen.startCooldownAndProceed()

	// Assert
	if !state.CooldownActive {
		t.Error("Expected cooldown to be active")
	}

	if state.CooldownStart == nil {
		t.Error("Expected cooldown start time to be set")
	}

	if !nextCalled {
		t.Error("Expected onNext to be called after starting cooldown")
	}
}

// TestTickerEntry_StartCooldown_InvalidData tests cooldown rejection with invalid data
func TestTickerEntry_StartCooldown_InvalidData(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
		// No ticker or strategy
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	nextCalled := false
	screen.SetNavCallbacks(
		func() error { nextCalled = true; return nil },
		nil,
		nil,
	)

	// Act
	screen.startCooldownAndProceed()

	// Assert
	if state.CooldownActive {
		t.Error("Expected cooldown NOT to be active with invalid data")
	}

	if nextCalled {
		t.Error("Expected onNext NOT to be called with invalid data")
	}
}

// TestTickerEntry_StrategyMetadataDisplay tests strategy metadata rendering
func TestTickerEntry_StrategyMetadataDisplay(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	screen.displayStrategyMetadata("Alt10")

	// Assert
	if len(screen.strategyInfo.Objects) == 0 {
		t.Error("Expected strategy metadata to be displayed")
	}
}

// TestTickerEntry_StrategyMetadataDisplay_Clear tests clearing metadata
func TestTickerEntry_StrategyMetadataDisplay_Clear(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{
		Sector: "Healthcare",
	}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Show metadata first
	screen.displayStrategyMetadata("Alt10")

	// Act - clear by selecting empty strategy
	screen.onStrategySelected("")

	// Assert
	if len(screen.strategyInfo.Objects) != 0 {
		t.Error("Expected strategy metadata to be cleared")
	}
}

// Helper function: contains checks if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper function: createTestPolicy creates a minimal policy for testing
func createTestPolicy() *models.Policy {
	return &models.Policy{
		PolicyID: "test-policy",
		Version:  "1.0.0",
		Defaults: models.PolicyDefaults{
			PortfolioHeatCap: 0.04,
			BucketHeatCap:    0.015,
			RiskPerTrade:     0.0075,
			CooldownSeconds:  120,
		},
		Sectors: []models.Sector{
			{
				Name:              "Healthcare",
				Priority:          1,
				Blocked:           false,
				Warning:           false,
				Notes:             "Best overall sector",
				AllowedStrategies: []string{"Alt10", "Alt46", "Alt43", "Alt39", "Alt28"},
			},
			{
				Name:              "Technology",
				Priority:          2,
				Blocked:           false,
				Warning:           false,
				Notes:             "Strong momentum sector",
				AllowedStrategies: []string{"Alt26", "Alt22", "Alt15", "Alt47", "Alt10"},
			},
		},
		Strategies: map[string]models.Strategy{
			"Alt10": {
				Label:              "Profit Targets",
				OptionsSuitability: "excellent",
				HoldWeeks:          "3-10",
				Notes:              "76.19% success rate",
			},
			"Alt46": {
				Label:              "Sector-Adaptive Parameters",
				OptionsSuitability: "excellent (healthcare)",
				HoldWeeks:          "3-12",
				Notes:              "Healthcare specialist",
			},
			"Alt43": {
				Label:              "Pyramiding",
				OptionsSuitability: "excellent",
				HoldWeeks:          "3-12",
				Notes:              "Add to winners",
			},
			"Alt39": {
				Label:              "Age-Based Targets",
				OptionsSuitability: "excellent",
				HoldWeeks:          "3-12",
				Notes:              "Time-adaptive exits",
			},
			"Alt28": {
				Label:              "ADX Filter",
				OptionsSuitability: "selective",
				HoldWeeks:          "variable",
				Notes:              "Diagnostic tool",
			},
			"Alt26": {
				Label:              "Profit Targets + Pyramiding",
				OptionsSuitability: "excellent",
				HoldWeeks:          "4-12",
				Notes:              "Combined strategy",
			},
			"Alt22": {
				Label:              "Parabolic SAR",
				OptionsSuitability: "excellent for tech",
				HoldWeeks:          "2-6",
				Notes:              "High churn",
			},
			"Alt15": {
				Label:              "Channel Breakouts",
				OptionsSuitability: "good",
				HoldWeeks:          "4-8",
				Notes:              "Breakout system",
			},
			"Alt47": {
				Label:              "Momentum-Scaled Sizing",
				OptionsSuitability: "stocks only",
				HoldWeeks:          "4-10",
				Notes:              "Ultra-low drawdowns",
			},
		},
	}
}
