package ui

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"tf-engine/internal/appcore"
	"tf-engine/internal/storage"
	"tf-engine/internal/ui/help"
	"tf-engine/internal/ui/screens"
)

// Screen represents a single screen in the workflow
type Screen interface {
	Render() fyne.CanvasObject
	Validate() bool
	GetName() string
}

// Navigator manages screen transitions and workflow state
type Navigator struct {
	screens      []Screen
	currentIndex int
	history      []int
	state        *appcore.AppState
	window       fyne.Window
}

// NewNavigator creates a new navigator with all 9 screens (8 workflow + 1 management)
func NewNavigator(state *appcore.AppState, window fyne.Window) *Navigator {
	nav := &Navigator{
		currentIndex: -1, // -1 = dashboard
		history:      []int{},
		state:        state,
		window:       window,
	}

	// Initialize all screens (8 workflow screens + trade management)
	nav.screens = []Screen{
		screens.NewSectorSelection(state, window),
		screens.NewScreenerLaunch(state, window),
		screens.NewTickerEntry(state, window),
		screens.NewChecklist(state, window),
		screens.NewPositionSizing(state, window),
		screens.NewHeatCheck(state, window),
		screens.NewTradeEntry(state, window),
		screens.NewCalendarWithFlags(state, window, state.FeatureFlags, nav), // Pass feature flags and navigator
		screens.NewTradeManagement(state, window, state.FeatureFlags),        // Screen 9: Phase 2 feature
	}

	// Set navigation callbacks on screens that support them
	nav.initializeCallbacks()

	return nav
}

// initializeCallbacks sets navigation callbacks on all screens
func (n *Navigator) initializeCallbacks() {
	for _, screen := range n.screens {
		// Use type assertion to set callbacks if screen supports it
		if s, ok := screen.(interface {
			SetNavCallbacks(onNext, onBack func() error, onCancel func())
		}); ok {
			s.SetNavCallbacks(n.Next, n.Back, n.Cancel)
		}
	}
}

// Next navigates to the next screen in the workflow
func (n *Navigator) Next() error {
	// Validate current screen before proceeding
	if n.currentIndex >= 0 && n.currentIndex < len(n.screens) {
		if !n.screens[n.currentIndex].Validate() {
			return errors.New("current screen validation failed")
		}
	}

	// Auto-save before navigation
	if err := n.AutoSave(); err != nil {
		return fmt.Errorf("auto-save failed: %w", err)
	}

	// Record history for back button
	n.history = append(n.history, n.currentIndex)

	// Move to next screen
	n.currentIndex++
	if n.currentIndex >= len(n.screens) {
		return errors.New("no more screens")
	}

	// Update state
	n.state.CurrentScreen = n.GetCurrentScreenName()

	// Render new screen
	n.window.SetContent(n.screens[n.currentIndex].Render())

	return nil
}

// Back navigates to the previous screen
func (n *Navigator) Back() error {
	if len(n.history) == 0 {
		return errors.New("no previous screen")
	}

	// Auto-save before navigation
	if err := n.AutoSave(); err != nil {
		return fmt.Errorf("auto-save failed: %w", err)
	}

	// Pop from history
	n.currentIndex = n.history[len(n.history)-1]
	n.history = n.history[:len(n.history)-1]

	// Update state
	n.state.CurrentScreen = n.GetCurrentScreenName()

	// Render previous screen
	if n.currentIndex == -1 {
		n.NavigateToDashboard()
	} else {
		n.window.SetContent(n.screens[n.currentIndex].Render())
	}

	return nil
}

// Cancel prompts for confirmation and returns to dashboard
func (n *Navigator) Cancel() {
	dialog.ShowConfirm(
		"Cancel Trade Entry?",
		"Your progress will be saved. Are you sure you want to return to dashboard?",
		func(confirmed bool) {
			if confirmed {
				n.AutoSave()
				n.NavigateToDashboard()
			}
		},
		n.window,
	)
}

// JumpToCalendar navigates directly to calendar view (read-only mode)
func (n *Navigator) JumpToCalendar() {
	// Auto-save current progress
	n.AutoSave()

	// Remember where we came from
	n.history = append(n.history, n.currentIndex)

	// Jump to calendar (screen index 7)
	n.currentIndex = 7
	n.state.CurrentScreen = "calendar"

	// Render calendar
	n.window.SetContent(n.screens[7].Render())
}

// JumpToTradeManagement navigates directly to trade management screen
func (n *Navigator) JumpToTradeManagement() {
	// Auto-save current progress
	n.AutoSave()

	// Remember where we came from
	n.history = append(n.history, n.currentIndex)

	// Jump to trade management (screen index 8)
	n.currentIndex = 8
	n.state.CurrentScreen = "trade_management"

	// Render trade management
	n.window.SetContent(n.screens[8].Render())
}

// NavigateToDashboard returns to the main dashboard
func (n *Navigator) NavigateToDashboard() {
	n.currentIndex = -1
	n.history = []int{}
	n.state.CurrentScreen = "dashboard"

	// Render dashboard (pass navigator so dashboard can navigate to screens)
	dashboard := NewDashboard(n.state, n.window, n)
	n.window.SetContent(dashboard.Render())
}

// NavigateToScreen jumps directly to a specific screen by index
func (n *Navigator) NavigateToScreen(index int) error {
	if index < 0 || index >= len(n.screens) {
		return fmt.Errorf("invalid screen index: %d", index)
	}

	// Auto-save before navigation
	if err := n.AutoSave(); err != nil {
		return fmt.Errorf("auto-save failed: %w", err)
	}

	// Record history
	n.history = append(n.history, n.currentIndex)

	// Navigate to screen
	n.currentIndex = index
	n.state.CurrentScreen = n.GetCurrentScreenName()
	n.window.SetContent(n.screens[n.currentIndex].Render())

	return nil
}

// AutoSave saves current trade progress
func (n *Navigator) AutoSave() error {
	if n.state.CurrentTrade == nil {
		return nil // Nothing to save
	}
	return storage.SaveInProgressTrade(n.state.CurrentTrade)
}

// ValidateCurrentScreen validates the current screen's data
func (n *Navigator) ValidateCurrentScreen() bool {
	if n.currentIndex < 0 || n.currentIndex >= len(n.screens) {
		return true // Dashboard or out of bounds
	}
	return n.screens[n.currentIndex].Validate()
}

// GetCurrentScreenName returns human-readable screen name
func (n *Navigator) GetCurrentScreenName() string {
	if n.currentIndex < 0 {
		return "dashboard"
	}

	if n.currentIndex < len(n.screens) {
		return n.screens[n.currentIndex].GetName()
	}
	return "unknown"
}

// GetCurrentIndex returns the current screen index
func (n *Navigator) GetCurrentIndex() int {
	return n.currentIndex
}

// CanGoBack returns true if there's navigation history
func (n *Navigator) CanGoBack() bool {
	return len(n.history) > 0
}

// GetHistoryDepth returns the number of screens in history
func (n *Navigator) GetHistoryDepth() int {
	return len(n.history)
}

// ClearHistory resets the navigation history
func (n *Navigator) ClearHistory() {
	n.history = []int{}
}

// ShowHelp displays context-sensitive help for the current screen
func (n *Navigator) ShowHelp() {
	screenName := n.GetCurrentScreenName()
	help.ShowHelpDialog(screenName, n.window)
}
