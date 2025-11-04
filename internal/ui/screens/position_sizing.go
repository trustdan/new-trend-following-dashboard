package screens

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tf-engine/internal/appcore"
)

// PositionSizing represents Screen 5: Position Size Calculator
type PositionSizing struct {
	state  *appcore.AppState
	window fyne.Window

	// Navigation callbacks
	onNext   func()
	onBack   func()
	onCancel func()

	// UI components
	convictionRadio  *widget.RadioGroup
	accountEntry     *widget.Entry
	riskPercentEntry *widget.Entry
	multiplierLabel  *widget.Label
	calculatedRisk   *widget.Label
	explanationLabel *widget.Label
	continueBtn      *widget.Button
}

// NewPositionSizing creates a new position sizing screen
func NewPositionSizing(state *appcore.AppState, window fyne.Window) *PositionSizing {
	return &PositionSizing{
		state:  state,
		window: window,
	}
}

// Validate checks if the screen's data is valid
func (s *PositionSizing) Validate() bool {
	// Conviction must be selected and position sizing calculated
	return s.state.CurrentTrade != nil &&
		s.state.CurrentTrade.Conviction >= 5 &&
		s.state.CurrentTrade.Conviction <= 8 &&
		s.state.CurrentTrade.SizingMultiplier > 0
}

// GetName returns the screen name
func (s *PositionSizing) GetName() string {
	return "position_sizing"
}

// SetNavCallbacks sets navigation callback functions
func (s *PositionSizing) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
	// Wrap error-returning callbacks to match internal func() signature
	s.onNext = func() {
		if onNext != nil {
			onNext()
		}
	}
	s.onBack = func() {
		if onBack != nil {
			onBack()
		}
	}
	s.onCancel = onCancel
}

// Render renders the position sizing UI
func (s *PositionSizing) Render() fyne.CanvasObject {
	// Header
	header := s.createHeader()

	// Info banner
	infoBanner := s.createInfoBanner()

	// Main form
	form := s.createForm()

	// Navigation buttons
	navButtons := s.createNavigationButtons()

	// Scrollable content
	scrollContent := container.NewVBox(
		infoBanner,
		form,
	)

	scroll := container.NewScroll(scrollContent)

	// Layout
	content := container.NewBorder(
		header,
		navButtons,
		nil,
		nil,
		scroll,
	)

	return content
}

// createHeader creates the screen header
func (s *PositionSizing) createHeader() fyne.CanvasObject {
	progress := widget.NewLabel("Step 5 of 8")
	progress.TextStyle = fyne.TextStyle{Italic: true}

	title := widget.NewLabel("Screen 5: Position Sizing Calculator")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Use poker-bet sizing based on trade conviction")
	subtitle.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		progress,
		title,
		subtitle,
		widget.NewSeparator(),
	)
}

// createInfoBanner creates the poker-bet sizing explanation banner
func (s *PositionSizing) createInfoBanner() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 230, B: 255, A: 255})

	icon := widget.NewIcon(theme.InfoIcon())

	text := widget.NewLabel(
		"🎲 Poker-Bet Sizing: Rate your conviction in this trade (5-8). " +
			"Higher conviction = larger position size. " +
			"Standard sizing is 7 (1.0× multiplier).",
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

// createForm creates the main position sizing form
func (s *PositionSizing) createForm() fyne.CanvasObject {
	// Conviction rating section
	convictionSection := s.createConvictionSection()

	// Account settings section
	accountSection := s.createAccountSection()

	// Calculation results section
	resultsSection := s.createResultsSection()

	return container.NewVBox(
		convictionSection,
		widget.NewSeparator(),
		accountSection,
		widget.NewSeparator(),
		resultsSection,
	)
}

// createConvictionSection creates the conviction rating selector
func (s *PositionSizing) createConvictionSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("💡 Trade Conviction Rating")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	// Conviction radio group
	s.convictionRadio = widget.NewRadioGroup(
		[]string{
			"5 - Weak conviction (0.5× size)",
			"6 - Below average (0.75× size)",
			"7 - Standard conviction (1.0× size)",
			"8 - Strong conviction (1.25× size)",
		},
		func(value string) {
			s.onConvictionChanged(value)
		},
	)

	// Restore previous selection if exists
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.Conviction >= 5 && s.state.CurrentTrade.Conviction <= 8 {
		convictionIdx := s.state.CurrentTrade.Conviction - 5
		if convictionIdx >= 0 && convictionIdx < len(s.convictionRadio.Options) {
			s.convictionRadio.Selected = s.convictionRadio.Options[convictionIdx]
		}
	} else {
		// Default to 7 (standard)
		s.convictionRadio.Selected = s.convictionRadio.Options[2]
	}

	explanation := widget.NewLabel(
		"Select your confidence level for this trade setup. Be honest - " +
			"overconfidence leads to oversizing and increased losses.",
	)
	explanation.Wrapping = fyne.TextWrapWord
	explanation.TextStyle = fyne.TextStyle{Italic: true}

	s.multiplierLabel = widget.NewLabel("")
	s.multiplierLabel.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		sectionTitle,
		s.convictionRadio,
		explanation,
		s.multiplierLabel,
	)
}

// createAccountSection creates the account size and risk percentage inputs
func (s *PositionSizing) createAccountSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("💰 Account Settings")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	// Account equity
	accountLabel := widget.NewLabel("Account Equity ($):")
	s.accountEntry = widget.NewEntry()
	s.accountEntry.SetPlaceHolder("e.g., 100000")

	// Restore previous value or use default
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.AccountEquity > 0 {
		s.accountEntry.SetText(fmt.Sprintf("%.0f", s.state.CurrentTrade.AccountEquity))
	} else {
		s.accountEntry.SetText("100000") // Default $100k
	}

	s.accountEntry.OnChanged = func(value string) {
		s.updateCalculation()
	}

	// Risk per trade percentage
	riskLabel := widget.NewLabel("Risk Per Trade (%):")
	s.riskPercentEntry = widget.NewEntry()
	s.riskPercentEntry.SetPlaceHolder("e.g., 0.75")

	// Use value from policy or previous trade
	riskPercent := s.state.Policy.Defaults.RiskPerTrade * 100 // Convert to percentage
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.RiskPerTrade > 0 {
		riskPercent = s.state.CurrentTrade.RiskPerTrade * 100
	}
	s.riskPercentEntry.SetText(fmt.Sprintf("%.2f", riskPercent))

	s.riskPercentEntry.OnChanged = func(value string) {
		s.updateCalculation()
	}

	accountGrid := container.NewGridWithColumns(2,
		accountLabel, s.accountEntry,
		riskLabel, s.riskPercentEntry,
	)

	return container.NewVBox(
		sectionTitle,
		accountGrid,
	)
}

// createResultsSection creates the calculation results display
func (s *PositionSizing) createResultsSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("📊 Calculated Risk")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	s.calculatedRisk = widget.NewLabel("")
	s.calculatedRisk.TextStyle = fyne.TextStyle{Bold: true}
	s.calculatedRisk.Alignment = fyne.TextAlignCenter

	s.explanationLabel = widget.NewLabel("")
	s.explanationLabel.Wrapping = fyne.TextWrapWord
	s.explanationLabel.TextStyle = fyne.TextStyle{Italic: true}
	s.explanationLabel.Alignment = fyne.TextAlignCenter

	// Initial calculation
	s.updateCalculation()

	resultCard := container.NewVBox(
		s.calculatedRisk,
		s.explanationLabel,
	)

	// Card with colored background
	bg := canvas.NewRectangle(color.RGBA{R: 240, G: 248, B: 255, A: 255})
	cardContent := container.NewStack(
		bg,
		container.NewPadded(resultCard),
	)

	return container.NewVBox(
		sectionTitle,
		cardContent,
	)
}

// createNavigationButtons creates the navigation button bar
func (s *PositionSizing) createNavigationButtons() fyne.CanvasObject {
	backBtn := widget.NewButton("← Back", func() {
		s.savePositionSizing()
		if s.onBack != nil {
			s.onBack()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		s.savePositionSizing()
		if s.onCancel != nil {
			s.onCancel()
		}
	})

	s.continueBtn = widget.NewButton("Continue →", func() {
		if s.Validate() {
			s.savePositionSizing()
			if s.onNext != nil {
				s.onNext()
			}
		}
	})
	s.continueBtn.Importance = widget.HighImportance

	// Initial validation state
	s.updateContinueButton()

	return container.NewBorder(
		nil,
		nil,
		container.NewHBox(backBtn, cancelBtn),
		s.continueBtn,
		nil,
	)
}

// onConvictionChanged handles conviction rating changes
func (s *PositionSizing) onConvictionChanged(value string) {
	// Extract conviction number from string like "7 - Standard conviction (1.0× size)"
	if len(value) == 0 {
		return
	}

	convictionStr := string(value[0])
	conviction, err := strconv.Atoi(convictionStr)
	if err != nil {
		return
	}

	// Get multiplier from policy
	multiplier := s.getMultiplier(conviction)

	// Update multiplier label
	s.multiplierLabel.SetText(fmt.Sprintf("Sizing Multiplier: %.2f×", multiplier))

	// Update calculation
	s.updateCalculation()
	s.updateContinueButton()
}

// getMultiplier returns the poker sizing multiplier for a given conviction level
func (s *PositionSizing) getMultiplier(conviction int) float64 {
	convictionKey := fmt.Sprintf("%d", conviction)
	if multiplier, exists := s.state.Policy.Checklist.PokerSizing[convictionKey]; exists {
		return multiplier
	}
	return 1.0 // Default to standard sizing
}

// updateCalculation updates the calculated risk display
func (s *PositionSizing) updateCalculation() {
	// Parse account equity
	accountStr := s.accountEntry.Text
	account, err := strconv.ParseFloat(accountStr, 64)
	if err != nil || account <= 0 {
		s.calculatedRisk.SetText("❌ Invalid account equity")
		return
	}

	// Parse risk percentage
	riskPercentStr := s.riskPercentEntry.Text
	riskPercent, err := strconv.ParseFloat(riskPercentStr, 64)
	if err != nil || riskPercent <= 0 {
		s.calculatedRisk.SetText("❌ Invalid risk percentage")
		return
	}

	// Get conviction level
	conviction := s.getSelectedConviction()
	if conviction == 0 {
		s.calculatedRisk.SetText("❌ Select conviction level")
		return
	}

	// Get multiplier
	multiplier := s.getMultiplier(conviction)

	// Calculate risk amount
	baseRisk := account * (riskPercent / 100.0)
	adjustedRisk := baseRisk * multiplier

	// Update display
	s.calculatedRisk.SetText(fmt.Sprintf("Risk Amount: $%.2f", adjustedRisk))

	s.explanationLabel.SetText(fmt.Sprintf(
		"Base risk: $%.2f (%.2f%% of $%.0f) × %.2f× conviction multiplier = $%.2f total risk",
		baseRisk, riskPercent, account, multiplier, adjustedRisk,
	))
}

// getSelectedConviction returns the currently selected conviction level (5-8)
func (s *PositionSizing) getSelectedConviction() int {
	if s.convictionRadio == nil || s.convictionRadio.Selected == "" {
		return 0
	}

	// Extract conviction from selected option
	selected := s.convictionRadio.Selected
	if len(selected) == 0 {
		return 0
	}

	convictionStr := string(selected[0])
	conviction, err := strconv.Atoi(convictionStr)
	if err != nil {
		return 0
	}

	return conviction
}

// updateContinueButton enables/disables continue button based on validation
func (s *PositionSizing) updateContinueButton() {
	if s.continueBtn == nil {
		return
	}

	if s.Validate() {
		s.continueBtn.Enable()
	} else {
		s.continueBtn.Disable()
	}
}

// savePositionSizing saves the position sizing data to the trade
func (s *PositionSizing) savePositionSizing() {
	if s.state.CurrentTrade == nil {
		return
	}

	// Parse and save account equity
	if accountStr := s.accountEntry.Text; accountStr != "" {
		if account, err := strconv.ParseFloat(accountStr, 64); err == nil {
			s.state.CurrentTrade.AccountEquity = account
		}
	}

	// Parse and save risk per trade
	if riskStr := s.riskPercentEntry.Text; riskStr != "" {
		if riskPercent, err := strconv.ParseFloat(riskStr, 64); err == nil {
			s.state.CurrentTrade.RiskPerTrade = riskPercent / 100.0 // Store as decimal
		}
	}

	// Save conviction and multiplier
	conviction := s.getSelectedConviction()
	if conviction >= 5 && conviction <= 8 {
		s.state.CurrentTrade.Conviction = conviction
		s.state.CurrentTrade.SizingMultiplier = s.getMultiplier(conviction)

		// Calculate and save max loss
		baseRisk := s.state.CurrentTrade.AccountEquity * s.state.CurrentTrade.RiskPerTrade
		s.state.CurrentTrade.MaxLoss = baseRisk * s.state.CurrentTrade.SizingMultiplier
	}
}
