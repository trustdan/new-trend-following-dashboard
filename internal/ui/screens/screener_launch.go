package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// ScreenerLaunch represents Screen 2: Launch FINVIZ Screeners
type ScreenerLaunch struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewScreenerLaunch creates a new screener launch screen
func NewScreenerLaunch(state *appcore.AppState, window fyne.Window) *ScreenerLaunch {
	return &ScreenerLaunch{
		state:  state,
		window: window,
	}
}

// Render renders the screener launch UI
func (s *ScreenerLaunch) Render() fyne.CanvasObject {
	title := widget.NewLabel("Launch FINVIZ Screener")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// TODO: Create screener buttons based on selected sector
	universeBtn := widget.NewButton("Launch Universe Screener", func() {
		// TODO: Open FINVIZ URL in browser
	})

	pullbackBtn := widget.NewButton("Launch Pullback Screener", func() {
		// TODO: Open FINVIZ URL in browser
	})

	breakoutBtn := widget.NewButton("Launch Breakout Screener", func() {
		// TODO: Open FINVIZ URL in browser
	})

	continueBtn := widget.NewButton("I've reviewed the screener →", func() {
		// TODO: Navigate to ticker entry
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Choose your screening approach:"),
		universeBtn,
		pullbackBtn,
		breakoutBtn,
		widget.NewSeparator(),
		continueBtn,
	)

	return container.NewPadded(content)
}
