# Alt43 - Volatility-Adaptive Targets

## ⭐ Record Healthcare ETF Performance

---

## 📊 Quick Stats

| Metric | Value |
|--------|-------|
| **XLV Performance** | Record (best healthcare ETF result) |
| **Options Suitability** | ⭐⭐⭐⭐⭐ Excellent (Stocks), ⭐⭐⭐ Selective (ETFs) |
| **Hold Time** | 3-12 weeks |
| **Best For** | Healthcare ETFs (XLV), individual stocks, variable volatility |
| **Not Recommended** | SPY (use Alt39 instead) |

---

## 🎯 What This Strategy Does

Alt43 uses **volatility-adaptive profit targets** based on ATR:

- **High volatility period** → Wider targets (let moves breathe)
- **Low volatility period** → Tighter targets (take profits quickly)
- **Dynamic adjustment** → Targets update with market conditions

### The Logic: Adapt to Market Volatility

**High ATR (volatile market):**
- Profit targets: 4N, 8N, 12N (wide)
- Rationale: Volatile markets need room to move

**Low ATR (quiet market):**
- Profit targets: 2N, 4N, 6N (tight)
- Rationale: Quiet markets = take profits before reversal

### Real-Time Adaptation:
Unlike fixed targets (Alt10), Alt43 recalculates targets daily based on current ATR.

---

## ✅ Best Use Cases

### Excellent For:
- **Healthcare ETFs** (XLV) - RECORD performance
- **Individual stocks** (UNH, MSFT, CAT, WMT) - Great
- **Variable volatility securities** - Adapts naturally
- **NOT for SPY** - Alt39 beats it on broad market

### Performance by Security:
- XLV (Healthcare ETF): RECORD (best healthcare ETF result)
- UNH (Healthcare stock): Excellent
- MSFT (Tech stock): Excellent
- CAT (Industrial): Excellent
- WMT (Consumer): Excellent
- SPY: Underperforms (use Alt39 instead)

---

## 📈 Options Trading Compatibility

### ⭐⭐⭐⭐⭐ EXCELLENT for Stock Options

**Why It's Perfect for Stock Options:**
- 3-12 week hold = 45-90 DTE sweet spot
- Volatility adaptation matches options IV
- High volatility = wider targets = more time for options
- Low volatility = tighter targets = exit before IV crush

### ⭐⭐⭐ SELECTIVE for ETF Options

**ETF Considerations:**
- Healthcare ETFs (XLV) = EXCELLENT (record performance)
- SPY = AVOID (use Alt39 instead)
- QQQ = OKAY (Alt22 or Alt26 better)

### Options Strategy Examples:

**XLV Bull Call Spread (Healthcare ETF):**
1. Enter on healthcare breakout (60-75 DTE)
2. Strategy adapts targets to XLV volatility
3. High IV: Hold for wider targets
4. Low IV: Exit quickly before crush
5. Record performance on XLV backtests

**Individual Stock Long Call (UNH, MSFT):**
1. Buy calls with 60-90 DTE
2. Volatile period: Patient hold (wide targets)
3. Quiet period: Quick exit (tight targets)
4. IV expansion = hold, IV contraction = exit

---

## 🎓 How to Use in TradingView

1. **Apply to healthcare ETF (XLV) or individual stock**
2. **Enter on Donchian breakout**
3. **Strategy calculates current ATR volatility:**
   - High ATR → Sets wide targets (4N/8N/12N)
   - Low ATR → Sets tight targets (2N/4N/6N)
4. **Targets update daily** based on rolling ATR
5. **Exit on target hit OR trail stop (2N)**

### Critical Understanding:
- Targets are NOT fixed (unlike Alt10)
- Targets adapt to current market volatility
- High volatility requires patience
- Low volatility requires speed

---

## 🟢 Sector Compatibility

| Sector | Rating | Notes |
|--------|--------|-------|
| **Healthcare (XLV)** | 🟢 EXCELLENT | RECORD performance on healthcare ETF |
| **Healthcare Stocks** | 🟢 Excellent | Great on UNH, JNJ, etc. |
| **Technology Stocks** | 🟢 Excellent | MSFT, GOOGL strong |
| **Industrials** | 🟢 Excellent | CAT performs well |
| **Consumer** | 🟢 Excellent | WMT strong |
| **SPY** | 🔴 Avoid | Alt39 better for broad market |
| **Utilities** | 🔴 AVOID | Mean-reverting sector |

---

## 💡 Pro Tips

### Volatility-Adaptive Philosophy:
- **High volatility = patience** → Don't exit too early
- **Low volatility = speed** → Take profits before reversal
- **Trust the adaptation** → Tested on 293 backtests
- **Healthcare ETFs = sweet spot** → Record XLV performance

### Risk Management:
- Position size for 2N stop loss
- High volatility = looser stop (more room)
- Low volatility = tighter stop (less room)
- Trail stop adapts with ATR changes

### Options-Specific Tips:
- **High IV:** Buy options (volatility-adaptive targets give time)
- **Low IV:** Sell spreads (tight targets = quick profits)
- **XLV specialty:** 60-75 DTE for healthcare ETF
- **Individual stocks:** 60-90 DTE for volatility swings

### Common Mistakes to Avoid:
- ❌ Using on SPY (Alt39 is better)
- ❌ Expecting fixed targets (they adapt daily)
- ❌ Exiting too early in high volatility
- ❌ Holding too long in low volatility

---

## 📊 Backtest Performance Summary

Based on 293 validated backtests across 21 securities:

- **XLV Performance:** RECORD (best healthcare ETF result)
- **Individual Stocks:** Excellent across UNH, MSFT, CAT, WMT
- **SPY Performance:** Underperforms (use Alt39 instead)
- **Volatility Adaptation:** Proven effective in variable conditions

**Key Insight:** Healthcare ETF specialist, but beats Alt39 on stocks (not SPY)

---

## 🎯 Why This Strategy Works

1. **Volatility adaptation** matches market conditions
2. **Record XLV performance** - Healthcare ETF specialist
3. **Variable targets** prevent premature exits (high vol) or late exits (low vol)
4. **Stock-friendly** - Great on individual names
5. **Options IV alignment** - Targets match implied volatility

---

## 📚 Related Strategies

- **Alt39 (Age-Based Targets)** - Better for SPY (time-based vs volatility-based)
- **Alt10 (Profit Targets)** - Fixed targets (3N/6N/9N) vs adaptive
- **Alt46 (Sector-Adaptive)** - Healthcare specialist (comparable XLV performance)
- **Baseline (Turtle Core)** - Benchmark comparison

---

## 🔬 Volatility-Adaptive Logic

### Target Calculation:

**Current ATR Measurement:**
```
ATR(14) = Average True Range over 14 periods
Volatility level = Current ATR vs 50-period average ATR
```

**High Volatility (ATR > avg):**
```
Profit targets: 4N, 8N, 12N (wide)
Stop loss: 2.5N (looser)
Hold time: 6-12 weeks (patient)
```

**Low Volatility (ATR < avg):**
```
Profit targets: 2N, 4N, 6N (tight)
Stop loss: 1.5N (tighter)
Hold time: 3-6 weeks (faster)
```

**Medium Volatility (ATR = avg):**
```
Profit targets: 3N, 6N, 9N (standard)
Stop loss: 2N (normal)
Hold time: 4-8 weeks (typical)
```

---

## 🚀 Quick Start Checklist

- [ ] Security is healthcare ETF (XLV) or individual stock
- [ ] NOT using on SPY (use Alt39 for SPY)
- [ ] Price above SMA200 (uptrend confirmed)
- [ ] Donchian breakout confirmed
- [ ] Understand targets adapt daily with volatility
- [ ] Position sized for volatility-adjusted stop
- [ ] Options expiration matches 45-90 DTE
- [ ] Ready for patient holding in high volatility
- [ ] Ready for quick exits in low volatility

---

## 📖 Additional Resources

- **[Pine Script Code](../seykota_alt43_volatility_adaptive_targets.pine)** - TradingView implementation
- **[Complete Strategies Guide](../PINE-SCRIPTS-GUIDE.md)** - All strategies comparison
- **[Healthcare Screeners](../../screeners/MASTER-SCREENER-GUIDE.md)** - Finding XLV candidates
- **[Policy Configuration](../../data/policy.v1.json)** - Strategy settings

---

**Strategy File:** `seykota_alt43_volatility_adaptive_targets.pine`
**Category:** Volatility-Adaptive
**Complexity:** Advanced
**Recommended For:** Healthcare ETF traders, individual stock traders

---

*Based on 293 validated backtests | Record XLV performance | Last updated: November 4, 2025*
