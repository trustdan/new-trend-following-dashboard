# Phase 6 Troubleshooting Summary

## Overview
This document tracks issues encountered while implementing Phase 6 (warning-based permission system) for the TF-Engine application.

---

## Issue Timeline

### Issue #1: Missing Screener URLs for Utilities/Energy Sectors
**Date:** 2025-11-04
**Symptom:** Screen 2 showed "No screeners available for this sector" for Utilities and Energy
**Root Cause:** `screener_urls` field missing in policy.v1.json for these sectors
**Fix Applied:** Added complete Finviz screener configurations (Universe, Pullback, Breakout, Golden Cross) for both sectors
**Status:** ✅ RESOLVED

---

### Issue #2: Missing Screener URLs for 6 Additional Sectors
**Date:** 2025-11-04
**Symptom:** Screen 2 showed "No screeners available for this sector" for Consumer Discretionary, Industrials, Communication Services, Consumer Defensive, Financials, and Real Estate
**Root Cause:** Only 4 out of 10 sectors had `screener_urls` configured in policy.v1.json
**Fix Applied:** Added screener URLs for all 6 missing sectors:
- Consumer Discretionary (`sec_consumercyclical`)
- Industrials (`sec_industrialgoods`)
- Communication Services (`sec_communication`)
- Consumer Defensive (`sec_consumergoods`)
- Financials (`sec_financial`)
- Real Estate (`sec_realestate`)

**Status:** ✅ RESOLVED

---

### Issue #3: Color Emoji Indicators Not Rendering
**Date:** 2025-11-04
**Symptom:** Strategy dropdown and Screen 3 info banner showed empty boxes or missing characters instead of 🟢🟡🔴 emoji indicators
**Root Cause:** Fyne UI framework or Windows font rendering doesn't support color emoji rendering properly
**Fix Applied:**
- Replaced emoji indicators with text-based indicators: `[GREEN]`, `[YELLOW]`, `[RED]`
- Updated parsing logic in `ticker_entry.go` to handle text prefixes
- Updated all user-facing text to reference text indicators

**Code Changes:**
```go
// Before
return "🟢"  // Green emoji
return "🟡"  // Yellow emoji
return "🔴"  // Red emoji

// After
return "[GREEN]"
return "[YELLOW]"
return "[RED]"
```

**Status:** ✅ FIX APPLIED (needs verification)

---

### Issue #4: Visual Color Coding Not Prominent Enough
**Date:** 2025-11-04
**Symptom:** User couldn't see color-coded suitability ratings clearly in dropdown
**User Request:** Use colored vertical bars on strategy metadata cards instead of relying on emoji/text indicators
**Fix Applied:**
- Enhanced `displayStrategyMetadata()` function to use dynamic color-coded backgrounds
- Increased left border width from 4px to 8px for prominence
- Added color-coded background tints:
  - **Green strategies:** Light green background (#F0FFF0) + dark green border (#00B450)
  - **Yellow strategies:** Light yellow background (#FFFAE6) + amber border (#DCB400)
  - **Red strategies:** Light red background (#FFF0F0) + red border (#C80000)
- Added "Sector Suitability" label showing rating and color explicitly

**Visual Example:**
```
┌─────────────────────────────────────────┐
│ 📋 Profit Targets (3N/6N/9N)            │ ← Green bg
│ Sector Suitability: EXCELLENT [GREEN]   │
│ Options Suitability: excellent          │
│ Typical Hold: 3-10 weeks               │
│ ────────────────────────────────────── │
│ Universal across stocks and ETFs.       │
└─────────────────────────────────────────┘
  ↑ 8px wide green border
```

**Status:** ✅ FIX APPLIED (needs verification)

---

### Issue #5: Strategy Dropdown Missing Entries
**Date:** 2025-11-04
**Initial Symptom:** Strategy dropdown showed only 4 strategies for Healthcare (expectation was 12)
**Root Causes:**
- Packaged build shipped with an out-of-date `dist/policy.v1.json` that contained only four strategies
- UX requirement clarified: Screen 3 should surface the **top five** sector-approved strategies (not all 12) while still exposing yellow/red options where appropriate
**Fix Applied:**
- Introduced `strategy_helpers.go` to rank strategies per sector and return a curated top-five list (allowed + best alternates)
- Updated `ticker_entry.go` to consume the curated list and reuse helper colour palettes
- Copied the latest `data/policy.v1.json` into `dist/policy.v1.json` during build packaging
- Extended unit tests to expect five entries and verify a mix of `[GREEN]`, `[YELLOW]`, `[RED]`

**Verification:** Dev run now logs `DEBUG: Returning 5 strategy labels` with the intended mix (Alt10, Alt46, Alt26, Alt22, Alt43). Screenshot confirmed dropdown aligns with new requirement.

**Status:** ✅ RESOLVED

---

### Issue #6: Color Coding "Not Translating"
**Date:** 2025-11-04
**Symptom:** User reported that colour cues were only appearing once they reached Screen 3
**Resolution:**
- Added coloured strategy badges (same helper as Screen 3) to Screen 1 sector cards and Screen 2 screener view so the guidance is visible earlier in the flow
- Reworded info banner on Screen 3 to clarify it now lists “Top five strategies”

**Status:** ✅ RESOLVED

---

## Current State

### What Should Be Working:
✅ All 10 sectors have complete screener URLs
✅ Screen 2 should show 4 screener buttons for all sectors
✅ Screens 1-3 show consistent colour badges for top five strategies per sector
✅ Dropdown entries include `[GREEN]`/`[YELLOW]`/`[RED]` prefixes for curated list
✅ Strategy metadata cards retain colour-coded backgrounds and borders
✅ Warning banners and acknowledgement logic trigger for yellow/red picks

### What's NOT Working:
🚧 Need packaging checklist to ensure `dist/policy.v1.json` stays in sync before release (tracked separately)

---

## Sector-Specific Strategy Suitability

### Healthcare (Best Performance)
- **GREEN (excellent/good):** Alt10, Alt46, Alt43, Alt39, Alt28, Baseline
- **YELLOW (marginal):** Alt26, Alt45, Alt47
- **RED (incompatible):** Alt22, Alt9

### Technology
- **GREEN:** Alt26, Alt22, Alt10, Alt47, Alt45, Alt46, Baseline, Alt39
- **YELLOW:** Alt43, Alt28, Alt9

### Consumer Discretionary
- **GREEN:** Alt10, Alt26, Alt9, Baseline
- **YELLOW:** Alt22, Alt47, Alt43, Alt46

### Utilities (Worst Performance - 0% Success)
- **RED:** ALL 11 strategies (Alt10, Alt46, Alt26, Baseline, Alt22, Alt28, Alt43, Alt47, Alt45, Alt9, Alt39)

### Energy (Mean-Reverting)
- **RED:** ALL 10 strategies (Alt10, Alt26, Baseline, Alt22, Alt28, Alt43, Alt46, Alt47, Alt45, Alt9)

---

## Build Information

**Last Build:** 2025-11-04 09:25
**Executable:** tf-engine.exe (41 MB)
**Go Version:** 1.21+
**Policy Version:** 1.0.0
**Policy Signature:** edb876e707187327681ceaf651d2a1e9ca2145c773a36ac2437a8b3c4ff7f264

---

## Next Steps & Open Items

- Automate policy syncing during build packaging to avoid future `dist/policy.v1.json` drift
- Gather user feedback on the top-five presentation; expand/adjust ranking logic if any sector needs a different mix

---

## Related Files

- [PHASE6_PROGRESS.md](PHASE6_PROGRESS.md) - Phase 6 implementation progress
- [CLAUDE.md](CLAUDE.md) - Project architecture and rules
- [data/policy.v1.json](data/policy.v1.json) - Policy configuration file
- [internal/ui/screens/ticker_entry.go](internal/ui/screens/ticker_entry.go) - Strategy selection screen implementation

---

**Last Updated:** 2025-11-04 10:45
