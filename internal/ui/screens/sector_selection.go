package screens

import (
	"fmt"
	"image/color"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// SectorSelection represents Screen 1: Sector Selection
type SectorSelection struct {
	state          *appcore.AppState
	window         fyne.Window
	onNext         func() error
	onBack         func() error
	onCancel       func()
	continueBtn    *widget.Button
	selectedSector string
}

// NewSectorSelection creates a new sector selection screen
func NewSectorSelection(state *appcore.AppState, window fyne.Window) *SectorSelection {
	return &SectorSelection{
		state:  state,
		window: window,
	}
}

// SetNavCallbacks sets navigation callback functions
func (s *SectorSelection) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
	s.onNext = onNext
	s.onBack = onBack
	s.onCancel = onCancel
}

// Validate checks if the screen's data is valid
func (s *SectorSelection) Validate() bool {
	// Sector must be selected
	return s.state.CurrentTrade != nil && s.state.CurrentTrade.Sector != ""
}

// GetName returns the screen name
func (s *SectorSelection) GetName() string {
	return "sector_selection"
}

// Render renders the sector selection UI
func (s *SectorSelection) Render() fyne.CanvasObject {
	// Header
	title := widget.NewLabel("Screen 1: Select Trading Sector")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Choose a sector based on 293 backtest results")
	subtitle.Alignment = fyne.TextAlignCenter

	progressLabel := widget.NewLabel("Step 1 of 8")
	progressLabel.Alignment = fyne.TextAlignCenter

	// Info banner
	infoBanner := s.createInfoBanner()

	// Create sector cards (sorted by priority)
	sectorCards := s.createSectorCards()

	// Navigation buttons
	s.continueBtn = widget.NewButton("Continue to Screener →", func() {
		if s.onNext != nil {
			if err := s.onNext(); err != nil {
				// Show error dialog if navigation fails
				s.showError(fmt.Sprintf("Navigation error: %v", err))
			}
		}
	})
	s.continueBtn.Importance = widget.HighImportance

	// Enable button only if sector is selected
	if !s.Validate() {
		s.continueBtn.Disable()
	}

	backBtn := widget.NewButton("← Dashboard", func() {
		if s.onBack != nil {
			s.onBack()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		if s.onCancel != nil {
			s.onCancel()
		}
	})

	navButtons := container.NewBorder(
		nil, nil,
		container.NewHBox(backBtn, cancelBtn),
		s.continueBtn,
	)

	// Main content
	content := container.NewBorder(
		// Top
		container.NewVBox(
			progressLabel,
			title,
			subtitle,
			widget.NewSeparator(),
			infoBanner,
			widget.NewSeparator(),
		),
		// Bottom
		container.NewVBox(
			widget.NewSeparator(),
			navButtons,
		),
		// Left, Right
		nil, nil,
		// Center
		container.NewVScroll(sectorCards),
	)

	return content
}

func (s *SectorSelection) createSectorCards() fyne.CanvasObject {
	if s.state.Policy == nil {
		return widget.NewLabel("Loading sectors...")
	}

	// Sort sectors by priority (lower number = higher priority)
	sectors := make([]models.Sector, len(s.state.Policy.Sectors))
	copy(sectors, s.state.Policy.Sectors)
	sort.Slice(sectors, func(i, j int) bool {
		return sectors[i].Priority < sectors[j].Priority
	})

	cards := container.NewVBox()

	for _, sector := range sectors {
		card := s.createSectorCard(sector)
		cards.Add(card)
		cards.Add(layout.NewSpacer()) // Add spacing between cards
	}

	return cards
}

func (s *SectorSelection) createSectorCard(sector models.Sector) fyne.CanvasObject {
	// Sector name with priority badge
	nameLabel := widget.NewLabel(fmt.Sprintf("Priority %d: %s", sector.Priority, sector.Name))
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Notes
	notesLabel := widget.NewLabel(sector.Notes)
	notesLabel.Wrapping = fyne.TextWrapWord

	// Strategy badges
	badgeHeader := widget.NewLabel("Top Strategy Fits:")
	badgeHeader.TextStyle = fyne.TextStyle{Bold: true}
	badgeList := buildCompactStrategyList(s.state.Policy, sector, 5)

	// Status indicator
	var statusLabel *widget.Label
	var statusColor color.Color

	if sector.Blocked {
		statusLabel = widget.NewLabel("❌ BLOCKED - Do Not Trade (0% backtest success)")
		statusColor = theme.ErrorColor()
		statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	} else if sector.Warning {
		statusLabel = widget.NewLabel("⚠️  WARNING - Use with Caution (marginal backtest results)")
		statusColor = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange
		statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		statusLabel = widget.NewLabel("✅ Approved for Trading (validated by backtests)")
		statusColor = color.RGBA{R: 0, G: 200, B: 100, A: 255} // Green
	}

	// Select button with improved styling
	selectBtn := widget.NewButton("Select This Sector", func() {
		s.selectSector(sector)
	})

	// Check if this sector is currently selected
	isSelected := s.state.CurrentTrade != nil && s.state.CurrentTrade.Sector == sector.Name

	if sector.Blocked {
		selectBtn.Disable()
	} else if isSelected {
		selectBtn.SetText("✓ Selected")
		selectBtn.Importance = widget.HighImportance
	}

	// Card background
	var cardBg *canvas.Rectangle
	if isSelected {
		// Highlight selected sector
		cardBg = canvas.NewRectangle(color.RGBA{R: 0, G: 100, B: 50, A: 50}) // Light green tint
	} else if sector.Blocked {
		// Grey out blocked sectors
		cardBg = canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 30}) // Grey tint
	} else if sector.Warning {
		// Amber tint for warnings
		cardBg = canvas.NewRectangle(color.RGBA{R: 255, G: 200, B: 100, A: 30}) // Amber tint
	} else {
		// Normal background
		cardBg = canvas.NewRectangle(color.Transparent)
	}

	// Status indicator bar (colored left border)
	statusBar := canvas.NewRectangle(statusColor)
	statusBar.SetMinSize(fyne.NewSize(4, 100))

	cardContent := container.NewVBox(
		nameLabel,
		notesLabel,
		badgeHeader,
		badgeList,
		statusLabel,
		selectBtn,
	)

	// Combine status bar with content
	cardWithBorder := container.NewBorder(
		nil, nil,
		statusBar,
		nil,
		cardContent,
	)

	// Stack background and content
	card := container.NewStack(
		cardBg,
		container.NewPadded(cardWithBorder),
	)

	return card
}

func (s *SectorSelection) selectSector(sector models.Sector) {
	// Phase 6: Check if sector has Utilities warning modal
	if sector.UtilitiesWarning != nil {
		s.showUtilitiesModal(sector)
		return
	}

	// Don't allow selection of blocked sectors (legacy check - shouldn't happen with Phase 6)
	if sector.Blocked {
		s.showError(fmt.Sprintf("%s is blocked for trading based on backtest results", sector.Name))
		return
	}

	// Initialize a new trade with selected sector
	if s.state.CurrentTrade == nil {
		s.state.CurrentTrade = &models.Trade{}
	}
	s.state.CurrentTrade.Sector = sector.Name
	s.selectedSector = sector.Name

	// Enable continue button now that sector is selected
	if s.continueBtn != nil {
		s.continueBtn.Enable()
	}

	// Refresh the entire screen to show selection state
	s.window.SetContent(s.Render())
}

// createInfoBanner creates an informational banner at the top
func (s *SectorSelection) createInfoBanner() fyne.CanvasObject {
	infoText := "Sectors are ranked by backtest performance. " +
		"Healthcare and Technology show the strongest results with 90%+ success rates. " +
		"Blocked sectors have failed backtests and should not be traded."

	infoLabel := widget.NewLabel(infoText)
	infoLabel.Wrapping = fyne.TextWrapWord

	infoIcon := widget.NewIcon(theme.InfoIcon())

	banner := container.NewBorder(
		nil, nil,
		infoIcon,
		nil,
		infoLabel,
	)

	// Light blue background for info banner
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 230, B: 255, A: 100})

	return container.NewStack(bg, container.NewPadded(banner))
}

// showUtilitiesModal displays a warning modal for Utilities sector
func (s *SectorSelection) showUtilitiesModal(sector models.Sector) {
	warning := sector.UtilitiesWarning

	// Title with icon
	titleLabel := widget.NewLabel(warning.Title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	// Warning message
	messageLabel := widget.NewLabel(warning.Message)
	messageLabel.Wrapping = fyne.TextWrapWord

	// Acknowledgement checkbox
	ackCheckbox := widget.NewCheck(warning.AcknowledgementText, nil)

	// Create a container for the content
	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		messageLabel,
		widget.NewSeparator(),
		ackCheckbox,
	)

	// Create custom dialog with buttons
	dialog := widget.NewModalPopUp(
		content,
		s.window.Canvas(),
	)

	// Go Back button
	goBackBtn := widget.NewButton("← Go Back", func() {
		dialog.Hide()
		fmt.Println("User declined Utilities sector (Go Back)")
	})

	// Continue Anyway button (initially disabled)
	continueBtn := widget.NewButton("Continue Anyway →", func() {
		if ackCheckbox.Checked {
			// Initialize trade with Utilities sector
			if s.state.CurrentTrade == nil {
				s.state.CurrentTrade = &models.Trade{}
			}
			s.state.CurrentTrade.Sector = sector.Name
			s.state.CurrentTrade.UtilitiesWarningAcknowledged = true
			s.selectedSector = sector.Name

			// Log warning acknowledgement
			fmt.Printf("⚠️  UTILITIES SECTOR WARNING ACKNOWLEDGED by user\n")

			// Enable continue button
			if s.continueBtn != nil {
				s.continueBtn.Enable()
			}

			// Hide modal and refresh screen
			dialog.Hide()
			s.window.SetContent(s.Render())
		}
	})
	continueBtn.Importance = widget.HighImportance
	continueBtn.Disable()

	// Enable Continue button when checkbox is checked
	ackCheckbox.OnChanged = func(checked bool) {
		if checked {
			continueBtn.Enable()
		} else {
			continueBtn.Disable()
		}
	}

	// Add buttons to dialog
	buttonRow := container.NewHBox(
		goBackBtn,
		layout.NewSpacer(),
		continueBtn,
	)

	content.Add(widget.NewSeparator())
	content.Add(buttonRow)

	// Set dialog size
	dialog.Resize(fyne.NewSize(600, 450))
	dialog.Show()
}

// showError displays an error message to the user
func (s *SectorSelection) showError(message string) {
	if s.window != nil {
		errorLabel := widget.NewLabel("Error: " + message)
		errorLabel.Wrapping = fyne.TextWrapWord

		// Simple error display - in production, this would be a proper dialog
		// For now, just log it since we don't want to disrupt the current screen
		fmt.Printf("ERROR: %s\n", message)
	}
}
