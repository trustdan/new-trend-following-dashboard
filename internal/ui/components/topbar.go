package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/appcore"
)

// TopBar provides quick navigation to key screens and reference materials
type TopBar struct {
	state     *appcore.AppState
	window    fyne.Window
	onSettings func()
	onCalendar func()
	onReferences func(refType string)
}

// NewTopBar creates a new top bar navigation component
func NewTopBar(state *appcore.AppState, window fyne.Window, onSettings, onCalendar func(), onReferences func(refType string)) *TopBar {
	return &TopBar{
		state:        state,
		window:       window,
		onSettings:   onSettings,
		onCalendar:   onCalendar,
		onReferences: onReferences,
	}
}

// Render creates the top bar UI
func (t *TopBar) Render() fyne.CanvasObject {
	// Settings button
	settingsBtn := widget.NewButton("Settings", func() {
		if t.onSettings != nil {
			t.onSettings()
		}
	})

	// Calendar button
	calendarBtn := widget.NewButton("Calendar", func() {
		if t.onCalendar != nil {
			t.onCalendar()
		}
	})

	// Pine Script Strategies dropdown
	strategiesMenu := fyne.NewMenu("Pine Scripts",
		fyne.NewMenuItem("📖 Complete Strategies Guide", func() {
			if t.onReferences != nil {
				t.onReferences("pine_scripts_guide")
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Alt10 - Profit Targets", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt10")
			}
		}),
		fyne.NewMenuItem("Alt22 - Parabolic SAR", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt22")
			}
		}),
		fyne.NewMenuItem("Alt26 - Fractional Pyramiding", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt26")
			}
		}),
		fyne.NewMenuItem("Alt28 - ADX Filter", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt28")
			}
		}),
		fyne.NewMenuItem("Alt39 - Age-Based Targets", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt39")
			}
		}),
		fyne.NewMenuItem("Alt43 - Volatility Adaptive", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt43")
			}
		}),
		fyne.NewMenuItem("Alt45 - Dual Momentum", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt45")
			}
		}),
		fyne.NewMenuItem("Alt46 - Sector Adaptive", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt46")
			}
		}),
		fyne.NewMenuItem("Alt47 - Momentum Scaled", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_alt47")
			}
		}),
		fyne.NewMenuItem("Turtle Core V2.2", func() {
			if t.onReferences != nil {
				t.onReferences("strategy_turtle_core")
			}
		}),
	)

	// Finviz Screeners dropdown
	screenersMenu := fyne.NewMenu("Screeners",
		fyne.NewMenuItem("Master Screener Guide", func() {
			if t.onReferences != nil {
				t.onReferences("screener_master")
			}
		}),
		fyne.NewMenuItem("Daily Cheat Sheet", func() {
			if t.onReferences != nil {
				t.onReferences("screener_daily")
			}
		}),
		fyne.NewMenuItem("Screener Decision Tree", func() {
			if t.onReferences != nil {
				t.onReferences("screener_decision_tree")
			}
		}),
		fyne.NewMenuItem("Start Here", func() {
			if t.onReferences != nil {
				t.onReferences("screener_start")
			}
		}),
	)

	// Create popup menus
	strategiesPopup := widget.NewPopUpMenu(strategiesMenu, t.window.Canvas())
	screenersPopup := widget.NewPopUpMenu(screenersMenu, t.window.Canvas())

	// Reference buttons that show popup menus
	strategiesBtn := widget.NewButton("Pine Scripts ▾", func() {
		strategiesPopup.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(settingsBtn))
	})

	screenersBtn := widget.NewButton("Screeners ▾", func() {
		screenersPopup.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(settingsBtn))
	})

	// Spacer to push everything to the left
	spacer := widget.NewLabel("")

	// Create horizontal container with buttons
	topBar := container.NewHBox(
		settingsBtn,
		calendarBtn,
		widget.NewSeparator(),
		strategiesBtn,
		screenersBtn,
		spacer,
	)

	return topBar
}
