package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// SectorSelection represents Screen 1: Sector Selection
type SectorSelection struct {
	state  *appcore.AppState
	window fyne.Window
}

// NewSectorSelection creates a new sector selection screen
func NewSectorSelection(state *appcore.AppState, window fyne.Window) *SectorSelection {
	return &SectorSelection{
		state:  state,
		window: window,
	}
}

// Render renders the sector selection UI
func (s *SectorSelection) Render() fyne.CanvasObject {
	title := widget.NewLabel("Select Trading Sector")
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Choose a sector based on 293 backtest results:")

	// Create sector cards
	sectorCards := s.createSectorCards()

	continueBtn := widget.NewButton("Continue →", func() {
		// TODO: Navigate to screener launch
	})
	continueBtn.Disable() // Enable when sector is selected

	backBtn := widget.NewButton("← Back to Dashboard", func() {
		// TODO: Navigate back to dashboard
	})

	content := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		title,
		subtitle,
		widget.NewSeparator(),
		sectorCards,
		widget.NewSeparator(),
		continueBtn,
	)

	return container.NewPadded(content)
}

func (s *SectorSelection) createSectorCards() fyne.CanvasObject {
	if s.state.Policy == nil {
		return widget.NewLabel("Loading sectors...")
	}

	cards := container.NewVBox()

	for _, sector := range s.state.Policy.Sectors {
		card := s.createSectorCard(sector)
		cards.Add(card)
	}

	return cards
}

func (s *SectorSelection) createSectorCard(sector models.Sector) fyne.CanvasObject {
	name := widget.NewLabel(sector.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}

	notes := widget.NewLabel(sector.Notes)

	var statusLabel *widget.Label
	if sector.Blocked {
		statusLabel = widget.NewLabel("❌ BLOCKED - Do Not Trade")
	} else if sector.Warning {
		statusLabel = widget.NewLabel("⚠ WARNING - Use with Caution")
	} else {
		statusLabel = widget.NewLabel("✓ Approved for Trading")
	}

	selectBtn := widget.NewButton("Select", func() {
		s.selectSector(sector)
	})

	if sector.Blocked {
		selectBtn.Disable()
	}

	card := container.NewVBox(
		name,
		notes,
		statusLabel,
		selectBtn,
	)

	return container.NewPadded(card)
}

func (s *SectorSelection) selectSector(sector models.Sector) {
	// Initialize a new trade with selected sector
	if s.state.CurrentTrade == nil {
		s.state.CurrentTrade = &models.Trade{}
	}
	s.state.CurrentTrade.Sector = sector.Name

	// TODO: Navigate to next screen
}
