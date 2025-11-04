# TF-ENGINE 2.0 - ARCHITECTURAL OVERVIEW

## Go + Fyne Trading Application for Windows

**Version:** 2.0
 **Target Platform:** Windows (Fyne is cross-platform but optimized for Windows)
 **Based On:** 293 backtests (21 securities × 14 strategies)
 **Date:** November 2025

------

## 🚨 ARCHITECTURAL RULES (MUST FOLLOW)

### Rule #1: No Feature Creep

**DO NOT create new features unless:**

- Explicitly requested by the user/architect, OR
- Approved in writing before implementation

**Rationale:** Previous version (TF-Engine 1.0) shows feature bloat without core functionality working. Focus on the workflow defined below.

### Rule #2: Data-Driven Design via Policy File

**ALL recommendations must be based on the 293 backtest results:**

- Healthcare: 92.31% strategy success → Priority #1
- Technology: 71.43% strategy success → Priority #2
- Utilities: 0% strategy success → BLOCKED
- Energy: 21.43% strategy success → WARNING

**Load sector/strategy policy from a signed, read-only `data/policy.v1.json`.** Users cannot edit in the UI; only the architect ships updates. On policy hash mismatch, run in safe mode (Healthcare/Tech allowed, Utilities blocked, Energy warned) and log the error.

**Rationale:** This aligns with research-driven updates without recompiling. The policy file is versioned, contains provenance metadata, and can be updated as new backtest data emerges. The application validates the policy file integrity at startup using SHA256 checksums.

### Rule #3: Anti-Impulsivity is Core

**The application MUST enforce:**

- 2-minute cooldown timer (cannot be skipped)
- 5 required checklist gates (all must pass)
- 3 optional checklist items (improve score)
- Heat check before allowing trade entry

**These are not optional features. They are requirements.**

### Rule #4: FINVIZ Integration Only

**DO NOT build a stock screener:**

- Use FINVIZ URLs with pre-configured parameters
- Open in system browser with chart mode (&v=211)
- Import ticker symbols manually (user types them)

**Rationale:** FINVIZ is proven, tested, and works. Don't reinvent it.

### Rule #5: Minimal Viable Product First

**Phase 1 (MVP) includes ONLY:**

1. Sector selection screen
2. Screener launch (opens browser)
3. Ticker + strategy selection
4. Checklist validation
5. Position sizing calculator
6. Heat check screen
7. Trade entry form
8. Calendar view (read-only)

**Phase 2 (Later) includes:**

- Trade editing/deletion
- Sample data mode
- Vimium keyboard nav
- Advanced analytics

### Rule #6: Save Early, Save Often

**Auto-save trade data:**

- After EACH screen completion
- To local JSON file (trades.json)
- Before ANY navigation
- On application exit

**Never lose user progress. Ever.**

### Rule #7: Styling is Functional, Not Fancy

**UI Requirements:**

- Day mode: Light green background (#E8F5E9), dark text (#1B5E20)
- Night mode: British Racing Green (#004225), light text (#E8F5E9)
- Large, readable fonts (minimum 14pt)
- High contrast (accessibility first)
- No animations, gradients, or "modern" effects

**Rationale:** User is color-aware (specified green themes), needs readability.

### Rule #8: Workflow is Linear (With Escape Hatches)

**Navigation:**

- FORWARD button advances to next screen (after validation)
- BACK button returns to previous screen (preserves data)
- CANCEL button returns to dashboard (confirms first)
- NO SKIP buttons (workflow must be followed)

**Exception:** Can jump to Calendar view from anywhere (read-only mode).

------

## 📐 APPLICATION ARCHITECTURE

### Tech Stack

**Language:** Go 1.21+
 **GUI Framework:** Fyne v2.4+
 **Data Storage:** JSON files (local filesystem)
 **Browser Integration:** `exec.Command` to open FINVIZ URLs
 **State Management:** Single global `AppState` struct

**Why Fyne?**

- Native feel on Windows
- Cross-platform (if needed later)
- Material Design by default (clean, modern)
- Good documentation
- Active community

**Why NOT other frameworks:**

- ❌ Electron (too heavy, web tech)
- ❌ Qt bindings (C++ complexity)
- ❌ GTK (Linux-focused)
- ❌ Native Windows API (not cross-platform)

------

### Directory Structure

```
tf-engine-2.0/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
│
├── cmd/
│   └── tf-engine/
│       └── main.go              # CLI entry if needed
│
├── internal/
│   ├── app/
│   │   ├── app.go               # Main application controller
│   │   ├── state.go             # Global state management
│   │   └── config.go            # Configuration loader
│   │
│   ├── data/
│   │   ├── models.go            # Data structures (Trade, Sector, etc.)
│   │   ├── storage.go           # JSON file persistence
│   │   ├── validation.go        # Business logic validation
│   │   └── defaults.go          # Default data (sectors, strategies)
│   │
│   ├── ui/
│   │   ├── theme.go             # Day/Night mode themes
│   │   ├── navigation.go        # Screen navigation controller
│   │   └── widgets/             # Custom Fyne widgets
│   │       ├── cooldown_timer.go
│   │       ├── heat_gauge.go
│   │       └── calendar_grid.go
│   │
│   ├── screens/
│   │   ├── dashboard.go         # Main dashboard (Image 1)
│   │   ├── sector_select.go     # Screen 1: Choose sector
│   │   ├── screener_launch.go   # Screen 2: Open FINVIZ
│   │   ├── ticker_strategy.go   # Screen 3: Enter ticker + strategy
│   │   ├── checklist.go         # Screen 4: Anti-impulsivity gates
│   │   ├── position_sizing.go   # Screen 5: Calculate position size
│   │   ├── heat_check.go        # Screen 6: Portfolio heat validation
│   │   ├── trade_entry.go       # Screen 7: Options strategy selection
│   │   ├── calendar.go          # Screen 8: Timeline grid view
│   │   └── _post_mvp/           # Phase 2 screens (not in MVP)
│   │       ├── trade_management.go  # Edit/delete trades (Phase 2)
│   │       └── welcome.go           # Welcome screen (Phase 2)
│   │
│   ├── finviz/
│   │   ├── urls.go              # Pre-configured FINVIZ screener URLs
│   │   └── launcher.go          # Browser launch logic
│   │
│   ├── calculations/
│   │   ├── position_sizing.go   # Poker-style betting calculator
│   │   ├── heat_calculator.go   # Portfolio heat aggregation
│   │   └── risk_metrics.go      # Risk/reward calculations
│   │
│   └── vimium/
│       ├── keybindings.go       # Keyboard navigation mode
│       └── command_parser.go    # Vimium-style command parsing
│
├── assets/
│   ├── welcome.txt              # Welcome screen text
│   ├── help.txt                 # Help documentation
│   └── sample_trades.json       # Sample data for testing
│
├── data/
│   ├── trades.json              # User's trade history (gitignored)
│   ├── settings.json            # User preferences (gitignored)
│   ├── policy.v1.json           # Sector rules, strategy whitelist, FINVIZ URLs (read-only, versioned)
│   ├── policy.v1.json.sha256    # Integrity checksum for policy validation
│   └── backtest_results.json    # Your 293 backtest data (read-only, legacy reference)
│
├── docs/
│   ├── USER_GUIDE.md            # End-user documentation
│   ├── DEVELOPMENT.md           # Developer setup guide
│   └── WORKFLOW.md              # Detailed screen-by-screen flow
│
├── tests/
│   ├── unit/                    # Unit tests (models, calculations)
│   ├── integration/             # Integration tests (screen flows)
│   └── e2e/                     # End-to-end tests (full workflow)
│
└── scripts/
    ├── build.bat                # Windows build script
    ├── run_dev.bat              # Development mode launcher
    └── generate_sample_data.go  # Create sample trades
```

------

## 🎨 UI/UX DESIGN

### Theme System

#### Day Mode (Light)

```go
type DayTheme struct {
    Background:     color.RGBA{232, 245, 233, 255}  // #E8F5E9 (Light Green)
    Surface:        color.RGBA{255, 255, 255, 255}  // White
    Primary:        color.RGBA{67, 160, 71, 255}    // #43A047 (Medium Green)
    PrimaryVariant: color.RGBA{27, 94, 32, 255}     // #1B5E20 (Dark Green)
    Secondary:      color.RGBA{129, 199, 132, 255}  // #81C784 (Light Green)
    Error:          color.RGBA{198, 40, 40, 255}    // #C62828 (Red)
    Text:           color.RGBA{27, 94, 32, 255}     // #1B5E20 (Dark Text)
    TextSecondary:  color.RGBA{76, 175, 80, 255}    // #4CAF50 (Medium Text)
}
```

#### Night Mode (Dark)

```go
type NightTheme struct {
    Background:     color.RGBA{0, 66, 37, 255}      // #004225 (British Racing Green)
    Surface:        color.RGBA{0, 77, 44, 255}      // #004D2C (Slightly Lighter)
    Primary:        color.RGBA{129, 199, 132, 255}  // #81C784 (Light Green)
    PrimaryVariant: color.RGBA{165, 214, 167, 255}  // #A5D6A7 (Lighter Green)
    Secondary:      color.RGBA{200, 230, 201, 255}  // #C8E6C9 (Very Light Green)
    Error:          color.RGBA{239, 83, 80, 255}    // #EF5350 (Light Red)
    Text:           color.RGBA{232, 245, 233, 255}  // #E8F5E9 (Light Text)
    TextSecondary:  color.RGBA{165, 214, 167, 255}  // #A5D6A7 (Medium Text)
}
```

### Font Specifications

**Primary Font:** "Roboto" (bundled with Fyne)
 **Monospace Font:** "Roboto Mono" (for ticker symbols, numbers)

**Sizes:**

- Heading 1: 24pt, Bold
- Heading 2: 20pt, SemiBold
- Body: 16pt, Regular ← **EFFECTIVE MINIMUM (use for all main text)**
- Small: 14pt, Regular (captions, labels only)
- Caption: 12pt, Regular (metadata only, use sparingly)

**CRITICAL:** Body text should always be 16pt. The 14pt "minimum" in earlier draft was too conservative. All readable text (buttons, labels, inputs, table cells) should be 16pt or larger.

### Layout Principles

**Grid System:**

- 12-column grid
- 16px gutter between columns
- Responsive breakpoints (not critical for desktop-only)

**Spacing:**

- Small: 8px
- Medium: 16px
- Large: 24px
- XLarge: 32px

**Button Sizes:**

- Small: 120px × 36px
- Medium: 160px × 44px (default)
- Large: 200px × 52px (primary actions)

**Input Fields:**

- Height: 44px
- Padding: 12px horizontal, 8px vertical
- Border: 2px solid (primary color when focused)

------

## 📊 DATA MODELS

### Core Data Structures

```go
// Trade represents a single trade entry
type Trade struct {
    ID                string    `json:"id"`                  // UUID
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
    
    // Screen 1: Sector Selection
    Sector            string    `json:"sector"`              // "Healthcare", "Technology", etc.
    
    // Screen 2: Screener (no data stored, just launches FINVIZ)
    
    // Screen 3: Ticker + Strategy
    Ticker            string    `json:"ticker"`              // "UNH", "MSFT", etc.
    Strategy          string    `json:"strategy"`            // "Alt43", "Alt46", etc.
    EntrySignalTime   time.Time `json:"entry_signal_time"`   // When user clicked "next"
    CooldownComplete  bool      `json:"cooldown_complete"`   // True after 2 minutes
    
    // Screen 4: Checklist
    ChecklistPassed   bool      `json:"checklist_passed"`    // True if all required gates pass
    ChecklistScore    int       `json:"checklist_score"`     // 0-8 (5 required + 3 optional)
    ChecklistResults  map[string]bool `json:"checklist_results"` // Gate name → pass/fail
    
    // Screen 5: Position Sizing
    AccountEquity     float64   `json:"account_equity"`      // Total account value
    RiskPerTrade      float64   `json:"risk_per_trade"`      // Percentage (e.g., 0.0075 = 0.75%)
    PositionSize      int       `json:"position_size"`       // Number of contracts
    MaxLoss           float64   `json:"max_loss"`            // Dollar amount at risk
    
    // Screen 6: Heat Check
    PortfolioHeat     float64   `json:"portfolio_heat"`      // Current total heat
    BucketHeat        float64   `json:"bucket_heat"`         // Heat in this sector
    HeatCheckPassed   bool      `json:"heat_check_passed"`   // True if under limits
    
    // Screen 7: Trade Entry
    OptionsStrategy   string    `json:"options_strategy"`    // "Bull Call Spread", etc.
    StrikePrice       float64   `json:"strike_price"`        // For single-leg strategies
    Strike1           float64   `json:"strike1"`             // For multi-leg strategies
    Strike2           float64   `json:"strike2"`
    Strike3           float64   `json:"strike3"`
    Strike4           float64   `json:"strike4"`
    ExpirationDate    time.Time `json:"expiration_date"`     // Option expiration
    Premium           float64   `json:"premium"`             // Total premium paid/received
    
    // Screen 8: Calendar (display only, no new data)
    
    // Exit Information (filled in later)
    ExitDate          *time.Time `json:"exit_date,omitempty"`
    ExitPrice         *float64   `json:"exit_price,omitempty"`
    ProfitLoss        *float64   `json:"profit_loss,omitempty"`
    WinLoss           *string    `json:"win_loss,omitempty"`  // "Win", "Loss", "Scratch"
}

// Sector configuration (loaded from policy.v1.json)
type Sector struct {
    Name              string    `json:"name"`                // "Healthcare"
    SuccessRate       float64   `json:"success_rate"`        // 0.9231 (92.31%)
    AllowedStrategies []string  `json:"allowed_strategies"`  // ["Alt43", "Alt46", "Alt39", "Alt10"]
    ScreenerURLs      map[string]string `json:"screener_urls"` // "universe" → URL, "pullback" → URL, etc.
    HeatCapPercent    float64   `json:"heat_cap_percent"`    // 0.015 (1.5% bucket cap)
    Blocked           bool      `json:"blocked"`             // True for Utilities (0% success)
    Warning           bool      `json:"warning"`             // True for Energy (21% success)
    WarningMessage    string    `json:"warning_message"`     // Shown to user if warning == true
}

// Strategy configuration (loaded from policy.v1.json)
type Strategy struct {
    Name              string    `json:"name"`                // "Alt43"
    Description       string    `json:"description"`         // "Volatility-Adaptive Targets"
    BestSectors       []string  `json:"best_sectors"`        // ["Healthcare", "Technology"]
    SuccessRate       float64   `json:"success_rate"`        // 0.6190 (61.90%)
    AvgTradeCount     int       `json:"avg_trade_count"`     // 20-60 trades typical
    HoldPeriodWeeks   string    `json:"hold_period_weeks"`   // "3-12 weeks"
    HoldPeriodMin     int       `json:"hold_period_min_weeks"` // 3 (minimum weeks)
    HoldPeriodMax     int       `json:"hold_period_max_weeks"` // 12 (maximum weeks)
    OptionsSuitable   bool      `json:"options_suitable"`    // True if holds match options timeframes
    OptionsWarning    string    `json:"options_warning"`     // "" if suitable, warning text if not
    RecommendedDTE    string    `json:"recommended_dte"`     // "30-45 DTE" for options traders
    
    // Best performers with this strategy
    BestPerformers    []Performance `json:"best_performers"`
}

// Performance record from backtests
type Performance struct {
    Security          string    `json:"security"`            // "XLV", "UNH", etc.
    Return            float64   `json:"return"`              // 0.3479 (34.79%)
    WinRate           float64   `json:"win_rate"`            // 0.64 (64%)
    TradeCount        int       `json:"trade_count"`         // 85
    ProfitFactor      float64   `json:"profit_factor"`       // 2.791
}

// Application settings
type Settings struct {
    ThemeMode         string    `json:"theme_mode"`          // "day" or "night"
    AccountEquity     float64   `json:"account_equity"`      // $100,000 default
    RiskPerTrade      float64   `json:"risk_per_trade"`      // 0.0075 (0.75%) default
    PortfolioHeatCap  float64   `json:"portfolio_heat_cap"`  // 0.04 (4%) default - GLOBAL CAP
    BucketHeatCap     float64   `json:"bucket_heat_cap"`     // 0.015 (1.5%) default - PER-SECTOR CAP
    SectorBuckets     map[string]float64 `json:"sector_buckets"` // Per-sector overrides (optional)
    VimiumEnabled     bool      `json:"vimium_enabled"`      // false default
    SampleDataMode    bool      `json:"sample_data_mode"`    // false default
    ShowWelcomeScreen bool      `json:"show_welcome_screen"` // true default (Phase 2)
}

// Policy file structure (loaded from policy.v1.json - READ ONLY)
type Policy struct {
    Version           string    `json:"version"`             // "1.0"
    GeneratedAt       time.Time `json:"generated_at"`        // When this policy was created
    SourceNotes       string    `json:"source_notes"`        // "Based on 293 backtests (2010-2025)"
    SHA256            string    `json:"sha256"`              // Checksum for integrity validation
    
    // Sector rules (from your 293 backtests)
    Sectors           []Sector  `json:"sectors"`             // All tradable + blocked sectors
    
    // Strategy rules (from your 293 backtests)
    Strategies        []Strategy `json:"strategies"`          // All available Pine Scripts
    
    // FINVIZ screener catalog (from Master Screener Guide)
    Screeners         map[string]ScreenerConfig `json:"screeners"` // "healthcare_universe" → config
    
    // Position sizing rules (poker-style multipliers)
    PositionSizing    PositionSizingRules `json:"position_sizing"`
}

// FINVIZ Screener configuration
type ScreenerConfig struct {
    Name              string    `json:"name"`                // "Healthcare Universe"
    URL               string    `json:"url"`                 // Full FINVIZ URL with &v=211
    Description       string    `json:"description"`         // "30-50 healthcare stocks in uptrends"
    ExpectedResults   string    `json:"expected_results"`    // "30-50 stocks"
    UseCase           string    `json:"use_case"`            // "weekly" or "daily"
    ProvenWinners     []string  `json:"proven_winners"`      // ["UNH", "LLY", "ABBV"]
}

// Position sizing rules from checklist score (Revision #7)
type PositionSizingRules struct {
    MinimumScore      int       `json:"minimum_score"`       // 5 (minimum to proceed)
    Multipliers       map[int]float64 `json:"multipliers"`   // 5 → 0.5, 6 → 0.75, 7 → 1.0, 8 → 1.25
    MinimumContracts  int       `json:"minimum_contracts"`   // 1 (always trade at least 1)
    Description       string    `json:"description"`         // Explanation of poker-style sizing
}

// Checklist gate
type ChecklistGate struct {
    Name              string    `json:"name"`                // "From Preset (SIG_REQ)"
    Description       string    `json:"description"`         // "Stock from FINVIZ screener"
    Required          bool      `json:"required"`            // true for required gates
    Order             int       `json:"order"`               // Display order
    HelpText          string    `json:"help_text"`           // Detailed explanation
}

// Options strategy configuration
type OptionsStrategy struct {
    Name              string    `json:"name"`                // "Bull Call Spread"
    Category          string    `json:"category"`            // "Bullish", "Bearish", "Neutral"
    LegsCount         int       `json:"legs_count"`          // 2 for spreads, 1 for single
    RiskLevel         string    `json:"risk_level"`          // "Low", "Medium", "High"
    Description       string    `json:"description"`         // Brief explanation
    BestFor           []string  `json:"best_for"`            // ["Bullish trends", "Limited risk"]
}
```

------

## 🔄 APPLICATION WORKFLOW

### State Machine

```
States:
- DASHBOARD       → Initial state
- SECTOR_SELECT   → Choosing trading sector
- SCREENER_LAUNCH → Opening FINVIZ in browser
- TICKER_ENTRY    → Entering ticker + strategy (cooldown active)
- CHECKLIST       → Validating checklist gates
- POSITION_SIZING → Calculating position size
- HEAT_CHECK      → Validating portfolio heat
- TRADE_ENTRY     → Entering trade details
- CALENDAR_VIEW   → Viewing timeline (read-only)
- TRADE_MGMT      → Editing/deleting trades

Transitions:
DASHBOARD → SECTOR_SELECT (click "Start New Trade")
SECTOR_SELECT → SCREENER_LAUNCH (select sector)
SCREENER_LAUNCH → TICKER_ENTRY (user returns from browser)
TICKER_ENTRY → CHECKLIST (cooldown timer expires + user clicks next)
CHECKLIST → POSITION_SIZING (all required gates pass)
POSITION_SIZING → HEAT_CHECK (position calculated)
HEAT_CHECK → TRADE_ENTRY (heat check passes)
TRADE_ENTRY → CALENDAR_VIEW (trade saved)
CALENDAR_VIEW → DASHBOARD (view complete)

Any State → DASHBOARD (click "Cancel" + confirm)
Any State → CALENDAR_VIEW (click "Calendar" button)

**Resume-in-Place (Revision #11):**
The application persists the current trade's state and screen position. When user clicks "Resume Session", they re-enter at the EXACT step they left, with all data preserved.

```go
type TradeInProgress struct {
    Trade         Trade     `json:"trade"`           // Current trade data
    CurrentScreen string    `json:"current_screen"`  // "CHECKLIST", "POSITION_SIZING", etc.
    LastUpdated   time.Time `json:"last_updated"`    // When user last interacted
    CooldownStart *time.Time `json:"cooldown_start"` // If in cooldown, when it started
}

func (a *App) SaveProgress() {
    if a.State.CurrentTrade == nil {
        return
    }
    
    progress := TradeInProgress{
        Trade:         *a.State.CurrentTrade,
        CurrentScreen: a.State.CurrentScreen.String(),
        LastUpdated:   time.Now(),
        CooldownStart: a.State.CooldownStartTime,
    }
    
    // Save to trades_in_progress.json
    SaveTradeProgress("data/trades_in_progress.json", progress)
}

func (a *App) ResumeSession() {
    progress, err := LoadTradeProgress("data/trades_in_progress.json")
    if err != nil {
        dialog.ShowInformation("No Session", "No trade in progress to resume", a.Window)
        return
    }
    
    // Restore trade state
    a.State.CurrentTrade = &progress.Trade
    
    // Restore cooldown if applicable
    if progress.CooldownStart != nil {
        elapsed := time.Since(*progress.CooldownStart)
        if elapsed < 2*time.Minute {
            // Cooldown still active, resume timer
            remaining := 2*time.Minute - elapsed
            a.StartCooldownTimer(remaining)
        }
    }
    
    // Navigate to saved screen
    a.NavigateTo(ScreenFromString(progress.CurrentScreen))
}
```

### Screen-by-Screen Specifications

#### Screen 0: Dashboard (Image 1)

**Purpose:** Home base, shows current state and launch points

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Start New Trade] [Resume Session ▼]  [Dark Mode] [Help] [VIM]│
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Dashboard                                                  │
│                                                             │
│  Account Settings              Heat Status                  │
│  ┌────────────────────┐       ┌──────────────────────┐    │
│  │ Equity: $100,000   │       │ Portfolio Heat: $0   │    │
│  │ Risk: 0.75%        │       │ Cap: 4.0% of $100k   │    │
│  │ Heat Cap: 4.0%     │       │ [████████░░░░] 50%   │    │
│  │ Bucket Cap: 1.5%   │       │                      │    │
│  │ [Edit Settings]    │       └──────────────────────┘    │
│  └────────────────────┘                                    │
│                                                             │
│  Open Positions (Alt 31)       Today's Candidates          │
│  ┌────────────────────┐       ┌──────────────────────┐    │
│  │ No open positions  │       │ Date: 2025-11-03     │    │
│  │                    │       │ No candidates found  │    │
│  │                    │       │                      │    │
│  │                    │       │ [Refresh]            │    │
│  └────────────────────┘       └──────────────────────┘    │
│                                                             │
│  [Sidebar]                     [Calendar View]              │
│  - Scanner                                                  │
│  - Checklist                                                │
│  - Position Sizing                                          │
│  - Heat Check                                               │
│  - Trade Entry                                              │
│  - Session History                                          │
└─────────────────────────────────────────────────────────────┘
```

**Components:**

- Top toolbar: Quick actions (Start, Resume, Dark Mode, Help, VIM toggle)
- Account settings box (display + edit button)
- Heat status gauge (visual indicator)
- Open positions count (from trades.json, filtered by ExitDate == nil)
- Today's candidates (placeholder for future FINVIZ integration)
- Left sidebar (navigation to all screens)
- Calendar view button (launches screen 8)

**Actions:**

- "Start New Trade" → Navigate to SECTOR_SELECT
- "Resume Session ▼" → Dropdown of incomplete trades → Resume from last screen
- "Dark Mode" toggle → Switch theme
- "Help" → Show help popup
- "VIM" toggle → Enable/disable Vimium mode
- Sidebar items → Navigate to specific screen (if allowed by state)

------

#### Screen 1: Sector Selection

**Purpose:** Choose which market sector to trade based on backtest success rates

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back to Dashboard] [Cancel]                    Step 1 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Select Trading Sector                                      │
│                                                             │
│  Choose a sector based on your 293 backtest results:       │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ✓ Healthcare                               92.31%  │   │
│  │   Most successful sector (12/13 strategies)        │   │
│  │   Best: Alt43 on XLV (+34.79%), Alt46 (+34.80%)   │   │
│  │   [View Screener]                                  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ✓ Technology                               71.43%  │   │
│  │   Strong performer (10/14 strategies)              │   │
│  │   Best: Alt15 on MSFT (+52.65%), Alt22 on QQQ     │   │
│  │   [View Screener]                                  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ⚠ Consumer Discretionary                   64.29%  │   │
│  │   Moderate success (9/14 strategies)               │   │
│  │   Best: AMZN (+35.82%), WMT (+31.69%)             │   │
│  │   [View Screener]                                  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ❌ UTILITIES                                 0.00%  │   │
│  │   NEVER TRADE - Zero profitable strategies         │   │
│  │   DUK lost money in ALL 14 tests                   │   │
│  │   [BLOCKED]                                        │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ⚠ Energy                                   21.43%  │   │
│  │   WARNING - Rarely profitable                      │   │
│  │   Only 3/14 strategies worked                      │   │
│  │   [Use with Caution]                               │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  [Continue →]                                              │
└─────────────────────────────────────────────────────────────┘
```

**Data Source:**

- Load from `data/policy.v1.json` (primary source)
- Calculate sector-level success rates from policy
- Sort by success rate (highest first)
- **ENFORCE BLACKOUT ZONE #1:** Block Utilities (success_rate == 0, blocked == true)
- **ENFORCE BLACKOUT ZONE #1:** Warn on Energy (success_rate < 30%, warning == true)

**Validation:**

- User MUST select a sector (no default)
- **Cannot proceed if blocked sector selected** (button disabled + tooltip explains)
- **Confirm dialog if warning sector selected:** "Energy has only 21.43% success rate. Are you sure?"

**Actions:**

- Click sector card → Select (highlight in green)
- "View Screener" → Preview FINVIZ URL (doesn't navigate yet)
- "Continue" → Save sector to Trade object → Navigate to SCREENER_LAUNCH

------

#### Screen 2: Screener Launch

**Purpose:** Open FINVIZ screener in browser, then return to app

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 2 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Launch FINVIZ Screener                                     │
│                                                             │
│  Selected Sector: Healthcare (92.31% success)               │
│                                                             │
│  Choose your screening approach:                            │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ WEEKLY SCREENERS (Universe Building)                │   │
│  │                                                    │   │
│  │ [Launch Healthcare Universe]                       │   │
│  │ • 30-50 healthcare stocks in uptrends              │   │
│  │ • Loads from policy: screener "healthcare_universe"│   │
│  │ • Proven winners: UNH, LLY, ABBV, JNJ              │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ DAILY SCREENERS (Situational Setups)                │   │
│  │                                                    │   │
│  │ [Launch Healthcare Pullback]                       │   │
│  │ • 0-10 stocks pulling back in uptrends             │   │
│  │ • RSI < 40, Price < SMA50, still > SMA200          │   │
│  │ • [Fix: Loosen RSI to <45] (one-click)            │   │
│  │                                                    │   │
│  │ [Launch Healthcare Breakout]                       │   │
│  │ • 0-15 stocks at 52-week highs                     │   │
│  │ • New high signal, strong momentum                 │   │
│  │ • [Fix: Switch to 20-day high] (one-click)        │   │
│  │                                                    │   │
│  │ [Launch Golden Cross]                              │   │
│  │ • 10-30 stocks with SMA50 > SMA200 recently       │   │
│  │ • Early trend development                          │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  What to do in FINVIZ:                                      │
│  1. Review the chart thumbnails                             │
│  2. Click on stocks with clean uptrends                     │
│  3. Verify trend looks like XLV (steady climb)              │
│  4. Note the ticker symbol of your chosen stock             │
│  5. Return to this app and enter the ticker                 │
│                                                             │
│  After reviewing FINVIZ, click Continue below:             │
│  [I've reviewed the screener results →]                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Screener URL Loading:**

```go
func GetScreenerURL(sector string, screenerType string, policy Policy) (string, error) {
    // Find sector in policy
    var sectorConfig *Sector
    for _, s := range policy.Sectors {
        if s.Name == sector {
            sectorConfig = &s
            break
        }
    }
    if sectorConfig == nil {
        return "", fmt.Errorf("sector %s not found in policy", sector)
    }
    
    // Get URL from sector's screener_urls map
    url, exists := sectorConfig.ScreenerURLs[screenerType]
    if !exists {
        return "", fmt.Errorf("screener type %s not found for sector %s", screenerType, sector)
    }
    
    return url, nil
}

// One-click fix-ups (Revision #11)
func GetLoosenedURL(baseURL string, fix string) string {
    switch fix {
    case "loosen_rsi":
        // Replace ta_rsi_os40 with ta_rsi_os45
        return strings.Replace(baseURL, "ta_rsi_os40", "ta_rsi_os45", 1)
    case "switch_to_20d":
        // Replace ta_highlow52w_nh with ta_highlow20d_nh
        return strings.Replace(baseURL, "ta_highlow52w_nh", "ta_highlow20d_nh", 1)
    default:
        return baseURL
    }
}
```

**Actions:**

- "Launch FINVIZ" → Open browser with URL from sectors.json
- "I've reviewed" → Navigate to TICKER_ENTRY (can click multiple times)

**Implementation:**

```go
// Example FINVIZ URL opening
func LaunchScreener(sector string, url string) error {
    cmd := exec.Command("cmd", "/c", "start", url)
    return cmd.Run()
}
```

------

#### Screen 3: Ticker + Strategy Entry (with Cooldown)

**Purpose:** Enter ticker symbol, select strategy, enforce 2-minute cooldown

**Layout (Before Cooldown):**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 3 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Enter Trade Details                                        │
│                                                             │
│  Sector: Healthcare (92.31% success)                        │
│                                                             │
│  Ticker Symbol:                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ [UNH_______________]                               │   │
│  └────────────────────────────────────────────────────┘   │
│  Example: UNH, LLY, ABBV, JNJ                              │
│                                                             │
│  Select Strategy (Pine Script):                             │
│  Only showing strategies validated for Healthcare sector    │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ○ Alt43 - Volatility-Adaptive (92% Healthcare)    │   │
│  │   XLV: +34.79% ⭐ RECORD | UNH: +30.92%           │   │
│  │   Hold: 3-12 weeks | DTE: 30-45 | ✓ Options OK   │   │
│  │                                                    │   │
│  │ ○ Alt46 - Sector-Adaptive (92% Healthcare)        │   │
│  │   XLV: +34.80% | UNH: +32.16% ⭐ UNH BEST         │   │
│  │   Hold: 3-12 weeks | DTE: 30-45 | ✓ Options OK   │   │
│  │                                                    │   │
│  │ ○ Alt39 - Age-Based Targets (92% Healthcare)      │   │
│  │   XLV: +29.70% | UNH: +27.07%                     │   │
│  │   Hold: 3-12 weeks | DTE: 30-45 | ✓ Options OK   │   │
│  │                                                    │   │
│  │ ○ Alt10 - Profit Targets (76% Universal)          │   │
│  │   Works everywhere | UNH: +33.13%                 │   │
│  │   Hold: 3-10 weeks | DTE: 30-60 | ✓ Options OK   │   │
│  │                                                    │   │
│  │ ○ Alt45 - Dual-Momentum (67% Universal)           │   │
│  │   Hold: 3-10 weeks | DTE: 30-60 | ✓ Options OK   │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ℹ Alt15 is NOT shown (Tech specialist, not for Healthcare) │
│                                                             │
│  [Confirm Selection]                                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Layout (During Cooldown):**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 3 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ⏱ Cooldown Period Active                                   │
│                                                             │
│  Trade Details:                                             │
│  • Sector: Healthcare                                       │
│  • Ticker: UNH                                              │
│  • Strategy: Alt43 (Volatility-Adaptive)                    │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │                                                    │   │
│  │              Time Remaining: 1:23                  │   │
│  │                                                    │   │
│  │  [████████████████████░░░░░░░░░░░░░░░░░░]  70%    │   │
│  │                                                    │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Why the wait?                                              │
│  The 2-minute cooldown prevents impulsive "I gotta get     │
│  in NOW!" decisions. Research shows that waiting just 2     │
│  minutes dramatically reduces emotional trading errors.     │
│                                                             │
│  Use this time to:                                          │
│  • Review the FINVIZ chart again                            │
│  • Check if there are any earnings announcements            │
│  • Verify the trade still makes sense                       │
│  • Take a breath and stay calm                              │
│                                                             │
│  [Continue →]  (disabled until timer expires)              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Validation:**

- Ticker must be 1-5 uppercase letters
- Strategy must be selected
- **Strategy must be in sector's allowed_strategies list** (loaded from policy.v1.json)
- **Options suitability check:** If strategy.options_suitable == false, show warning
- **Hold period check:** If strategy.hold_period_max_weeks > 12, show warning for options traders
- Cooldown cannot be skipped (button disabled until timer expires)

**Strategy Filtering Logic:**

```go
func GetAllowedStrategies(sector Sector, policy Policy) []Strategy {
    var allowed []Strategy
    for _, strategyName := range sector.AllowedStrategies {
        for _, strategy := range policy.Strategies {
            if strategy.Name == strategyName {
                allowed = append(allowed, strategy)
                break
            }
        }
    }
    return allowed
}
```

**Implementation:**

```go
type CooldownTimer struct {
    StartTime     time.Time
    Duration      time.Duration // 2 minutes
    OnComplete    func()
    ticker        *time.Ticker
}

func (ct *CooldownTimer) Start() {
    ct.StartTime = time.Now()
    ct.ticker = time.NewTicker(1 * time.Second)
    go func() {
        for range ct.ticker.C {
            elapsed := time.Since(ct.StartTime)
            if elapsed >= ct.Duration {
                ct.ticker.Stop()
                ct.OnComplete()
                return
            }
            // Update UI with remaining time
        }
    }()
}
```

------

#### Screen 4: Checklist Validation (Images 2-6)

**Purpose:** Anti-impulsivity gates to prevent bad trades

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 4 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Pre-Trade Checklist                                        │
│                                                             │
│  Trade: UNH (Healthcare) using Alt43                        │
│                                                             │
│  Required Gates (All Must Pass)                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ ✓ From Preset (SIG_REQ)                           │   │
│  │   Stock came from FINVIZ screener [?]             │   │
│  │                                                    │   │
│  │ ✗ Trend Confirmed (RISK_REQ)                      │   │
│  │   Price broke above 55-day high [?]               │   │
│  │   [Mark as Pass] [Mark as Fail]                   │   │
│  │                                                    │   │
│  │ ? Liquidity OK (OPT_REQ)                          │   │
│  │   Weekly options available [?]                    │   │
│  │   [Mark as Pass] [Mark as Fail]                   │   │
│  │                                                    │   │
│  │ ? TV Confirm (EXIT_REQ)                           │   │
│  │   Exit plan confirmed in TradingView [?]          │   │
│  │   [Mark as Pass] [Mark as Fail]                   │   │
│  │                                                    │   │
│  │ ? Earnings OK (BEHAV_REQ)                         │   │
│  │   No earnings for 5 days [?]                      │   │
│  │   [Mark as Pass] [Mark as Fail]                   │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Optional Quality Items (Improve Score)                     │
│  ┌────────────────────────────────────────────────────┐   │
│  │ □ Regime OK (SPY > 200SMA)                        │   │
│  │ □ No Chase (< 2N above 20-EMA)                    │   │
│  │ □ Journal Entry Written                            │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Checklist Score: 3/8                                       │
│  Status: ❌ Cannot proceed (not all required gates pass)   │
│                                                             │
│  [Evaluate Checklist]                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Gate Definitions (from Images 3-6):**

1. **From Preset (SIG_REQ)** - REQUIRED
   - Auto-pass if ticker came from FINVIZ screener (we track this)
   - Manual fail if user typed ticker randomly
2. **Trend Confirmed (RISK_REQ)** - REQUIRED
   - Check: Price broke above 55-day high?
   - User clicks "Mark as Pass" or "Mark as Fail"
   - Help text: Shows Image 4 content
3. **Liquidity OK (OPT_REQ)** - REQUIRED
   - Check: Weekly options available?
   - User verifies on broker platform
   - Help text: "Open interest > 100 on ATM strikes"
4. **TV Confirm (EXIT_REQ)** - REQUIRED
   - Check: Exit plan confirmed in TradingView?
   - Help text: Shows Image 5 content (10-day low or stop loss)
5. **Earnings OK (BEHAV_REQ)** - REQUIRED
   - Check: No earnings in next 5 days?
   - Help text: Shows Image 6 content (2-minute rule)
6. **Regime OK** - OPTIONAL
   - Check: SPY > 200-day SMA?
   - Bonus point if true
7. **No Chase** - OPTIONAL
   - Check: Price < 2 ATR above 20-day EMA?
   - Bonus point if true
8. **Journal Entry Written** - OPTIONAL
   - Check: User wrote rationale?
   - Bonus point if true

**Validation:**

- All 5 required gates must pass
- Score = passed gates / 8 (0.375 minimum to proceed)
- If fail, show dialog: "Fix the failing gates or cancel trade"

**Actions:**

- "?" icon → Show detailed help popup for that gate
- "Mark as Pass/Fail" → Toggle gate state
- "Evaluate Checklist" → Validate all gates → Proceed or block

------

#### Screen 5: Position Sizing

**Purpose:** Calculate position size using poker-style betting

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 5 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Calculate Position Size                                    │
│                                                             │
│  Trade: UNH (Healthcare) using Alt43                        │
│  Checklist Score: 7/8 (Excellent)                          │
│                                                             │
│  Account Settings                                           │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Account Equity:    $ [100000___]                   │   │
│  │ Risk per Trade:      [0.75____] %                  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Trade-Specific Inputs                                      │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Entry Price:       $ [525.00__]                    │   │
│  │ Stop Loss:         $ [510.00__]                    │   │
│  │ Distance to Stop:  $ 15.00 (2.86%)                 │   │
│  │                                                    │   │
│  │ Option Premium:    $ [8.50____] per contract       │   │
│  │ Contracts Control:   100 shares per contract       │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Position Size Calculation                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Risk Amount:       $ 750 (0.75% of $100k)          │   │
│  │ Risk per Share:    $ 15.00 (entry - stop)          │   │
│  │ Share Quantity:    50 shares (750 / 15)            │   │
│  │                                                    │   │
│  │ For options:                                       │   │
│  │ Contracts:         1 contract (50 / 100)           │   │
│  │ Total Cost:        $ 850 (1 × $8.50 × 100)         │   │
│  │ Max Loss:          $ 850 (premium paid)            │   │
│  │                                                    │   │
│  │ ⚠ Actual risk ($850) > Target risk ($750)         │   │
│  │   Adjust by using wider stop or fractional sizing  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Poker-Style Bet Sizing (Based on Checklist Score)         │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Score 5/8 (Minimum): 0.5x sizing (half position)   │   │
│  │ Score 6/8 (Good):    0.75x sizing                   │   │
│  │ Score 7/8 (Great):   1.0x sizing ← YOUR SCORE      │   │
│  │ Score 8/8 (Perfect): 1.25x sizing                   │   │
│  │                                                    │   │
│  │ Recommended: 1 contract (full size)                 │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  [Calculate] [Accept & Continue →]                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Calculation Logic:**

```go
// Position sizing multipliers (from policy.v1.json)
// These are FIRST-CLASS business rules, not arbitrary
const (
    MinimumChecklistScore = 5  // Cannot proceed below this
)

var PositionMultipliers = map[int]float64{
    5: 0.5,   // Minimum score → Half position
    6: 0.75,  // Good score → Three-quarter position
    7: 1.0,   // Great score → Full position
    8: 1.25,  // Perfect score → Oversized position (125%)
}

func CalculatePositionSize(equity, riskPercent, entryPrice, stopPrice, optionPremium float64, checklistScore int) PositionSize {
    // Validate checklist score
    if checklistScore < MinimumChecklistScore {
        return PositionSize{Error: "Checklist score too low (minimum 5 required)"}
    }
    
    // Base risk amount
    riskAmount := equity * riskPercent
    
    // Risk per share
    riskPerShare := math.Abs(entryPrice - stopPrice)
    if riskPerShare == 0 {
        return PositionSize{Error: "Entry and stop cannot be equal"}
    }
    
    // Share quantity
    shareQty := int(riskAmount / riskPerShare)
    
    // Contracts (round down)
    contracts := shareQty / 100
    if contracts == 0 {
        contracts = 1 // Minimum 1 contract
    }
    
    // Actual cost for options
    totalCost := float64(contracts) * optionPremium * 100
    
    // Get poker-style multiplier from policy
    multiplier, exists := PositionMultipliers[checklistScore]
    if !exists {
        multiplier = 1.0 // Default to full position if score not in table
    }
    
    // Apply multiplier
    adjustedContracts := int(float64(contracts) * multiplier)
    if adjustedContracts < 1 {
        adjustedContracts = 1 // Minimum 1 contract always
    }
    
    adjustedCost := float64(adjustedContracts) * optionPremium * 100
    
    return PositionSize{
        Contracts:         adjustedContracts,
        TotalCost:         adjustedCost,
        MaxLoss:           adjustedCost, // For long options
        RiskAmount:        riskAmount,
        Multiplier:        multiplier,
        ChecklistScore:    checklistScore,
        BaseContracts:     contracts,
        AdjustmentApplied: fmt.Sprintf("Score %d → %.2fx multiplier = %d contracts", 
            checklistScore, multiplier, adjustedContracts),
    }
}
```

**Multiplier Table (Document for Testing):**

| Checklist Score | Multiplier | Position Size | Rationale                              |
| --------------- | ---------- | ------------- | -------------------------------------- |
| 5 (Minimum)     | 0.5x       | Half          | Barely passed gates, reduce risk       |
| 6 (Good)        | 0.75x      | Three-quarter | Solid setup, normal risk reduction     |
| 7 (Great)       | 1.0x       | Full          | Excellent setup, full position         |
| 8 (Perfect)     | 1.25x      | Oversized     | All gates + bonuses, increase exposure |

**This table is loaded from policy.v1.json and cannot be changed by users.**

**Validation:**

- All inputs must be positive numbers
- Stop loss must be below entry (for longs)
- Option premium must be realistic (> 0)

------

#### Screen 6: Heat Check

**Purpose:** Validate portfolio heat before allowing trade

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 6 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Portfolio Heat Check                                       │
│                                                             │
│  Current Heat Status                                        │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Account Equity:      $100,000                      │   │
│  │ Portfolio Heat Cap:  4.0%                          │   │
│  │ Maximum Heat:        $4,000                        │   │
│  │                                                    │   │
│  │ Current Heat:        $1,500                        │   │
│  │ Heat Percentage:     1.5%                          │   │
│  │                                                    │   │
│  │ [███████░░░░░░░░░░░░░░░░░░░░░] 37.5%             │   │
│  │                                                    │   │
│  │ Status: ✓ Under Cap (room for $2,500 more)       │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Sector-Level Heat (Bucket Caps)                           │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Healthcare:                                        │   │
│  │   Current: $500 (0.5%)                            │   │
│  │   Cap: $1,500 (1.5%)                              │   │
│  │   [████░░░░░░░░░░] 33% ✓                         │   │
│  │                                                    │   │
│  │ Technology:                                        │   │
│  │   Current: $750 (0.75%)                           │   │
│  │   Cap: $1,500 (1.5%)                              │   │
│  │   [█████░░░░░░░░] 50% ✓                          │   │
│  │                                                    │   │
│  │ Consumer Discretionary:                            │   │
│  │   Current: $250 (0.25%)                           │   │
│  │   Cap: $1,500 (1.5%)                              │   │
│  │   [██░░░░░░░░░░░] 17% ✓                          │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  New Trade Impact                                           │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Sector: Healthcare                                 │   │
│  │ Position Size: 1 contract                          │   │
│  │ Max Loss: $850                                     │   │
│  │                                                    │   │
│  │ After This Trade:                                  │   │
│  │   Portfolio Heat: $2,350 (2.35%) ✓                │   │
│  │   Healthcare Heat: $1,350 (1.35%) ✓               │   │
│  │                                                    │   │
│  │ Status: ✓ APPROVED - Within all limits            │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  [Continue to Trade Entry →]                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Heat Calculation:**

```go
func CalculatePortfolioHeat(trades []Trade) HeatStatus {
    totalHeat := 0.0
    sectorHeat := make(map[string]float64)
    
    for _, trade := range trades {
        if trade.ExitDate == nil { // Open position
            heat := trade.MaxLoss
            totalHeat += heat
            sectorHeat[trade.Sector] += heat
        }
    }
    
    return HeatStatus{
        TotalHeat:   totalHeat,
        SectorHeat:  sectorHeat,
    }
}

func ValidateHeat(current HeatStatus, newTrade Trade, settings Settings, policy Policy) (bool, string) {
    // ENFORCE BLACKOUT ZONE #2: Deny heat allocation to blocked sectors
    sector := policy.GetSector(newTrade.Sector)
    if sector.Blocked {
        return false, fmt.Sprintf("%s is BLOCKED - cannot allocate heat (0%% success rate)", newTrade.Sector)
    }
    
    // Check portfolio cap (GLOBAL)
    newTotal := current.TotalHeat + newTrade.MaxLoss
    portfolioCap := settings.AccountEquity * settings.PortfolioHeatCap
    if newTotal > portfolioCap {
        return false, fmt.Sprintf("Portfolio heat ($%.0f) would exceed GLOBAL cap ($%.0f). Current: $%.0f, New trade: $%.0f", 
            newTotal, portfolioCap, current.TotalHeat, newTrade.MaxLoss)
    }
    
    // Check sector bucket cap (PER-SECTOR)
    newSectorHeat := current.SectorHeat[newTrade.Sector] + newTrade.MaxLoss
    bucketCap := settings.AccountEquity * settings.BucketHeatCap
    
    // Use sector-specific override if exists
    if override, ok := settings.SectorBuckets[newTrade.Sector]; ok {
        bucketCap = settings.AccountEquity * override
    }
    
    if newSectorHeat > bucketCap {
        return false, fmt.Sprintf("%s heat ($%.0f) would exceed SECTOR bucket cap ($%.0f). Current: $%.0f, New trade: $%.0f", 
            newTrade.Sector, newSectorHeat, bucketCap, current.SectorHeat[newTrade.Sector], newTrade.MaxLoss)
    }
    
    return true, "APPROVED - Within all limits"
}
```

**Validation:**

- Must be under portfolio heat cap (GLOBAL: 4% default)
- Must be under sector bucket cap (PER-SECTOR: 1.5% default)
- **MUST NOT be blocked sector** (enforcement point #2)
- **If fail, BLOCK trade entry** → Cannot proceed to Screen 7
- Show error dialog with exact numbers and what to do (close positions or wait)

------

#### Screen 7: Trade Entry

**Purpose:** Select options strategy and enter trade details

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back] [Cancel]                                 Step 7 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Enter Trade Details                                        │
│                                                             │
│  Trade Summary                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Sector: Healthcare (92.31% success)                │   │
│  │ Ticker: UNH                                        │   │
│  │ Strategy: Alt43 (Volatility-Adaptive)             │   │
│  │ Position Size: 1 contract                          │   │
│  │ Max Loss: $850                                     │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Options Strategy                                           │
│  ┌────────────────────────────────────────────────────┐   │
│  │ [Dropdown: Select Options Strategy ▼]             │   │
│  │                                                    │   │
│  │ Bullish Strategies:                                │   │
│  │   • Bull Call Spread                               │   │
│  │   • Bull Put Credit Spread                         │   │
│  │   • Long Call                                      │   │
│  │   • Covered Call                                   │   │
│  │                                                    │   │
│  │ Bearish Strategies:                                │   │
│  │   • Bear Put Spread                                │   │
│  │   • Bear Call Credit Spread                        │   │
│  │   • Long Put                                       │   │
│  │   • Cash-Secured Put                               │   │
│  │                                                    │   │
│  │ Neutral Strategies:                                │   │
│  │   • Iron Condor                                    │   │
│  │   • Iron Butterfly                                 │   │
│  │   • Calendar Spread (Call/Put)                     │   │
│  │   • Diagonal Spread (Call/Put)                     │   │
│  │   • [+ Show All 24 Strategies]                     │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  Selected: Bull Call Spread                                 │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Description: Buy lower strike call, sell higher    │   │
│  │ Risk Level: Limited                                │   │
│  │ Max Profit: Difference in strikes - net premium    │   │
│  │ Max Loss: Net premium paid                         │   │
│  │                                                    │   │
│  │ Strike Prices:                                     │   │
│  │   Buy:  $ [520____] (lower strike)                │   │
│  │   Sell: $ [530____] (higher strike)               │   │
│  │   Spread Width: $10                                │   │
│  │                                                    │   │
│  │ Expiration: [2025-12-20_____] ▼                   │   │
│  │   DTE: 47 days                                     │   │
│  │   Recommended: 30-45 DTE for Alt43                 │   │
│  │                                                    │   │
│  │ Premiums:                                          │   │
│  │   Buy Call:  $ [12.50__]                          │   │
│  │   Sell Call: $ [8.20___]                          │   │
│  │   Net Debit: $ 4.30 per share                     │   │
│  │   Total Cost: $ 430 (1 contract)                  │   │
│  │                                                    │   │
│  │ Risk/Reward:                                       │   │
│  │   Max Loss: $430 (net debit)                      │   │
│  │   Max Profit: $570 (spread - debit)              │   │
│  │   Risk/Reward: 1:1.33                             │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  [Save Trade & View Calendar →]                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Dynamic Fields Based on Strategy:**

- Single-leg (Long Call/Put): 1 strike price
- Spreads: 2 strike prices
- Iron Condor/Butterfly: 4 strike prices
- Calendar/Diagonal: 2 strikes + 2 expirations

**Validation:**

- All required fields must be filled
- Strike prices must be logical (buy < sell for call spread)
- Expiration must be in future
- **DTE validation:** Calculate days to expiration and compare to strategy.recommended_dte
  - If outside recommended range, show warning: "Alt43 recommends 30-45 DTE, you selected 60 DTE"
  - Allow proceeding but mark with warning flag
- **Hold period validation:** Calculate weeks to expiration and compare to strategy.hold_period_min_weeks and hold_period_max_weeks
  - If outside range, show warning: "Alt15 has 15+ week holds, not ideal for options (time decay)"
- **ENFORCE BLACKOUT ZONE #3:** Before saving, verify sector.Blocked == false
  - If blocked, refuse save and show error: "Cannot trade Utilities (0% success rate)"
  - This is belt-and-suspenders (should never happen if Screen 1 and Screen 6 work)

**DTE Validation Logic:**

```go
func ValidateDTE(expirationDate time.Time, strategy Strategy) (bool, string) {
    daysToExpiration := int(time.Until(expirationDate).Hours() / 24)
    
    // Parse recommended DTE (e.g., "30-45 DTE")
    recommendedMin, recommendedMax := parseRecommendedDTE(strategy.RecommendedDTE)
    
    if daysToExpiration < recommendedMin || daysToExpiration > recommendedMax {
        return false, fmt.Sprintf("%s recommends %s, you selected %d DTE. Proceed with caution.", 
            strategy.Name, strategy.RecommendedDTE, daysToExpiration)
    }
    
    return true, "DTE within recommended range"
}
```

------

#### Screen 8: Calendar View

**Purpose:** Visual timeline of all trades across sectors

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Dashboard] [New Trade]                         Step 8 of 8 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Trade Calendar                                             │
│                                                             │
│  Timeline: Nov 20, 2024 ← → Dec 24, 2025                   │
│  [< Previous] [Today] [Next >]                             │
│                                                             │
│  Sectors  │ 2wk ago │ Last wk │ This wk │ Next wk │ ...   │
│  ─────────┼─────────┼─────────┼─────────┼─────────┼───── │
│  Health   │         │  [UNH]  │ [======]│ [======]│ ...   │
│           │         │  Alt43  │  $850   │  $850   │       │
│           │         │  47 DTE │  40 DTE │  33 DTE │       │
│  ─────────┼─────────┼─────────┼─────────┼─────────┼───── │
│  Tech     │ [MSFT]  │ [======]│ [======]│ [==]    │       │
│           │  Alt15  │  $1200  │  $1200  │  $1200  │       │
│           │  Closed │  14 DTE │   7 DTE │   EXIT  │       │
│  ─────────┼─────────┼─────────┼─────────┼─────────┼───── │
│  Consumer │         │         │ [AMZN]  │ [======]│ ...   │
│           │         │         │  Alt10  │  $650   │       │
│           │         │         │  28 DTE │  21 DTE │       │
│  ─────────┼─────────┼─────────┼─────────┼─────────┼───── │
│                                                             │
│  Legend:                                                    │
│  [====] Open position (bar length = time to expiration)    │
│  [TICK] New position (this week)                           │
│  Closed Position shown in gray                             │
│                                                             │
│  Click any bar for trade details                           │
│                                                             │
│  Heat Status:                                               │
│  Portfolio: $2,700 / $4,000 (67.5%)                        │
│  Healthcare: $850 / $1,500 (56.7%)                         │
│  Technology: $1,200 / $1,500 (80%) ⚠ Near limit          │
│  Consumer: $650 / $1,500 (43.3%)                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Implementation Details:**

- **X-axis: TIME (2 weeks past → 12 weeks future)** ← Horizontal axis
- **Y-axis: SECTORS** ← Vertical axis (Healthcare, Tech, Consumer, etc.)
- **Ticker symbols are BAR LABELS, not axes** (displayed on/inside the bars)
- Each trade = horizontal bar spanning from CreatedAt to ExpirationDate
- Bar color = sector color (green shades per sector)
- Bar length = days until expiration (proportional to time)
- Hover = tooltip with trade details (ticker, strategy, DTE, heat)
- Click = open trade details popup (full trade information)

**CRITICAL: Do NOT put tickers on the Y-axis. Sectors only on Y-axis.**

**Grid Rendering:**

```go
type CalendarGrid struct {
    Trades      []Trade
    StartDate   time.Time  // CLAMPED: today - 14 days
    EndDate     time.Time  // CLAMPED: today + 84 days (12 weeks)
    CellWidth   float32    // Pixels per week
    RowHeight   float32    // Pixels per sector row
}

func (cg *CalendarGrid) Initialize() {
    // Clamp render window to prevent excessive scrolling
    today := time.Now()
    cg.StartDate = today.AddDate(0, 0, -14)  // 2 weeks back
    cg.EndDate = today.AddDate(0, 0, 84)     // 12 weeks forward
}

func (cg *CalendarGrid) RenderTrade(trade Trade) {
    // Calculate bar position
    startX := cg.GetXPosition(trade.CreatedAt)
    endX := cg.GetXPosition(trade.ExpirationDate)
    y := cg.GetYPosition(trade.Sector)
    
    // Determine bar color (gray for closed trades)
    var barColor color.Color
    if trade.ExitDate != nil {
        barColor = color.RGBA{128, 128, 128, 255} // Gray for closed
    } else {
        barColor = GetSectorColor(trade.Sector)   // Sector color for open
    }
    
    // Draw bar
    bar := canvas.NewRectangle(barColor)
    bar.Move(fyne.NewPos(startX, y))
    bar.Resize(fyne.NewSize(endX - startX, cg.RowHeight))
    
    // Add ticker label (on the bar, not on Y-axis)
    label := widget.NewLabel(trade.Ticker)
    label.Move(fyne.NewPos(startX + 5, y + 5))
    
    // Add heat indicator (small text on bar)
    heatLabel := widget.NewLabel(fmt.Sprintf("$%.0f", trade.MaxLoss))
    heatLabel.Move(fyne.NewPos(startX + 5, y + 20))
}

// Render heat summary at bottom of calendar (Revision #4)
func (cg *CalendarGrid) RenderHeatSummary(trades []Trade, settings Settings) *fyne.Container {
    // Calculate current heat by sector
    sectorHeat := make(map[string]float64)
    totalHeat := 0.0
    
    for _, trade := range trades {
        if trade.ExitDate == nil { // Open positions only
            sectorHeat[trade.Sector] += trade.MaxLoss
            totalHeat += trade.MaxLoss
        }
    }
    
    // Build summary widget
    lines := []string{
        fmt.Sprintf("Portfolio Heat: $%.0f / $%.0f (%.1f%%)", 
            totalHeat, 
            settings.AccountEquity * settings.PortfolioHeatCap,
            (totalHeat / (settings.AccountEquity * settings.PortfolioHeatCap)) * 100),
    }
    
    for sector, heat := range sectorHeat {
        cap := settings.AccountEquity * settings.BucketHeatCap
        pct := (heat / cap) * 100
        status := "✓"
        if pct > 80 {
            status = "⚠ Near limit"
        }
        lines = append(lines, fmt.Sprintf("%s: $%.0f / $%.0f (%.1f%%) %s", 
            sector, heat, cap, pct, status))
    }
    
    return container.NewVBox(widget.NewLabel(strings.Join(lines, "\n")))
}
```

------

#### Screen 9: Trade Management

**Purpose:** Edit or delete existing trades (Phase 2 feature)

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│ [Back to Dashboard]                                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Trade Management                                           │
│                                                             │
│  All Trades                                                 │
│  [Filter: ▼ All] [Sort: ▼ Newest First]                   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ UNH (Healthcare) - Alt43                           │   │
│  │ Opened: 2025-11-03 | Expires: 2025-12-20         │   │
│  │ Status: ✓ Open | Heat: $850                       │   │
│  │ [Edit] [Close Trade] [Delete]                     │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │ MSFT (Technology) - Alt15                          │   │
│  │ Opened: 2025-10-27 | Closed: 2025-11-01          │   │
│  │ Status: ✓ Closed | P/L: +$425 (Win)              │   │
│  │ [View Details] [Delete]                            │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  [Export to CSV]                                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

------

### Additional Screens

#### Welcome Screen (on startup)

**Layout:**

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│                                                             │
│           🎯 TF-ENGINE 2.0                                  │
│           Trend-Following Trading System                     │
│                                                             │
│                                                             │
│  Built on 293 backtests across 21 securities × 14 strategies│
│                                                             │
│  Key Features:                                              │
│  • Sector-based screening (Healthcare: 92% success)         │
│  • Anti-impulsivity checklist                               │
│  • Position sizing calculator                               │
│  • Heat management                                          │
│  • Trade calendar visualization                             │
│                                                             │
│                                                             │
│  Quick Start:                                               │
│  1. Click "Start New Trade"                                 │
│  2. Select your sector (Healthcare recommended)             │
│  3. Follow the 8-step workflow                              │
│  4. Let the system enforce discipline                       │
│                                                             │
│                                                             │
│  [Start New Trade] [View Dashboard] [Help]                 │
│                                                             │
│  [ ] Don't show this again                                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

------

## 💾 DATA PERSISTENCE

### JSON File Structure

**policy.v1.json (READ ONLY - shipped with application):**

```json
{
  "version": "1.0",
  "generated_at": "2025-11-03T00:00:00Z",
  "source_notes": "Based on 293 backtests (21 securities × 14 strategies, 2010-2025). Healthcare: 92.31% success, Tech: 71.43%, Utilities: 0% (blocked).",
  "sha256": "calculated_at_runtime",
  "sectors": [
    {
      "name": "Healthcare",
      "success_rate": 0.9231,
      "allowed_strategies": ["Alt43", "Alt46", "Alt39", "Alt10", "Alt45"],
      "screener_urls": {
        "universe": "https://finviz.com/screener.ashx?v=211&...",
        "pullback": "https://finviz.com/screener.ashx?v=211&...",
        "breakout": "https://finviz.com/screener.ashx?v=211&..."
      },
      "heat_cap_percent": 0.015,
      "blocked": false,
      "warning": false
    },
    {
      "name": "Utilities",
      "success_rate": 0.0000,
      "allowed_strategies": [],
      "screener_urls": {},
      "heat_cap_percent": 0.0,
      "blocked": true,
      "warning": false,
      "warning_message": "NEVER TRADE - Zero profitable strategies in 28 backtests"
    }
  ],
  "strategies": [
    {
      "name": "Alt43",
      "description": "Volatility-Adaptive Targets",
      "best_sectors": ["Healthcare", "Technology"],
      "success_rate": 0.6190,
      "hold_period_weeks": "3-12 weeks",
      "hold_period_min_weeks": 3,
      "hold_period_max_weeks": 12,
      "options_suitable": true,
      "options_warning": "",
      "recommended_dte": "30-45 DTE"
    }
  ],
  "screeners": {
    "healthcare_universe": {
      "name": "Healthcare Universe",
      "url": "https://finviz.com/screener.ashx?v=211&f=sec_healthcare,cap_midover,sh_avgvol_o500,sh_price_o50,fa_epsyoy_pos,fa_sales5years_pos,fa_roe_pos,ta_sma200_pa&ft=4",
      "description": "30-50 healthcare stocks in confirmed uptrends",
      "expected_results": "30-50 stocks",
      "use_case": "weekly",
      "proven_winners": ["UNH", "LLY", "ABBV", "JNJ"]
    }
  },
  "position_sizing": {
    "minimum_score": 5,
    "multipliers": {
      "5": 0.5,
      "6": 0.75,
      "7": 1.0,
      "8": 1.25
    },
    "minimum_contracts": 1,
    "description": "Poker-style sizing: Better checklist score = larger position"
  }
}
```

**policy.v1.json.sha256 (integrity checksum):**

```
a1b2c3d4e5f6... (SHA256 hash of policy.v1.json)
```

**trades.json:**

```json
{
  "version": "2.0",
  "last_updated": "2025-11-03T10:30:00Z",
  "trades": [
    {
      "id": "trade-001",
      "created_at": "2025-11-03T09:15:00Z",
      "sector": "Healthcare",
      "ticker": "UNH",
      "strategy": "Alt43",
      "position_size": 1,
      "max_loss": 850.00,
      "options_strategy": "Bull Call Spread",
      "expiration_date": "2025-12-20T16:00:00Z",
      "checklist_passed": true,
      "checklist_score": 7,
      "exit_date": null,
      "profit_loss": null
    }
  ]
}
```

**settings.json:**

```json
{
  "version": "2.0",
  "theme_mode": "night",
  "account_equity": 100000.00,
  "risk_per_trade": 0.0075,
  "portfolio_heat_cap": 0.04,
  "bucket_heat_cap": 0.015,
  "sector_buckets": {},
  "vimium_enabled": false,
  "sample_data_mode": false,
  "show_welcome_screen": true
}
```

**backtest_results.json (LEGACY - kept for reference):** This file is no longer the primary source of truth. All sector/strategy rules are now in `policy.v1.json`. This file is kept for historical reference and can be used to regenerate policy files if needed.

### Policy File Integrity Checking

**At Application Startup:**

```go
func LoadAndValidatePolicy(policyPath, checksumPath string) (*Policy, error) {
    // Read policy file
    policyData, err := os.ReadFile(policyPath)
    if err != nil {
        log.Error("Policy file not found, entering safe mode")
        return GetSafeModePolicy(), nil
    }
    
    // Calculate SHA256
    hash := sha256.Sum256(policyData)
    calculatedHash := hex.EncodeToString(hash[:])
    
    // Read expected checksum
    expectedHash, err := os.ReadFile(checksumPath)
    if err != nil {
        log.Warn("Checksum file not found, using calculated hash")
    } else if string(expectedHash) != calculatedHash {
        log.Error("Policy checksum mismatch! Entering safe mode")
        return GetSafeModePolicy(), nil
    }
    
    // Parse policy
    var policy Policy
    if err := json.Unmarshal(policyData, &policy); err != nil {
        log.Error("Policy parse failed, entering safe mode")
        return GetSafeModePolicy(), nil
    }
    
    log.Info("Policy loaded successfully", "version", policy.Version)
    return &policy, nil
}

// Safe mode policy (fallback if integrity check fails)
func GetSafeModePolicy() *Policy {
    return &Policy{
        Version:     "safe-mode",
        GeneratedAt: time.Now(),
        SourceNotes: "Safe mode fallback - Healthcare/Tech allowed, Utilities blocked",
        Sectors: []Sector{
            {
                Name:              "Healthcare",
                SuccessRate:       0.9231,
                AllowedStrategies: []string{"Alt43", "Alt46", "Alt10"},
                Blocked:           false,
            },
            {
                Name:              "Technology", 
                SuccessRate:       0.7143,
                AllowedStrategies: []string{"Alt15", "Alt22", "Alt10"},
                Blocked:           false,
            },
            {
                Name:    "Utilities",
                Blocked: true,
            },
            {
                Name:           "Energy",
                Warning:        true,
                WarningMessage: "Low success rate (21.43%), use with caution",
            },
        },
    }
}
```

**When to Use Safe Mode:**

- Policy file missing or corrupted
- Checksum mismatch (file was tampered with)
- JSON parse errors
- Invalid policy version

**Safe Mode Behavior:**

- Healthcare and Technology sectors enabled with core strategies
- Utilities completely blocked
- Energy shows warning
- User notified via dialog: "Running in safe mode. Contact support."
- Application logs the error for troubleshooting

------

### Auto-Save Logic

**When to Save:**

1. After EACH screen completion (before navigation)
2. On application exit (graceful shutdown)
3. Every 60 seconds (background auto-save)

**What to Save:**

- Current trade-in-progress (even if incomplete)
- All completed trades
- User settings changes

**How to Save:**

```go
func (a *App) SaveTrade(trade *Trade) error {
    // Load existing trades
    data, err := LoadTrades("data/trades.json")
    if err != nil {
        return err
    }
    
    // Update or append
    found := false
    for i, t := range data.Trades {
        if t.ID == trade.ID {
            data.Trades[i] = *trade
            found = true
            break
        }
    }
    if !found {
        data.Trades = append(data.Trades, *trade)
    }
    
    // Save to disk
    data.LastUpdated = time.Now()
    return SaveTrades("data/trades.json", data)
}

// Background auto-save goroutine
func (a *App) StartAutoSave() {
    ticker := time.NewTicker(60 * time.Second)
    go func() {
        for range ticker.C {
            if a.State.CurrentTrade != nil {
                a.SaveTrade(a.State.CurrentTrade)
            }
        }
    }()
}
```

------

## ⌨️ VIMIUM MODE

### Keyboard Navigation

**When Enabled:**

- All buttons get keyboard shortcuts
- Modal popup shows shortcuts
- Navigate without mouse

**Key Bindings:**

```
General:
  ?       - Show keyboard shortcuts
  Esc     - Cancel current action
  Enter   - Confirm / Next
  Backspace - Go back
  
Navigation:
  j/k     - Scroll down/up
  g g     - Go to top
  G       - Go to bottom
  
Dashboard:
  n       - Start new trade
  r       - Resume session
  c       - View calendar
  h       - Toggle help
  
Workflow:
  1-8     - Jump to screen 1-8 (if allowed)
  Space   - Continue to next screen
  b       - Go back
  q       - Cancel (with confirmation)
```

**Implementation:**

```go
type VimiumMode struct {
    Enabled     bool
    Bindings    map[rune]func()
    HelpVisible bool
}

func (vm *VimiumMode) HandleKey(key rune) {
    if !vm.Enabled {
        return
    }
    
    if action, ok := vm.Bindings[key]; ok {
        action()
    }
}

// Register bindings
func (a *App) SetupVimium() {
    a.Vimium.Bindings = map[rune]func(){
        'n': func() { a.NavigateTo(SECTOR_SELECT) },
        'c': func() { a.NavigateTo(CALENDAR_VIEW) },
        'h': func() { a.ShowHelp() },
        '?': func() { a.ShowKeyboardShortcuts() },
    }
}
```

------

## 🧪 SAMPLE DATA MODE

### Purpose

- Test UI without real trades
- Demo the application
- Training mode for new users

### Implementation

```go
func LoadSampleData() []Trade {
    return []Trade{
        {
            ID:               "sample-001",
            Sector:           "Healthcare",
            Ticker:           "UNH",
            Strategy:         "Alt43",
            PositionSize:     1,
            MaxLoss:          850,
            OptionsStrategy:  "Bull Call Spread",
            ExpirationDate:   time.Now().AddDate(0, 0, 47),
            ChecklistPassed:  true,
            ChecklistScore:   7,
        },
        {
            ID:               "sample-002",
            Sector:           "Technology",
            Ticker:           "MSFT",
            Strategy:         "Alt15",
            PositionSize:     2,
            MaxLoss:          1200,
            OptionsStrategy:  "LEAPS Diagonal",
            ExpirationDate:   time.Now().AddDate(0, 0, 14),
            ChecklistPassed:  true,
            ChecklistScore:   8,
        },
        // ... more sample trades
    }
}

// Toggle sample mode
func (a *App) ToggleSampleMode() {
    a.Settings.SampleDataMode = !a.Settings.SampleDataMode
    if a.Settings.SampleDataMode {
        a.State.Trades = LoadSampleData()
    } else {
        a.State.Trades = LoadTrades("data/trades.json")
    }
    a.RefreshUI()
}
```

------

## 🏗️ DEVELOPMENT PHASES

### Phase 1: MVP (Weeks 1-4)

**Goal:** Core workflow from sector selection to calendar view

**Week 1:** Project setup + Dashboard + Sector selection

- Set up Go module, Fyne dependencies
- Implement theme system (day/night)
- Build dashboard screen
- Build sector selection screen
- Load backtest_results.json

**Week 2:** Screener launch + Ticker entry + Cooldown

- Implement FINVIZ URL launcher
- Build ticker + strategy screen
- Implement 2-minute cooldown timer
- Auto-save after each screen

**Week 3:** Checklist + Position sizing

- Build checklist screen with 5+3 gates
- Implement poker-style position calculator
- Validation logic for all inputs

**Week 4:** Heat check + Trade entry + Calendar

- Build heat calculator
- Build options strategy form (all 24 types)
- Build calendar grid view
- End-to-end workflow testing

### Phase 2: Enhancement (Weeks 5-6)

**Goal:** Polish, sample data, trade management

**Week 5:** Trade management + Sample data

- Build trade edit/delete screen
- Implement sample data mode
- Export to CSV functionality

**Week 6:** Vimium + Help + Welcome

- Implement keyboard navigation
- Build help system
- Polish welcome screen
- Bug fixes and UI improvements

### Phase 3: Testing & Deployment (Week 7-8)

**Goal:** Comprehensive testing and Windows deployment

**Week 7:** Testing

- Unit tests (data models, calculations)
- Integration tests (screen flows)
- User acceptance testing

**Week 8:** Deployment

- Windows installer (using Fyne package)
- User guide documentation
- Final bug fixes
- Launch

------

## 🔧 TECHNICAL SPECIFICATIONS

### Dependencies

**Required:**

```
go 1.21+
fyne.io/fyne/v2 v2.4+
```

**Optional:**

```
github.com/google/uuid (for trade IDs)
github.com/stretchr/testify (for testing)
```

### Build Commands

**Development:**

```bash
go run main.go
```

**Production (Windows):**

```bash
fyne package -os windows -icon assets/icon.png
```

**Testing:**

```bash
go test ./...
```

### Performance Targets

**Startup Time:** < 2 seconds
 **Screen Transition:** < 100ms
 **Calendar Render:** < 500ms (for 100 trades)
 **Memory Usage:** < 100MB
 **File Save:** < 50ms

### Error Handling

**Categories:**

1. **User Errors:** Show friendly dialog, allow correction
2. **Data Errors:** Show error, offer to reset or restore backup
3. **System Errors:** Log to file, show technical details, offer support contact
4. **Policy Errors:** Hash mismatch → fall back to last known good policy → enter safe mode

**Example:**

```go
func (a *App) HandleError(err error, severity string) {
    switch severity {
    case "user":
        dialog.ShowError(err, a.Window)
    case "data":
        dialog.ShowConfirm(
            "Data Error",
            fmt.Sprintf("Error: %v. Reset to defaults?", err),
            func(ok bool) {
                if ok {
                    a.ResetToDefaults()
                }
            },
            a.Window,
        )
    case "system":
        a.LogError(err)
        dialog.ShowError(
            fmt.Errorf("System error: %v. Please contact support.", err),
            a.Window,
        )
    case "policy":
        // Policy file corrupted or tampered with
        a.LogError(fmt.Errorf("Policy integrity check failed: %v", err))
        
        // Try to load last known good policy from backup
        lastGood, backupErr := a.LoadBackupPolicy()
        if backupErr == nil {
            a.State.Policy = lastGood
            dialog.ShowInformation(
                "Safe Mode",
                "Policy file integrity check failed. Running with last known good policy. Please update policy file.",
                a.Window,
            )
        } else {
            // No backup available, use hardcoded safe mode
            a.State.Policy = GetSafeModePolicy()
            dialog.ShowInformation(
                "Safe Mode",
                "Policy file corrupted. Running in safe mode (Healthcare/Tech only). Contact support.",
                a.Window,
            )
        }
    }
}
```

------

## 📝 CODING STANDARDS

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golint` for linting
- Keep functions under 50 lines
- Use descriptive variable names

### File Organization

- One struct per file (main structs)
- Group related functions together
- Keep UI code separate from business logic
- Use interfaces for testability

### Comments

- Document all exported functions
- Explain "why" not "what"
- Use godoc format

### Testing

- Minimum 80% code coverage
- Test all calculations (position sizing, heat)
- Mock external dependencies (FINVIZ, file system)

------

## 🚀 DEPLOYMENT

### Windows Installer

**Using Fyne packaging:**

```bash
# Generate Windows .exe with icon
fyne package -os windows -icon assets/icon.png -name "TF-Engine"

# Output: TF-Engine.exe
```

**Include in installer:**

- TF-Engine.exe
- data/ folder (with backtest_results.json)
- assets/ folder (with help.txt, sample_trades.json)
- README.txt (quick start guide)

### First-Run Experience

1. Launch TF-Engine.exe
2. Welcome screen appears
3. User clicks "Start New Trade"
4. Creates data/trades.json automatically
5. Creates data/settings.json with defaults

------

## ✅ DONE CRITERIA

### MVP is "Done" when:

- [ ] User can complete full workflow (8 screens)
- [ ] All gates enforced (cooldown, checklist, heat check)
- [ ] Trades saved automatically
- [ ] Calendar view shows timeline correctly
- [ ] Day/night themes work
- [ ] No data loss (auto-save works)
- [ ] FINVIZ URLs open correctly
- [ ] Position sizing calculator accurate
- [ ] Heat calculator accurate

### Phase 2 is "Done" when:

- [ ] Trade management works (edit/delete)
- [ ] Sample data mode toggles correctly
- [ ] Vimium keyboard navigation works
- [ ] Help system accessible
- [ ] Welcome screen shows on first run

### Phase 3 is "Done" when:

- [ ] Windows installer tested
- [ ] User guide written
- [ ] All tests passing (>80% coverage)
- [ ] No critical bugs
- [ ] Deployed to your Windows machine

------

## 🎯 SUCCESS METRICS

**User Behavior:**

- Trade completion rate: >90% (users finish workflow once started)
- Cooldown skip attempts: 0 (enforced correctly)
- Checklist bypass attempts: 0 (enforced correctly)
- Heat limit violations: 0 (enforced correctly)

**Technical:**

- Crash rate: <0.1% of sessions
- Data loss rate: 0% (auto-save works)
- Startup time: <2 seconds
- Average workflow time: 5-8 minutes (including cooldown)

**Business:**

- Trades entered match backtest recommendations: >80%
- User follows healthcare/tech priority: >70% of trades
- User avoids utilities/energy: 100% (blocked/warned)

------

## 📞 SUPPORT & MAINTENANCE

### Bug Reporting

- GitHub Issues (if open source)
- Email support (if private)
- Include: OS version, Go version, error logs

### Updates

- Minor updates: Bug fixes, UI tweaks
- Major updates: New features (per architectural rules)
- Versioning: Semantic (v2.1.0, v2.2.0, etc.)

------

## 🏁 FINAL NOTES

This architecture is designed to:

1. **Enforce discipline** (cooldown, checklist, heat limits)
2. **Leverage your data** (293 backtests drive all recommendations)
3. **Stay simple** (linear workflow, no feature bloat)
4. **Be maintainable** (clear structure, good separation of concerns)
5. **Respect your constraints** (no extra features without approval)

**Remember:** The goal is not to build the fanciest trading app. The goal is to build an app that helps you follow your proven system consistently.

**Your 293 backtests showed what works. This app enforces it.**

------

**Ready to build? Start with Phase 1, Week 1: Project setup + Dashboard.**
