package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/ui"
)

func main() {
	// Create Fyne application
	fyneApp := app.NewWithID("com.tfengine.v2")
	fyneApp.Settings().SetTheme(&ui.TFEngineTheme{Mode: "day"})

	// Create main window
	window := fyneApp.NewWindow("TF-Engine 2.0 - Trend Following Trading System")
	window.Resize(fyne.NewSize(1200, 800))

	// Initialize application state
	appState := appcore.NewAppState()

	// Load policy file
	if err := appState.LoadPolicy("data/policy.v1.json"); err != nil {
		// Fall back to safe mode
		appState.UseSafeMode()
	}

	// Create dashboard (initial screen)
	dashboard := ui.NewDashboard(appState, window)

	// Set content and show
	window.SetContent(container.NewStack(dashboard.Render()))
	window.ShowAndRun()
}
