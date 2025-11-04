package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// Checklist represents Screen 4: Anti-Impulsivity Checklist
type Checklist struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewChecklist creates a new checklist screen
func NewChecklist(state *appcore.AppState, window fyne.Window) *Checklist {
	return &Checklist{
		state:  state,
		window: window,
	}
}

// Render renders the checklist UI
func (c *Checklist) Render() fyne.CanvasObject {
	title := widget.NewLabel("Pre-Trade Checklist")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// TODO: Create checklist items from policy
	requiredGates := container.NewVBox(
		widget.NewCheck("From Preset (SIG_REQ)", func(bool) {}),
		widget.NewCheck("Trend Confirmed (RISK_REQ)", func(bool) {}),
		widget.NewCheck("Liquidity OK (OPT_REQ)", func(bool) {}),
		widget.NewCheck("TV Confirm (EXIT_REQ)", func(bool) {}),
		widget.NewCheck("Earnings OK (BEHAV_REQ)", func(bool) {}),
	)

	optionalGates := container.NewVBox(
		widget.NewCheck("Regime OK", func(bool) {}),
		widget.NewCheck("No Chase", func(bool) {}),
		widget.NewCheck("Journal Entry Written", func(bool) {}),
	)

	evaluateBtn := widget.NewButton("Evaluate Checklist", func() {
		// TODO: Validate all required gates passed
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Required Gates (All Must Pass)"),
		requiredGates,
		widget.NewLabel("Optional Quality Items"),
		optionalGates,
		widget.NewSeparator(),
		evaluateBtn,
	)

	return container.NewPadded(content)
}

// Validate checks if the screen's data is valid
func (s *Checklist) Validate() bool {
	// Checklist must be completed and cooldown must be done
	return s.state.CurrentTrade != nil &&
		s.state.CurrentTrade.ChecklistPassed &&
		s.state.IsCooldownComplete()
}

// GetName returns the screen name
func (s *Checklist) GetName() string {
	return "checklist"
}
