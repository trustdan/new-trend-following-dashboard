package screens

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// ============================================================================
// PHASE 6: Warning-Based Trading System Tests
// ============================================================================

// TestTickerEntry_GetSuitability_GreenStrategy tests green (excellent/good) strategy rating
func TestTickerEntry_GetSuitability_GreenStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	suitability := screen.getSuitability("Alt10", "Healthcare")

	// Assert
	if suitability.Rating != "excellent" {
		t.Errorf("Expected rating 'excellent', got '%s'", suitability.Rating)
	}
	if suitability.Color != "green" {
		t.Errorf("Expected color 'green', got '%s'", suitability.Color)
	}
	if suitability.RequireAcknowledgement {
		t.Error("Expected RequireAcknowledgement to be false for green strategy")
	}
}

// TestTickerEntry_GetSuitability_YellowStrategy tests yellow (marginal) strategy rating
func TestTickerEntry_GetSuitability_YellowStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	suitability := screen.getSuitability("Alt26", "Healthcare")

	// Assert
	if suitability.Rating != "marginal" {
		t.Errorf("Expected rating 'marginal', got '%s'", suitability.Rating)
	}
	if suitability.Color != "yellow" {
		t.Errorf("Expected color 'yellow', got '%s'", suitability.Color)
	}
	if !suitability.RequireAcknowledgement {
		t.Error("Expected RequireAcknowledgement to be true for yellow strategy")
	}
}

// TestTickerEntry_GetSuitability_RedStrategy tests red (incompatible) strategy rating
func TestTickerEntry_GetSuitability_RedStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	suitability := screen.getSuitability("Alt22", "Healthcare")

	// Assert
	if suitability.Rating != "incompatible" {
		t.Errorf("Expected rating 'incompatible', got '%s'", suitability.Rating)
	}
	if suitability.Color != "red" {
		t.Errorf("Expected color 'red', got '%s'", suitability.Color)
	}
	if !suitability.RequireAcknowledgement {
		t.Error("Expected RequireAcknowledgement to be true for red strategy")
	}
}

// TestTickerEntry_GetSuitability_NotFound tests default fallback for missing rating
func TestTickerEntry_GetSuitability_NotFound(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act - request strategy not in suitability map
	suitability := screen.getSuitability("Alt999", "Healthcare")

	// Assert - should default to marginal
	if suitability.Rating != "marginal" {
		t.Errorf("Expected default rating 'marginal', got '%s'", suitability.Rating)
	}
	if suitability.Color != "yellow" {
		t.Errorf("Expected default color 'yellow', got '%s'", suitability.Color)
	}
	if !suitability.RequireAcknowledgement {
		t.Error("Expected default RequireAcknowledgement to be true")
	}
}

// TestTickerEntry_GetColorIndicator tests color indicator mapping
func TestTickerEntry_GetColorIndicator(t *testing.T) {
	tests := []struct {
		color    string
		expected string
	}{
		{"green", "[GREEN]"},
		{"yellow", "[YELLOW]"},
		{"red", "[RED]"},
		{"unknown", "[UNKNOWN]"},
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			indicator := getColorIndicatorText(tt.color)
			if indicator != tt.expected {
				t.Errorf("getColorIndicatorText(%s) = %s, want %s", tt.color, indicator, tt.expected)
			}
		})
	}
}

// TestTickerEntry_GetAllStrategiesWithIndicators tests ALL strategies are shown
func TestTickerEntry_GetAllStrategiesWithIndicators(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)

	// Act
	strategies := screen.getAllStrategiesWithIndicators()

	// Assert
	expectedCount := 5
	if len(strategies) != expectedCount {
		t.Fatalf("Expected %d strategies, got %d", expectedCount, len(strategies))
	}

	var hasGreen, hasYellow, hasRed bool
	foundAlt10 := false
	foundAlt26 := false
	foundAlt22 := false

	for _, s := range strategies {
		if strings.HasPrefix(s, "[GREEN]") {
			hasGreen = true
		}
		if strings.HasPrefix(s, "[YELLOW]") {
			hasYellow = true
		}
		if strings.HasPrefix(s, "[RED]") {
			hasRed = true
		}
		if strings.Contains(s, "Alt10") {
			foundAlt10 = true
		}
		if strings.Contains(s, "Alt26") {
			foundAlt26 = true
		}
		if strings.Contains(s, "Alt22") {
			foundAlt22 = true
		}
	}

	if !hasGreen {
		t.Error("Expected at least one [GREEN] strategy")
	}
	if !hasYellow {
		t.Error("Expected at least one [YELLOW] strategy")
	}
	if !hasRed {
		t.Error("Expected at least one [RED] strategy")
	}
	if !foundAlt10 {
		t.Error("Expected Alt10 to be included in dropdown")
	}
	if !foundAlt26 {
		t.Error("Expected Alt26 to be included in dropdown")
	}
	if !foundAlt22 {
		t.Error("Expected Alt22 to be included in dropdown")
	}
}

// TestTickerEntry_OnStrategySelected_GreenStrategy tests green strategy selection (no warning)
func TestTickerEntry_OnStrategySelected_GreenStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare", Ticker: "UNH"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	strategies := screen.getAllStrategiesWithIndicators()
	alt10 := findStrategyLabel(t, strategies, "Alt10")
	screen.onStrategySelected(alt10)

	// Assert
	if state.CurrentTrade.Strategy != "Alt10" {
		t.Errorf("Expected strategy 'Alt10', got '%s'", state.CurrentTrade.Strategy)
	}
	if state.CurrentTrade.StrategySuitability != "excellent" {
		t.Errorf("Expected suitability 'excellent', got '%s'", state.CurrentTrade.StrategySuitability)
	}

	// Green strategy - Continue button should be enabled (no acknowledgement required)
	if screen.continueBtn.Disabled() {
		t.Error("Expected continue button to be enabled for green strategy")
	}

	// Warning banner should be hidden
	if len(screen.warningBanner.Objects) != 0 {
		t.Error("Expected warning banner to be hidden for green strategy")
	}

	// Acknowledgement checkbox should be hidden
	if screen.ackCheckbox.Visible() {
		t.Error("Expected acknowledgement checkbox to be hidden for green strategy")
	}
}

// TestTickerEntry_OnStrategySelected_YellowStrategy tests yellow strategy selection (warning + ack)
func TestTickerEntry_OnStrategySelected_YellowStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare", Ticker: "UNH"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	strategies := screen.getAllStrategiesWithIndicators()
	alt26 := findStrategyLabel(t, strategies, "Alt26")
	screen.onStrategySelected(alt26)

	// Assert
	if state.CurrentTrade.Strategy != "Alt26" {
		t.Errorf("Expected strategy 'Alt26', got '%s'", state.CurrentTrade.Strategy)
	}
	if state.CurrentTrade.StrategySuitability != "marginal" {
		t.Errorf("Expected suitability 'marginal', got '%s'", state.CurrentTrade.StrategySuitability)
	}

	// Yellow strategy - Continue button should be DISABLED until acknowledged
	if !screen.continueBtn.Disabled() {
		t.Error("Expected continue button to be disabled for yellow strategy (not acknowledged)")
	}

	// Warning banner should be visible
	if len(screen.warningBanner.Objects) == 0 {
		t.Error("Expected warning banner to be visible for yellow strategy")
	}

	// Acknowledgement checkbox should be visible and unchecked
	if !screen.ackCheckbox.Visible() {
		t.Error("Expected acknowledgement checkbox to be visible for yellow strategy")
	}
	if screen.ackCheckbox.Checked {
		t.Error("Expected acknowledgement checkbox to be unchecked initially")
	}
}

// TestTickerEntry_OnStrategySelected_RedStrategy tests red strategy selection (strong warning + ack)
func TestTickerEntry_OnStrategySelected_RedStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare", Ticker: "UNH"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Act
	strategies := screen.getAllStrategiesWithIndicators()
	alt22 := findStrategyLabel(t, strategies, "Alt22")
	screen.onStrategySelected(alt22)

	// Assert
	if state.CurrentTrade.Strategy != "Alt22" {
		t.Errorf("Expected strategy 'Alt22', got '%s'", state.CurrentTrade.Strategy)
	}
	if state.CurrentTrade.StrategySuitability != "incompatible" {
		t.Errorf("Expected suitability 'incompatible', got '%s'", state.CurrentTrade.StrategySuitability)
	}

	// Red strategy - Continue button should be DISABLED until acknowledged
	if !screen.continueBtn.Disabled() {
		t.Error("Expected continue button to be disabled for red strategy (not acknowledged)")
	}

	// Warning banner should be visible
	if len(screen.warningBanner.Objects) == 0 {
		t.Error("Expected warning banner to be visible for red strategy")
	}

	// Acknowledgement checkbox should be visible
	if !screen.ackCheckbox.Visible() {
		t.Error("Expected acknowledgement checkbox to be visible for red strategy")
	}
}

// TestTickerEntry_AcknowledgementCheckbox_EnablesContinue tests checkbox enables Continue button
func TestTickerEntry_AcknowledgementCheckbox_EnablesContinue(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare", Ticker: "UNH"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Select yellow strategy (requires acknowledgement)
	strategies := screen.getAllStrategiesWithIndicators()
	alt26 := findStrategyLabel(t, strategies, "Alt26")
	screen.onStrategySelected(alt26)

	// Verify Continue button is initially disabled
	if !screen.continueBtn.Disabled() {
		t.Fatal("Expected continue button to be disabled before acknowledgement")
	}

	// Act - check acknowledgement checkbox
	screen.ackCheckbox.SetChecked(true)
	screen.updateContinueButton()

	// Assert - Continue button should now be enabled
	if screen.continueBtn.Disabled() {
		t.Error("Expected continue button to be enabled after acknowledgement")
	}
}

// TestTickerEntry_StartCooldown_SetsAcknowledgementFlag tests acknowledgement flag is set on proceed
func TestTickerEntry_StartCooldown_SetsAcknowledgementFlag(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = createTestPolicy()
	state.CurrentTrade = &models.Trade{Sector: "Healthcare", Ticker: "UNH"}

	window := test.NewWindow(nil)
	defer window.Close()

	screen := NewTickerEntry(state, window)
	screen.Render() // Initialize UI components

	// Select yellow strategy and acknowledge
	strategies := screen.getAllStrategiesWithIndicators()
	alt26 := findStrategyLabel(t, strategies, "Alt26")
	screen.onStrategySelected(alt26)
	screen.ackCheckbox.SetChecked(true)
	screen.updateContinueButton()

	nextCalled := false
	screen.SetNavCallbacks(
		func() error { nextCalled = true; return nil },
		nil,
		nil,
	)

	// Act
	screen.startCooldownAndProceed()

	// Assert
	if !state.CurrentTrade.StrategyWarningAcknowledged {
		t.Error("Expected StrategyWarningAcknowledged to be true after proceeding")
	}
	if !nextCalled {
		t.Error("Expected onNext to be called")
	}
}

// TestTickerEntry_UpdateContinueButton_Phase6Logic tests Phase 6 button enable/disable logic
func TestTickerEntry_UpdateContinueButton_Phase6Logic(t *testing.T) {
	tests := []struct {
		name            string
		sector          string
		ticker          string
		strategy        string
		suitability     string
		acknowledged    bool
		expectedEnabled bool
	}{
		{
			name:            "Green strategy + valid data = enabled",
			sector:          "Healthcare",
			ticker:          "UNH",
			strategy:        "Alt10",
			suitability:     "excellent",
			acknowledged:    false, // Not required for green
			expectedEnabled: true,
		},
		{
			name:            "Yellow strategy + acknowledged = enabled",
			sector:          "Healthcare",
			ticker:          "UNH",
			strategy:        "Alt26",
			suitability:     "marginal",
			acknowledged:    true,
			expectedEnabled: true,
		},
		{
			name:            "Yellow strategy + NOT acknowledged = disabled",
			sector:          "Healthcare",
			ticker:          "UNH",
			strategy:        "Alt26",
			suitability:     "marginal",
			acknowledged:    false,
			expectedEnabled: false,
		},
		{
			name:            "Red strategy + acknowledged = enabled",
			sector:          "Healthcare",
			ticker:          "UNH",
			strategy:        "Alt22",
			suitability:     "incompatible",
			acknowledged:    true,
			expectedEnabled: true,
		},
		{
			name:            "Red strategy + NOT acknowledged = disabled",
			sector:          "Healthcare",
			ticker:          "UNH",
			strategy:        "Alt22",
			suitability:     "incompatible",
			acknowledged:    false,
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.Policy = createTestPolicy()
			state.CurrentTrade = &models.Trade{
				Sector:   tt.sector,
				Ticker:   tt.ticker,
				Strategy: tt.strategy,
			}

			window := test.NewWindow(nil)
			defer window.Close()

			screen := NewTickerEntry(state, window)
			screen.Render() // Initialize UI components

			// Set acknowledgement checkbox state
			screen.ackCheckbox.SetChecked(tt.acknowledged)

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

func findStrategyLabel(t *testing.T, options []string, strategyID string) string {
	t.Helper()
	for _, option := range options {
		if strings.Contains(option, strategyID) {
			return option
		}
	}
	t.Fatalf("strategy %s not found in options: %v", strategyID, options)
	return ""
}
