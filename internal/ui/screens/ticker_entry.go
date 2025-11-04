package screens

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// TickerEntry represents Screen 3: Ticker Entry & Strategy Selection
type TickerEntry struct {
	state  *appcore.AppState
	window fyne.Window

	// Navigation callbacks
	onNext   func()
	onBack   func()
	onCancel func()

	// UI components
	tickerEntry    *widget.Entry
	strategySelect *widget.Select
	strategyInfo   *fyne.Container
	continueBtn    *widget.Button

	// Phase 6: Warning system components
	warningBanner *fyne.Container
	ackCheckbox   *widget.Check
}

// NewTickerEntry creates a new ticker entry screen
func NewTickerEntry(state *appcore.AppState, window fyne.Window) *TickerEntry {
	return &TickerEntry{
		state:  state,
		window: window,
	}
}

// Validate checks if the screen's data is valid
func (s *TickerEntry) Validate() bool {
	// Ticker and strategy must be selected
	return s.state.CurrentTrade != nil &&
		s.state.CurrentTrade.Ticker != "" &&
		s.state.CurrentTrade.Strategy != ""
}

// GetName returns the screen name
func (s *TickerEntry) GetName() string {
	return "ticker_entry"
}

// SetNavCallbacks sets navigation callback functions
func (s *TickerEntry) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
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

// Render renders the ticker entry UI
func (s *TickerEntry) Render() fyne.CanvasObject {
	// Header
	header := s.createHeader()

	// Info banner
	infoBanner := s.createInfoBanner()

	// Main form
	form := s.createForm()

	// Navigation buttons
	navButtons := s.createNavigationButtons()

	// Layout
	content := container.NewBorder(
		container.NewVBox(header, infoBanner),
		navButtons,
		nil,
		nil,
		container.NewPadded(form),
	)

	return content
}

// createHeader creates the screen header
func (s *TickerEntry) createHeader() fyne.CanvasObject {
	progress := widget.NewLabel("Step 3 of 8")
	progress.TextStyle = fyne.TextStyle{Italic: true}

	title := widget.NewLabel("Screen 3: Enter Ticker & Strategy")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Select ticker symbol and trading strategy")
	subtitle.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		progress,
		title,
		subtitle,
		widget.NewSeparator(),
	)
}

// createInfoBanner creates the information banner
func (s *TickerEntry) createInfoBanner() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 230, B: 255, A: 255})

	icon := widget.NewIcon(theme.InfoIcon())

	// Get sector name for context
	sectorName := "Unknown"
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.Sector != "" {
		sectorName = s.state.CurrentTrade.Sector
	}

	text := widget.NewLabel(
		fmt.Sprintf("ALL strategies are shown with color-coded suitability for %s sector. "+
			"🟢 = Excellent/Good | 🟡 = Marginal (requires acknowledgement) | 🔴 = Incompatible (strong warning). "+
			"Enter a ticker you found in the FINVIZ screener. "+
			"A 120-second cooldown will start when you proceed.",
			sectorName),
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

// createForm creates the main ticker/strategy form
func (s *TickerEntry) createForm() fyne.CanvasObject {
	// Ticker input section
	tickerLabel := widget.NewLabel("Ticker Symbol:")
	tickerLabel.TextStyle = fyne.TextStyle{Bold: true}

	s.tickerEntry = widget.NewEntry()
	s.tickerEntry.SetPlaceHolder("e.g., UNH, MSFT, CAT")

	// Convert to uppercase on change
	s.tickerEntry.OnChanged = func(value string) {
		upper := strings.ToUpper(value)
		if upper != value {
			s.tickerEntry.SetText(upper)
			return
		}

		// Update trade state
		if s.state.CurrentTrade != nil {
			s.state.CurrentTrade.Ticker = upper
			s.updateContinueButton()
		}
	}

	tickerHelp := widget.NewLabel("Ticker from FINVIZ screener (1-5 characters)")
	tickerHelp.TextStyle = fyne.TextStyle{Italic: true}

	// Strategy dropdown section
	strategyLabel := widget.NewLabel("Pine Script Strategy:")
	strategyLabel.TextStyle = fyne.TextStyle{Bold: true}

	s.strategySelect = s.createStrategyDropdown()

	strategyHelp := widget.NewLabel("ALL strategies shown - green = good, yellow/red require acknowledgement")
	strategyHelp.TextStyle = fyne.TextStyle{Italic: true}

	// Strategy metadata display (initially hidden)
	s.strategyInfo = container.NewVBox()

	// Phase 6: Warning banner (initially hidden)
	s.warningBanner = container.NewVBox()

	// Phase 6: Acknowledgement checkbox (initially hidden)
	s.ackCheckbox = widget.NewCheck("", func(checked bool) {
		s.updateContinueButton()
	})
	s.ackCheckbox.Hide()

	// Combine form elements
	form := container.NewVBox(
		tickerLabel,
		s.tickerEntry,
		tickerHelp,
		layout.NewSpacer(),
		strategyLabel,
		s.strategySelect,
		strategyHelp,
		layout.NewSpacer(),
		s.strategyInfo,
		s.warningBanner,
		s.ackCheckbox,
	)

	return form
}

// createStrategyDropdown creates strategy dropdown showing ALL strategies with color indicators
func (s *TickerEntry) createStrategyDropdown() *widget.Select {
	strategies := s.getAllStrategiesWithIndicators()

	dropdown := widget.NewSelect(strategies, func(value string) {
		s.onStrategySelected(value)
	})

	dropdown.PlaceHolder = "Select a strategy..."

	return dropdown
}

// getAllStrategiesWithIndicators returns ALL strategies with color-coded suitability indicators
func (s *TickerEntry) getAllStrategiesWithIndicators() []string {
	if s.state.CurrentTrade == nil || s.state.CurrentTrade.Sector == "" {
		return []string{}
	}

	var strategyLabels []string

	// Show ALL strategies from policy with color indicators
	for stratID, strategy := range s.state.Policy.Strategies {
		suitability := s.getSuitability(stratID, s.state.CurrentTrade.Sector)
		indicator := s.getColorIndicator(suitability.Color)

		// Format: "🟢 Alt10 - Profit Targets (3N/6N/9N)"
		label := fmt.Sprintf("%s %s - %s", indicator, stratID, strategy.Label)
		strategyLabels = append(strategyLabels, label)
	}

	return strategyLabels
}

// getSuitability returns the suitability rating for a strategy in a given sector
func (s *TickerEntry) getSuitability(strategyID, sector string) models.StrategySuitability {
	// Find the sector in policy
	for _, sec := range s.state.Policy.Sectors {
		if sec.Name == sector {
			// Check if strategy_suitability exists for this strategy
			if suitability, exists := sec.StrategySuitability[strategyID]; exists {
				return suitability
			}
		}
	}

	// Default to marginal if not found (requires acknowledgement)
	return models.StrategySuitability{
		Rating:                 "marginal",
		Color:                  "yellow",
		Rationale:              "Suitability data not available for this sector/strategy combination",
		RequireAcknowledgement: true,
	}
}

// getColorIndicator returns the emoji indicator for a color
func (s *TickerEntry) getColorIndicator(colorName string) string {
	switch colorName {
	case "green":
		return "🟢"
	case "yellow":
		return "🟡"
	case "red":
		return "🔴"
	default:
		return "⚪"
	}
}

// onStrategySelected handles strategy selection
func (s *TickerEntry) onStrategySelected(value string) {
	if value == "" {
		s.strategyInfo.Objects = nil
		s.strategyInfo.Refresh()
		s.hideWarningBanner()
		s.ackCheckbox.Hide()
		return
	}

	// Extract strategy ID from "🟢 Alt10 - Profit Targets" format
	// Remove color indicator emoji first
	value = strings.TrimPrefix(value, "🟢 ")
	value = strings.TrimPrefix(value, "🟡 ")
	value = strings.TrimPrefix(value, "🔴 ")
	value = strings.TrimPrefix(value, "⚪ ")

	parts := strings.Split(value, " - ")
	if len(parts) == 0 {
		return
	}
	strategyID := strings.TrimSpace(parts[0])

	// Get suitability for this strategy/sector combination
	suitability := s.getSuitability(strategyID, s.state.CurrentTrade.Sector)

	// Update trade state
	if s.state.CurrentTrade != nil {
		s.state.CurrentTrade.Strategy = strategyID
		s.state.CurrentTrade.StrategySuitability = suitability.Rating
		s.state.CurrentTrade.StrategyWarningAcknowledged = false // Reset acknowledgement
	}

	// Display strategy metadata
	s.displayStrategyMetadata(strategyID)

	// Phase 6: Show warning banner and acknowledgement for yellow/red strategies
	if suitability.RequireAcknowledgement {
		s.showWarningBanner(suitability)
		s.ackCheckbox.Show()
		s.ackCheckbox.SetChecked(false)
		s.continueBtn.Disable()

		// Log warning display
		fmt.Printf("⚠️  Strategy warning displayed: %s in %s (rating: %s)\n",
			strategyID, s.state.CurrentTrade.Sector, suitability.Rating)
	} else {
		// Green strategy - no warning needed
		s.hideWarningBanner()
		s.ackCheckbox.Hide()
		s.updateContinueButton()
	}
}

// displayStrategyMetadata shows strategy details below dropdown
func (s *TickerEntry) displayStrategyMetadata(strategyID string) {
	strategy, exists := s.state.Policy.Strategies[strategyID]
	if !exists {
		return
	}

	// Create metadata card
	bg := canvas.NewRectangle(color.RGBA{R: 240, G: 255, B: 240, A: 255})

	// Left border (green)
	border := canvas.NewRectangle(color.RGBA{R: 0, G: 180, B: 80, A: 255})
	border.SetMinSize(fyne.NewSize(4, 100))

	// Metadata content
	label := widget.NewLabel(fmt.Sprintf("📋 %s", strategy.Label))
	label.TextStyle = fyne.TextStyle{Bold: true}

	suitability := widget.NewLabel(fmt.Sprintf("Options Suitability: %s", strategy.OptionsSuitability))

	holdWeeks := widget.NewLabel(fmt.Sprintf("Typical Hold: %s weeks", strategy.HoldWeeks))

	notes := widget.NewLabel(strategy.Notes)
	notes.Wrapping = fyne.TextWrapWord
	notes.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(
		label,
		suitability,
		holdWeeks,
		widget.NewSeparator(),
		notes,
	)

	card := container.NewStack(
		bg,
		container.NewBorder(nil, nil, border, nil,
			container.NewPadded(content),
		),
	)

	// Update strategy info container
	s.strategyInfo.Objects = []fyne.CanvasObject{card}
	s.strategyInfo.Refresh()
}

// createNavigationButtons creates Continue/Back/Cancel buttons
func (s *TickerEntry) createNavigationButtons() fyne.CanvasObject {
	// Back button
	backBtn := widget.NewButton("← Back to Screener", func() {
		if s.onBack != nil {
			s.onBack()
		}
	})

	// Cancel button
	cancelBtn := widget.NewButton("Cancel", func() {
		if s.onCancel != nil {
			s.onCancel()
		}
	})

	// Continue button (starts cooldown)
	s.continueBtn = widget.NewButton("Continue to Checklist →", func() {
		s.startCooldownAndProceed()
	})
	s.continueBtn.Importance = widget.HighImportance
	s.continueBtn.Disable() // Initially disabled

	return container.NewBorder(
		widget.NewSeparator(),
		nil,
		nil,
		nil,
		container.NewHBox(
			backBtn,
			cancelBtn,
			layout.NewSpacer(),
			s.continueBtn,
		),
	)
}

// updateContinueButton enables/disables continue button based on validation
func (s *TickerEntry) updateContinueButton() {
	if s.continueBtn == nil {
		return
	}

	// Phase 6: Check if acknowledgement is required
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.Strategy != "" {
		suitability := s.getSuitability(s.state.CurrentTrade.Strategy, s.state.CurrentTrade.Sector)

		if suitability.RequireAcknowledgement {
			// Yellow/red strategy - require acknowledgement checkbox
			if s.Validate() && s.ackCheckbox.Checked {
				s.continueBtn.Enable()
			} else {
				s.continueBtn.Disable()
			}
			return
		}
	}

	// Green strategy or no strategy selected - normal validation
	if s.Validate() {
		s.continueBtn.Enable()
	} else {
		s.continueBtn.Disable()
	}
}

// startCooldownAndProceed starts the cooldown timer and proceeds to next screen
func (s *TickerEntry) startCooldownAndProceed() {
	if !s.Validate() {
		s.showError("Please enter ticker and select strategy")
		return
	}

	// Phase 6: Log warning acknowledgement if applicable
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.Strategy != "" {
		suitability := s.getSuitability(s.state.CurrentTrade.Strategy, s.state.CurrentTrade.Sector)

		if suitability.RequireAcknowledgement && s.ackCheckbox.Checked {
			s.state.CurrentTrade.StrategyWarningAcknowledged = true

			// Log warning acknowledgement
			fmt.Printf("⚠️  Strategy warning ACKNOWLEDGED: %s in %s (rating: %s)\n",
				s.state.CurrentTrade.Strategy,
				s.state.CurrentTrade.Sector,
				suitability.Rating)
		}
	}

	// Start cooldown in app state (120 seconds hardcoded in AppState)
	s.state.StartCooldown()

	// Log cooldown start
	fmt.Printf("✓ Cooldown started: 120 seconds for %s (%s)\n",
		s.state.CurrentTrade.Ticker,
		s.state.CurrentTrade.Strategy,
	)

	// Proceed to next screen
	if s.onNext != nil {
		s.onNext()
	}
}

// showWarningBanner displays a warning banner for yellow/red strategies
func (s *TickerEntry) showWarningBanner(suitability models.StrategySuitability) {
	// Determine warning severity color
	var bgColor color.Color
	var icon string
	var title string

	if suitability.Color == "red" {
		bgColor = color.RGBA{R: 255, G: 200, B: 200, A: 255} // Light red
		icon = "🔴"
		title = "INCOMPATIBLE STRATEGY WARNING"
	} else {
		bgColor = color.RGBA{R: 255, G: 240, B: 200, A: 255} // Light yellow/amber
		icon = "🟡"
		title = "MARGINAL STRATEGY WARNING"
	}

	bg := canvas.NewRectangle(bgColor)

	// Warning content
	titleLabel := widget.NewLabel(fmt.Sprintf("%s %s", icon, title))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	rationaleLabel := widget.NewLabel(suitability.Rationale)
	rationaleLabel.Wrapping = fyne.TextWrapWord

	// Acknowledgement text for checkbox
	s.ackCheckbox.Text = fmt.Sprintf("I acknowledge this strategy is %s for this sector and understand the risks", suitability.Rating)
	s.ackCheckbox.Refresh()

	content := container.NewVBox(
		titleLabel,
		rationaleLabel,
	)

	banner := container.NewStack(
		bg,
		container.NewPadded(content),
	)

	s.warningBanner.Objects = []fyne.CanvasObject{banner}
	s.warningBanner.Refresh()
}

// hideWarningBanner hides the warning banner
func (s *TickerEntry) hideWarningBanner() {
	s.warningBanner.Objects = nil
	s.warningBanner.Refresh()
}

// showError displays an error message (console for now)
func (s *TickerEntry) showError(message string) {
	fmt.Printf("❌ Error: %s\n", message)
}
