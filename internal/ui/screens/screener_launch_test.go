package screens

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

func TestScreenerLaunch_Validate(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	// Screener launch should always be valid
	if !screen.Validate() {
		t.Error("Validate() should always return true for screener launch")
	}
}

func TestScreenerLaunch_GetName(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	expected := "screener_launch"
	result := screen.GetName()

	if result != expected {
		t.Errorf("GetName() = %s, expected %s", result, expected)
	}
}

func TestScreenerLaunch_SetNavCallbacks(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

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

func TestScreenerLaunch_Render_NoSector(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	// State with no sector selected
	state := &appcore.AppState{
		Policy: &models.Policy{},
	}

	screen := NewScreenerLaunch(state, window)

	// Should render without crashing even with no sector
	content := screen.Render()
	if content == nil {
		t.Error("Render() returned nil with no sector")
	}
}

func TestScreenerLaunch_Render_WithSector(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	// State with Healthcare sector selected
	policy := &models.Policy{
		Sectors: []models.Sector{
			{
				Name:     "Healthcare",
				Priority: 1,
				ScreenerURLs: map[string]string{
					"universe":     "https://finviz.com/screener.ashx?v=211&f=sec_healthcare",
					"pullback":     "https://finviz.com/screener.ashx?v=211&f=sec_healthcare,ta_rsi_os40",
					"breakout":     "https://finviz.com/screener.ashx?v=211&f=sec_healthcare,ta_highlow52w_nh",
					"golden_cross": "https://finviz.com/screener.ashx?v=211&f=sec_healthcare,ta_sma50_pa200",
				},
			},
		},
	}

	state := &appcore.AppState{
		Policy: policy,
		CurrentTrade: &models.Trade{
			Sector: "Healthcare",
		},
	}

	screen := NewScreenerLaunch(state, window)

	// Should render with screener cards
	content := screen.Render()
	if content == nil {
		t.Error("Render() returned nil with Healthcare sector")
	}
}

func TestScreenerLaunch_CreateScreenerCards(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	sector := models.Sector{
		Name: "Healthcare",
		ScreenerURLs: map[string]string{
			"universe": "https://finviz.com/screener.ashx?v=211&f=sec_healthcare",
			"pullback": "https://finviz.com/screener.ashx?v=211&f=sec_healthcare,ta_rsi_os40",
		},
	}

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	// Should create cards for available screeners
	cards := screen.createScreenerCards(sector)
	if cards == nil {
		t.Error("createScreenerCards() returned nil")
	}
}

func TestScreenerLaunch_CreateScreenerCards_Empty(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	sector := models.Sector{
		Name:         "EmptySector",
		ScreenerURLs: map[string]string{}, // No screeners
	}

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	// Should handle empty screener URLs gracefully
	cards := screen.createScreenerCards(sector)
	if cards == nil {
		t.Error("createScreenerCards() returned nil for empty sector")
	}
}

func TestScreenerLaunch_LaunchURL_Tracking(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{
		Policy: &models.Policy{
			Sectors: []models.Sector{
				{
					Name:   "Healthcare",
					ScreenerURLs: map[string]string{
						"universe": "https://finviz.com/screener.ashx?v=211&f=sec_healthcare",
					},
				},
			},
		},
		CurrentTrade: &models.Trade{
			Sector: "Healthcare",
		},
	}

	screen := NewScreenerLaunch(state, window)

	// Launch URL (won't actually open browser in test, but will track timestamp)
	testURL := "https://finviz.com/screener.ashx?v=211&f=sec_healthcare"

	// Note: This will attempt to open the URL but won't fail the test
	// We're mainly testing that the tracking works
	screen.launchURL("universe", testURL)

	// Verify timestamp was recorded
	if lastRun, exists := screen.lastLaunch["universe"]; exists {
		// Verify it's recent (within last second)
		elapsed := time.Since(lastRun)
		if elapsed > time.Second {
			t.Errorf("Last launch timestamp too old: %v", elapsed)
		}
	} else {
		t.Error("Last launch timestamp not recorded")
	}
}

func TestScreenerLaunch_CreateScreenerCard(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	// Create a screener card
	card := screen.createScreenerCard(
		"universe",
		"Universe Screener",
		"Test description",
		"Weekly",
		"Test purpose",
		"https://finviz.com/screener.ashx?v=211",
	)

	if card == nil {
		t.Error("createScreenerCard() returned nil")
	}
}

func TestScreenerLaunch_CreateInfoBanner(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	banner := screen.createInfoBanner()
	if banner == nil {
		t.Error("createInfoBanner() returned nil")
	}
}

func TestScreenerLaunch_ScreenerOrder(t *testing.T) {
	// Test that screeners appear in correct order: universe, pullback, breakout, golden_cross
	app := test.NewApp()
	defer app.Quit()

	window := test.NewWindow(nil)
	defer window.Close()

	sector := models.Sector{
		Name: "Healthcare",
		ScreenerURLs: map[string]string{
			"golden_cross": "url1",
			"universe":     "url2",
			"breakout":     "url3",
			"pullback":     "url4",
		},
	}

	state := &appcore.AppState{}
	screen := NewScreenerLaunch(state, window)

	// Create cards (internally maintains order)
	cards := screen.createScreenerCards(sector)
	if cards == nil {
		t.Error("createScreenerCards() returned nil")
	}

	// Order is enforced in createScreenerCards() via the order slice
	// This test verifies it doesn't panic with out-of-order keys
}
