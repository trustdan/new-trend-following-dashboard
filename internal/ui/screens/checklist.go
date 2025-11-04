package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tf-engine/internal/appcore"
	"tf-engine/internal/widgets"
)

// Checklist represents Screen 4: Anti-Impulsivity Checklist
type Checklist struct {
	state  *appcore.AppState
	window fyne.Window

	// Navigation callbacks
	onNext   func()
	onBack   func()
	onCancel func()

	// UI components
	requiredChecks map[string]*widget.Check
	optionalChecks map[string]*widget.Check
	cooldownTimer  *widgets.CooldownTimer
	continueBtn    *widget.Button
	validationMsg  *widget.Label
}

// ChecklistLabels maps policy IDs to human-readable labels and descriptions
var ChecklistLabels = map[string]struct {
	Label       string
	Description string
}{
	// Required items
	"SIG_REQ": {
		Label:       "Signal Requirements Met",
		Description: "Price above SMA200 with confirmed Donchian breakout or strategy-specific entry signal",
	},
	"RISK_REQ": {
		Label:       "Risk Parameters Acceptable",
		Description: "Stop-loss placement is acceptable for account size and risk tolerance",
	},
	"OPT_REQ": {
		Label:       "Options Alignment Verified",
		Description: "Strategy hold duration matches options expiration window (avoid theta decay)",
	},
	"EXIT_REQ": {
		Label:       "Exit Plan Defined",
		Description: "Clear profit targets and stop-loss levels defined before entry",
	},
	"BEHAV_REQ": {
		Label:       "Behavioral Check Passed",
		Description: "Emotionally calm, not chasing price, following systematic process",
	},
	// Optional items
	"REGIME_OK": {
		Label:       "Market Regime Favorable",
		Description: "SPY/QQQ above SMA200, sector strength confirmed in universe screener",
	},
	"NO_CHASE": {
		Label:       "Not Chasing Price",
		Description: "Entry within acceptable distance from signal (not late to party)",
	},
	"JOURNAL_DONE": {
		Label:       "Trade Documented",
		Description: "Entry rationale and setup notes captured in trading journal",
	},
}

// NewChecklist creates a new checklist screen
func NewChecklist(state *appcore.AppState, window fyne.Window) *Checklist {
	return &Checklist{
		state:          state,
		window:         window,
		requiredChecks: make(map[string]*widget.Check),
		optionalChecks: make(map[string]*widget.Check),
	}
}

// Validate checks if the screen's data is valid
func (s *Checklist) Validate() bool {
	// All required checkboxes must be checked
	for _, check := range s.requiredChecks {
		if !check.Checked {
			return false
		}
	}

	// Cooldown must be complete
	if s.cooldownTimer != nil && !s.cooldownTimer.IsComplete() {
		return false
	}

	return true
}

// GetName returns the screen name
func (s *Checklist) GetName() string {
	return "checklist"
}

// SetNavCallbacks sets navigation callback functions
func (s *Checklist) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
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

// Render renders the checklist UI
func (s *Checklist) Render() fyne.CanvasObject {
	// Header
	header := s.createHeader()

	// Warning banner
	warningBanner := s.createWarningBanner()

	// Cooldown timer
	cooldownSection := s.createCooldownSection()

	// Checklist sections
	requiredSection := s.createRequiredSection()
	optionalSection := s.createOptionalSection()

	// Validation message
	s.validationMsg = widget.NewLabel("")
	s.validationMsg.TextStyle = fyne.TextStyle{Italic: true}
	s.validationMsg.Alignment = fyne.TextAlignCenter

	// Navigation buttons
	navButtons := s.createNavigationButtons()

	// Scrollable content
	scrollContent := container.NewVBox(
		warningBanner,
		cooldownSection,
		widget.NewSeparator(),
		requiredSection,
		widget.NewSeparator(),
		optionalSection,
		widget.NewSeparator(),
		s.validationMsg,
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
func (s *Checklist) createHeader() fyne.CanvasObject {
	progress := widget.NewLabel("Step 4 of 8")
	progress.TextStyle = fyne.TextStyle{Italic: true}

	title := widget.NewLabel("Screen 4: Anti-Impulsivity Checklist")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Complete all required criteria before proceeding")
	subtitle.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		progress,
		title,
		subtitle,
		widget.NewSeparator(),
	)
}

// createWarningBanner creates the behavioral finance warning banner
func (s *Checklist) createWarningBanner() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 255, G: 200, B: 100, A: 255})

	icon := widget.NewIcon(theme.WarningIcon())

	text := widget.NewLabel(
		"⚠️ Behavioral Finance Guardrail: This checklist prevents impulsive trading. " +
			"You must complete all 5 required items AND wait for the cooldown timer. " +
			"This is non-negotiable.",
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

// createCooldownSection creates the cooldown timer display
func (s *Checklist) createCooldownSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("⏱️ Cooldown Timer")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	// Check if cooldown was started in previous screen
	if s.state.CurrentTrade != nil && !s.state.CurrentTrade.CooldownStartTime.IsZero() {
		// Calculate cooldown duration from policy
		cooldownDuration := time.Duration(s.state.Policy.Defaults.CooldownSeconds) * time.Second

		// Create timer from start time
		s.cooldownTimer = widgets.NewCooldownTimerFromTime(
			cooldownDuration,
			s.state.CurrentTrade.CooldownStartTime,
			func() {
				// On complete, update validation state
				s.updateValidation()
			},
		)
	} else {
		// No cooldown started - show error message
		errorLabel := widget.NewLabel("⚠️ Cooldown not started. Return to Ticker Entry screen.")
		errorLabel.TextStyle = fyne.TextStyle{Bold: true}
		return container.NewVBox(
			sectionTitle,
			errorLabel,
		)
	}

	explanation := widget.NewLabel(
		"You must wait for the cooldown to complete before proceeding. " +
			"This enforces deliberate decision-making and prevents emotional trades.",
	)
	explanation.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		sectionTitle,
		s.cooldownTimer,
		explanation,
	)
}

// createRequiredSection creates the required checklist items
func (s *Checklist) createRequiredSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("✓ Required Items (All 5 Must Be Checked)")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	items := container.NewVBox()

	// Get required items from policy
	requiredItems := s.state.Policy.Checklist.Required

	for _, itemID := range requiredItems {
		checkItem := s.createChecklistItem(itemID, true)
		items.Add(checkItem)
	}

	return container.NewVBox(
		sectionTitle,
		items,
	)
}

// createOptionalSection creates the optional checklist items
func (s *Checklist) createOptionalSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("➕ Optional Items (Recommended But Not Required)")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	items := container.NewVBox()

	// Get optional items from policy
	optionalItems := s.state.Policy.Checklist.Optional

	for _, itemID := range optionalItems {
		checkItem := s.createChecklistItem(itemID, false)
		items.Add(checkItem)
	}

	return container.NewVBox(
		sectionTitle,
		items,
	)
}

// createChecklistItem creates a single checklist item with label and description
func (s *Checklist) createChecklistItem(itemID string, required bool) fyne.CanvasObject {
	labelInfo, exists := ChecklistLabels[itemID]
	if !exists {
		labelInfo = struct {
			Label       string
			Description string
		}{
			Label:       itemID,
			Description: "No description available",
		}
	}

	// Create checkbox
	check := widget.NewCheck(labelInfo.Label, func(checked bool) {
		// Update validation state when checkbox changes
		s.updateValidation()
	})

	// Restore checked state if it exists
	if s.state.CurrentTrade != nil {
		if required {
			if checkedItems, exists := s.state.CurrentTrade.ChecklistRequired[itemID]; exists {
				check.Checked = checkedItems
			}
		} else {
			if checkedItems, exists := s.state.CurrentTrade.ChecklistOptional[itemID]; exists {
				check.Checked = checkedItems
			}
		}
	}

	// Store checkbox reference
	if required {
		s.requiredChecks[itemID] = check
	} else {
		s.optionalChecks[itemID] = check
	}

	// Description label
	description := widget.NewLabel(labelInfo.Description)
	description.Wrapping = fyne.TextWrapWord
	description.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		check,
		container.NewPadded(description),
	)
}

// createNavigationButtons creates the navigation button bar
func (s *Checklist) createNavigationButtons() fyne.CanvasObject {
	backBtn := widget.NewButton("← Back", func() {
		s.saveChecklistState()
		if s.onBack != nil {
			s.onBack()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		s.saveChecklistState()
		if s.onCancel != nil {
			s.onCancel()
		}
	})

	s.continueBtn = widget.NewButton("Continue →", func() {
		if s.Validate() {
			s.saveChecklistState()
			if s.onNext != nil {
				s.onNext()
			}
		}
	})
	s.continueBtn.Importance = widget.HighImportance

	// Initial validation state
	s.updateValidation()

	return container.NewBorder(
		nil,
		nil,
		container.NewHBox(backBtn, cancelBtn),
		s.continueBtn,
		nil,
	)
}

// updateValidation updates the continue button state and validation message
func (s *Checklist) updateValidation() {
	isValid := s.Validate()
	s.continueBtn.Disabled = !isValid

	if !isValid {
		// Build specific error message
		uncheckedCount := 0
		for _, check := range s.requiredChecks {
			if !check.Checked {
				uncheckedCount++
			}
		}

		cooldownComplete := s.cooldownTimer != nil && s.cooldownTimer.IsComplete()

		if uncheckedCount > 0 && !cooldownComplete {
			s.validationMsg.SetText(fmt.Sprintf(
				"❌ %d required items unchecked • Cooldown in progress",
				uncheckedCount,
			))
		} else if uncheckedCount > 0 {
			s.validationMsg.SetText(fmt.Sprintf(
				"❌ %d required items must be checked",
				uncheckedCount,
			))
		} else if !cooldownComplete {
			s.validationMsg.SetText("⏱️ Waiting for cooldown to complete...")
		}
	} else {
		s.validationMsg.SetText("✓ All requirements met - Ready to continue")
	}

	s.continueBtn.Refresh()
}

// saveChecklistState saves the current checkbox states to the trade
func (s *Checklist) saveChecklistState() {
	if s.state.CurrentTrade == nil {
		return
	}

	// Initialize maps if needed
	if s.state.CurrentTrade.ChecklistRequired == nil {
		s.state.CurrentTrade.ChecklistRequired = make(map[string]bool)
	}
	if s.state.CurrentTrade.ChecklistOptional == nil {
		s.state.CurrentTrade.ChecklistOptional = make(map[string]bool)
	}

	// Save required items
	for itemID, check := range s.requiredChecks {
		s.state.CurrentTrade.ChecklistRequired[itemID] = check.Checked
	}

	// Save optional items
	for itemID, check := range s.optionalChecks {
		s.state.CurrentTrade.ChecklistOptional[itemID] = check.Checked
	}

	// Mark checklist as passed if all required items checked
	s.state.CurrentTrade.ChecklistPassed = s.Validate()
}
