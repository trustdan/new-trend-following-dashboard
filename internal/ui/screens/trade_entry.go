package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// TradeEntry represents Screen 7: Options Strategy Selection
type TradeEntry struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewTradeEntry creates a new trade entry screen
func NewTradeEntry(state *appcore.AppState, window fyne.Window) *TradeEntry {
	return &TradeEntry{
		state:  state,
		window: window,
	}
}

// Render renders the trade entry UI
func (t *TradeEntry) Render() fyne.CanvasObject {
	title := widget.NewLabel("Enter Trade Details")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Options strategy selection
	strategies := []string{
		"Bull Call Spread",
		"Bear Put Spread",
		"Bull Put Credit Spread",
		"Bear Call Credit Spread",
		"Long Call",
		"Long Put",
		"Iron Condor",
		"Iron Butterfly",
		// ... more strategies
	}

	strategySelect := widget.NewSelect(strategies, func(value string) {
		// Strategy selected
	})

	saveBtn := widget.NewButton("Save Trade & View Calendar →", func() {
		// TODO: Save trade and show calendar
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Select Options Strategy:"),
		strategySelect,
		widget.NewSeparator(),
		saveBtn,
	)

	return container.NewPadded(content)
}
