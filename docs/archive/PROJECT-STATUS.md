# TF-Engine 2.0 - Project Status

**Created:** November 3, 2025
**Status:** Skeleton Complete, Ready for Development
**Next Phase:** Build out screen implementations

---

## ✅ What's Been Created

### 1. Project Structure
```
tf-engine/
├── main.go                              # ✅ Entry point (complete)
├── go.mod                               # ✅ Dependencies configured
├── .gitignore                           # ✅ Prevents committing user data
│
├── internal/
│   ├── appcore/
│   │   └── state.go                     # ✅ Global app state + cooldown logic
│   ├── models/
│   │   ├── policy.go                    # ✅ Policy loading + safe mode
│   │   ├── trade.go                     # ✅ Trade data structure
│   │   └── settings.go                  # ✅ User preferences
│   └── ui/
│       ├── theme.go                     # ✅ Day/night mode themes
│       ├── dashboard.go                 # ✅ Main dashboard (stub)
│       └── screens/                     # ✅ All 8 screens (stubs)
│           ├── sector_selection.go
│           ├── screener_launch.go
│           ├── ticker_entry.go
│           ├── checklist.go
│           ├── position_sizing.go
│           ├── heat_check.go
│           ├── trade_entry.go
│           └── calendar.go
│
├── data/
│   └── policy.v1.json                   # ✅ Already exists (your config)
│
├── scripts/
│   ├── build.bat                        # ✅ Windows build script
│   └── run_dev.bat                      # ✅ Development runner
│
├── CLAUDE.md                            # ✅ Guidance for future Claude instances
├── README-DEV.md                        # ✅ Developer quick start
└── PROJECT-STATUS.md                    # ✅ This file
```

### 2. Key Features Implemented

#### ✅ Policy-Driven Architecture
- Loads from `data/policy.v1.json`
- Falls back to safe mode if policy missing/corrupted
- Sector filtering based on backtest success rates
- Strategy-sector mapping enforced

#### ✅ Application State Management
- `AppState` struct holds current trade, settings, policy
- Cooldown timer logic (120 seconds)
- Methods for checking cooldown completion

#### ✅ Data Models
- `Trade` struct with all 8 screens' data fields
- `Policy` struct matching your JSON schema
- `Settings` struct for user preferences
- `Sector` and `Strategy` types

#### ✅ UI Foundation
- Custom theme with day/night modes
- Green color palette (light green for day, British racing green for night)
- Dashboard layout structure
- All 8 screen files created with basic structure

#### ✅ Development Tools
- Build script (`scripts/build.bat`)
- Dev run script (`scripts/run_dev.bat`)
- `.gitignore` configured
- Go module initialized with Fyne v2.7.0

---

## 🚧 What Still Needs Implementation

### High Priority (MVP)

1. **Screen Navigation System**
   - Linear flow between screens 1→8
   - Back button functionality
   - Cancel with confirmation dialog

2. **Screen 2: FINVIZ Integration**
   - Launch browser with screener URLs from policy
   - Windows-specific `exec.Command` implementation

3. **Screen 3: Cooldown Timer UI**
   - Visual countdown display (progress bar + time remaining)
   - Disable "Continue" button until timer expires
   - Persist cooldown state across app restarts

4. **Screen 4: Checklist Validation**
   - Load gates from policy checklist
   - Track which gates are checked
   - Calculate score (0-8)
   - Enforce minimum score of 5 required gates

5. **Screen 5: Position Size Calculator**
   - Poker-style multipliers from policy
   - Risk calculation: (equity × risk%) / (entry - stop)
   - Contract quantity calculation
   - Display adjusted position based on checklist score

6. **Screen 6: Heat Check Logic**
   - Calculate current portfolio heat from all open trades
   - Calculate sector bucket heat
   - Validate against caps (4% portfolio, 1.5% sector)
   - Block trade if limits exceeded

7. **Screen 7: Options Strategy Form**
   - Dynamic fields based on strategy type
   - Strike price validation
   - Expiration date picker
   - DTE calculation and warning

8. **Screen 8: Calendar Timeline Widget**
   - Custom Fyne widget
   - Y-axis: Sectors (from policy)
   - X-axis: Time (-14 days to +84 days)
   - Horizontal bars for each trade
   - Ticker labels on bars
   - Heat summary at bottom

9. **Data Persistence**
   - Save trade to `data/trades.json`
   - Load all trades on startup
   - Auto-save after each screen
   - Resume in-progress trade

### Medium Priority (Phase 2)

10. **Screen 9: Trade Management**
    - Table of all trades
    - Edit trade details
    - Delete trade functionality
    - Filter: All / Active / Closed

11. **Help System**
    - Help button with question mark icon
    - Popup dialogs with guidance
    - Welcome screen on first startup

12. **Sample Data Generator**
    - Create 8-12 realistic sample trades
    - Vary sectors, dates, strategies
    - Populate for testing calendar view

13. **Vimium Mode**
    - Keyboard shortcuts
    - Toggle on/off
    - Visual indicator when active

### Low Priority (Polish)

14. **Theme Refinement**
    - Fine-tune colors for contrast
    - Test readability on both modes
    - Add hover states, focus indicators

15. **Error Handling**
    - User-friendly error messages
    - Validation feedback
    - Network connectivity checks for FINVIZ

16. **Testing**
    - Unit tests for calculations
    - Integration tests for screen flows
    - Manual test checklist

---

## 📋 Recommended Development Order

### Week 1: Core Navigation + Screens 1-3
1. Implement screen navigation system
2. Complete Screen 1 (Sector Selection) with real policy data
3. Complete Screen 2 (FINVIZ launcher)
4. Complete Screen 3 (Ticker Entry + Cooldown Timer)

**Milestone:** Can select sector, launch screener, enter ticker, wait 2 minutes

### Week 2: Screens 4-6 (Business Logic)
5. Complete Screen 4 (Checklist with validation)
6. Complete Screen 5 (Position Sizing with poker multipliers)
7. Complete Screen 6 (Heat Check with enforcement)

**Milestone:** Can complete anti-impulsivity workflow and enforce portfolio limits

### Week 3: Screens 7-8 (Trade Entry + Visualization)
8. Complete Screen 7 (Options Strategy entry)
9. Build Screen 8 (Calendar timeline widget)
10. Implement data persistence (JSON I/O)

**Milestone:** Can enter full trade and see it on calendar

### Week 4: Polish + Phase 2 Features
11. Add Screen 9 (Trade Management)
12. Add help system and welcome screen
13. Add sample data generator
14. Testing and bug fixes

**Milestone:** MVP complete and usable for live trading

---

## 🏃 Quick Start for Development

### Run the App (Development Mode)
```bash
# Option 1: Use the script
scripts\run_dev.bat

# Option 2: Manual
go run main.go
```

### Build Executable
```bash
# Option 1: Use the script
scripts\build.bat

# Option 2: Manual
go build -o tf-engine.exe
```

### Test the Skeleton
The current skeleton should launch a Fyne window showing:
- Dashboard with buttons
- "Start New Trade" button (currently does nothing)
- Account and heat status (placeholder text)

**Expected behavior:** Window opens, shows green-themed UI, buttons are visible but mostly non-functional (stubs).

---

## 📝 Development Notes

### Policy-Driven Reminder
All sector rules, strategies, and business logic come from `data/policy.v1.json`. Never hardcode:
- Sector names or success rates
- Strategy names or allowed sectors
- Heat caps or risk percentages
- Checklist items or scoring

### Anti-Impulsivity is Non-Negotiable
The 120-second cooldown, 5-gate checklist, and heat limits are **core requirements**, not optional features. Do not allow shortcuts.

### Auto-Save Everything
After EVERY screen transition, persist the current trade state. User should never lose progress.

### No Feature Creep
Stick to the 8-screen workflow. Don't add features that aren't in the architectural docs without explicit approval.

---

## 🎯 Success Criteria

The application is **complete** when:

1. ✅ User can progress through all 8 screens sequentially
2. ✅ Cooldown timer enforces 120-second wait
3. ✅ Checklist requires all 5 gates to pass
4. ✅ Heat check blocks trades that exceed caps
5. ✅ Calendar shows all trades with correct dates
6. ✅ Trades persist across app restarts
7. ✅ Blocked sectors (Utilities) cannot be selected
8. ✅ Day/night mode works with readable text

The application is **ready for live trading** when:

1. ✅ Sample data mode works (can test without real trades)
2. ✅ Help system explains each screen
3. ✅ Vimium mode enables keyboard navigation
4. ✅ User has tested complete workflow 3+ times
5. ✅ No crashes or data loss during testing

---

## 🚀 Ready to Build!

The skeleton is complete. Everything compiles. Dependencies are installed.

**Next step:** Pick Screen 1 (Sector Selection) and make it fully functional with real policy data, then move to Screen 2, and so on.

Good luck! 🎯
