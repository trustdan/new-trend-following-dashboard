# TF-Engine GUI: Overview Plan
## The Forest Over the Trees

**TF = Trend Following** - Systematic trading using Donchian breakouts, ATR-based sizing, and mechanical exits

**Created:** 2025-10-29
**Purpose:** Strategic foundation for the entire GUI trading platform
**Philosophy:** Anti-impulsivity through systematic discipline enforcement
**Status:** 🎯 Planning Phase

---

## Table of Contents

1. [Vision & Philosophy](#vision--philosophy)
2. [System Architecture](#system-architecture)
3. [User Workflow](#user-workflow)
4. [Technical Components](#technical-components)
5. [Implementation Strategy](#implementation-strategy)
6. [Proof-of-Concept Approach](#proof-of-concept-approach)
7. [Behavioral Specifications (Gherkin)](#behavioral-specifications-gherkin)
8. [Success Criteria](#success-criteria)
9. [Risk Management](#risk-management)
10. [Next Steps](#next-steps)

---

## Vision & Philosophy

### The Purpose

This is **not** a trading platform. This is a **discipline enforcement system** that happens to be implemented as software.

**What is Trend Following?**
- Follow existing price trends (the "tide") rather than predict reversals
- Enter on breakouts (new highs/lows over N-bar period)
- Exit mechanically (opposite breakout or ATR-based stop)
- Add to winners (pyramid), cut losers quickly
- No opinion, no prediction - just systematic execution

From `docs/anti-impulsivity.md`:

> **Trade the tide, not the splash.** A breakout and a mechanical exit are the core.

**The TF-Engine (Trend Following Engine)** implements this philosophy through:
- 55-bar Donchian breakouts (Ed Seykota's System-2)
- ATR-based position sizing (Van Tharp method)
- 10-bar opposite Donchian exits (mechanical)
- 0.5×N pyramiding (add to winners)
- 2×N protective stops (cut losers)

The system exists to make impulsive trading **impossible** while making systematic, trend-following trading **effortless**.

### Core Principles (Non-Negotiable)

1. **Friction where it matters** - Hard gates for signal, risk, liquidity, exit, behavior
2. **Nudge for better trades** - Optional quality items affect score, not permission
3. **Immediate feedback** - Large 3-state banner (RED/YELLOW/GREEN) updates live
4. **Journal while you decide** - One-click logging of full trade plan
5. **Calendar awareness** - Rolling 10-week sector view prevents basket crowding
6. **No backdoors** - Cannot bypass gates, skip cooldowns, or override caps

### The 5 Hard Gates (Enforced by Backend)

Every trade must pass ALL gates before execution. These are **immutable**.

```gherkin
Given a trader wants to execute a trade
Then ALL of the following gates must be GREEN:

Gate 1: Signal Confirmed
  - 55-bar Donchian breakout (long: close > donHi[1], short: close < donLo[1])
  - Verified against TradingView chart using Ed-Seykota.pine script
  - No subjective pattern analysis allowed

Gate 2: Risk/Size Calculated
  - Per-unit risk = equity × risk% (typically 0.75-1.0%)
  - Stop distance = 2.0 × N (where N = ATR(20))
  - Initial stop = entry ± (2.0 × N)
  - Shares/contracts = floor(unitRisk$ ÷ stopDistance)
  - Add-on levels = entry + (k × 0.5 × N) for k=1,2,3

Gate 3: Options Requirements Met (if applicable)
  - DTE: 60-90 days at entry
  - Roll/close trigger: ~21 DTE
  - Liquidity: bid-ask spread < 10% of mid
  - Open interest: > 100 contracts
  - If stock trade: liquidity check via avg volume

Gate 4: Exits Defined
  - Primary: 10-bar opposite Donchian (donLo[1] for longs, donHi[1] for shorts)
  - Secondary: Initial stop = entry ± (2.0 × N)
  - Exit trigger: CLOSER of Donchian or 2×N stop
  - No discretionary exits allowed

Gate 5: Behavior Constraints Honored
  - 2-minute cool-off elapsed since checklist evaluation
  - Ticker not on sector cooldown (from recent loss)
  - Portfolio heat < 4.0% of equity
  - Sector bucket heat < 1.5% of equity
  - No intraday overrides permitted
```

### Banner System (Visual Discipline Enforcement)

**Size:** Minimum 20% of screen height. **Cannot be missed.**

**Visual Design:** Sleek, modern gradients with smooth transitions. No flat colors.

```
┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    🛑  DO NOT TRADE  🛑                ║    │
│  ║    One or more REQUIRED gates failed   ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Deep Red → Dark Crimson]           │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    ⚠️  CAUTION  ⚠️                     ║    │
│  ║    Quality score below threshold       ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Amber → Golden Yellow]             │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  ╔════════════════════════════════════════╗    │
│  ║                                        ║    │
│  ║    ✓  OK TO TRADE  ✓                  ║    │
│  ║    All gates pass • Quality met        ║    │
│  ║                                        ║    │
│  ╚════════════════════════════════════════╝    │
│  [Gradient: Emerald → Forest Green]            │
└─────────────────────────────────────────────────┘
```

**State Logic:**
- RED: Any required gate fails → **Absolutely no trade**
- YELLOW: All required pass, quality score < threshold (default 3.0) → **Caution**
- GREEN: All required pass, quality score ≥ threshold → **Proceed**

**Gradient Specifications:**
- **RED Banner:** Linear gradient from `#DC2626` (deep red) to `#991B1B` (dark crimson)
- **YELLOW Banner:** Linear gradient from `#F59E0B` (amber) to `#FBBF24` (golden)
- **GREEN Banner:** Linear gradient from `#10B981` (emerald) to `#059669` (forest green)
- **Transitions:** Smooth 0.3s ease-in-out animation when state changes
- **Text:** Large, bold, white text with subtle shadow for readability
- **Border:** Subtle rounded corners with soft glow effect in banner color

---

## System Architecture

### High-Level Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                    DESKTOP APPLICATION                         │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              FYNE FRONTEND (Native GUI)                  │ │
│  │                                                          │ │
│  │  Screens:                                                │ │
│  │  - Dashboard (Positions, Candidates, Heat)               │ │
│  │  - FINVIZ Scanner (One-click daily scan)                 │ │
│  │  - Candidate Import (Review & approve tickers)           │ │
│  │  - Checklist (Required gates + quality items)            │ │
│  │  - Position Sizer (ATR-based Van Tharp calculations)     │ │
│  │  - Heat Monitor (Portfolio 4% / Bucket 1.5% caps)        │ │
│  │  - Trade Entry (Final 5-gate validation)                 │ │
│  │  - Calendar View (10-week sector diversification)        │ │
│  │  - TradingView Integration (Open charts w/ script)       │ │
│  │                                                          │ │
│  │  State Management: Direct backend calls (in-process)     │ │
│  │  UI Framework: Fyne v2 (native Go GUI)                   │ │
│  │  Styling: Custom theme (British Racing Green accent)     │ │
│  └──────────────────────────────────────────────────────────┘ │
│                           │                                    │
│                           │ Direct Function Calls              │
│                           │                                    │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │                 GO BACKEND (tf-engine)                   │ │
│  │                                                          │ │
│  │  ┌────────────────┐  ┌────────────────┐                │ │
│  │  │  HTTP Server   │  │  Static Assets │                │ │
│  │  │  (Gin/Chi)     │  │  (Embedded)    │                │ │
│  │  └────────────────┘  └────────────────┘                │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │              DOMAIN LOGIC (✅ Complete)            │ │ │
│  │  │                                                    │ │ │
│  │  │  - Position Sizing (stock, opt-delta, opt-contracts) │ │
│  │  │  - Checklist Evaluation (6 required + 4 optional)  │ │ │
│  │  │  - Heat Management (portfolio & sector caps)       │ │ │
│  │  │  - Gate Enforcement (all 5 gates)                  │ │ │
│  │  │  - Settings Management (equity, risk%, caps)       │ │ │
│  │  │  - Candidate Management (FINVIZ imports)           │ │ │
│  │  │  - Cooldown Tracking (sector & ticker)             │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │          STORAGE (SQLite - ✅ Complete)            │ │ │
│  │  │                                                    │ │ │
│  │  │  Tables: settings, positions, decisions,          │ │ │
│  │  │          candidates, cooldowns, evaluations       │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │        SCRAPE (FINVIZ - ✅ Complete)               │ │ │
│  │  │                                                    │ │ │
│  │  │  - Parse FINVIZ screener pages                     │ │ │
│  │  │  - Extract ticker symbols                          │ │ │
│  │  │  - Support presets (TF_BREAKOUT_LONG, etc.)        │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  Database: trading.db (SQLite, single file)                   │
│  Deployment: Single binary (tf-gui.exe on Windows)            │
│  Opens: Native desktop window (no browser required)           │
└───────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Backend (✅ Complete - tf-engine):**
- Language: Go 1.24+
- Database: SQLite via modernc.org/sqlite
- Web Scraping: golang.org/x/net
- Logging: Structured logging with levels (DEBUG, INFO, WARN, ERROR)
  - Log to file: `logs/tf-gui.log` (rotated daily)
  - Log to console: Configurable verbosity
  - Feature usage tracking (which features are used, how often)
  - Error tracking (full stack traces, context)
- Testing: Go standard testing + testify

**Frontend (✅ Complete - Fyne GUI):**
- Language: Go
- Framework: Fyne v2 (native cross-platform GUI)
- UI Components: Fyne widgets + custom containers
- Styling: Custom theme (British Racing Green #00352B accent)
- Color System: Day/night mode with theme toggle
- Design Language: Modern Material Design with custom colors
- State: Direct backend function calls (in-process, no HTTP)
- Logging: File logging with rotation
  - User actions tracked (button clicks, navigation, form submissions)
  - Component interactions logged (tab switches, dialogs, errors)
  - Performance metrics (calculation times, render delays)
  - Feature usage analytics (which screens visited, how often)
- Build: Go build (cross-compilation support)
- Output: Single binary executable

**Integration:**
- Frontend → Backend: Direct function calls (same process)
- No HTTP layer needed (Fyne GUI calls domain logic directly)
- Deployment: Single binary (tf-gui.exe on Windows)
- Desktop: Native window with Fyne canvas rendering

**Development Environment:**
- OS: Linux (WSL2/Kali) for development
- Build Target: Windows .exe (cross-compiled)
- Git: Windows only (per `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md`)
- Handoff: Export zip → Windows Git repo
- Logging: Comprehensive logging for all features (debug, feature evaluation, performance)

---

## Visual Design Philosophy

### Modern, Clean, and Functional

**Core Principle:** The interface should be clear and purposeful. Clean design with strategic color usage encourages regular engagement with the systematic process.

**Design Language:**
- **Strategic color usage** - British Racing Green (#00352B) for accents and important elements
- **Material Design principles** - Clean, modern aesthetic with Fyne's native widgets
- **Clear visual hierarchy** - Important elements (banner, buttons) are prominent
- **Generous spacing** - Padding and margins for readability
- **Consistent typography** - Fyne's built-in font system with bold for emphasis

### Color System

**Custom Fyne Theme - Day Mode (Light):**
```go
// Background colors
Background: color.White                    // Main background
Card: color.NRGBA{R: 249, G: 250, B: 251, A: 255}  // Cards, panels

// Text colors
Foreground: color.Black                    // Main text
PlaceHolder: color.Gray{Y: 0xAA}          // Placeholder text

// State colors
Primary: color.NRGBA{R: 0, G: 53, B: 43, A: 255}   // British Racing Green
Success: color.NRGBA{R: 16, G: 185, B: 129, A: 255} // GREEN banner
Warning: color.NRGBA{R: 245, G: 158, B: 11, A: 255} // YELLOW banner
Error: color.NRGBA{R: 220, G: 38, B: 38, A: 255}    // RED banner
```

**Custom Fyne Theme - Night Mode (Dark):**
```go
// Background colors
Background: color.NRGBA{R: 15, G: 23, B: 42, A: 255}  // slate-900
Card: color.NRGBA{R: 30, G: 41, B: 59, A: 255}       // slate-800

// Text colors
Foreground: color.NRGBA{R: 241, G: 245, B: 249, A: 255}  // slate-100
PlaceHolder: color.NRGBA{R: 148, G: 163, B: 184, A: 255} // slate-400

// State colors (same as light mode for consistency)
Primary: color.NRGBA{R: 0, G: 53, B: 43, A: 255}   // British Racing Green
Success: color.NRGBA{R: 16, G: 185, B: 129, A: 255}
Warning: color.NRGBA{R: 245, G: 158, B: 11, A: 255}
Error: color.NRGBA{R: 220, G: 38, B: 38, A: 255}
```

**Theme Toggle:**
- Icon: Sun ☀️ (light mode) / Moon 🌙 (dark mode)
- Location: Top toolbar with other controls
- Transition: Immediate (Fyne's built-in theme switching)
- Persistence: Saved to app preferences
- Default: System preference detected via Fyne

### Component Design Guidelines

**Buttons:**
```go
// Primary Action (e.g., "SAVE GO DECISION")
widget.NewButtonWithIcon("Save GO", theme.ConfirmIcon(), handler)
- Importance: widget.HighImportance (green background)
- Text: Bold, uppercase
- Icon: theme icon (confirm, cancel, etc.)
- Disabled state: Grayed out automatically

// Danger Action (e.g., "SAVE NO-GO DECISION")
widget.NewButtonWithIcon("Save NO-GO", theme.CancelIcon(), handler)
- Importance: widget.DangerImportance (red background)

// Secondary Action
widget.NewButton("Cancel", handler)
- Importance: widget.LowImportance (transparent background)
```

**Cards/Panels:**
```go
container.NewPadded(
    widget.NewCard("Title", "Subtitle", content)
)
- Uses Fyne's Card widget for consistent styling
- Automatic borders and shadows
- Padding: Standard Fyne padding (theme.Padding())
```

**Forms:**
```go
// Input Fields
widget.NewEntry()
- Placeholder text supported
- Validation via OnChanged
- Focus highlighting automatic

// Checkboxes
widget.NewCheck("Label", onChanged)
- Material Design styling
- Check animation built-in

// Dropdowns
widget.NewSelect(options, onChanged)
- Native dropdown styling
- Keyboard navigation
```

**Banner (The Star of the Show):**
```go
// Large colored rectangle with text overlay
canvas.NewRectangle(bannerColor)  // RED/YELLOW/GREEN
canvas.NewText("BANNER MESSAGE", color.White)
- Height: 150-200 pixels (fixed)
- Width: Full container width
- Centered text, 24-32pt bold
- Color changes based on state
- Icon + text layout
```

**Tables:**
```go
widget.NewTable(
    length func() (int, int),
    create func() fyne.CanvasObject,
    update func(widget.TableCellID, fyne.CanvasObject)
)
- Built-in alternating row colors (via theme)
- Header row support
- Cell padding: theme.Padding()
```

### Icon Usage

**Icon Library:** Fyne's built-in theme icons + custom icons as needed

**Sizes:**
- Standard: theme.IconSize() (typically 20px)
- Large: For banners and headers (32-48px custom)

**Usage:**
```go
theme.ConfirmIcon()  // Checkmark
theme.CancelIcon()   // X mark
theme.InfoIcon()     // Information
theme.WarningIcon()  // Warning triangle
theme.ErrorIcon()    // Error/alert
theme.HomeIcon()     // Dashboard
theme.DocumentIcon() // Checklist
// etc.
```

### Spacing System

Fyne uses theme-based spacing:

```go
theme.Padding()       // Standard padding (4px)
theme.InnerPadding()  // Inner padding
theme.LineSpacing()   // Line spacing
```

Custom spacing where needed:
```go
container.NewPadded(content)  // Automatic padding
container.NewVBox(           // Vertical with spacing
    widget.NewLabel("Item 1"),
    widget.NewSeparator(),
    widget.NewLabel("Item 2"),
)
```

### Window Sizing

```go
// Minimum window size
window.Resize(fyne.NewSize(1200, 800))
window.SetFixedSize(false)  // Allow resizing

// Optimized for desktop trading
// Minimum: 1200x800
// Recommended: 1400x900 or larger
```

### Typography

Fyne uses theme-based typography:

```go
// Text sizes
theme.TextSize()              // Standard (14px)
theme.TextHeadingSize()       // Headings (18px)
theme.TextSubHeadingSize()    // Subheadings (16px)

// Custom sizes for banner
canvas.NewText("BANNER", color.White)
text.TextSize = 32  // Large banner text
text.TextStyle = fyne.TextStyle{Bold: true}
```

---

## User Workflow

### Daily Trading Routine (The Happy Path)

This is the **intended workflow** that the system enforces:

```gherkin
Feature: Daily trend-following workflow
  As a systematic trader
  I want to execute a disciplined daily routine
  So that I trade the tide, not the splash

  Background:
    Given I am using the TF-Engine GUI application
    And the backend tf-engine is running
    And I have my TradingView charts ready
    And my Ed-Seykota.pine script is loaded on TradingView

  Scenario: Morning scan and candidate identification
    When I open the Dashboard at market open
    Then I see my current open positions with heat % and stops
    And I see yesterday's candidates (if any remain)
    And I see the total portfolio heat vs 4% cap
    And I see any sector buckets on cooldown

    When I click "Run Daily FINVIZ Scan"
    Then the backend scrapes my preset FINVIZ screener
    And returns a list of tickers meeting technical criteria
    And displays count: "Found 23 candidates from TF_BREAKOUT_LONG"

    When I click "Review Candidates"
    Then I see a table of tickers with basic info
    And I can check/uncheck tickers to import
    And I can see which sectors are represented
    And I can see which are on cooldown (grayed out)

    When I select 12 tickers and click "Import Selected"
    Then those tickers are saved to the candidates table
    And they appear in my "Today's Candidates" list
    And I see notification: "12 candidates imported for 2025-10-29"

  Scenario: Chart analysis in TradingView
    Given I have imported 12 candidates
    When I click "Open in TradingView" next to ticker AAPL
    Then a new browser tab opens to AAPL chart on TradingView
    And the Ed-Seykota.pine script is already applied (if preset in TV)

    # Manual step (outside system):
    And I manually verify the 55-bar Donchian breakout
    And I note the current N (ATR) value from the script
    And I note the entry price (current close or limit)
    And I identify the sector bucket (e.g., "Tech/Comm")

    When I return to the TF-Engine GUI
    Then I am ready to fill out the checklist

  Scenario: Checklist evaluation (gates + quality)
    Given I have analyzed AAPL in TradingView
    When I navigate to the Checklist screen

    And I enter ticker: "AAPL"
    And I enter sector: "Tech/Comm" (dropdown)
    And I enter entry price: 180.50
    And I enter N (ATR): 2.35 (from TradingView)

    Then the system pre-calculates:
      - Stop distance: 2.0 × 2.35 = 4.70
      - Initial stop: 180.50 - 4.70 = 175.80 (for long)
      - Add-on levels: [182.68, 184.85, 187.03] (entry + k×0.5×N)

    # Required Gates (must check all):
    When I check "✓ Signal: 55-bar breakout confirmed (close > donHi[1])"
    And I check "✓ Risk/Size: Will use 2×N stop, add every 0.5×N, max 4 units"
    And I check "✓ Liquidity: Stock avg volume > 1M shares (or options OI > 100)"
    And I check "✓ Exits: Will exit on 10-bar opposite Donchian OR 2×N stop"
    And I check "✓ Behavior: Not on cooldown, heat caps OK, will honor 2-min timer"

    # Optional Quality Items (improve score):
    And I check "✓ Regime OK: SPY > 200 SMA (market favorable for longs)"
    And I check "✓ No Chase: Entry within 2N of 20-EMA (not extended)"
    And I check "✓ Earnings OK: No earnings within next 2 weeks (for long options)"
    And I enter journal note: "Clean breakout, strong volume, sector underweight"

    Then the quality score updates to 4 / 4 optional items
    And the banner changes from RED → YELLOW → GREEN
    And the banner displays: "🟢 OK TO TRADE 🟢"

    When I click "Save Evaluation"
    Then the backend records the evaluation timestamp
    And the 2-minute impulse brake timer starts
    And I see: "Evaluation saved. Cool-off period: 2 minutes remaining."

  Scenario: Position sizing calculation
    Given my checklist evaluation is GREEN
    And 2 minutes have NOT yet elapsed (still in cool-off)

    When I navigate to the Position Sizing screen
    And the form is pre-filled with:
      - Ticker: AAPL
      - Entry: 180.50
      - N (ATR): 2.35
      - Method: "stock" (dropdown: stock, opt-delta-atr, opt-contracts)
      - Risk % per unit: 0.75% (from settings)
      - Equity: $100,000 (from settings)
      - Max units: 4 (from settings)

    When I click "Calculate Position Size"
    Then the backend calculates using Van Tharp method:
      - Risk per unit: $100,000 × 0.0075 = $750
      - Stop distance: 2.0 × 2.35 = 4.70
      - Shares per unit: floor(750 ÷ 4.70) = 159 shares
      - Actual risk: 159 × 4.70 = $747.30 (≤ $750 ✓)
      - Initial stop: $175.80
      - Add-on schedule:
        * Unit 1: 159 shares @ 180.50 (now)
        * Unit 2: 159 shares @ 182.68 (+0.5N)
        * Unit 3: 159 shares @ 184.85 (+1.0N)
        * Unit 4: 159 shares @ 187.03 (+1.5N)

    And the results are displayed prominently
    And I see a warning if shares × price > 25% of equity (concentration risk)
    And I can adjust entry or method and recalculate

    When I click "Save Position Plan"
    Then the plan is saved to the database
    And it's linked to my checklist evaluation

  Scenario: Heat check before final entry
    Given my position sizing is calculated (159 shares AAPL, $747 risk)
    And my 2-minute cool-off is still active

    When I navigate to the Heat Check screen
    And I see my current portfolio heat: $2,890 / $4,000 cap (72.25%)
    And I see my Tech/Comm bucket heat: $1,125 / $1,500 cap (75%)

    When I click "Check Heat for AAPL Trade"
    Then the backend calculates:
      - New portfolio heat: $2,890 + $747 = $3,637 (90.93%)
      - Portfolio cap: $4,000 (4% of equity)
      - Remaining capacity: $363 (9.07%)
      - Result: ✓ WITHIN CAP

      - New bucket heat: $1,125 + $747 = $1,872 (124.8%)
      - Bucket cap: $1,500 (1.5% of equity)
      - Overage: $372 (24.8%)
      - Result: ✗ EXCEEDS BUCKET CAP

    Then I see a RED warning:
      "⚠️ BUCKET CAP EXCEEDED
       Tech/Comm bucket would be $1,872 (124.8% of $1,500 cap)
       Overage: $372

       Options:
       1. Reduce position size to max 96 shares (risk = $451)
       2. Close an existing Tech/Comm position first
       3. Choose a different ticker from another sector"

    And the "Proceed to Trade Entry" button is DISABLED
    And the banner remains GREEN but shows: "⚠️ Heat cap issue - see Heat Check"

  Scenario: Adjust position size to respect heat caps
    Given the bucket cap is exceeded by $372
    When I return to Position Sizing screen
    And I click "Calculate Max Shares for Heat Cap"

    Then the backend calculates:
      - Remaining bucket capacity: $1,500 - $1,125 = $375
      - Max shares: floor(375 ÷ 4.70) = 79 shares
      - Actual risk: 79 × 4.70 = $371.30 (≤ $375 ✓)

    And the form updates to show: 79 shares per unit
    And I see: "Position reduced to respect heat cap (79 shares instead of 159)"

    When I click "Save Adjusted Position Plan"
    Then the plan is updated with reduced size

    When I return to Heat Check and click "Recheck"
    Then I see:
      - New bucket heat: $1,125 + $371 = $1,496 (99.73%)
      - Result: ✓ WITHIN CAP

    And the "Proceed to Trade Entry" button becomes ENABLED

  Scenario: Final trade entry with 5-gate validation
    Given my checklist is GREEN
    And my position sizing is calculated (79 shares)
    And my heat check passed
    And 2 minutes have elapsed since checklist evaluation

    When I navigate to Trade Entry screen
    Then I see a summary of my trade plan:
      - Ticker: AAPL
      - Direction: LONG
      - Entry: $180.50
      - Shares: 79 per unit (max 4 units)
      - Initial stop: $175.80 (2×N)
      - Risk per unit: $371.30
      - Add-on levels: [182.68, 184.85, 187.03]
      - Sector: Tech/Comm
      - Exit plan: 10-bar Donchian OR 2×N stop, whichever closer

    When I click "Run Final Gate Check"
    Then the backend validates ALL 5 gates:

      Gate 1: Banner Status
        - Current: GREEN ✓
        - All required checks confirmed ✓
        - Quality score ≥ threshold ✓
        - Result: PASS ✓

      Gate 2: Impulse Brake (2-minute cool-off)
        - Last evaluation: 2025-10-29 09:32:15
        - Current time: 2025-10-29 09:34:27
        - Elapsed: 2 min 12 sec ✓
        - Result: PASS ✓

      Gate 3: Cooldown Status
        - Ticker AAPL on cooldown? NO ✓
        - Bucket Tech/Comm on cooldown? NO ✓
        - Result: PASS ✓

      Gate 4: Heat Caps
        - Portfolio heat: 90.93% of cap ✓
        - Bucket heat: 99.73% of cap ✓
        - Result: PASS ✓

      Gate 5: Sizing Completed
        - Position plan saved? YES ✓
        - Risk calculated? YES ($371.30) ✓
        - Result: PASS ✓

    And all gates show GREEN ✓
    And the "SAVE GO DECISION" button is ENABLED (green, large)
    And the "SAVE NO-GO DECISION" button is also visible (for journaling rejections)

    When I click "SAVE GO DECISION"
    Then the backend saves a decision record with:
      - Timestamp
      - Ticker, sector, entry, stop, shares
      - All gate results (all PASS)
      - Banner status (GREEN)
      - Quality score (4/4)
      - Journal note
      - Risk amount

    And I see confirmation: "✓ GO decision saved for AAPL"
    And the ticker appears in my "Ready to Execute" list
    And I can now manually execute the trade in my broker

    # System does NOT execute trades - this is intentional
    # Trader must manually enter in broker (final human verification)

  Scenario: Logging a NO-GO decision
    Given I complete analysis of ticker XYZ
    But the heat check shows portfolio at 95% of cap

    When I navigate to Trade Entry screen
    And I click "Run Final Gate Check"
    Then Gate 4 (Heat Caps) shows FAIL ✗
    And the "SAVE GO DECISION" button is DISABLED

    When I click "SAVE NO-GO DECISION"
    And I enter reason: "Portfolio heat at 95%, no capacity for new position"

    Then the backend saves a NO-GO decision record
    And I see: "✓ NO-GO decision logged for XYZ"
    And this is journaled for future review
    And the ticker is removed from candidates

  Scenario: Position appears in calendar view
    Given I saved a GO decision for AAPL (Tech/Comm)
    When I navigate to Calendar view
    Then I see a 10-week grid (2 weeks back + 8 weeks forward)
    And each row is a sector bucket
    And each column is a week (Mon-Sun date range)

    And I see AAPL appear in:
      - Tech/Comm row
      - Week of 2025-10-27 column (entry week)
      - Plus subsequent weeks until expected exit (~60 days if options)

    And I can see sector diversification at a glance
    And I can identify if too many trades in one sector/week
    And I can spot gaps (weeks with no trades)
```

### Additional Workflows

**Portfolio Management:**
- View all open positions with current heat
- Update stops as positions move (trailing 10-bar Donchian)
- Close positions (log exit decision)
- Add units (if favorable move ≥ 0.5N from last add)

**Settings Management:**
- Update equity (monthly or after significant P&L change)
- Adjust risk% per unit (0.5-1.0% typical range)
- Configure heat caps (portfolio 4%, bucket 1.5% defaults)
- Manage sector bucket definitions

**FINVIZ Screener Presets:**
- Create custom FINVIZ screener URLs
- Save as named presets (e.g., "TF_BREAKOUT_LONG", "TF_PULLBACK_LONG")
- One-click scan using saved preset

---

## Technical Components

### Backend Functions (Direct Calls - No HTTP)

The Fyne GUI calls backend domain logic directly (in-process, no HTTP layer needed):

**Settings:**
```go
storage.GetAllSettings(db) -> map[string]string
storage.UpdateSetting(db, key, value) -> error
```

**Candidates:**
```go
scrape.ScrapeFinviz(presetURL) -> []string (tickers)
storage.AddCandidates(db, tickers, date) -> error
storage.GetCandidates(db, date) -> []Candidate
storage.DeleteCandidate(db, ticker) -> error
```

**Checklist:**
```go
domain.EvaluateChecklist(params) -> ChecklistResult
  // Returns: banner, missing_count, score, etc.
storage.SaveEvaluation(db, ticker, result) -> error
storage.GetLastEvaluation(db, ticker) -> Evaluation
```

**Position Sizing:**
```go
domain.CalculateSizeStock(ticker, entry, atr, k, riskPct, equity) -> SizeResult
domain.CalculateSizeOptDeltaATR(...) -> SizeResult
domain.CalculateSizeOptContracts(...) -> SizeResult
  // Returns: shares, risk, stops, add-ons
```

**Heat:**
```go
domain.CheckHeat(db, riskAmount, bucket, equity, portfolioCap, bucketCap) -> HeatResult
  // Returns: current heat, new heat, caps, pass/fail
```

**Gates:**
```go
domain.CheckGates(db, ticker, banner, riskDollars, bucket, equity, caps) -> GateResults
  // Returns: all 5 gate results + final pass/fail
```

**Decisions:**
```go
storage.SaveDecision(db, decision) -> error
storage.GetDecisions(db, date) -> []Decision
storage.GetDecision(db, id) -> Decision
```

**Positions:**
```go
storage.GetPositions(db) -> []Position
storage.AddPosition(db, position) -> error
storage.UpdatePosition(db, id, updates) -> error
storage.ClosePosition(db, id, exitPrice, exitDate) -> error
```

**Calendar:**
```go
storage.GetPositionsForCalendar(db, weeks) -> map[string][]Position
  // Returns positions grouped by sector × week
```

### Frontend Screens (Fyne Implementation)

**Layout:**
- `main.go` - Main application with window and tabs
- `theme.go` - Custom Fyne theme (British Racing Green + dark/light modes)
- `buildMainUI()` - Creates tabbed interface with navigation
- Large banner displayed at top of relevant screens (RED/YELLOW/GREEN)

**Dashboard:**
- `dashboard.go` - Main landing screen (buildDashboardScreen)
- Position table using widget.NewTable
- Heat gauges using progress bars
- Candidates summary card

**FINVIZ Scan:**
- `scanner.go` - Scan trigger and results (buildScannerScreen)
- Button to run FINVIZ scrape
- Results table with checkbox selection
- Import button

**Checklist:**
- `checklist.go` - Main checklist form (buildChecklistScreen)
- 5 required checkboxes (widget.NewCheck)
- 4 optional quality checkboxes
- Journal text area (widget.NewMultiLineEntry)
- Banner component showing RED/YELLOW/GREEN state

**Position Sizing:**
- `position_sizing.go` - Sizing calculator (buildPositionSizingScreen)
- Form with entry fields
- Method dropdown (Stock, Options delta-ATR, Options contracts)
- Results display in card
- Add-on schedule table

**Heat Check:**
- `heat_check.go` - Heat analysis (buildHeatCheckScreen)
- Portfolio heat progress bar
- Bucket heat progress bars
- Warning dialogs if caps exceeded
- Suggestions for resolution

**Trade Entry:**
- `trade_entry.go` - Final gate check (buildTradeEntryScreen)
- Trade summary card
- 5 gate results display
- Large GO/NO-GO buttons

**Calendar:**
- `calendar.go` - 10-week sector grid (buildCalendarScreen)
- Custom grid layout with containers
- Cells colored by heat level
- Tooltips on hover

**TradingView Integration:**
- `utils.go` - Helper function to open browser
- Constructs URL and opens in default browser
- Available from multiple screens (candidates, checklist, etc.)

**Utility Helpers:**
- `widgets.go` - Custom widgets and helpers
- Dialog helpers (ShowStyledInformation, ShowStyledConfirm)
- Banner creation function
- Progress bar builders

### State Management (Fyne)

**State Structs (`internal/gui/state/`):**

```go
// state.go - Central app state
type AppState struct {
    DB *storage.DB

    // Settings
    Settings map[string]string

    // Candidates
    Candidates []domain.Candidate
    CandidatesDate string

    // Current Checklist Evaluation
    CurrentChecklist struct {
        Ticker        string
        Banner        string // "RED", "YELLOW", "GREEN"
        RequiredGates domain.ChecklistGates
        QualityItems  domain.ChecklistQuality
        Score         int
        EvaluatedAt   time.Time
    }

    // Positions and Heat
    Positions      []domain.Position
    PortfolioHeat  float64
    BucketHeat     map[string]float64

    // UI State
    Banner struct {
        State   string // "RED", "YELLOW", "GREEN"
        Message string
    }
    Notifications []Notification
    IsLoading     bool

    // Callbacks for UI updates
    OnSettingsChanged    func()
    OnCandidatesChanged  func()
    OnChecklistChanged   func()
    OnPositionsChanged   func()
    OnBannerChanged      func()
}

// notification.go
type Notification struct {
    Type    string // "success", "warning", "error", "info"
    Message string
    Time    time.Time
}

// Methods on AppState to update and notify
func (s *AppState) SetBanner(state, message string) {
    s.Banner.State = state
    s.Banner.Message = message
    if s.OnBannerChanged != nil {
        s.OnBannerChanged()
    }
}

func (s *AppState) LoadSettings() error {
    settings, err := storage.GetAllSettings(s.DB)
    if err != nil {
        return err
    }
    s.Settings = settings
    if s.OnSettingsChanged != nil {
        s.OnSettingsChanged()
    }
    return nil
}

func (s *AppState) LoadPositions() error {
    positions, err := storage.GetAllPositions(s.DB)
    if err != nil {
        return err
    }
    s.Positions = positions

    // Calculate portfolio heat
    s.PortfolioHeat = 0
    s.BucketHeat = make(map[string]float64)
    for _, pos := range positions {
        s.PortfolioHeat += pos.Risk
        s.BucketHeat[pos.Bucket] += pos.Risk
    }

    if s.OnPositionsChanged != nil {
        s.OnPositionsChanged()
    }
    return nil
}
```

**Direct Backend Calls (No API Layer):**

```go
// In screen code, directly call backend functions

// Load settings
settings, err := storage.GetAllSettings(app.DB)
if err != nil {
    dialog.ShowError(err, window)
    return
}

// Update settings
err = storage.SaveSetting(app.DB, "equity", "100000")
if err != nil {
    dialog.ShowError(err, window)
    return
}

// Scan candidates
candidates, err := domain.ScanCandidates(preset)
if err != nil {
    dialog.ShowError(err, window)
    return
}

// Import candidates
err = storage.ImportCandidates(app.DB, tickers, date)
if err != nil {
    dialog.ShowError(err, window)
    return
}

// Evaluate checklist
result, err := domain.EvaluateChecklist(domain.ChecklistParams{
    FromPreset: true,
    Trend:      true,
    Liquidity:  true,
    Timeframe:  true,
    Earnings:   true,
    Journal:    journalText,
})
if err != nil {
    dialog.ShowError(err, window)
    return
}

// Update app state
app.CurrentChecklist.Banner = result.Banner
app.CurrentChecklist.Score = result.Score
app.SetBanner(result.Banner, result.Message)
```

---

## Implementation Strategy

### Phase 0: Foundation & Proof-of-Concept (Week 1-2)

**Goals:**
- Validate Go + Fyne architecture
- Prove direct backend function calls work from GUI
- Test cross-compilation to Windows .exe
- Establish build and deployment workflow

**Deliverables:**
1. **Fyne Basic App**
   - Simple Fyne app that calls backend functions directly
   - Display settings from database
   - Edit and save settings
   - Duration: 2-3 days

2. **Custom Theme Implementation**
   - British Racing Green accent color (#00352B)
   - Day mode and dark mode color schemes
   - Custom fonts and spacing
   - Duration: 1-2 days

3. **Build Pipeline**
   - `scripts/build-windows.sh` (cross-compile .exe from Linux)
   - Windows .exe with embedded database support
   - Test: Run .exe on Windows, GUI launches natively
   - Duration: 1-2 days

**Success Criteria:**
- ✓ Fyne app displays real data from SQLite database
- ✓ Direct function calls to backend domain logic work
- ✓ Custom British Racing Green theme applied
- ✓ Windows .exe runs natively (no browser required)
- ✓ Database operations work correctly
- ✓ Documentation updated (LLM-update.md, PROGRESS.md)

---

### Phase 1: Dashboard & FINVIZ Scanner (Week 3-4)

**Goals:**
- Build core navigation and layout using Fyne
- Implement FINVIZ scanning workflow
- Display candidates and positions

**Deliverables:**
1. **Layout Components**
   - Main window with tab navigation (widget.NewAppTabs)
   - Header with app name and session info
   - Tab-based navigation (Dashboard, Checklist, Sizing, Heat, Entry, Calendar)

2. **Dashboard Screen**
   - Display open positions (direct call: storage.GetAllPositions)
   - Show portfolio heat (calculate from positions)
   - Show today's candidates count
   - Quick stats: equity, available capacity
   - Uses widget.NewTable for positions list

3. **FINVIZ Scanner**
   - One-click "Run Daily Scan" button
   - Directly calls domain.ScanCandidates(preset)
   - Displays results in widget.NewTable
   - Checkbox selection for import (widget.NewCheck per row)
   - "Import Selected" button calls storage.ImportCandidates

4. **No Backend API Layer Needed**
   - All screens call backend functions directly
   - No HTTP server, no JSON marshalling
   - Direct access to domain and storage packages

**Success Criteria:**
- ✓ Dashboard displays real data from SQLite database
- ✓ FINVIZ scan button fetches candidates via domain.ScanCandidates
- ✓ Candidate import saves to database via storage.ImportCandidates
- ✓ Tab navigation works between all screens

---

### Phase 2: Checklist & Position Sizing (Week 5-6)

**Goals:**
- Implement 5 gates + quality items checklist using Fyne
- Display 3-state banner (custom widget)
- Calculate position sizing
- Enforce 2-minute cool-off

**Deliverables:**
1. **Checklist Screen**
   - Form: ticker entry, ATR entry, sector dropdown (widget.NewForm)
   - 5 required checkboxes (widget.NewCheck)
   - 4 optional checkboxes (widget.NewCheck)
   - Journal text area (widget.NewMultiLineEntry)
   - Large banner at top (custom container.New with colored background)
   - "Save Evaluation" button (widget.NewButton)

2. **Banner Widget**
   - Custom Fyne widget, minimum 20% screen height
   - RED: Any required unchecked (color.NRGBA{R: 220, G: 38, B: 38, A: 255})
   - YELLOW: All required, score < 3 (color.NRGBA{R: 245, G: 158, B: 11, A: 255})
   - GREEN: All required, score ≥ 3 (color.NRGBA{R: 16, G: 185, B: 129, A: 255})
   - Updates live as checkboxes change (via OnChanged callbacks)

3. **Position Sizing Screen**
   - Form: pre-filled from checklist evaluation
   - Method dropdown (widget.NewSelect with stock/opt-delta/opt-contracts)
   - "Calculate" button calls domain.CalculateSizeStock/Opt directly
   - Results display: shares, risk, stops, add-ons
   - "Save Position Plan" button

4. **Direct Backend Calls (No API Layer)**
   - domain.EvaluateChecklist(params) for checklist
   - domain.CalculateSizeStock/Opt for position sizing
   - storage.SaveEvaluation for timestamp (2-min timer)

**Success Criteria:**
- ✓ Banner changes color based on checklist state in real-time
- ✓ Checklist evaluation saves to database via storage layer
- ✓ 2-minute timer starts upon evaluation save
- ✓ Position sizing calculates correctly (matches Van Tharp method)
- ✓ Results match existing CLI output exactly

---

### Phase 3: Heat Check & Trade Entry (Week 7-8)

**Goals:**
- Implement heat cap validation using Fyne
- Build 5-gate final check
- Save GO/NO-GO decisions

**Deliverables:**
1. **Heat Check Screen**
   - Display current portfolio heat (widget.NewLabel)
   - Display current bucket heat (widget.NewTable)
   - "Check Heat" button calls domain.CheckHeat directly
   - RED warning if cap exceeded (custom container with red background)
   - Suggestions to resolve (widget.NewList)

2. **Trade Entry Screen**
   - Trade summary display (widget.NewCard with all details)
   - "Run Final Gate Check" button calls domain.CheckGates
   - Display all 5 gate results (widget.NewAccordion with pass/fail icons)
   - Large "SAVE GO DECISION" button (disabled until all pass, widget.NewButton)
   - "SAVE NO-GO DECISION" button (always enabled, widget.NewButton)

3. **Gate Validation**
   - Gate 1: Banner GREEN (check app state)
   - Gate 2: 2-min elapsed (check time.Since)
   - Gate 3: Not on cooldown (storage.CheckCooldown)
   - Gate 4: Heat caps OK (domain.CheckHeat)
   - Gate 5: Sizing done (check app state)

4. **Direct Backend Calls (No API Layer)**
   - domain.CheckHeat(db, risk, bucket, equity, caps)
   - domain.CheckGates(db, ticker, banner, risk, bucket, equity, caps)
   - storage.SaveDecision(db, decisionData)

**Success Criteria:**
- ✓ Heat check prevents cap violations
- ✓ All 5 gates must pass for GO decision
- ✓ 2-minute timer is enforced (no bypassing)
- ✓ NO-GO decisions are logged to database
- ✓ No backdoors exist (GO button disabled if any gate fails)

---

### Phase 4: Calendar & TradingView Integration (Week 9-10)

**Goals:**
- Visualize sector diversification using Fyne
- Integrate with TradingView charts
- Polish UI/UX

**Deliverables:**
1. **Calendar View**
   - 10-week grid (2 back + 8 forward) using custom Fyne canvas
   - Rows: sector buckets (widget.NewLabel for headers)
   - Columns: weeks (Mon-Sun)
   - Cells: tickers active in that sector/week (custom canvas objects)
   - Color coding: heat levels using custom color gradients

2. **TradingView Integration**
   - "Open in TradingView" buttons on candidates (widget.NewButton)
   - Construct URL with ticker symbol
   - Open in default browser using fyne.CurrentApp().OpenURL()
   - Documentation for applying Ed-Seykota.pine script

3. **UI Polish**
   - Keyboard shortcuts (window.Canvas().SetOnTypedKey)
   - Better loading states (widget.NewProgressBarInfinite)
   - Error handling improvements (dialog.ShowError)
   - VIM keybindings (already implemented in current Fyne GUI)

4. **Direct Backend Call (No API Layer)**
   - storage.GetCalendarData(db, weeks) returns position data
   - GUI builds calendar grid from raw data

**Success Criteria:**
- ✓ Calendar displays positions correctly
- ✓ TradingView links work
- ✓ UI feels fast and responsive
- ✓ No obvious bugs or UX issues

---

### Phase 5: Testing & Packaging (Week 11-12)

**Goals:**
- Comprehensive testing
- Windows installer
- User documentation

**Deliverables:**
1. **Testing**
   - Go unit tests for GUI components
   - Backend domain logic tests (already exist)
   - End-to-end workflow testing (manual)
   - Test all 5 gates enforcement
   - Test heat cap edge cases

2. **Windows Packaging**
   - Use `fyne package` to create .exe with icon
   - Bundle with database initialization
   - Desktop shortcut creation (fyne package -os windows -icon icon.png)
   - Optional: MSI installer using WiX or NSIS
   - Uninstaller (standard Windows uninstall)

3. **Documentation**
   - User guide (screenshots of Fyne GUI)
   - Quick start (5-minute setup)
   - Trading workflow guide
   - Troubleshooting
   - FAQ

**Success Criteria:**
- ✓ All critical paths tested
- ✓ Windows .exe runs on clean PC without dependencies
- ✓ User can complete full workflow in < 10 minutes
- ✓ Documentation is clear and complete

---

## Proof-of-Concept Approach

### Why Fyne for Production?

**Rationale:**
1. **Fast validation** - Prove Go GUI integration works (1-2 days)
2. **Desktop-native feel** - Pure desktop application, no browser required
3. **Single binary** - True single-file deployment (.exe)
4. **No HTTP layer** - Direct function calls, simpler architecture

**Fyne POC Scope (Minimal):**
```
1. Create simple Fyne window
2. Display: "TF-Engine Settings"
3. Show equity, risk%, caps from database
4. Button: "Refresh" (re-query database)
5. Button: "Update" (save changes)
6. Status: "Saved successfully" or error
```

**Fyne POC Success (COMPLETED):**
- ✓ Window opens
- ✓ Data loads from SQLite
- ✓ Updates save to SQLite
- ✓ Compiles to Windows .exe
- ✓ Runs on Windows without dependencies
- ✓ Custom British Racing Green theme implemented
- ✓ VIM keybindings implemented
- ✓ Dark mode toggle working

**Fyne Chosen for Production:**

**Why Fyne over Svelte?**
1. **Simpler Architecture** - No HTTP server, no JSON marshalling, direct function calls
2. **True Desktop App** - Native window, OS integration, no browser dependencies
3. **Single Binary** - One .exe file with everything embedded
4. **Go Ecosystem** - Same language as backend, better type safety
5. **Performance** - Direct memory access, no HTTP overhead
6. **Security** - No localhost port exposure, no CORS issues
7. **Offline First** - Works without network, no web stack vulnerabilities
8. **Proven** - Current GUI already working with all core features

**Current Fyne Implementation:**
- ✓ All 6 tabs implemented (Dashboard, Checklist, Sizing, Heat, Entry, Calendar)
- ✓ Custom theme with British Racing Green accent
- ✓ VIM mode with Vimium-style hints
- ✓ Dark mode toggle
- ✓ Info dialogs with explanations
- ✓ TradingView integration
- ✓ Working with SQLite database
- ✓ Browser opens to app when .exe runs
- ✓ Cross-compiles to Windows .exe

**Decision Made: Fyne**

Fyne was chosen for production because:
- Desktop-native feel is critical for trading application
- Simpler architecture (no HTTP layer, no JSON)
- Proven working implementation with all core features
- Single .exe deployment
- Better performance (direct function calls)
- Same language as backend (Go)

---

## Behavioral Specifications (Gherkin)

### Feature: Anti-Impulsivity Enforcement

```gherkin
Feature: System prevents impulsive trading
  As the trading system
  I must make it impossible to trade impulsively
  So that only systematic, trend-following trades are executed

  Background:
    Given the trader has analyzed a potential trade
    And the backend tf-engine is operational

  Scenario: Cannot bypass checklist evaluation
    When the trader tries to go directly to Trade Entry
    And the checklist has not been evaluated
    Then the system shows: "Complete checklist first"
    And the Trade Entry screen is READ-ONLY
    And the "SAVE GO DECISION" button is DISABLED

  Scenario: Cannot skip 2-minute cool-off
    Given the checklist was evaluated at 09:30:00
    And the current time is 09:31:45 (1 min 45 sec later)
    When the trader tries to save a GO decision
    Then the system shows: "Cool-off period: 15 seconds remaining"
    And the "SAVE GO DECISION" button is DISABLED
    And the countdown timer displays: 00:15

    When the current time becomes 09:32:01 (2 min 1 sec later)
    Then the countdown timer shows: 00:00
    And the "SAVE GO DECISION" button becomes ENABLED

  Scenario: Cannot exceed portfolio heat cap
    Given the portfolio heat is $3,800
    And the portfolio cap is $4,000 (4% of $100k)
    And the trader proposes a trade with $300 risk
    When the heat check is performed
    Then the new portfolio heat would be $4,100
    And this exceeds the cap by $100
    And the system shows RED warning: "Portfolio cap exceeded by $100"
    And the "SAVE GO DECISION" button is DISABLED
    And the system suggests: "Reduce position size or close existing position"

  Scenario: Cannot exceed bucket heat cap
    Given the Tech/Comm bucket heat is $1,400
    And the bucket cap is $1,500 (1.5% of $100k)
    And the trader proposes an AAPL trade (Tech/Comm) with $200 risk
    When the heat check is performed
    Then the new bucket heat would be $1,600
    And this exceeds the cap by $100
    And the system shows RED warning: "Tech/Comm bucket cap exceeded by $100"
    And the "SAVE GO DECISION" button is DISABLED

  Scenario: Cannot trade ticker on cooldown
    Given ticker AAPL is on cooldown until 2025-11-05
    And today is 2025-10-29
    When the trader selects AAPL for analysis
    Then the system shows: "⚠️ AAPL on cooldown until 2025-11-05 (7 days)"
    And the checklist form is DISABLED
    And the system suggests: "Choose a different ticker"

  Scenario: Cannot save GO decision with RED banner
    Given the checklist evaluation resulted in RED banner
    Because the "Signal" gate is unchecked
    When the trader navigates to Trade Entry
    Then Gate 1 (Banner Status) shows: FAIL ✗
    And the overall gate check is: FAIL
    And the "SAVE GO DECISION" button is DISABLED
    And only "SAVE NO-GO DECISION" is enabled

  Scenario: Cannot save GO decision with incomplete position sizing
    Given the checklist is GREEN
    And the 2-minute timer has elapsed
    And heat caps are OK
    And no cooldowns
    But position sizing has NOT been calculated
    When the trader tries to save a GO decision
    Then Gate 5 (Sizing Completed) shows: FAIL ✗
    And the "SAVE GO DECISION" button is DISABLED
    And the system shows: "Complete position sizing first"
```

### Feature: Systematic Workflow Enforcement

```gherkin
Feature: Daily trading workflow
  As a systematic trader
  I must follow the complete workflow
  So that I don't skip critical analysis steps

  Scenario: Morning scan workflow
    Given I open the app at market open
    When I navigate to Dashboard
    Then I see the "Run Daily FINVIZ Scan" button prominently

    When I click "Run Daily FINVIZ Scan"
    Then the backend scrapes my TF_BREAKOUT_LONG preset
    And returns candidates: ["AAPL", "MSFT", "NVDA", ...]
    And displays: "Found 23 candidates"

    When I click "Review Candidates"
    Then I see a table with columns: [Ticker, Sector, Last Close, Volume]
    And I can check/uncheck tickers
    And tickers on cooldown are grayed out

    When I select 12 tickers and click "Import Selected"
    Then those 12 tickers are saved to candidates table
    And they appear in "Today's Candidates" list

  Scenario: Chart analysis integration
    Given I have imported candidates
    When I click "Open in TradingView" next to AAPL
    Then a new browser tab opens to:
      https://www.tradingview.com/chart/?symbol=AAPL
    And I manually apply my Ed-Seykota.pine script (if not default)
    And I verify the 55-bar Donchian breakout visually
    And I note the current N (ATR) value: 2.35
    And I note the entry price: 180.50

  Scenario: Complete trade evaluation
    Given I have analyzed AAPL in TradingView
    When I fill out the checklist with:
      | Field         | Value      |
      | Ticker        | AAPL       |
      | Entry         | 180.50     |
      | ATR (N)       | 2.35       |
      | Sector        | Tech/Comm  |
    And I check all 5 required gates
    And I check 4 quality items
    And I enter journal note
    Then the banner turns GREEN

    When I save the evaluation
    Then the 2-minute timer starts
    And I see: "Cool-off: 2:00 remaining"

    When I calculate position sizing
    Then I get: 159 shares per unit, $747 risk

    When I check heat
    And all caps are OK
    And 2 minutes have elapsed
    Then all 5 gates pass
    And "SAVE GO DECISION" is enabled

    When I save the GO decision
    Then it's logged to the database
    And appears in "Ready to Execute" list
```

### Feature: Banner State Transitions

```gherkin
Feature: Banner reflects trade readiness
  As a visual discipline tool
  The banner must clearly show whether I can trade
  So that I never miss a failed gate

  Scenario: Initial state is RED
    Given I open a new checklist
    Then the banner shows: "🔴 DO NOT TRADE 🔴"
    And the banner is red background
    And the message says: "Complete required checklist items"

  Scenario: Transition to YELLOW
    Given the banner is RED
    When I check all 5 required gates
    But I check only 1 quality item (score = 1)
    And the quality threshold is 3
    Then the banner shows: "🟡 CAUTION 🟡"
    And the banner is yellow background
    And the message says: "Quality score below threshold (1/3)"

  Scenario: Transition to GREEN
    Given the banner is YELLOW
    When I check 2 more quality items (score now 3)
    Then the banner shows: "🟢 OK TO TRADE 🟢"
    And the banner is green background
    And the message says: "All gates pass, quality score met"

  Scenario: Regression from GREEN to YELLOW
    Given the banner is GREEN (all required, score 3/3)
    When I uncheck 1 quality item (score now 2/3)
    Then the banner immediately shows: "🟡 CAUTION 🟡"

  Scenario: Regression from YELLOW to RED
    Given the banner is YELLOW
    When I uncheck 1 required gate
    Then the banner immediately shows: "🔴 DO NOT TRADE 🔴"
```

### Feature: TradingView Integration

```gherkin
Feature: Seamless TradingView chart access
  As a trader
  I want to quickly open candidates in TradingView
  So that I can verify signals with my Pine script

  Scenario: Open ticker in TradingView
    Given I am viewing my candidates list
    When I click "Open in TV" for ticker AAPL
    Then a new browser tab opens
    And the URL is: https://www.tradingview.com/chart/?symbol=AAPL
    And I can manually load my Ed-Seykota.pine script

  Scenario: Verify Donchian breakout in TradingView
    Given I have AAPL chart open with Ed-Seykota.pine script
    And the script shows entryLen = 55, exitLen = 10
    When I look at the current bar
    Then I see if close > donHi[1] (for long breakout)
    And I note the N value (ATR) displayed by the script
    And I note the plotted stop level (entry - 2×N)

  Scenario: Return to app with parameters
    Given I verified the breakout in TradingView
    And N = 2.35, entry = 180.50
    When I return to the TF-Engine GUI
    And I enter those values in the checklist
    Then the system pre-calculates:
      - Stop: 180.50 - (2.0 × 2.35) = 175.80
      - Add1: 180.50 + (0.5 × 2.35) = 181.68
      - Add2: 180.50 + (1.0 × 2.35) = 182.85
      - Add3: 180.50 + (1.5 × 2.35) = 184.03
```

---

## Success Criteria

### Functional Requirements ✅

**Must Have:**
- [ ] FINVIZ scanner imports candidates via one click
- [ ] Checklist enforces all 5 required gates
- [ ] Banner visually indicates RED/YELLOW/GREEN state with **gradients** (large, obvious, beautiful)
- [ ] **Day/night mode toggle** with smooth transition (button in header)
- [ ] Theme preference persists to localStorage
- [ ] All UI elements use **gradients** (no flat colors)
- [ ] Smooth animations on all state changes (0.3s ease-in-out)
- [ ] 2-minute cool-off timer is enforced (cannot bypass)
- [ ] Position sizing calculates using Van Tharp method (matches backend CLI)
- [ ] Heat caps prevent trades exceeding 4% portfolio / 1.5% bucket
- [ ] All 5 gates must pass before GO decision saves
- [ ] NO-GO decisions are logged for journaling
- [ ] Calendar shows 10-week sector × week grid
- [ ] TradingView links open charts for analysis
- [ ] Settings persist in SQLite database
- [ ] Positions, decisions, candidates, cooldowns all persist
- [ ] **Comprehensive logging** throughout application
  - Backend: All API calls, database operations, calculations logged
  - Frontend: User actions, API calls, component lifecycle, errors logged
  - Log files accessible for debugging
  - Feature usage tracked to identify useful vs problematic features
  - Performance metrics captured for optimization

**Should Have:**
- [ ] Keyboard shortcuts for common actions (with visual hints)
- [ ] Loading states with animated gradient spinners
- [ ] Error messages with helpful suggestions and gradient accents
- [ ] Success confirmations with slide-in animations
- [ ] Auto-refresh after data changes
- [ ] Tooltips with smooth fade-in on hover
- [ ] Breadcrumb navigation with gradient dividers

**Nice to Have:**
- [ ] Mobile responsive (for viewing, not trading)
- [ ] Export decisions to CSV with styled download button
- [ ] Trade log search and filter with animated results
- [ ] Performance charts (P&L over time) with gradient fill
- [ ] Browser notifications (optional)
- [ ] Custom gradient presets (user can choose accent colors)
- [ ] Confetti animation on successful GO decision save 🎉

### Non-Functional Requirements ✅

**Performance:**
- [ ] Backend API responses < 100ms for calculations
- [ ] Database queries < 50ms
- [ ] FINVIZ scan < 5 seconds (network dependent)
- [ ] UI feels instant (no perceived lag)
- [ ] Page load < 2 seconds

**Reliability:**
- [ ] No crashes during normal operation
- [ ] Database corruption protection (SQLite WAL mode)
- [ ] Graceful error handling (API failures)
- [ ] Input validation prevents bad data

**Security:**
- [ ] Runs locally only (no external API calls except FINVIZ)
- [ ] No sensitive data in logs
- [ ] Database file permissions restricted (user only)
- [ ] No SQL injection vulnerabilities

**Usability:**
- [ ] Intuitive navigation (can complete workflow without docs)
- [ ] Large click targets (especially GO decision button)
- [ ] Clear visual hierarchy
- [ ] Consistent terminology (matches trading literature)
- [ ] Helpful error messages with suggestions

**Maintainability:**
- [ ] Code is well-commented
- [ ] API endpoints follow REST conventions
- [ ] Components are modular and reusable
- [ ] Documentation is up to date
- [ ] Build process is automated

**Deployability:**
- [ ] Single .exe for Windows (includes frontend + backend)
- [ ] No installation required (just run .exe)
- [ ] Opens browser automatically OR embedded webview
- [ ] Database created on first run
- [ ] Minimal dependencies (Go + SQLite only)

---

## Risk Management

### Technical Risks

**Risk: Fyne GUI complexity for custom widgets**
- **Likelihood:** Medium
- **Impact:** Medium
- **Mitigation:** Use standard Fyne widgets where possible; build custom widgets only when needed
- **Contingency:** Simplify UI design to use standard widgets only
- **Status:** RESOLVED - Custom widgets (banner, calendar) successfully implemented

**Risk: Cross-compilation to Windows fails**
- **Likelihood:** Low (Go is mature)
- **Impact:** Critical
- **Mitigation:** Test early and often; use pure Go (no cgo)
- **Contingency:** Develop on Windows directly (abandon Linux-first workflow)
- **Status:** RESOLVED - Cross-compilation working perfectly (Linux → Windows .exe)

**Risk: FINVIZ changes page structure (scraper breaks)**
- **Likelihood:** Medium (over time)
- **Impact:** Medium (can still use manual import)
- **Mitigation:** Document scraper logic; make selectors configurable
- **Contingency:** Fallback to manual CSV import; consider FINVIZ API (paid)

**Risk: SQLite modernc.org/sqlite driver issues on Windows**
- **Likelihood:** Low
- **Impact:** Medium
- **Mitigation:** Use modernc.org/sqlite (pure Go, no cgo required)
- **Contingency:** Switch to mattn/go-sqlite3 with cgo enabled
- **Status:** RESOLVED - modernc.org/sqlite working on Windows

### Process Risks

**Risk: Scope creep (adding non-essential features)**
- **Likelihood:** High
- **Impact:** Medium (delays delivery)
- **Mitigation:** Strict adherence to anti-impulsivity principles; reject flexibility
- **Contingency:** Timebox phases; cut nice-to-haves

**Risk: Linux-first workflow causes Windows issues**
- **Likelihood:** Medium
- **Impact:** Medium (deployment friction)
- **Mitigation:** Test on Windows early and often; automate builds
- **Contingency:** Switch to Windows development if too painful

**Risk: Documentation falls behind code**
- **Likelihood:** High
- **Impact:** Low (but violates RULES.md)
- **Mitigation:** Daily Loop mandates updating LLM-update.md and PROGRESS.md
- **Contingency:** Dedicated doc sprint before each phase completion

### User Risks

**Risk: User finds workaround to bypass gates**
- **Likelihood:** Low (gates enforced in backend)
- **Impact:** Critical (defeats purpose)
- **Mitigation:** Backend validates all gates; UI cannot override
- **Contingency:** Log all decision attempts; review for patterns

**Risk: User doesn't understand the workflow**
- **Likelihood:** Medium
- **Impact:** Medium (frustration, abandonment)
- **Mitigation:** Clear documentation; in-app tooltips; guided tutorial
- **Contingency:** Video walkthrough; one-on-one training

**Risk: User trades manually outside system**
- **Likelihood:** High (cannot prevent)
- **Impact:** Low (system still useful for analysis)
- **Mitigation:** Make system so easy that manual tracking is more work
- **Contingency:** Accept this; system is a tool, not a jail

---

## Next Steps

### Immediate Actions (This Week)

1. **Review and Approve This Plan**
   - Read thoroughly
   - Ask questions
   - Identify gaps
   - Confirm priorities

2. **Create Step-by-Step Sub-Plans**
   - `plans/phase0-foundation-poc.md` - Fyne POC and setup
   - `plans/phase1-dashboard-finviz.md` - Dashboard and scanning
   - `plans/phase2-checklist-sizing.md` - Checklist and sizing
   - `plans/phase3-heat-gates.md` - Heat check and trade entry
   - `plans/phase4-calendar-polish.md` - Calendar and final polish
   - `plans/phase5-testing-packaging.md` - Testing and Windows installer

3. **Set Up GUI Project Structure (COMPLETED)**
   - ✓ Created `internal/gui/` directory structure
   - ✓ Implemented all 6 main tabs
   - ✓ Custom British Racing Green theme
   - ✓ VIM mode with Vimium-style hints
   - ✓ Dark mode toggle
   - ✓ Working with SQLite database

4. **Fyne Implementation (COMPLETED)**
   - ✓ Created Fyne desktop app
   - ✓ Integrated with backend domain logic (direct calls)
   - ✓ Database access working
   - ✓ Cross-compilation to Windows working
   - ✓ Custom widgets implemented

### Week 1-2: Foundation (COMPLETED)

- ✓ Fyne implementation complete
- ✓ Custom theme implemented (British Racing Green)
- ✓ Build pipeline working (cross-compile Linux → Windows)
- ✓ Decision: Fyne chosen for production
- ✓ All 6 tabs implemented and working

### Week 3-4: Dashboard & Scanner

- [ ] Layout and routing (2 days)
- [ ] Dashboard screen (3 days)
- [ ] FINVIZ scanner (3 days)
- [ ] Backend API endpoints (2 days)

### Week 5-6: Checklist & Sizing

- [ ] Checklist form (3 days)
- [ ] Banner component (2 days)
- [ ] Position sizing (3 days)
- [ ] 2-minute timer (2 days)

### Week 7-8: Heat & Gates

- [ ] Heat check (3 days)
- [ ] Trade entry screen (3 days)
- [ ] 5-gate validation (2 days)
- [ ] Decision saving (2 days)

### Week 9-10: Calendar & Polish

- [ ] Calendar view (3 days)
- [ ] TradingView integration (1 day)
- [ ] UI polish (3 days)
- [ ] Bug fixes (3 days)

### Week 11-12: Testing & Packaging

- [ ] Comprehensive testing (5 days)
- [ ] Windows installer (3 days)
- [ ] Documentation (2 days)

---

## Appendices

### A. Ed-Seykota.pine Script Parameters

From `reference/Ed-Seykota.pine`:

**Core Parameters:**
- `entryLen = 55` - Donchian entry lookback (System-2)
- `exitLen = 10` - Donchian exit lookback
- `nLen = 20` - N = ATR length
- `stopN = 2.0` - Initial stop in N multiples
- `addStepN = 0.5` - Add every X × N
- `maxUnits = 4` - Max units (including initial)
- `riskPct = 1.0` - Risk % of equity per unit

**Optional Filters:**
- `useMarket = false` - Market regime filter (SPY > 200 SMA)
- `timeExitMode = "None"` - Time exit (None/Close/Roll)
- `minVol = 0` - Minimum 20-bar avg volume

**Calculation:**
```pine
N = ta.atr(nLen)  // ATR(20)
donHi = ta.highest(high, entryLen)  // 55-bar high
donLo = ta.lowest(low, entryLen)    // 55-bar low
stopL = strategy.position_avg_price - stopN * N_entry  // Long stop
stopS = strategy.position_avg_price + stopN * N_entry  // Short stop
```

**Signal:**
```pine
longBreak = allowLong and liqOK and longRegOK and (close > donHiPrev)
shortBreak = allowShort and liqOK and shortRegOK and (close < donLoPrev)
```

**Position Sizing:**
```pine
sharesForUnit(_equity, _Nentry) =>
    riskDollars = _equity * (riskPct/100.0)
    perShareRisk = math.max(stopN * _Nentry, syminfo.mintick)
    math.max(1, math.floor(riskDollars / perShareRisk))
```

### B. Backend API Reference (Existing Domain Logic)

**Position Sizing:**
- `domain.CalculateSizeStock(ticker, entry, atr, k, riskPct, equity)`
- `domain.CalculateSizeOptDeltaATR(ticker, entry, atr, delta, k, riskPct, equity)`
- `domain.CalculateSizeOptContracts(ticker, entry, atr, contracts, debit, k, riskPct, equity)`

**Checklist:**
- `domain.EvaluateChecklist(fromPreset, trend, liquidity, timeframe, earnings, journal, ...)`
- Returns: `{ banner: "GREEN"|"YELLOW"|"RED", missing_count, missing_items, score, ... }`

**Heat:**
- `domain.CheckHeat(db, riskAmount, bucket, equity, portfolioCap, bucketCap)`
- Returns: `{ current_portfolio_heat, new_portfolio_heat, portfolio_cap_exceeded, bucket_cap_exceeded, ... }`

**Gates:**
- `domain.CheckGates(db, ticker, banner, riskDollars, bucket, equity, portfolioCap, bucketCap)`
- Returns: `{ gate1_pass, gate2_pass, gate3_pass, gate4_pass, gate5_pass, all_gates_pass }`

**Database:**
- `storage.NewDB(path)` - Open/create SQLite database
- `storage.GetSettings()`, `storage.UpdateSettings()`
- `storage.GetPositions()`, `storage.AddPosition()`, `storage.UpdatePosition()`, `storage.ClosePosition()`
- `storage.GetCandidates(date)`, `storage.AddCandidates(tickers, date)`
- `storage.SaveDecision()`, `storage.GetDecisions()`
- `storage.GetCooldowns()`, `storage.AddCooldown()`, `storage.RemoveCooldown()`

**FINVIZ Scraper:**
- `scrape.ScrapeFinviz(url)` - Returns list of tickers
- Presets: `TF_BREAKOUT_LONG`, `TF_BREAKOUT_SHORT`, etc.

### C. Reference Documents

**Must Read:**
1. `docs/project/WHY.md` - Core philosophy (why this system exists)
2. `docs/anti-impulsivity.md` - Design principles
3. `docs/dev/DEVELOPMENT_PHILOSOPHY.md` - How we build
4. `docs/dev/CLAUDE_RULES.md` - Development standards
5. `CLAUDE.md` - Guidance for future AI sessions
6. `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md` - Operating rules (Linux-first, no Git in Linux)

**Reference:**
- `reference/Ed-Seykota.pine` - TradingView Pine Script
- `FRESH_START_PLAN.md` - Original GUI plan (superseded by this)
- `PROJECT_HISTORY.md` - Why we abandoned Excel/VBA
- `docs/PROJECT_STATUS.md` - M24 completion status

**Ongoing:**
- `docs/LLM-update.md` - Session-by-session log (always current)
- `docs/PROGRESS.md` - Narrative progress and decisions

---

**End of Overview Plan**

**Status:** 📋 Ready for Review
**Next:** Approve plan → Create phase-specific sub-plans → Begin Fyne POC
**Updated:** 2025-10-29
