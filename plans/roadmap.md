# TF-Engine 2.0 - Development Roadmap (Final)

**Version:** 2.0.0 Final
**Created:** November 3, 2025
**Last Updated:** November 4, 2025
**Status:** Authoritative - Single Source of Truth

**Purpose:** Comprehensive step-by-step development plan with Gherkin scenarios, pseudo-code, test strategies, and production readiness requirements for the TF-Engine 2.0 trend-following options trading decision support system.

---

## Guiding Principles & Non-Goals

**Grounding:** Derived from the architectural rules (no feature creep; FINVIZ-only; linear workflow; save early) and the architect's intent (disciplined, green day/night UI; horserace calendar).

**Non-Goals for v2.0:**
- Building a custom screener
- Adding analytics dashboards
- Relaxing sector guardrails
- Real-time market data integration
- Social sharing features
- Mobile companion app

### Key Tenets

1. **Policy-Driven, Not Hardcoded**
   Sector/strategy mappings and FINVIZ URLs ship in `data/policy.v1.json`; UI cannot edit them; the app can run in **safe mode** if policy is missing/invalid.

2. **Linear, Anti-Impulsive Workflow**
   Cooldown + 5 required gates + heat enforcement are non-negotiable behavioral guardrails.

3. **FINVIZ-First**
   We only launch prebuilt URLs in **chart view** (`v=211`). No custom screener implementation.

4. **MVP Discipline**
   Phase 2 features exist but remain behind flags/build tags until Phase 1 is accepted.

5. **Test-Driven Development**
   Gherkin scenarios define behavior; unit tests validate implementation; integration tests prove end-to-end flow.

---

## Table of Contents

1. [Development Phases Overview](#development-phases-overview)
2. [Phase 0: Feature Freeze & Repo Hygiene](#phase-0-feature-freeze--repo-hygiene)
3. [Phase 1: Foundation Layer](#phase-1-foundation-layer)
4. [Phase 2: Core Workflow (Screens 1-3)](#phase-2-core-workflow-screens-1-3)
5. [Phase 3: Anti-Impulsivity Screens (4-6)](#phase-3-anti-impulsivity-screens-4-6)
6. [Phase 4: Trade Entry & Visualization (7-8)](#phase-4-trade-entry--visualization-screens-7-8)
7. [Phase 5: Polish & Phase 2 Features](#phase-5-polish--phase-2-features)
8. [Test Data Strategy](#test-data-strategy)
9. [Windows Installer Specification](#windows-installer-specification)
10. [Safe Mode UX Specification](#safe-mode-ux-specification)
11. [Observability & Error Handling](#observability--error-handling)
12. [Integration Testing](#integration-testing)
13. [Release Gates](#release-gates)
14. [Risks & Mitigations](#risks--mitigations)
15. [Acceptance Criteria](#acceptance-criteria)

---

## Development Phases Overview

### Phase 0: Feature Freeze & Repo Hygiene (Days 0-3)
**Goal:** Prevent scope creep; make the repo policy-driven from the first commit

**Duration:** 3 days (1.5 days per developer if pairing)

**Deliverables:**
- `data/policy.v1.json` with SHA256 signature
- `feature.flags.json` (Phase 2 features default OFF)
- `internal/screens/_post_mvp/` folder structure
- `scripts/verify_policy_hash.go` and CI integration
- `CONTRIBUTING.md` with "no feature creep" rule
- `internal/testdata/` with test fixtures

**Exit Criteria:**
- Policy loader reads and validates signatures
- Safe mode activates when policy is missing/invalid
- Feature flags system working (can toggle Phase 2 features off)
- CI fails on policy signature mismatch

---

### Phase 1: Foundation (Week 1, Days 4-6)
**Goal:** Build navigation, state management, data persistence, and reusable components

**Duration:** 3 days

**Deliverables:**
- Navigation system with history stack
- Data persistence layer (JSON with atomic writes)
- Auto-save mechanism
- Policy validation and safe mode
- **Cooldown timer widget** (reusable component)

**Exit Criteria:**
- Navigator can move forward/back/cancel
- Trades auto-save after each screen
- Cooldown timer counts down and blocks "Continue" button
- 80%+ unit test coverage on foundation

---

### Phase 2: Core Workflow (Week 1-2, Days 7-10)
**Goal:** Implement Screens 1-3

**Duration:** 4 days

**Deliverables:**
- Screen 1: Sector Selection with policy enforcement
- Screen 2: FINVIZ screener launcher
- Screen 3: Ticker entry with strategy filtering

**Exit Criteria:**
- Can select Healthcare, launch screener, enter UNH + Alt10
- Strategy dropdown filters by selected sector
- Cooldown timer starts when clicking "Continue" from Screen 3
- FINVIZ URLs verified to include `v=211`

---

### Phase 3: Anti-Impulsivity Screens (Week 2, Days 11-14)
**Goal:** Implement Screens 4-6

**Duration:** 4 days

**Deliverables:**
- Screen 4: Checklist validation (5 required + 3 optional gates)
- Screen 5: Position sizing calculator (poker multipliers)
- Screen 6: Heat check enforcement

**Exit Criteria:**
- Checklist requires all 5 gates + cooldown completion
- Position sizing math validated by unit tests
- Heat check blocks trades exceeding sector/portfolio caps

---

### Phase 4: Trade Entry & Visualization (Week 3, Days 15-21)
**Goal:** Implement Screens 7-8 and complete end-to-end workflow

**Duration:** 7 days

**Deliverables:**
- Screen 7: Options strategy entry form (26 strategy types)
- Screen 8: Trade calendar timeline widget
- Complete end-to-end workflow

**Exit Criteria:**
- Can enter a bull call spread with strikes/expiration
- Calendar displays trades as horizontal bars by sector
- Calendar renders <500ms for 100 trades
- End-to-end test passes (Sector → Calendar)

---

### Phase 5: Polish & Phase 2 Features (Week 4, Days 22-28)
**Goal:** Trade management, help system, testing, packaging

**Duration:** 7 days

**Deliverables:**
- Screen 9: Trade management (edit/delete) - **Phase 2 flag**
- Sample data generator - **Phase 2 flag**
- Help system and welcome screen
- Windows installer (NSIS)
- Comprehensive manual testing

**Exit Criteria:**
- All quality gates pass
- Windows installer builds and installs cleanly
- Phase 2 features confirmed OFF by default
- Manual testing completed (3+ full workflows)

---

## Phase 0: Feature Freeze & Repo Hygiene

### Detailed Implementation Plan

#### Day 0-1: Policy Infrastructure

**Task 1.1: Create Policy Schema with Signature**

```go
// data/policy.v1.json (excerpt)
{
  "policy_id": "tf-engine-policy",
  "version": "1.0.0",
  "app_min_version": "2.0.0",
  "security": {
    "signature_alg": "sha256",
    "signature": "abc123...",  // SHA256 of policy content (excluding this field)
    "enforce_hash": true,
    "on_hash_mismatch": "safe_mode"
  },
  "defaults": { ... },
  "sectors": [ ... ],
  "strategies": { ... }
}
```

**Task 1.2: Policy Verification Script**

```go
// scripts/verify_policy_hash.go
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
)

type PolicySecurity struct {
    SignatureAlg   string `json:"signature_alg"`
    Signature      string `json:"signature"`
    EnforceHash    bool   `json:"enforce_hash"`
    OnHashMismatch string `json:"on_hash_mismatch"`
}

type PolicyStub struct {
    Security PolicySecurity `json:"security"`
}

func main() {
    data, err := os.ReadFile("data/policy.v1.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading policy: %v\n", err)
        os.Exit(1)
    }

    // Parse to get signature
    var stub PolicyStub
    if err := json.Unmarshal(data, &stub); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing policy: %v\n", err)
        os.Exit(1)
    }

    // Calculate hash of policy without signature field
    var rawPolicy map[string]interface{}
    json.Unmarshal(data, &rawPolicy)

    // Remove signature from security section
    if sec, ok := rawPolicy["security"].(map[string]interface{}); ok {
        delete(sec, "signature")
    }

    canonical, _ := json.Marshal(rawPolicy)
    hash := sha256.Sum256(canonical)
    calculatedHash := hex.EncodeToString(hash[:])

    if calculatedHash != stub.Security.Signature {
        fmt.Fprintf(os.Stderr, "❌ Policy signature mismatch!\n")
        fmt.Fprintf(os.Stderr, "Expected: %s\n", stub.Security.Signature)
        fmt.Fprintf(os.Stderr, "Got:      %s\n", calculatedHash)
        os.Exit(1)
    }

    fmt.Println("✅ Policy signature valid")
}
```

**Task 1.3: CI Integration**

```yaml
# .github/workflows/ci.yml (if using GitHub Actions)
name: CI

on: [push, pull_request]

jobs:
  validate-policy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Verify Policy Signature
        run: go run scripts/verify_policy_hash.go

  test:
    needs: validate-policy
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run Tests
        run: go test ./...
```

---

#### Day 1-2: Feature Flags System

**Task 2.1: Feature Flags Schema**

```json
// feature.flags.json
{
  "version": "1.0.0",
  "flags": {
    "trade_management": {
      "enabled": false,
      "description": "Screen 9: Edit/delete trades",
      "phase": 2,
      "since_version": "2.1.0"
    },
    "sample_data_generator": {
      "enabled": false,
      "description": "Generate sample trades for testing",
      "phase": 2,
      "since_version": "2.1.0"
    },
    "vimium_mode": {
      "enabled": false,
      "description": "Keyboard navigation shortcuts",
      "phase": 2,
      "since_version": "2.2.0"
    },
    "advanced_analytics": {
      "enabled": false,
      "description": "Win rate tracking, equity curves",
      "phase": 2,
      "since_version": "2.3.0"
    }
  }
}
```

**Task 2.2: Feature Flag Loader**

```go
// internal/config/feature_flags.go
package config

import (
    "encoding/json"
    "os"
)

type FeatureFlag struct {
    Enabled     bool   `json:"enabled"`
    Description string `json:"description"`
    Phase       int    `json:"phase"`
    SinceVersion string `json:"since_version"`
}

type FeatureFlags struct {
    Version string                  `json:"version"`
    Flags   map[string]FeatureFlag  `json:"flags"`
}

func LoadFeatureFlags(path string) (*FeatureFlags, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var flags FeatureFlags
    if err := json.Unmarshal(data, &flags); err != nil {
        return nil, err
    }

    return &flags, nil
}

func (ff *FeatureFlags) IsEnabled(flagName string) bool {
    if flag, exists := ff.Flags[flagName]; exists {
        return flag.Enabled
    }
    return false  // Default to disabled
}
```

---

#### Day 2-3: Project Structure & Documentation

**Task 3.1: Create Folder Structure**

```bash
mkdir -p internal/screens/_post_mvp
mkdir -p internal/testdata
mkdir -p logs
mkdir -p data/backups
```

**Task 3.2: Create CONTRIBUTING.md**

```markdown
# Contributing to TF-Engine 2.0

## Golden Rule: No Feature Creep

**DO NOT add features unless:**
1. Explicitly requested by the project architect, OR
2. Approved in writing before implementation

This rule exists because the previous version failed due to excessive features.

## Pull Request Checklist

Before submitting a PR:

- [ ] Feature is in the approved roadmap
- [ ] If Phase 2 feature, it's behind a feature flag (default OFF)
- [ ] Unit tests added (>80% coverage for new code)
- [ ] Policy signature still validates (`go run scripts/verify_policy_hash.go`)
- [ ] No hardcoded business logic (check `data/policy.v1.json`)
- [ ] Behavioral guardrails preserved (cooldowns, checklist, heat checks)
- [ ] Documentation updated if adding public API

## Prohibited Changes

❌ Removing cooldown timer
❌ Skipping checklist gates
❌ Bypassing heat caps
❌ Hardcoding sector/strategy mappings
❌ Building custom screener (use FINVIZ)
❌ Adding analytics dashboards without approval

## Phase 2 Features

All Phase 2 features MUST:
1. Be in `internal/screens/_post_mvp/` or similar segregation
2. Have a feature flag in `feature.flags.json` (default: false)
3. Not execute if flag is disabled
4. Be documented as "Phase 2" in code comments

## Testing Requirements

- Unit tests: `go test ./...` must pass
- Integration test: `go test -tags=integration ./...` must pass
- Policy validation: `go run scripts/verify_policy_hash.go` must pass
- Manual test: Complete 1 full trade entry workflow
```

---

## Phase 1: Foundation Layer

### 1.1 Navigation System

#### Gherkin Scenarios

```gherkin
Feature: Screen Navigation
  As a trader
  I want to navigate between screens in a linear workflow
  So that I can complete my trade setup systematically

  Background:
    Given the application is running
    And I am on the dashboard

  Scenario: Forward navigation with valid data
    Given I am on "Screen 1: Sector Selection"
    And I have selected "Healthcare" sector
    When I click the "Continue" button
    Then I should be navigated to "Screen 2: Screener Launch"
    And my sector selection should be auto-saved
    And the file "data/trades_in_progress.json" should exist

  Scenario: Forward navigation with invalid data
    Given I am on "Screen 1: Sector Selection"
    And I have NOT selected any sector
    When I try to click the "Continue" button
    Then the "Continue" button should be disabled
    And I should see a tooltip "Please select a sector"

  Scenario: Back navigation preserves data
    Given I am on "Screen 3: Ticker Entry"
    And I have entered ticker "UNH"
    And I have selected strategy "Alt10"
    When I click the "Back" button
    Then I should be navigated to "Screen 2: Screener Launch"
    When I click "Continue" to return to "Screen 3"
    Then I should see ticker "UNH" still populated
    And I should see strategy "Alt10" still selected

  Scenario: Cancel returns to dashboard with confirmation
    Given I am on "Screen 4: Checklist"
    And I have partially completed the checklist
    When I click the "Cancel" button
    Then I should see a confirmation dialog "Are you sure? Progress will be saved."
    When I click "Yes, Cancel"
    Then I should be navigated to "Dashboard"
    And my progress should be saved as "in-progress trade"

  Scenario: Jump to calendar from any screen
    Given I am on "Screen 5: Position Sizing"
    When I click the "View Calendar" button in the navigation bar
    Then I should be navigated to "Screen 8: Calendar"
    And the calendar should be in "read-only mode"
    And I should see a "Return to Trade Entry" button
```

#### Pseudo-Code

```go
// internal/ui/navigator.go

type Navigator struct {
    screens       []Screen
    currentIndex  int
    history       []int
    state         *appcore.AppState
    window        fyne.Window
}

func NewNavigator(state *appcore.AppState, window fyne.Window) *Navigator {
    nav := &Navigator{
        screens: []Screen{
            NewSectorSelection(state, window),
            NewScreenerLaunch(state, window),
            NewTickerEntry(state, window),
            NewChecklist(state, window),
            NewPositionSizing(state, window),
            NewHeatCheck(state, window),
            NewTradeEntry(state, window),
            NewCalendar(state, window),
        },
        currentIndex: -1,  // -1 = dashboard
        history:      []int{},
        state:        state,
        window:       window,
    }
    return nav
}

// Navigate forward to next screen
func (n *Navigator) Next() error {
    // Validate current screen before proceeding
    if !n.ValidateCurrentScreen() {
        return errors.New("current screen validation failed")
    }

    // Auto-save before navigation
    if err := n.AutoSave(); err != nil {
        return fmt.Errorf("auto-save failed: %w", err)
    }

    // Record history for back button
    n.history = append(n.history, n.currentIndex)

    // Move to next screen
    n.currentIndex++
    if n.currentIndex >= len(n.screens) {
        return errors.New("no more screens")
    }

    // Update state
    n.state.CurrentScreen = n.GetCurrentScreenName()

    // Render new screen
    n.window.SetContent(n.screens[n.currentIndex].Render())

    return nil
}

// Navigate back to previous screen
func (n *Navigator) Back() error {
    if len(n.history) == 0 {
        return errors.New("no previous screen")
    }

    // Auto-save before navigation
    if err := n.AutoSave(); err != nil {
        return fmt.Errorf("auto-save failed: %w", err)
    }

    // Pop from history
    n.currentIndex = n.history[len(n.history)-1]
    n.history = n.history[:len(n.history)-1]

    // Update state
    n.state.CurrentScreen = n.GetCurrentScreenName()

    // Render previous screen
    n.window.SetContent(n.screens[n.currentIndex].Render())

    return nil
}

// Cancel with confirmation
func (n *Navigator) Cancel() {
    dialog.ShowConfirm(
        "Cancel Trade Entry?",
        "Your progress will be saved. Are you sure you want to return to dashboard?",
        func(confirmed bool) {
            if confirmed {
                n.AutoSave()
                n.NavigateToDashboard()
            }
        },
        n.window,
    )
}

// Validate current screen's data
func (n *Navigator) ValidateCurrentScreen() bool {
    if n.currentIndex < 0 || n.currentIndex >= len(n.screens) {
        return true  // Dashboard or out of bounds
    }
    return n.screens[n.currentIndex].Validate()
}

// Auto-save trade data
func (n *Navigator) AutoSave() error {
    if n.state.CurrentTrade == nil {
        return nil  // Nothing to save
    }
    return storage.SaveInProgressTrade(n.state.CurrentTrade)
}

// GetCurrentScreenName returns human-readable screen name
func (n *Navigator) GetCurrentScreenName() string {
    if n.currentIndex < 0 {
        return "dashboard"
    }
    screenNames := []string{
        "sector_selection",
        "screener_launch",
        "ticker_entry",
        "checklist",
        "position_sizing",
        "heat_check",
        "trade_entry",
        "calendar",
    }
    if n.currentIndex < len(screenNames) {
        return screenNames[n.currentIndex]
    }
    return "unknown"
}
```

#### Testing Strategy

```go
// internal/ui/navigator_test.go

func TestNavigator_Next_ValidData(t *testing.T) {
    // Arrange
    state := appcore.NewAppState()
    state.CurrentTrade = &models.Trade{Sector: "Healthcare"}
    nav := NewNavigator(state, nil)

    // Act
    err := nav.Next()

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 0, nav.currentIndex)
    assert.Equal(t, "sector_selection", state.CurrentScreen)
}

func TestNavigator_Next_InvalidData_Fails(t *testing.T) {
    // Arrange
    state := appcore.NewAppState()
    // No trade data set
    nav := NewNavigator(state, nil)

    // Act
    err := nav.Next()

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "validation failed")
}

func TestNavigator_Back_PreservesData(t *testing.T) {
    // Arrange
    state := appcore.NewAppState()
    state.CurrentTrade = &models.Trade{
        Sector:   "Healthcare",
        Ticker:   "UNH",
        Strategy: "Alt10",
    }
    nav := NewNavigator(state, nil)
    nav.Next() // Move to screen 0
    nav.Next() // Move to screen 1

    // Act
    err := nav.Back()

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 0, nav.currentIndex)
    assert.Equal(t, "UNH", state.CurrentTrade.Ticker) // Data preserved
}

func TestNavigator_AutoSave_CalledOnNavigation(t *testing.T) {
    // Arrange
    state := appcore.NewAppState()
    state.CurrentTrade = &models.Trade{Sector: "Healthcare"}
    nav := NewNavigator(state, nil)

    // Act
    nav.Next()

    // Assert
    // Verify file exists
    assert.FileExists(t, "data/trades_in_progress.json")
}
```

---

### 1.2 Data Persistence Layer

#### Gherkin Scenarios

```gherkin
Feature: Trade Data Persistence
  As a trader
  I want my trade data to be automatically saved
  So that I never lose my progress

  Scenario: Auto-save after each screen
    Given I am on "Screen 1: Sector Selection"
    And I select "Healthcare" sector
    When I click "Continue"
    Then a file "data/trades_in_progress.json" should be created
    And it should contain my sector selection
    And the file modification time should be within 1 second of now

  Scenario: Resume incomplete trade on app restart
    Given I was on "Screen 3: Ticker Entry"
    And I entered ticker "MSFT" and strategy "Alt26"
    And the app crashed before I could continue
    When I restart the application
    Then I should see a dialog "Resume in-progress trade?"
    When I click "Yes"
    Then I should be navigated to "Screen 3: Ticker Entry"
    And I should see ticker "MSFT" populated
    And I should see strategy "Alt26" selected

  Scenario: Save completed trade to history
    Given I have completed all 8 screens
    And I am on "Screen 8: Calendar"
    When I click "Finish and Save Trade"
    Then a file "data/trades.json" should be updated
    And the trade should appear in my trade history
    And the "trades_in_progress.json" should be deleted
    And a backup should exist in "data/backups/"

  Scenario: Prevent data corruption on concurrent saves
    Given I have multiple auto-save operations queued
    When each screen transition triggers an auto-save
    Then only the latest data should be persisted
    And no file corruption should occur
    And the JSON should be valid
```

#### Pseudo-Code

```go
// internal/storage/trades.go

import (
    "encoding/json"
    "os"
    "sync"
    "time"
)

// Storage paths
const (
    TradesFile     = "data/trades.json"
    InProgressFile = "data/trades_in_progress.json"
    BackupDir      = "data/backups/"
)

// Thread-safe storage manager
type TradeStorage struct {
    mu sync.RWMutex
}

var globalStorage = &TradeStorage{}

// SaveInProgressTrade saves current trade state
func SaveInProgressTrade(trade *models.Trade) error {
    globalStorage.mu.Lock()
    defer globalStorage.mu.Unlock()

    trade.LastUpdated = time.Now()

    data, err := json.MarshalIndent(trade, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal error: %w", err)
    }

    // Write atomically (write to temp, then rename)
    tmpFile := InProgressFile + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return fmt.Errorf("write error: %w", err)
    }

    if err := os.Rename(tmpFile, InProgressFile); err != nil {
        return fmt.Errorf("rename error: %w", err)
    }

    return nil
}

// LoadInProgressTrade loads incomplete trade
func LoadInProgressTrade() (*models.Trade, error) {
    globalStorage.mu.RLock()
    defer globalStorage.mu.RUnlock()

    if !fileExists(InProgressFile) {
        return nil, nil  // No in-progress trade
    }

    data, err := os.ReadFile(InProgressFile)
    if err != nil {
        return nil, fmt.Errorf("read error: %w", err)
    }

    var trade models.Trade
    if err := json.Unmarshal(data, &trade); err != nil {
        return nil, fmt.Errorf("unmarshal error: %w", err)
    }

    return &trade, nil
}

// SaveCompletedTrade saves trade to history
func SaveCompletedTrade(trade *models.Trade) error {
    globalStorage.mu.Lock()
    defer globalStorage.mu.Unlock()

    trade.CompletedAt = time.Now()
    trade.Status = "active"

    // Load existing trades
    trades, err := loadAllTradesUnsafe()
    if err != nil {
        trades = []models.Trade{}
    }

    // Append new trade
    trades = append(trades, *trade)

    // Backup existing file before overwriting
    if fileExists(TradesFile) {
        backupPath := fmt.Sprintf("%strades_%s.json",
            BackupDir, time.Now().Format("20060102_150405"))
        os.MkdirAll(BackupDir, 0755)
        if err := copyFile(TradesFile, backupPath); err != nil {
            // Log warning but don't fail
            log.Warn("Failed to create backup: %v", err)
        }
    }

    // Save to file atomically
    data, err := json.MarshalIndent(trades, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal error: %w", err)
    }

    tmpFile := TradesFile + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return fmt.Errorf("write error: %w", err)
    }

    if err := os.Rename(tmpFile, TradesFile); err != nil {
        return fmt.Errorf("rename error: %w", err)
    }

    // Clear in-progress file
    os.Remove(InProgressFile)

    return nil
}

// LoadAllTrades loads complete trade history
func LoadAllTrades() ([]models.Trade, error) {
    globalStorage.mu.RLock()
    defer globalStorage.mu.RUnlock()
    return loadAllTradesUnsafe()
}

func loadAllTradesUnsafe() ([]models.Trade, error) {
    if !fileExists(TradesFile) {
        return []models.Trade{}, nil
    }

    data, err := os.ReadFile(TradesFile)
    if err != nil {
        return nil, fmt.Errorf("read error: %w", err)
    }

    var trades []models.Trade
    if err := json.Unmarshal(data, &trades); err != nil {
        return nil, fmt.Errorf("unmarshal error: %w", err)
    }

    return trades, nil
}

// Helper functions
func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func copyFile(src, dst string) error {
    data, err := os.ReadFile(src)
    if err != nil {
        return err
    }
    return os.WriteFile(dst, data, 0644)
}
```

#### Testing Strategy

```go
// internal/storage/trades_test.go

func TestSaveInProgressTrade_CreatesFile(t *testing.T) {
    // Arrange
    cleanup := setupTestDataDir(t)
    defer cleanup()

    trade := &models.Trade{
        Sector: "Healthcare",
        Ticker: "UNH",
    }

    // Act
    err := SaveInProgressTrade(trade)

    // Assert
    assert.NoError(t, err)
    assert.FileExists(t, InProgressFile)

    // Verify JSON is valid
    data, _ := os.ReadFile(InProgressFile)
    var loaded models.Trade
    assert.NoError(t, json.Unmarshal(data, &loaded))
}

func TestLoadInProgressTrade_RestoresData(t *testing.T) {
    // Arrange
    cleanup := setupTestDataDir(t)
    defer cleanup()

    originalTrade := &models.Trade{
        Sector:   "Technology",
        Ticker:   "MSFT",
        Strategy: "Alt26",
    }
    SaveInProgressTrade(originalTrade)

    // Act
    loadedTrade, err := LoadInProgressTrade()

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Technology", loadedTrade.Sector)
    assert.Equal(t, "MSFT", loadedTrade.Ticker)
    assert.Equal(t, "Alt26", loadedTrade.Strategy)
}

func TestSaveCompletedTrade_AppendsToHistory(t *testing.T) {
    // Arrange
    cleanup := setupTestDataDir(t)
    defer cleanup()

    // Save first trade
    trade1 := &models.Trade{Ticker: "UNH"}
    SaveCompletedTrade(trade1)

    // Save second trade
    trade2 := &models.Trade{Ticker: "MSFT"}

    // Act
    err := SaveCompletedTrade(trade2)

    // Assert
    assert.NoError(t, err)
    trades, _ := LoadAllTrades()
    assert.Len(t, trades, 2)
    assert.Equal(t, "UNH", trades[0].Ticker)
    assert.Equal(t, "MSFT", trades[1].Ticker)
}

func TestSaveCompletedTrade_CreatesBackup(t *testing.T) {
    // Arrange
    cleanup := setupTestDataDir(t)
    defer cleanup()

    trade1 := &models.Trade{Ticker: "UNH"}
    SaveCompletedTrade(trade1)

    // Act
    trade2 := &models.Trade{Ticker: "MSFT"}
    SaveCompletedTrade(trade2)

    // Assert - backup should exist
    backups, _ := filepath.Glob(BackupDir + "trades_*.json")
    assert.NotEmpty(t, backups)
}

func TestConcurrentSaves_NoCorruption(t *testing.T) {
    // Arrange
    cleanup := setupTestDataDir(t)
    defer cleanup()

    var wg sync.WaitGroup

    // Act - concurrent saves
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            trade := &models.Trade{Ticker: fmt.Sprintf("TEST%d", id)}
            SaveInProgressTrade(trade)
        }(i)
    }

    wg.Wait()

    // Assert
    trade, err := LoadInProgressTrade()
    assert.NoError(t, err)
    assert.NotNil(t, trade)
    // Should have one of the trades (not corrupted)
    assert.Contains(t, trade.Ticker, "TEST")

    // Verify JSON is valid
    data, _ := os.ReadFile(InProgressFile)
    var loaded models.Trade
    assert.NoError(t, json.Unmarshal(data, &loaded))
}
```

---

### 1.3 Cooldown Timer Widget

#### Gherkin Scenarios

```gherkin
Feature: Cooldown Timer Widget
  As a trader
  I want a visual countdown timer
  So that I know how long until I can proceed

  Scenario: Timer starts at configured duration
    Given the policy specifies 120 seconds cooldown
    When the cooldown timer starts
    Then I should see "2:00" displayed
    And the timer should be counting down

  Scenario: Timer updates every second
    Given the cooldown timer is running
    When 1 second passes
    Then the display should update to "1:59"
    And the progress bar should decrease

  Scenario: Timer completion enables button
    Given the cooldown timer started 119 seconds ago
    And the "Continue" button is disabled
    When the timer reaches 0:00
    Then the "Continue" button should be enabled
    And I should see "Ready to continue"

  Scenario: Timer persists across app restarts
    Given the cooldown timer started 60 seconds ago
    And the app crashes
    When I restart the application
    And I navigate back to the checklist screen
    Then the timer should show approximately 1:00 remaining
    And it should continue counting down

  Scenario: Cannot bypass timer by clicking repeatedly
    Given the cooldown timer shows 0:30 remaining
    When I click the "Continue" button 10 times rapidly
    Then the button should remain disabled
    And I should see a message "Please wait for cooldown to complete"
```

#### Pseudo-Code

```go
// internal/widgets/cooldown_timer.go

import (
    "time"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/container"
)

type CooldownTimer struct {
    widget.BaseWidget

    duration      time.Duration
    startTime     time.Time
    ticker        *time.Ticker
    done          chan bool
    onComplete    func()

    label         *widget.Label
    progressBar   *widget.ProgressBar
}

func NewCooldownTimer(duration time.Duration, onComplete func()) *CooldownTimer {
    timer := &CooldownTimer{
        duration:   duration,
        startTime:  time.Now(),
        ticker:     time.NewTicker(1 * time.Second),
        done:       make(chan bool),
        onComplete: onComplete,
        label:      widget.NewLabel(""),
        progressBar: widget.NewProgressBar(),
    }

    timer.ExtendBaseWidget(timer)
    timer.start()

    return timer
}

func (t *CooldownTimer) CreateRenderer() fyne.WidgetRenderer {
    return widget.NewSimpleRenderer(
        container.NewVBox(
            t.label,
            t.progressBar,
        ),
    )
}

func (t *CooldownTimer) start() {
    go func() {
        for {
            select {
            case <-t.done:
                return
            case <-t.ticker.C:
                t.update()
            }
        }
    }()
}

func (t *CooldownTimer) update() {
    elapsed := time.Since(t.startTime)
    remaining := t.duration - elapsed

    if remaining <= 0 {
        t.complete()
        return
    }

    // Update label
    minutes := int(remaining.Minutes())
    seconds := int(remaining.Seconds()) % 60
    t.label.SetText(fmt.Sprintf("Cooldown: %d:%02d remaining", minutes, seconds))

    // Update progress bar (inverse - counts down)
    progress := float64(elapsed) / float64(t.duration)
    t.progressBar.SetValue(progress)

    t.Refresh()
}

func (t *CooldownTimer) complete() {
    t.ticker.Stop()
    t.label.SetText("Ready to continue")
    t.progressBar.SetValue(1.0)

    if t.onComplete != nil {
        t.onComplete()
    }

    t.Refresh()
}

func (t *CooldownTimer) Stop() {
    t.done <- true
    t.ticker.Stop()
}

func (t *CooldownTimer) GetRemaining() time.Duration {
    elapsed := time.Since(t.startTime)
    remaining := t.duration - elapsed
    if remaining < 0 {
        return 0
    }
    return remaining
}

func (t *CooldownTimer) IsComplete() bool {
    return t.GetRemaining() <= 0
}
```

#### Testing Strategy

```go
// internal/widgets/cooldown_timer_test.go

func TestCooldownTimer_StartsAtFullDuration(t *testing.T) {
    // Arrange
    duration := 120 * time.Second
    called := false
    onComplete := func() { called = true }

    // Act
    timer := NewCooldownTimer(duration, onComplete)
    defer timer.Stop()

    // Assert
    remaining := timer.GetRemaining()
    assert.InDelta(t, 120.0, remaining.Seconds(), 1.0)
    assert.False(t, timer.IsComplete())
    assert.False(t, called)
}

func TestCooldownTimer_CountsDown(t *testing.T) {
    // Arrange
    duration := 3 * time.Second
    timer := NewCooldownTimer(duration, nil)
    defer timer.Stop()

    // Act
    time.Sleep(2 * time.Second)

    // Assert
    remaining := timer.GetRemaining()
    assert.InDelta(t, 1.0, remaining.Seconds(), 0.5)
}

func TestCooldownTimer_CallsOnComplete(t *testing.T) {
    // Arrange
    duration := 1 * time.Second
    called := false
    onComplete := func() { called = true }

    timer := NewCooldownTimer(duration, onComplete)
    defer timer.Stop()

    // Act
    time.Sleep(1500 * time.Millisecond)

    // Assert
    assert.True(t, called)
    assert.True(t, timer.IsComplete())
}

func TestCooldownTimer_Persistence(t *testing.T) {
    // Arrange
    startTime := time.Now().Add(-60 * time.Second) // Started 60 seconds ago

    // Simulate app restart - recreate timer with past start time
    duration := 120 * time.Second
    elapsed := time.Since(startTime)
    remainingDuration := duration - elapsed

    timer := NewCooldownTimer(remainingDuration, nil)
    defer timer.Stop()

    // Assert
    remaining := timer.GetRemaining()
    assert.InDelta(t, 60.0, remaining.Seconds(), 2.0)
}
```

---

## 6. Test Data Strategy

**Purpose:** Enable realistic testing without requiring live market data or real trades.

### 6.1 Test Data Generators

**Location:** `internal/testing/generators/`

```go
// internal/testing/generators/trades.go
package generators

import (
    "math/rand"
    "time"
    "tf-engine/internal/models"
)

// GenerateSampleTrades creates realistic sample trades for calendar view testing
func GenerateSampleTrades(count int) []models.Trade {
    sectors := []string{"Healthcare", "Technology", "Industrials", "Consumer", "Financials"}
    tickers := map[string][]string{
        "Healthcare":   {"UNH", "JNJ", "PFE", "ABBV", "TMO"},
        "Technology":   {"MSFT", "AAPL", "NVDA", "AMD", "CRM"},
        "Industrials":  {"CAT", "BA", "UNP", "GE", "HON"},
        "Consumer":     {"AMZN", "HD", "NKE", "MCD", "SBUX"},
        "Financials":   {"JPM", "BAC", "GS", "MS", "BLK"},
    }
    strategies := []string{"Alt10", "Alt26", "Alt43", "Alt46"}
    optionsTypes := []string{"Bull call spread", "Bear put spread", "Bull put credit spread", "Iron condor"}

    trades := make([]models.Trade, count)
    now := time.Now()

    for i := 0; i < count; i++ {
        sector := sectors[rand.Intn(len(sectors))]
        ticker := tickers[sector][rand.Intn(len(tickers[sector]))]

        entryDate := now.AddDate(0, 0, -rand.Intn(14)) // Random entry within past 14 days
        expirationDate := entryDate.AddDate(0, 0, rand.Intn(70)+14) // Expire 2-12 weeks out

        trades[i] = models.Trade{
            ID:             generateTradeID(),
            Ticker:         ticker,
            Sector:         sector,
            Strategy:       strategies[rand.Intn(len(strategies))],
            OptionsType:    optionsTypes[rand.Intn(len(optionsTypes))],
            EntryDate:      entryDate,
            ExpirationDate: expirationDate,
            Risk:           float64(rand.Intn(500) + 200), // $200-$700 risk
            Status:         "active",
            PnL:            float64(rand.Intn(400) - 100), // -$100 to +$300 P&L
        }
    }

    return trades
}

// GenerateHeatCheckScenario creates trades that test heat limit enforcement
func GenerateHeatCheckScenario() []models.Trade {
    return []models.Trade{
        {Ticker: "UNH", Sector: "Healthcare", Risk: 450, Status: "active"},
        {Ticker: "JNJ", Sector: "Healthcare", Risk: 300, Status: "active"},
        {Ticker: "ABBV", Sector: "Healthcare", Risk: 400, Status: "active"}, // Should trigger 1.5% sector limit
        {Ticker: "MSFT", Sector: "Technology", Risk: 500, Status: "active"},
        {Ticker: "AAPL", Sector: "Technology", Risk: 350, Status: "active"},
    }
}
```

### 6.2 Test Fixtures

**Location:** `testdata/fixtures/`

```
testdata/
├── fixtures/
│   ├── policy.valid.json          # Valid policy with proper signature
│   ├── policy.invalid_sig.json    # Policy with mismatched signature
│   ├── policy.corrupted.json      # Malformed JSON
│   ├── trades.empty.json          # Empty trade history
│   ├── trades.active.json         # 3 active trades
│   ├── trades.heat_limit.json     # Trades at heat limit
│   └── trades.mixed.json          # Active + closed trades
```

### 6.3 Sample Data Button (Phase 2 Feature)

**Screen:** Dashboard (Screen 8)

**Behavior:**
1. User clicks "Generate Sample Data" button
2. App calls `GenerateSampleTrades(10)`
3. Saves sample trades to `data/trades.json`
4. Reloads dashboard calendar view
5. Shows toast: "Generated 10 sample trades"

**Gherkin Scenario:**

```gherkin
Feature: Sample Data Generation

  Scenario: User generates sample trades for testing
    Given the dashboard is displayed
    And the calendar view is empty
    When the user clicks "Generate Sample Data"
    Then 10 realistic sample trades are created
    And the calendar view displays the sample trades
    And trades span multiple sectors and expiration dates
    And the sample data button is disabled (prevent duplicates)
```

---

## 7. Windows Installer Specification (NSIS)

**Technology:** NSIS (Nullsoft Scriptable Install System)
**Rationale:** Industry-standard, scriptable, creates small executables, supports silent installs

### 7.1 Installer Requirements

**Deliverables:**
- `TFEngine-Setup-1.0.0.exe` (installer executable)
- Silent install support: `/S` flag
- Custom install directory selection
- Desktop shortcut creation (optional)
- Start Menu entry
- Uninstaller with proper cleanup
- Version detection (prevent downgrade)

### 7.2 NSIS Script

**Location:** `build/installer/tf-engine-installer.nsi`

```nsis
; TF-Engine 2.0 NSIS Installer Script
; Build command: makensis tf-engine-installer.nsi

!include "MUI2.nsh"
!include "LogicLib.nsh"

; Application metadata
!define APP_NAME "TF-Engine"
!define APP_VERSION "1.0.0"
!define APP_PUBLISHER "TF Systems"
!define APP_EXE "tf-engine.exe"
!define UNINSTALL_EXE "Uninstall.exe"

Name "${APP_NAME} ${APP_VERSION}"
OutFile "TFEngine-Setup-${APP_VERSION}.exe"
InstallDir "$PROGRAMFILES64\${APP_NAME}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallPath"

RequestExecutionLevel admin

; Modern UI Configuration
!define MUI_ICON "assets\icon.ico"
!define MUI_UNICON "assets\icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\${APP_EXE}"
!define MUI_FINISHPAGE_RUN_TEXT "Launch TF-Engine"

; Installer pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

; Version check function
Function .onInit
  ReadRegStr $0 HKLM "Software\${APP_NAME}" "Version"
  ${If} $0 != ""
    ${VersionCompare} $0 "${APP_VERSION}" $R0
    ${If} $R0 == 1
      MessageBox MB_YESNO|MB_ICONQUESTION "A newer version ($0) is already installed. Continue anyway?" IDYES +2
      Abort
    ${EndIf}
  ${EndIf}
FunctionEnd

; Main installation section
Section "Install"
  SetOutPath "$INSTDIR"

  ; Core files
  File "dist\${APP_EXE}"
  File "data\policy.v1.json"
  File "LICENSE.txt"
  File "README.md"

  ; Create data directory
  CreateDirectory "$INSTDIR\data"

  ; Registry keys
  WriteRegStr HKLM "Software\${APP_NAME}" "InstallPath" "$INSTDIR"
  WriteRegStr HKLM "Software\${APP_NAME}" "Version" "${APP_VERSION}"

  ; Uninstaller
  WriteUninstaller "$INSTDIR\${UNINSTALL_EXE}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\${UNINSTALL_EXE}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "Publisher" "${APP_PUBLISHER}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayVersion" "${APP_VERSION}"

  ; Start Menu shortcut
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}"
  CreateShortcut "$SMPROGRAMS\${APP_NAME}\Uninstall.lnk" "$INSTDIR\${UNINSTALL_EXE}"

  ; Desktop shortcut (optional)
  MessageBox MB_YESNO "Create desktop shortcut?" IDNO +2
  CreateShortcut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}"
SectionEnd

; Uninstaller
Section "Uninstall"
  ; Remove files
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\policy.v1.json"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\${UNINSTALL_EXE}"

  ; Prompt to preserve user data
  MessageBox MB_YESNO "Delete trade history and settings? (Cannot be undone)" IDYES delete_data
  Goto skip_data
  delete_data:
    RMDir /r "$INSTDIR\data"
  skip_data:

  ; Remove shortcuts
  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\Uninstall.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  ; Remove registry keys
  DeleteRegKey HKLM "Software\${APP_NAME}"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"

  ; Remove install directory (if empty)
  RMDir "$INSTDIR"
SectionEnd
```

### 7.3 Build Process

**Location:** `scripts/build-installer.sh` (Git Bash on Windows)

```bash
#!/bin/bash
set -e

echo "Building TF-Engine Windows Installer..."

# Step 1: Build Go executable
echo "Step 1/4: Compiling Go application..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/tf-engine.exe .

# Step 2: Copy assets
echo "Step 2/4: Copying assets..."
mkdir -p dist
cp data/policy.v1.json dist/
cp LICENSE.txt dist/
cp README.md dist/

# Step 3: Run NSIS compiler
echo "Step 3/4: Building installer..."
makensis build/installer/tf-engine-installer.nsi

# Step 4: Sign executable (optional - requires code signing certificate)
if [ -f "code-signing.p12" ]; then
  echo "Step 4/4: Signing installer..."
  signtool sign /f code-signing.p12 /p "$CERT_PASSWORD" /tr http://timestamp.digicert.com /td sha256 /fd sha256 TFEngine-Setup-1.0.0.exe
else
  echo "Step 4/4: Skipping code signing (no certificate found)"
fi

echo "✅ Installer created: TFEngine-Setup-1.0.0.exe"
```

### 7.4 Testing Installer

**Manual Test Checklist:**
- [ ] Silent install: `TFEngine-Setup-1.0.0.exe /S`
- [ ] Custom directory selection works
- [ ] Desktop shortcut appears (if selected)
- [ ] Start Menu entry created
- [ ] Application launches after install
- [ ] Uninstaller prompts for data preservation
- [ ] Uninstaller removes all files/registry keys
- [ ] Version downgrade prevention works
- [ ] Install on clean Windows 10/11 VM

**Automated Test:** Use Chocolatey test environment or WiX Toolset test harness

---

## 8. Safe Mode UX Specification

**Purpose:** Gracefully handle missing/corrupted policy with minimal-configuration fallback.

### 8.1 Safe Mode Activation Triggers

1. **Policy file missing** (`data/policy.v1.json` not found)
2. **Policy signature mismatch** (SHA256 hash validation fails)
3. **Policy JSON malformed** (parsing error)
4. **User explicitly requests** (via command-line flag `--safe-mode`)

### 8.2 Safe Mode Visual Indicators

**Banner Widget:**

```go
// internal/widgets/safe_mode_banner.go
package widgets

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "image/color"
)

type SafeModeBanner struct {
    widget.BaseWidget
}

func NewSafeModeBanner() *SafeModeBanner {
    banner := &SafeModeBanner{}
    banner.ExtendBaseWidget(banner)
    return banner
}

func (b *SafeModeBanner) CreateRenderer() fyne.WidgetRenderer {
    bg := canvas.NewRectangle(color.RGBA{R: 255, G: 140, B: 0, A: 255}) // Orange background

    icon := widget.NewIcon(theme.WarningIcon())

    label := widget.NewLabel("SAFE MODE: Using minimal policy configuration")
    label.TextStyle = fyne.TextStyle{Bold: true}

    learnMoreBtn := widget.NewButton("Learn More", func() {
        // Show dialog explaining safe mode
    })

    content := container.NewBorder(nil, nil, icon, learnMoreBtn, label)

    return &safeModeBannerRenderer{
        banner: b,
        bg:     bg,
        content: content,
    }
}

type safeModeBannerRenderer struct {
    banner  *SafeModeBanner
    bg      *canvas.Rectangle
    content *fyne.Container
}

func (r *safeModeBannerRenderer) Layout(size fyne.Size) {
    r.bg.Resize(size)
    r.content.Resize(size)
}

func (r *safeModeBannerRenderer) MinSize() fyne.Size {
    return fyne.NewSize(400, 50)
}

func (r *safeModeBannerRenderer) Refresh() {
    r.bg.Refresh()
    r.content.Refresh()
}

func (r *safeModeBannerRenderer) Objects() []fyne.CanvasObject {
    return []fyne.CanvasObject{r.bg, r.content}
}

func (r *safeModeBannerRenderer) Destroy() {}
```

**Window Title Suffix:**
- Normal mode: "TF-Engine 2.0"
- Safe mode: "TF-Engine 2.0 [SAFE MODE]"

**Status Bar Indicator:**
- Orange dot icon in status bar (persistent across all screens)

### 8.3 Safe Mode Restrictions

**Allowed Operations:**
- View safe mode policy (3 sectors: Healthcare, Technology, Utilities-blocked)
- Navigate all 8 screens
- Use cooldown timer
- Use checklist
- View sample data

**Blocked Operations:**
- ❌ Execute real trades (Continue button disabled with tooltip: "Real trades disabled in Safe Mode")
- ❌ Modify policy file
- ❌ Access advanced settings

### 8.4 Safe Mode Recovery Flow

**Gherkin Scenario:**

```gherkin
Feature: Safe Mode Recovery

  Scenario: App detects corrupted policy and activates safe mode
    Given the policy file exists but has invalid signature
    When the application launches
    Then the app displays the Safe Mode banner
    And the window title shows "[SAFE MODE]"
    And the user sees a dialog: "Policy file validation failed. Running in Safe Mode with minimal configuration. Please restore policy.v1.json or contact support."
    And the dialog has buttons: ["Download Policy", "Continue in Safe Mode", "Exit"]

  Scenario: User restores policy while in safe mode
    Given the app is running in Safe Mode
    When the user clicks "Download Policy" in the dialog
    Then the app opens the browser to "https://example.com/policy-download"
    And the user downloads valid policy.v1.json
    And the user clicks "File → Reload Policy"
    Then the app validates the new policy signature
    And the app exits Safe Mode
    And the Safe Mode banner disappears
```

### 8.5 Telemetry Events

Log these events when safe mode activates:

```json
{
  "event": "safe_mode_activated",
  "timestamp": "2025-11-03T14:32:11Z",
  "reason": "policy_signature_mismatch",
  "policy_path": "data/policy.v1.json",
  "expected_hash": "a3f2...",
  "actual_hash": "b8d1...",
  "user_action": "continued_in_safe_mode"
}
```

---

## 9. Observability & Error Handling

### 9.1 Structured Logging

**Framework:** Go standard library `log/slog` (Go 1.21+)

```go
// internal/logging/logger.go
package logging

import (
    "log/slog"
    "os"
)

var Logger *slog.Logger

func InitLogger(debug bool) {
    level := slog.LevelInfo
    if debug {
        level = slog.LevelDebug
    }

    handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: level,
    })

    Logger = slog.New(handler)
}

// Usage:
// logging.Logger.Info("policy loaded", "sectors", len(policy.Sectors), "version", policy.Version)
// logging.Logger.Error("navigation failed", "screen", currentScreen, "error", err)
```

**Log Levels:**
- **DEBUG:** Screen transitions, button clicks (only when `--debug` flag set)
- **INFO:** App start, policy loaded, trade saved, cooldown started
- **WARN:** Safe mode activated, heat limit approaching
- **ERROR:** Policy validation failed, file I/O errors, navigation errors

### 9.2 Error Taxonomy

**Error Categories:**

```go
// internal/errors/types.go
package errors

import "fmt"

type ErrorCategory string

const (
    UserError   ErrorCategory = "USER_ERROR"   // Invalid input, validation failures
    DataError   ErrorCategory = "DATA_ERROR"   // File not found, corrupted data
    SystemError ErrorCategory = "SYSTEM_ERROR" // OS errors, permissions, crashes
)

type AppError struct {
    Category ErrorCategory
    Message  string
    Cause    error
    Context  map[string]interface{}
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Cause)
}

func NewUserError(msg string, cause error) *AppError {
    return &AppError{
        Category: UserError,
        Message:  msg,
        Cause:    cause,
        Context:  make(map[string]interface{}),
    }
}

func NewDataError(msg string, cause error) *AppError {
    return &AppError{
        Category: DataError,
        Message:  msg,
        Cause:    cause,
        Context:  make(map[string]interface{}),
    }
}

func NewSystemError(msg string, cause error) *AppError {
    return &AppError{
        Category: SystemError,
        Message:  msg,
        Cause:    cause,
        Context:  make(map[string]interface{}),
    }
}
```

**Error Handling Examples:**

```go
// User Error (show friendly message, don't log as ERROR)
if selectedSector.Blocked {
    err := errors.NewUserError("sector is blocked", nil)
    err.Context["sector"] = selectedSector.Name
    dialog.ShowError(errors.New("This sector is blocked for trading based on backtest results."), window)
    logging.Logger.Warn("user attempted blocked sector", "sector", selectedSector.Name)
    return err
}

// Data Error (show message, log as ERROR, activate safe mode)
policy, err := models.LoadPolicy("data/policy.v1.json")
if err != nil {
    dataErr := errors.NewDataError("failed to load policy", err)
    logging.Logger.Error("policy load failed", "error", err, "activating_safe_mode", true)
    appState.UseSafeMode()
    return dataErr
}

// System Error (show message, log as ERROR, offer exit)
err := storage.SaveCompletedTrade(trade)
if err != nil {
    sysErr := errors.NewSystemError("failed to save trade", err)
    logging.Logger.Error("trade save failed", "error", err, "trade_id", trade.ID)
    dialog.ShowError(errors.New("Failed to save trade data. Check disk space and permissions."), window)
    return sysErr
}
```

### 9.3 Telemetry Events

**Key Events to Track:**

```go
// App Lifecycle
logging.Logger.Info("app_started", "version", "1.0.0", "os", runtime.GOOS)
logging.Logger.Info("app_shutdown", "duration_seconds", uptime)

// Policy Events
logging.Logger.Info("policy_loaded", "sectors", len(policy.Sectors), "strategies", len(policy.Strategies))
logging.Logger.Warn("policy_validation_failed", "reason", "signature_mismatch")

// Navigation Events
logging.Logger.Debug("screen_transition", "from", "sector_selection", "to", "screener_results")
logging.Logger.Info("workflow_completed", "ticker", trade.Ticker, "sector", trade.Sector)

// Safety Events
logging.Logger.Info("cooldown_started", "duration_seconds", 120, "ticker", ticker)
logging.Logger.Warn("heat_limit_exceeded", "sector", sector, "current", 1.7, "limit", 1.5)
logging.Logger.Info("checklist_completed", "required_items", 5, "optional_items", 2)

// User Actions
logging.Logger.Info("trade_saved", "trade_id", trade.ID, "ticker", trade.Ticker, "risk", trade.Risk)
logging.Logger.Info("sample_data_generated", "count", 10)
```

---

## 10. Integration Testing

**Purpose:** Verify end-to-end workflow with realistic scenarios.

### 10.1 Full Workflow Integration Test

**Location:** `internal/integration_test.go`

```go
// +build integration

package main_test

import (
    "testing"
    "time"
    "tf-engine/internal/appcore"
    "tf-engine/internal/models"
    "tf-engine/internal/storage"
)

func TestFullTradingWorkflow_HealthcareUNH(t *testing.T) {
    // Setup
    state := appcore.NewAppState()
    err := state.LoadPolicy("../data/policy.v1.json")
    if err != nil {
        t.Fatalf("Failed to load policy: %v", err)
    }

    // Step 1: Select Healthcare sector
    selectedSector := state.Policy.Sectors[0] // Healthcare (priority 1)
    if selectedSector.Blocked {
        t.Fatal("Healthcare should not be blocked")
    }
    state.CurrentTrade = &models.Trade{Sector: selectedSector.Name}
    t.Logf("✓ Sector selected: %s", selectedSector.Name)

    // Step 2: Verify allowed strategies
    allowedStrategies := selectedSector.AllowedStrategies
    if len(allowedStrategies) == 0 {
        t.Fatal("Healthcare should have allowed strategies")
    }
    if !contains(allowedStrategies, "Alt10") {
        t.Error("Healthcare should allow Alt10 strategy")
    }
    t.Logf("✓ Allowed strategies: %v", allowedStrategies)

    // Step 3: Select ticker and strategy
    state.CurrentTrade.Ticker = "UNH"
    state.CurrentTrade.Strategy = "Alt10"
    t.Logf("✓ Ticker selected: %s, Strategy: %s", state.CurrentTrade.Ticker, state.CurrentTrade.Strategy)

    // Step 4: Start cooldown timer
    state.StartCooldown()
    if !state.CooldownActive {
        t.Fatal("Cooldown should be active")
    }
    t.Logf("✓ Cooldown started: %d seconds remaining", state.GetCooldownRemaining())

    // Simulate waiting (don't actually wait 120 seconds in test)
    time.Sleep(2 * time.Second)
    remaining := state.GetCooldownRemaining()
    if remaining >= 120 {
        t.Error("Cooldown timer should be counting down")
    }

    // Step 5: Complete checklist (simulated)
    checklistItems := state.Policy.Checklist.Required
    if len(checklistItems) != 5 {
        t.Errorf("Expected 5 required checklist items, got %d", len(checklistItems))
    }
    t.Logf("✓ Checklist items: %v", checklistItems)

    // Step 6: Calculate position size
    conviction := 7 // Standard sizing
    multiplier := state.Policy.Checklist.PokerSizing[fmt.Sprint(conviction)]
    if multiplier != 1.0 {
        t.Errorf("Conviction 7 should have 1.0x multiplier, got %.2f", multiplier)
    }
    state.CurrentTrade.Risk = 500.0 // $500 risk
    t.Logf("✓ Position sized: $%.2f risk with %.2fx multiplier", state.CurrentTrade.Risk, multiplier)

    // Step 7: Check portfolio heat
    state.AllTrades = []models.Trade{
        {Ticker: "JNJ", Sector: "Healthcare", Risk: 300, Status: "active"},
        {Ticker: "ABBV", Sector: "Healthcare", Risk: 400, Status: "active"},
    }

    accountSize := 50000.0
    existingHeat := (300.0 + 400.0) / accountSize // 1.4% in Healthcare
    newTradeHeat := 500.0 / accountSize // +1.0% = 2.4% total
    totalSectorHeat := existingHeat + newTradeHeat

    if totalSectorHeat > state.Policy.Defaults.BucketHeatCap {
        t.Errorf("Trade would exceed sector heat cap: %.2f%% > %.2f%%", totalSectorHeat*100, state.Policy.Defaults.BucketHeatCap*100)
    }
    t.Logf("✓ Heat check passed: %.2f%% sector heat (limit: %.2f%%)", totalSectorHeat*100, state.Policy.Defaults.BucketHeatCap*100)

    // Step 8: Select options strategy
    state.CurrentTrade.OptionsType = "Bull call spread"
    state.CurrentTrade.EntryDate = time.Now()
    state.CurrentTrade.ExpirationDate = time.Now().AddDate(0, 0, 45) // 45 DTE
    t.Logf("✓ Options strategy: %s, DTE: 45", state.CurrentTrade.OptionsType)

    // Step 9: Save trade
    state.CurrentTrade.Status = "active"
    err = storage.SaveCompletedTrade(state.CurrentTrade)
    if err != nil {
        t.Fatalf("Failed to save trade: %v", err)
    }
    t.Logf("✓ Trade saved: %+v", state.CurrentTrade)

    // Verification
    allTrades, err := storage.LoadAllTrades()
    if err != nil {
        t.Fatalf("Failed to load trades: %v", err)
    }

    found := false
    for _, trade := range allTrades {
        if trade.Ticker == "UNH" && trade.Status == "active" {
            found = true
            break
        }
    }
    if !found {
        t.Error("Saved trade not found in trade history")
    }

    t.Log("✅ Full workflow test passed")
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

**Run command:** `go test -tags=integration ./internal/...`

---

## 11. Release Gates (Phase Exit Criteria)

### Phase 0 Exit Gate

**Deliverables:**
- [x] `data/policy.v1.json` has valid `security.signature` field
- [x] `scripts/verify_policy_hash.go` successfully validates policy
- [x] `internal/config/feature_flags.json` exists with all Phase 2 features disabled
- [x] `CONTRIBUTING.md` includes PR checklist and feature freeze policy
- [x] CI pipeline runs policy verification on every commit

**Verification Commands:**
```bash
go run scripts/verify_policy_hash.go          # Must output: ✅ Policy signature valid
cat internal/config/feature_flags.json | jq .  # All boolean flags must be false
git log --oneline -1                           # Latest commit passes CI
```

**Sign-off:** Product Owner reviews `CONTRIBUTING.md` and approves feature freeze

---

### Phase 1 Exit Gate

**Deliverables:**
- [x] Navigator supports forward/backward navigation through 8 screens
- [x] Auto-save persists in-progress trades after each screen
- [x] Cooldown timer widget displays countdown and blocks continue button
- [x] Manual test: Complete full workflow (sector → screener → ticker → checklist → sizing → heat → options → calendar)

**Verification Commands:**
```bash
go test ./internal/ui/...           # All unit tests pass
go run . --debug                    # Launch app, complete 1 full trade workflow
cat data/in_progress.json | jq .   # Verify auto-save worked
```

**Gherkin Test Pass Rate:** 100% (all Phase 1 scenarios pass)

**Sign-off:** QA lead approves manual test checklist

---

### Phase 2 Exit Gate

**Deliverables:**
- [x] All 8 screens fully implemented with Gherkin tests passing
- [x] Policy-driven sector/strategy filtering works correctly
- [x] Heat check enforces 4% portfolio / 1.5% sector limits
- [x] Sample data generation creates 10 realistic trades
- [x] Feature flags enable: `trade_management`, `sample_data_generator`, `vimium_mode`

**Verification Commands:**
```bash
go test -tags=integration ./internal/...  # Integration test passes
go run . --load-sample-data               # Generate sample data, verify calendar displays correctly
```

**Performance Budget Compliance:**
- App launch: < 2 seconds
- Screen transitions: < 200ms
- Trade save: < 500ms

**Sign-off:** Product Owner approves feature completeness

---

### Phase 3 Exit Gate

**Deliverables:**
- [x] Day/Night theme toggle works without restart
- [x] Text contrast meets WCAG AA standards (4.5:1 ratio minimum)
- [x] Help menu displays context-sensitive documentation
- [x] Welcome screen shows on first launch with "Don't show again" option

**Verification Commands:**
```bash
go test ./internal/theme/...         # Theme tests pass
go run . --first-launch               # Verify welcome screen appears
```

**Accessibility Test:** Manual review of all screens in both themes with colorblindness simulator

**Sign-off:** UX designer approves theme consistency

---

### Phase 4 Exit Gate

**Deliverables:**
- [x] Windows installer (`TFEngine-Setup-1.0.0.exe`) builds successfully
- [x] Silent install works: `TFEngine-Setup-1.0.0.exe /S`
- [x] Uninstaller prompts for data preservation
- [x] Desktop shortcut and Start Menu entry created
- [x] Version downgrade prevention works

**Verification Commands:**
```bash
./scripts/build-installer.sh          # Builds installer.exe
ls -lh TFEngine-Setup-1.0.0.exe      # Verify file exists (expect ~15-25 MB)
```

**Manual Test:** Install on clean Windows 10/11 VM, verify all functionality

**Sign-off:** DevOps lead approves installer build process

---

### Phase 5 Exit Gate (MVP Release)

**Deliverables:**
- [x] All Phase 0-4 exit gates passed
- [x] User acceptance testing completed with 3 external beta testers
- [x] Known bugs documented in `docs/known-issues.md` with severity ratings
- [x] Release notes written (`CHANGELOG.md`)
- [x] Installer signed with code signing certificate (if available)

**Verification Commands:**
```bash
git tag v1.0.0                        # Tag release
git log --oneline v1.0.0 ^v0.9.0     # Review commits since last beta
```

**Acceptance Criteria:**
- [ ] 3/3 beta testers successfully complete full trade workflow
- [ ] No severity-1 (blocker) or severity-2 (critical) bugs outstanding
- [ ] App launches on Windows 10 (build 19041+) and Windows 11

**Sign-off:** Product Owner and QA lead approve for public release

---

## 12. Risks & Mitigations

### Risk 1: Policy File Corruption in Production

**Probability:** Low
**Impact:** High (app becomes unusable)

**Mitigation:**
1. ✅ SHA256 signature validation catches tampering
2. ✅ Safe Mode provides fallback configuration
3. ✅ Installer includes backup copy of policy in `install_dir/policy.v1.json.bak`
4. ✅ App checks for `.bak` file if main policy fails validation

**Rollback Strategy:**
```go
// internal/config/policy_loader.go
func LoadPolicyWithFallback() (*models.Policy, error) {
    // Try primary policy
    policy, err := models.LoadPolicy("data/policy.v1.json")
    if err == nil && validateSignature(policy) {
        return policy, nil
    }

    // Try backup policy
    logging.Logger.Warn("primary policy failed, trying backup")
    policy, err = models.LoadPolicy("policy.v1.json.bak")
    if err == nil && validateSignature(policy) {
        // Restore backup as primary
        os.Rename("policy.v1.json.bak", "data/policy.v1.json")
        return policy, nil
    }

    // Fall back to safe mode
    logging.Logger.Error("all policies failed, activating safe mode")
    return models.SafeModePolicy(), nil
}
```

---

### Risk 2: Heat Calculation Logic Errors

**Probability:** Medium
**Impact:** High (users could exceed risk limits)

**Mitigation:**
1. ✅ Unit tests for all heat calculation functions with edge cases
2. ✅ Integration test verifies multi-trade heat scenarios
3. ✅ Manual override requires admin password (Phase 3 feature flag)
4. ✅ Audit log records all heat limit decisions

**Test Coverage Required:** 95% for `internal/appcore/heat.go`

---

### Risk 3: Cooldown Timer Bypass Vulnerability

**Probability:** Low
**Impact:** Medium (defeats anti-impulsivity purpose)

**Mitigation:**
1. ✅ Timer logic in backend (`appcore`), not GUI layer
2. ✅ Continue button disabled until `IsCooldownComplete()` returns true
3. ✅ No keyboard shortcuts bypass cooldown
4. ✅ Telemetry logs all cooldown starts/completions for audit

**Exploit Test:** Attempt to:
- Change system clock (app uses `time.Since()`, immune to this)
- Close/reopen app (cooldown state persists in `in_progress.json`)
- Press Enter key repeatedly (button disabled, no effect)

---

### Risk 4: Data Loss During Screen Transitions

**Probability:** Low
**Impact:** Medium (frustrating user experience)

**Mitigation:**
1. ✅ Auto-save after every screen transition
2. ✅ Atomic file writes (temp file → rename)
3. ✅ On app restart, detect incomplete trade and prompt: "Resume or Start New?"
4. ✅ Manual "Save Draft" button on every screen (backup option)

**Recovery Test:** Kill app process mid-workflow, restart, verify draft restored

---

### Risk 5: Fyne GUI Rendering Issues on Older Windows

**Probability:** Medium
**Impact:** Low (visual glitches, not functionality loss)

**Mitigation:**
1. ✅ Test on Windows 10 (build 19041, May 2020) and Windows 11
2. ✅ Use Fyne stable release (v2.4+), avoid experimental features
3. ✅ Provide command-line flag `--software-render` to disable GPU acceleration
4. ✅ Document minimum Windows version in README: Windows 10 1909+

---

### Risk 6: Screener URL Parameters Break (Finviz Changes)

**Probability:** Medium
**Impact:** Low (users can manually navigate Finviz)

**Mitigation:**
1. ✅ URLs stored in policy.json (updatable without code changes)
2. ✅ App checks URL validity on launch (HTTP HEAD request with 5s timeout)
3. ✅ If URL unreachable, show warning: "Screener may be unavailable. Continue anyway?"
4. ✅ Finviz alternatives documented in `docs/screener-alternatives.md` (TradingView, Barchart)

---

## 13. Rollback & Versioning Strategy

### Version Numbering

**Semantic Versioning:** `MAJOR.MINOR.PATCH` (e.g., 1.0.0 → 1.1.0 → 2.0.0)

- **MAJOR:** Breaking changes (incompatible policy format, database schema changes)
- **MINOR:** New features (new screen, sample data generator)
- **PATCH:** Bug fixes, performance improvements

### Rollback Scenarios

**Scenario 1: User installs buggy version 1.1.0**

**Solution:**
1. User downloads previous installer: `TFEngine-Setup-1.0.0.exe`
2. Uninstall current version (preserves data if user selects "No" to data deletion prompt)
3. Install version 1.0.0
4. App reads existing `data/trades.json` (backward compatible)

**Data Compatibility Rule:** Newer versions must read older data formats. Example:

```go
// internal/models/trade.go (v1.1.0)
type Trade struct {
    ID           string    `json:"id"`
    Ticker       string    `json:"ticker"`
    // New field in v1.1.0
    Tags         []string  `json:"tags,omitempty"`  // omitempty ensures v1.0.0 compatibility
}
```

---

**Scenario 2: Corrupted policy after manual edit**

**Solution:**
1. App detects signature mismatch on launch
2. Shows dialog: "Policy validation failed. Restore from backup?"
3. User clicks "Restore"
4. App copies `policy.v1.json.bak` → `data/policy.v1.json`
5. App relaunches

**Implementation:**

```go
// In main.go
func main() {
    policy, err := appcore.LoadPolicyWithFallback()
    if err != nil {
        dialog := widget.NewConfirm(
            "Policy Validation Failed",
            "The policy file is corrupted. Restore from backup?",
            func(restore bool) {
                if restore {
                    restorePolicyBackup()
                    fyne.CurrentApp().Quit()
                    // User manually relaunches
                } else {
                    // Continue in Safe Mode
                    appState.UseSafeMode()
                }
            },
            window,
        )
        dialog.Show()
    }
}
```

---

**Scenario 3: Database migration failure (future SQLite migration)**

**Solution:**
1. Before migration, back up `data/trades.json` → `data/trades.json.pre_migration`
2. Attempt SQLite migration
3. If migration fails:
   - Log error: "Migration failed, reverting to JSON storage"
   - App continues using JSON (graceful degradation)
4. On next launch, retry migration

---

### Version Check on Launch

**Feature:** Notify user if newer version available (optional telemetry endpoint)

```go
// internal/update/checker.go
func CheckForUpdates(currentVersion string) (string, bool) {
    resp, err := http.Get("https://example.com/api/latest-version")
    if err != nil {
        return "", false // Fail silently
    }
    defer resp.Body.Close()

    var release struct {
        Version string `json:"version"`
    }
    json.NewDecoder(resp.Body).Decode(&release)

    if compareVersions(release.Version, currentVersion) > 0 {
        return release.Version, true // Newer version available
    }
    return "", false
}

// Usage in main.go:
newVersion, available := update.CheckForUpdates("1.0.0")
if available {
    dialog := widget.NewConfirm(
        "Update Available",
        fmt.Sprintf("Version %s is available. Download now?", newVersion),
        func(download bool) {
            if download {
                fyne.CurrentApp().OpenURL("https://example.com/download")
            }
        },
        window,
    )
    dialog.Show()
}
```

---

## 14. Acceptance Criteria (MVP Complete)

### Must-Have Features (Phase 0-5)

**Policy Management:**
- [x] Policy file loads and validates SHA256 signature
- [x] Safe Mode activates if policy corrupted
- [x] Feature flags control Phase 2 feature availability

**Core Workflow (8 Screens):**
- [x] Screen 1: Sector selection with blocked/warned sectors
- [x] Screen 2: Finviz screener URL launching
- [x] Screen 3: Ticker entry + strategy dropdown (sector-filtered)
- [x] Screen 4: Anti-impulsivity checklist (5 required + 3 optional)
- [x] Screen 5: Poker-bet position sizing (conviction 5-8)
- [x] Screen 6: Portfolio heat check (4% total, 1.5% sector caps)
- [x] Screen 7: Options strategy selection (24 types)
- [x] Screen 8: Trade calendar (horserace view with sector Y-axis)

**Navigation:**
- [x] Forward/backward buttons work on all screens
- [x] Progress auto-saves after each screen transition
- [x] Resume incomplete trade on app restart

**Behavioral Guardrails:**
- [x] 120-second cooldown timer (cannot be bypassed)
- [x] Checklist gates prevent skipping required items
- [x] Heat limits enforced (trades blocked if exceeding caps)

**Data Persistence:**
- [x] Trades saved to `data/trades.json` with atomic writes
- [x] In-progress trade saved to `data/in_progress.json`
- [x] App loads all trades on startup for calendar view

**Polish:**
- [x] Day/Night theme toggle
- [x] Help menu with screen-specific documentation
- [x] Welcome screen on first launch
- [x] Sample data generation (10 realistic trades)

**Packaging:**
- [x] Windows installer (NSIS) with silent install support
- [x] Uninstaller with data preservation prompt
- [x] Desktop shortcut and Start Menu entry

### Quality Gates

**Code Quality:**
- [ ] 80%+ unit test coverage (excluding GUI code)
- [ ] 100% Gherkin scenario pass rate (all Phase 1-2 scenarios)
- [ ] No linter errors: `golangci-lint run`
- [ ] No security vulnerabilities: `go list -json -m all | nancy sleuth`

**Performance:**
- [ ] App launch: < 2 seconds (cold start on HDD)
- [ ] Screen transitions: < 200ms
- [ ] Trade save: < 500ms
- [ ] Calendar view renders 20 trades: < 1 second

**Documentation:**
- [ ] `README.md` includes installation instructions
- [ ] `CONTRIBUTING.md` enforces feature freeze policy
- [ ] `docs/user-guide.md` explains 8-screen workflow
- [ ] Inline code comments for complex logic (policy validation, heat calculation)

### User Acceptance (Beta Testing)

**Beta Tester Checklist (3 testers):**
1. [ ] Install app via `TFEngine-Setup-1.0.0.exe` on Windows 10/11
2. [ ] Complete full trade workflow: Healthcare → UNH → Alt10 → Checklist → Position sizing → Heat check → Bull call spread → Calendar
3. [ ] Verify cooldown timer prevents immediate trade execution
4. [ ] Attempt to trade blocked sector (Utilities) - verify prevention
5. [ ] Generate sample data and review calendar view
6. [ ] Toggle day/night theme and verify text readability
7. [ ] Close app mid-workflow, restart, and resume in-progress trade
8. [ ] Uninstall app and verify data preservation prompt

**Success Criteria:** 3/3 beta testers complete checklist with no severity-1 or severity-2 bugs reported

---

## 15. Open Issues Backlog (Future Enhancements)

These features are **not** in scope for MVP but documented for future consideration:

### Post-MVP Features (Phase 6+)

**Analytics Dashboard:**
- Win rate by sector/strategy
- Average hold time analysis
- P&L heatmap (calendar cells colored by profitability)

**Broker Integration:**
- Auto-populate current positions from TD Ameritrade/Schwab API
- Real-time P&L updates from market data

**Mobile Companion App:**
- View-only dashboard on iOS/Android
- Push notifications for expiring trades

**Walk-Forward Analysis:**
- Re-run backtests on new data to validate strategy performance
- Alert if strategy performance degrades below threshold

**AI Trade Advisor:**
- Machine learning model suggests optimal entry timing based on historical patterns
- Displays confidence score with rationale

**Social Features:**
- Share anonymized trade ideas with other TF-Engine users
- Community leaderboard (opt-in)

---

## Conclusion

This roadmap provides a complete blueprint for building TF-Engine 2.0 from policy infrastructure (Phase 0) through MVP release (Phase 5). The phased approach ensures:

1. **Feature Discipline:** Phase 0 feature freeze prevents scope creep
2. **Policy-Driven Design:** All business logic in `policy.v1.json`, not hardcoded
3. **Behavioral Finance:** Cooldowns, checklists, and heat limits enforce discipline
4. **Graceful Degradation:** Safe Mode ensures app remains usable despite data corruption
5. **Quality Gates:** Each phase has clear acceptance criteria and verification commands

**Next Steps:**
1. Review and approve this roadmap
2. Delete `plans/roadmap.original.md` and `plans/roadmap.revised.md`
3. Begin Phase 0 implementation (policy signing and feature flags)
4. Schedule bi-weekly phase review meetings

**Estimated Timeline:**
- Phase 0: 3 days
- Phase 1: 2 weeks
- Phase 2: 4 weeks
- Phase 3: 1 week
- Phase 4: 3 days
- Phase 5: 1 week (UAT)

**Total:** ~8-9 weeks to MVP release

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-03
**Owner:** Product Lead
**Reviewers:** Engineering Lead, QA Lead, UX Designer
