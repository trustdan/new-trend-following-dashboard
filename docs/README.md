# Documentation Archive

This directory contains archived development documentation and historical progress files.

---

## Directory Structure

### `/archive/`

Contains historical completion reports, progress updates, and fix documentation from the development process.

**Contents:**
- **Phase Completion Files**: PHASE_0 through PHASE_6 completion reports
- **Screen Implementation Files**: Individual screen completion reports (Screens 1-3)
- **Fix Documentation**: Bug fixes, CI/CD improvements, installer builds
- **Test Reports**: Test coverage and results from development
- **Legacy Files**: Old project status and README files

---

## Archive Categories

### Development Phases

The application was built in phases:
- **Phase 0**: Project setup and planning
- **Phase 1**: Core architecture and state management
- **Phase 2**: Screen 1-2 implementation (Sector selection, Screener results)
- **Phase 3**: Screen 3 implementation (Ticker entry with Phase 6 updates)
- **Phase 4**: Screen 4 implementation (Checklist and cooldown)
- **Phase 5**: Screen 5-6 implementation (Position sizing, portfolio heat)
- **Phase 6**: Warning system implementation (replaced hard blocks)

### Bug Fixes and Improvements

- Button navigation fixes
- CI/CD race condition resolution
- Logging and GUI improvements
- Position sizing calculator fixes
- UI polish and screener updates

### Testing Documentation

- Integration testing reports
- Manual testing progress
- Test coverage analysis
- Test results summaries

### Installer Documentation

- Windows installer build guides (superseded by BUILD_GUIDE.md in root)
- Installer completion reports

---

## Why Files Are Archived

These files document the development journey but are no longer needed for day-to-day development. They've been archived to:

1. **Reduce clutter** in the root directory
2. **Preserve history** for reference if needed
3. **Keep active docs visible** (BUILD_GUIDE.md, CLAUDE.md, etc.)

---

## Active Documentation (In Root Directory)

For current development, refer to these files in the project root:

- **[README.md](../README.md)** - Project overview and getting started
- **[CLAUDE.md](../CLAUDE.md)** - AI assistant instructions and architecture
- **[architects-intent.md](../architects-intent.md)** - Original design intent
- **[architectural-overview.md](../architectural-overview.md)** - System architecture
- **[BUILD_GUIDE.md](../BUILD_GUIDE.md)** - Build and release documentation
- **[BUILD_QUICK_REFERENCE.txt](../BUILD_QUICK_REFERENCE.txt)** - Quick build cheat sheet
- **[CONTRIBUTING.md](../CONTRIBUTING.md)** - Contribution guidelines
- **[BEARISH_SCREENER_EVALUATION.md](../BEARISH_SCREENER_EVALUATION.md)** - Bearish screener analysis
- **[BEARISH_SCREENERS_FIXED.md](../BEARISH_SCREENERS_FIXED.md)** - Bearish screener fixes

---

## Searching the Archive

Use grep or file search to find specific topics:

```bash
# Find all references to "utilities sector"
grep -r "utilities" archive/

# Find test-related documentation
ls archive/ | grep -i test

# Find Phase 6 documentation
ls archive/ | grep PHASE6
```

---

## Archive Index

### Phase Completions (Chronological)
1. PHASE_0_COMPLETE.md
2. PHASE_1_COMPLETE.md
3. PHASE_2_PROGRESS.md
4. PHASE_2_COMPLETE.md
5. PHASE_3_COMPLETE.md
6. PHASE_4_COMPLETE.md
7. PHASE_5_COMPLETE.md
8. PHASE_5_COMPLETION_STATUS.md
9. PHASE_5_SUMMARY.md
10. PHASE6_PROGRESS.md
11. PHASE6_SUMMARY.md
12. PHASE6_TROUBLESHOOTING.md

### Screen Implementations
- SCREEN_1_COMPLETE.md (Sector selection)
- SCREEN_2_COMPLETE.md (Screener results)
- SCREEN_3_COMPLETE.md (Ticker entry with strategy selection)

### Bug Fixes
- BUTTON_FIX_COMPLETE.md
- CI_FIX_SUMMARY.md
- CI_RACE_CONDITION_FIX.md
- LOGGING_AND_GUI_FIX.md
- POSITION_SIZING_FIX_COMPLETE.md
- UI_POLISH_FIXES_COMPLETE.md

### Testing
- INTEGRATION_TESTING_COMPLETE.md
- MANUAL_TESTING_PROGRESS.md
- TEST_COVERAGE_REPORT.md
- TEST_RESULTS.md

### Installer
- INSTALLER_BUILD_COMPLETE.md
- INSTALLER_GUIDE.md (superseded by BUILD_GUIDE.md)
- WINDOWS_INSTALLER_COMPLETE.md

### Legacy
- PROJECT-STATUS.md
- README-DEV.md
- SCREENER_AND_NOTIFICATION_UPDATES.md

---

## Note on Build Documentation

The archived installer files (INSTALLER_GUIDE.md, WINDOWS_INSTALLER_COMPLETE.md) have been superseded by:
- **[BUILD_GUIDE.md](../BUILD_GUIDE.md)** (comprehensive)
- **[BUILD_QUICK_REFERENCE.txt](../BUILD_QUICK_REFERENCE.txt)** (cheat sheet)

Use the new build documentation in the root directory for current builds.

---

**Last Updated:** November 4, 2025
**Archive Created:** November 4, 2025
**Purpose:** Preserve development history while keeping root directory clean
