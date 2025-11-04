# Screener and Notification Updates

**Date:** November 4, 2025
**Status:** Complete

## Summary

Added bearish screeners for all sectors and implemented Windows toast notifications when the cooldown timer completes.

---

## 1. Bearish Screeners Added

### What Changed
Doubled the available screeners by adding bearish versions alongside the existing bullish ones.

### New Screeners Per Sector

**Bearish Screeners:**
1. **Universe (Bearish)** - Stocks in long-term downtrends (below SMA200)
2. **Bounce (Bearish)** - Overbought stocks in downtrends (RSI > 60, price below SMA200)
3. **Breakdown (Bearish)** - New 52-week lows with momentum
4. **Death Cross (Bearish)** - SMA50 crossing below SMA200

### Finviz Filter Changes

| Bullish Filter | Bearish Equivalent |
|----------------|-------------------|
| `ta_sma200_pa` (price above SMA200) | `ta_sma200_pb` (price below SMA200) |
| `ta_sma50_pb` (price below SMA50) | `ta_sma50_pa` (price above SMA50) |
| `ta_rsi_os40` (RSI oversold < 40) | `ta_rsi_ob60` (RSI overbought > 60) |
| `ta_highlow52w_nh` (new 52-week high) | `ta_highlow52w_nl` (new 52-week low) |
| `ta_sma50_pa200` (SMA50 above SMA200) | `ta_sma50_pb200` (SMA50 below SMA200) |
| `ta_pattern_tlsupport` (trendline support) | `ta_pattern_tlresistance` (trendline resistance) |

### Files Modified
- [data/policy.v1.json](data/policy.v1.json) - Added bearish screener URLs to all sectors
- [internal/ui/screens/screener_launch.go](internal/ui/screens/screener_launch.go) - Updated UI to display bullish/bearish sections

### UI Changes

**Before:**
- Single list of 4 screeners per sector

**After:**
- **Bullish Section** (🟢 blue-bordered cards):
  - Universe Screener (Bullish)
  - Pullback Screener (Bullish)
  - Breakout Screener (Bullish)
  - Golden Cross Screener (Bullish)

- **Bearish Section** (🔴 red-bordered cards):
  - Universe Screener (Bearish)
  - Bounce Screener (Bearish)
  - Breakdown Screener (Bearish)
  - Death Cross Screener (Bearish)

Each sector now has **8 screeners total** (4 bullish + 4 bearish).

---

## 2. Neutral Screeners - Decision Made

### Analysis
Neutral screeners (consolidation, range-bound) were considered but **intentionally excluded** for the following reasons:

1. **Trend-following incompatibility**: The app's strategies are designed for directional momentum, not range-bound markets
2. **Anti-signal bias**: Neutral conditions are when trend-following strategies perform worst
3. **User confusion**: Including neutral screeners would encourage trading in unfavorable conditions
4. **Finviz limitations**: Hard to filter for "neutral" patterns consistently

### Recommendation
Focus remains on **strong bullish or bearish directional setups only**. This aligns with the research showing trend-following strategies fail in choppy/mean-reverting environments (e.g., Utilities 0% success rate).

---

## 3. Windows Toast Notification

### What Changed
When the cooldown timer completes on the Anti-Impulsivity Checklist screen, the app now sends a **Windows toast notification** in the bottom-right corner of the screen.

### Notification Details
- **Title:** "Trade Cooldown Complete!"
- **Content:** "You're cleared hot to proceed with [TICKER] trade"
- **Example:** "You're cleared hot to proceed with UNH trade"

### Implementation Details
- Uses Fyne's native `fyne.Notification` API
- Automatically includes the ticker symbol from `state.CurrentTrade.Ticker`
- Fires when cooldown timer reaches 0 seconds
- Appears as native Windows notification (bottom-right taskbar area)

### Files Modified
- [internal/ui/screens/checklist.go](internal/ui/screens/checklist.go) - Added notification callback to cooldown timer completion

### Code Added
```go
notification := &fyne.Notification{
    Title:   "Trade Cooldown Complete!",
    Content: fmt.Sprintf("You're cleared hot to proceed with %s trade", ticker),
}

app := fyne.CurrentApp()
if app != nil {
    app.SendNotification(notification)
}
```

---

## Testing Status

### Build Status
- ✅ Application compiles successfully with no errors
- ✅ All Go tests pass

### Manual Testing Required
1. **Screener UI**: Navigate to Screen 2 (Screener Launch) and verify:
   - Bullish section appears with 🟢 blue borders
   - Bearish section appears with 🔴 red borders
   - All 8 screeners display correctly per sector
   - Finviz URLs open correctly in browser

2. **Toast Notification**: Start a trade flow and verify:
   - Enter ticker symbol (e.g., "UNH")
   - Wait for cooldown timer to complete (5 minutes default)
   - Verify toast notification appears in bottom-right of Windows screen
   - Verify ticker symbol appears correctly in notification

---

## Sector Coverage

All **10 sectors** now have bearish screeners:
1. Healthcare ✅
2. Technology ✅
3. Consumer Discretionary ✅
4. Industrials ✅
5. Communication Services ✅
6. Consumer Defensive ✅
7. Financials ✅
8. Real Estate ✅
9. Energy ✅
10. Utilities ✅

---

## User Impact

### Benefits
1. **Expanded opportunities**: Users can now trade bearish/put strategies systematically
2. **Symmetrical workflow**: Same rigorous screening for shorts as longs
3. **Better awareness**: Toast notification ensures users don't miss cooldown completion
4. **Professional UX**: Native Windows notifications feel polished and integrated

### No Breaking Changes
- Existing bullish screeners unchanged
- Policy file maintains backward compatibility
- UI gracefully handles missing bearish URLs (won't crash if sector lacks them)

---

## Next Steps (Optional Future Enhancements)

1. **Notification settings**: Allow users to customize notification text or disable
2. **Sound alerts**: Add optional audio chime when cooldown completes
3. **Screener favorites**: Let users bookmark most-used screeners
4. **Mobile notifications**: If mobile companion app is built, sync notifications

---

**Status:** Ready for user testing and feedback
