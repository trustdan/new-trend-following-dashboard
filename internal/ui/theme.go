package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// TFEngineTheme implements custom day/night mode themes
type TFEngineTheme struct {
	Mode string // "day" or "night"
}

// NewTFEngineTheme creates a new theme with the specified mode
func NewTFEngineTheme() *TFEngineTheme {
	return &TFEngineTheme{Mode: "day"}
}

// Day mode colors
var (
	DayBackground = color.RGBA{232, 245, 233, 255} // #E8F5E9 Light Green
	DayPrimary    = color.RGBA{67, 160, 71, 255}   // #43A047 Medium Green
	DayText       = color.RGBA{27, 94, 32, 255}    // #1B5E20 Dark Text
)

// Night mode colors
var (
	NightBackground = color.RGBA{0, 66, 37, 255}     // #004225 British Racing Green
	NightPrimary    = color.RGBA{129, 199, 132, 255} // #81C784 Light Green
	NightText       = color.RGBA{232, 245, 233, 255} // #E8F5E9 Light Text
)

func (t *TFEngineTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.Mode == "night" {
		switch name {
		case theme.ColorNameBackground:
			return NightBackground
		case theme.ColorNamePrimary:
			return NightPrimary
		case theme.ColorNameForeground:
			return NightText
		}
	} else {
		switch name {
		case theme.ColorNameBackground:
			return DayBackground
		case theme.ColorNamePrimary:
			return DayPrimary
		case theme.ColorNameForeground:
			return DayText
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
