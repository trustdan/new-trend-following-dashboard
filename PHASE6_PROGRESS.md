# Phase 6: Warning-Based Trading System - Implementation Progress

**Status:** UI Implementation Complete - Testing & Documentation Pending
**Started:** November 4, 2025
**Last Updated:** November 4, 2025 (Updated: UI complete)

---

## Overview

Phase 6 shifts the architecture from prescriptive (hard-blocking) to permissive (warning-based) to allow user flexibility while preserving research-driven guardrails.

### Key Changes

1. **Unblocked Utilities Sector** - Changed from `"blocked": true` to `"blocked": false` with strong warning modal
2. **Strategy Filtering Removed** - All strategies now shown with color-coded suitability ratings
3. **Warning System Added** - Green/yellow/red indicators with acknowledgement checkboxes
4. **Telemetry Tracking** - Log all warning overrides for future analysis

---

## ✅ Completed Tasks

### 1. Policy Schema Updates (policy.v1.json)

**Status:** ✅ COMPLETE

Added `strategy_suitability` ratings for all 10 sectors:
- Healthcare: 11 strategies rated (5 green, 4 yellow, 2 red)
- Technology: 11 strategies rated (7 green, 3 yellow, 1 red)
- Consumer Discretionary: 8 strategies rated (4 green, 4 yellow)
- Industrials: 8 strategies rated (6 green, 2 yellow)
- Communication Services: 8 strategies rated (6 green, 2 yellow)
- Consumer Defensive: 8 strategies rated (6 green, 1 yellow, 1 red)
- Financials: 7 strategies rated (4 green, 3 yellow)
- Real Estate: 7 strategies rated (1 green, 2 yellow, 4 red)
- Energy: 10 strategies rated (ALL red - mean-reverting sector)
- Utilities: 11 strategies rated (ALL red - 0% backtest success)

**Color Coding:**
- **Green (excellent/good):** No acknowledgement required
- **Yellow (marginal):** Acknowledgement checkbox required
- **Red (incompatible):** Strong acknowledgement checkbox required

### 2. Utilities Warning Configuration

**Status:** ✅ COMPLETE

Added `utilities_warning` to Utilities sector:
```json
{
  "title": "⚠️ Utilities Sector Warning",
  "message": "We have ZERO successful backtests with Utilities...",
  "acknowledgement_text": "I understand Utilities has 0% backtest success..."
}
```

**Key Statistics in Warning:**
- Alt10 (Profit Targets): -12.4%
- Alt46 (Sector-Adaptive): -6.2%
- Alt26 (Fractional Pyramid): -8.9%
- All other strategies: 0% success rate

### 3. Data Model Updates

**Status:** ✅ COMPLETE

#### models/policy.go
Added new structs:
```go
type StrategySuitability struct {
    Rating                 string // "excellent", "good", "marginal", "incompatible"
    Color                  string // "green", "yellow", "red"
    Rationale              string
    RequireAcknowledgement bool
}

type UtilitiesWarning struct {
    Title               string
    Message             string
    AcknowledgementText string
}
```

Updated `Sector` struct with:
- `StrategySuitability map[string]StrategySuitability`
- `UtilitiesWarning *UtilitiesWarning`

#### models/trade.go
Added tracking fields:
```go
StrategySuitability          string // "excellent", "good", "marginal", "incompatible"
StrategyWarningAcknowledged  bool   // True if user acknowledged yellow/red warning
UtilitiesWarningAcknowledged bool   // True if user acknowledged Utilities warning
```

---

## 🚧 Pending Implementation

### 4. UI Updates - ticker_entry.go

**Status:** ✅ COMPLETE

**Requirements:**
1. **Show ALL Strategies** - Remove sector filtering from strategy dropdown
2. **Color Indicators** - Display green/yellow/red indicator next to each strategy
3. **Suitability Info** - Show rationale on hover/selection
4. **Warning Banner** - Display when yellow/red strategy selected
5. **Acknowledgement Checkbox** - Required for yellow/red strategies
6. **Continue Button Logic** - Disable until acknowledgement checked

**Implementation Approach:**
```go
func (s *TickerEntry) createStrategyDropdown() *widget.Select {
    var options []string

    // Get ALL strategies (no filtering)
    for stratID, strategy := range s.policy.Strategies {
        suitability := s.getSuitability(stratID, s.trade.Sector)
        indicator := getIndicator(suitability.Color)

        option := fmt.Sprintf("%s %s - %s", indicator, stratID, strategy.Label)
        options = append(options, option)
    }

    return widget.NewSelect(options, s.onStrategySelected)
}

func (s *TickerEntry) onStrategySelected(selected string) {
    stratID := extractStrategyID(selected)
    suitability := s.getSuitability(stratID, s.trade.Sector)

    s.trade.Strategy = stratID
    s.trade.StrategySuitability = suitability.Rating

    if suitability.RequireAcknowledgement {
        s.showWarningBanner(suitability)
        s.ackCheckbox.Show()
        s.ackCheckbox.SetChecked(false)
        s.continueBtn.Disable()
    } else {
        s.hideWarningBanner()
        s.ackCheckbox.Hide()
        s.continueBtn.Enable()
    }
}

func (s *TickerEntry) getSuitability(stratID, sector string) StrategySuitability {
    for _, sec := range s.policy.Sectors {
        if sec.Name == sector {
            if suit, exists := sec.StrategySuitability[stratID]; exists {
                return suit
            }
        }
    }

    // Default to marginal if not found
    return StrategySuitability{
        Rating: "marginal",
        Color: "yellow",
        Rationale: "Suitability data not available for this combination",
        RequireAcknowledgement: true,
    }
}
```

### 5. UI Updates - sector_selection.go

**Status:** ✅ COMPLETE

**Requirements:**
1. **Utilities Modal** - Show modal when Utilities selected
2. **Modal Content** - Display backtest failure statistics
3. **Acknowledgement Checkbox** - Required to proceed
4. **Go Back Button** - Allow user to cancel
5. **Continue Anyway Button** - Disabled until acknowledged

**Implementation Approach:**
```go
func (s *SectorSelection) onSectorClicked(sector Sector) {
    if sector.UtilitiesWarning != nil {
        s.showUtilitiesModal(sector)
        return
    }

    // Normal flow for other sectors
    s.trade.Sector = sector.Name
    s.navigator.Next()
}

func (s *SectorSelection) showUtilitiesModal(sector Sector) {
    warning := sector.UtilitiesWarning

    ackCheckbox := widget.NewCheck(warning.AcknowledgementText, nil)

    content := container.NewVBox(
        widget.NewLabelWithStyle(warning.Title,
            fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
        widget.NewLabel(warning.Message),
        ackCheckbox,
    )

    dialog := dialog.NewCustom(
        "Utilities Warning",
        "Close",
        content,
        s.window,
    )

    dialog.SetButtons([]widget.Button{
        {Text: "Go Back", OnTapped: func() { dialog.Hide() }},
        {Text: "Continue Anyway", OnTapped: func() {
            if ackCheckbox.Checked {
                s.trade.Sector = sector.Name
                s.trade.UtilitiesWarningAcknowledged = true
                dialog.Hide()
                s.navigator.Next()
            } else {
                dialog.ShowInformation("Required",
                    "You must acknowledge the warning", s.window)
            }
        }},
    })

    dialog.Show()
}
```

### 6. Telemetry Logging

**Status:** ⏳ PENDING

**Requirements:**
Add logging in ticker_entry.go and sector_selection.go:

```go
// When yellow/red strategy selected
if suitability.RequireAcknowledgement {
    log.Warn("strategy_warning_displayed",
        "sector", s.trade.Sector,
        "ticker", s.trade.Ticker,
        "strategy", stratID,
        "suitability", suitability.Rating,
    )
}

// When user acknowledges warning
if s.ackCheckbox.Checked {
    log.Warn("strategy_warning_acknowledged",
        "sector", s.trade.Sector,
        "ticker", s.trade.Ticker,
        "strategy", s.trade.Strategy,
        "suitability", s.trade.StrategySuitability,
    )
    s.trade.StrategyWarningAcknowledged = true
}

// When Utilities selected
log.Warn("utilities_sector_selected",
    "ticker", s.trade.Ticker,
    "strategy", s.trade.Strategy,
    "acknowledged", s.trade.UtilitiesWarningAcknowledged,
)
```

### 7. Unit Tests

**Status:** ⏳ PENDING

**Test Coverage Needed:**
- `policy_test.go` - Test strategy_suitability parsing
- `trade_test.go` - Test acknowledgement field persistence
- `ticker_entry_test.go` - Test color indicator display
- `sector_selection_test.go` - Test Utilities modal behavior

**Example Test:**
```go
func TestStrategySuitability_AllRatings(t *testing.T) {
    policy := loadTestPolicy()
    healthcare := policy.Sectors[0]

    // Test green strategy
    assert.Equal(t, "excellent", healthcare.StrategySuitability["Alt10"].Rating)
    assert.Equal(t, "green", healthcare.StrategySuitability["Alt10"].Color)
    assert.False(t, healthcare.StrategySuitability["Alt10"].RequireAcknowledgement)

    // Test yellow strategy
    assert.Equal(t, "marginal", healthcare.StrategySuitability["Alt26"].Rating)
    assert.Equal(t, "yellow", healthcare.StrategySuitability["Alt26"].Color)
    assert.True(t, healthcare.StrategySuitability["Alt26"].RequireAcknowledgement)

    // Test red strategy
    assert.Equal(t, "incompatible", healthcare.StrategySuitability["Alt22"].Rating)
    assert.Equal(t, "red", healthcare.StrategySuitability["Alt22"].Color)
    assert.True(t, healthcare.StrategySuitability["Alt22"].RequireAcknowledgement)
}

func TestUtilitiesWarning_Exists(t *testing.T) {
    policy := loadTestPolicy()
    utilities := findSector(policy, "Utilities")

    assert.NotNil(t, utilities.UtilitiesWarning)
    assert.Contains(t, utilities.UtilitiesWarning.Message, "ZERO successful backtests")
    assert.Equal(t, "⚠️ Utilities Sector Warning", utilities.UtilitiesWarning.Title)
}
```

### 8. Integration Tests

**Status:** ⏳ PENDING

Test full workflow with warning overrides:
- Select Healthcare → Select incompatible strategy (Alt22) → Verify acknowledgement required
- Select Utilities → Verify modal shows → Acknowledge → Select any strategy → Verify all red
- Track that warning acknowledgements persist in saved trade

### 9. Documentation Updates

**Status:** ⏳ PENDING

**Files to Update:**
1. `CLAUDE.md` - Remove "Relaxing sector guardrails" from non-goals (line 19)
2. `CLAUDE.md` - Update Rule #2 to clarify warnings vs hard blocks
3. `README.md` - Document warning system for users
4. `architects-intent.md` - Note architectural shift rationale

### 10. Manual Testing

**Status:** ⏳ PENDING

**Test Scenarios:**
- [ ] Green strategy (Healthcare + Alt10) - No warning, immediate continue
- [ ] Yellow strategy (Healthcare + Alt26) - Warning banner, checkbox required
- [ ] Red strategy (Healthcare + Alt22) - Strong warning, explicit acknowledgement
- [ ] Utilities sector - Modal shows, must acknowledge to proceed
- [ ] All Utilities strategies show red indicators
- [ ] Go Back from Utilities modal returns to sector selection
- [ ] Warning acknowledgements persist in saved trade data

---

## Risk Mitigation

### Risk 1: Users Ignore Warnings and Lose Money
**Mitigation:** Make warning text explicit and track override outcomes

### Risk 2: Color-blind Accessibility
**Mitigation:** Use icons (✓ ⚠ ✗) in addition to colors

### Risk 3: Warning Fatigue
**Mitigation:** Only show warnings for yellow/red (not green)

---

## Next Steps

1. ✅ **Verify Data Models** - Compile project to ensure no syntax errors
2. ✅ **Implement ticker_entry.go** - Color-coded strategy display
3. ✅ **Implement sector_selection.go** - Utilities warning modal
4. ⏳ **Add Telemetry** - Log warning overrides (basic console logging added, can enhance later)
5. ⏳ **Write Tests** - Unit and integration tests
6. ⏳ **Update Documentation** - CLAUDE.md, README.md
7. ⏳ **Manual Testing** - Full workflow validation

---

## Estimated Remaining Time

- UI Implementation (ticker_entry.go + sector_selection.go): **4-6 hours**
- Telemetry Integration: **1 hour**
- Unit Tests: **2-3 hours**
- Integration Tests: **1-2 hours**
- Documentation: **1 hour**
- Manual Testing: **1-2 hours**

**Total:** 10-15 hours remaining

---

## Exit Criteria

Phase 6 is complete when:

- [x] All sectors have strategy_suitability ratings in policy.json
- [x] Utilities has utilities_warning configuration
- [x] Data models updated (policy.go, trade.go)
- [x] ticker_entry.go shows ALL strategies with color indicators
- [x] Acknowledgement checkboxes appear for yellow/red strategies
- [x] Utilities modal blocks without acknowledgement
- [x] Basic telemetry logging (console output for warnings)
- [ ] Unit tests pass (100% of Gherkin scenarios)
- [ ] Integration tests pass (warning override workflows)
- [ ] Manual testing completed (all scenarios)
- [ ] Documentation updated (CLAUDE.md non-goals)

---

**For Questions or Issues:** Reference `plans/roadmap.md` Phase 6 section (lines 2560-3199)
