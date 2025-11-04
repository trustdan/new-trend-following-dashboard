package screens

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// TestTradeEntry_NewTradeEntry tests screen initialization
func TestTradeEntry_NewTradeEntry(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()

	// Act
	screen := NewTradeEntry(state, window)

	// Assert
	if screen == nil {
		t.Fatal("NewTradeEntry returned nil")
	}
	if screen.state == nil {
		t.Error("state is nil")
	}
	if screen.strategySelect == nil {
		t.Error("strategySelect is nil")
	}
	if screen.strike1Entry == nil {
		t.Error("strike1Entry is nil")
	}
	if screen.expirationDate == nil {
		t.Error("expirationDate is nil")
	}
}

// TestTradeEntry_GetRequiredStrikes tests strike count determination
func TestTradeEntry_GetRequiredStrikes(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	tests := []struct {
		strategy      string
		expectedCount int
	}{
		// Single-leg strategies
		{"Long call", 1},
		{"Long put", 1},
		{"Covered call", 1},
		{"Cash-secured put", 1},

		// Two-leg strategies
		{"Bull call spread", 2},
		{"Bear put spread", 2},
		{"Bull put credit spread", 2},
		{"Bear call credit spread", 2},
		{"Straddle", 2},
		{"Strangle", 2},

		// Three-leg strategies
		{"Long put butterfly", 3},
		{"Long call butterfly", 3},
		{"Call ratio backspread", 3},
		{"Put ratio backspread", 3},
		{"Call broken wing", 3},
		{"Put broken wing", 3},

		// Four-leg strategies
		{"Iron butterfly", 4},
		{"Iron condor", 4},
		{"Inverse iron butterfly", 4},
		{"Inverse iron condor", 4},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.strategy, func(t *testing.T) {
			count := screen.getRequiredStrikes(tt.strategy)
			if count != tt.expectedCount {
				t.Errorf("getRequiredStrikes(%s) = %d, want %d", tt.strategy, count, tt.expectedCount)
			}
		})
	}
}

// TestTradeEntry_GetRequiredStrikes_UnknownStrategy tests default behavior
func TestTradeEntry_GetRequiredStrikes_UnknownStrategy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act
	count := screen.getRequiredStrikes("Unknown Strategy")

	// Assert (should default to 2-leg)
	if count != 2 {
		t.Errorf("getRequiredStrikes(unknown) = %d, want 2 (default)", count)
	}
}

// TestTradeEntry_Validate tests validation logic
func TestTradeEntry_Validate(t *testing.T) {
	tests := []struct {
		name           string
		setupTrade     func() *models.Trade
		expirationText string
		want           bool
	}{
		{
			name: "Valid trade with all fields",
			setupTrade: func() *models.Trade {
				return &models.Trade{
					Sector:          "Healthcare",
					Ticker:          "UNH",
					Strategy:        "Alt10",
					OptionsStrategy: "Bull call spread",
				}
			},
			expirationText: "45",
			want:           true,
		},
		{
			name: "Invalid - no options strategy",
			setupTrade: func() *models.Trade {
				return &models.Trade{
					Sector:   "Healthcare",
					Ticker:   "UNH",
					Strategy: "Alt10",
				}
			},
			expirationText: "45",
			want:           false,
		},
		{
			name: "Invalid - no expiration date",
			setupTrade: func() *models.Trade {
				return &models.Trade{
					Sector:          "Healthcare",
					Ticker:          "UNH",
					Strategy:        "Alt10",
					OptionsStrategy: "Bull call spread",
				}
			},
			expirationText: "",
			want:           false,
		},
		{
			name: "Invalid - nil trade",
			setupTrade: func() *models.Trade {
				return nil
			},
			expirationText: "45",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.CurrentTrade = tt.setupTrade()
			window := test.NewWindow(nil)
			defer window.Close()

			screen := NewTradeEntry(state, window)
			screen.expirationDate.SetText(tt.expirationText)

			// Act
			result := screen.Validate()

			// Assert
			if result != tt.want {
				t.Errorf("Validate() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestTradeEntry_GetName tests screen name
func TestTradeEntry_GetName(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act
	name := screen.GetName()

	// Assert
	if name != "trade_entry" {
		t.Errorf("GetName() = %s, want trade_entry", name)
	}
}

// TestTradeEntry_OnStrategySelected tests UI updates on strategy selection
func TestTradeEntry_OnStrategySelected(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.CurrentTrade = &models.Trade{}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act - select a 4-leg strategy
	screen.onStrategySelected("Iron condor")

	// Assert
	if state.CurrentTrade.OptionsStrategy != "Iron condor" {
		t.Errorf("OptionsStrategy not set in CurrentTrade")
	}

	// Strike container should have 8 elements (4 labels + 4 entry fields)
	if len(screen.strikeContainer.Objects) != 8 {
		t.Errorf("Strike container has %d objects, want 8 for 4-leg strategy", len(screen.strikeContainer.Objects))
	}
}

// TestTradeEntry_OnStrategySelected_TwoLeg tests 2-leg strategy
func TestTradeEntry_OnStrategySelected_TwoLeg(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.CurrentTrade = &models.Trade{}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act - select a 2-leg strategy
	screen.onStrategySelected("Bull call spread")

	// Assert
	// Strike container should have 4 elements (2 labels + 2 entry fields)
	if len(screen.strikeContainer.Objects) != 4 {
		t.Errorf("Strike container has %d objects, want 4 for 2-leg strategy", len(screen.strikeContainer.Objects))
	}
}

// TestTradeEntry_Render tests screen rendering
func TestTradeEntry_Render(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.CurrentTrade = &models.Trade{
		Ticker:   "UNH",
		Sector:   "Healthcare",
		Strategy: "Alt10",
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}
}

// TestTradeEntry_Render_WithExistingTrade tests pre-population
func TestTradeEntry_Render_WithExistingTrade(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	expiration := time.Now().AddDate(0, 0, 45)
	state.CurrentTrade = &models.Trade{
		Ticker:          "MSFT",
		Sector:          "Technology",
		Strategy:        "Alt26",
		OptionsStrategy: "Bull call spread",
		Strike1:         450.0,
		Strike2:         460.0,
		ExpirationDate:  expiration,
		Premium:         2.50,
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}

	// Verify strategy dropdown is pre-selected
	if screen.strategySelect.Selected != "Bull call spread" {
		t.Errorf("Strategy not pre-selected, got %s", screen.strategySelect.Selected)
	}

	// Verify strike fields are populated
	if screen.strike1Entry.Text != "450.00" {
		t.Errorf("Strike1 not pre-populated, got %s", screen.strike1Entry.Text)
	}
	if screen.strike2Entry.Text != "460.00" {
		t.Errorf("Strike2 not pre-populated, got %s", screen.strike2Entry.Text)
	}

	// Verify premium is populated
	if screen.premiumEntry.Text != "2.50" {
		t.Errorf("Premium not pre-populated, got %s", screen.premiumEntry.Text)
	}
}

// TestTradeEntry_AllStrategiesAvailable tests all 26 strategies are in dropdown
func TestTradeEntry_AllStrategiesAvailable(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewTradeEntry(state, window)

	expectedStrategies := []string{
		"Bull call spread",
		"Bear put spread",
		"Bull put credit spread",
		"Bear call credit spread",
		"Long call",
		"Long put",
		"Covered call",
		"Cash-secured put",
		"Iron butterfly",
		"Iron condor",
		"Long put butterfly",
		"Long call butterfly",
		"Calendar call spread",
		"Calendar put spread",
		"Diagonal call spread",
		"Diagonal put spread",
		"Inverse iron butterfly",
		"Inverse iron condor",
		"Short put butterfly",
		"Short call butterfly",
		"Straddle",
		"Strangle",
		"Call ratio backspread",
		"Put ratio backspread",
		"Call broken wing",
		"Put broken wing",
	}

	// Act
	availableStrategies := screen.strategySelect.Options

	// Assert
	if len(availableStrategies) != 26 {
		t.Errorf("Expected 26 strategies, got %d", len(availableStrategies))
	}

	// Verify all expected strategies are present
	for _, expected := range expectedStrategies {
		found := false
		for _, available := range availableStrategies {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Strategy %s not found in dropdown", expected)
		}
	}
}
