package screens

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/storage"
)

// TradeEntry represents Screen 7: Options Strategy Selection
type TradeEntry struct {
	state  *appcore.AppState
	window fyne.Window

	// UI components
	strategySelect *widget.Select
	strike1Entry   *widget.Entry
	strike2Entry   *widget.Entry
	strike3Entry   *widget.Entry
	strike4Entry   *widget.Entry
	expirationDate *widget.Entry
	premiumEntry   *widget.Entry
	saveBtn        *widget.Button

	// Dynamic containers for conditional fields
	strikeContainer *fyne.Container
}

// NewTradeEntry creates a new trade entry screen
func NewTradeEntry(state *appcore.AppState, window fyne.Window) *TradeEntry {
	te := &TradeEntry{
		state:  state,
		window: window,
	}

	te.initializeComponents()
	return te
}

// initializeComponents sets up all UI components
func (t *TradeEntry) initializeComponents() {
	// All 26 options strategies
	strategies := []string{
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

	t.strategySelect = widget.NewSelect(strategies, func(value string) {
		t.onStrategySelected(value)
	})
	t.strategySelect.PlaceHolder = "Select an options strategy..."

	// Strike price entries
	t.strike1Entry = widget.NewEntry()
	t.strike1Entry.SetPlaceHolder("Strike 1 (e.g., 450)")

	t.strike2Entry = widget.NewEntry()
	t.strike2Entry.SetPlaceHolder("Strike 2 (e.g., 460)")

	t.strike3Entry = widget.NewEntry()
	t.strike3Entry.SetPlaceHolder("Strike 3 (e.g., 440)")

	t.strike4Entry = widget.NewEntry()
	t.strike4Entry.SetPlaceHolder("Strike 4 (e.g., 470)")

	// Expiration date (DTE format)
	t.expirationDate = widget.NewEntry()
	t.expirationDate.SetPlaceHolder("Days to expiration (e.g., 45)")

	// Premium
	t.premiumEntry = widget.NewEntry()
	t.premiumEntry.SetPlaceHolder("Total premium (e.g., 2.50)")

	// Save button
	t.saveBtn = widget.NewButton("Save Trade & View Calendar →", func() {
		t.saveTrade()
	})
	t.saveBtn.Importance = widget.HighImportance

	// Initially show only 2 strikes (most common)
	t.strikeContainer = container.NewVBox()
}

// onStrategySelected updates the UI based on selected strategy
func (t *TradeEntry) onStrategySelected(strategy string) {
	if t.state.CurrentTrade != nil {
		t.state.CurrentTrade.OptionsStrategy = strategy
	}

	// Clear and rebuild strike container based on strategy requirements
	t.strikeContainer.Objects = nil

	// Determine how many strikes this strategy needs
	strikeCount := t.getRequiredStrikes(strategy)

	// Add strike fields
	if strikeCount >= 1 {
		t.strikeContainer.Add(widget.NewLabel("Strike Price 1:"))
		t.strikeContainer.Add(t.strike1Entry)
	}
	if strikeCount >= 2 {
		t.strikeContainer.Add(widget.NewLabel("Strike Price 2:"))
		t.strikeContainer.Add(t.strike2Entry)
	}
	if strikeCount >= 3 {
		t.strikeContainer.Add(widget.NewLabel("Strike Price 3:"))
		t.strikeContainer.Add(t.strike3Entry)
	}
	if strikeCount >= 4 {
		t.strikeContainer.Add(widget.NewLabel("Strike Price 4:"))
		t.strikeContainer.Add(t.strike4Entry)
	}

	t.strikeContainer.Refresh()
}

// getRequiredStrikes returns the number of strikes needed for a strategy
func (t *TradeEntry) getRequiredStrikes(strategy string) int {
	// Single-leg strategies (1 strike)
	singleLeg := []string{
		"Long call",
		"Long put",
		"Covered call",
		"Cash-secured put",
	}

	// Two-leg strategies (2 strikes)
	twoLeg := []string{
		"Bull call spread",
		"Bear put spread",
		"Bull put credit spread",
		"Bear call credit spread",
		"Calendar call spread",
		"Calendar put spread",
		"Diagonal call spread",
		"Diagonal put spread",
		"Straddle",
		"Strangle",
	}

	// Three-leg strategies (3 strikes)
	threeLeg := []string{
		"Long put butterfly",
		"Long call butterfly",
		"Short put butterfly",
		"Short call butterfly",
		"Call ratio backspread",
		"Put ratio backspread",
		"Call broken wing",
		"Put broken wing",
	}

	// Four-leg strategies (4 strikes)
	fourLeg := []string{
		"Iron butterfly",
		"Iron condor",
		"Inverse iron butterfly",
		"Inverse iron condor",
	}

	// Check which category
	for _, s := range singleLeg {
		if s == strategy {
			return 1
		}
	}
	for _, s := range twoLeg {
		if s == strategy {
			return 2
		}
	}
	for _, s := range threeLeg {
		if s == strategy {
			return 3
		}
	}
	for _, s := range fourLeg {
		if s == strategy {
			return 4
		}
	}

	return 2 // Default to 2-leg spread
}

// Render renders the trade entry UI
func (t *TradeEntry) Render() fyne.CanvasObject {
	title := widget.NewLabel("Screen 7: Options Strategy Selection")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Show current trade summary
	var summaryText string
	if t.state.CurrentTrade != nil {
		summaryText = fmt.Sprintf("Ticker: %s | Sector: %s | Strategy: %s | Conviction: %d | Risk: $%.2f",
			t.state.CurrentTrade.Ticker,
			t.state.CurrentTrade.Sector,
			t.state.CurrentTrade.Strategy,
			t.state.CurrentTrade.Conviction,
			t.state.CurrentTrade.MaxLoss,
		)
	} else {
		summaryText = "No trade in progress"
	}
	summary := widget.NewLabel(summaryText)

	// Instructions
	instructions := widget.NewLabel("Select your options structure and enter the trade details:")

	// Main form
	form := container.NewVBox(
		widget.NewLabel("Options Strategy:"),
		t.strategySelect,
		widget.NewSeparator(),

		// Dynamic strike fields (populated by onStrategySelected)
		t.strikeContainer,
		widget.NewSeparator(),

		widget.NewLabel("Expiration Date:"),
		t.expirationDate,
		widget.NewLabel("(Enter days to expiration, e.g., 45 for 45 DTE)"),
		widget.NewSeparator(),

		widget.NewLabel("Total Premium:"),
		t.premiumEntry,
		widget.NewLabel("(Credit received or debit paid per contract)"),
		widget.NewSeparator(),

		t.saveBtn,
	)

	// Pre-populate if strategy already selected
	if t.state.CurrentTrade != nil && t.state.CurrentTrade.OptionsStrategy != "" {
		t.strategySelect.SetSelected(t.state.CurrentTrade.OptionsStrategy)

		if t.state.CurrentTrade.Strike1 > 0 {
			t.strike1Entry.SetText(fmt.Sprintf("%.2f", t.state.CurrentTrade.Strike1))
		}
		if t.state.CurrentTrade.Strike2 > 0 {
			t.strike2Entry.SetText(fmt.Sprintf("%.2f", t.state.CurrentTrade.Strike2))
		}
		if t.state.CurrentTrade.Strike3 > 0 {
			t.strike3Entry.SetText(fmt.Sprintf("%.2f", t.state.CurrentTrade.Strike3))
		}
		if t.state.CurrentTrade.Strike4 > 0 {
			t.strike4Entry.SetText(fmt.Sprintf("%.2f", t.state.CurrentTrade.Strike4))
		}

		if !t.state.CurrentTrade.ExpirationDate.IsZero() {
			dte := int(time.Until(t.state.CurrentTrade.ExpirationDate).Hours() / 24)
			t.expirationDate.SetText(fmt.Sprintf("%d", dte))
		}

		if t.state.CurrentTrade.Premium > 0 {
			t.premiumEntry.SetText(fmt.Sprintf("%.2f", t.state.CurrentTrade.Premium))
		}
	}

	content := container.NewVBox(
		title,
		summary,
		widget.NewSeparator(),
		instructions,
		form,
	)

	return container.NewPadded(content)
}

// saveTrade validates and saves the completed trade
func (t *TradeEntry) saveTrade() {
	if t.state.CurrentTrade == nil {
		dialog.ShowError(fmt.Errorf("no trade in progress"), t.window)
		return
	}

	// Validate all required fields
	if t.strategySelect.Selected == "" {
		dialog.ShowError(fmt.Errorf("please select an options strategy"), t.window)
		return
	}

	// Parse strikes
	var err error
	if t.strike1Entry.Text != "" {
		t.state.CurrentTrade.Strike1, err = strconv.ParseFloat(t.strike1Entry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid strike 1: %v", err), t.window)
			return
		}
	}

	if t.strike2Entry.Text != "" {
		t.state.CurrentTrade.Strike2, err = strconv.ParseFloat(t.strike2Entry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid strike 2: %v", err), t.window)
			return
		}
	}

	if t.strike3Entry.Text != "" {
		t.state.CurrentTrade.Strike3, err = strconv.ParseFloat(t.strike3Entry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid strike 3: %v", err), t.window)
			return
		}
	}

	if t.strike4Entry.Text != "" {
		t.state.CurrentTrade.Strike4, err = strconv.ParseFloat(t.strike4Entry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid strike 4: %v", err), t.window)
			return
		}
	}

	// Parse expiration date (DTE format)
	if t.expirationDate.Text == "" {
		dialog.ShowError(fmt.Errorf("please enter days to expiration"), t.window)
		return
	}

	dte, err := strconv.Atoi(t.expirationDate.Text)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid days to expiration: %v", err), t.window)
		return
	}
	t.state.CurrentTrade.ExpirationDate = time.Now().AddDate(0, 0, dte)

	// Parse premium
	if t.premiumEntry.Text != "" {
		t.state.CurrentTrade.Premium, err = strconv.ParseFloat(t.premiumEntry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid premium: %v", err), t.window)
			return
		}
	}

	// Save trade to storage
	t.state.CurrentTrade.UpdatedAt = time.Now()
	err = storage.SaveCompletedTrade(t.state.CurrentTrade)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to save trade: %v", err), t.window)
		return
	}

	// Add to AllTrades for calendar display
	t.state.AllTrades = append(t.state.AllTrades, *t.state.CurrentTrade)

	// Show success message
	dialog.ShowInformation(
		"Trade Saved",
		fmt.Sprintf("Trade saved successfully!\n\nTicker: %s\nOptions: %s\nExpiration: %s",
			t.state.CurrentTrade.Ticker,
			t.state.CurrentTrade.OptionsStrategy,
			t.state.CurrentTrade.ExpirationDate.Format("2006-01-02"),
		),
		t.window,
	)

	// Clear current trade (ready for next one)
	t.state.CurrentTrade = nil

	// Navigation to calendar will be handled by navigator
}

// Validate checks if the screen's data is valid
func (t *TradeEntry) Validate() bool {
	// Options strategy must be selected and expiration date entered
	if t.state.CurrentTrade == nil {
		return false
	}

	return t.state.CurrentTrade.OptionsStrategy != "" &&
		t.expirationDate.Text != ""
}

// GetName returns the screen name
func (t *TradeEntry) GetName() string {
	return "trade_entry"
}
