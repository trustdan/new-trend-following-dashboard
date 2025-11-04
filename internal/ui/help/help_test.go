package help

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestGetHelpForScreen_SectorSelection(t *testing.T) {
	// Act
	help := GetHelpForScreen("sector_selection")

	// Assert
	if help.Title == "" {
		t.Error("Title should not be empty")
	}
	if help.Description == "" {
		t.Error("Description should not be empty")
	}
	if len(help.Steps) == 0 {
		t.Error("Steps should not be empty")
	}
	if len(help.Tips) == 0 {
		t.Error("Tips should not be empty")
	}
}

func TestGetHelpForScreen_AllScreens(t *testing.T) {
	screens := []string{
		"sector_selection",
		"screener_launch",
		"ticker_entry",
		"checklist",
		"position_sizing",
		"heat_check",
		"trade_entry",
		"calendar",
		"trade_management",
	}

	for _, screen := range screens {
		t.Run(screen, func(t *testing.T) {
			help := GetHelpForScreen(screen)

			if help.Title == "" {
				t.Errorf("%s: Title should not be empty", screen)
			}
			if help.Description == "" {
				t.Errorf("%s: Description should not be empty", screen)
			}
			if len(help.Steps) == 0 {
				t.Errorf("%s: Steps should not be empty", screen)
			}
			if len(help.Tips) == 0 {
				t.Errorf("%s: Tips should not be empty", screen)
			}
		})
	}
}

func TestGetHelpForScreen_UnknownScreen(t *testing.T) {
	// Act
	help := GetHelpForScreen("unknown_screen")

	// Assert - should return generic help
	if help.Title == "" {
		t.Error("Generic help title should not be empty")
	}
	if help.Description == "" {
		t.Error("Generic help description should not be empty")
	}
	if len(help.Steps) == 0 {
		t.Error("Generic help steps should not be empty")
	}
	if len(help.Tips) == 0 {
		t.Error("Generic help tips should not be empty")
	}
}

func TestGetHelpForScreen_TickerEntry_HasCooldownInfo(t *testing.T) {
	// Act
	help := GetHelpForScreen("ticker_entry")

	// Assert - should mention cooldown
	found := false
	allText := help.Title + help.Description
	for _, step := range help.Steps {
		allText += step
	}
	for _, tip := range help.Tips {
		allText += tip
		if contains(tip, "Cooldown") || contains(tip, "cooldown") {
			found = true
		}
	}

	if !found {
		t.Error("Ticker entry help should mention cooldown timer")
	}
}

func TestGetHelpForScreen_HeatCheck_HasLimits(t *testing.T) {
	// Act
	help := GetHelpForScreen("heat_check")

	// Assert - should mention 4% and 1.5% limits
	allText := help.Description
	for _, step := range help.Steps {
		allText += step
	}

	if !contains(allText, "4%") {
		t.Error("Heat check help should mention 4% portfolio limit")
	}
	if !contains(allText, "1.5%") {
		t.Error("Heat check help should mention 1.5% sector limit")
	}
}

func TestGetHelpForScreen_PositionSizing_HasConvictionScale(t *testing.T) {
	// Act
	help := GetHelpForScreen("position_sizing")

	// Assert - should mention conviction scale 5-8
	allText := help.Description
	for _, step := range help.Steps {
		allText += step
	}

	if !contains(allText, "5-8") {
		t.Error("Position sizing help should mention 5-8 conviction scale")
	}
}

func TestShowHelpDialog_DoesNotPanic(t *testing.T) {
	// Arrange
	window := test.NewWindow(nil)

	// Act & Assert - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowHelpDialog panicked: %v", r)
		}
	}()

	ShowHelpDialog("sector_selection", window)
}

func TestShowWelcomeScreen_DoesNotPanic(t *testing.T) {
	// Arrange
	window := test.NewWindow(nil)
	callback := func(dontShowAgain bool) {
		// Callback for testing
	}

	// Act & Assert - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowWelcomeScreen panicked: %v", r)
		}
	}()

	ShowWelcomeScreen(window, callback)
}

func TestGetHelpForScreen_Checklist_Has5RequiredItems(t *testing.T) {
	// Act
	help := GetHelpForScreen("checklist")

	// Assert
	allText := help.Description
	for _, step := range help.Steps {
		allText += step
	}

	if !contains(allText, "5 REQUIRED") {
		t.Error("Checklist help should mention 5 required items")
	}
}

func TestGetHelpForScreen_Calendar_HasColorCoding(t *testing.T) {
	// Act
	help := GetHelpForScreen("calendar")

	// Assert - should explain color coding
	allText := help.Description
	for _, tip := range help.Tips {
		allText += tip
	}

	colors := []string{"Green", "Red", "Yellow", "Blue"}
	for _, color := range colors {
		if !contains(allText, color) {
			t.Errorf("Calendar help should mention %s color coding", color)
		}
	}
}

func TestGetHelpForScreen_TradeManagement_MentionsFeatureFlag(t *testing.T) {
	// Act
	help := GetHelpForScreen("trade_management")

	// Assert
	allText := help.Title + help.Description
	for _, tip := range help.Tips {
		allText += tip
	}

	if !contains(allText, "feature") || !contains(allText, "flag") {
		t.Error("Trade management help should mention feature flag requirement")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
