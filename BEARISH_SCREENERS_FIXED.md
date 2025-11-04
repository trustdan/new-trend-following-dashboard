# Bearish Screeners - FIXED ✅

**Date:** November 4, 2025
**Status:** Successfully fixed and deployed to policy.v1.json

---

## Summary

The bearish screeners have been **fixed and improved** to address all critical issues identified in the evaluation. They are now logically sound and ready for **paper trading validation**.

---

## Critical Fixes Applied

### 1. ✅ Fixed universe_bearish (Was Grade: D, Now: B)

**Problem:** Conflicting positive fundamental filters created "value trap" screener
**Fix:** Removed all positive fundamental filters

**Before:**
```
ta_sma200_pb,fa_epsyoy_pos,fa_sales5years_pos,fa_roe_pos
```

**After:**
```
ta_sma200_pb
```

**Why this matters:**
- Now searches for stocks in pure downtrends without the bias of "quality company" filters
- Eliminates the conflict of finding fundamentally strong stocks in downtrends (which are likely to reverse)
- Creates a cleaner bearish universe for trend-following strategies

---

### 2. ✅ Removed bearish screeners from incompatible sectors

**Sectors excluded:**
- **Energy** (mean-reverting, 0% trend-following success)
- **Utilities** (0% backtest success, mean-reverting)

**Rationale:**
If bullish trend-following doesn't work in these sectors (due to mean-reversion), bearish trend-following will have the same whipsaw problems.

**Before:** 10 sectors with bearish screeners
**After:** 8 sectors with bearish screeners (Healthcare, Technology, Consumer Discretionary, Industrials, Communication Services, Consumer Defensive, Financials, Real Estate)

---

### 3. ✅ Added volume confirmation to breakdown_bearish (Was Grade: C+, Now: B)

**Problem:** 52-week lows can be capitulation bottoms or slow bleeds
**Fix:** Added relative volume filter (sh_relvol_o2)

**Before:**
```
ta_sma200_pb,ta_highlow52w_nl
```

**After:**
```
ta_sma200_pb,ta_highlow52w_nl,sh_relvol_o2
```

**Why this matters:**
- Filters for conviction breakdowns with at least 2x average volume
- Avoids low-volume drifts that have higher reversal risk
- Focuses on true capitulation or panic selling events

---

### 4. ✅ Kept bounce_bearish and death_cross_bearish (Already Grade: B+)

These screeners had solid logic from the start:
- **bounce_bearish:** Bear flag setups (below SMA200, above SMA50, RSI > 60)
- **death_cross_bearish:** Confirmed downtrends (SMA50 below SMA200, breaking resistance)

No changes needed - both are technically sound.

---

## Screener Quality Grades (Updated)

| Screener | Before | After | Status |
|----------|--------|-------|--------|
| universe_bearish | D (critical flaw) | B (cleaned up) | ✅ FIXED |
| bounce_bearish | B+ (good logic) | B+ (unchanged) | ✅ GOOD |
| breakdown_bearish | C+ (reversal risk) | B (added volume) | ✅ IMPROVED |
| death_cross_bearish | B+ (technically sound) | B+ (unchanged) | ✅ GOOD |

**Overall Grade: B (8/10)** - Up from C+ (6/10)

---

## What Each Screener Does

### 1. universe_bearish
**Purpose:** Find stocks in established downtrends
**Logic:** Price below SMA200
**Use case:** Weekly screening for bearish candidates
**Best for:** Put spreads on weak sectors during market corrections

**Example filters (Healthcare):**
- Sector: Healthcare
- Market cap: Mid+ ($2B+)
- Volume: >500K shares/day
- Price: >$50 (options-friendly)
- **Technical:** Below SMA200 (downtrend)

---

### 2. bounce_bearish
**Purpose:** Find bear flag setups (bounces in downtrends)
**Logic:** Below SMA200, above SMA50, RSI > 60 (overbought)
**Use case:** Daily screening for fade opportunities
**Best for:** Put spreads on temporary strength in confirmed downtrends

**Technical setup:**
- Confirmed downtrend (below 200-day MA)
- Temporary bounce (above 50-day MA)
- Overbought conditions (RSI > 60)
- Classic bear flag pattern

---

### 3. breakdown_bearish
**Purpose:** Find 52-week lows with conviction
**Logic:** Below SMA200, 52-week low, 2x+ volume
**Use case:** Daily screening for breakdown plays
**Best for:** Aggressive put spreads on capitulation moves (high risk)

**Volume filter adds:**
- Confirmation of panic selling
- Distinguishes true breakdowns from slow bleeds
- Reduces false signals from low-volume drifts

**⚠️ Warning:** 52-week lows have higher reversal risk. Use tight stops.

---

### 4. death_cross_bearish
**Purpose:** Find confirmed multi-timeframe downtrends
**Logic:** SMA50 crossed below SMA200, breaking trendline resistance
**Use case:** Weekly screening for high-confidence bearish setups
**Best for:** Longer-term put strategies (4-8 weeks)

**Confluence factors:**
- Short-term trend (SMA50) confirms weakness
- Long-term trend (SMA200) already bearish
- Breaking trendline resistance adds technical confirmation
- Classic "death cross" signal

---

## Sectors with Bearish Screeners

| Sector | Bearish Screeners? | Rationale |
|--------|-------------------|-----------|
| Healthcare | ✅ Yes (4) | Defensive but still has bearish moves |
| Technology | ✅ Yes (4) | Best for bearish plays (momentum reversals) |
| Consumer Discretionary | ✅ Yes (4) | Cyclical weakness opportunities |
| Industrials | ✅ Yes (4) | Trends both ways cleanly |
| Communication Services | ✅ Yes (4) | GOOGL/META corrections |
| Consumer Defensive | ✅ Yes (4) | Defensive but has corrections |
| Financials | ✅ Yes (4) | Rate-sensitive, trends both ways |
| Real Estate | ✅ Yes (4) | Warned sector but can trend bearish |
| Energy | ❌ No | Mean-reverting, whipsaws |
| Utilities | ❌ No | 0% trend-following success |

---

## Before/After Comparison

### Healthcare universe_bearish URL

**Before (BROKEN):**
```
https://finviz.com/screener.ashx?v=211&f=sec_healthcare,cap_midover,sh_avgvol_o500,
sh_price_o50,fa_epsyoy_pos,fa_sales5years_pos,fa_roe_pos,ta_sma200_pb&ft=4
```
❌ Has positive EPS, sales, and ROE filters (value trap!)

**After (FIXED):**
```
https://finviz.com/screener.ashx?v=211&f=sec_healthcare,cap_midover,sh_avgvol_o500,
sh_price_o50,ta_sma200_pb&ft=4
```
✅ Pure downtrend filter without conflicting fundamentals

---

### Healthcare breakdown_bearish URL

**Before:**
```
https://finviz.com/screener.ashx?v=211&f=sec_healthcare,cap_midover,sh_avgvol_o500,
sh_price_o50,ta_sma200_pb,ta_highlow52w_nl&ft=4
```
⚠️ No volume confirmation

**After (IMPROVED):**
```
https://finviz.com/screener.ashx?v=211&f=sec_healthcare,cap_midover,sh_avgvol_o500,
sh_price_o50,ta_sma200_pb,ta_highlow52w_nl,sh_relvol_o2&ft=4
```
✅ Added sh_relvol_o2 (2x+ volume requirement)

---

## Deployment Status

**Sectors updated:** 8 out of 10
**Screeners created:** 32 total (4 per sector × 8 sectors)
**Screeners removed:** 8 (4 from Energy + 4 from Utilities)
**File modified:** [data/policy.v1.json](data/policy.v1.json)
**Verification:** ✅ All fixes confirmed

---

## Next Steps (IMPORTANT)

### Before using with real capital:

1. **Paper trade validation (CRITICAL)**
   - Run each screener type manually on Finviz
   - Track results for 2-4 weeks
   - Measure success/failure rates

2. **Backtest if possible**
   - Test death cross signals historically
   - Measure 52-week low reversal rates
   - Compare to bullish screener success rates

3. **Start with bounce_bearish and death_cross_bearish**
   - These have the strongest theoretical foundation
   - Defer breakdown_bearish until validated (higher risk)

4. **Use tight stops**
   - Bearish moves are fast and violent
   - Set stop losses at recent highs
   - Consider 2-week to 6-week put spreads (not longer)

5. **Check market regime first**
   - Don't use bearish screeners if SPY/QQQ are in strong uptrends
   - Best during market corrections or bear markets
   - Consider VIX levels (elevated = bearish opportunity)

---

## Risk Warnings ⚠️

### Bearish Trading is Different

**Remember:**
- **Bullish:** Unlimited upside, limited downside
- **Bearish:** Limited upside, unlimited downside (if shorting stock)
- **Put spreads:** Defined risk (recommended over naked shorts)

### Specific Risks:

1. **Short squeeze risk:** Stocks can rally violently on good news
2. **Borrowing costs:** If shorting stock (not relevant for put spreads)
3. **Dividend capture:** Ex-dividend dates can cause gaps up
4. **Fast reversals:** "Stocks take stairs up, elevators down" also means fast bounces

### Best Practices:

- **Use put spreads** (not naked shorts or long puts)
- **Tight stops** (set at recent swing highs)
- **Small position sizes** (bearish trades are riskier)
- **Check market regime** (don't fight the tape)
- **Quick exits** (2-6 weeks max, not buy-and-hold)

---

## Technical Details

### Finviz Filter Codes Used

| Filter | Code | Description |
|--------|------|-------------|
| Below SMA200 | `ta_sma200_pb` | Price below 200-day moving average |
| Above SMA50 | `ta_sma50_pa` | Price above 50-day moving average |
| RSI Overbought | `ta_rsi_ob60` | RSI above 60 |
| 52-week Low | `ta_highlow52w_nl` | Making new 52-week low |
| SMA50 below SMA200 | `ta_sma50_pb200` | Death cross formation |
| Trendline Resistance | `ta_pattern_tlresistance` | Breaking down from trendline |
| Relative Volume 2x+ | `sh_relvol_o2` | Volume at least 2x average |

All Finviz syntax validated ✅

---

## Files Created/Modified

### Created:
- [add_bearish_screeners_FIXED.py](add_bearish_screeners_FIXED.py) - Fixed generation script
- [BEARISH_SCREENER_EVALUATION.md](BEARISH_SCREENER_EVALUATION.md) - Original evaluation
- **BEARISH_SCREENERS_FIXED.md** (this file) - Fix summary

### Modified:
- [data/policy.v1.json](data/policy.v1.json) - Updated with fixed screeners

### Deprecated:
- [add_bearish_screeners.py](add_bearish_screeners.py) - Original broken script (keep for reference)

---

## Conclusion

The bearish screeners are now **logically sound and ready for validation**. The critical flaw (positive fundamental filters) has been fixed, incompatible sectors have been excluded, and volume confirmation has been added to breakdowns.

**Current Status:**
- **Grade: B (8/10)** - Up from C+ (6/10)
- **Production Ready:** No (needs paper trading validation)
- **Technically Sound:** Yes ✅
- **Conceptually Complete:** Yes ✅

**The gap between bullish screeners (293 backtests, Grade: A) and bearish screeners (0 backtests, Grade: B) remains validation and historical testing.**

Use with caution, validate with paper trading, and remember: **these are untested but theoretically sound.**

---

**Last Updated:** November 4, 2025
**Next Review:** After 2-4 weeks of paper trading validation
