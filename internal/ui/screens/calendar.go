package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
)

// Calendar represents Screen 8: Trade Calendar View
type Calendar struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewCalendar creates a new calendar screen
func NewCalendar(state *appcore.AppState, window fyne.Window) *Calendar {
	return &Calendar{
		state:  state,
		window: window,
	}
}

// Render renders the calendar UI
func (c *Calendar) Render() fyne.CanvasObject {
	title := widget.NewLabel("Trade Calendar")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// TODO: Implement horserace timeline view
	// Y-axis: Sectors
	// X-axis: Time (-14 days to +84 days)
	// Bars: Trades labeled with ticker symbols

	placeholder := widget.NewLabel("Calendar view (to be implemented)")

	content := container.NewVBox(
		title,
		placeholder,
	)

	return container.NewPadded(content)
}

// Validate checks if the screen's data is valid
func (s *Calendar) Validate() bool {
	// Calendar is always valid (display-only)
	return true
}

// GetName returns the screen name
func (s *Calendar) GetName() string {
	return "calendar"
}
