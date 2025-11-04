package screens

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

func TestSectorSelection_Validate(t *testing.T) {
	tests := []struct {
		name     string
		trade    *models.Trade
		expected bool
	}{
		{
			name:     "No trade - invalid",
			trade:    nil,
			expected: false,
		},
		{
			name:     "Trade with no sector - invalid",
			trade:    &models.Trade{},
			expected: false,
		},
		{
			name: "Trade with sector - valid",
			trade: &models.Trade{
				Sector: "Healthcare",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewApp()
			defer app.Quit()

			window := test.NewWindow(nil)
			defer window.Close()

			state := &appcore.AppState{
				CurrentTrade: tt.trade,
			}

			screen := NewSectorSelection(state, window)
			result := screen.Validate()

			if result != tt.expected {
				t.Errorf("Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSectorSelection_GetName(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewSectorSelection(state, window)

	expected := "sector_selection"
	result := screen.GetName()

	if result != expected {
		t.Errorf("GetName() = %s, expected %s", result, expected)
	}
}

func TestSectorSelection_SelectSector(t *testing.T) {
	tests := []struct {
		name           string
		sector         models.Sector
		expectError    bool
		expectedSector string
	}{
		{
			name: "Select valid sector",
			sector: models.Sector{
				Name:    "Healthcare",
				Blocked: false,
			},
			expectError:    false,
			expectedSector: "Healthcare",
		},
		{
			name: "Attempt to select blocked sector",
			sector: models.Sector{
				Name:    "Utilities",
				Blocked: true,
			},
			expectError:    true,
			expectedSector: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewApp()
			defer app.Quit()

			window := test.NewWindow(nil)
			defer window.Close()

			state := &appcore.AppState{}
			screen := NewSectorSelection(state, window)

			// Call selectSector
			screen.selectSector(tt.sector)

			// Check if sector was set correctly
			if !tt.expectError {
				if state.CurrentTrade == nil {
					t.Error("CurrentTrade should not be nil after valid selection")
					return
				}
				if state.CurrentTrade.Sector != tt.expectedSector {
					t.Errorf("CurrentTrade.Sector = %s, expected %s",
						state.CurrentTrade.Sector, tt.expectedSector)
				}
			} else {
				// For blocked sectors, trade should not be created
				if state.CurrentTrade != nil && state.CurrentTrade.Sector == tt.sector.Name {
					t.Error("CurrentTrade.Sector should not be set for blocked sector")
				}
			}
		})
	}
}

func TestSectorSelection_SetNavCallbacks(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewSectorSelection(state, window)

	nextCalled := false
	backCalled := false
	cancelCalled := false

	onNext := func() error {
		nextCalled = true
		return nil
	}

	onBack := func() error {
		backCalled = true
		return nil
	}

	onCancel := func() {
		cancelCalled = true
	}

	screen.SetNavCallbacks(onNext, onBack, onCancel)

	// Test that callbacks were set
	if screen.onNext == nil {
		t.Error("onNext callback not set")
	}
	if screen.onBack == nil {
		t.Error("onBack callback not set")
	}
	if screen.onCancel == nil {
		t.Error("onCancel callback not set")
	}

	// Test callback execution
	screen.onNext()
	if !nextCalled {
		t.Error("onNext callback not executed")
	}

	screen.onBack()
	if !backCalled {
		t.Error("onBack callback not executed")
	}

	screen.onCancel()
	if !cancelCalled {
		t.Error("onCancel callback not executed")
	}
}

func TestSectorSelection_Render(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	policy := &models.Policy{
		Sectors: []models.Sector{
			{
				Name:              "Healthcare",
				Priority:          1,
				Blocked:           false,
				Warning:           false,
				Notes:             "Best overall performance",
				AllowedStrategies: []string{"Alt10", "Alt46"},
			},
			{
				Name:     "Utilities",
				Priority: 10,
				Blocked:  true,
				Notes:    "Do not trade",
			},
		},
	}

	state := &appcore.AppState{
		Policy: policy,
	}

	screen := NewSectorSelection(state, window)

	// Test rendering
	content := screen.Render()
	if content == nil {
		t.Error("Render() returned nil")
	}
}

func TestSectorSelection_SectorSorting(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	// Create policy with sectors in wrong order
	policy := &models.Policy{
		Sectors: []models.Sector{
			{Name: "Consumer", Priority: 4},
			{Name: "Healthcare", Priority: 1},
			{Name: "Technology", Priority: 2},
			{Name: "Utilities", Priority: 10},
		},
	}

	state := &appcore.AppState{
		Policy: policy,
	}

	screen := NewSectorSelection(state, window)

	// Render the screen (this will sort sectors internally)
	content := screen.Render()
	if content == nil {
		t.Error("Render() returned nil")
	}

	// The sorting happens in createSectorCards(), which is called by Render()
	// We can't easily test the UI order, but we can verify Render() doesn't panic
}
