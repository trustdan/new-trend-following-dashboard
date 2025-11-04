package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/config"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
)

// TradeManagement represents Screen 9: Trade Management (Phase 2 Feature)
type TradeManagement struct {
	state        *appcore.AppState
	window       fyne.Window
	featureFlags *config.FeatureFlags
	filterStatus string // "all", "active", "closed"
}

// NewTradeManagement creates a new trade management screen
func NewTradeManagement(state *appcore.AppState, window fyne.Window, featureFlags *config.FeatureFlags) *TradeManagement {
	return &TradeManagement{
		state:        state,
		window:       window,
		featureFlags: featureFlags,
		filterStatus: "all",
	}
}

// Render renders the trade management UI
func (tm *TradeManagement) Render() fyne.CanvasObject {
	title := widget.NewLabel("Screen 9: Trade Management")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Check if feature is enabled
	if tm.featureFlags != nil && !tm.featureFlags.IsEnabled("trade_management") {
		return tm.renderDisabledState()
	}

	// Filter dropdown
	filterLabel := widget.NewLabel("Filter:")
	filterSelect := widget.NewSelect([]string{"Show All", "Active Only", "Closed Only"}, func(value string) {
		switch value {
		case "Show All":
			tm.filterStatus = "all"
		case "Active Only":
			tm.filterStatus = "active"
		case "Closed Only":
			tm.filterStatus = "closed"
		}
		tm.window.SetContent(tm.Render())
	})
	filterSelect.Selected = "Show All"

	filterBar := container.NewHBox(filterLabel, filterSelect)

	// Load trades
	trades := tm.getFilteredTrades()

	// Create trades table
	tradesTable := tm.createTradesTable(trades)

	// Action buttons
	backBtn := widget.NewButton("Back to Calendar", func() {
		// Navigator will handle navigation back
	})

	buttons := container.NewHBox(backBtn)

	content := container.NewVBox(
		title,
		widget.NewLabel("View, edit, or delete your trade history"),
		widget.NewSeparator(),
		filterBar,
		widget.NewSeparator(),
		tradesTable,
		widget.NewSeparator(),
		buttons,
	)

	return container.NewScroll(content)
}

// renderDisabledState shows a message when feature flag is OFF
func (tm *TradeManagement) renderDisabledState() fyne.CanvasObject {
	title := widget.NewLabel("Screen 9: Trade Management")
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("This feature is currently disabled.")
	message.Wrapping = fyne.TextWrapWord

	flag := tm.featureFlags.GetFlag("trade_management")
	var details string
	if flag != nil {
		details = fmt.Sprintf("Feature: %s\nPhase: %d\nAvailable in version: %s",
			flag.Description, flag.Phase, flag.SinceVersion)
	}

	detailsLabel := widget.NewLabel(details)
	detailsLabel.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Back to Calendar", func() {
		// Navigator will handle navigation back
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		message,
		detailsLabel,
		widget.NewSeparator(),
		backBtn,
	)

	return container.NewCenter(content)
}

// createTradesTable creates a table widget with trade data
func (tm *TradeManagement) createTradesTable(trades []models.Trade) fyne.CanvasObject {
	if len(trades) == 0 {
		return widget.NewLabel("No trades to display")
	}

	// Create container for table rows
	rows := []fyne.CanvasObject{}

	// Header row
	header := container.NewHBox(
		tm.createTableCell("Date", 100, true),
		tm.createTableCell("Ticker", 80, true),
		tm.createTableCell("Sector", 120, true),
		tm.createTableCell("Strategy", 100, true),
		tm.createTableCell("Options", 150, true),
		tm.createTableCell("P&L", 80, true),
		tm.createTableCell("Status", 80, true),
		tm.createTableCell("Actions", 150, true),
	)
	rows = append(rows, header)
	rows = append(rows, widget.NewSeparator())

	// Data rows
	for i := range trades {
		trade := &trades[i] // Capture pointer for closures
		row := tm.createTradeRow(trade)
		rows = append(rows, row)
		if i < len(trades)-1 {
			rows = append(rows, widget.NewSeparator())
		}
	}

	return container.NewVBox(rows...)
}

// createTradeRow creates a single table row for a trade
func (tm *TradeManagement) createTradeRow(trade *models.Trade) fyne.CanvasObject {
	// Format date
	dateStr := trade.CreatedAt.Format("2006-01-02")

	// Format P&L
	pnl := trade.GetPnL()
	pnlStr := fmt.Sprintf("$%.2f", pnl)
	if pnl > 0 {
		pnlStr = "+" + pnlStr
	}

	// Create cells
	dateCell := tm.createTableCell(dateStr, 100, false)
	tickerCell := tm.createTableCell(trade.Ticker, 80, false)
	sectorCell := tm.createTableCell(trade.Sector, 120, false)
	strategyCell := tm.createTableCell(trade.Strategy, 100, false)
	optionsCell := tm.createTableCell(trade.OptionsStrategy, 150, false)
	pnlCell := tm.createTableCell(pnlStr, 80, false)
	statusCell := tm.createTableCell(trade.GetStatus(), 80, false)

	// Action buttons
	editBtn := widget.NewButton("Edit", func() {
		tm.editTrade(trade)
	})
	editBtn.Importance = widget.LowImportance

	deleteBtn := widget.NewButton("Delete", func() {
		tm.confirmDeleteTrade(trade)
	})
	deleteBtn.Importance = widget.DangerImportance

	actionsCell := container.NewHBox(editBtn, deleteBtn)

	row := container.NewHBox(
		dateCell,
		tickerCell,
		sectorCell,
		strategyCell,
		optionsCell,
		pnlCell,
		statusCell,
		actionsCell,
	)

	return row
}

// createTableCell creates a table cell with fixed width
func (tm *TradeManagement) createTableCell(text string, width float32, bold bool) fyne.CanvasObject {
	label := widget.NewLabel(text)
	if bold {
		label.TextStyle = fyne.TextStyle{Bold: true}
	}

	cell := container.NewMax(label)
	cell.Resize(fyne.NewSize(width, 0))

	return cell
}

// getFilteredTrades returns trades based on current filter
func (tm *TradeManagement) getFilteredTrades() []models.Trade {
	// Load all trades
	allTrades, err := storage.LoadAllTrades()
	if err != nil {
		return []models.Trade{}
	}

	// Apply filter
	if tm.filterStatus == "all" {
		return allTrades
	}

	filtered := []models.Trade{}
	for _, trade := range allTrades {
		status := trade.GetStatus()
		if tm.filterStatus == "active" && status == "active" {
			filtered = append(filtered, trade)
		} else if tm.filterStatus == "closed" && status == "closed" {
			filtered = append(filtered, trade)
		}
	}

	return filtered
}

// editTrade opens a dialog to edit trade details
func (tm *TradeManagement) editTrade(trade *models.Trade) {
	// Create form fields
	tickerEntry := widget.NewEntry()
	tickerEntry.SetText(trade.Ticker)

	pnlEntry := widget.NewEntry()
	pnlEntry.SetText(fmt.Sprintf("%.2f", trade.GetPnL()))

	statusSelect := widget.NewSelect([]string{"active", "closed", "expired"}, nil)
	statusSelect.Selected = trade.GetStatus()

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Ticker", Widget: tickerEntry},
			{Text: "P&L ($)", Widget: pnlEntry},
			{Text: "Status", Widget: statusSelect},
		},
		OnSubmit: func() {
			// Update trade
			trade.Ticker = tickerEntry.Text

			// Parse P&L
			var pnl float64
			fmt.Sscanf(pnlEntry.Text, "%f", &pnl)
			trade.ProfitLoss = &pnl

			trade.Status = statusSelect.Selected

			// Save trade
			if err := tm.updateTrade(trade); err != nil {
				dialog.ShowError(err, tm.window)
				return
			}

			// Refresh display
			tm.window.SetContent(tm.Render())
		},
		OnCancel: func() {
			// Dialog will close automatically
		},
	}

	// Show dialog
	dialog.ShowForm("Edit Trade", "Save", "Cancel", form.Items, func(submitted bool) {
		if submitted {
			form.OnSubmit()
		}
	}, tm.window)
}

// confirmDeleteTrade shows confirmation dialog before deleting
func (tm *TradeManagement) confirmDeleteTrade(trade *models.Trade) {
	message := fmt.Sprintf("Are you sure you want to delete this trade?\n\nTicker: %s\nSector: %s\nEntry: %s\n\nThis action cannot be undone.",
		trade.Ticker, trade.Sector, trade.CreatedAt.Format("2006-01-02"))

	dialog.ShowConfirm("Delete Trade", message, func(confirmed bool) {
		if confirmed {
			if err := tm.deleteTrade(trade); err != nil {
				dialog.ShowError(err, tm.window)
				return
			}

			// Refresh display
			tm.window.SetContent(tm.Render())
		}
	}, tm.window)
}

// updateTrade saves updated trade to storage
func (tm *TradeManagement) updateTrade(trade *models.Trade) error {
	// Load all trades
	allTrades, err := storage.LoadAllTrades()
	if err != nil {
		return fmt.Errorf("failed to load trades: %w", err)
	}

	// Find and update the trade
	found := false
	for i := range allTrades {
		if allTrades[i].ID == trade.ID {
			allTrades[i] = *trade
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("trade not found: %s", trade.ID)
	}

	// Save all trades
	return storage.SaveAllTrades(allTrades)
}

// deleteTrade removes a trade from storage
func (tm *TradeManagement) deleteTrade(trade *models.Trade) error {
	// Load all trades
	allTrades, err := storage.LoadAllTrades()
	if err != nil {
		return fmt.Errorf("failed to load trades: %w", err)
	}

	// Filter out the deleted trade
	filteredTrades := []models.Trade{}
	for _, t := range allTrades {
		if t.ID != trade.ID {
			filteredTrades = append(filteredTrades, t)
		}
	}

	// Save filtered trades
	return storage.SaveAllTrades(filteredTrades)
}

// Validate validates the screen state (not used for read-only screen)
func (tm *TradeManagement) Validate() bool {
	return true
}

// GetName returns the screen name
func (tm *TradeManagement) GetName() string {
	return "trade_management"
}
