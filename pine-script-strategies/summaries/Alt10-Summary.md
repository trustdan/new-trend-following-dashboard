# Alt10 - Profit Targets (3N/6N/9N)

## ⭐ #1 Overall Winner - 76% Success Rate

---

## 📊 Quick Stats

| Metric | Value |
|--------|-------|
| **Success Rate** | 76.19% (highest of all strategies) |
| **Options Suitability** | ⭐⭐⭐⭐⭐ Excellent |
| **Hold Time** | 3-10 weeks |
| **Best For** | Universal - works on stocks and ETFs |
| **Performance** | Healthcare +33.13%, Technology strong, SPY excellent |

---

## 🎯 What This Strategy Does

Alt10 uses **systematic profit targets** to lock in gains at three specific levels:

1. **3N profit target** - Take 1/3 position off the table
2. **6N profit target** - Take another 1/3 position off
3. **9N profit target** - Exit final 1/3 position
4. **Stop loss at 2N** - Trailing stop below highest close

**N = ATR (Average True Range)** - A volatility-based measurement

### Example Trade:
- Entry: Stock breaks out at $100, ATR = $2
- 3N target: Exit 1/3 at $106 (3 × $2 ATR)
- 6N target: Exit 1/3 at $112 (6 × $2 ATR)
- 9N target: Exit final 1/3 at $118 (9 × $2 ATR)
- Stop loss: Trail at 2N ($4) below highest close

---

## ✅ Best Use Cases

### Excellent For:
- **Healthcare stocks/ETFs** (UNH, XLV) - +33.13% in backtests
- **SPY/QQQ** - Broad market exposure
- **Consumer Cyclical** (XLY) - +21.39%
- **Any trending sector** - Universal applicability

### Performance by Security:
- UNH (Healthcare): +33.13%
- XLV (Healthcare ETF): Excellent
- SPY: Strong performance
- QQQ: Strong performance
- XLY (Consumer): +21.39%

---

## 📈 Options Trading Compatibility

### ⭐⭐⭐⭐⭐ EXCELLENT for Options

**Why It's Perfect for Options:**
- 3-10 week hold time matches 45-75 DTE options perfectly
- Systematic profit-taking aligns with theta decay management
- Scale out at each target level (close 1/3 of contracts at each stop)

### Options Strategy Examples:

**Bull Call Spread:**
1. Enter spread on Donchian breakout
2. Close 1/3 at 3N (30% profit)
3. Close 1/3 at 6N (60% profit)
4. Close final 1/3 at 9N (90% profit)

**Long Call:**
1. Buy calls with 60-75 DTE
2. Scale out: 33% at 3N, 33% at 6N, 34% at 9N
3. Trail stop at 2N below peak

---

## 🎓 How to Use in TradingView

1. **Apply the strategy** to a stock above SMA200
2. **Wait for Donchian breakout** (20-period high)
3. **Enter position** on breakout confirmation
4. **Set profit alerts:**
   - Alert at entry price + 3N
   - Alert at entry price + 6N
   - Alert at entry price + 9N
5. **Trail stop loss** at 2N below highest close
6. **Scale out** 1/3 position at each profit target

---

## 🟢 Sector Compatibility

| Sector | Rating | Notes |
|--------|--------|-------|
| **Healthcare** | 🟢 Excellent | +33.13% - Best sector performance |
| **Technology** | 🟢 Excellent | Strong performance on tech stocks |
| **Consumer Cyclical** | 🟢 Excellent | +21.39% on XLY |
| **Industrials** | 🟢 Good | Solid performance |
| **Financials** | 🟢 Good | Works well |
| **Energy** | 🟡 Marginal | Use with caution |
| **Utilities** | 🔴 AVOID | -12.4% - Mean-reverting sector |

---

## 💡 Pro Tips

### Risk Management:
- **Never risk more than 2% per trade** - Position size based on 2N stop
- **Scale out systematically** - Discipline prevents emotional decisions
- **Trail stop rigorously** - Protects profits in reversals

### Options-Specific Tips:
- Match expiration to 3-10 week hold: **45-75 DTE sweet spot**
- Close earliest expiring leg first at 3N
- Keep longest expiring leg for 9N target
- Don't hold past 7 DTE (theta decay acceleration)

### Common Mistakes to Avoid:
- ❌ Holding through all three targets hoping for 9N (market may reverse)
- ❌ Moving profit targets (discipline is critical)
- ❌ Ignoring the 2N stop loss (protect capital first)
- ❌ Using on Utilities sector (0% success rate)

---

## 📊 Backtest Performance Summary

Based on 293 validated backtests across 21 securities:

- **Overall Success Rate:** 76.19% (highest of all strategies)
- **Healthcare Performance:** +33.13% (UNH), Record on XLV
- **Consumer Cyclical:** +21.39% (XLY)
- **Technology:** Strong across MSFT, GOOGL
- **Broad Market:** Excellent on SPY, QQQ

**Risk-Adjusted Returns:** Best Sharpe ratio among all tested strategies

---

## 🎯 Why This Strategy Works

1. **Systematic profit-taking** prevents holding winners too long
2. **Volatility-based targets** adapt to market conditions
3. **Trailing stop** protects accumulated gains
4. **Universal applicability** - works across sectors and securities
5. **Proven track record** - 76% win rate in rigorous backtesting

---

## 📚 Related Strategies

- **Alt39 (Age-Based Targets)** - Alternative time-adaptive approach for SPY
- **Alt46 (Sector-Adaptive)** - Healthcare specialist with similar performance
- **Alt26 (Fractional Pyramid)** - Better for SPY, but Alt10 more universal

---

## 🚀 Quick Start Checklist

- [ ] Security is in trending sector (Healthcare, Technology, Consumer)
- [ ] Price above SMA200 (long-term uptrend confirmed)
- [ ] Donchian 20-period breakout confirmed
- [ ] Position sized for 2N stop loss (2% account risk)
- [ ] Options expiration matches 45-75 DTE
- [ ] Profit target alerts set at 3N, 6N, 9N
- [ ] Plan to scale out 1/3 at each target
- [ ] Trailing stop set at 2N below highest close

---

## 📖 Additional Resources

- **[Pine Script Code](../alt10.pine)** - TradingView implementation
- **[Complete Strategies Guide](../PINE-SCRIPTS-GUIDE.md)** - All strategies comparison
- **[Master Screener Guide](../../screeners/MASTER-SCREENER-GUIDE.md)** - Finding candidates
- **[Policy Configuration](../../data/policy.v1.json)** - Sector suitability settings

---

**Strategy File:** `alt10.pine`
**Category:** Profit Targets
**Complexity:** Intermediate
**Recommended For:** All traders (universal strategy)

---

*Based on 293 validated backtests | 99.74% data quality | Last updated: November 4, 2025*
