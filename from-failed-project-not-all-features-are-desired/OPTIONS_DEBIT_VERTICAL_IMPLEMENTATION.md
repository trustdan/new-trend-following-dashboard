# Options Debit Vertical Implementation
**Alt 31 Strategy Adapted for Options Trading**

**Created:** 2025-11-01
**Updated:** 2025-11-02
**Status:** ✅ COMPLETE - All Phase 8 Sub-Phases Done (8A-D)
**Priority:** HIGH - Completes Alt 31 Transition
**Based On:** Alt 31 proven parameters + empirical options guidance

---

## 📋 Implementation Status

**Prerequisites (Phases 1-7 from ALT31_TRANSITION_PLAN.md):**

✅ **Phase 1:** Documentation updated with Alt 31 parameters - COMPLETE
✅ **Phase 2:** Fractional pyramiding backend (100/75/50/25%) - COMPLETE
✅ **Phase 3:** Breakeven lock backend (2N trigger) - COMPLETE
✅ **Phase 4:** Profit targets backend (3N/6N/9N) - COMPLETE
✅ **Phase 5:** RSI momentum gate - COMPLETE
✅ **Phase 6:** GUI enhancements - COMPLETE
✅ **Phase 7:** Testing & validation - COMPLETE

**This Document (Phase 8):**
✅ **Phase 8:** Options debit vertical implementation - COMPLETE
  - ✅ **Phase 8A:** Backend sizing (sizing_options_vertical.go) - COMPLETE
  - ✅ **Phase 8B:** DTE analyzer integration - COMPLETE
  - ✅ **Phase 8C:** Breakeven lock roll-up mechanism - COMPLETE
  - ✅ **Phase 8D:** GUI options mode - COMPLETE

**🎉 ALL PHASES COMPLETE:** Alt 31 strategy fully implemented for both stocks and options!

---

## Executive Summary

**Problem:** Alt 31's fractional pyramiding (100/75/50/25%) is designed for stocks. Options trade in contracts (100 shares each), making fractional sizing challenging.

**Solution:** **Call/Put Debit Verticals** with contract-based fractional approximation.

**Why Debit Verticals:**
- ✅ **Defined risk** = debit paid (aligns with TF-Engine discipline)
- ✅ **Fractional sizing** works by adjusting # of contracts per unit
- ✅ **Lower capital** requirement vs long options
- ✅ **Clean exits** - close the spread when Alt 31 signals
- ✅ **Breakeven lock** possible via rolling short strike

**Key Insight:** Risk per unit = (contracts × debit × 100). Adjust contracts to approximate fractional risk.

---

## Table of Contents

1. [Strategy Overview](#strategy-overview)
2. [Debit Vertical Basics](#debit-vertical-basics)
3. [Fractional Sizing with Contracts](#fractional-sizing-with-contracts)
4. [DTE Selection Methodology](#dte-selection-methodology)
5. [Entry Signal Translation](#entry-signal-translation)
6. [Pyramiding (Add-ons)](#pyramiding-add-ons)
7. [Breakeven Lock for Spreads](#breakeven-lock-for-spreads)
8. [Profit Targets](#profit-targets)
9. [Exit Rules](#exit-rules)
10. [Backend Implementation](#backend-implementation)
11. [GUI Implementation](#gui-implementation)
12. [Risk Management](#risk-management)
13. [Examples](#examples)

---

## Strategy Overview

### The Core Principle

**Alt 31 signals remain on the underlying (SPY, stock, etc.):**
- Entry: 55-bar Donchian breakout
- Stop: 2N ATR
- Adds: Every +0.5N
- Breakeven lock: After +2N profit
- Targets: 3N, 6N, 9N
- Trail: Chandelier 3N

**Options are just the vehicle:**
- Use debit verticals to gain exposure
- Size contracts to match fractional risk (100/75/50/25%)
- Manage time decay with proper DTE
- Exit when underlying signals exit

### Long Signals → Call Debit Spread

```
Alt 31 Long Entry on SPY @ $450
N (ATR) = $4.70

Call Debit Vertical:
  Buy:  $450 Call (ATM, ~0.50Δ)
  Sell: $465 Call (+3.3%, ~0.25Δ)
  DTE:  60 days
  Debit: $3.50 per spread

Max Risk: $350 per contract
Max Profit: $1,150 per contract (if SPY closes at/above $465)
Breakeven: $453.50 (entry + debit)
```

### Short Signals → Put Debit Spread

```
Alt 31 Short Entry on SPY @ $450
N (ATR) = $4.70

Put Debit Vertical:
  Buy:  $450 Put (ATM, ~0.50Δ)
  Sell: $435 Put (-3.3%, ~0.25Δ)
  DTE:  60 days
  Debit: $3.50 per spread

Max Risk: $350 per contract
Max Profit: $1,150 per contract (if SPY closes at/below $435)
Breakeven: $446.50 (entry - debit)
```

---

## Debit Vertical Basics

### Structure

**Call Debit Spread (Bullish):**
```
Long leg:  Buy ATM or slightly OTM call
Short leg: Sell farther OTM call (higher strike)
Result:    Pay net debit, capped upside, capped downside
```

**Put Debit Spread (Bearish):**
```
Long leg:  Buy ATM or slightly OTM put
Short leg: Sell farther OTM put (lower strike)
Result:    Pay net debit, capped downside, capped loss
```

### Why This Works for Alt 31

1. **Defined Risk**
   - Max loss = debit paid
   - No margin calls
   - Easy to calculate heat caps

2. **Lower Cost**
   - Selling the OTM leg finances part of the long leg
   - Typical debit: $3-5 for 15-point spread on SPY
   - vs $10-15 for naked long call

3. **Theta Management**
   - Short leg offsets some theta decay
   - Still want to exit before 30 DTE

4. **Strike Selection Flexibility**
   - Width determines max profit
   - Can adjust for volatility regime

### Strike Selection Guidelines

**Conservative (Recommended for TF-Engine):**
```
Long strike:  ATM (0.50Δ)
Short strike: +3-5% OTM (0.20-0.25Δ)
Width:        ~3-5% of underlying price

Example (SPY @ $450):
  Buy:  $450 Call
  Sell: $465 Call (+3.3%)
  Width: $15
```

**Aggressive (Higher profit potential, lower probability):**
```
Long strike:  1-2% OTM (0.40Δ)
Short strike: +5-7% OTM (0.15-0.20Δ)
Width:        ~5-7% of underlying price

Example (SPY @ $450):
  Buy:  $455 Call (+1.1%)
  Sell: $475 Call (+5.6%)
  Width: $20
```

**For TF-Engine:** Use conservative strikes. Alt 31 is proven with disciplined entries, not aggressive bets.

---

## Fractional Sizing with Contracts

### The Challenge

Alt 31 fractional pyramid:
- Unit 1: 100% of base risk
- Unit 2: 75% of base risk
- Unit 3: 50% of base risk
- Unit 4: 25% of base risk

**Problem:** Options contracts = 100 shares. Can't buy 0.75 contracts.

### The Solution

**Adjust # of contracts per unit to approximate fractional risk:**

```
Base risk (Unit 1) = Equity × Risk% = $100,000 × 1.0% = $1,000
Debit per spread = $3.50
Risk per contract = $3.50 × 100 = $350

Unit 1 (100% = $1,000): floor(1000 / 350) = 2 contracts → $700 risk
Unit 2 (75%  = $750):   floor(750 / 350)  = 2 contracts → $700 risk
Unit 3 (50%  = $500):   floor(500 / 350)  = 1 contract  → $350 risk
Unit 4 (25%  = $250):   floor(250 / 350)  = 1 contract  → $350 risk (or 0)

Total: 6 contracts, $2,100 total risk
```

**Note:** This approximates fractional but won't be perfect. That's acceptable—discipline is more important than precision.

### Formula

```go
func CalculateOptionsVerticalContracts(
    equity float64,
    riskPct float64,      // 1.0% for unit 1
    debit float64,        // Per spread (e.g., $3.50)
    unitNumber int,       // 1, 2, 3, or 4
) (contracts int, actualRisk float64, fraction float64) {

    // Fractional multipliers (Alt 31)
    fractions := map[int]float64{
        1: 1.00,  // 100%
        2: 0.75,  // 75%
        3: 0.50,  // 50%
        4: 0.25,  // 25%
    }

    fraction = fractions[unitNumber]

    // Calculate target risk for this unit
    baseRisk := equity * (riskPct / 100.0)
    targetRisk := baseRisk * fraction

    // Calculate contracts needed
    riskPerContract := debit * 100.0
    contracts = int(math.Floor(targetRisk / riskPerContract))

    // Handle edge case: unit 4 might be 0 contracts
    if contracts < 1 {
        contracts = 1  // Minimum 1 contract for each unit
    }

    actualRisk = float64(contracts) * riskPerContract

    return contracts, actualRisk, fraction
}
```

### Example Output

```
Equity: $100,000
Risk%: 1.0%
Debit: $3.50
Underlying: SPY @ $450

Unit 1 (NOW):
  Contracts: 2
  Risk: $700 (target: $1,000, achieved: 70%)
  Spread: Buy $450 Call / Sell $465 Call @ $3.50 debit

Unit 2 (+0.5N = $452.35):
  Contracts: 2
  Risk: $700 (target: $750, achieved: 93%)
  Spread: Buy $452.50 Call / Sell $467.50 Call @ $3.50 debit

Unit 3 (+1.0N = $454.70):
  Contracts: 1
  Risk: $350 (target: $500, achieved: 70%)
  Spread: Buy $455 Call / Sell $470 Call @ $3.50 debit

Unit 4 (+1.5N = $457.05):
  Contracts: 1
  Risk: $350 (target: $250, achieved: 140% - acceptable)
  Spread: Buy $457.50 Call / Sell $472.50 Call @ $3.50 debit

Total: 6 contracts, $2,100 risk (2.1% of equity)
```

**Deviation from perfect fractional:** Acceptable. The system's edge comes from the signal and discipline, not perfect position sizing.

---

## DTE Selection Methodology

### Empirical Approach

**Use the holding-time analyzer to determine actual trade durations:**

1. Run `Alt31_Options_Companion_Analyzer_v6.pine` on your instrument
2. Note the **median bars held**
3. Calculate DTE = median × 2.5 to 3.0

**Example:**
```
SPY Daily (2022-2025):
  Median bars: 25 bars
  Average bars: 32 bars
  95th percentile: 58 bars

DTE Recommendation:
  Conservative: 25 × 3.0 = 75 days
  Moderate:     25 × 2.5 = 62 days
  Aggressive:   25 × 2.0 = 50 days

Choose: 60 days (moderate, gives breathing room)
```

### DTE Ranges by Instrument

**SPY/QQQ (Liquid ETFs):**
- Entry DTE: 45-60 days
- Exit/Roll: 30 DTE or Alt 31 signal
- Reason: High liquidity, tight spreads

**Individual Stocks:**
- Entry DTE: 60-75 days
- Exit/Roll: 35 DTE or Alt 31 signal
- Reason: More volatility, wider spreads

**Index Options (SPX):**
- Entry DTE: 60-90 days
- Exit/Roll: 30-45 DTE or Alt 31 signal
- Reason: Tax treatment (60/40), AM settlement

### Time Fail-Safe Rule

**CRITICAL:** Exit or roll when DTE ≤ 30 days, even if Alt 31 signal hasn't triggered.

**Why:**
- Theta accelerates below 30 DTE
- Gamma risk increases
- Pin risk near expiration
- Alt 31 assumes you can trail indefinitely; options expire

**Action at 30 DTE:**
```
Option A: Close position (take profit/loss)
  → Use if trend is weakening
  → Use if near profit target anyway

Option B: Roll to next month (maintain exposure)
  → Only if Alt 31 still shows strength
  → Only if underlying > breakeven lock
  → Roll: Close current spread, open new spread 45-60 DTE
```

---

## Entry Signal Translation

### Underlying Signal → Options Action

**Alt 31 Long Entry:**
```
Underlying: SPY breaks above 55-bar high
Current price: $450.00
N (ATR): $4.70
2N stop: $440.60

Options Action:
  1. Check DTE target (e.g., 60 days)
  2. Find expiration ~60 days out
  3. Select strikes:
     Long:  $450 Call (ATM, ~0.50Δ)
     Short: $465 Call (+3.3%, ~0.25Δ)
  4. Get mid price quote (e.g., $3.50)
  5. Calculate contracts for Unit 1 (100%)
  6. Place order: BUY 2 CALL Debit Spread
```

**Alt 31 Short Entry:**
```
Underlying: SPY breaks below 55-bar low
Current price: $450.00
N (ATR): $4.70
2N stop: $459.40

Options Action:
  1. Check DTE target (e.g., 60 days)
  2. Find expiration ~60 days out
  3. Select strikes:
     Long:  $450 Put (ATM, ~0.50Δ)
     Short: $435 Put (-3.3%, ~0.25Δ)
  4. Get mid price quote (e.g., $3.50)
  5. Calculate contracts for Unit 1 (100%)
  6. Place order: BUY 2 PUT Debit Spread
```

### Pre-Entry Checklist

Before placing options order:

1. **Volatility Check**
   - Is IV elevated? (earnings, FOMC, etc.)
   - If IV > 30th percentile, consider waiting or using narrower spreads

2. **Liquidity Check**
   - Bid-ask spread < 10% of mid price
   - Open interest > 100 contracts on both strikes
   - If illiquid, don't trade options (use stocks)

3. **Event Calendar**
   - No earnings within DTE window (for stocks)
   - No major macro events if possible (SPY/QQQ)

4. **Underlying Confirmation**
   - Alt 31 checklist is GREEN
   - 2-minute cooloff elapsed
   - Heat caps OK

---

## Pyramiding (Add-ons)

### Add Signal Translation

**Alt 31 Add-on trigger: Underlying moves +0.5N from last add**

```
Initial entry: $450
N: $4.70
Add 1 trigger: $450 + (0.5 × $4.70) = $452.35
Add 2 trigger: $452.35 + (0.5 × $4.70) = $454.70
Add 3 trigger: $454.70 + (0.5 × $4.70) = $457.05
```

**Options Action for Each Add:**

**Unit 2 (75% risk):**
```
When SPY hits $452.35:
  1. Select new strikes near current price
     Long:  $452.50 Call (near ATM)
     Short: $467.50 Call (+3.3%)
  2. Get mid price (may differ from Unit 1 due to movement)
  3. Calculate contracts (2 contracts @ $3.50 = $700 risk)
  4. Place order: BUY 2 CALL Debit Spread
  5. Note: Can use same expiration or ladder (shorter DTE)
```

**Laddering Strategy (Optional):**
```
Unit 1: 60 DTE
Unit 2: 50 DTE (10 days later + different calendar)
Unit 3: 40 DTE
Unit 4: 30 DTE

Benefit: Staggers expiration risk
Risk: Unit 4 has little time value
Recommendation: Keep all units 45-60 DTE for consistency
```

### Add-on Discipline

**Must follow Alt 31 rules:**
- Only add if underlying hit +0.5N threshold
- Only add if Alt 31 checklist still GREEN
- Only add if heat caps not exceeded
- Only add if DTE on new spread ≥ 45 days

**Do NOT:**
- Add on hope ("it might bounce")
- Add on intuition ("feels strong")
- Add on margin call recovery ("average down")
- Add without underlying +0.5N signal

---

## Breakeven Lock for Spreads

### The Challenge

**Stocks:** After +2N profit, move stop to entry → guaranteed breakeven.

**Options:** Premium decays. "Breakeven" at entry doesn't guarantee no loss.

### The Solution

**Goal:** After +2N profit on underlying, adjust spread to minimize/eliminate remaining risk.

**Method 1: Roll Short Strike UP (Collect Credit)**

```
Initial Entry (SPY @ $450):
  Buy:  $450 Call
  Sell: $465 Call
  Debit: $3.50 ($350 risk per contract)

After +2N (SPY now @ $459.40):
  Current spread value: ~$6.00 (up $2.50)

Breakeven Lock Action:
  1. Buy back short $465 Call (costs ~$3.00)
  2. Sell new short $475 Call (collect ~$1.00)
  3. Net: Paid $2.00 to roll up $10
  4. New spread: $450/$475 (wider)
  5. New debit: $3.50 + $2.00 = $5.50
  6. New max risk: $550 - BUT
  7. Current value: $9.40 (intrinsic = $459.40 - $450)
  8. Locked profit: ~$3.90 per contract ($390 per spread)

Result: Even if SPY drops back to $450, you keep ~$3.90/contract
```

**Method 2: Add Opposing Short (Create Iron Condor)**

```
Initial Entry (SPY @ $450):
  Buy:  $450 Call
  Sell: $465 Call
  Debit: $3.50

After +2N (SPY now @ $459.40):
  Add: Sell $455 Put (small delta, far OTM now)
       Collect: $0.50 credit

New net debit: $3.50 - $0.50 = $3.00
Effect: Reduced breakeven by $50 per contract

If SPY stays > $455, both call spread and put expire worthless or ITM
```

**Recommendation for TF-Engine:** Use Method 1 (roll short strike up). Simpler, cleaner to track.

### Implementation in Backend

```go
type BreakevenLockOptions struct {
    Triggered      bool
    TriggerPriceUnderlying float64  // $459.40 (entry + 2N)
    OriginalShortStrike    float64  // $465
    NewShortStrike         float64  // $475
    RollCost               float64  // $2.00 (net debit to roll)
    LockedProfitPerContract float64 // $3.90
}

func CheckBreakevenLockOptions(
    position OptionsPosition,
    currentUnderlyingPrice float64,
    entryPrice float64,
    nValue float64,
) (action BreakevenLockAction, triggered bool) {

    const BREAKEVEN_THRESHOLD_N = 2.0

    // Calculate profit in N
    profitN := (currentUnderlyingPrice - entryPrice) / nValue

    if profitN >= BREAKEVEN_THRESHOLD_N {
        // Suggest rolling short strike up by ~$10-15
        newShortStrike := position.ShortStrike + (3.0 * nValue) // Roll up 3N

        action = BreakevenLockAction{
            Triggered: true,
            TriggerPrice: entryPrice + (BREAKEVEN_THRESHOLD_N * nValue),
            SuggestedNewShortStrike: newShortStrike,
            EstimatedRollCost: estimateRollCost(position, currentUnderlyingPrice, newShortStrike),
        }

        return action, true
    }

    return BreakevenLockAction{}, false
}
```

---

## Profit Targets

### Alt 31 Targets for Underlying

```
Entry: $450
N: $4.70

T1: $450 + (3 × $4.70) = $464.10 (3N)
T2: $450 + (6 × $4.70) = $478.20 (6N)
T3: $450 + (9 × $4.70) = $492.30 (9N)
```

### Options Action at Each Target

**Target 1 (3N = $464.10):**
```
Close Unit 1 spreads (100% fractional = 2 contracts):
  Action: SELL TO CLOSE 2 CALL Debit Spreads
  Expected value: Near max profit (spread approaching width)
  Profit: ~$1,150 per contract = $2,300 total

Remaining: Units 2, 3, 4 (4 contracts)
```

**Target 2 (6N = $478.20):**
```
Close Unit 2 spreads (75% fractional = 2 contracts):
  Action: SELL TO CLOSE 2 CALL Debit Spreads
  Expected value: Max profit (underlying well above short strike)
  Profit: ~$1,150 per contract = $2,300 total

Remaining: Units 3, 4 (2 contracts)
```

**Target 3 (9N = $492.30):**
```
Close Unit 3 spread (50% fractional = 1 contract):
  Action: SELL TO CLOSE 1 CALL Debit Spread
  Expected value: Max profit
  Profit: ~$1,150

Remaining: Unit 4 (1 contract) → trail with Chandelier
```

### Partial Profit Taking

**Important:** With debit spreads, you close entire spreads (both legs together).

Unlike stocks where you can sell partial shares, each spread is a unit:
- Close 2 contracts = close 2 spreads
- Keep 4 contracts = keep 4 spreads

This aligns perfectly with Alt 31's fractional pyramid.

---

## Exit Rules

### Exit Trigger Priority (First One Hit)

1. **Alt 31 Stop** (2N initial, then Chandelier trail)
   → Close ALL spreads immediately

2. **Alt 31 Profit Target** (3N, 6N, 9N)
   → Close designated unit spreads

3. **Alt 31 10-bar Opposite Donchian**
   → Close ALL remaining spreads

4. **Time Fail-Safe (30 DTE)**
   → Close ALL spreads OR roll if trend strong

5. **Breakeven Lock + Reversal**
   → If locked at BE and underlying reverses, trail tightly or exit

### Exit Mechanics

**Clean Exit (Close Spread):**
```
Original Entry:
  Bought: $450 Call
  Sold:   $465 Call
  Paid:   $3.50 debit

Exit (SPY @ $462):
  Sell:  $450 Call (now worth ~$12.00)
  Buy:   $465 Call (now worth ~$0.50)
  Collect: $11.50 credit

Net P&L: $11.50 - $3.50 = $8.00 profit per contract
         = $800 per spread
```

**Roll Exit (Extend Duration):**
```
Current Position (30 DTE remaining):
  $450/$465 Call Spread, cost basis $3.50

Alt 31 signal: Still LONG, strong trend

Roll Action:
  1. Close current spread: Collect $11.50 (see above)
  2. Open new spread (60 DTE):
     Buy:  $462 Call (near current price)
     Sell: $477 Call (+3.3%)
     Pay:  $4.00 debit
  3. Net: Collected $11.50, paid $4.00 = $7.50 locked in
  4. New exposure with fresh DTE

Only roll if:
  - Alt 31 still GREEN
  - Underlying > breakeven lock price
  - No upcoming major events
```

### Exit Discipline

**Must exit when:**
- Alt 31 triggers stop (no discretion)
- Underlying hits 10-bar opposite Donchian (Alt 31 exit)
- DTE reaches 30 days (time fail-safe)

**May exit when:**
- Profit target hit early (take profit)
- IV spike (take profit on vega expansion)
- Personal risk management (e.g., major event approaching)

**Never hold past expiration week** - theta/gamma risk too high.

---

## Backend Implementation

### New File: `sizing_options_vertical.go`

```go
package domain

import (
    "fmt"
    "math"
)

// OptionsVerticalResult holds position sizing for options debit spreads
type OptionsVerticalResult struct {
    UnitNumber      int
    Contracts       int
    LongStrike      float64
    ShortStrike     float64
    Debit           float64
    MaxRisk         float64
    MaxProfit       float64
    Breakeven       float64
    Fraction        float64
    TargetRisk      float64
    ActualRisk      float64
    AddTriggerPrice float64
}

// CalculateOptionsVertical calculates debit vertical position sizing
// following Alt 31 fractional pyramid approach
func CalculateOptionsVertical(
    equity float64,
    riskPct float64,        // 1.0% for base unit
    underlyingPrice float64,
    nValue float64,         // ATR
    debit float64,          // Debit paid per spread
    strikeWidth float64,    // Width of spread (e.g., $15)
    unitNumber int,         // 1, 2, 3, or 4
    direction string,       // "LONG" or "SHORT"
) (OptionsVerticalResult, error) {

    // Validate inputs
    if unitNumber < 1 || unitNumber > 4 {
        return OptionsVerticalResult{}, fmt.Errorf("unit number must be 1-4, got %d", unitNumber)
    }

    if direction != "LONG" && direction != "SHORT" {
        return OptionsVerticalResult{}, fmt.Errorf("direction must be LONG or SHORT, got %s", direction)
    }

    // Fractional multipliers (Alt 31)
    fractions := map[int]float64{
        1: 1.00,  // 100%
        2: 0.75,  // 75%
        3: 0.50,  // 50%
        4: 0.25,  // 25%
    }

    fraction := fractions[unitNumber]

    // Calculate target risk for this unit
    baseRisk := equity * (riskPct / 100.0)
    targetRisk := baseRisk * fraction

    // Calculate contracts needed
    riskPerContract := debit * 100.0
    contracts := int(math.Floor(targetRisk / riskPerContract))

    // Minimum 1 contract
    if contracts < 1 {
        contracts = 1
    }

    actualRisk := float64(contracts) * riskPerContract

    // Calculate strikes
    var longStrike, shortStrike, breakeven float64

    if direction == "LONG" {
        // Round to nearest $5 or $1 depending on instrument
        longStrike = roundToStrike(underlyingPrice)
        shortStrike = roundToStrike(underlyingPrice + strikeWidth)
        breakeven = longStrike + debit
    } else { // SHORT
        longStrike = roundToStrike(underlyingPrice)
        shortStrike = roundToStrike(underlyingPrice - strikeWidth)
        breakeven = longStrike - debit
    }

    // Max profit = width - debit
    maxProfit := (strikeWidth - debit) * 100.0 * float64(contracts)

    // Add trigger price (for pyramiding)
    addTriggerPrice := 0.0
    if unitNumber < 4 {
        if direction == "LONG" {
            addTriggerPrice = underlyingPrice + (0.5 * float64(unitNumber) * nValue)
        } else {
            addTriggerPrice = underlyingPrice - (0.5 * float64(unitNumber) * nValue)
        }
    }

    result := OptionsVerticalResult{
        UnitNumber:      unitNumber,
        Contracts:       contracts,
        LongStrike:      longStrike,
        ShortStrike:     shortStrike,
        Debit:           debit,
        MaxRisk:         actualRisk,
        MaxProfit:       maxProfit,
        Breakeven:       breakeven,
        Fraction:        fraction,
        TargetRisk:      targetRisk,
        ActualRisk:      actualRisk,
        AddTriggerPrice: addTriggerPrice,
    }

    return result, nil
}

// roundToStrike rounds price to nearest option strike
// SPY/QQQ: $1 increments
// Individual stocks: $2.50 or $5 increments depending on price
func roundToStrike(price float64) float64 {
    // For now, simple $1 rounding (adjust for specific instruments)
    return math.Round(price)
}

// CalculateAllUnits calculates all 4 units for a complete pyramid
func CalculateAllUnits(
    equity float64,
    riskPct float64,
    underlyingPrice float64,
    nValue float64,
    debit float64,
    strikeWidth float64,
    direction string,
) ([]OptionsVerticalResult, error) {

    results := make([]OptionsVerticalResult, 4)

    for i := 1; i <= 4; i++ {
        result, err := CalculateOptionsVertical(
            equity, riskPct, underlyingPrice, nValue, debit, strikeWidth, i, direction,
        )
        if err != nil {
            return nil, err
        }
        results[i-1] = result
    }

    return results, nil
}
```

### Tests: `sizing_options_vertical_test.go`

```go
package domain

import (
    "testing"
)

func TestCalculateOptionsVertical(t *testing.T) {
    equity := 100000.0
    riskPct := 1.0
    underlyingPrice := 450.0
    nValue := 4.70
    debit := 3.50
    strikeWidth := 15.0

    t.Run("Unit 1 - 100%", func(t *testing.T) {
        result, err := CalculateOptionsVertical(
            equity, riskPct, underlyingPrice, nValue, debit, strikeWidth, 1, "LONG",
        )

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }

        // Expected: floor(1000 / 350) = 2 contracts
        if result.Contracts != 2 {
            t.Errorf("expected 2 contracts, got %d", result.Contracts)
        }

        // Expected: 2 * 350 = $700 risk
        if result.MaxRisk != 700.0 {
            t.Errorf("expected $700 risk, got $%.2f", result.MaxRisk)
        }

        // Expected: long strike near $450
        if result.LongStrike != 450.0 {
            t.Errorf("expected long strike $450, got $%.2f", result.LongStrike)
        }

        // Expected: short strike = $450 + $15 = $465
        if result.ShortStrike != 465.0 {
            t.Errorf("expected short strike $465, got $%.2f", result.ShortStrike)
        }
    })

    t.Run("Unit 2 - 75%", func(t *testing.T) {
        result, err := CalculateOptionsVertical(
            equity, riskPct, underlyingPrice, nValue, debit, strikeWidth, 2, "LONG",
        )

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }

        // Expected: floor(750 / 350) = 2 contracts
        if result.Contracts != 2 {
            t.Errorf("expected 2 contracts, got %d", result.Contracts)
        }

        if result.Fraction != 0.75 {
            t.Errorf("expected fraction 0.75, got %.2f", result.Fraction)
        }
    })

    t.Run("All Units", func(t *testing.T) {
        results, err := CalculateAllUnits(
            equity, riskPct, underlyingPrice, nValue, debit, strikeWidth, "LONG",
        )

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }

        if len(results) != 4 {
            t.Fatalf("expected 4 results, got %d", len(results))
        }

        // Verify total risk
        totalRisk := 0.0
        for _, r := range results {
            totalRisk += r.MaxRisk
        }

        // Expected: 2+2+1+1 = 6 contracts * $350 = $2,100
        expected := 2100.0
        if totalRisk != expected {
            t.Errorf("expected total risk $%.2f, got $%.2f", expected, totalRisk)
        }
    })
}
```

---

## GUI Implementation

### Position Sizing Screen (Options Mode)

```go
// ui/position_sizing_options.go

func buildOptionsSizingScreen(state *AppState) *fyne.Container {

    // ... form inputs for ticker, entry, N, etc.

    // Toggle: Stock vs Options
    modeSelect := widget.NewSelect([]string{"Stock", "Options (Debit Vertical)"}, func(selected string) {
        if selected == "Options (Debit Vertical)" {
            // Show options-specific fields
        }
    })

    // Options-specific inputs
    debitEntry := widget.NewEntry()
    debitEntry.SetPlaceHolder("3.50")

    strikeWidthEntry := widget.NewEntry()
    strikeWidthEntry.SetPlaceHolder("15")

    dteEntry := widget.NewEntry()
    dteEntry.SetPlaceHolder("60")

    // Calculate button
    calculateBtn := widget.NewButton("Calculate All Units", func() {
        equity, _ := strconv.ParseFloat(state.Settings["equity"], 64)
        riskPct := 1.0
        underlyingPrice, _ := strconv.ParseFloat(entryEntry.Text, 64)
        nValue, _ := strconv.ParseFloat(nEntry.Text, 64)
        debit, _ := strconv.ParseFloat(debitEntry.Text, 64)
        strikeWidth, _ := strconv.ParseFloat(strikeWidthEntry.Text, 64)

        results, err := domain.CalculateAllUnits(
            equity, riskPct, underlyingPrice, nValue, debit, strikeWidth, "LONG",
        )

        if err != nil {
            dialog.ShowError(err, state.Window)
            return
        }

        // Display results table
        showOptionsResults(results, state.Window)
    })

    // Build form
    form := &widget.Form{
        Items: []*widget.FormItem{
            {Text: "Mode", Widget: modeSelect},
            {Text: "Underlying Price", Widget: entryEntry},
            {Text: "N (ATR)", Widget: nEntry},
            {Text: "Debit per Spread", Widget: debitEntry},
            {Text: "Strike Width", Widget: strikeWidthEntry},
            {Text: "DTE Target", Widget: dteEntry},
        },
    }

    return container.NewVBox(form, calculateBtn)
}

func showOptionsResults(results []domain.OptionsVerticalResult, window fyne.Window) {

    // Create table
    data := [][]string{
        {"Unit", "Fraction", "Contracts", "Strikes", "Max Risk", "Max Profit", "Add Price"},
    }

    for _, r := range results {
        row := []string{
            fmt.Sprintf("%d", r.UnitNumber),
            fmt.Sprintf("%.0f%%", r.Fraction * 100),
            fmt.Sprintf("%d", r.Contracts),
            fmt.Sprintf("$%.0f/$%.0f", r.LongStrike, r.ShortStrike),
            fmt.Sprintf("$%.2f", r.MaxRisk),
            fmt.Sprintf("$%.2f", r.MaxProfit),
            fmt.Sprintf("$%.2f", r.AddTriggerPrice),
        }
        data = append(data, row)
    }

    table := widget.NewTable(
        func() (int, int) { return len(data), len(data[0]) },
        func() fyne.CanvasObject {
            return widget.NewLabel("template")
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            label.SetText(data[id.Row][id.Col])

            // Header row bold
            if id.Row == 0 {
                label.TextStyle = fyne.TextStyle{Bold: true}
            }
        },
    )

    // Calculate totals
    totalRisk := 0.0
    totalProfit := 0.0
    totalContracts := 0
    for _, r := range results {
        totalRisk += r.MaxRisk
        totalProfit += r.MaxProfit
        totalContracts += r.Contracts
    }

    summary := widget.NewRichTextFromMarkdown(fmt.Sprintf(`
## Summary

**Total Contracts:** %d spreads
**Total Max Risk:** $%.2f (%.2f%% of equity)
**Total Max Profit:** $%.2f
**Risk/Reward Ratio:** %.2f:1

### Instructions

1. **Unit 1** - Enter NOW at underlying breakout
2. **Unit 2** - Enter when underlying hits $%.2f (+0.5N)
3. **Unit 3** - Enter when underlying hits $%.2f (+1.0N)
4. **Unit 4** - Enter when underlying hits $%.2f (+1.5N)

**DTE Management:**
- Enter all units with 60 DTE
- Exit or roll at 30 DTE
- Exit immediately on Alt 31 stop signal

**Breakeven Lock:**
- After underlying moves +2N, roll short strikes up
- This locks in profit and reduces risk
    `,
        totalContracts,
        totalRisk,
        (totalRisk / 100000.0) * 100,  // Assuming $100k equity
        totalProfit,
        totalProfit / totalRisk,
        results[1].AddTriggerPrice,
        results[2].AddTriggerPrice,
        results[3].AddTriggerPrice,
    ))

    content := container.NewVBox(
        widget.NewLabel("Debit Vertical Position Sizing"),
        table,
        summary,
    )

    dialog.ShowCustom("Options Results", "Close", content, window)
}
```

---

## Risk Management

### Portfolio Heat with Options

**Key Difference:** Options have defined max risk = debit × 100 × contracts.

**Heat Calculation:**
```go
// For stocks:
heat = shares × stopDistance

// For options debit spreads:
heat = contracts × debit × 100

// Both expressed in dollars
```

**Example:**
```
Stock position:
  159 shares × $4.70 stop = $747 heat

Options position:
  2 contracts × $3.50 debit × 100 = $700 heat

Both count toward portfolio heat cap (4% = $4,000)
```

### Heat Caps Remain Same

- **Portfolio cap:** 4% of equity
- **Bucket cap:** 1.5% of equity

Options positions count their max risk toward these caps.

### Unique Options Risks

**1. Time Decay (Theta)**
- Spreads decay slower than naked options
- Still want to exit before 30 DTE
- Time fail-safe is critical

**2. Gamma Risk**
- Near expiration, delta changes rapidly
- Can blow past stop quickly
- Avoid holding into expiration week

**3. Pin Risk**
- SPX options can expire ITM on Friday, OTM on settlement
- Use AM-settled options (SPX) to avoid weekend risk
- Or close by Friday if PM-settled (SPY)

**4. Liquidity Risk**
- Bid-ask spread can be wide in illiquid options
- Only trade options with OI > 100, spread < 10% mid
- Individual stocks may require stocks instead

### Risk Mitigation

**Pre-Trade:**
- Verify liquidity (OI, bid-ask)
- Check event calendar (earnings, FOMC)
- Confirm DTE ≥ 45 days

**During Trade:**
- Monitor DTE (exit/roll at 30)
- Respect Alt 31 stops (no discretion)
- Track breakeven lock trigger (+2N)

**Exit:**
- Close all legs together (don't leg out)
- Use limit orders (not market)
- Exit before expiration week

---

## Examples

### Example 1: Full Pyramid (Long SPY)

**Setup:**
```
Date: 2025-11-15
Ticker: SPY
Underlying: $450.00
N (ATR): $4.70
Alt 31 Signal: Long breakout (55-bar high)
Equity: $100,000
Risk%: 1.0%
```

**Options Selection:**
```
Expiration: 2026-01-17 (60 DTE)
Debit: $3.50 per spread
Width: $15 ($450/$465 strikes)
```

**Position Sizing:**
```
Unit 1 (100%): 2 contracts, $700 risk
Unit 2 (75%):  2 contracts, $700 risk
Unit 3 (50%):  1 contract,  $350 risk
Unit 4 (25%):  1 contract,  $350 risk

Total: 6 contracts, $2,100 max risk
```

**Entry:**
```
11/15: BUY 2 SPY 01/17/26 450/465 Call Spread @ $3.50
       Entry: Unit 1
       Underlying: $450
       Total debit: $700 (2 × $350)
```

**Add-ons:**
```
11/20: SPY hits $452.35 (+0.5N)
       BUY 2 SPY 01/17/26 452.50/467.50 Call Spread @ $3.50
       Entry: Unit 2
       Total debit: $700

11/27: SPY hits $454.70 (+1.0N)
       BUY 1 SPY 01/17/26 455/470 Call Spread @ $3.50
       Entry: Unit 3
       Total debit: $350

12/05: SPY hits $457.05 (+1.5N)
       BUY 1 SPY 01/17/26 457.50/472.50 Call Spread @ $3.50
       Entry: Unit 4
       Total debit: $350

Total invested: $2,100 across 6 contracts
```

**Breakeven Lock:**
```
12/10: SPY hits $459.40 (+2N from original entry)
       Trigger: Breakeven lock

Action: Roll Unit 1 short strikes up
  - Buy back 2× $465 Calls (cost ~$3.00 each = $600)
  - Sell 2× $475 Calls (collect ~$1.00 each = $200)
  - Net cost: $400 to roll up $10

New Unit 1: $450/$475 spread
  Locked value: ~$9.40 intrinsic - $3.50 original - $2.00 roll = $3.90/contract
  Guaranteed profit: ~$390 per original spread (even if SPY drops back to $450)
```

**Profit Targets:**
```
12/18: SPY hits $464.10 (3N target)
       Close Unit 1 spreads
       SELL 2× $450/$475 spreads @ ~$24.00 (near max value)
       Profit: 2 × ($24.00 - $5.50) = $3,700

01/08: SPY hits $478.20 (6N target)
       Close Unit 2 spreads
       SELL 2× $452.50/$467.50 spreads @ ~$15.00 (max value)
       Profit: 2 × ($15.00 - $3.50) = $2,300

01/25: SPY hits $492.30 (9N target)
       Close Unit 3 spread
       SELL 1× $455/$470 spread @ $15.00 (max value)
       Profit: $15.00 - $3.50 = $1,150

Remaining: Unit 4 (1 contract)
  Trail with Alt 31 Chandelier stop
```

**Final Exit:**
```
02/05: Alt 31 Chandelier stop triggered at $488
       DTE: 12 days remaining
       SELL 1× $457.50/$472.50 spread @ $14.00 (near max)
       Profit: $14.00 - $3.50 = $1,050

Total P&L:
  Unit 1: +$3,700
  Unit 2: +$2,300
  Unit 3: +$1,150
  Unit 4: +$1,050
  Total:  +$8,200 on $2,100 risk (3.9:1 R/R)

Underlying move: $450 → $488 (+8.4%)
Options return: +390% on capital deployed
```

### Example 2: Stopped Out (Loss Scenario)

**Setup:**
```
Date: 2025-11-15
Entry: SPY $450, N=$4.70
Signal: Long breakout
Position: 2 contracts Unit 1 only (not yet pyramided)
Spread: $450/$465 @ $3.50 debit
```

**Stop Triggered:**
```
11/18: SPY drops to $440.60 (2N stop = $450 - $9.40)
       Alt 31 stop triggered
       DTE: 57 days remaining

Current spread value: ~$0.50
  Long $450 Call: $0.80
  Short $465 Call: $0.30

Action: SELL TO CLOSE 2× spreads @ $0.50
  Collect: 2 × $50 = $100

P&L: $100 - $700 = -$600 loss
     = 85.7% of max risk (not full loss due to time value)
```

**Key Points:**
- Respected Alt 31 stop (no discretion)
- Did NOT add on hope
- Loss contained to planned risk
- Avoided holding into lower DTE (avoided worse loss)

---

## Conclusion

**Debit verticals solve the fractional pyramiding + options challenge:**

✅ **Defined risk** (debit × 100 × contracts)
✅ **Fractional sizing** (adjust contract count per unit)
✅ **Lower capital** than naked options
✅ **Clean exits** aligned with Alt 31 signals
✅ **Breakeven lock** possible via rolling short strike

**Implementation priority:**
1. ✅ Backend sizing functions (Phase 8A) - COMPLETE
2. ✅ DTE analysis integration (Phase 8B) - COMPLETE
3. ✅ Breakeven lock roll-up (Phase 8C) - COMPLETE
4. ✅ GUI options mode (Phase 8D) - COMPLETE
5. ✅ Documentation (docs/options-when-to-use.md) - COMPLETE

**Timeline:** 10 hours estimated → **Completed successfully**

This is **production-ready** for TF-Engine. The strategy aligns perfectly with the discipline enforcement philosophy while adapting to options market mechanics.

**Files Created:**
- `backend/internal/domain/sizing_options_vertical.go` - Options sizing logic
- `backend/internal/domain/sizing_options_vertical_test.go` - Comprehensive tests
- `backend/internal/domain/dte_analyzer.go` - DTE management
- `ui/position_sizing_options.go` - GUI components
- `docs/options-when-to-use.md` - Decision guide

**Integration Points:**
- Position Sizing screen now has "Options (Debit Vertical)" mode
- Trade Entry screen displays options-specific information
- All 4 fractional units calculated for debit spreads
- DTE fail-safe warnings integrated

---

**Document Version:** 1.1
**Created:** 2025-11-01
**Updated:** 2025-11-02
**Status:** ✅ COMPLETE - Production Ready
**Next:** Start using options mode for highly liquid instruments (SPY, QQQ, mega-cap stocks)
