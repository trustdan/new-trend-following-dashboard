package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/appcore"
	"tf-engine/internal/storage"
	"tf-engine/internal/testing/generators"
	"tf-engine/internal/ui/help"
	"tf-engine/internal/ui/screens"
)

// Dashboard represents the main dashboard screen
type Dashboard struct {
	state     *appcore.AppState
	window    fyne.Window
	navigator *Navigator
}

// NewDashboard creates a new dashboard
func NewDashboard(state *appcore.AppState, window fyne.Window, navigator *Navigator) *Dashboard {
	return &Dashboard{
		state:     state,
		window:    window,
		navigator: navigator,
	}
}

// Render renders the dashboard UI
func (d *Dashboard) Render() fyne.CanvasObject {
	title := widget.NewLabel("TF-Engine 2.0 - Dashboard")
	title.TextStyle = fyne.TextStyle{Bold: true}

	startButton := widget.NewButton("Start New Trade", func() {
		if d.navigator != nil {
			d.navigator.NavigateToScreen(0) // Navigate to Sector Selection
		}
	})

	resumeButton := widget.NewButton("Resume Session", func() {
		// TODO: Resume from saved progress
	})

	calendarButton := widget.NewButton("View Calendar", func() {
		if d.navigator != nil {
			d.navigator.JumpToCalendar()
		}
	})

	// Trade Management button (Phase 2 feature)
	manageTradesButton := widget.NewButton("Manage Trades", func() {
		if d.navigator != nil {
			d.navigator.JumpToTradeManagement()
		}
	})

	// Check if Trade Management feature is enabled
	manageTradesEnabled := d.state.FeatureFlags != nil && d.state.FeatureFlags.IsEnabled("trade_management")
	if !manageTradesEnabled {
		manageTradesButton.Disable()
	}

	// Sample Data Generator button (Phase 2 feature)
	sampleDataButton := widget.NewButton("Generate Sample Data", func() {
		d.generateSampleData()
	})

	// Check if Sample Data Generator feature is enabled
	sampleDataEnabled := d.state.FeatureFlags != nil && d.state.FeatureFlags.IsEnabled("sample_data_generator")
	if !sampleDataEnabled {
		sampleDataButton.Disable()
	}

	// Advanced Analytics button (Phase 2 feature)
	analyticsButton := widget.NewButton("📊 View Analytics", func() {
		if d.navigator != nil {
			d.navigator.JumpToAnalytics()
		}
	})

	// Check if Advanced Analytics feature is enabled
	analyticsEnabled := d.state.FeatureFlags != nil && d.state.FeatureFlags.IsEnabled("advanced_analytics")
	if !analyticsEnabled {
		analyticsButton.Disable()
	}

	// Phase 2 features label
	phase2Label := widget.NewLabel("Phase 2 Features:")
	phase2Label.TextStyle = fyne.TextStyle{Italic: true}

	// Vim Mode toggle button (Phase 2 feature)
	var vimModeButton *widget.Button
	if d.navigator != nil && d.navigator.GetVimiumManager() != nil {
		vimModeButton = d.navigator.GetVimiumManager().GetToggleButton()
	} else {
		vimModeButton = widget.NewButton("⌨️ Vim Mode", func() {})
		vimModeButton.Disable()
	}

	helpButton := widget.NewButton("Help", func() {
		help.ShowHelpDialog("welcome", d.window)
	})

	settingsButton := widget.NewButton("⚙️ Settings", func() {
		d.showSettings()
	})

	// Account info (from Settings)
	accountEquity := 25000.0
	riskPercent := 2.8
	if d.state.Settings != nil {
		accountEquity = d.state.Settings.AccountEquity
		riskPercent = d.state.Settings.RiskPerTrade * 100
	}
	accountInfo := widget.NewLabel(fmt.Sprintf("Account Equity: $%.0f\nRisk per Trade: %.2f%%",
		accountEquity, riskPercent))

	// Heat status
	heatInfo := widget.NewLabel("Portfolio Heat: 0%\nAvailable: 4.0%")

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		startButton,
		resumeButton,
		calendarButton,
		settingsButton,
		helpButton,
		widget.NewSeparator(),
		phase2Label,
		manageTradesButton,
		sampleDataButton,
		analyticsButton,
		vimModeButton,
		widget.NewSeparator(),
		widget.NewLabel("Account Settings"),
		accountInfo,
		widget.NewSeparator(),
		widget.NewLabel("Heat Status"),
		heatInfo,
	)

	return container.NewPadded(content)
}

// generateSampleData generates sample trades and saves them
func (d *Dashboard) generateSampleData() {
	// Confirm with user
	dialog.ShowConfirm(
		"Generate Sample Data",
		"This will create 10 sample trades for testing. Continue?",
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// Generate 10 sample trades
			sampleTrades := generators.GenerateSampleTrades(10)

			// Save each trade
			for _, trade := range sampleTrades {
				if err := storage.SaveCompletedTrade(&trade); err != nil {
					dialog.ShowError(
						fmt.Errorf("Failed to save sample trade: %v", err),
						d.window,
					)
					return
				}
			}

			// Reload trades into state
			allTrades, err := storage.LoadAllTrades()
			if err != nil {
				dialog.ShowError(
					fmt.Errorf("Failed to reload trades: %v", err),
					d.window,
				)
				return
			}
			d.state.AllTrades = allTrades

			// Show success message
			dialog.ShowInformation(
				"Success",
				fmt.Sprintf("Generated %d sample trades successfully!\n\nClick 'View Calendar' to see them.", len(sampleTrades)),
				d.window,
			)
		},
		d.window,
	)
}

// showSettings displays the settings screen
func (d *Dashboard) showSettings() {
	settingsScreen := screens.NewSettings(d.state, d.window)

	// Set back callback to return to dashboard
	settingsScreen.SetBackCallback(func() {
		// Refresh dashboard to show updated settings
		d.window.SetContent(d.Render())
	})

	d.window.SetContent(settingsScreen.Render())
}
