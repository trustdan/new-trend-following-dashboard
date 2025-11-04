package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// Dashboard represents the main dashboard screen
type Dashboard struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewDashboard creates a new dashboard
func NewDashboard(state *appcore.AppState, window fyne.Window) *Dashboard {
	return &Dashboard{
		state:  state,
		window: window,
	}
}

// Render renders the dashboard UI
func (d *Dashboard) Render() fyne.CanvasObject {
	title := widget.NewLabel("TF-Engine 2.0 - Dashboard")
	title.TextStyle = fyne.TextStyle{Bold: true}

	startButton := widget.NewButton("Start New Trade", func() {
		d.navigateToSectorSelection()
	})

	resumeButton := widget.NewButton("Resume Session", func() {
		// TODO: Resume from saved progress
	})

	calendarButton := widget.NewButton("View Calendar", func() {
		// TODO: Show calendar view
	})

	helpButton := widget.NewButton("Help", func() {
		// TODO: Show help dialog
	})

	// Account info
	accountInfo := widget.NewLabel("Account Equity: $100,000\nRisk per Trade: 0.75%")

	// Heat status
	heatInfo := widget.NewLabel("Portfolio Heat: 0%\nAvailable: 4.0%")

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		startButton,
		resumeButton,
		calendarButton,
		helpButton,
		widget.NewSeparator(),
		widget.NewLabel("Account Settings"),
		accountInfo,
		widget.NewSeparator(),
		widget.NewLabel("Heat Status"),
		heatInfo,
	)

	return container.NewPadded(content)
}

func (d *Dashboard) navigateToSectorSelection() {
	// TODO: Navigate to sector selection screen
	// Will be implemented with proper navigation
}
