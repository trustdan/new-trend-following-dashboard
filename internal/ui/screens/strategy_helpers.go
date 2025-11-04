package screens

import (
	"image/color"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/models"
)

type strategyOption struct {
	ID          string
	Suitability models.StrategySuitability
	Strategy    models.Strategy
}

func findSector(policy *models.Policy, name string) *models.Sector {
	if policy == nil || name == "" {
		return nil
	}
	for i := range policy.Sectors {
		if policy.Sectors[i].Name == name {
			return &policy.Sectors[i]
		}
	}
	return nil
}

func getStrategySuitability(sector *models.Sector, strategyID string) models.StrategySuitability {
	if sector == nil {
		return defaultSuitability()
	}

	if suitability, exists := sector.StrategySuitability[strategyID]; exists {
		return suitability
	}

	return defaultSuitability()
}

func defaultSuitability() models.StrategySuitability {
	return models.StrategySuitability{
		Rating:                 "marginal",
		Color:                  "yellow",
		Rationale:              "Suitability data not available for this sector/strategy combination",
		RequireAcknowledgement: true,
	}
}

func buildStrategyOption(policy *models.Policy, sector *models.Sector, strategyID string) (strategyOption, bool) {
	if policy == nil {
		return strategyOption{}, false
	}
	strategy, exists := policy.Strategies[strategyID]
	if !exists {
		return strategyOption{}, false
	}

	suitability := getStrategySuitability(sector, strategyID)

	return strategyOption{
		ID:          strategyID,
		Suitability: suitability,
		Strategy:    strategy,
	}, true
}

func selectTopStrategies(policy *models.Policy, sector *models.Sector, count int) []strategyOption {
	if policy == nil || sector == nil || count <= 0 {
		return nil
	}

	results := make([]strategyOption, 0, count)
	added := make(map[string]struct{})

	addOption := func(id string) {
		if _, exists := added[id]; exists {
			return
		}
		if option, ok := buildStrategyOption(policy, sector, id); ok {
			results = append(results, option)
			added[id] = struct{}{}
		}
	}

	for _, id := range sector.AllowedStrategies {
		addOption(id)
		if len(results) >= count {
			return results[:count]
		}
	}

	remaining := collectRemainingOptions(policy, sector, added)

	for _, option := range remaining {
		if len(results) >= count {
			break
		}
		results = append(results, option)
		added[option.ID] = struct{}{}
	}

	if len(results) < count {
		for id := range policy.Strategies {
			if len(results) >= count {
				break
			}
			if _, exists := added[id]; exists {
				continue
			}
			if option, ok := buildStrategyOption(policy, sector, id); ok {
				results = append(results, option)
				added[id] = struct{}{}
			}
		}
	}

	if len(results) > count {
		return results[:count]
	}
	return results
}

func collectRemainingOptions(policy *models.Policy, sector *models.Sector, added map[string]struct{}) []strategyOption {
	options := make([]strategyOption, 0, len(sector.StrategySuitability))
	for id := range sector.StrategySuitability {
		if _, exists := added[id]; exists {
			continue
		}
		if option, ok := buildStrategyOption(policy, sector, id); ok {
			options = append(options, option)
		}
	}

	sort.Slice(options, func(i, j int) bool {
		left := options[i]
		right := options[j]

		colorRank := map[string]int{
			"green":  0,
			"yellow": 1,
			"red":    2,
		}

		ratingRank := map[string]int{
			"excellent":    0,
			"good":         1,
			"marginal":     2,
			"incompatible": 3,
		}

		leftColor, leftExists := colorRank[strings.ToLower(left.Suitability.Color)]
		if !leftExists {
			leftColor = 3
		}
		rightColor, rightExists := colorRank[strings.ToLower(right.Suitability.Color)]
		if !rightExists {
			rightColor = 3
		}
		if leftColor != rightColor {
			return leftColor < rightColor
		}

		leftRating, leftRatingExists := ratingRank[strings.ToLower(left.Suitability.Rating)]
		if !leftRatingExists {
			leftRating = 4
		}
		rightRating, rightRatingExists := ratingRank[strings.ToLower(right.Suitability.Rating)]
		if !rightRatingExists {
			rightRating = 4
		}
		if leftRating != rightRating {
			return leftRating < rightRating
		}

		return left.ID < right.ID
	})

	return options
}

func getColorIndicatorText(colorName string) string {
	switch strings.ToLower(colorName) {
	case "green":
		return "[GREEN]"
	case "yellow":
		return "[YELLOW]"
	case "red":
		return "[RED]"
	default:
		return "[UNKNOWN]"
	}
}

func getColorPalette(colorName string) (color.Color, color.Color) {
	switch strings.ToLower(colorName) {
	case "green":
		return color.RGBA{R: 240, G: 255, B: 240, A: 255}, color.RGBA{R: 0, G: 180, B: 80, A: 255}
	case "yellow":
		return color.RGBA{R: 255, G: 250, B: 230, A: 255}, color.RGBA{R: 220, G: 180, B: 0, A: 255}
	case "red":
		return color.RGBA{R: 255, G: 240, B: 240, A: 255}, color.RGBA{R: 200, G: 0, B: 0, A: 255}
	default:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}, color.RGBA{R: 128, G: 128, B: 128, A: 255}
	}
}

func newStrategyBadge(option strategyOption) fyne.CanvasObject {
	bgColor, borderColor := getColorPalette(option.Suitability.Color)

	background := canvas.NewRectangle(bgColor)
	background.SetMinSize(fyne.NewSize(0, 50))

	border := canvas.NewRectangle(borderColor)
	border.SetMinSize(fyne.NewSize(6, 50))

	title := widget.NewLabel(getColorIndicatorText(option.Suitability.Color) + " " + option.ID + " - " + option.Strategy.Label)
	title.TextStyle = fyne.TextStyle{Bold: true}

	rating := strings.ToUpper(option.Suitability.Rating)
	ratingLabel := widget.NewLabel("Sector Suitability: " + rating)
	ratingLabel.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(
		title,
		ratingLabel,
	)

	badge := container.NewStack(
		background,
		container.NewBorder(nil, nil, border, nil, container.NewPadded(content)),
	)

	return badge
}

func buildStrategyBadges(policy *models.Policy, sector models.Sector, count int) fyne.CanvasObject {
	if policy == nil {
		return widget.NewLabel("Strategy data not loaded")
	}

	options := selectTopStrategies(policy, &sector, count)
	if len(options) == 0 {
		return widget.NewLabel("No strategy guidance available")
	}

	badges := make([]fyne.CanvasObject, 0, len(options))
	for _, option := range options {
		badges = append(badges, newStrategyBadge(option))
	}

	return container.NewVBox(badges...)
}

// getSymbolForColor returns a unicode symbol for the given color
func getSymbolForColor(colorName string) string {
	switch strings.ToLower(colorName) {
	case "green":
		return "✓" // Checkmark for good strategies
	case "yellow":
		return "⚠" // Warning for marginal strategies
	case "red":
		return "🛑" // Red stop sign for incompatible strategies
	default:
		return "?"
	}
}

// buildCompactStrategyList creates a space-efficient text display of strategy suitability
func buildCompactStrategyList(policy *models.Policy, sector models.Sector, count int) fyne.CanvasObject {
	if policy == nil {
		return widget.NewLabel("Strategy data not loaded")
	}

	options := selectTopStrategies(policy, &sector, count)
	if len(options) == 0 {
		return widget.NewLabel("No strategy guidance available")
	}

	// Group strategies by color
	greenStrategies := []string{}
	yellowStrategies := []string{}
	redStrategies := []string{}

	for _, option := range options {
		strategyName := option.ID + " - " + option.Strategy.Label
		switch strings.ToLower(option.Suitability.Color) {
		case "green":
			greenStrategies = append(greenStrategies, strategyName)
		case "yellow":
			yellowStrategies = append(yellowStrategies, strategyName)
		case "red":
			redStrategies = append(redStrategies, strategyName)
		}
	}

	// Build compact display lines
	lines := make([]fyne.CanvasObject, 0, 3)

	if len(greenStrategies) > 0 {
		greenLabel := widget.NewLabel("✓ " + strings.Join(greenStrategies, " • "))
		greenLabel.Wrapping = fyne.TextWrapWord
		lines = append(lines, greenLabel)
	}

	if len(yellowStrategies) > 0 {
		yellowLabel := widget.NewLabel("⚠ " + strings.Join(yellowStrategies, " • "))
		yellowLabel.Wrapping = fyne.TextWrapWord
		lines = append(lines, yellowLabel)
	}

	if len(redStrategies) > 0 {
		redLabel := widget.NewLabel("🛑 " + strings.Join(redStrategies, " • "))
		redLabel.Wrapping = fyne.TextWrapWord
		lines = append(lines, redLabel)
	}

	return container.NewVBox(lines...)
}
