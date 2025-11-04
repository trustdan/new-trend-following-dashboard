# Bearish Screener Quality Evaluation

## Executive Summary

The bearish screeners have **logical filter inversions** but contain **critical conceptual flaws** that make them less reliable than the bullish screeners. Grade: **C+ (needs significant revision)**

---

## Screener-by-Screener Analysis

### 1. Universe_Bearish ❌ **MAJOR FLAW**

**Current Logic:**
- Price below SMA200 (bearish trend) ✓
- Positive EPS growth (fa_epsyoy_pos) ❌
- Positive 5-year sales growth (fa_sales5years_pos) ❌
- Positive ROE (fa_roe_pos) ❌

**Problem:**
This creates a **"value trap" screener**, not a bearish momentum screener. You're looking for fundamentally strong companies in downtrends, which are more likely to:
- Reverse and recover (bad for bearish trades)
- Be temporarily oversold quality names
- Whipsaw bearish positions

**Fix Required:**
Remove positive fundamental filters OR replace with deteriorating fundamentals:
```
Remove: fa_epsyoy_pos, fa_sales5years_pos, fa_roe_pos
Consider adding: fa_epsqoq_neg (negative quarterly EPS growth)
```

**Severity:** 🔴 Critical - This undermines the entire bearish universe concept

---

### 2. Bounce_Bearish ✅ **GOOD LOGIC**

**Current Logic:**
- Price below SMA200 (established downtrend) ✓
- Price above SMA50 (temporary strength) ✓
- RSI > 60 (overbought in downtrend) ✓

**Assessment:**
This is a **bear flag / failed rally setup**. Technically sound for:
- Shorting bounces in downtrends
- Put spreads on temporary strength
- Counter-trend fades

**Improvement Suggestions:**
- Consider adding volume filter for conviction (sh_relvol_o2)
- Consider momentum filter (ta_change_d for daily weakness)

**Grade:** ✅ B+ (solid concept, minor improvements possible)

---

### 3. Breakdown_Bearish ⚠️ **MIXED QUALITY**

**Current Logic:**
- Price below SMA200 (downtrend) ✓
- 52-week low (ta_highlow52w_nl) ✓

**Assessment:**
Technically inverted correctly, but **52-week lows have different dynamics than 52-week highs**:

**Bullish Breakout (52w high):**
- Often continuation patterns
- Buyers in control
- Trend-following friendly

**Bearish Breakdown (52w low):**
- Often capitulation bottoms
- Higher reversal risk
- Value buyers step in
- "Catching a falling knife" problem

**Historical Issue:**
Many legendary losers (Enron, Lehman, etc.) made 52-week lows repeatedly before zero. But many quality names (AAPL 2019, MSFT 2022) also hit 52-week lows before massive recoveries.

**Improvement Needed:**
Add filters to distinguish failing companies from temporarily weak quality names:
- High debt ratios (fa_debteq_high)
- Declining revenues (fa_sales5years_neg)
- Negative profit margins

**Grade:** ⚠️ C+ (dangerous without additional filters)

---

### 4. Death_Cross_Bearish ✅ **GOOD LOGIC**

**Current Logic:**
- Price below SMA200 (downtrend) ✓
- SMA50 below SMA200 (ta_sma50_pb200) ✓
- Breaking trendline resistance (ta_pattern_tlresistance) ✓

**Assessment:**
This is the **mirror opposite of golden cross** and is technically sound:
- Confirms downtrend on multiple timeframes
- Trendline resistance adds confluence
- Classic technical setup

**Consideration:**
Death crosses are often **lagging signals** (just like golden crosses). By the time SMA50 crosses below SMA200, much of the move is done.

**Grade:** ✅ B+ (solid but lagging indicator)

---

## Fundamental Problem: Sector Applicability

### Critical Issue ⚠️

The script applies bearish screeners to **ALL sectors**, including:

1. **Utilities** (0% trend-following success)
   - Mean-reverting sector
   - Bearish trend-following likely just as bad as bullish
   - Should utilities even have bearish screeners?

2. **Energy** (mean-reverting, warned sector)
   - Already struggles with bullish trend-following
   - Bearish setups face same whipsaw problems
   - Requires mean-reversion approach instead

3. **Healthcare** (92% bullish success rate)
   - Defensive sector with strong uptrends
   - Bearish opportunities are rare and brief
   - Risk/reward heavily skewed to long side

### Sector-Specific Recommendation:

**Strong Bearish Candidates:**
- Technology (after momentum exhaustion)
- Consumer Discretionary (recession plays)
- Communication Services (momentum reversals)

**Weak Bearish Candidates:**
- Healthcare (defensive, resilient)
- Consumer Defensive (anti-cyclical)
- Utilities (mean-reverting, avoid trend-following entirely)

---

## Data Quality Comparison

### Bullish Screeners (Original)
- Backed by **293 validated backtests**
- Sector-strategy mapping proven
- 99.74% data quality
- Success rates: Healthcare 92%, Technology strong

### Bearish Screeners (New)
- **ZERO backtest validation** ⚠️
- No historical performance data
- Unproven sector applicability
- Success rates: **Unknown**

**Critical Gap:** You're deploying bearish screeners without the same rigorous research foundation that powers the bullish side.

---

## Missing Elements

### 1. Bearish Strategy Mapping ❌

The bullish system knows:
- Alt10 works on Healthcare (+33.13%)
- Alt22 loves QQQ momentum
- Alt26 best for SPY

**But bearish side has:**
- No strategy-to-bearish-setup mapping
- No guidance on which Pine Script strategies work for shorting
- No options structure recommendations (put spreads vs long puts)

### 2. Bearish Risk Considerations ❌

**Bullish trades have unlimited upside, limited downside.**
**Bearish trades have limited upside, unlimited downside.**

The screeners don't account for:
- Short squeeze risk
- Borrowing costs (if shorting)
- Dividend capture risk
- "Stocks take stairs up, elevators down" (fast reversals)

### 3. Market Regime Context ❌

Bearish screening should consider:
- Is SPY/QQQ in downtrend? (market regime)
- Is VIX elevated? (fear environment)
- Are defensive sectors outperforming? (rotation signal)

Current screeners are "isolated" without market context.

---

## Recommended Fixes

### Priority 1: Fix Universe_Bearish (Critical)

**Remove conflicting fundamental filters:**
```python
# BAD (current)
bearish_urls['universe_bearish'] = bullish_urls['universe'].replace('ta_sma200_pa', 'ta_sma200_pb')
# This keeps: fa_epsyoy_pos, fa_sales5years_pos, fa_roe_pos

# GOOD (proposed)
# Remove ALL fundamental filters from bearish universe, OR
# Add negative fundamental filters: fa_epsyoq_neg, fa_debteq_high
```

### Priority 2: Add Sector-Specific Logic

**Don't create bearish screeners for:**
- Utilities (0% trend-following success)
- Energy (mean-reverting warned sector)
- Healthcare (defensive, 92% bullish success)

**Focus bearish screeners on:**
- Technology (momentum reversals)
- Consumer Discretionary (cyclical weakness)
- Communication Services (GOOGL/META corrections)

### Priority 3: Validate with Backtesting

**Before deploying bearish screeners operationally:**
1. Run Pine Script strategies in "short mode" (if supported)
2. Test death cross signals historically
3. Measure 52-week low reversal rates
4. Compare bearish vs bullish success rates by sector

### Priority 4: Add Bearish Strategy Guidance

**Update policy.json with:**
```json
"bearish_strategies": {
  "death_cross_put_spreads": {
    "best_sectors": ["Technology", "Consumer Discretionary"],
    "avoid_sectors": ["Healthcare", "Utilities"],
    "typical_hold": "2-6 weeks",
    "risk_note": "Use tight stops, bearish moves are fast"
  }
}
```

---

## Finviz Technical Notes

### Filters Used (Validated) ✅

| Filter | Code | Correct? |
|--------|------|----------|
| Below SMA200 | `ta_sma200_pb` | ✅ Yes |
| Above SMA50 | `ta_sma50_pa` | ✅ Yes |
| RSI Overbought | `ta_rsi_ob60` | ✅ Yes |
| 52-week Low | `ta_highlow52w_nl` | ✅ Yes |
| SMA50 below SMA200 | `ta_sma50_pb200` | ✅ Yes |
| Trendline Resistance | `ta_pattern_tlresistance` | ✅ Yes |

**Technical Implementation:** ✅ Correct Finviz syntax

---

## Final Verdict

### Grades by Screener:

1. **Universe_Bearish:** ❌ **D** (fundamental filters conflict)
2. **Bounce_Bearish:** ✅ **B+** (solid bear flag logic)
3. **Breakdown_Bearish:** ⚠️ **C+** (reversal risk, needs filters)
4. **Death_Cross_Bearish:** ✅ **B+** (technically sound)

### Overall Assessment:

**Quality Score: 6/10** (C+ grade)

**Why not higher?**
- Universe_bearish has critical flaw (positive fundamentals)
- No backtest validation (bullish side has 293 tests)
- Applied to ALL sectors (including incompatible ones)
- Missing strategy mapping and risk guidance
- No market regime context

**Why not lower?**
- Technical filter inversions are correct
- Bounce_bearish and death_cross_bearish are conceptually sound
- Good starting point that needs refinement

---

## Action Items

### Before Using Bearish Screeners in Production:

- [ ] **CRITICAL:** Fix universe_bearish fundamental filter conflict
- [ ] Remove bearish screeners from Utilities and Energy sectors
- [ ] Add sector-specific bearish strategy guidance to policy.json
- [ ] Validate with historical data (even basic tests)
- [ ] Add market regime checks (SPY/QQQ trend confirmation)
- [ ] Document bearish risk considerations in CLAUDE.md
- [ ] Test on paper trading before live capital

### Can Use Immediately (with caution):

- ✅ **bounce_bearish** (for bear flag setups in Technology/Discretionary)
- ✅ **death_cross_bearish** (for confirmed downtrends with confluence)

### Do NOT Use Yet:

- ❌ **universe_bearish** (fix fundamental filters first)
- ❌ **breakdown_bearish** (needs additional safety filters)

---

## Conclusion

The bearish screeners are **logically structured** but **conceptually incomplete**. They're like building a race car with the steering wheel installed backwards—technically functional but dangerous without fixes.

**Recommendation:** Fix universe_bearish immediately, add sector filtering logic, and validate with at least basic historical testing before relying on these for real trades.

The bullish screeners took 293 backtests to validate. The bearish screeners took 10 minutes of Python. **That's the quality difference.**
