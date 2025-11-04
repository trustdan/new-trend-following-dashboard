# TF-Engine GUI Development Roadmap

**TF = Trend Following** - Systematic Donchian breakout trading system

**⚠️ IMPORTANT: Technology Change - Fyne Chosen Over Svelte**

This roadmap was originally written for a Svelte + Go architecture. **The project has since adopted Fyne (native Go GUI) for the frontend.** Many steps below reference Svelte, HTTP APIs, and web-based UI - these should be interpreted in terms of Fyne equivalents:
- "SvelteKit component" → "Fyne screen/widget"
- "HTTP API endpoint" → "Direct function call to backend"
- "Svelte store" → "AppState struct with callbacks"
- "CSS/Tailwind" → "Fyne theme and custom widgets"

**Current Status:** Phase 0 completed with Fyne. All 6 main tabs implemented, custom British Racing Green theme working, VIM mode implemented, cross-compilation working.

**Purpose:** This roadmap serves as the wireframe for all phase-specific step documents that will guide the implementation of the TF-Engine GUI application.

**Structure:** Each phase is broken down into discrete step documents (e.g., `phase0-step1.md`, `phase0-step2.md`, etc.) that can be executed sequentially.

**Timeline:** 12 weeks (3 months) from start to production-ready application

**Approach:** Build incrementally, test continuously, maintain documentation daily (per `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md`)

**Logging Philosophy:** Comprehensive logging from day one. Every feature, API call, user action, and component lifecycle event should be logged. This enables:
- **Debugging:** Detailed logs make troubleshooting easier
- **Feature Evaluation:** Track which features are used, which cause problems
- **Performance Monitoring:** Identify slow operations and bottlenecks
- **Pruning Decisions:** Remove features that cause headaches (data-driven decisions)
- **User Behavior:** Understand actual workflow vs intended workflow

---

## Table of Contents

- [Phase 0: Foundation & Proof-of-Concept](#phase-0-foundation--proof-of-concept) (Week 1-2) - 4 steps
- [Phase 1: Dashboard & FINVIZ Scanner](#phase-1-dashboard--finviz-scanner) (Week 3-4) - 5 steps
- [Phase 2: Checklist & Position Sizing](#phase-2-checklist--position-sizing) (Week 5-6) - 5 steps
- [Phase 3: Heat Check & Trade Entry](#phase-3-heat-check--trade-entry) (Week 7-8) - 5 steps
- [Phase 4: Calendar & Polish](#phase-4-calendar--polish) (Week 9-10) - 4 steps
- [Phase 5: Testing & Packaging](#phase-5-testing--packaging) (Week 11-12) - 5 steps

**Total Steps:** 28 detailed implementation documents

---

## Phase 0: Foundation & Proof-of-Concept

**Duration:** Week 1-2 (10 days)
**Goal:** Validate technology choices and establish build pipeline
**Success:** Working POCs + chosen tech stack + automated build process

### Step 1: Development Environment Setup

**File:** `phase0-step1-dev-environment.md`
**Duration:** 1 day
**Dependencies:** None

Set up the complete development environment on Linux (WSL2/Kali). Install Go 1.24+, verify the existing tf-engine backend compiles and tests pass. Install development tools (VSCode/Cursor, Go extensions). Install Fyne prerequisites (X11 libraries for development). Configure the workspace according to `1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md` (no Git in Linux). Verify cross-compilation to Windows works with a simple "Hello World" Go binary. Set up logging directories (`logs/`) and configure log rotation. Document all versions and paths for reproducibility.

**Deliverables:**
- Go environment configured and verified
- Fyne prerequisites installed (X11 development libraries)
- Backend compiles successfully (`go build` and `go test ./...` pass)
- Cross-compilation to Windows verified
- VSCode/Cursor workspace configured
- Logging directories created (`logs/`)
- Documentation: Environment setup guide

---

### Step 2: Fyne Proof-of-Concept

**File:** `phase0-step2-fyne-poc.md`
**Duration:** 2-3 days
**Dependencies:** Step 1

Build a minimal desktop GUI application using Fyne to validate the "pure Go desktop app" approach. Create a simple window that displays current settings from the SQLite database (equity, risk%, caps). Add "Refresh" and "Update" buttons that call the existing backend domain logic directly (no HTTP). Test embedding resources. Cross-compile to Windows .exe and verify it runs without dependencies. This POC proves we CAN build a desktop app, giving us a fallback option if the Svelte approach fails.

**Deliverables:**
- Fyne app displays data from trading.db
- Buttons call backend functions directly (in-process)
- Cross-compiles to Windows .exe
- Runs on Windows without installation
- Documentation: Fyne POC results, pros/cons analysis

---

### Step 3: Custom Theme Implementation

**File:** `phase0-step3-custom-theme.md`
**Duration:** 1-2 days
**Dependencies:** Step 2

Implement custom Fyne theme with British Racing Green accent color (#00352B). Create both day mode and dark mode color schemes. Define custom fonts, spacing, and widget styles. Implement theme toggle functionality with state persistence. Test theme on all standard Fyne widgets to ensure consistency. This step establishes the visual identity of the application.

**Deliverables:**
- Custom Fyne theme with British Racing Green
- Day mode color scheme defined
- Dark mode color scheme defined
- Theme toggle button working
- Theme preference persists between sessions
- All widgets use custom theme
- Documentation: Theme implementation guide

---

### Step 4: Build Pipeline & Cross-Compilation

**File:** `phase0-step4-build-pipeline.md`
**Duration:** 1 day
**Dependencies:** Steps 2, 3

Establish the production build pipeline for Fyne. Create scripts: `build-windows.sh` (cross-compiles Linux → Windows .exe with Fyne bundling), and `export-for-windows.sh` (creates zip for Windows handoff per RULES.md). Create GUI project structure (`internal/gui/screens/`, `internal/gui/widgets/`, `internal/gui/state/`, `internal/gui/theme/`). Test the complete build pipeline end-to-end: build on Linux, transfer to Windows, verify execution. Configure logging for GUI events. Document the build process.

**Deliverables:**
- Build scripts created and tested (Linux → Windows cross-compile)
- GUI project structure established in `internal/gui/`
- Cross-compilation verified working
- Export-to-Windows workflow validated
- Logging configured for GUI events
- Documentation: Build pipeline guide, Fyne architecture decisions

---

## Phase 1: Dashboard & FINVIZ Scanner

**Duration:** Week 3-4 (10 days)
**Goal:** Build core navigation and implement daily workflow starting point
**Success:** User can scan FINVIZ, import candidates, view dashboard

### Step 5: Tab Navigation & State Management

**File:** `phase1-step5-navigation-state.md`
**Duration:** 2 days
**Dependencies:** Phase 0 complete

Implement the main tab navigation structure using Fyne's AppTabs widget. Create 6 main tabs: Dashboard, Checklist, Position Sizing, Heat Check, Trade Entry, Calendar. Set up centralized application state management using a Go struct (`AppState`) that holds database connection, current settings, positions, candidates, and UI state. Implement state update callbacks to refresh UI when data changes. **Add comprehensive logging:** log every tab switch, state update, database operation, and user action with timestamps and context. Log to `logs/tf-gui.log` with rotation. Create helper functions for common state operations (LoadSettings, LoadPositions, LoadCandidates). Test navigation between all tabs.

**Deliverables:**
- Tab navigation working (6 main tabs using Fyne AppTabs)
- AppState struct defined and initialized with database connection
- State management callbacks implemented
- **Comprehensive logging of tab switches and state changes**
- **Performance metrics logged (database query times)**
- **Log rotation configured (daily)**
- Unit tests for state management
- Documentation: State management guide, logging format specification

---

### Step 6: Application Layout & Navigation

**File:** `phase1-step6-layout-navigation.md`
**Duration:** 2 days
**Dependencies:** Step 5

Build the application shell that will house all screens. Create the main `App.svelte` component with routing (SvelteKit's file-based routing). Implement `Header.svelte` with app title, theme toggle button (sun/moon icon), and settings icon. Create `Navigation.svelte` as a sidebar with links to Dashboard, Scanner, Checklist, Position Sizing, Heat Check, Trade Entry, and Calendar. Style the navigation with gradient accents. Implement the theme toggle functionality that switches between day/night mode with smooth 0.3s transitions. Ensure theme preference persists to localStorage. Create placeholder routes (`+page.svelte` files) for each main screen. Test navigation between screens with smooth transitions. **Add frontend logging:** Log every navigation event (from screen, to screen, timestamp), theme toggle, and component mount/unmount. Create a logging utility (`lib/logger.ts`) that logs to browser console with color-coding and timestamps. Track navigation patterns to understand user flow.

**Deliverables:**
- App layout with header and sidebar
- Theme toggle working (day/night mode)
- Navigation to all main screens
- Smooth page transitions
- Theme persistence to localStorage
- Responsive layout (desktop-first)
- **Frontend logging utility (`lib/logger.ts`)**
- **Navigation and theme events logged**
- Documentation: Component structure guide, logging conventions

---

### Step 7: Dashboard Screen

**File:** `phase1-step7-dashboard.md`
**Duration:** 2 days
**Dependencies:** Step 6

Build the Dashboard as the main landing page. Display key information: portfolio summary (equity, total heat, available capacity), list of open positions (ticker, entry, current stop, risk amount, days held), today's candidates count, and any active cooldowns. Fetch data from `GET /api/positions` and `GET /api/candidates`. Create reusable components: `Card.svelte` for content containers, `Badge.svelte` for status indicators, `PositionTable.svelte` for displaying positions. Implement the large banner at the top showing current overall system status (uses the same gradient banner as checklist, but shows portfolio health). Add "Quick Actions" section with buttons to navigate to Scanner or Checklist. Style everything with gradients, shadows, and the established design system. Implement loading states with animated gradient spinners.

**Deliverables:**
- Dashboard displays real data from backend
- Portfolio summary cards with gradients
- Position table with alternating row colors
- Large status banner at top
- Quick action buttons
- Loading states with spinners
- Documentation: Dashboard component guide

---

### Step 8: FINVIZ Scanner Implementation

**File:** `phase1-step8-finviz-scanner.md`
**Duration:** 2-3 days
**Dependencies:** Step 7

Create the FINVIZ Scanner screen that allows one-click daily scanning. Build `FINVIZScanner.svelte` with a large "Run Daily Scan" button (gradient background). When clicked, call `POST /api/candidates/scan` with the preset name (e.g., "TF_BREAKOUT_LONG"). Display the scan results in a table: ticker, last close, volume, sector. The backend already has the FINVIZ scraper, so this is primarily a UI task. Show loading state during scan (can take 3-5 seconds). Display success message with count of candidates found. Handle errors gracefully (network issues, FINVIZ changes). Allow user to review results before importing. Add a "Preset Manager" section where user can view/select different FINVIZ presets. Style the scan button prominently with a gradient and hover effect.

**Deliverables:**
- FINVIZ scanner screen with prominent button
- Scan results table
- Loading state during scan
- Success/error notifications
- Preset selector (if multiple presets exist)
- Documentation: Scanner workflow guide

---

### Step 9: Candidate Import & Review

**File:** `phase1-step9-candidate-import.md`
**Duration:** 2 days
**Dependencies:** Step 8

Build the candidate import workflow. After scanning, display candidates in a table with checkboxes for selection. Show ticker, sector, last close, volume. Highlight any tickers that are on cooldown (grayed out, not selectable). Add "Select All" / "Deselect All" buttons. Display sector distribution (e.g., "3 Tech/Comm, 2 Energy, 1 Finance"). The large "Import Selected" button should be gradient-styled and disabled until at least one ticker is selected. When clicked, call `POST /api/candidates/import` with selected tickers and today's date. Show success notification: "12 candidates imported for 2025-10-29". Update the Dashboard to reflect new candidates. Create the `CandidateList.svelte` component for reuse elsewhere. Implement filters (by sector, minimum volume, etc.).

**Deliverables:**
- Candidate review table with checkboxes
- Cooldown indicators (grayed out)
- Sector distribution summary
- Import button with validation
- Success notification
- Dashboard updates after import
- Documentation: Candidate import workflow

---

## Phase 2: Checklist & Position Sizing

**Duration:** Week 5-6 (10 days)
**Goal:** Implement the 5 gates checklist and position sizing calculator
**Success:** Banner changes color based on checklist, sizing calculates correctly

### Step 10: Banner Component

**File:** `phase2-step10-banner-component.md`
**Duration:** 2 days
**Dependencies:** Phase 1 complete

Build the centerpiece of the application: the large gradient banner. Create `Banner.svelte` as a highly reusable component that displays RED/YELLOW/GREEN states with smooth gradient backgrounds. The banner should be minimum 20% of viewport height (at least 150px), centered, with rounded corners and a soft glow effect. Implement three gradient backgrounds: RED (#DC2626 → #991B1B), YELLOW (#F59E0B → #FBBF24), GREEN (#10B981 → #059669). The banner should accept props: `state` ('RED' | 'YELLOW' | 'GREEN'), `message` (main text), and `details` (subtext). Add smooth 0.3s transition animations when state changes. Include a subtle pulse effect on state transitions. The text should be large (36px), bold, white, with a subtle shadow. Add appropriate icons (🛑 for RED, ⚠️ for YELLOW, ✓ for GREEN). Make the banner responsive but always prominent.

**Deliverables:**
- Banner component with three gradient states
- Smooth transition animations (0.3s ease-in-out)
- Pulse effect on state change
- Customizable via props
- Responsive sizing (min 150px height)
- Documentation: Banner component API and usage guide

---

### Step 11: Checklist Form & Required Gates

**File:** `phase2-step11-checklist-form.md`
**Duration:** 2-3 days
**Dependencies:** Step 10

Create the Checklist screen with the large banner at the top. Build a form with input fields: ticker (text), entry price (number), ATR/N (number), sector (dropdown from predefined list), and structure (dropdown: "Stock", "Call", "Put", "Call Spread", "Put Spread"). Below the inputs, display the 5 required gates as large checkboxes with custom styling (gradient when checked). Each gate should have clear text: "✓ Signal: 55-bar Donchian breakout confirmed", "✓ Risk/Size: 2×N stop, 0.5×N adds, max 4 units", "✓ Liquidity: Avg volume >1M shares OR options OI >100", "✓ Exits: 10-bar Donchian OR 2×N stop", "✓ Behavior: Not on cooldown, heat OK, 2-min timer honored". As the user checks/unchecks boxes, the banner should update live. If any required gate is unchecked, banner is RED. If all required but quality score is low, banner is YELLOW. When banner turns GREEN, start the 2-minute timer (displayed prominently). Add a "Save Evaluation" button that's only enabled when banner is GREEN. **Log all checklist interactions:** each checkbox change (which gate, checked/unchecked, timestamp), banner state transitions (RED→YELLOW→GREEN with reasons), form field changes, save attempts. Track how often each gate is the bottleneck (last one to be checked).

**Deliverables:**
- Checklist form with all input fields
- 5 required gates as styled checkboxes
- Banner updates live as checkboxes change
- Sector dropdown populated from backend
- Structure dropdown for trade type
- Real-time banner state updates
- **Detailed logging of all checklist interactions and banner transitions**
- Documentation: Checklist workflow guide, interaction logs format

---

### Step 12: Quality Items & Scoring

**File:** `phase2-step12-quality-scoring.md`
**Duration:** 1-2 days
**Dependencies:** Step 11

Add the optional quality items section below the required gates. Create 4 optional checkboxes: "Regime OK (SPY > 200 SMA for longs)", "No Chase (entry within 2N of 20-EMA)", "Earnings OK (no earnings within next 2 weeks)", and "Journal Note (why this trade now?)". The journal note should be a textarea where the user can write their reasoning. Each optional checkbox adds 1 point to the quality score. Display the score prominently: "Quality Score: 3 / 4". Also display the threshold (configurable in settings, default 3.0). The banner should transition from YELLOW to GREEN when score ≥ threshold (assuming all required gates pass). Store the quality score in the checklist evaluation. The quality items should be styled similarly to required gates but with a different visual treatment (perhaps a blue accent instead of red to distinguish them).

**Deliverables:**
- 4 optional quality checkboxes
- Journal note textarea
- Quality score calculation and display
- Threshold comparison
- Banner updates based on score
- Visual distinction between required and optional
- Documentation: Quality scoring logic

---

### Step 13: Position Sizing Calculator

**File:** `phase2-step13-position-sizing.md`
**Duration:** 2-3 days
**Dependencies:** Step 12

Build the Position Sizing screen with a form pre-filled from the checklist data (ticker, entry, ATR, sector). Add input fields for method (dropdown: "stock", "opt-delta-atr", "opt-contracts"), K multiple (number, default 2.0), max units (number, default 4). Display current settings from the backend (equity, risk% per unit). When the user clicks "Calculate Position Size", call `POST /api/size/calculate` with all parameters. Display the results prominently in a card with gradient border: shares/contracts per unit, risk $ per unit, initial stop price, and the add-on schedule (Add1, Add2, Add3 prices based on 0.5N increments). Show the full pyramid plan visually. For options trades, display additional fields: delta (for opt-delta-atr method) or contracts and debit (for opt-contracts method). Add warnings if the position is too concentrated (>25% of equity). Style the results card with gradients and make the numbers large and obvious. Add a "Save Position Plan" button that saves to the database and links to the checklist evaluation.

**Deliverables:**
- Position sizing form with all methods
- Integration with backend sizing API
- Results display with add-on schedule
- Warnings for concentration risk
- Separate handling for stock vs options
- Save button to persist plan
- Documentation: Position sizing methodology

---

### Step 14: 2-Minute Cool-Off Timer

**File:** `phase2-step14-cooloff-timer.md`
**Duration:** 1-2 days
**Dependencies:** Step 13

Implement the 2-minute impulse brake timer. When the user clicks "Save Evaluation" on the Checklist screen (only enabled when banner is GREEN), record the evaluation timestamp in the database via `POST /api/checklist/evaluate`. Start a countdown timer displayed prominently on the screen: "Cool-off period: 2:00 remaining". The timer should count down every second (2:00, 1:59, 1:58, ..., 0:01, 0:00). While the timer is active, certain actions are disabled (specifically, the "SAVE GO DECISION" button on the Trade Entry screen). The timer display should be styled with a gradient accent and large text. When the timer reaches 0:00, show a subtle notification: "Cool-off complete. You may proceed to trade entry." Store the timer state in a Svelte store so it persists across screen navigation. The backend will validate the 2-minute elapsed time when gates are checked, but the UI timer provides user feedback.

**Deliverables:**
- 2-minute countdown timer UI
- Timer starts on "Save Evaluation"
- Timer persists across screen navigation (Svelte store)
- Visual feedback when timer completes
- Integration with checklist evaluation endpoint
- Documentation: Cool-off timer workflow

---

## Phase 3: Heat Check & Trade Entry

**Duration:** Week 7-8 (10 days)
**Goal:** Implement heat cap validation and final 5-gate check
**Success:** Cannot save GO decision if gates fail or caps exceeded

### Step 15: Heat Check Screen

**File:** `phase3-step15-heat-check.md`
**Duration:** 2-3 days
**Dependencies:** Phase 2 complete

Build the Heat Check screen to validate that the proposed trade won't exceed portfolio or sector bucket caps. Display the current portfolio heat as a visual gauge (progress bar or circular gauge with gradient fill showing % of 4% cap). Display all sector buckets with their current heat levels. When the user clicks "Check Heat for This Trade", call `POST /api/heat/check` with the ticker, risk amount (from position sizing), and sector bucket. The backend calculates: current portfolio heat, new portfolio heat (current + proposed), portfolio cap (equity × 4%), current bucket heat, new bucket heat, and bucket cap (equity × 1.5%). Display the results in a clear card. If either cap would be exceeded, show a large RED warning with gradient background: "⚠️ PORTFOLIO CAP EXCEEDED" or "⚠️ BUCKET CAP EXCEEDED". Show the exact overage amount. Provide suggestions: "Reduce position size to X shares" or "Close existing position in this sector". If caps are OK, show a GREEN success message. Add a "Calculate Max Shares" button that determines the maximum position size that fits within the caps.

**Deliverables:**
- Heat check screen with visual gauges
- Current portfolio and bucket heat display
- Heat calculation for proposed trade
- RED warnings for cap violations
- Suggestions to resolve cap issues
- "Calculate Max Shares" helper
- Documentation: Heat management guide

---

### Step 16: Trade Entry Screen & Summary

**File:** `phase3-step16-trade-entry.md`
**Duration:** 2 days
**Dependencies:** Step 15

Create the Trade Entry screen where the final decision is made. This screen should be visually impressive and serious - it's the culmination of the entire workflow. Display a comprehensive trade summary card with gradient border showing: ticker, direction (LONG/SHORT), entry price, shares/contracts per unit, max units, initial stop, risk $ per unit, total max risk (if all units added), add-on schedule, sector bucket, exit plan (10-bar Donchian OR 2×N), and quality score. Below the summary, add a large section titled "Final Gate Check" with a prominent "RUN FINAL GATE CHECK" button (gradient). When clicked, this will trigger the 5-gate validation (implemented in the next step). Reserve space for displaying the gate results. At the bottom, two large action buttons: "SAVE GO DECISION" (green gradient, initially disabled) and "SAVE NO-GO DECISION" (red gradient, always enabled for journaling rejections).

**Deliverables:**
- Trade entry screen with complete summary
- Trade plan displayed prominently
- "Run Final Gate Check" button
- Space for gate results display
- GO/NO-GO decision buttons (styled)
- Documentation: Trade entry workflow

---

### Step 17: 5 Gates Validation

**File:** `phase3-step17-gates-validation.md`
**Duration:** 2-3 days
**Dependencies:** Step 16

Implement the heart of the discipline enforcement system: the 5 gates check. When the user clicks "RUN FINAL GATE CHECK", call `POST /api/gates/check` with all trade data. The backend validates each gate and returns pass/fail for each. Display the results in a visually clear format with icons and gradient accents:

**Gate 1: Banner Status** - ✓ PASS (GREEN) or ✗ FAIL (not green)
**Gate 2: Impulse Brake** - ✓ PASS (2+ min elapsed) or ✗ FAIL (timer still active)
**Gate 3: Cooldown Status** - ✓ PASS (not on cooldown) or ✗ FAIL (ticker or bucket on cooldown)
**Gate 4: Heat Caps** - ✓ PASS (within caps) or ✗ FAIL (exceeds portfolio or bucket cap)
**Gate 5: Sizing Completed** - ✓ PASS (plan saved) or ✗ FAIL (sizing not done)

If ALL gates pass, show a large "ALL GATES PASS ✓" message with green gradient and enable the "SAVE GO DECISION" button. If ANY gate fails, show "GATES FAILED ✗" with red gradient and keep "SAVE GO DECISION" disabled. The "SAVE NO-GO DECISION" button is always enabled to allow journaling of rejected trades. **Log every gate check:** which gates passed/failed, reasons for failure, timing data (how long since checklist eval, current heat levels), and whether user proceeds to save decision. Track gate failure patterns to identify if specific gates are too strict or causing frustration. Backend should log detailed gate validation logic and calculations.

**Deliverables:**
- 5 gates validation implementation
- Visual display of each gate result
- Clear pass/fail indicators
- "SAVE GO DECISION" enabled only if all pass
- Integration with backend gates API
- **Comprehensive logging of all gate checks and results**
- **Gate failure pattern tracking for feature evaluation**
- Documentation: Gates enforcement logic, gate check logs format

---

### Step 18: Decision Saving & Journaling

**File:** `phase3-step18-decision-saving.md`
**Duration:** 1-2 days
**Dependencies:** Step 17

Implement the decision saving functionality. When "SAVE GO DECISION" is clicked (only enabled if all gates pass), call `POST /api/decisions` with the full trade record: timestamp, ticker, sector, direction, entry, shares, max units, risk $ per unit, total max risk, initial stop, add-on prices, structure, quality score, journal note, banner status, and all gate results. Show a success notification with green gradient: "✓ GO decision saved for AAPL". The decision should appear in a "Ready to Execute" list (could be on Dashboard or a separate screen). When "SAVE NO-GO DECISION" is clicked, show a modal asking for the rejection reason. Save the NO-GO decision with all the same data plus the reason. This is important for journaling - tracking what you DON'T trade is as valuable as tracking what you DO. Clear the current checklist/sizing data and reset to allow starting a new analysis. Add visual feedback: subtle confetti animation on successful GO decision (optional but fun).

**Deliverables:**
- GO decision saving with full data
- NO-GO decision with reason modal
- Success notifications
- Ready to Execute list
- Reset workflow for next trade
- Optional: Confetti animation on GO decision
- Documentation: Decision workflow guide

---

### Step 19: Integration Testing for Core Workflow

**File:** `phase3-step19-integration-testing.md`
**Duration:** 1-2 days
**Dependencies:** Steps 15-18

Perform comprehensive integration testing of the complete workflow from FINVIZ scan to decision saving. Create test scenarios covering the happy path (all gates pass) and multiple failure paths (each gate failing individually). Test the UI behavior: banner transitions, timer countdown, heat warnings, gate validation, decision saving. Verify data persistence - close and reopen the app, ensure candidates, evaluations, and decisions are all saved. Test edge cases: exactly at heat cap, zero candidates found, network errors during scan, etc. Fix any bugs discovered. Update component integration tests. Ensure all API calls handle errors gracefully. Test the complete workflow on Linux and prepare for Windows testing (export zip). Document any known issues or limitations.

**Deliverables:**
- Test scenarios document
- Integration tests passing
- Bug fixes for issues found
- Error handling verification
- Persistence testing
- Documentation: Testing results, known issues

---

## Phase 4: Calendar & Polish

**Duration:** Week 9-10 (10 days)
**Goal:** Add calendar view and polish the entire application
**Success:** Professional, production-ready UI with no obvious issues

### Step 20: Calendar View Implementation

**File:** `phase4-step20-calendar-view.md`
**Duration:** 3-4 days
**Dependencies:** Phase 3 complete

Build the Calendar screen showing a rolling 10-week view (2 weeks back + 8 weeks forward) with sectors as rows and weeks as columns. Each cell represents a sector × week combination and displays tickers that are active in that bucket during that week. Create `Calendar.svelte` as the main component and `CalendarCell.svelte` for individual cells. Calculate week boundaries (Monday-Sunday). Fetch data from `GET /api/calendar?weeks=10` which returns positions grouped by sector and week. Color-code cells based on heat level: low heat (green tint), medium (yellow tint), high (red tint). Show ticker symbols in each cell, stacked if multiple positions. Add tooltips on hover showing position details (entry date, risk amount, expected exit). Highlight the current week. Make the calendar horizontally scrollable if needed. Add a legend explaining the color coding. The calendar provides an at-a-glance view of sector diversification and helps identify basket crowding.

**Deliverables:**
- Calendar grid (10 weeks × N sectors)
- Positions displayed in correct cells
- Heat-based color coding
- Tooltips with position details
- Current week highlight
- Responsive layout with scroll
- Documentation: Calendar view guide

---

### Step 21: TradingView Integration

**File:** `phase4-step21-tradingview-integration.md`
**Duration:** 1 day
**Dependencies:** Step 20

Implement the TradingView integration for quick chart access. Add "Open in TradingView" buttons next to tickers throughout the app (Candidates list, Dashboard positions, Checklist form). When clicked, construct a TradingView URL: `https://www.tradingview.com/chart/?symbol={ticker}` and open it in a new browser tab. Create a `TradingViewLink.svelte` component that's reusable. Style the button with a subtle gradient and the TradingView logo/icon. Add a settings option to customize the TradingView URL template (some users may have different chart layouts or scripts they want to load automatically). In the documentation, explain how to set up the Ed-Seykota.pine script as a default on TradingView so it loads automatically. This integration makes the workflow seamless: see candidate → open chart → verify signal → return to app → fill checklist.

**Deliverables:**
- TradingView link component
- Links added to Candidates, Dashboard, Checklist
- URL opens in new tab
- Customizable URL template in settings
- Documentation: TradingView integration setup guide

---

### Step 22: UI Polish & Refinements

**File:** `phase4-step22-ui-polish.md`
**Duration:** 3-4 days
**Dependencies:** Steps 20-21

Polish the entire application to production quality. Review every screen for visual consistency, proper spacing, alignment, and gradient usage. Ensure all buttons have proper hover effects (lift + shadow). Add micro-interactions: checkboxes animate when checked, success messages slide in from top, error messages shake slightly. Implement keyboard shortcuts: Escape closes modals, Tab navigation works properly, Enter submits forms. Add loading skeletons (animated placeholders) instead of just spinners. Improve error messages to be more helpful with specific suggestions. Add tooltips to explain fields and features. Ensure all text is readable (proper contrast) in both day and night modes. Fix any layout issues on different screen sizes. Add breadcrumb navigation at the top of each screen. Test the theme toggle thoroughly - ensure all components update correctly. Fix any visual glitches or jarring transitions. **Add a debug panel (dev mode only):** View recent logs, clear logs, export logs to file. Add performance monitoring overlay (optional, togglable) showing API call times and render performance. **Review logs from previous steps to identify any features causing frequent errors or user frustration** - prepare list of candidates for removal or simplification.

**Deliverables:**
- Visual consistency across all screens
- Micro-interactions and animations
- Keyboard shortcuts implemented
- Loading skeletons
- Improved error messages
- Tooltips on complex fields
- Breadcrumb navigation
- Theme toggle refinements
- **Debug panel for viewing/exporting logs (dev mode)**
- **Performance monitoring overlay (optional)**
- **Feature evaluation report based on logs (which features problematic?)**
- Documentation: UI/UX improvements list, debug panel usage, features to reconsider

---

### Step 23: Performance Optimization

**File:** `phase4-step23-performance.md`
**Duration:** 1-2 days
**Dependencies:** Step 22

Optimize the application for performance. **Review performance logs collected in previous steps:** Identify slow API endpoints, sluggish UI components, long database queries. Profile the Svelte app to identify any slow components or unnecessary re-renders. Implement lazy loading for routes that aren't visited immediately. Optimize API calls: implement caching where appropriate (e.g., settings, candidate lists that don't change often). Reduce bundle size: analyze the built static files and remove unused dependencies. Ensure the Go backend responds quickly (<100ms for calculations, <500ms for database queries). Add pagination if candidate lists or decision logs get very long. Test on lower-end hardware to ensure it runs smoothly. Optimize images and icons (use SVGs where possible, compress PNGs). Measure and document load times, interaction responsiveness, and time to interactive. Set performance budgets for future development. **Use logged performance metrics to prioritize optimization efforts (data-driven decisions).**

**Deliverables:**
- **Performance analysis from logs (slowest operations identified)**
- Performance profiling results
- Lazy loading implemented
- API caching strategy
- Bundle size optimized
- Database query optimization
- Performance metrics documented
- **Performance improvement report (before/after metrics)**
- Documentation: Performance benchmarks, optimization guide

---

## Phase 5: Testing & Packaging

**Duration:** Week 11-12 (10 days)
**Goal:** Comprehensive testing and production-ready Windows installer
**Success:** Zero critical bugs, smooth installation, complete documentation

### Step 24: Comprehensive Testing Suite

**File:** `phase5-step24-comprehensive-testing.md`
**Duration:** 3 days
**Dependencies:** Phase 4 complete

Create and execute a comprehensive testing suite. Write unit tests for critical Svelte components (Banner, Checklist, PositionSizer, etc.) using Vitest. Write integration tests for the complete workflow: scan → import → checklist → sizing → heat → gates → decision. Test all edge cases: empty databases, maximum heat, failed scans, network errors, invalid inputs, boundary conditions (exactly at caps, zero risk, etc.). Create a test plan document with all scenarios and expected outcomes. Execute manual testing of the complete workflow multiple times, varying parameters. Test both day and night modes thoroughly. Test keyboard navigation and accessibility. Run the Go backend tests (`go test ./...`) to ensure no regressions. Test on Windows (export zip, copy to Windows, run .exe, verify all functionality). Log all bugs and issues in a tracking document.

**Deliverables:**
- Frontend unit tests (Vitest)
- Integration test suite
- Manual test plan executed
- Edge case testing complete
- Backend regression tests passing
- Windows functionality verified
- Bug tracking document
- Documentation: Testing report

---

### Step 25: Bug Fixing Sprint

**File:** `phase5-step25-bug-fixing.md`
**Duration:** 2-3 days
**Dependencies:** Step 24

Address all bugs and issues discovered during testing. Prioritize by severity: critical (app crashes, data loss, gates not enforced) → high (incorrect calculations, UI broken) → medium (cosmetic issues, minor UX problems) → low (nice-to-haves). Fix critical and high severity bugs first. For each bug: reproduce it, identify root cause, implement fix, verify fix, add regression test. Update documentation if the fix changes behavior. Test fixes on both Linux and Windows. If any bugs can't be fixed immediately, document them in KNOWN_LIMITATIONS.md with workarounds. Aim for zero critical bugs and minimal high-severity bugs before proceeding to packaging. Re-run the test suite after fixes to ensure no new regressions.

**Deliverables:**
- All critical bugs fixed
- High-severity bugs addressed
- Bug fix regression tests
- Updated KNOWN_LIMITATIONS.md
- Verification testing passed
- Documentation: Bug fix log

---

### Step 26: Windows Installer Creation

**File:** `phase5-step26-windows-installer.md`
**Duration:** 2-3 days
**Dependencies:** Step 25

Create a professional Windows installer for easy deployment. Choose between WiX Toolset (for .msi installer) or NSIS (for .exe installer). Preference: WiX v4 for professional .msi creation. Create the installer configuration file (`.wxs` for WiX or `.nsi` for NSIS) that includes: the app.exe binary, installation directory (Program Files or user's choice), Start Menu shortcut, Desktop shortcut (optional), uninstaller registration. Add an application icon (create or find a suitable icon for trend-following/trading). Set proper permissions and file associations. Test the installer: install on clean Windows VM, verify shortcuts work, run the app, verify database is created in correct location (AppData), uninstall and verify clean removal. Sign the installer (optional, requires code signing certificate). Create installation instructions document.

**Deliverables:**
- Windows installer (.msi or .exe)
- Application icon
- Start Menu and Desktop shortcuts
- Uninstaller
- Installation tested on clean Windows
- Documentation: Installation guide

---

### Step 27: User Documentation

**File:** `phase5-step27-user-documentation.md`
**Duration:** 2 days
**Dependencies:** Step 26

Create complete user-facing documentation. Write a comprehensive User Guide covering: installation, first-time setup (initialize database, configure settings), daily workflow (scan → TradingView → checklist → sizing → heat → gates → decision), understanding the banner states, what each gate means, how to interpret heat warnings, TradingView integration setup, theme customization, troubleshooting common issues. Include screenshots for each major screen. Create a Quick Start guide (one-page) for users who want to get started immediately. Write FAQ covering: what is trend-following, why the 2-minute timer, can I bypass gates (no), how to reset after a mistake, where is the database stored, how to backup data, how to export decisions to CSV (if implemented). Create tooltips and in-app help text. Write TROUBLESHOOTING.md with common error messages and solutions.

**Deliverables:**
- User Guide (comprehensive, with screenshots)
- Quick Start guide (one-page)
- FAQ document
- TROUBLESHOOTING.md
- In-app help text/tooltips
- Documentation: User documentation complete

---

### Step 28: Final Validation & Release

**File:** `phase5-step28-final-validation.md`
**Duration:** 1-2 days
**Dependencies:** Steps 24-27

Perform final validation before declaring the project production-ready. Execute the complete workflow from scratch on a clean Windows machine: download installer, install app, run first time, initialize database, configure settings, run FINVIZ scan, import candidates, open TradingView, complete checklist, calculate sizing, check heat, run gates, save GO decision. Verify every step works as expected. Test the NO-GO path as well. Review all documentation for accuracy and completeness. Ensure CLAUDE.md is up to date for future AI sessions. Update PROGRESS.md and LLM-update.md with final status. Create RELEASE_NOTES.md describing features, known limitations, system requirements. Build the final release package: installer + documentation + README. Create a checksum (SHA256) of the installer. Tag the release version (v1.0.0). Celebrate! 🎉 Create a post-release plan for bug reports and future enhancements.

**Deliverables:**
- Final validation complete on clean Windows
- All documentation reviewed and updated
- RELEASE_NOTES.md created
- Release package built (installer + docs)
- SHA256 checksum
- Version tagged (v1.0.0)
- Post-release plan
- Documentation: Final validation report

---

## Summary

**Total Duration:** 12 weeks (3 months)
**Total Steps:** 28 detailed implementation documents
**Phase Breakdown:**
- Phase 0: 4 steps (Foundation & POC)
- Phase 1: 5 steps (Dashboard & FINVIZ)
- Phase 2: 5 steps (Checklist & Sizing)
- Phase 3: 5 steps (Heat & Gates)
- Phase 4: 4 steps (Calendar & Polish)
- Phase 5: 5 steps (Testing & Packaging)

**Approach:**
Each step document (`phaseX-stepY-description.md`) will contain:
- Clear objectives and deliverables
- Step-by-step instructions
- Code examples and snippets
- Testing criteria
- Documentation requirements
- Dependencies and prerequisites
- Estimated time to complete
- Links to relevant overview-plan.md sections

**Next Actions:**
1. Review and approve this roadmap
2. Begin creating Phase 0 step documents (steps 1-4)
3. Execute each step sequentially
4. Update PROGRESS.md and LLM-update.md after each step
5. Export to Windows for testing at end of each phase
6. Celebrate milestones along the way! 🚀

---

**End of Roadmap**

**Status:** 📋 Ready for Review
**Next:** Approve roadmap → Create Phase 0 step documents → Begin implementation
**Updated:** 2025-10-29
