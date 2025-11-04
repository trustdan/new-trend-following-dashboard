# Manual Testing Progress Tracker

**Testing Date:** November 4, 2025
**Installer Version:** 1.0.0
**Tester:** Dan
**Status:** 🔄 In Progress

---

## Testing Overview

This document tracks progress through the 37-item manual testing checklist from [INSTALLER_BUILD_COMPLETE.md](INSTALLER_BUILD_COMPLETE.md).

**Test Categories:**
- Installation Tests (6 items)
- Application Launch Tests (3 items)
- Workflow Tests (3 items)
- Reinstallation Tests (2 items)
- Uninstallation Tests (2 items)
- Edge Case Tests (4 items)

**Current Progress:** 0/37 tests completed (0%)

---

## Installation Tests

### ✅ Test 1: Fresh Install Test
**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Passed | ❌ Failed

**Steps:**
1. [ ] Double-click `TFEngine-Setup-1.0.0.exe`
2. [ ] Accept UAC prompt (admin elevation)
3. [ ] Welcome screen appears
4. [ ] License agreement accepted
5. [ ] Installation directory shown (default: `C:\Program Files\TF-Engine`)
6. [ ] Click Install button
7. [ ] Choose "Yes" for desktop shortcut
8. [ ] Check "Launch TF-Engine" on finish screen
9. [ ] Click "Finish"
10. [ ] Application launches to Sector Selection screen

**Expected Result:** Application launches to Sector Selection screen
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Any issues, screenshots, or observations]_

---

### ✅ Test 2: Verify Installation
**Status:** ⬜ Not Started

**Files to Check in `C:\Program Files\TF-Engine`:**
- [ ] tf-engine.exe
- [ ] policy.v1.json
- [ ] README.md
- [ ] LICENSE.txt
- [ ] feature.flags.json
- [ ] data/ directory (empty)
- [ ] data/ui/ directory (empty)
- [ ] Uninstall.exe

**Start Menu Check:**
- [ ] Start → TF-Engine folder exists
- [ ] "TF-Engine" shortcut exists
- [ ] "Uninstall" shortcut exists

**Desktop Check:**
- [ ] "TF-Engine" shortcut exists

**Registry Check (regedit):**
- [ ] `HKLM\Software\TF-Engine\InstallPath` = installation directory
- [ ] `HKLM\Software\TF-Engine\Version` = "1.0.0"

**Add/Remove Programs Check:**
- [ ] Settings → Apps → Installed apps
- [ ] "TF-Engine" appears in list
- [ ] Shows version 1.0.0, publisher "TF Systems"

**Expected Result:** All files, shortcuts, and registry entries present
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Any missing files or incorrect paths]_

---

### ✅ Test 3: Feature Flags Verification
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Open `C:\Program Files\TF-Engine\feature.flags.json`
2. [ ] Verify contents match expected defaults:
   ```json
   {
     "trade_management": false,
     "sample_data_generator": false,
     "vimium_mode": false,
     "advanced_analytics": false
   }
   ```

**Expected Result:** All Phase 2 features are OFF by default
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Any unexpected flags]_

---

## Application Launch Tests

### ✅ Test 4: Launch from Start Menu
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Start → TF-Engine → TF-Engine
2. [ ] Application opens to Sector Selection screen
3. [ ] No errors in console/logs

**Expected Result:** Application opens normally
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Check logs/tf-engine_YYYY-MM-DD_HH-MM-SS.log]_

---

### ✅ Test 5: Launch from Desktop
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Double-click Desktop shortcut
2. [ ] Application opens normally

**Expected Result:** Application opens to Sector Selection screen
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Any shortcut errors]_

---

### ✅ Test 6: Launch from Executable
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Navigate to `C:\Program Files\TF-Engine`
2. [ ] Double-click tf-engine.exe
3. [ ] Application opens normally

**Expected Result:** Application opens to Sector Selection screen
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Any permission errors]_

---

## Workflow Tests

### ✅ Test 7: Complete One Trade Entry (End-to-End)
**Status:** ⬜ Not Started

**Steps:**
1. [ ] **Screen 1:** Select "Healthcare" sector
2. [ ] **Screen 2:** Click "Universe Screener" (opens browser)
3. [ ] **Screen 3:** Enter ticker "UNH", select strategy
4. [ ] **Screen 4:** Complete checklist (wait for 120-second cooldown)
5. [ ] **Screen 5:** Select poker sizing (e.g., "7 - Standard")
6. [ ] **Screen 6:** Pass heat check (should be under 4% portfolio, 1.5% sector)
7. [ ] **Screen 7:** Select options strategy (e.g., "Bull call spread")
8. [ ] **Screen 8:** View trade on calendar
9. [ ] Close application
10. [ ] Relaunch application
11. [ ] Trade should still appear on calendar (auto-save worked)

**Expected Result:** Trade persists across app restart
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Check data/ui/trades_in_progress.json for persistence]_

---

### ✅ Test 8: Sample Data Generation (Feature Flag)
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Navigate to Calendar screen
2. [ ] Verify "Generate Sample Data" button is **NOT visible** (feature flagged OFF)
3. [ ] Close application
4. [ ] Edit `C:\Program Files\TF-Engine\feature.flags.json`
5. [ ] Set `"sample_data_generator": true`
6. [ ] Save and restart app
7. [ ] Navigate to Calendar screen
8. [ ] Click "Generate Sample Data"
9. [ ] Confirm dialog → Yes
10. [ ] Calendar populates with 10 sample trades

**Expected Result:** Feature flag controls button visibility correctly
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Verify sample trades appear in calendar]_

---

### ✅ Test 9: Help System
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Click "?" icon in top-right corner on Sector Selection screen
2. [ ] Help dialog appears with context-sensitive content
3. [ ] Close help dialog
4. [ ] Navigate to Checklist screen
5. [ ] Click "?" icon again
6. [ ] Help dialog shows different content (checklist-specific)
7. [ ] Navigate to Calendar screen
8. [ ] Click "?" icon again
9. [ ] Help dialog shows calendar-specific content

**Expected Result:** Help content is context-sensitive
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Note any missing or incorrect help text]_

---

## Reinstallation Tests

### ✅ Test 10: Install Over Existing Version
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Run `TFEngine-Setup-1.0.0.exe` again (over existing installation)
2. [ ] Dialog appears: "TF-Engine is already installed (version 1.0.0). Continue?"
3. [ ] Click "Yes"
4. [ ] Installer completes
5. [ ] Application still works
6. [ ] Existing trade data is preserved (check calendar)

**Expected Result:** Reinstallation preserves user data
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Verify no data loss]_

---

### ✅ Test 11: Custom Directory Installation
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Uninstall current version (Start → TF-Engine → Uninstall)
2. [ ] Run installer again
3. [ ] Choose custom directory: `C:\MyApps\TF-Engine`
4. [ ] Complete installation
5. [ ] Verify files are in custom directory
6. [ ] Application launches normally from custom path

**Expected Result:** Application works from custom installation directory
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Check shortcuts point to correct custom path]_

---

## Uninstallation Tests

### ✅ Test 12: Uninstall via Start Menu (Preserve Data)
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Create sample trade data first (add at least one trade)
2. [ ] Start → TF-Engine → Uninstall
3. [ ] Confirm uninstallation
4. [ ] Dialog: "Delete trade history and settings?"
5. [ ] Click "No"
6. [ ] Verify `C:\Program Files\TF-Engine\data\` still exists
7. [ ] Verify all other files removed (tf-engine.exe, policy.v1.json, etc.)
8. [ ] Verify shortcuts removed from Start Menu and Desktop
9. [ ] Verify registry keys removed (regedit)

**Expected Result:** Data preserved, program files removed
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Check data directory contents]_

---

### ✅ Test 13: Uninstall via Add/Remove Programs (Delete Data)
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Reinstall application
2. [ ] Create sample trade data (add at least one trade)
3. [ ] Settings → Apps → Installed apps → TF-Engine → Uninstall
4. [ ] Dialog: "Delete trade history and settings?"
5. [ ] Click "Yes"
6. [ ] Verify `C:\Program Files\TF-Engine\` completely removed
7. [ ] Verify all shortcuts removed
8. [ ] Verify registry keys removed

**Expected Result:** Complete removal including user data
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Verify no orphaned files]_

---

## Edge Case Tests

### ✅ Test 14: Low Disk Space
**Status:** ⬜ Not Started | ⏭️ Skipped (Optional)

**Steps:**
1. [ ] Simulate low disk space (<10 MB free)
2. [ ] Run installer
3. [ ] Installer should warn about insufficient space

**Expected Result:** Installer detects and warns about low disk space
**Actual Result:** _[Fill in after testing, or mark Skipped]_
**Notes:** _[This test is optional and may be difficult to simulate]_

---

### ✅ Test 15: Network Drive Installation
**Status:** ⬜ Not Started | ⏭️ Skipped (Optional)

**Steps:**
1. [ ] Map network drive (e.g., Z:\)
2. [ ] Try installing to network drive
3. [ ] Application should work (but not recommended)

**Expected Result:** Installation works but may have performance issues
**Actual Result:** _[Fill in after testing, or mark Skipped]_
**Notes:** _[Network drive testing is optional]_

---

### ✅ Test 16: Non-Admin User
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Log in as standard user (non-administrator)
2. [ ] Run installer
3. [ ] UAC prompt should appear requesting admin password
4. [ ] Enter admin credentials
5. [ ] Installation should complete successfully

**Expected Result:** Installer prompts for admin elevation
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Verify standard users can install with admin password]_

---

### ✅ Test 17: Antivirus Scan
**Status:** ⬜ Not Started

**Steps:**
1. [ ] Scan `TFEngine-Setup-1.0.0.exe` with Windows Defender
2. [ ] Right-click installer → Scan with Windows Defender
3. [ ] Should pass with no threats detected

**Expected Result:** No false positives from antivirus
**Actual Result:** _[Fill in after testing]_
**Notes:** _[Note any antivirus warnings]_

---

## Bug Tracking

**Bugs Found During Testing:**

### Bug #1: [Title]
- **Severity:** Critical | High | Medium | Low
- **Found In:** [Test number/name]
- **Description:** [What happened]
- **Steps to Reproduce:**
  1. [Step 1]
  2. [Step 2]
- **Expected:** [What should happen]
- **Actual:** [What actually happened]
- **Status:** Open | Fixed | Won't Fix

---

## Test Summary

**Date Completed:** _[Fill in when finished]_

**Overall Status:**
- ✅ Passed: 0/17 tests
- ❌ Failed: 0/17 tests
- ⏭️ Skipped: 0/17 tests (edge cases optional)

**Critical Issues:** _[List any blocking bugs]_

**Recommendation:**
- [ ] Ready for beta testing
- [ ] Needs bug fixes before beta
- [ ] Major issues found - rebuild required

---

## Next Steps After Testing

### If All Tests Pass ✅
1. Mark installer as "Ready for Beta Testing"
2. Recruit 3-5 beta testers
3. Distribute installer + testing instructions
4. Monitor feedback

### If Critical Bugs Found ❌
1. Document all bugs in GitHub issues
2. Prioritize fixes (P0/P1/P2)
3. Fix critical bugs
4. Rebuild installer
5. Retest failed tests

---

**Last Updated:** November 4, 2025
**Next Review:** After completing all tests