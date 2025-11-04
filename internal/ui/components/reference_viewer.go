package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/logging"
)

// ReferenceViewer displays pine script strategies and screener guides
type ReferenceViewer struct {
	window fyne.Window
}

// NewReferenceViewer creates a new reference viewer
func NewReferenceViewer(window fyne.Window) *ReferenceViewer {
	return &ReferenceViewer{
		window: window,
	}
}

// ShowReference displays the requested reference material
func (rv *ReferenceViewer) ShowReference(refType string) {
	var filePath string
	var title string

	// Map reference types to file paths
	switch refType {
	// Pine Script Strategy Guide
	case "pine_scripts_guide":
		filePath = "pine-script-strategies/PINE-SCRIPTS-GUIDE.md"
		title = "Pine Script Strategies - Complete Guide"

	// Pine Script Strategies
	case "strategy_alt10":
		filePath = "pine-script-strategies/alt10.pine"
		title = "Alt10 - Profit Targets Strategy"
	case "strategy_alt22":
		filePath = "pine-script-strategies/alt22.pine"
		title = "Alt22 - Parabolic SAR Strategy"
	case "strategy_alt26":
		filePath = "pine-script-strategies/alt26.pine"
		title = "Alt26 - Fractional Pyramiding Strategy"
	case "strategy_alt28":
		filePath = "pine-script-strategies/alt28.pine"
		title = "Alt28 - ADX Filter Strategy"
	case "strategy_alt39":
		filePath = "pine-script-strategies/alt39.pine"
		title = "Alt39 - Age-Based Targets Strategy"
	case "strategy_alt43":
		filePath = "pine-script-strategies/seykota_alt43_volatility_adaptive_targets.pine"
		title = "Alt43 - Volatility Adaptive Targets Strategy"
	case "strategy_alt45":
		filePath = "pine-script-strategies/seykota_alt45_dual_momentum_confirmation.pine"
		title = "Alt45 - Dual Momentum Confirmation Strategy"
	case "strategy_alt46":
		filePath = "pine-script-strategies/seykota_alt46_sector_adaptive_parameters.pine"
		title = "Alt46 - Sector Adaptive Parameters Strategy"
	case "strategy_alt47":
		filePath = "pine-script-strategies/seykota_alt47_momentum_scaled_sizing.pine"
		title = "Alt47 - Momentum Scaled Sizing Strategy"
	case "strategy_turtle_core":
		filePath = "pine-script-strategies/turtle-core-v2.2.pine"
		title = "Turtle Core V2.2 Strategy (Baseline)"

	// Finviz Screeners
	case "screener_master":
		filePath = "screeners/MASTER-SCREENER-GUIDE.md"
		title = "Master Finviz Screener Guide"
	case "screener_daily":
		filePath = "screeners/DAILY-CHEAT-SHEET.md"
		title = "Daily Screener Cheat Sheet"
	case "screener_decision_tree":
		filePath = "screeners/SCREENER-DECISION-TREE.md"
		title = "Screener Decision Tree"
	case "screener_start":
		filePath = "screeners/START-HERE.md"
		title = "Finviz Screeners - Start Here"

	default:
		dialog.ShowError(fmt.Errorf("unknown reference type: %s", refType), rv.window)
		return
	}

	// Read the file
	content, err := rv.readFile(filePath)
	if err != nil {
		logging.ErrorLogger.Printf("Failed to read reference file %s: %v", filePath, err)
		dialog.ShowError(fmt.Errorf("failed to load reference: %v", err), rv.window)
		return
	}

	// Create formatted content display
	rv.showDialog(title, content)
}

// readFile reads a file and returns its content
func (rv *ReferenceViewer) readFile(filePath string) (string, error) {
	// Try multiple possible locations
	locations := []string{
		filePath,                        // Direct path
		filepath.Join("..", filePath),   // One level up
		filepath.Join("../..", filePath), // Two levels up
	}

	// Also check relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		locations = append(locations, filepath.Join(exeDir, filePath))
	}

	var lastErr error
	for _, loc := range locations {
		content, err := os.ReadFile(loc)
		if err == nil {
			logging.DebugLogger.Printf("Found reference file at: %s", loc)
			return string(content), nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("file not found in any location: %w", lastErr)
}

// showDialog displays the reference content in a scrollable dialog
func (rv *ReferenceViewer) showDialog(title, content string) {
	// Detect if content is Pine Script or Markdown
	isPineScript := strings.HasSuffix(title, "Strategy") && !strings.Contains(title, "Guide")

	// Add header explaining the content
	var header string
	if isPineScript {
		header = "📊 Pine Script Strategy Code\n\n" +
			"Copy this code into TradingView's Pine Editor to use this strategy.\n" +
			"Click the Pine Editor tab, paste the code, and click 'Add to Chart'.\n\n" +
			strings.Repeat("─", 80) + "\n\n"
	} else {
		header = "" // No header for markdown guides - they have their own formatting
	}

	fullContent := header + content

	// Create a rich text widget for displaying the content
	textWidget := widget.NewRichTextFromMarkdown(fullContent)
	textWidget.Wrapping = fyne.TextWrapWord

	// Create scrollable container
	scroll := container.NewScroll(textWidget)
	scroll.SetMinSize(fyne.NewSize(800, 600))

	// Create dialog with custom content
	dlg := dialog.NewCustom(title, "Close", scroll, rv.window)
	dlg.Resize(fyne.NewSize(900, 700))
	dlg.Show()
}

// ShowStrategy is a convenience method for showing pine script strategies
func (rv *ReferenceViewer) ShowStrategy(strategyID string) {
	rv.ShowReference("strategy_" + strings.ToLower(strategyID))
}

// ShowScreenerGuide is a convenience method for showing screener guides
func (rv *ReferenceViewer) ShowScreenerGuide(guideType string) {
	rv.ShowReference("screener_" + strings.ToLower(guideType))
}
