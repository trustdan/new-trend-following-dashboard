package screens

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/analytics"
	"tf-engine/internal/appcore"
	"tf-engine/internal/config"
	"tf-engine/internal/storage"
)

// Analytics represents the analytics dashboard (Phase 2 Feature)
type Analytics struct {
	state        *appcore.AppState
	window       fyne.Window
	featureFlags *config.FeatureFlags
}

// NewAnalytics creates a new analytics screen
func NewAnalytics(state *appcore.AppState, window fyne.Window, featureFlags *config.FeatureFlags) *Analytics {
	return &Analytics{
		state:        state,
		window:       window,
		featureFlags: featureFlags,
	}
}

// Render renders the analytics UI
func (a *Analytics) Render() fyne.CanvasObject {
	title := widget.NewLabel("📊 Advanced Analytics")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Check if feature is enabled
	if a.featureFlags != nil && !a.featureFlags.IsEnabled("advanced_analytics") {
		return a.renderDisabledState()
	}

	// Load trades
	trades, err := storage.LoadAllTrades()
	if err != nil {
		return a.renderError("Failed to load trades: " + err.Error())
	}

	if len(trades) == 0 {
		return a.renderEmpty()
	}

	// Calculate statistics
	overallStats := analytics.CalculateTradeStats(trades)
	sectorStats := analytics.CalculateSectorStats(trades)
	strategyStats := analytics.CalculateStrategyStats(trades)
	equityCurve := analytics.CalculateEquityCurve(trades)

	// Create UI sections
	overallSection := a.renderOverallStats(overallStats)
	sectorSection := a.renderSectorStats(sectorStats)
	strategySection := a.renderStrategyStats(strategyStats)
	equityCurveSection := a.renderEquityCurve(equityCurve)

	// Back button
	backBtn := widget.NewButton("Back to Dashboard", func() {
		// Navigator will handle navigation
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		overallSection,
		widget.NewSeparator(),
		sectorSection,
		widget.NewSeparator(),
		strategySection,
		widget.NewSeparator(),
		equityCurveSection,
		widget.NewSeparator(),
		backBtn,
	)

	return container.NewScroll(content)
}

// renderDisabledState shows a message when feature flag is OFF
func (a *Analytics) renderDisabledState() fyne.CanvasObject {
	title := widget.NewLabel("📊 Advanced Analytics")
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("This feature is currently disabled.")
	message.Wrapping = fyne.TextWrapWord

	flag := a.featureFlags.GetFlag("advanced_analytics")
	var details string
	if flag != nil {
		details = fmt.Sprintf("Feature: %s\nPhase: %d\nAvailable in version: %s",
			flag.Description, flag.Phase, flag.SinceVersion)
	}

	detailsLabel := widget.NewLabel(details)
	detailsLabel.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Back to Dashboard", func() {
		// Navigator will handle navigation
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

// renderError displays an error message
func (a *Analytics) renderError(errorMsg string) fyne.CanvasObject {
	title := widget.NewLabel("📊 Advanced Analytics")
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("Error: " + errorMsg)
	message.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Back to Dashboard", func() {
		// Navigator will handle navigation
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		message,
		widget.NewSeparator(),
		backBtn,
	)

	return container.NewCenter(content)
}

// renderEmpty displays a message when no trades exist
func (a *Analytics) renderEmpty() fyne.CanvasObject {
	title := widget.NewLabel("📊 Advanced Analytics")
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("No completed trades to analyze. Start trading to see your performance statistics!")
	message.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Back to Dashboard", func() {
		// Navigator will handle navigation
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		message,
		widget.NewSeparator(),
		backBtn,
	)

	return container.NewCenter(content)
}

// renderOverallStats displays overall performance statistics
func (a *Analytics) renderOverallStats(stats analytics.TradeStats) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Overall Performance", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Format statistics
	items := []fyne.CanvasObject{
		a.createStatRow("Total Trades", fmt.Sprintf("%d", stats.TotalTrades)),
		a.createStatRow("Win Rate", fmt.Sprintf("%.1f%% (%d/%d)", stats.WinRate, stats.WinningTrades, stats.TotalTrades)),
		a.createStatRow("Total P&L", a.formatPnL(stats.TotalPnL)),
		a.createStatRow("Average P&L", a.formatPnL(stats.AveragePnL)),
		a.createStatRow("Average Win", a.formatPnL(stats.AverageWin)),
		a.createStatRow("Average Loss", a.formatPnL(stats.AverageLoss)),
		a.createStatRow("Largest Win", a.formatPnL(stats.LargestWin)),
		a.createStatRow("Largest Loss", a.formatPnL(stats.LargestLoss)),
		a.createStatRow("Profit Factor", fmt.Sprintf("%.2f", stats.ProfitFactor)),
		a.createStatRow("Max Drawdown", fmt.Sprintf("$%.2f (%.1f%%)", stats.MaxDrawdown, stats.MaxDrawdownPct)),
		a.createStatRow("Current Streak", fmt.Sprintf("%d trades", stats.CurrentStreak)),
		a.createStatRow("Longest Win Streak", fmt.Sprintf("%d trades", stats.LongestWinStreak)),
		a.createStatRow("Longest Loss Streak", fmt.Sprintf("%d trades", stats.LongestLossStreak)),
	}

	content := container.NewVBox(header)
	for _, item := range items {
		content.Add(item)
	}

	return content
}

// renderSectorStats displays performance by sector
func (a *Analytics) renderSectorStats(stats []analytics.SectorStats) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Performance by Sector", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	if len(stats) == 0 {
		return container.NewVBox(header, widget.NewLabel("No sector data available"))
	}

	// Create table-like display
	rows := []fyne.CanvasObject{header}

	// Table header
	headerRow := container.NewHBox(
		a.createTableCell("Sector", 120, true),
		a.createTableCell("Trades", 60, true),
		a.createTableCell("Win Rate", 80, true),
		a.createTableCell("Total P&L", 100, true),
		a.createTableCell("Avg P&L", 100, true),
	)
	rows = append(rows, headerRow)
	rows = append(rows, widget.NewSeparator())

	// Data rows
	for _, stat := range stats {
		row := container.NewHBox(
			a.createTableCell(stat.Sector, 120, false),
			a.createTableCell(fmt.Sprintf("%d", stat.TotalTrades), 60, false),
			a.createTableCell(fmt.Sprintf("%.1f%%", stat.WinRate), 80, false),
			a.createTableCell(a.formatPnL(stat.TotalPnL), 100, false),
			a.createTableCell(a.formatPnL(stat.AveragePnL), 100, false),
		)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

// renderStrategyStats displays performance by strategy
func (a *Analytics) renderStrategyStats(stats []analytics.StrategyStats) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Performance by Strategy", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	if len(stats) == 0 {
		return container.NewVBox(header, widget.NewLabel("No strategy data available"))
	}

	// Create table-like display
	rows := []fyne.CanvasObject{header}

	// Table header
	headerRow := container.NewHBox(
		a.createTableCell("Strategy", 120, true),
		a.createTableCell("Trades", 60, true),
		a.createTableCell("Win Rate", 80, true),
		a.createTableCell("Total P&L", 100, true),
		a.createTableCell("Avg P&L", 100, true),
	)
	rows = append(rows, headerRow)
	rows = append(rows, widget.NewSeparator())

	// Data rows
	for _, stat := range stats {
		row := container.NewHBox(
			a.createTableCell(stat.Strategy, 120, false),
			a.createTableCell(fmt.Sprintf("%d", stat.TotalTrades), 60, false),
			a.createTableCell(fmt.Sprintf("%.1f%%", stat.WinRate), 80, false),
			a.createTableCell(a.formatPnL(stat.TotalPnL), 100, false),
			a.createTableCell(a.formatPnL(stat.AveragePnL), 100, false),
		)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

// renderEquityCurve displays the equity curve chart
func (a *Analytics) renderEquityCurve(curve []analytics.EquityCurvePoint) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Equity Curve", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	if len(curve) == 0 {
		return container.NewVBox(header, widget.NewLabel("No equity curve data available"))
	}

	// For now, display a simple text representation
	// TODO: Use a charting library for visual equity curve
	curveText := fmt.Sprintf("Equity curve with %d data points\n", len(curve))
	curveText += fmt.Sprintf("Starting: $%.2f\n", curve[0].Equity)
	curveText += fmt.Sprintf("Ending: $%.2f\n", curve[len(curve)-1].Equity)
	curveText += fmt.Sprintf("Change: %s\n", a.formatPnL(curve[len(curve)-1].Equity-curve[0].Equity))

	curveLabel := widget.NewLabel(curveText)
	curveLabel.Wrapping = fyne.TextWrapWord

	// Simple chart placeholder
	chartPlaceholder := canvas.NewRectangle(color.NRGBA{R: 50, G: 50, B: 50, A: 255})
	chartPlaceholder.FillColor = theme.BackgroundColor()
	chartPlaceholder.SetMinSize(fyne.NewSize(600, 200))

	note := widget.NewLabel("📈 Visual equity curve chart coming in future update")
	note.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(header, curveLabel, chartPlaceholder, note)
}

// Helper functions

func (a *Analytics) createStatRow(label, value string) fyne.CanvasObject {
	labelWidget := widget.NewLabel(label + ":")
	labelWidget.TextStyle = fyne.TextStyle{Bold: true}

	valueWidget := widget.NewLabel(value)

	return container.NewHBox(labelWidget, valueWidget)
}

func (a *Analytics) createTableCell(text string, width float32, bold bool) fyne.CanvasObject {
	label := widget.NewLabel(text)
	if bold {
		label.TextStyle = fyne.TextStyle{Bold: true}
	}

	cell := container.NewMax(label)
	cell.Resize(fyne.NewSize(width, 0))

	return cell
}

func (a *Analytics) formatPnL(pnl float64) string {
	if pnl > 0 {
		return fmt.Sprintf("+$%.2f", pnl)
	}
	return fmt.Sprintf("$%.2f", pnl)
}

// Validate validates the screen state (not used for read-only screen)
func (a *Analytics) Validate() bool {
	return true
}

// GetName returns the screen name
func (a *Analytics) GetName() string {
	return "analytics"
}
