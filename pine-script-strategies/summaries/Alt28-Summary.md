# Alt28 - ADX Filter

## ⭐ Diagnostic Tool for Trend Strength

---

## 📊 Quick Stats

**Primary Use:** Diagnostic tool / filter (not standalone strategy)

**Options Suitability:** ⭐⭐⭐ Selective

**Hold Time:** Variable (trend-dependent)

**Best For:** Confirming trend strength, filtering entries

**Weakness:** Over-filters in choppy markets

---

## 🎯 What This Strategy Does

Alt28 uses **ADX (Average Directional Index)** to filter entries:

- **ADX > 25** → Strong trend, enter trades
- **ADX < 25** → Weak trend, stay out
- **Diagnostic tool** → Confirms trend before committing capital
- **Filter function** → Combine with other strategies

### The ADX Indicator:
Measures trend strength on 0-100 scale:
- **0-20:** Weak/absent trend (range-bound)
- **20-25:** Emerging trend
- **25-50:** Strong trend (tradeable)
- **50-75:** Very strong trend
- **75-100:** Extremely strong trend (rare)

### Not a Standalone Strategy:
Best used WITH other strategies to confirm trend strength before entry.

---

## ✅ Best Use Cases

### Excellent For:
- **Confirming trend strength** - Is this stock really trending?
- **Healthcare trending names** (UNH, CAT, MSFT)
- **Filtering false breakouts** - ADX > 25 = real move
- **Diagnostic analysis** - Why did my trade fail? Check ADX

### Performance by Security:
- UNH (Healthcare): Good when trending
- CAT (Industrial): Strong on trends
- MSFT (Tech): Solid confirmation

---

## ❌ Weakness: Over-Filtering

| Market Type | ADX Behavior | Result |
|-------------|--------------|---------|
| **Strong Trend** | ADX > 30 | Great (enters winning trades) |
| **Choppy Market** | ADX 15-25 | Over-filters (misses recoveries) |
| **Range-Bound** | ADX < 20 | Correctly avoids (good) |
| **Energy Sector** | ADX 18-24 | Filters everything (too conservative) |

**Key Weakness:** Can miss early trend entries waiting for ADX confirmation

---

## 📈 Options Trading Compatibility

### ⭐⭐⭐ SELECTIVE for Options

**Why It's Selective:**
- Variable hold time (trend-dependent)
- Can filter out early entries (miss option value acceleration)
- Best for far-dated options (60-90 DTE)
- Useful as confirmation tool, not primary strategy

### Options Strategy Examples:

**Healthcare Bull Call Spread (UNH) with ADX Confirmation:**
1. Identify UNH breakout
2. Check ADX: Must be > 25 (strong trend confirmed)
3. Enter spread with 60-75 DTE
4. Hold while ADX stays elevated
5. Exit if ADX drops below 20 (trend weakening)

**Diagnostic Use (Post-Trade Analysis):**
1. Trade failed on MSFT
2. Check ADX at entry: Was only 18 (weak trend)
3. Lesson: Wait for ADX > 25 next time
4. Improves future trade selection

---

## 🎓 How to Use in TradingView

### As Standalone Strategy:
1. **Apply to trending stock** (healthcare, tech)
2. **Wait for ADX > 25** (strong trend confirmation)
3. **Enter on Donchian breakout** + ADX filter
4. **Exit when ADX < 20** (trend exhausted)

### As Filter/Diagnostic Tool (Recommended):
1. **Use Alt10 or Alt46 as primary strategy**
2. **Check ADX before entry:**
   - ADX > 30 = Strong trend (full position)
   - ADX 25-30 = Good trend (normal position)
   - ADX 20-25 = Weak trend (reduce size)
   - ADX < 20 = No trend (skip trade)
3. **Improves win rate** by filtering weak trends

---

## 🟢 Sector Compatibility

| Sector | Rating | Notes |
|--------|--------|-------|
| **Healthcare (trending)** | 🟢 Good | UNH when strong trend confirmed |
| **Technology (trending)** | 🟢 Good | MSFT, CAT when trending |
| **Industrials** | 🟡 Okay | Works on strong trends only |
| **Energy** | 🔴 Over-filters | Choppy sector = constant low ADX |
| **Utilities** | 🔴 AVOID | Mean-reverting (ADX always low) |
| **Choppy Markets** | 🔴 Over-filters | Misses recoveries |

---

## 💡 Pro Tips

### Best Use Cases:
- **Diagnostic tool** - Post-trade analysis (why did it fail?)
- **Filter for other strategies** - Combine with Alt10/Alt46
- **Trend confirmation** - Is this breakout real?
- **Position sizing aid** - Larger size on higher ADX

### Risk Management:
- Don't rely solely on ADX for entries
- Use as confirmation, not primary signal
- Strong ADX (>30) = increase position size
- Weak ADX (<25) = reduce or skip

### Options-Specific Tips:
- **60-90 DTE best** - Variable hold times
- **Use as filter** - Check ADX before buying options
- **ADX < 20** = Skip the trade (weak trend won't move options)
- **ADX > 30** = Full position (strong trend will move)

### Common Mistakes to Avoid:
- ❌ Using as standalone strategy (it's a filter)
- ❌ Waiting too long for ADX confirmation (miss entries)
- ❌ Using in choppy sectors (over-filters)
- ❌ Ignoring ADX < 20 warnings (weak trends fail)

---

## 📊 Backtest Performance Summary

Based on 293 validated backtests across 21 securities:

- **Trending Markets:** Good (confirms real trends)
- **Healthcare Stocks:** Solid on UNH, CAT, MSFT
- **Choppy Sectors:** Over-filters (Energy, Real Estate)
- **Primary Value:** Diagnostic tool, not primary strategy

**Key Insight:** Best as filter/confirmation, not standalone approach

---

## 🎯 Why This Strategy Works (As Filter)

1. **Trend confirmation** - Filters false breakouts
2. **Diagnostic value** - Explains trade failures
3. **Position sizing aid** - Scale with trend strength
4. **Improves win rate** - Only trade strong trends
5. **Risk reduction** - Avoids choppy markets

---

## 🎯 Why This Strategy FAILS (Standalone)

1. **Over-filtering** - Misses early trend entries
2. **Choppy markets** - Constant low ADX sidelines capital
3. **Energy sector** - Choppy = permanent filter
4. **Late entries** - Waiting for ADX > 25 = worse fills
5. **Better alternatives** - Alt10, Alt46 standalone

---

## 📚 Related Strategies

- **Alt10 (Profit Targets)** - Use ADX as filter for Alt10 entries
- **Alt46 (Sector-Adaptive)** - Combine with ADX for healthcare confirmation
- **Baseline (Turtle Core)** - Add ADX filter to Turtle entries
- **Alt43 (Volatility-Adaptive)** - Alternative adaptive approach

---

## 🔬 ADX Calculation & Interpretation

### ADX Formula:
```
+DI (Positive Directional Indicator) = Smoothed +DM / ATR
-DI (Negative Directional Indicator) = Smoothed -DM / ATR
DX (Directional Movement Index) = abs(+DI - -DI) / (+DI + -DI) × 100
ADX = Smoothed average of DX over 14 periods
```

### Interpretation:
```
ADX 0-20: No trend / Range-bound
  Action: Avoid trend-following trades

ADX 20-25: Weak trend emerging
  Action: Small positions or wait

ADX 25-50: Strong trend
  Action: Enter trades, normal position sizing

ADX 50-75: Very strong trend
  Action: Full positions, hold patiently

ADX 75-100: Extremely strong trend (rare)
  Action: Maximum positions, let it run
```

### As Filter Example (MSFT):

**Scenario 1: Strong ADX**
- MSFT breakout at $400
- ADX = 35 (strong trend)
- **Action:** Enter full position (trend confirmed)
- **Result:** Trend continues, profitable trade

**Scenario 2: Weak ADX**
- MSFT breakout at $400
- ADX = 18 (weak/absent trend)
- **Action:** Skip trade (no trend confirmation)
- **Result:** Price whipsaws, avoided losing trade

---

## 🚀 Quick Start Checklist

### As Standalone Strategy:
- [ ] Security is trending strongly (healthcare, tech)
- [ ] ADX > 25 (trend confirmed)
- [ ] Donchian breakout confirmed
- [ ] Position sized normally
- [ ] Exit when ADX < 20 (trend exhausted)

### As Filter/Diagnostic Tool (Recommended):
- [ ] Using Alt10 or Alt46 as primary strategy
- [ ] Check ADX before every entry
- [ ] ADX > 25 = proceed with entry
- [ ] ADX < 25 = reduce size or skip
- [ ] Post-trade: Review ADX to understand why trade succeeded/failed

---

## 📖 Additional Resources

- **[Pine Script Code](../alt28.pine)** - TradingView implementation
- **[Complete Strategies Guide](../PINE-SCRIPTS-GUIDE.md)** - All strategies comparison
- **[Healthcare Screeners](../../screeners/MASTER-SCREENER-GUIDE.md)** - Finding trending candidates
- **[Policy Configuration](../../data/policy.v1.json)** - Strategy settings

---

**Strategy File:** `alt28.pine`
**Category:** Trend Strength Filter / Diagnostic
**Complexity:** Intermediate
**Recommended For:** As filter for other strategies, diagnostic analysis

---

*Based on 293 validated backtests | Best as filter, not standalone | Last updated: November 4, 2025*
