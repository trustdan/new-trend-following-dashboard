package help

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// HelpContent contains the help text for each screen
type HelpContent struct {
	Title       string
	Description string
	Steps       []string
	Tips        []string
}

// GetHelpForScreen returns context-sensitive help content
func GetHelpForScreen(screenName string) HelpContent {
	helpMap := map[string]HelpContent{
		"sector_selection": {
			Title:       "Screen 1: Sector Selection",
			Description: "Choose which market sector you want to trade. Strategy performance varies significantly by sector based on 293 backtests.",
			Steps: []string{
				"1. Review sector performance indicators (green = good backtest results)",
				"2. Select a sector by clicking on it",
				"3. Click 'Continue' to proceed to screener",
			},
			Tips: []string{
				"⚠ Blocked sectors (e.g., Utilities) have 0% backtest success - trades are prevented",
				"⚠ Warned sectors (Energy, Real Estate) show caution icon - review carefully",
				"✓ Healthcare and Technology sectors show strongest trend-following performance",
			},
		},
		"screener_launch": {
			Title:       "Screen 2: Screener Launch",
			Description: "Launch FINVIZ screeners to find trade candidates. Use Universe (weekly) for broad search, Situational (daily) for specific setups.",
			Steps: []string{
				"1. Click 'Universe Screener' for broad weekly scan (30-60 quality stocks)",
				"2. Or click situational screeners for specific patterns:",
				"   • Pullback: Oversold stocks in uptrends (RSI < 40)",
				"   • Breakout: New 52-week highs",
				"   • Golden Cross: SMA50 crossing above SMA200",
				"3. Review results in your browser",
				"4. Return here and enter your chosen ticker on next screen",
			},
			Tips: []string{
				"✓ Screeners open in chart view (v=211 parameter) for easy pattern recognition",
				"✓ Universe screeners run Monday mornings, Situational daily before market open",
				"✓ Look for stocks above SMA200 with confirmed Donchian breakouts",
			},
		},
		"ticker_entry": {
			Title:       "Screen 3: Ticker & Strategy Entry",
			Description: "Enter your chosen ticker and select a strategy. The 5-minute cooldown timer starts when you proceed.",
			Steps: []string{
				"1. Type ticker symbol (e.g., UNH, MSFT)",
				"2. Select strategy from dropdown (filtered by sector)",
				"3. Click 'Continue' to start anti-impulsivity cooldown",
			},
			Tips: []string{
				"⏱ Cooldown timer cannot be bypassed - this prevents impulsive trades",
				"✓ Strategy dropdown shows ONLY strategies validated for selected sector",
				"✓ Alt10 (Profit Targets) has highest success rate: 76.19% across all sectors",
			},
		},
		"checklist": {
			Title:       "Screen 4: Anti-Impulsivity Checklist",
			Description: "Complete required criteria before trading. This forced pause improves decision quality.",
			Steps: []string{
				"1. Wait for cooldown timer to complete (5 minutes)",
				"2. Check all 5 REQUIRED items:",
				"   • Signal confirmed (price above SMA200, Donchian breakout)",
				"   • Risk acceptable (stop-loss placement reasonable)",
				"   • Options match strategy (hold time vs expiration)",
				"   • Exit plan clear (profit targets & stop defined)",
				"   • Emotional readiness (calm, not revenge trading)",
				"3. Optionally check OPTIONAL items for higher conviction",
				"4. Click 'Continue' once all required items checked",
			},
			Tips: []string{
				"✓ Think of this as a pre-flight checklist - catches mistakes before takeoff",
				"✓ Optional items increase conviction score → affects position sizing",
				"⚠ If you can't check all 5 required items, cancel the trade",
			},
		},
		"position_sizing": {
			Title:       "Screen 5: Position Sizing (Poker-Bet System)",
			Description: "Size position based on conviction level. Uses poker-bet multipliers to scale risk.",
			Steps: []string{
				"1. Rate your conviction: 5-8 scale",
				"   • 5 = Weak signal (0.5× sizing)",
				"   • 6 = Below average (0.75× sizing)",
				"   • 7 = Standard (1.0× sizing) ← Most common",
				"   • 8 = Strong signal (1.25× sizing)",
				"2. Review calculated position size",
				"3. Adjust if needed based on account size",
				"4. Click 'Continue' to proceed to heat check",
			},
			Tips: []string{
				"✓ Most trades should be conviction 7 (standard sizing)",
				"⚠ Don't use 8 unless ALL checklist items passed AND strong technical setup",
				"✓ When in doubt, use 6 (75% sizing) - reduces risk on marginal trades",
			},
		},
		"heat_check": {
			Title:       "Screen 6: Portfolio Heat Check",
			Description: "Enforce diversification limits to prevent concentration risk. Trades exceeding caps are blocked.",
			Steps: []string{
				"1. Review current portfolio heat by sector",
				"2. Check if new trade exceeds limits:",
				"   • 4% total portfolio heat (all sectors combined)",
				"   • 1.5% per-sector heat (e.g., max 1.5% in Healthcare)",
				"3. If limits exceeded, close existing position or skip trade",
				"4. If limits OK, click 'Continue' to proceed",
			},
			Tips: []string{
				"✓ Heat limits are non-negotiable - prevent overconcentration disasters",
				"⚠ If at limits, look for trades in different sector to maintain diversification",
				"✓ Closing losing position frees up heat for new opportunity",
			},
		},
		"trade_entry": {
			Title:       "Screen 7: Options Strategy Entry",
			Description: "Select specific options structure and enter strike prices.",
			Steps: []string{
				"1. Select options strategy type (24 available)",
				"2. Enter strike prices based on strategy:",
				"   • Bull call spread: Long strike + Short strike",
				"   • Iron condor: 4 strikes",
				"   • Single options: 1 strike",
				"3. Enter expiration date (typically 30-60 days out)",
				"4. Enter premium paid/received",
				"5. Click 'Continue' to add trade to calendar",
			},
			Tips: []string{
				"✓ Bull call spreads: Most common for trend-following (limited risk)",
				"✓ Expiration should match strategy hold time (Alt10 = 3-10 weeks)",
				"⚠ Avoid expiration <30 days (time decay too fast)",
			},
		},
		"calendar": {
			Title:       "Screen 8: Trade Calendar (Horserace View)",
			Description: "Visualize all trades across time and sectors. Green/red bars show P&L, yellow = expiring soon.",
			Steps: []string{
				"1. Review timeline (-14 days to +84 days)",
				"2. Identify sector concentration (vertical stacking)",
				"3. Check for expiring trades (yellow bars <7 days)",
				"4. Click '+ New Trade' to add another trade",
				"5. Click 'Refresh' to reload data",
			},
			Tips: []string{
				"✓ Horizontal bars = trade duration (entry → expiration)",
				"✓ Color coding: Blue=Active, Green=Profitable, Red=Losing, Yellow=Expiring",
				"✓ Y-axis groups by sector for easy concentration checks",
			},
		},
		"trade_management": {
			Title:       "Screen 9: Trade Management (Phase 2 Feature)",
			Description: "Edit or delete past trades. Filter by status.",
			Steps: []string{
				"1. Select filter: Show All / Active Only / Closed Only",
				"2. Click 'Edit' to modify trade details (P&L, status)",
				"3. Click 'Delete' to remove trade (requires confirmation)",
				"4. Click 'Back to Calendar' to return",
			},
			Tips: []string{
				"⚠ This feature must be enabled via feature.flags.json",
				"✓ Use 'Active Only' filter to focus on open positions",
				"✓ Edit P&L when closing trades manually",
			},
		},
	}

	content, exists := helpMap[screenName]
	if !exists {
		return HelpContent{
			Title:       "TF-Engine 2.0 Help",
			Description: "Trend-Following Options Trading Decision Support System",
			Steps: []string{
				"This application guides you through an 8-screen workflow:",
				"1. Sector Selection - Choose market sector",
				"2. Screener Launch - Find trade candidates with FINVIZ",
				"3. Ticker Entry - Enter symbol and select strategy",
				"4. Checklist - Complete anti-impulsivity criteria",
				"5. Position Sizing - Size based on conviction (poker-bet system)",
				"6. Heat Check - Verify portfolio diversification limits",
				"7. Trade Entry - Enter options strike prices and expiration",
				"8. Calendar - View all trades on timeline",
			},
			Tips: []string{
				"✓ Based on 293 backtests across 14 strategies and 21 securities",
				"✓ Behavioral guardrails prevent emotional trading (cooldown, checklist, heat limits)",
				"✓ Sector-first workflow: Strategy performance is sector-dependent",
				"✓ Use Help menu on any screen for context-specific guidance",
			},
		}
	}

	return content
}

// ShowHelpDialog displays a help dialog for the current screen
func ShowHelpDialog(screenName string, window fyne.Window) {
	help := GetHelpForScreen(screenName)

	// Title
	title := widget.NewLabel(help.Title)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Wrapping = fyne.TextWrapWord

	// Description
	description := widget.NewLabel(help.Description)
	description.Wrapping = fyne.TextWrapWord

	// Steps section
	stepsLabel := widget.NewLabel("Steps:")
	stepsLabel.TextStyle = fyne.TextStyle{Bold: true}

	stepsText := ""
	for _, step := range help.Steps {
		stepsText += step + "\n"
	}
	steps := widget.NewLabel(stepsText)
	steps.Wrapping = fyne.TextWrapWord

	// Tips section
	tipsLabel := widget.NewLabel("Tips:")
	tipsLabel.TextStyle = fyne.TextStyle{Bold: true}

	tipsText := ""
	for _, tip := range help.Tips {
		tipsText += tip + "\n"
	}
	tips := widget.NewLabel(tipsText)
	tips.Wrapping = fyne.TextWrapWord

	// Content container
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		description,
		widget.NewSeparator(),
		stepsLabel,
		steps,
		widget.NewSeparator(),
		tipsLabel,
		tips,
	)

	// Scrollable container
	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	// Show dialog
	dialog.ShowCustom("Help", "Close", scroll, window)
}

// ShowWelcomeScreen displays the welcome/onboarding screen
func ShowWelcomeScreen(window fyne.Window, onComplete func(dontShowAgain bool)) {
	// Title
	title := widget.NewLabel("Welcome to TF-Engine 2.0")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Main description
	desc := widget.NewLabel(`
TF-Engine is a Trend-Following Options Trading Decision Support System based on 293 validated backtests.

Key Features:
• Sector-First Workflow: Strategy performance is highly sector-dependent
• Anti-Impulsivity Guardrails: 120s cooldown, 5-gate checklist, portfolio heat limits
• Poker-Bet Position Sizing: Scale risk based on conviction (5-8 rating)
• Horserace Calendar: Visualize all trades across time and sectors

This app helps you make systematic, disciplined trading decisions.`)
	desc.Wrapping = fyne.TextWrapWord

	// 8-screen workflow
	workflowLabel := widget.NewLabel("8-Screen Workflow:")
	workflowLabel.TextStyle = fyne.TextStyle{Bold: true}

	workflow := widget.NewLabel(`
1. Sector Selection - Choose Healthcare, Technology, etc.
2. Screener Launch - Find candidates with FINVIZ
3. Ticker & Strategy Entry - Select ticker and strategy
4. Anti-Impulsivity Checklist - 5 required gates + cooldown
5. Position Sizing - Conviction-based (poker system)
6. Portfolio Heat Check - Enforce diversification limits
7. Options Trade Entry - Enter strikes and expiration
8. Trade Calendar - Visualize portfolio timeline`)
	workflow.Wrapping = fyne.TextWrapWord

	// Don't show again checkbox
	dontShowCheck := widget.NewCheck("Don't show this again", nil)

	// Buttons
	getStartedBtn := widget.NewButton("Get Started", func() {
		onComplete(dontShowCheck.Checked)
	})
	getStartedBtn.Importance = widget.HighImportance

	helpBtn := widget.NewButton("Learn More", func() {
		ShowHelpDialog("welcome", window)
	})

	buttons := container.NewHBox(getStartedBtn, helpBtn)

	// Content
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		desc,
		widget.NewSeparator(),
		workflowLabel,
		workflow,
		widget.NewSeparator(),
		dontShowCheck,
		buttons,
	)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(650, 500))

	// Show as custom dialog
	welcomeDialog := dialog.NewCustom("Welcome", "Close", scroll, window)
	welcomeDialog.Show()
}
