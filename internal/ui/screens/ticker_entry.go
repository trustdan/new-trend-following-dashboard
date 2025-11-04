package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// TickerEntry represents Screen 3: Ticker + Strategy Entry (with cooldown)
type TickerEntry struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewTickerEntry creates a new ticker entry screen
func NewTickerEntry(state *appcore.AppState, window fyne.Window) *TickerEntry {
	return &TickerEntry{
		state:  state,
		window: window,
	}
}

// Render renders the ticker entry UI
func (t *TickerEntry) Render() fyne.CanvasObject {
	title := widget.NewLabel("Enter Trade Details")
	title.TextStyle = fyne.TextStyle{Bold: true}

	tickerEntry := widget.NewEntry()
	tickerEntry.SetPlaceHolder("Enter ticker (e.g., UNH, MSFT)")

	strategyLabel := widget.NewLabel("Select Strategy (Pine Script):")

	// TODO: Filter strategies by selected sector
	strategySelect := widget.NewSelect([]string{"Alt10", "Alt26", "Alt43"}, func(value string) {
		// Strategy selected
	})

	confirmBtn := widget.NewButton("Confirm Selection", func() {
		t.state.StartCooldown()
		// TODO: Show cooldown screen
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Ticker Symbol:"),
		tickerEntry,
		strategyLabel,
		strategySelect,
		widget.NewSeparator(),
		confirmBtn,
	)

	return container.NewPadded(content)
}
