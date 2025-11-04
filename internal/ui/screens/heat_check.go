package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// HeatCheck represents Screen 6: Portfolio Heat Validation
type HeatCheck struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewHeatCheck creates a new heat check screen
func NewHeatCheck(state *appcore.AppState, window fyne.Window) *HeatCheck {
	return &HeatCheck{
		state:  state,
		window: window,
	}
}

// Render renders the heat check UI
func (h *HeatCheck) Render() fyne.CanvasObject {
	title := widget.NewLabel("Portfolio Heat Check")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// TODO: Calculate and display heat metrics
	heatInfo := widget.NewLabel("Portfolio Heat: 0% / 4.0%\nHealthcare Heat: 0% / 1.5%")

	continueBtn := widget.NewButton("Continue to Trade Entry →", func() {
		// TODO: Navigate to trade entry
	})

	content := container.NewVBox(
		title,
		heatInfo,
		continueBtn,
	)

	return container.NewPadded(content)
}
