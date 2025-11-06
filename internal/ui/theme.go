package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// TFEngineTheme implements custom day/night mode themes
type TFEngineTheme struct {
	Mode   string       // "day" or "night"
	window fyne.Window  // Reference to window for refreshing UI
}

// NewTFEngineTheme creates a new theme with the specified mode
func NewTFEngineTheme(window fyne.Window) *TFEngineTheme {
	return &TFEngineTheme{
		Mode:   "day",
		window: window,
	}
}

// Day mode colors
var (
	DayBackground      = color.RGBA{240, 248, 241, 255} // #F0F8F1 Very Light Green
	DayPrimary         = color.RGBA{46, 125, 50, 255}   // #2E7D32 Forest Green
	DayText            = color.RGBA{27, 94, 32, 255}    // #1B5E20 Dark Text
	DayButton          = color.RGBA{232, 245, 233, 255} // #E8F5E9 Light Green for buttons
	DayDisabledText    = color.RGBA{158, 158, 158, 255} // #9E9E9E Grey
)

// Night mode colors (proper contrast for readability)
var (
	NightBackground     = color.RGBA{15, 76, 58, 255}    // #0F4C3A Dark British Racing Green
	NightPrimary        = color.RGBA{102, 187, 106, 255} // #66BB6A Medium Green
	NightText           = color.RGBA{232, 245, 233, 255} // #E8F5E9 Very Light Text
	NightButton         = color.RGBA{27, 94, 32, 255}    // #1B5E20 Dark Green for buttons
	NightButtonText     = color.RGBA{232, 245, 233, 255} // #E8F5E9 Light Text on buttons
	NightDisabledText   = color.RGBA{117, 117, 117, 255} // #757575 Medium Grey
)

func (t *TFEngineTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.Mode == "night" {
		switch name {
		case theme.ColorNameBackground:
			return NightBackground
		case theme.ColorNameButton:
			return NightButton
		case theme.ColorNameDisabledButton:
			return color.RGBA{40, 40, 40, 255} // Dark grey for disabled buttons
		case theme.ColorNameForeground:
			return NightText
		case theme.ColorNameDisabled:
			return NightDisabledText
		case theme.ColorNamePrimary:
			return NightPrimary
		case theme.ColorNameHover:
			return color.RGBA{46, 125, 50, 255} // Lighter green on hover
		case theme.ColorNameFocus:
			return NightPrimary
		case theme.ColorNamePlaceHolder:
			return color.RGBA{158, 158, 158, 255} // Grey placeholder
		case theme.ColorNamePressed:
			return color.RGBA{27, 94, 32, 255} // Dark green when pressed
		case theme.ColorNameScrollBar:
			return color.RGBA{102, 187, 106, 255} // Medium green scrollbar
		case theme.ColorNameShadow:
			return color.RGBA{0, 0, 0, 100} // Semi-transparent black
		}
	} else {
		switch name {
		case theme.ColorNameBackground:
			return DayBackground
		case theme.ColorNameButton:
			return DayButton
		case theme.ColorNameDisabledButton:
			return color.RGBA{224, 224, 224, 255} // Light grey for disabled buttons
		case theme.ColorNameForeground:
			return DayText
		case theme.ColorNameDisabled:
			return DayDisabledText
		case theme.ColorNamePrimary:
			return DayPrimary
		case theme.ColorNameHover:
			return color.RGBA{76, 175, 80, 255} // Brighter green on hover
		case theme.ColorNameFocus:
			return DayPrimary
		case theme.ColorNamePlaceHolder:
			return color.RGBA{158, 158, 158, 255} // Grey placeholder
		case theme.ColorNamePressed:
			return color.RGBA{27, 94, 32, 255} // Darker green when pressed
		case theme.ColorNameScrollBar:
			return color.RGBA{67, 160, 71, 255} // Medium green scrollbar
		case theme.ColorNameShadow:
			return color.RGBA{0, 0, 0, 50} // Semi-transparent black
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *TFEngineTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *TFEngineTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *TFEngineTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16 // Body text: 16pt minimum
	case theme.SizeNameHeadingText:
		return 24 // H1: 24pt
	case theme.SizeNameSubHeadingText:
		return 20 // H2: 20pt
	}
	return theme.DefaultTheme().Size(name)
}

// ToggleMode switches between day and night mode
func (t *TFEngineTheme) ToggleMode() string {
	if t.Mode == "day" {
		t.Mode = "night"
	} else {
		t.Mode = "day"
	}

	// Refresh the UI to apply new theme
	if t.window != nil {
		fyne.CurrentApp().Settings().SetTheme(t)
		t.window.Content().Refresh()
	}

	return t.Mode
}

// GetMode returns the current theme mode
func (t *TFEngineTheme) GetMode() string {
	return t.Mode
}
