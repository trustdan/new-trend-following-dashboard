package screens

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tf-engine/internal/appcore"
	"tf-engine/internal/storage"
)

// Settings represents the Settings/Configuration screen
type Settings struct {
	state  *appcore.AppState
	window fyne.Window

	// Navigation callbacks
	onBack func()

	// UI components
	accountEntry     *widget.Entry
	riskPercentEntry *widget.Entry
	themeSelect      *widget.Select
}

// NewSettings creates a new settings screen
func NewSettings(state *appcore.AppState, window fyne.Window) *Settings {
	return &Settings{
		state:  state,
		window: window,
	}
}

// SetBackCallback sets the back navigation callback
func (s *Settings) SetBackCallback(onBack func()) {
	s.onBack = onBack
}

// Render renders the settings UI
func (s *Settings) Render() fyne.CanvasObject {
	// Header
	title := widget.NewLabel("⚙️ Account Settings")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Configure your trading account parameters")
	subtitle.Alignment = fyne.TextAlignCenter

	// Info banner
	infoBanner := s.createInfoBanner()

	// Settings form
	form := s.createSettingsForm()

	// Save button
	saveBtn := widget.NewButton("Save Settings", func() {
		s.saveSettings()
	})
	saveBtn.Importance = widget.HighImportance

	// Back to Dashboard button
	backBtn := widget.NewButton("← Back to Dashboard", func() {
		if s.onBack != nil {
			s.onBack()
		}
	})

	buttons := container.NewBorder(
		nil, nil,
		backBtn,
		saveBtn,
		nil,
	)

	content := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		infoBanner,
		widget.NewSeparator(),
		form,
		widget.NewSeparator(),
		buttons,
	)

	return container.NewPadded(container.NewScroll(content))
}

// createInfoBanner creates an informational banner
func (s *Settings) createInfoBanner() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 230, B: 255, A: 255})

	icon := widget.NewIcon(theme.InfoIcon())

	text := widget.NewLabel(
		"These settings will be used as defaults for all new trades. " +
			"Standard bet size = Account Equity × Risk Per Trade.",
	)
	text.Wrapping = fyne.TextWrapWord

	content := container.NewBorder(
		nil, nil,
		container.NewPadded(icon),
		nil,
		container.NewPadded(text),
	)

	return container.NewStack(
		bg,
		container.NewPadded(content),
	)
}

// createSettingsForm creates the settings input form
func (s *Settings) createSettingsForm() fyne.CanvasObject {
	// Account Equity
	accountLabel := widget.NewLabel("Account Equity ($):")
	accountLabel.TextStyle = fyne.TextStyle{Bold: true}

	s.accountEntry = widget.NewEntry()
	s.accountEntry.SetPlaceHolder("e.g., 25000")

	if s.state.Settings != nil {
		s.accountEntry.SetText(fmt.Sprintf("%.0f", s.state.Settings.AccountEquity))
	}

	accountHelp := widget.NewLabel("Your total trading capital")
	accountHelp.TextStyle = fyne.TextStyle{Italic: true}

	// Risk Per Trade
	riskLabel := widget.NewLabel("Risk Per Trade (%):")
	riskLabel.TextStyle = fyne.TextStyle{Bold: true}

	s.riskPercentEntry = widget.NewEntry()
	s.riskPercentEntry.SetPlaceHolder("e.g., 2.80")

	if s.state.Settings != nil {
		s.riskPercentEntry.SetText(fmt.Sprintf("%.2f", s.state.Settings.RiskPerTrade*100))
	}

	riskHelp := widget.NewLabel("Percentage of account to risk per trade (standard conviction)")
	riskHelp.TextStyle = fyne.TextStyle{Italic: true}

	// Preview calculation
	previewLabel := s.createPreviewLabel()

	// Theme selection
	themeLabel := widget.NewLabel("Theme:")
	themeLabel.TextStyle = fyne.TextStyle{Bold: true}

	s.themeSelect = widget.NewSelect(
		[]string{"Day Mode", "Night Mode"},
		func(selected string) {
			// Theme change handler
		},
	)

	if s.state.Settings != nil {
		if s.state.Settings.ThemeMode == "night" {
			s.themeSelect.Selected = "Night Mode"
		} else {
			s.themeSelect.Selected = "Day Mode"
		}
	} else {
		s.themeSelect.Selected = "Day Mode"
	}

	// Add change listeners for preview
	s.accountEntry.OnChanged = func(value string) {
		s.updatePreview(previewLabel)
	}
	s.riskPercentEntry.OnChanged = func(value string) {
		s.updatePreview(previewLabel)
	}

	form := container.NewVBox(
		accountLabel,
		s.accountEntry,
		accountHelp,
		widget.NewSeparator(),

		riskLabel,
		s.riskPercentEntry,
		riskHelp,
		widget.NewSeparator(),

		previewLabel,
		widget.NewSeparator(),

		themeLabel,
		s.themeSelect,
	)

	return form
}

// createPreviewLabel creates the preview calculation label
func (s *Settings) createPreviewLabel() *widget.Label {
	previewLabel := widget.NewLabel("")
	previewLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	previewLabel.Alignment = fyne.TextAlignCenter

	s.updatePreview(previewLabel)

	return previewLabel
}

// updatePreview updates the preview calculation
func (s *Settings) updatePreview(label *widget.Label) {
	// Parse account equity
	accountStr := s.accountEntry.Text
	account, err := strconv.ParseFloat(accountStr, 64)
	if err != nil || account <= 0 {
		label.SetText("📊 Preview: Enter valid account equity")
		return
	}

	// Parse risk percentage
	riskPercentStr := s.riskPercentEntry.Text
	riskPercent, err := strconv.ParseFloat(riskPercentStr, 64)
	if err != nil || riskPercent <= 0 {
		label.SetText("📊 Preview: Enter valid risk percentage")
		return
	}

	// Calculate standard bet size
	standardBet := account * (riskPercent / 100.0)

	label.SetText(fmt.Sprintf(
		"📊 Preview: Standard bet size = $%.2f (at 1.0× conviction)",
		standardBet,
	))
}

// saveSettings saves the settings to AppState and persists to disk
func (s *Settings) saveSettings() {
	if s.state.Settings == nil {
		return
	}

	// Parse and validate account equity
	accountStr := s.accountEntry.Text
	account, err := strconv.ParseFloat(accountStr, 64)
	if err != nil || account <= 0 {
		dialog.ShowError(
			fmt.Errorf("Invalid account equity: %s", accountStr),
			s.window,
		)
		return
	}

	// Parse and validate risk percentage
	riskPercentStr := s.riskPercentEntry.Text
	riskPercent, err := strconv.ParseFloat(riskPercentStr, 64)
	if err != nil || riskPercent <= 0 {
		dialog.ShowError(
			fmt.Errorf("Invalid risk percentage: %s", riskPercentStr),
			s.window,
		)
		return
	}

	// Update settings
	s.state.Settings.AccountEquity = account
	s.state.Settings.RiskPerTrade = riskPercent / 100.0 // Store as decimal

	// Update theme
	if s.themeSelect.Selected == "Night Mode" {
		s.state.Settings.ThemeMode = "night"
	} else {
		s.state.Settings.ThemeMode = "day"
	}

	// Save to disk
	if err := storage.SaveSettings(s.state.Settings); err != nil {
		dialog.ShowError(
			fmt.Errorf("Failed to save settings: %v", err),
			s.window,
		)
		return
	}

	// Show success message
	dialog.ShowInformation(
		"Settings Saved",
		fmt.Sprintf(
			"Settings saved successfully!\n\n"+
				"Account Equity: $%.2f\n"+
				"Risk Per Trade: %.2f%%\n"+
				"Standard Bet: $%.2f",
			account,
			riskPercent,
			account*(riskPercent/100.0),
		),
		s.window,
	)
}
