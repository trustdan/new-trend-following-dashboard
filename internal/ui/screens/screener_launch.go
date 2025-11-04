package screens

import (
	"fmt"
	"image/color"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
)

// ScreenerLaunch represents Screen 2: Launch FINVIZ Screeners
type ScreenerLaunch struct {
	state       *appcore.AppState
	window      fyne.Window
	onNext      func() error
	onBack      func() error
	onCancel    func()
	lastLaunch  map[string]time.Time // Track last launch time per screener
	continueBtn *widget.Button
}

// NewScreenerLaunch creates a new screener launch screen
func NewScreenerLaunch(state *appcore.AppState, window fyne.Window) *ScreenerLaunch {
	return &ScreenerLaunch{
		state:      state,
		window:     window,
		lastLaunch: make(map[string]time.Time),
	}
}

// SetNavCallbacks sets navigation callback functions
func (s *ScreenerLaunch) SetNavCallbacks(onNext, onBack func() error, onCancel func()) {
	s.onNext = onNext
	s.onBack = onBack
	s.onCancel = onCancel
}

// Render renders the screener launch UI
func (s *ScreenerLaunch) Render() fyne.CanvasObject {
	// Header
	title := widget.NewLabel("Screen 2: Launch FINVIZ Screener")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Get selected sector
	var selectedSector *models.Sector
	if s.state.CurrentTrade != nil && s.state.CurrentTrade.Sector != "" {
		for i, sector := range s.state.Policy.Sectors {
			if sector.Name == s.state.CurrentTrade.Sector {
				selectedSector = &s.state.Policy.Sectors[i]
				break
			}
		}
	}

	var sectorName string
	if s.state.CurrentTrade != nil {
		sectorName = s.state.CurrentTrade.Sector
	} else {
		sectorName = "Unknown"
	}

	subtitle := widget.NewLabel(fmt.Sprintf("Launch screeners for %s sector", sectorName))
	subtitle.Alignment = fyne.TextAlignCenter

	progressLabel := widget.NewLabel("Step 2 of 8")
	progressLabel.Alignment = fyne.TextAlignCenter

	// Info banner and strategy badges
	infoBanner := s.createInfoBanner()

	var badgeHeader fyne.CanvasObject
	var badgeList fyne.CanvasObject
	if selectedSector != nil {
		header := widget.NewLabel("Top Strategy Fits:")
		header.TextStyle = fyne.TextStyle{Bold: true}
		badgeHeader = header
		badgeList = buildCompactStrategyList(s.state.Policy, *selectedSector, 5)
	}

	// Create screener cards
	var screenerCards fyne.CanvasObject
	if selectedSector != nil {
		screenerCards = s.createScreenerCards(*selectedSector)
	} else {
		screenerCards = widget.NewLabel("Error: No sector selected")
	}

	// Navigation buttons
	s.continueBtn = widget.NewButton("Continue to Ticker Entry →", func() {
		if s.onNext != nil {
			if err := s.onNext(); err != nil {
				s.showError(fmt.Sprintf("Navigation error: %v", err))
			}
		}
	})
	s.continueBtn.Importance = widget.HighImportance

	backBtn := widget.NewButton("← Back to Sector", func() {
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
	topItems := []fyne.CanvasObject{
		progressLabel,
		title,
		subtitle,
		widget.NewSeparator(),
		infoBanner,
	}
	if badgeHeader != nil && badgeList != nil {
		topItems = append(topItems,
			widget.NewSeparator(),
			badgeHeader,
			badgeList,
		)
	}
	topItems = append(topItems, widget.NewSeparator())

	content := container.NewBorder(
		container.NewVBox(topItems...),
		// Bottom
		container.NewVBox(
			widget.NewSeparator(),
			navButtons,
		),
		// Left, Right
		nil, nil,
		// Center
		container.NewVScroll(screenerCards),
	)

	return content
}

// Validate checks if the screen's data is valid
func (s *ScreenerLaunch) Validate() bool {
	// Screener launch always valid (just opens URL)
	return true
}

// GetName returns the screen name
func (s *ScreenerLaunch) GetName() string {
	return "screener_launch"
}

// createScreenerCards creates cards for all available screeners
func (s *ScreenerLaunch) createScreenerCards(sector models.Sector) fyne.CanvasObject {
	cards := container.NewVBox()

	// Define screener metadata
	screenerInfo := map[string]struct {
		title       string
		description string
		frequency   string
		purpose     string
		direction   string // "bullish" or "bearish"
	}{
		// Bullish screeners
		"universe": {
			title:       "Universe Screener (Bullish)",
			description: "Find 30-60 quality stocks in long-term uptrends",
			frequency:   "Run: Weekly (Monday mornings)",
			purpose:     "Build your watch list of trendable stocks",
			direction:   "bullish",
		},
		"pullback": {
			title:       "Pullback Screener (Bullish)",
			description: "Oversold stocks in uptrends (RSI < 40, price above SMA200)",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Find stocks retracing into support levels",
			direction:   "bullish",
		},
		"breakout": {
			title:       "Breakout Screener (Bullish)",
			description: "New 52-week highs with strong momentum",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Catch momentum breakouts early",
			direction:   "bullish",
		},
		"golden_cross": {
			title:       "Golden Cross Screener (Bullish)",
			description: "SMA50 crossing above SMA200 (bullish trend confirmation)",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Identify major trend reversals",
			direction:   "bullish",
		},
		// Bearish screeners
		"universe_bearish": {
			title:       "Universe Screener (Bearish)",
			description: "Find 30-60 quality stocks in long-term downtrends",
			frequency:   "Run: Weekly (Monday mornings)",
			purpose:     "Build your watch list of shortable stocks",
			direction:   "bearish",
		},
		"bounce_bearish": {
			title:       "Bounce Screener (Bearish)",
			description: "Overbought stocks in downtrends (RSI > 60, price below SMA200)",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Find stocks bouncing into resistance levels",
			direction:   "bearish",
		},
		"breakdown_bearish": {
			title:       "Breakdown Screener (Bearish)",
			description: "New 52-week lows with strong momentum",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Catch momentum breakdowns early",
			direction:   "bearish",
		},
		"death_cross_bearish": {
			title:       "Death Cross Screener (Bearish)",
			description: "SMA50 crossing below SMA200 (bearish trend confirmation)",
			frequency:   "Run: Daily (before market open)",
			purpose:     "Identify major trend reversals to downside",
			direction:   "bearish",
		},
	}

	// Bullish section header
	bullishHeader := widget.NewLabel("🟢 Bullish Screeners (Long Opportunities)")
	bullishHeader.TextStyle = fyne.TextStyle{Bold: true}
	cards.Add(bullishHeader)
	cards.Add(widget.NewSeparator())

	// Create bullish screener cards
	bullishOrder := []string{"universe", "pullback", "breakout", "golden_cross"}
	bullishCount := 0
	for _, key := range bullishOrder {
		if screenerURL, exists := sector.ScreenerURLs[key]; exists {
			info := screenerInfo[key]
			card := s.createScreenerCard(key, info.title, info.description, info.frequency, info.purpose, screenerURL, "bullish")
			cards.Add(card)
			bullishCount++
		}
	}

	// Add spacer between sections
	cards.Add(widget.NewLabel(""))

	// Bearish section header
	bearishHeader := widget.NewLabel("🔴 Bearish Screeners (Short/Put Opportunities)")
	bearishHeader.TextStyle = fyne.TextStyle{Bold: true}
	cards.Add(bearishHeader)
	cards.Add(widget.NewSeparator())

	// Create bearish screener cards
	bearishOrder := []string{"universe_bearish", "bounce_bearish", "breakdown_bearish", "death_cross_bearish"}
	bearishCount := 0
	for _, key := range bearishOrder {
		if screenerURL, exists := sector.ScreenerURLs[key]; exists {
			info := screenerInfo[key]
			card := s.createScreenerCard(key, info.title, info.description, info.frequency, info.purpose, screenerURL, "bearish")
			cards.Add(card)
			bearishCount++
		}
	}

	if bullishCount == 0 && bearishCount == 0 {
		return widget.NewLabel("No screeners available for this sector")
	}

	return cards
}

// createScreenerCard creates a single screener card
func (s *ScreenerLaunch) createScreenerCard(key, title, description, frequency, purpose, screenerURL, direction string) fyne.CanvasObject {
	// Title
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Description
	descLabel := widget.NewLabel(description)
	descLabel.Wrapping = fyne.TextWrapWord

	// Frequency
	freqLabel := widget.NewLabel(frequency)
	freqLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Purpose
	purposeLabel := widget.NewLabel("Purpose: " + purpose)

	// Sort info - need to strip _bearish suffix for lookup
	sortKey := key
	if len(key) > 8 && key[len(key)-8:] == "_bearish" {
		sortKey = key[:len(key)-8] // Remove _bearish suffix for sort lookup
	}

	sortLabel := widget.NewLabel("")
	if s.state.Policy != nil && s.state.Policy.ScreenerSorting != nil {
		if sortConfig, exists := s.state.Policy.ScreenerSorting[sortKey]; exists {
			sortDirection := "ascending"
			if sortConfig.Direction == "desc" {
				sortDirection = "descending"
			}
			sortLabel.SetText(fmt.Sprintf("📊 Sorted by: %s (%s) - %s",
				sortConfig.SortBy, sortDirection, sortConfig.Rationale))
			sortLabel.TextStyle = fyne.TextStyle{Italic: true}
		}
	}

	// Last run timestamp
	lastRunLabel := widget.NewLabel("")
	if lastRun, exists := s.lastLaunch[key]; exists {
		elapsed := time.Since(lastRun)
		if elapsed < time.Minute {
			lastRunLabel.SetText(fmt.Sprintf("✓ Launched %d seconds ago", int(elapsed.Seconds())))
		} else if elapsed < time.Hour {
			lastRunLabel.SetText(fmt.Sprintf("✓ Launched %d minutes ago", int(elapsed.Minutes())))
		} else {
			lastRunLabel.SetText(fmt.Sprintf("✓ Launched %s", lastRun.Format("3:04 PM")))
		}
	}

	// Launch button
	launchBtn := widget.NewButton("🔗 Open in Browser", func() {
		s.launchURL(key, screenerURL, direction)
	})
	launchBtn.Importance = widget.HighImportance

	// Card content
	cardContent := container.NewVBox(
		titleLabel,
		descLabel,
		freqLabel,
		purposeLabel,
		sortLabel,
		lastRunLabel,
		launchBtn,
	)

	// Color-coded left border based on direction
	var leftBorder *canvas.Rectangle
	var bg *canvas.Rectangle
	if direction == "bearish" {
		leftBorder = canvas.NewRectangle(color.RGBA{R: 220, G: 53, B: 69, A: 255}) // Red
		bg = canvas.NewRectangle(color.RGBA{R: 255, G: 240, B: 240, A: 255})       // Very light red
	} else {
		leftBorder = canvas.NewRectangle(color.RGBA{R: 0, G: 120, B: 215, A: 255}) // Blue
		bg = canvas.NewRectangle(color.RGBA{R: 240, G: 248, B: 255, A: 255})       // Very light blue
	}
	leftBorder.SetMinSize(fyne.NewSize(4, 100))

	cardWithBorder := container.NewBorder(
		nil, nil,
		leftBorder,
		nil,
		cardContent,
	)

	card := container.NewStack(
		bg,
		container.NewPadded(cardWithBorder),
	)

	return card
}

// launchURL opens a FINVIZ screener URL in the default browser
func (s *ScreenerLaunch) launchURL(key, rawURL, direction string) {
	// Parse and validate URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		s.showError(fmt.Sprintf("Invalid URL for %s screener: %v", key, err))
		return
	}

	// Verify v=211 parameter is present (chart view)
	query := parsedURL.Query()
	if query.Get("v") != "211" {
		s.showError(fmt.Sprintf("Warning: URL missing v=211 chart view parameter"))
	}

	// Add screener-specific sorting if configured
	if s.state.Policy != nil && s.state.Policy.ScreenerSorting != nil {
		if sortConfig, exists := s.state.Policy.ScreenerSorting[key]; exists {
			sortParam := sortConfig.SortBy
			// Prepend minus for descending sort
			if sortConfig.Direction == "desc" {
				sortParam = "-" + sortParam
			}
			query.Set("o", sortParam)
			parsedURL.RawQuery = query.Encode()
		}
	}

	// Open URL in default browser
	urlObj, err := url.Parse(parsedURL.String())
	if err != nil {
		s.showError(fmt.Sprintf("Failed to parse URL: %v", err))
		return
	}

	// Use Fyne's OpenURL to launch browser
	if s.window != nil {
		app := fyne.CurrentApp()
		if app != nil {
			err = app.OpenURL(urlObj)
			if err != nil {
				s.showError(fmt.Sprintf("Failed to open URL: %v", err))
				return
			}
		}

		// Store direction in current trade for direction-aware checklist
		if s.state.CurrentTrade != nil {
			s.state.CurrentTrade.Direction = direction
		}

		// Record launch time
		s.lastLaunch[key] = time.Now()

		// Refresh screen to show timestamp
		s.window.SetContent(s.Render())
	}
}

// createInfoBanner creates an informational banner
func (s *ScreenerLaunch) createInfoBanner() fyne.CanvasObject {
	infoText := "FINVIZ screeners are pre-configured with filters for your selected sector. " +
		"Click any screener to open it in your browser. " +
		"Review the results, then proceed to enter your ticker symbol in the next screen."

	infoLabel := widget.NewLabel(infoText)
	infoLabel.Wrapping = fyne.TextWrapWord

	infoIcon := widget.NewIcon(theme.InfoIcon())

	banner := container.NewBorder(
		nil, nil,
		infoIcon,
		nil,
		infoLabel,
	)

	// Light blue background
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 230, B: 255, A: 100})

	return container.NewStack(bg, container.NewPadded(banner))
}

// showError displays an error message
func (s *ScreenerLaunch) showError(message string) {
	fmt.Printf("ERROR: %s\n", message)
}
