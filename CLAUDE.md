# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

# GOLDEN RULES

Always review the rules when starting a new conversation.

Keep the project README.md up to date but not bloated.

Keep a running lessons-learned.md file to keep tabs on when a feature just isn't coming together or we're getting repeat errors that won't fix.

Please refer to and use the /.claude/agents as appropriate, depending on the task at hand.  Don't over-rely on them over your own innate judgment though.  They're a tool, not the solution to everything.  

---

---

## CRITICAL ARCHITECTURAL RULES

### Rule #1: No Feature Creep
**DO NOT create new extra features if not explicitly asked for or at least approved by the user architect.**

This rule exists because the previous project failed due to excessive features. This application must stay ruthlessly focused on the core workflow. When the user requests a feature, implement exactly what was asked—no more, no less.

### Rule #2: Policy-Driven Design
**The application is controlled by [data/policy.v1.json](data/policy.v1.json)—never hardcode business logic.**

All sector configurations, strategy mappings, portfolio limits, and behavioral rules live in this JSON file. The app reads and enforces these policies at runtime. This allows the user to update trading rules without recompiling code.

### Rule #3: Anti-Impulsivity is Non-Negotiable
**Every trade must pass through cooldown timers, checklists, and heat checks.**

This isn't optional UX polish—these are behavioral finance guardrails that prevent emotional trading mistakes. Never allow shortcuts or "skip" buttons for these safety mechanisms.

---

## PROJECT OVERVIEW

**Project Type:** Trend-Following Options Trading Decision Support System
**Target Platform:** Windows Desktop Application
**Technology Stack:** Go + Fyne GUI framework (planned)
**Current Phase:** Planning/Architecture (research phase completed)

### Purpose
This application guides options traders through a systematic decision workflow based on 293 validated backtests across 14 trend-following strategies and 21 securities. The research discovered that strategy performance is highly sector-dependent, leading to a sector-first workflow design.

### Research Foundation
The app is informed by extensive backtesting research documented in [DISCOVERIES_AND_LEARNINGS.md](DISCOVERIES_AND_LEARNINGS.md):
- **293 validated backtests** (99.74% data quality)
- **Key finding:** Alt10 (Profit Targets) achieved 76.19% success rate
- **Critical insight:** Healthcare 92% success rate, Utilities 0% success rate
- **Python validation:** Logistic regression confirms profit targets (4.47× odds) and pyramiding (3.96× odds) drive profitability
- **Sector truth:** Not all sectors are tradeable with trend-following strategies

---

## CORE ARCHITECTURE

### Master Configuration: policy.v1.json

This JSON file is the single source of truth for all application behavior:

```json
{
  "sectors": [
    {
      "name": "Healthcare",
      "priority": 1,
      "blocked": false,
      "warning": false,
      "allowed_strategies": ["Alt10", "Alt46", "Alt43", "Alt39", "Alt28"],
      "screener_urls": {
        "universe": "...",
        "pullback": "...",
        "breakout": "...",
        "golden_cross": "..."
      }
    }
  ],
  "strategies": {
    "Alt10": {
      "label": "Profit Targets (3N/6N/9N)",
      "options_suitability": "excellent",
      "hold_weeks": "3-10"
    }
  },
  "checklist": {
    "required": ["SIG_REQ", "RISK_REQ", "OPT_REQ", "EXIT_REQ", "BEHAV_REQ"],
    "poker_sizing": { "5": 0.5, "6": 0.75, "7": 1.0, "8": 1.25 }
  },
  "defaults": {
    "portfolio_heat_cap": 0.04,
    "bucket_heat_cap": 0.015,
    "cooldown_seconds": 120
  }
}
```

**Key Design Principles:**
- Sectors define which strategies are allowed (not all strategies work everywhere)
- Blocked sectors (Utilities) prevent trades entirely
- Warned sectors (Energy, Real Estate) show cautionary messages
- Portfolio heat limits: 4% total, 1.5% per sector
- Anti-impulsivity: 120-second cooldown before trade execution

---

## APPLICATION WORKFLOW (9 Screens)

The app guides users through a sequential decision process. Progress auto-saves after each screen.

### Screen 1: Sector Selection
**Purpose:** Choose which market sector to trade
**Logic:** Present sectors from policy.json, sorted by priority
**Constraints:**
- Grey out blocked sectors (e.g., Utilities with 0% backtest success)
- Show warning icon for marginal sectors (Energy, Real Estate)
- Display sector notes from policy (e.g., "Healthcare: 92% strategy success")

### Screen 2: Screener Results
**Purpose:** Launch Finviz screeners to find trade candidates
**Logic:** Display buttons for each screener URL defined in selected sector
**Screener Types:**
- **Universe** (weekly): 30-60 quality stocks in uptrends
- **Pullback** (daily): Oversold stocks in uptrends (RSI < 40)
- **Breakout** (daily): New 52-week highs
- **Golden Cross** (daily): SMA50 crossing above SMA200

**Implementation Note:** URLs open in default browser with `v=211` parameter for chart view

### Screen 3: Ticker Entry & Strategy Selection
**Purpose:** Enter ticker symbol and select matching strategy
**Logic:**
1. User types ticker (e.g., "UNH", "MSFT")
2. Dropdown populates with ONLY strategies allowed for the selected sector (from `allowed_strategies` array in policy.json)
3. Display strategy metadata: label, options suitability, typical hold weeks
4. Start 120-second anti-impulsivity cooldown when user proceeds

**Critical:** Strategy dropdown is sector-filtered. Healthcare shows different strategies than Technology.

### Screen 4: Anti-Impulsivity Checklist
**Purpose:** Force deliberate decision-making through required criteria
**Logic:**
- Display 5 required checkboxes (from `checklist.required` in policy)
- Display 3 optional checkboxes (from `checklist.optional`)
- Enable "Continue" button only when all 5 required items checked
- Cooldown timer displays remaining seconds

**Checklist Examples:**
- SIG_REQ: "Is price above SMA200 with confirmed Donchian breakout?"
- RISK_REQ: "Is stop-loss placement acceptable for account size?"
- OPT_REQ: "Does strategy hold duration match options expiration?"
- EXIT_REQ: "Do I have clear profit targets and exit plan?"
- BEHAV_REQ: "Am I emotionally calm and following my system?"

### Screen 5: Position Sizing Calculator
**Purpose:** Calculate position size using poker-bet sizing principles
**Logic:**
1. User rates trade conviction: 5-8 (from `checklist.poker_sizing`)
2. App multiplies base risk by sizing multiplier:
   - 5 → 0.5× (weak conviction)
   - 6 → 0.75×
   - 7 → 1.0× (standard)
   - 8 → 1.25× (strong conviction)
3. Display: "Risk $X per contract based on Y% account risk"

**Formula:** `contracts = (account_size × sizing_multiplier × risk_per_trade) / (strike - stop_loss)`

### Screen 6: Portfolio Heat Check
**Purpose:** Enforce diversification limits and prevent concentration risk
**Logic:**
1. Calculate current portfolio allocation by sector
2. Display heat bars: Healthcare 1.2% / 1.5% cap, Technology 0.8% / 1.5% cap, etc.
3. Check if new trade would exceed:
   - Sector cap: 1.5% (from `bucket_heat_cap` in policy)
   - Portfolio cap: 4% (from `portfolio_heat_cap`)
4. Block trade if limits exceeded; allow user to close existing position first

**Visual Design:** Green bars (safe), yellow bars (approaching limit), red bars (at limit)

### Screen 7: Options Strategy Selection
**Purpose:** Select specific options structure
**Options List (24 types):**
- Bull call spread
- Bear put spread
- Bull put credit spread
- Bear call credit spread
- Long call
- Long put
- Covered call
- Cash-secured put
- Iron butterfly
- Iron condor
- Long put butterfly
- Long call butterfly
- Calendar call spread
- Calendar put spread
- Diagonal call spread
- Diagonal put spread
- Inverse iron butterfly
- Inverse iron condor
- Short put butterfly
- Short call butterfly
- Straddle
- Strangle
- Call ratio backspread
- Put ratio backspread
- Call broken wing
- Put broken wing

**UI:** Dropdown or grid layout with strategy diagrams

### Screen 8: Trade Calendar (Dashboard - "Horserace View")
**Purpose:** Visualize all active trades across time and sectors
**Layout:**
- **Y-axis:** Sectors (Healthcare, Technology, Industrials, etc.)
- **X-axis:** Time (-14 days to +84 days from current date)
- **Bars:** Each trade represented as horizontal bar showing entry → expiration
- **Bar Label:** Ticker symbol (e.g., "UNH", "MSFT")

**Example Visual:**
```
Healthcare    [--------UNH butterfly-----]    [--XLV call spread--]
Technology           [---MSFT call spread-------------]
Industrials     [--CAT put spread--]
Consumer            (no active trades)
```

**Configuration:** From `calendar` object in policy.json:
- `past_days: 14` - Show 2 weeks of history
- `future_days: 84` - Show 12 weeks forward
- `y_axis: "sector"` - Group by sector
- `bar_label: "ticker"` - Label bars with ticker symbols

**Color Coding:**
- Green bars: Profitable trades (mark-to-market)
- Red bars: Losing trades
- Yellow bars: Expiring within 7 days

### Screen 9: Trade Management
**Purpose:** Edit or delete past trades
**Logic:**
- Display table of all trades (active + closed)
- Columns: Date, Ticker, Sector, Strategy (Pine Script), Options Type, P&L, Status
- Actions: Edit button, Delete button
- Filter: Show All / Active Only / Closed Only

**Data Persistence:** Trades stored in JSON file or SQLite database, loaded on startup

---

## DEVELOPMENT GUIDELINES

### Technology Stack (Planned)

**Language:** Go
**GUI Framework:** Fyne (cross-platform, but targeting Windows)
**Data Storage:** JSON files or SQLite for trade history
**Configuration:** JSON (policy.v1.json)

### Build Commands

```bash
# Install Fyne
go get fyne.io/fyne/v2

# Run development version
go run .

# Build Windows executable
go build -o trend-following-dashboard.exe

# Build with Fyne bundler (includes assets)
fyne package -os windows -icon icon.png
```

### Project Structure (To Be Created)

```
new-trend-following-dashboard/
├── main.go                      # Entry point, app initialization
├── config/
│   └── policy.go               # Load and parse policy.v1.json
├── screens/
│   ├── sector_selection.go     # Screen 1
│   ├── screener_results.go     # Screen 2
│   ├── ticker_entry.go         # Screen 3
│   ├── checklist.go            # Screen 4
│   ├── position_sizing.go      # Screen 5
│   ├── heat_check.go           # Screen 6
│   ├── options_strategy.go     # Screen 7
│   ├── trade_calendar.go       # Screen 8
│   └── trade_management.go     # Screen 9
├── models/
│   ├── sector.go               # Sector data structure
│   ├── strategy.go             # Strategy data structure
│   └── trade.go                # Trade data structure
├── storage/
│   └── trades.go               # Persistence layer (JSON/SQLite)
├── widgets/
│   ├── heat_bar.go             # Custom portfolio heat visualization
│   └── calendar_view.go        # Horserace timeline widget
├── data/
│   └── policy.v1.json          # Master configuration (exists)
└── assets/
    ├── icon.png
    └── welcome_screen.png
```

### Testing Strategy

Since this is a GUI application with behavioral safety features:

1. **Unit Tests:** Test policy.json parsing, heat calculations, position sizing math
2. **Integration Tests:** Test screen transitions, data persistence, cooldown timers
3. **Manual Testing:** Use sample data feature to populate calendar with test trades
4. **User Acceptance:** Test with real market scenarios before going live

---

## STYLING REQUIREMENTS

### Color Palette

**Day Mode:**
- Background: Light grey/white
- Primary accent: Light green (#6FCF97 or similar)
- Text: Dark grey/black
- Warnings: Amber
- Errors: Red

**Night Mode:**
- Background: Dark grey (#1E1E1E)
- Primary accent: Forest green / British racing green (#0F4C3A or #004225)
- Text: Light grey/white
- Warnings: Amber (brighter)
- Errors: Red (brighter)

**Critical Design Principle:** Text contrast is paramount. Numbers and letters must be highly readable against all backgrounds. Avoid low-contrast green-on-green combinations.

### UI Style Guidelines

- **Modern and sleek:** Avoid outdated Windows XP look
- **Minimal:** Don't clutter screens with unnecessary information
- **Sequential flow:** Users should naturally progress through screens 1→9
- **Persistent navigation:** Tab bar at top shows current position in workflow
- **Help accessibility:** Question mark icon always visible in top-right corner

---

## KEY FEATURES

### Auto-Save Functionality
**Requirement:** Progress must auto-save after each screen transition.

**Implementation:**
- When user clicks "Continue" on any screen, persist current trade state
- Use JSON file: `trades_in_progress.json` or SQLite temp table
- On app restart, detect incomplete trade and offer: "Resume?" or "Start New?"

**Why Critical:** Users may close app mid-workflow; losing progress creates frustration

### Sample Data Generation
**Requirement:** Button to populate calendar with realistic sample trades

**Implementation:**
- Generate 8-12 sample trades across different sectors
- Vary entry dates (-14 days to -1 day)
- Vary expiration dates (+7 days to +70 days)
- Mix profitable and losing trades
- Use actual tickers from policy.json sector screeners

**Why Critical:** Allows testing calendar view without executing real trades

### Vimium Mode (Optional)
**Requirement:** Keyboard navigation inspired by Vimium browser extension

**Implementation:**
- Toggle with `Ctrl+V` or button
- When enabled: Display keyboard shortcuts on screen
- Example shortcuts:
  - `j/k` - Navigate up/down in lists
  - `h/l` - Previous/next screen
  - `Enter` - Select highlighted item
  - `Esc` - Cancel/back

**Why Critical:** Power users want keyboard efficiency

### Welcome Screen
**Requirement:** Show on first startup with option to view later

**Implementation:**
- Modal dialog on app launch (unless "Don't show again" checked)
- Content: Brief overview of 9-screen workflow, links to documentation
- Menu item: Help → Welcome Screen (allows re-display)

---

## BEHAVIORAL FINANCE PRINCIPLES

### The Anti-Impulsivity System

**Problem:** Traders make emotional decisions during market volatility
**Solution:** Multi-layered behavioral guardrails

#### Layer 1: Cooldown Timer (120 seconds)
Starts when user clicks "Continue" from Ticker Entry screen. Forces user to wait 2 minutes before proceeding to position sizing. This interrupts the impulse to "chase" a moving price.

#### Layer 2: Required Checklist
User must consciously check 5 required criteria. Acts as a "pre-flight checklist" similar to aviation safety protocols. Prevents trading on incomplete analysis.

#### Layer 3: Portfolio Heat Limits
Mathematical constraints that prevent overconcentration. Even if user *wants* to overweight a sector, the app blocks it. Removes willpower from the equation.

#### Layer 4: Sector Blacklisting
Utilities sector completely blocked based on 0% backtest success rate. App won't even show strategy dropdown for blocked sectors. Prevents users from "trying anyway" against data.

**Design Philosophy:** Don't rely on trader discipline—engineer discipline into the workflow.

---

## SECTOR-STRATEGY MAPPING LOGIC

### How Strategy Filtering Works

When user selects a sector in Screen 1, the app remembers this choice. When user reaches Screen 3 (Strategy Selection), the dropdown is populated from:

```go
// Pseudocode
selectedSector := userChoice // "Healthcare"
policy := loadPolicy("data/policy.v1.json")

for _, sector := range policy.Sectors {
    if sector.Name == selectedSector {
        allowedStrategies := sector.AllowedStrategies
        // allowedStrategies = ["Alt10", "Alt46", "Alt43", "Alt39", "Alt28"]

        // Populate dropdown with only these strategies
        for _, stratID := range allowedStrategies {
            strategy := policy.Strategies[stratID]
            dropdown.Add(strategy.Label) // "Profit Targets (3N/6N/9N)"
        }
    }
}
```

### Why This Mapping Matters

Research proved that strategy performance is sector-dependent:
- **Healthcare:** Alt10 +33.13%, Alt46 +32.16% (excellent)
- **Utilities:** Alt10 -12.4%, Alt46 -6.2% (catastrophic)

If the app allowed Alt10 on Utilities, users would lose money despite "following the system." The policy.json enforces data-driven constraints.

### Updating Strategy Mappings

When new backtest data becomes available:
1. Run backtests for new sector/strategy combination
2. If strategy achieves >60% success rate, add to `allowed_strategies` array
3. Update policy.v1.json
4. App automatically reflects changes on next launch (no code changes needed)

---

## FINVIZ SCREENER INTEGRATION

### Three-Tier Screening Framework

The app doesn't run screeners itself—it launches pre-configured Finviz URLs in the user's browser:

**Tier 1: Universe Screeners** (Weekly)
Purpose: Find 30-60 quality stocks in long-term uptrends
User runs: Monday morning
App behavior: Opens Finviz URL with filters like:
- Sector: Healthcare
- Market Cap: Mid+ (>$2B)
- Price: Above SMA200 (critical filter)
- Volume: >500K
- Fundamentals: Positive EPS growth, ROE, sales

**Tier 2: Situational Screeners** (Daily)
Purpose: Find 0-10 stocks with specific entry setups TODAY
User runs: Before market open
App behavior: Opens Finviz URL with additional filters:
- **Pullback:** Price above SMA200, below SMA50, RSI < 40
- **Breakout:** Price at 52-week high
- **Golden Cross:** SMA50 just crossed above SMA200

**Tier 3: Execution** (Real-time)
Purpose: Apply Pine Script strategies to situational candidates
User action: Copy tickers from Finviz, paste into TradingView, run strategy
App behavior: N/A (happens outside the app)

### Why External Screeners?

1. **Finviz is best-in-class** for stock screening with 70+ filters
2. **Real-time data** without API costs or data licensing
3. **Visual chart view** (v=211 parameter) helps users identify patterns
4. **No maintenance burden** of building a custom screener

**Implementation:** Use Go's `exec.Command` or Fyne's `OpenURL()` to launch browser

---

## COMMON DEVELOPMENT PATTERNS

### Loading Policy Configuration

```go
type Policy struct {
    Sectors    []Sector            `json:"sectors"`
    Strategies map[string]Strategy `json:"strategies"`
    Checklist  Checklist           `json:"checklist"`
    Defaults   Defaults            `json:"defaults"`
    Calendar   Calendar            `json:"calendar"`
}

func LoadPolicy(path string) (*Policy, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var policy Policy
    err = json.Unmarshal(data, &policy)
    return &policy, err
}
```

### Calculating Portfolio Heat

```go
func CalculateHeatByBucket(trades []Trade) map[string]float64 {
    heatMap := make(map[string]float64)

    for _, trade := range trades {
        if trade.Status == "active" {
            heatMap[trade.Sector] += trade.Risk / accountSize
        }
    }

    return heatMap
}

func CheckHeatLimit(newTrade Trade, policy Policy) bool {
    currentHeat := CalculateHeatByBucket(getActiveTrades())
    bucketHeat := currentHeat[newTrade.Sector] + newTrade.Risk/accountSize

    return bucketHeat <= policy.Defaults.BucketHeatCap
}
```

### Implementing Cooldown Timer

```go
func StartCooldown(duration int, callback func()) {
    ticker := time.NewTicker(1 * time.Second)
    remaining := duration

    go func() {
        for range ticker.C {
            remaining--
            UpdateUI(remaining) // Display countdown

            if remaining <= 0 {
                ticker.Stop()
                callback() // Enable "Continue" button
                return
            }
        }
    }()
}
```

---

## DOCUMENTATION REFERENCES

### For Understanding the Research

- **[DISCOVERIES_AND_LEARNINGS.md](DISCOVERIES_AND_LEARNINGS.md)**: Complete 293-backtest analysis with strategy rankings
- **[README.md](README.md)**: Project overview and methodology
- **[MASTER-SCREENER-GUIDE.md](MASTER-SCREENER-GUIDE.md)**: Three-tier Finviz screening framework

### For Understanding User Intent

- **[architects-intent.md](architects-intent.md)**: Original user requirements and feature specifications
- **[data/policy.v1.json](data/policy.v1.json)**: Current configuration and sector/strategy mappings

### For Learning from Failure

- **from-failed-project-not-all-features-are-desired/**: Previous project attempt with lessons learned
  - Key lesson: Feature creep kills projects—stay focused on core workflow

---

## FUTURE ENHANCEMENTS (Not Implemented Yet)

These are potential features mentioned in documentation but NOT yet approved for implementation:

- Walk-forward analysis integration
- Automated position tracking via broker API
- Real-time P&L updates from market data
- Mobile companion app (iOS/Android)
- Backtesting interface within the app
- Machine learning for trade prediction
- Social sharing of trade ideas

**Remember Rule #1:** Do not implement these without explicit user approval. The goal is a focused tool that solves one problem well, not a Swiss Army knife that does everything poorly.

---

## WORKING WITH THIS CODEBASE

### When Starting a New Feature

1. **Read architects-intent.md** to understand user's original vision
2. **Check policy.v1.json** to see if feature is configuration-driven
3. **Review DISCOVERIES_AND_LEARNINGS.md** for research context
4. **Ask user for approval** if feature wasn't explicitly requested

### When Modifying Existing Code

1. **Preserve behavioral guardrails** (cooldowns, heat checks, checklists)
2. **Maintain policy-driven design** (don't hardcode business rules)
3. **Test with sample data** before real trades
4. **Document breaking changes** in commit messages

### When User Reports a Bug

1. **Verify against policy.json** - Is "bug" actually correct behavior per policy?
2. **Check sector-strategy mapping** - Are they using a blocked combination?
3. **Review heat calculations** - Are limits being enforced correctly?
4. **Test cooldown timers** - Do they complete full 120 seconds?

---

## SUCCESS METRICS

The application is successful if:

1. **Users complete trades systematically** (don't skip checklist steps)
2. **Portfolio heat stays below 4%** (diversification enforced)
3. **Users trade winning sectors** (Healthcare/Tech, avoid Utilities/Energy)
4. **Impulsive trades decrease** (cooldowns and checklists working)
5. **Users can visualize trade lifecycle** (calendar view provides clarity)

The application is failing if:

1. Users circumvent behavioral guardrails
2. Users concentrate risk in losing sectors
3. Users abandon the workflow mid-stream
4. Calendar view is confusing or cluttered
5. Policy changes require code modifications

---

**Last Updated:** November 3, 2025
**Policy Version:** 1.0.0
**Application Status:** Planning phase (research completed, implementation pending)
