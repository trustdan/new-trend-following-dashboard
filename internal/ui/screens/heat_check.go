package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// HeatCheck represents Screen 6: Portfolio Heat Validation
type HeatCheck struct {
	state  *appcore.AppState
	window fyne.Window

	// Navigation callbacks
	onNext   func()
	onBack   func()
	onCancel func()

	// UI components
	portfolioHeatBar *canvas.Rectangle
	portfolioLabel   *widget.Label
	sectorHeatBars   map[string]*HeatBar
	warningLabel     *widget.Label
	continueBtn      *widget.Button
}

// HeatBar represents a visual heat bar for a sector
type HeatBar struct {
	sector       string
	currentHeat  float64
	newTradeHeat float64
	cap          float64
	container    *fyne.Container
}

// NewHeatCheck creates a new heat check screen
func NewHeatCheck(state *appcore.AppState, window fyne.Window) *HeatCheck {
	return &HeatCheck{
		state:          state,
		window:         window,
		sectorHeatBars: make(map[string]*HeatBar),
	}
}

// Validate checks if the screen's data is valid
func (s *HeatCheck) Validate() bool {
	// Heat check must pass (not exceed limits)
	if s.state.CurrentTrade == nil {
		return false
	}

	// Calculate if new trade would exceed limits
	portfolioHeat, sectorHeat := s.calculateHeat()

	// Check portfolio cap (4%)
	if portfolioHeat > s.state.Policy.Defaults.PortfolioHeatCap {
		return false
	}

	// Check sector cap (1.5% or sector-specific)
	sectorCap := s.getSectorCap(s.state.CurrentTrade.Sector)
	if sectorHeat > sectorCap {
		return false
	}

	return true
}

// GetName returns the screen name
func (s *HeatCheck) GetName() string {
	return "heat_check"
}

// SetNavCallbacks sets navigation callback functions
func (s *HeatCheck) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
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

// Render renders the heat check UI
func (s *HeatCheck) Render() fyne.CanvasObject {
	// Header
	header := s.createHeader()

	// Info banner
	infoBanner := s.createInfoBanner()

	// Portfolio heat section
	portfolioSection := s.createPortfolioHeatSection()

	// Sector heat sections
	sectorSections := s.createSectorHeatSections()

	// Warning/approval message
	s.warningLabel = widget.NewLabel("")
	s.warningLabel.Wrapping = fyne.TextWrapWord
	s.warningLabel.Alignment = fyne.TextAlignCenter
	s.warningLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Navigation buttons
	navButtons := s.createNavigationButtons()

	// Update validation state
	s.updateValidation()

	// Scrollable content
	scrollContent := container.NewVBox(
		infoBanner,
		portfolioSection,
		widget.NewSeparator(),
	)

	// Add sector sections
	for _, section := range sectorSections {
		scrollContent.Add(section)
		scrollContent.Add(widget.NewSeparator())
	}

	scrollContent.Add(s.warningLabel)

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
func (s *HeatCheck) createHeader() fyne.CanvasObject {
	progress := widget.NewLabel("Step 6 of 8")
	progress.TextStyle = fyne.TextStyle{Italic: true}

	title := widget.NewLabel("Screen 6: Portfolio Heat Check")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Enforce diversification and concentration limits")
	subtitle.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		progress,
		title,
		subtitle,
		widget.NewSeparator(),
	)
}

// createInfoBanner creates the heat check explanation banner
func (s *HeatCheck) createInfoBanner() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 255, G: 200, B: 100, A: 255})

	icon := widget.NewIcon(theme.WarningIcon())

	text := widget.NewLabel(
		"⚠️ Heat Check: This screen prevents overconcentration in any sector. " +
			"Portfolio cap: 4% total risk. Sector cap: 1.5% per sector. " +
			"If limits are exceeded, you must close existing positions first.",
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

// createPortfolioHeatSection creates the overall portfolio heat display
func (s *HeatCheck) createPortfolioHeatSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("📊 Portfolio-Wide Heat")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}

	// Calculate portfolio heat
	portfolioHeat, _ := s.calculateHeat()
	portfolioCap := s.state.Policy.Defaults.PortfolioHeatCap

	// Create heat bar
	heatPercent := portfolioHeat / portfolioCap
	barColor := s.getHeatColor(heatPercent)

	heatBar := canvas.NewRectangle(barColor)
	heatBar.SetMinSize(fyne.NewSize(float32(heatPercent*400), 30))

	capLine := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 128})
	capLine.SetMinSize(fyne.NewSize(2, 40))

	barContainer := container.NewHBox(
		heatBar,
	)

	// Labels
	currentLabel := widget.NewLabel(fmt.Sprintf("Current: %.2f%%", portfolioHeat*100))
	capLabel := widget.NewLabel(fmt.Sprintf("Cap: %.2f%%", portfolioCap*100))

	return container.NewVBox(
		sectionTitle,
		barContainer,
		container.NewHBox(currentLabel, widget.NewLabel("/"), capLabel),
	)
}

// createSectorHeatSections creates heat bars for each sector
func (s *HeatCheck) createSectorHeatSections() []fyne.CanvasObject {
	sections := []fyne.CanvasObject{}

	sectionTitle := widget.NewLabel("🏥 Sector Heat Breakdown")
	sectionTitle.TextStyle = fyne.TextStyle{Bold: true}
	sections = append(sections, sectionTitle)

	// Get all sectors from policy
	for _, sector := range s.state.Policy.Sectors {
		if sector.Blocked {
			continue // Skip blocked sectors
		}

		// Calculate heat for this sector
		currentHeat := s.calculateSectorHeat(sector.Name, false)
		totalHeat := s.calculateSectorHeat(sector.Name, true)
		newTradeHeat := totalHeat - currentHeat
		sectorCap := s.getSectorCap(sector.Name)

		// Only show sectors with active trades or the current sector
		if currentHeat > 0 || sector.Name == s.state.CurrentTrade.Sector {
			heatBar := s.createSectorHeatBar(sector.Name, currentHeat, newTradeHeat, sectorCap)
			sections = append(sections, heatBar)
		}
	}

	return sections
}

// createSectorHeatBar creates a visual heat bar for a sector
func (s *HeatCheck) createSectorHeatBar(sectorName string, currentHeat, newTradeHeat, cap float64) fyne.CanvasObject {
	// Sector label
	label := widget.NewLabel(fmt.Sprintf("%s:", sectorName))
	label.TextStyle = fyne.TextStyle{Bold: true}

	// Current heat bar
	currentPercent := currentHeat / cap
	currentColor := s.getHeatColor(currentPercent)
	currentBar := canvas.NewRectangle(currentColor)
	currentBar.SetMinSize(fyne.NewSize(float32(currentPercent*300), 20))

	// New trade heat bar (if applicable)
	totalPercent := (currentHeat + newTradeHeat) / cap
	newTradeBar := canvas.NewRectangle(color.RGBA{R: 255, G: 200, B: 0, A: 180})
	newTradeBar.SetMinSize(fyne.NewSize(float32((totalPercent-currentPercent)*300), 20))

	// Heat values
	heatLabel := widget.NewLabel(fmt.Sprintf(
		"%.2f%% → %.2f%% / %.2f%%",
		currentHeat*100,
		(currentHeat+newTradeHeat)*100,
		cap*100,
	))

	// Status icon
	var statusIcon *widget.Label
	if totalPercent > 1.0 {
		statusIcon = widget.NewLabel("❌ Exceeds limit")
		statusIcon.TextStyle = fyne.TextStyle{Bold: true}
	} else if totalPercent > 0.8 {
		statusIcon = widget.NewLabel("⚠️ Approaching limit")
	} else {
		statusIcon = widget.NewLabel("✓ Safe")
	}

	barContainer := container.NewHBox(
		currentBar,
		newTradeBar,
	)

	return container.NewVBox(
		container.NewHBox(label, statusIcon),
		barContainer,
		heatLabel,
	)
}

// createNavigationButtons creates the navigation button bar
func (s *HeatCheck) createNavigationButtons() fyne.CanvasObject {
	backBtn := widget.NewButton("← Back", func() {
		s.saveHeatCheck()
		if s.onBack != nil {
			s.onBack()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		s.saveHeatCheck()
		if s.onCancel != nil {
			s.onCancel()
		}
	})

	s.continueBtn = widget.NewButton("Continue →", func() {
		if s.Validate() {
			s.saveHeatCheck()
			if s.onNext != nil {
				s.onNext()
			}
		}
	})
	s.continueBtn.Importance = widget.HighImportance

	return container.NewBorder(
		nil,
		nil,
		container.NewHBox(backBtn, cancelBtn),
		s.continueBtn,
		nil,
	)
}

// calculateHeat calculates portfolio-wide and sector-specific heat
func (s *HeatCheck) calculateHeat() (portfolioHeat, sectorHeat float64) {
	if s.state.CurrentTrade == nil || s.state.CurrentTrade.AccountEquity == 0 {
		return 0, 0
	}

	// TODO: Load active trades from storage
	// For now, use empty slice (will be implemented in storage integration)
	activeTrades := []models.Trade{}

	accountSize := s.state.CurrentTrade.AccountEquity

	// Calculate existing heat
	for _, trade := range activeTrades {
		tradeHeat := trade.MaxLoss / accountSize
		portfolioHeat += tradeHeat

		if trade.Sector == s.state.CurrentTrade.Sector {
			sectorHeat += tradeHeat
		}
	}

	// Add new trade heat
	newTradeHeat := s.state.CurrentTrade.MaxLoss / accountSize
	portfolioHeat += newTradeHeat
	sectorHeat += newTradeHeat

	return portfolioHeat, sectorHeat
}

// calculateSectorHeat calculates heat for a specific sector
func (s *HeatCheck) calculateSectorHeat(sectorName string, includeNewTrade bool) float64 {
	if s.state.CurrentTrade == nil || s.state.CurrentTrade.AccountEquity == 0 {
		return 0
	}

	// TODO: Load active trades from storage
	activeTrades := []models.Trade{}

	accountSize := s.state.CurrentTrade.AccountEquity
	var heat float64

	// Calculate existing heat
	for _, trade := range activeTrades {
		if trade.Sector == sectorName {
			heat += trade.MaxLoss / accountSize
		}
	}

	// Add new trade heat if applicable
	if includeNewTrade && s.state.CurrentTrade.Sector == sectorName {
		heat += s.state.CurrentTrade.MaxLoss / accountSize
	}

	return heat
}

// getSectorCap returns the heat cap for a specific sector
func (s *HeatCheck) getSectorCap(sectorName string) float64 {
	// Find sector in policy
	for _, sector := range s.state.Policy.Sectors {
		if sector.Name == sectorName {
			// Use sector-specific cap if available, otherwise use default bucket cap
			if sector.HeatCapPercent > 0 {
				return sector.HeatCapPercent
			}
			return s.state.Policy.Defaults.BucketHeatCap
		}
	}

	// Default fallback
	return s.state.Policy.Defaults.BucketHeatCap
}

// getHeatColor returns the appropriate color for a heat percentage
func (s *HeatCheck) getHeatColor(heatPercent float64) color.Color {
	if heatPercent > 1.0 {
		return color.RGBA{R: 220, G: 0, B: 0, A: 255} // Red - exceeds limit
	} else if heatPercent > 0.8 {
		return color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange - approaching limit
	} else if heatPercent > 0.5 {
		return color.RGBA{R: 255, G: 215, B: 0, A: 255} // Yellow - moderate
	} else {
		return color.RGBA{R: 0, G: 180, B: 80, A: 255} // Green - safe
	}
}

// updateValidation updates the continue button state and warning message
func (s *HeatCheck) updateValidation() {
	isValid := s.Validate()

	if s.continueBtn != nil {
		if isValid {
			s.continueBtn.Enable()
		} else {
			s.continueBtn.Disable()
		}
	}

	if s.warningLabel != nil {
		if !isValid {
			portfolioHeat, sectorHeat := s.calculateHeat()
			sectorCap := s.getSectorCap(s.state.CurrentTrade.Sector)

			if portfolioHeat > s.state.Policy.Defaults.PortfolioHeatCap {
				s.warningLabel.SetText(fmt.Sprintf(
					"❌ Portfolio Heat Exceeded: %.2f%% > %.2f%% cap\n"+
						"You must close existing positions to proceed.",
					portfolioHeat*100,
					s.state.Policy.Defaults.PortfolioHeatCap*100,
				))
			} else if sectorHeat > sectorCap {
				s.warningLabel.SetText(fmt.Sprintf(
					"❌ %s Sector Heat Exceeded: %.2f%% > %.2f%% cap\n"+
						"You must close existing %s positions to proceed.",
					s.state.CurrentTrade.Sector,
					sectorHeat*100,
					sectorCap*100,
					s.state.CurrentTrade.Sector,
				))
			}
		} else {
			s.warningLabel.SetText("✓ Heat checks passed - Safe to proceed")
		}
	}
}

// saveHeatCheck saves the heat check results to the trade
func (s *HeatCheck) saveHeatCheck() {
	if s.state.CurrentTrade == nil {
		return
	}

	// Calculate and save heat values
	portfolioHeat, sectorHeat := s.calculateHeat()
	s.state.CurrentTrade.PortfolioHeat = portfolioHeat
	s.state.CurrentTrade.BucketHeat = sectorHeat
	s.state.CurrentTrade.HeatCheckPassed = s.Validate()
}
