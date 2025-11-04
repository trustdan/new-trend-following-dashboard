package screens

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"tf-engine/internal/appcore"
	"tf-engine/internal/config"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
	"tf-engine/internal/testing/generators"
)

// Calendar represents Screen 8: Trade Calendar View (Horserace Timeline)
type Calendar struct {
	state        *appcore.AppState
	window       fyne.Window
	featureFlags *config.FeatureFlags
}

// NewCalendar creates a new calendar screen
func NewCalendar(state *appcore.AppState, window fyne.Window) *Calendar {
	return &Calendar{
		state:  state,
		window: window,
	}
}

// NewCalendarWithFlags creates a new calendar screen with feature flags
func NewCalendarWithFlags(state *appcore.AppState, window fyne.Window, featureFlags *config.FeatureFlags) *Calendar {
	return &Calendar{
		state:        state,
		window:       window,
		featureFlags: featureFlags,
	}
}

// Render renders the calendar UI
func (c *Calendar) Render() fyne.CanvasObject {
	title := widget.NewLabel("Screen 8: Trade Calendar (Horserace Timeline)")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Configuration from policy
	pastDays := 14
	futureDays := 84
	if c.state.Policy != nil {
		if c.state.Policy.Calendar.PastDays > 0 {
			pastDays = c.state.Policy.Calendar.PastDays
		}
		if c.state.Policy.Calendar.FutureDays > 0 {
			futureDays = c.state.Policy.Calendar.FutureDays
		}
	}

	// Calculate time range
	now := time.Now()
	startDate := now.AddDate(0, 0, -pastDays)
	endDate := now.AddDate(0, 0, futureDays)

	// Get all sectors from policy
	sectors := c.getSectors()

	// Summary stats
	activeTrades := c.countActiveTrades()
	totalRisk := c.calculateTotalRisk()
	portfolioHeat := 0.0
	if c.state.Policy != nil {
		// TODO: Get actual account size from settings
		accountSize := 50000.0
		portfolioHeat = totalRisk / accountSize
	}

	summaryText := fmt.Sprintf("Active Trades: %d | Total Risk: $%.2f | Portfolio Heat: %.2f%%",
		activeTrades, totalRisk, portfolioHeat*100)
	summary := widget.NewLabel(summaryText)

	// Create the timeline visualization
	timeline := c.createTimeline(sectors, startDate, endDate)

	// Legend
	legend := c.createLegend()

	// Action buttons
	newTradeBtn := widget.NewButton("+ New Trade", func() {
		// Navigator will handle going back to Screen 1
	})
	newTradeBtn.Importance = widget.HighImportance

	refreshBtn := widget.NewButton("Refresh", func() {
		// Reload all trades and refresh display
		c.loadAllTrades()
		c.window.SetContent(c.Render())
	})

	// Sample data button (Phase 2 feature)
	buttonsContainer := container.NewHBox(newTradeBtn, refreshBtn)

	if c.featureFlags != nil && c.featureFlags.IsEnabled("sample_data_generator") {
		sampleDataBtn := widget.NewButton("Generate Sample Data", func() {
			c.generateSampleData()
		})
		sampleDataBtn.Importance = widget.LowImportance
		buttonsContainer.Add(sampleDataBtn)
	}

	buttons := buttonsContainer

	content := container.NewVBox(
		title,
		summary,
		widget.NewSeparator(),
		legend,
		widget.NewSeparator(),
		timeline,
		widget.NewSeparator(),
		buttons,
	)

	scrollContainer := container.NewScroll(content)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	return container.NewPadded(scrollContainer)
}

// createTimeline creates the horserace timeline visualization
func (c *Calendar) createTimeline(sectors []string, startDate, endDate time.Time) fyne.CanvasObject {
	// Timeline dimensions
	const (
		sectorRowHeight = 80
		timelineWidth   = 1200
		leftMargin      = 120
		topMargin       = 40
		barHeight       = 30
	)

	// Calculate total canvas height
	canvasHeight := float32(len(sectors)*sectorRowHeight + topMargin + 20)
	canvasWidth := float32(timelineWidth + leftMargin + 20)

	// Create container for all elements
	var elements []fyne.CanvasObject

	// Background
	bg := canvas.NewRectangle(color.RGBA{R: 250, G: 250, B: 250, A: 255})
	bg.Resize(fyne.NewSize(canvasWidth, canvasHeight))
	elements = append(elements, bg)

	// Draw time axis (X-axis)
	elements = append(elements, c.drawTimeAxis(startDate, endDate, leftMargin, timelineWidth)...)

	// Draw sector rows (Y-axis)
	for i, sector := range sectors {
		y := float32(topMargin + i*sectorRowHeight)

		// Sector label
		sectorLabel := canvas.NewText(sector, color.Black)
		sectorLabel.TextSize = 12
		sectorLabel.TextStyle = fyne.TextStyle{Bold: true}
		sectorLabel.Move(fyne.NewPos(10, y+10))
		elements = append(elements, sectorLabel)

		// Sector row background (alternating colors)
		var rowColor color.Color
		if i%2 == 0 {
			rowColor = color.RGBA{R: 245, G: 245, B: 245, A: 255}
		} else {
			rowColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		rowBg := canvas.NewRectangle(rowColor)
		rowBg.Move(fyne.NewPos(leftMargin, y))
		rowBg.Resize(fyne.NewSize(float32(timelineWidth), float32(sectorRowHeight)))
		elements = append(elements, rowBg)

		// Draw trades for this sector
		trades := c.getTradesForSector(sector)
		for _, trade := range trades {
			tradeBar := c.createTradeBar(trade, startDate, endDate, leftMargin, timelineWidth, y, barHeight)
			elements = append(elements, tradeBar...)
		}
	}

	// "Today" line
	todayLine := c.drawTodayLine(startDate, endDate, leftMargin, timelineWidth, canvasHeight, topMargin)
	elements = append(elements, todayLine...)

	// Create container with all elements
	timeline := container.NewWithoutLayout(elements...)
	timeline.Resize(fyne.NewSize(canvasWidth, canvasHeight))

	return timeline
}

// drawTimeAxis draws the time axis labels
func (c *Calendar) drawTimeAxis(startDate, endDate time.Time, leftMargin float32, timelineWidth int) []fyne.CanvasObject {
	var elements []fyne.CanvasObject

	// Draw labels every 7 days
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	pixelsPerDay := float32(timelineWidth) / float32(totalDays)

	for i := 0; i <= totalDays; i += 7 {
		date := startDate.AddDate(0, 0, i)
		x := leftMargin + float32(i)*pixelsPerDay

		// Vertical grid line
		gridLine := canvas.NewLine(color.RGBA{R: 200, G: 200, B: 200, A: 255})
		gridLine.StrokeWidth = 1
		gridLine.Position1 = fyne.NewPos(x, 20)
		gridLine.Position2 = fyne.NewPos(x, 600) // Extend to bottom
		elements = append(elements, gridLine)

		// Date label
		dateLabel := canvas.NewText(date.Format("Jan 2"), color.RGBA{R: 100, G: 100, B: 100, A: 255})
		dateLabel.TextSize = 10
		dateLabel.Move(fyne.NewPos(x-20, 5))
		elements = append(elements, dateLabel)
	}

	return elements
}

// drawTodayLine draws a red vertical line at today's date
func (c *Calendar) drawTodayLine(startDate, endDate time.Time, leftMargin float32, timelineWidth int, canvasHeight, topMargin float32) []fyne.CanvasObject {
	var elements []fyne.CanvasObject

	now := time.Now()
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	daysSinceStart := int(now.Sub(startDate).Hours() / 24)

	if daysSinceStart >= 0 && daysSinceStart <= totalDays {
		pixelsPerDay := float32(timelineWidth) / float32(totalDays)
		x := leftMargin + float32(daysSinceStart)*pixelsPerDay

		// Red vertical line
		todayLine := canvas.NewLine(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		todayLine.StrokeWidth = 2
		todayLine.Position1 = fyne.NewPos(x, topMargin)
		todayLine.Position2 = fyne.NewPos(x, canvasHeight-20)
		elements = append(elements, todayLine)

		// "TODAY" label
		todayLabel := canvas.NewText("TODAY", color.RGBA{R: 255, G: 0, B: 0, A: 255})
		todayLabel.TextSize = 10
		todayLabel.TextStyle = fyne.TextStyle{Bold: true}
		todayLabel.Move(fyne.NewPos(x-20, topMargin-15))
		elements = append(elements, todayLabel)
	}

	return elements
}

// createTradeBar creates a visual bar for a single trade
func (c *Calendar) createTradeBar(trade models.Trade, startDate, endDate time.Time, leftMargin float32, timelineWidth int, y, barHeight float32) []fyne.CanvasObject {
	var elements []fyne.CanvasObject

	// Calculate position and width
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	pixelsPerDay := float32(timelineWidth) / float32(totalDays)

	entryDate := trade.CreatedAt
	expirationDate := trade.ExpirationDate

	daysSinceStart := int(entryDate.Sub(startDate).Hours() / 24)
	daysToExpiration := int(expirationDate.Sub(entryDate).Hours() / 24)

	// Only draw if trade is visible in timeline
	if daysSinceStart < 0 || daysSinceStart > totalDays {
		return elements
	}

	x := leftMargin + float32(daysSinceStart)*pixelsPerDay
	width := float32(daysToExpiration) * pixelsPerDay

	// Ensure minimum width for visibility
	if width < 40 {
		width = 40
	}

	// Determine bar color based on trade status and profitability
	barColor := c.getTradeBarColor(trade)

	// Draw trade bar
	tradeBar := canvas.NewRectangle(barColor)
	tradeBar.Move(fyne.NewPos(x, y+15))
	tradeBar.Resize(fyne.NewSize(width, barHeight))
	elements = append(elements, tradeBar)

	// Trade label (ticker symbol)
	tradeLabel := canvas.NewText(trade.Ticker, color.White)
	tradeLabel.TextSize = 11
	tradeLabel.TextStyle = fyne.TextStyle{Bold: true}
	tradeLabel.Move(fyne.NewPos(x+5, y+20))
	elements = append(elements, tradeLabel)

	// Options strategy label (smaller)
	strategyLabel := canvas.NewText(trade.OptionsStrategy, color.RGBA{R: 255, G: 255, B: 255, A: 200})
	strategyLabel.TextSize = 9
	strategyLabel.Move(fyne.NewPos(x+5, y+32))
	elements = append(elements, strategyLabel)

	return elements
}

// getTradeBarColor determines the color of a trade bar
func (c *Calendar) getTradeBarColor(trade models.Trade) color.Color {
	now := time.Now()
	daysToExpiration := int(trade.ExpirationDate.Sub(now).Hours() / 24)

	// Yellow: Expiring within 7 days
	if daysToExpiration >= 0 && daysToExpiration <= 7 {
		return color.RGBA{R: 255, G: 193, B: 7, A: 255} // Amber/Gold
	}

	// Red: Expired (past expiration date)
	if daysToExpiration < 0 {
		return color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red
	}

	// Green: Profitable (if P&L data available)
	if trade.ProfitLoss != nil && *trade.ProfitLoss > 0 {
		return color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green
	}

	// Red: Losing (if P&L data available)
	if trade.ProfitLoss != nil && *trade.ProfitLoss < 0 {
		return color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red
	}

	// Blue: Active trade (default)
	return color.RGBA{R: 13, G: 110, B: 253, A: 255} // Blue
}

// createLegend creates the color legend
func (c *Calendar) createLegend() fyne.CanvasObject {
	blueDot := canvas.NewCircle(color.RGBA{R: 13, G: 110, B: 253, A: 255})
	blueDot.Resize(fyne.NewSize(15, 15))

	greenDot := canvas.NewCircle(color.RGBA{R: 40, G: 167, B: 69, A: 255})
	greenDot.Resize(fyne.NewSize(15, 15))

	redDot := canvas.NewCircle(color.RGBA{R: 220, G: 53, B: 69, A: 255})
	redDot.Resize(fyne.NewSize(15, 15))

	yellowDot := canvas.NewCircle(color.RGBA{R: 255, G: 193, B: 7, A: 255})
	yellowDot.Resize(fyne.NewSize(15, 15))

	legend := container.NewHBox(
		widget.NewLabel("Legend:"),
		blueDot,
		widget.NewLabel("Active"),
		greenDot,
		widget.NewLabel("Profitable"),
		redDot,
		widget.NewLabel("Losing/Expired"),
		yellowDot,
		widget.NewLabel("Expiring Soon (<7 days)"),
	)

	return legend
}

// getSectors returns list of sectors from policy
func (c *Calendar) getSectors() []string {
	if c.state.Policy == nil || len(c.state.Policy.Sectors) == 0 {
		// Default sectors if policy not loaded
		return []string{
			"Healthcare",
			"Technology",
			"Industrials",
			"Consumer",
			"Financials",
			"Energy",
			"Real Estate",
			"Utilities",
		}
	}

	sectors := make([]string, len(c.state.Policy.Sectors))
	for i, sector := range c.state.Policy.Sectors {
		sectors[i] = sector.Name
	}
	return sectors
}

// getTradesForSector returns all trades for a specific sector
func (c *Calendar) getTradesForSector(sector string) []models.Trade {
	var trades []models.Trade
	for _, trade := range c.state.AllTrades {
		if trade.Sector == sector {
			trades = append(trades, trade)
		}
	}
	return trades
}

// countActiveTrades counts trades that haven't expired yet
func (c *Calendar) countActiveTrades() int {
	count := 0
	now := time.Now()
	for _, trade := range c.state.AllTrades {
		if trade.ExpirationDate.After(now) {
			count++
		}
	}
	return count
}

// calculateTotalRisk sums up max loss for all active trades
func (c *Calendar) calculateTotalRisk() float64 {
	total := 0.0
	now := time.Now()
	for _, trade := range c.state.AllTrades {
		if trade.ExpirationDate.After(now) {
			total += trade.MaxLoss
		}
	}
	return total
}

// loadAllTrades loads trades from storage
func (c *Calendar) loadAllTrades() {
	// TODO: Load from storage
	// trades, err := storage.LoadAllTrades()
	// if err == nil {
	//     c.state.AllTrades = trades
	// }
}

// generateSampleData generates sample trades and saves them to storage (Phase 2 feature)
func (c *Calendar) generateSampleData() {
	// Confirm with user
	dialog.ShowConfirm(
		"Generate Sample Data",
		"This will generate 10 sample trades for testing. Continue?",
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// Generate sample trades
			sampleTrades := generators.GenerateSampleTrades(10)

			// Save to storage
			if err := storage.SaveAllTrades(sampleTrades); err != nil {
				dialog.ShowError(err, c.window)
				return
			}

			// Update state
			c.state.AllTrades = sampleTrades

			// Show success message
			dialog.ShowInformation(
				"Success",
				fmt.Sprintf("Generated %d sample trades successfully!", len(sampleTrades)),
				c.window,
			)

			// Refresh display
			c.window.SetContent(c.Render())
		},
		c.window,
	)
}

// Validate checks if the screen's data is valid
func (c *Calendar) Validate() bool {
	// Calendar is always valid (display-only)
	return true
}

// GetName returns the screen name
func (c *Calendar) GetName() string {
	return "calendar"
}
