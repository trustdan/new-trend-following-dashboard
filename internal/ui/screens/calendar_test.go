package screens

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"image/color"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// TestCalendar_NewCalendar tests screen initialization
func TestCalendar_NewCalendar(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()

	// Act
	screen := NewCalendar(state, window)

	// Assert
	if screen == nil {
		t.Fatal("NewCalendar returned nil")
	}
	if screen.state == nil {
		t.Error("state is nil")
	}
	if screen.window == nil {
		t.Error("window is nil")
	}
}

// TestCalendar_GetSectors tests sector loading
func TestCalendar_GetSectors(t *testing.T) {
	tests := []struct {
		name          string
		setupPolicy   func() *models.Policy
		expectedCount int
		shouldContain string
	}{
		{
			name: "Policy with sectors",
			setupPolicy: func() *models.Policy {
				return &models.Policy{
					Sectors: []models.Sector{
						{Name: "Healthcare"},
						{Name: "Technology"},
						{Name: "Industrials"},
					},
				}
			},
			expectedCount: 3,
			shouldContain: "Healthcare",
		},
		{
			name: "No policy - uses defaults",
			setupPolicy: func() *models.Policy {
				return nil
			},
			expectedCount: 8, // Default sectors
			shouldContain: "Healthcare",
		},
		{
			name: "Policy with no sectors - uses defaults",
			setupPolicy: func() *models.Policy {
				return &models.Policy{
					Sectors: []models.Sector{},
				}
			},
			expectedCount: 8,
			shouldContain: "Healthcare",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			state.Policy = tt.setupPolicy()
			window := test.NewWindow(nil)
			defer window.Close()
			screen := NewCalendar(state, window)

			// Act
			sectors := screen.getSectors()

			// Assert
			if len(sectors) != tt.expectedCount {
				t.Errorf("getSectors() returned %d sectors, want %d", len(sectors), tt.expectedCount)
			}

			// Check if expected sector is present
			found := false
			for _, sector := range sectors {
				if sector == tt.shouldContain {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected sector %s not found in results", tt.shouldContain)
			}
		})
	}
}

// TestCalendar_GetTradesForSector tests trade filtering
func TestCalendar_GetTradesForSector(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.AllTrades = []models.Trade{
		{Ticker: "UNH", Sector: "Healthcare"},
		{Ticker: "JNJ", Sector: "Healthcare"},
		{Ticker: "MSFT", Sector: "Technology"},
		{Ticker: "AAPL", Sector: "Technology"},
		{Ticker: "CAT", Sector: "Industrials"},
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	tests := []struct {
		sector        string
		expectedCount int
	}{
		{"Healthcare", 2},
		{"Technology", 2},
		{"Industrials", 1},
		{"Consumer", 0},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.sector, func(t *testing.T) {
			trades := screen.getTradesForSector(tt.sector)
			if len(trades) != tt.expectedCount {
				t.Errorf("getTradesForSector(%s) returned %d trades, want %d", tt.sector, len(trades), tt.expectedCount)
			}
		})
	}
}

// TestCalendar_CountActiveTrades tests active trade counting
func TestCalendar_CountActiveTrades(t *testing.T) {
	// Arrange
	now := time.Now()
	state := appcore.NewAppState()
	state.AllTrades = []models.Trade{
		{Ticker: "UNH", ExpirationDate: now.AddDate(0, 0, 30)},  // Active (30 days out)
		{Ticker: "JNJ", ExpirationDate: now.AddDate(0, 0, 45)},  // Active (45 days out)
		{Ticker: "MSFT", ExpirationDate: now.AddDate(0, 0, -5)}, // Expired (5 days ago)
		{Ticker: "AAPL", ExpirationDate: now.AddDate(0, 0, 10)}, // Active (10 days out)
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	count := screen.countActiveTrades()

	// Assert
	if count != 3 {
		t.Errorf("countActiveTrades() = %d, want 3 (3 active, 1 expired)", count)
	}
}

// TestCalendar_CalculateTotalRisk tests risk summation
func TestCalendar_CalculateTotalRisk(t *testing.T) {
	// Arrange
	now := time.Now()
	state := appcore.NewAppState()
	state.AllTrades = []models.Trade{
		{Ticker: "UNH", ExpirationDate: now.AddDate(0, 0, 30), MaxLoss: 500.0},  // Active
		{Ticker: "JNJ", ExpirationDate: now.AddDate(0, 0, 45), MaxLoss: 300.0},  // Active
		{Ticker: "MSFT", ExpirationDate: now.AddDate(0, 0, -5), MaxLoss: 200.0}, // Expired - should not count
		{Ticker: "AAPL", ExpirationDate: now.AddDate(0, 0, 10), MaxLoss: 400.0}, // Active
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	totalRisk := screen.calculateTotalRisk()

	// Assert
	expectedRisk := 500.0 + 300.0 + 400.0 // Only active trades
	if totalRisk != expectedRisk {
		t.Errorf("calculateTotalRisk() = %.2f, want %.2f", totalRisk, expectedRisk)
	}
}

// TestCalendar_GetTradeBarColor tests color determination logic
func TestCalendar_GetTradeBarColor(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		trade         models.Trade
		expectedColor color.Color
		description   string
	}{
		{
			name: "Expiring soon (yellow)",
			trade: models.Trade{
				ExpirationDate: now.AddDate(0, 0, 5), // 5 days out
			},
			expectedColor: color.RGBA{R: 255, G: 193, B: 7, A: 255},
			description:   "Yellow for expiring within 7 days",
		},
		{
			name: "Expired (red)",
			trade: models.Trade{
				ExpirationDate: now.AddDate(0, 0, -10), // 10 days ago
			},
			expectedColor: color.RGBA{R: 220, G: 53, B: 69, A: 255},
			description:   "Red for expired trades",
		},
		{
			name: "Profitable (green)",
			trade: models.Trade{
				ExpirationDate: now.AddDate(0, 0, 30),
				ProfitLoss:     floatPtr(150.0), // Profitable
			},
			expectedColor: color.RGBA{R: 40, G: 167, B: 69, A: 255},
			description:   "Green for profitable trades",
		},
		{
			name: "Losing (red)",
			trade: models.Trade{
				ExpirationDate: now.AddDate(0, 0, 30),
				ProfitLoss:     floatPtr(-50.0), // Losing
			},
			expectedColor: color.RGBA{R: 220, G: 53, B: 69, A: 255},
			description:   "Red for losing trades",
		},
		{
			name: "Active default (blue)",
			trade: models.Trade{
				ExpirationDate: now.AddDate(0, 0, 30),
				ProfitLoss:     nil, // No P&L data
			},
			expectedColor: color.RGBA{R: 13, G: 110, B: 253, A: 255},
			description:   "Blue for active trades with unknown P&L",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			state := appcore.NewAppState()
			window := test.NewWindow(nil)
			defer window.Close()
			screen := NewCalendar(state, window)

			// Act
			resultColor := screen.getTradeBarColor(tt.trade)

			// Assert
			if !colorsEqual(resultColor, tt.expectedColor) {
				t.Errorf("getTradeBarColor() for %s returned wrong color: got %v, want %v",
					tt.description, resultColor, tt.expectedColor)
			}
		})
	}
}

// TestCalendar_Validate tests validation always passes
func TestCalendar_Validate(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	result := screen.Validate()

	// Assert
	if !result {
		t.Error("Calendar Validate() should always return true (display-only screen)")
	}
}

// TestCalendar_GetName tests screen name
func TestCalendar_GetName(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	name := screen.GetName()

	// Assert
	if name != "calendar" {
		t.Errorf("GetName() = %s, want calendar", name)
	}
}

// TestCalendar_Render tests basic rendering
func TestCalendar_Render(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.AllTrades = []models.Trade{
		{
			Ticker:         "UNH",
			Sector:         "Healthcare",
			ExpirationDate: time.Now().AddDate(0, 0, 45),
			CreatedAt:      time.Now().AddDate(0, 0, -5),
			MaxLoss:        500.0,
		},
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}
}

// TestCalendar_Render_WithPolicy tests rendering with policy configuration
func TestCalendar_Render_WithPolicy(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.Policy = &models.Policy{
		Sectors: []models.Sector{
			{Name: "Healthcare"},
			{Name: "Technology"},
		},
		Calendar: models.CalendarConfig{
			PastDays:   21,
			FutureDays: 90,
		},
	}
	state.AllTrades = []models.Trade{
		{
			Ticker:         "UNH",
			Sector:         "Healthcare",
			ExpirationDate: time.Now().AddDate(0, 0, 45),
			CreatedAt:      time.Now().AddDate(0, 0, -5),
			MaxLoss:        500.0,
		},
	}
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil")
	}
}

// TestCalendar_Render_NoTrades tests rendering with empty trade list
func TestCalendar_Render_NoTrades(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	state.AllTrades = []models.Trade{} // Empty
	window := test.NewWindow(nil)
	defer window.Close()
	screen := NewCalendar(state, window)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render() returned nil even with no trades")
	}
}

// Helper function to create float pointer
func floatPtr(f float64) *float64 {
	return &f
}

// Helper function to compare colors
func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
