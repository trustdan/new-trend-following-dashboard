package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// PositionSizing represents Screen 5: Position Size Calculator
type PositionSizing struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewPositionSizing creates a new position sizing screen
func NewPositionSizing(state *appcore.AppState, window fyne.Window) *PositionSizing {
	return &PositionSizing{
		state:  state,
		window: window,
	}
}

// Render renders the position sizing UI
func (p *PositionSizing) Render() fyne.CanvasObject {
	title := widget.NewLabel("Calculate Position Size")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// TODO: Add position sizing calculator
	accountEquityEntry := widget.NewEntry()
	accountEquityEntry.SetText("100000")

	riskPerTradeEntry := widget.NewEntry()
	riskPerTradeEntry.SetText("0.75")

	calculateBtn := widget.NewButton("Calculate", func() {
		// TODO: Calculate position size
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Account Equity:"),
		accountEquityEntry,
		widget.NewLabel("Risk per Trade (%):"),
		riskPerTradeEntry,
		calculateBtn,
	)

	return container.NewPadded(content)
}
