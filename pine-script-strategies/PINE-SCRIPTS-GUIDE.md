# PINE SCRIPT STRATEGIES GUIDE
## Complete Reference for TF-Engine 2.0 Strategies
**Based on 293 validated backtests across 21 securities**

---

## 🎯 QUICK STRATEGY SELECTOR

### Need a recommendation? Start here:

| Your Situation | Best Strategy | Why |
|---------------|---------------|-----|
| **Healthcare stocks/ETFs** | Alt10 or Alt46 | 33% returns in backtests |
| **Tech momentum stocks (QQQ, GOOGL)** | Alt22 or Alt26 | Parabolic SAR loves breakouts |
| **SPY/broad market ETFs** | Alt39 or Alt26 | Best SPY performance |
| **Low drawdown priority** | Alt47 or Alt26 | Ultra-low drawdowns |
| **Don't know what to choose** | Baseline (Turtle Core) | Reliable benchmark across all sectors |

---

## 📊 COMPLETE STRATEGY REFERENCE TABLE

| Strategy | Label | Best For | Options Suitability | Hold Time | Win Rate* |
|----------|-------|----------|---------------------|-----------|-----------|
| **Alt10** | Profit Targets (3N/6N/9N) | Universal - works everywhere | ⭐⭐⭐⭐⭐ Excellent | 3-10 weeks | 76% |
| **Alt26** | Fractional Pyramid | SPY, Tech momentum | ⭐⭐⭐⭐⭐ Excellent | 2-8 weeks | High |
| **Alt43** | Volatility-Adaptive Targets | Healthcare ETFs, stocks | ⭐⭐⭐⭐⭐ Excellent (stocks) | 3-12 weeks | High |
| **Alt46** | Sector-Adaptive Parameters | Healthcare specialist | ⭐⭐⭐⭐⭐ Excellent (healthcare) | 3-12 weeks | 92%** |
| **Alt22** | Parabolic SAR | QQQ, tech breakouts | ⭐⭐⭐⭐ Good (tech only) | 2-6 weeks | Medium |
| **Alt47** | Momentum-Scaled Sizing | Individual stocks only | ⭐⭐⭐⭐⭐ Excellent (stocks) | 4-10 weeks | High |
| **Alt39** | Age-Based Targets | SPY, QQQ, healthcare | ⭐⭐⭐⭐⭐ Excellent | 3-12 weeks | High |
| **Alt45** | Dual-Momentum Confirmation | Stocks, QQQ (not SPY) | ⭐⭐⭐⭐ Good | 3-12 weeks | High |
| **Alt28** | ADX Filter | Trending healthcare/tech | ⭐⭐⭐ Selective | Variable | Medium |
| **Alt9** | Time Exit (40 bars) | Growth stocks only | ⭐⭐ Selective | 8 weeks | Low |
| **Baseline** | Turtle Core v2.2 | Benchmark - all sectors | ⭐⭐⭐ Good | 4-12 weeks | Medium |

*Win rate from 293 backtest dataset
**Healthcare sector success rate

---

## 🏆 THE TOP 5 STRATEGIES (Ranked by Performance)

### 1. Alt10 - Profit Targets (3N/6N/9N)
**⭐ #1 Overall Winner - 76% Success Rate**

**What It Does:**
- Takes profits at 3N, 6N, and 9N (where N = ATR)
- Systematic scaling out preserves gains
- Stop loss at 2N trailing

**Best Use Cases:**
- Healthcare: +33.13% (UNH, XLV)
- SPY/QQQ broad market
- ANY sector with clean trends

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT
- Perfect for multi-leg spreads
- Scale out at each target (close 1/3 at each level)
- 3-10 week hold matches LEAPS perfectly

**Sector Performance:**
- 🟢 Healthcare: +33.13%
- 🟢 Technology: Strong
- 🟢 Consumer Cyclical: +21.39%
- 🔴 Utilities: -12.4% (AVOID)

**How to Use in TradingView:**
1. Apply to any trending stock above SMA200
2. Enter on Donchian breakout
3. Scale out 1/3 position at 3N, 6N, 9N profit targets
4. Trail stop at 2N below highest close

**Pro Tips:**
- Universal strategy - works on 90% of securities
- Best risk-adjusted returns across all tests
- If unsure, choose this one

---

### 2. Alt46 - Sector-Adaptive Parameters
**⭐ Healthcare Specialist - 92% Sector Success Rate**

**What It Does:**
- Adjusts entry/exit thresholds based on sector volatility
- Healthcare uses tighter parameters
- Tech uses looser parameters for momentum

**Best Use Cases:**
- Healthcare stocks/ETFs (XLV, UNH): +32.16%
- Tech mega-caps (MSFT, GOOGL)
- When you know your sector

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT (Healthcare)
- Sector-optimized hold times
- Healthcare = 4-10 weeks (perfect for options)
- Tech = 2-8 weeks (faster momentum)

**Sector Performance:**
- 🟢 Healthcare: +32.16% (near-record XLV)
- 🟢 Technology: Strong
- 🟡 Other sectors: Marginal
- 🔴 Utilities: -6.2% (AVOID)

**How to Use in TradingView:**
1. Select security from preferred sector
2. Strategy auto-adjusts parameters
3. Healthcare = conservative, Tech = aggressive
4. Follow signals exactly as given

**Pro Tips:**
- Use for healthcare trades FIRST
- Backtested specifically for sector optimization
- Near-record results on healthcare ETFs

---

### 3. Alt26 - Fractional Pyramid (100%→75%→50%→25%)
**⭐ Best SPY Result - Ultra-Low Drawdowns**

**What It Does:**
- Adds to winners in fractional increments
- Start 100%, add 75%, add 50%, add 25%
- Builds bigger positions in confirmed trends

**Best Use Cases:**
- SPY (best result in entire backtest suite)
- QQQ, XLV
- Tech momentum stocks (MSFT, WMT)

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT
- Perfect for options pyramiding
- Add contracts as profit grows
- 2-8 week hold = ideal for 30-60 DTE options

**Sector Performance:**
- 🟢 Technology: BEST
- 🟢 Healthcare: Strong
- 🟢 SPY: Best SPY result
- 🟡 Defensive sectors: Okay

**How to Use in TradingView:**
1. Enter first position on Donchian breakout
2. Add 75% more on next Donchian breakout
3. Add 50% more on subsequent breakout
4. Add final 25% if trend continues
5. Exit all on stop loss

**Pro Tips:**
- Requires strong trending moves
- Low drawdowns make it psychologically easier
- Best for building large positions

---

### 4. Alt43 - Volatility-Adaptive Targets
**⭐ Record Healthcare ETF Performance**

**What It Does:**
- Adjusts profit targets based on ATR volatility
- Higher volatility = wider targets
- Lower volatility = tighter targets

**Best Use Cases:**
- Healthcare ETFs (XLV) - RECORD performance
- Individual stocks (UNH, MSFT, CAT, WMT)
- NOT for SPY (underperforms vs Alt39)

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT (Stocks), ⭐⭐⭐ SELECTIVE (ETFs)
- Great for stocks with variable volatility
- Healthcare ETFs shine here
- 3-12 week holds work for LEAPS

**Sector Performance:**
- 🟢 Healthcare ETFs: RECORD (XLV best)
- 🟢 Individual stocks: Excellent
- 🔴 SPY: Underperforms (use Alt39 instead)

**How to Use in TradingView:**
1. Apply to healthcare ETFs or stocks
2. Strategy calculates volatility-adjusted targets
3. Wider targets in high volatility
4. Tighter targets in low volatility
5. Exit on target hit or stop loss

**Pro Tips:**
- Healthcare ETF specialist
- Avoid on SPY (Alt39 is better there)
- Best for volatile individual names

---

### 5. Alt39 - Age-Based Targets
**⭐ Time-Adaptive Exits - Great on Slow Grinders**

**What It Does:**
- Profit targets widen as trade gets older
- Young trades = tight targets (take profits fast)
- Old trades = loose targets (let winners run)

**Best Use Cases:**
- SPY (excellent - better than Alt43)
- QQQ
- Healthcare (XLV, UNH)
- Slow-grinding stocks

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT
- Perfect for broad market ETFs
- Time-adaptive matches options decay
- 3-12 week holds = 30-90 DTE sweet spot

**Sector Performance:**
- 🟢 SPY: Excellent (better than Alt43)
- 🟢 QQQ: Excellent
- 🟢 Healthcare: Great
- 🟢 Slow grinders: Perfect

**How to Use in TradingView:**
1. Apply to SPY, QQQ, or healthcare
2. Enter on Donchian breakout
3. Young trades: Take profits quickly
4. Old trades: Let winners run longer
5. Exit on age-based target or stop

**Pro Tips:**
- Best SPY strategy for volatility-adaptive exits
- Beats Alt43 on broad market ETFs
- Time-adaptive logic matches options decay

---

## 📈 SPECIALIZED STRATEGIES (Use in Specific Scenarios)

### Alt22 - Parabolic SAR
**Best for:** Tech momentum breakouts (QQQ, GOOGL, AMZN)

**What It Does:**
- Uses Parabolic SAR indicator for entries/exits
- Fast-moving, high-churn strategy
- Catches explosive momentum moves

**Options Compatibility:** ⭐⭐⭐⭐ GOOD (Tech only)
- Excellent for QQQ and tech momentum
- AVOID defensive sectors (healthcare, utilities)
- 2-6 week holds = 30-45 DTE options

**When to Use:**
- Tech stocks breaking out to new highs
- QQQ in strong uptrend
- High relative volume momentum

**When to AVOID:**
- Healthcare (defensive sector)
- Utilities (mean-reverting)
- Low volatility environments

**Pro Tip:** High churn = more trades but more commissions. Works best in momentum markets.

---

### Alt47 - Momentum-Scaled Sizing
**Best for:** Individual stocks with ultra-low drawdown priority

**What It Does:**
- Scales position size based on momentum strength
- Strong momentum = bigger size
- Weak momentum = smaller size

**Options Compatibility:** ⭐⭐⭐⭐⭐ EXCELLENT (Stocks only)
- Ultra-low drawdowns (-3% to -5%)
- Perfect for risk-averse traders
- AVOID ALL ETFs (catastrophic on ETFs)

**When to Use:**
- Individual stocks (MSFT, CAT, PLD, WMT)
- Low drawdown requirement
- Steady compounding priority

**When to AVOID:**
- SPY, QQQ, or ANY ETF
- High volatility environments

**Pro Tip:** Best drawdown profile of all strategies on stocks. Never use on ETFs.

---

### Alt45 - Dual-Momentum Confirmation
**Best for:** Stocks, QQQ (not SPY)

**What It Does:**
- Requires BOTH price momentum AND relative strength
- Filters out weak breakouts
- 2nd-best consistency across backtests

**Options Compatibility:** ⭐⭐⭐⭐ GOOD (Stocks), ⭐⭐ WEAK (SPY)
- Excellent on individual stocks
- Great on QQQ
- Weak on SPY (use Alt39 instead)

**When to Use:**
- Quality over quantity (fewer but better trades)
- Stocks with strong relative strength
- QQQ, XLV, CAT, GOOGL, WMT

**When to AVOID:**
- SPY (use Alt39 or Alt10 instead)
- Low volatility markets

**Pro Tip:** Best for avoiding false breakouts. Fewer trades but higher quality.

---

### Alt28 - ADX Filter
**Best for:** Diagnostic tool for trending markets

**What It Does:**
- Only enters when ADX > threshold (strong trend)
- Filters out choppy, range-bound markets
- Diagnostic tool more than standalone strategy

**Options Compatibility:** ⭐⭐⭐ SELECTIVE
- Good for confirming trend strength
- Healthcare trending names (UNH, CAT, MSFT)
- Can over-filter choppy sectors

**When to Use:**
- Confirming trend strength before entry
- Healthcare/tech trending stocks
- As filter COMBINED with other strategies

**When to AVOID:**
- Choppy sectors (Energy, Real Estate)
- As standalone strategy

**Pro Tip:** Best used as diagnostic tool. Over-filters in sideways markets.

---

### Alt9 - Time Exit (40 bars)
**Best for:** Growth stocks (not a primary strategy)

**What It Does:**
- Forces exit after 40 bars (8 weeks)
- Prevents holding losers too long
- Fixed time-based risk management

**Options Compatibility:** ⭐⭐ SELECTIVE (Stocks only)
- Okay on growth stocks (AMZN, WMT, UNH, CAT)
- CATASTROPHIC on SPY (-20%+)
- AVOID utilities (0% success rate)

**When to Use:**
- Growth stocks with momentum
- When you want forced exit discipline
- Non-ETF securities

**When to AVOID:**
- SPY (catastrophic results)
- Utilities (0% success rate)
- Defensive sectors
- As primary strategy

**Pro Tip:** Not recommended as primary strategy. Use Alt10, Alt39, or Baseline instead.

---

## 🛡️ BASELINE STRATEGY

### Turtle Core v2.2 (Baseline)
**The Reliable Benchmark**

**What It Does:**
- Classic Turtle Trading rules
- Donchian breakout entries
- Regular exits with pyramiding
- Time-tested trend-following

**Options Compatibility:** ⭐⭐⭐ GOOD
- Reliable across most sectors
- 4-12 week holds = 45-90 DTE options
- Solid benchmark performance

**Best Use Cases:**
- MSFT, AMZN, WMT, CAT, XLV
- When you want proven, conservative approach
- Benchmark for comparing other strategies

**When to Use:**
- Don't want to over-think strategy selection
- Want reliable, proven approach
- Benchmarking other strategies

**Pro Tip:** Use this if you're new to trend-following. It's the foundation all other strategies build on.

---

## 🚫 STRATEGIES TO AVOID

### Alt20 - Asymmetric Long/Short
**❌ DO NOT USE - Catastrophic Performance**

- Disabled in policy.json
- Catastrophic losses across all assets
- Long/short asymmetry doesn't work in trending markets
- Use long-only strategies instead

---

## 🎓 HOW TO USE THIS GUIDE

### Step 1: Identify Your Sector
- **Healthcare?** → Alt10, Alt46, Alt43, Alt39, Alt28
- **Technology?** → Alt26, Alt22, Alt47, Alt10
- **SPY/Broad Market?** → Alt39, Alt26, Alt10
- **Don't know?** → Baseline (Turtle Core)

### Step 2: Match Your Goals
- **Best overall performance?** → Alt10
- **Lowest drawdowns?** → Alt47 (stocks), Alt26 (ETFs)
- **Options-friendly holds?** → Alt10, Alt26, Alt39, Alt43
- **Sector specialist?** → Alt46 (healthcare)

### Step 3: Check Sector Compatibility
- Look at policy.json `strategy_suitability` ratings
- 🟢 Green = Excellent/Good (no warning)
- 🟡 Yellow = Marginal (acknowledgement required)
- 🔴 Red = Incompatible (strong warning)

### Step 4: Apply in TradingView
1. Copy Pine Script code from respective .pine file
2. Open TradingView Pine Editor
3. Paste code and click "Add to Chart"
4. Follow strategy signals for entries/exits

---

## 📚 SECTOR COMPATIBILITY MATRIX

| Strategy | Healthcare | Technology | Consumer | Industrials | Energy | Utilities |
|----------|-----------|------------|----------|-------------|--------|-----------|
| Alt10 | 🟢 Best | 🟢 Great | 🟢 Great | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt46 | 🟢 BEST | 🟢 Good | 🟡 Okay | 🟡 Okay | 🟡 Marginal | 🔴 AVOID |
| Alt43 | 🟢 Record | 🟢 Good | 🟢 Good | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt39 | 🟢 Great | 🟢 Great | 🟢 Good | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt26 | 🟡 Marginal | 🟢 BEST | 🟢 Great | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt22 | 🔴 Weak | 🟢 BEST | 🟡 Okay | 🟡 Okay | 🟡 Marginal | 🔴 AVOID |
| Alt47 | 🟡 Marginal | 🟢 Great | 🟢 Good | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt45 | 🟡 Marginal | 🟢 Good | 🟢 Good | 🟢 Good | 🟡 Marginal | 🔴 AVOID |
| Alt28 | 🟢 Good | 🟢 Good | 🟡 Okay | 🟡 Okay | 🔴 Over-filters | 🔴 AVOID |
| Alt9 | 🔴 Weak | 🟡 Okay | 🟡 Okay | 🟡 Okay | 🔴 Weak | 🔴 AVOID |
| Baseline | 🟢 Good | 🟢 Good | 🟢 Good | 🟢 Good | 🟡 Marginal | 🔴 AVOID |

**Key:**
- 🟢 Green = Excellent/Good performance (recommended)
- 🟡 Yellow = Marginal performance (use with caution)
- 🔴 Red = Poor/Incompatible (avoid)

---

## 💡 PRACTICAL TIPS

### For Options Traders
1. **Best strategies for options:** Alt10, Alt26, Alt39, Alt43, Alt46
2. **Match hold time to expiration:**
   - 2-6 week strategies → 30-45 DTE
   - 3-10 week strategies → 45-75 DTE
   - 4-12 week strategies → 60-90 DTE
3. **Scale out at profit targets** (especially Alt10)
4. **Avoid short-hold strategies** (Alt22) for far-dated options

### For Low Drawdown Priority
1. **Best:** Alt47 (stocks only) - Ultra-low drawdowns
2. **Second:** Alt26 - Low drawdowns on ETFs
3. **Third:** Alt45 - Consistent, fewer false breakouts
4. **Avoid:** Alt22 (high churn), Alt9 (time-based risk)

### For Sector Specialists
1. **Healthcare:** Alt46 (92% sector success), Alt10 (+33.13%)
2. **Technology:** Alt22 (QQQ specialist), Alt26 (best SPY)
3. **SPY/Broad Market:** Alt39 (beats Alt43), Alt26, Alt10
4. **Individual Stocks:** Alt47 (ultra-low DD), Alt10, Alt39

### For Beginners
1. **Start with:** Baseline (Turtle Core v2.2)
2. **Graduate to:** Alt10 (universal winner)
3. **Specialize with:** Alt46 (healthcare), Alt26 (tech)
4. **Avoid:** Alt9, Alt20, Alt22 (until experienced)

---

## 🔬 RESEARCH FOUNDATION

This guide is based on:
- **293 validated backtests** across 21 securities
- **14 trend-following strategies** tested
- **21 securities:** Individual stocks + ETFs (SPY, QQQ, XLV, XLK, etc.)
- **99.74% data quality** (only 2 corrupted tests)

### Key Research Findings:
1. **Profit targets work:** Alt10 achieved 76.19% success rate
2. **Pyramiding works:** Alt26 had best SPY result
3. **Sector matters:** Healthcare 92% success, Utilities 0% success
4. **Time-based exits fail:** Alt9 catastrophic on SPY (-20%+)
5. **Long/short asymmetry fails:** Alt20 disabled (catastrophic)

### Statistical Validation:
- Python logistic regression confirms:
  - Profit targets: 4.47× odds of success
  - Pyramiding: 3.96× odds of success
  - Parabolic SAR: 0.156× odds (hurts on wrong assets)

---

## 📖 RELATED DOCUMENTATION

- **[DISCOVERIES_AND_LEARNINGS.md](../backtesting-lessons/DISCOVERIES_AND_LEARNINGS.md)** - Complete 293-backtest analysis
- **[MASTER-SCREENER-GUIDE.md](../screeners/MASTER-SCREENER-GUIDE.md)** - How to find trade candidates
- **[policy.v1.json](../data/policy.v1.json)** - Current strategy configurations
- **[architects-intent.md](../architects-intent.md)** - Original design intent

---

## 🎯 FINAL RECOMMENDATIONS

### If you only remember 3 things:

1. **Alt10 is the universal winner** - 76% success rate, works everywhere except Utilities
2. **Sector matters more than strategy** - Healthcare 92% success, Utilities 0% success
3. **Match hold time to options expiration** - 3-10 week strategies = 45-75 DTE sweet spot

### Default Strategy Hierarchy:
1. **First choice:** Alt10 (universal)
2. **Healthcare specialist:** Alt46
3. **Tech/momentum:** Alt22 or Alt26
4. **SPY/broad market:** Alt39 or Alt26
5. **Don't know:** Baseline (Turtle Core)

---

**Last Updated:** November 4, 2025
**Policy Version:** 1.0.0
**Based on:** 293 validated backtests across 21 securities
**Data Quality:** 99.74%
