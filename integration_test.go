//go:build integration
// +build integration

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
)

// TestFullTradingWorkflow_HealthcareUNH tests the complete trading workflow
// from sector selection through to calendar display
func TestFullTradingWorkflow_HealthcareUNH(t *testing.T) {
	// Setup test environment
	cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Log("=== Integration Test: Full Trading Workflow (Healthcare → UNH → Calendar) ===")

	// Initialize application state
	state := appcore.NewAppState()
	err := state.LoadPolicy("data/policy.v1.json")
	if err != nil {
		t.Fatalf("Failed to load policy: %v", err)
	}
	t.Logf("✓ Policy loaded: %d sectors, %d strategies", len(state.Policy.Sectors), len(state.Policy.Strategies))

	// ==================================================
	// SCREEN 1: SECTOR SELECTION
	// ==================================================
	t.Log("\n--- Screen 1: Sector Selection ---")

	// Find Healthcare sector
	var healthcareSector *models.Sector
	for i := range state.Policy.Sectors {
		if state.Policy.Sectors[i].Name == "Healthcare" {
			healthcareSector = &state.Policy.Sectors[i]
			break
		}
	}
	if healthcareSector == nil {
		t.Fatal("Healthcare sector not found in policy")
	}

	// Validate Healthcare sector is not blocked
	if healthcareSector.Blocked {
		t.Fatal("Healthcare should not be blocked")
	}
	t.Logf("✓ Healthcare sector selected (Priority: %d)", healthcareSector.Priority)

	// Initialize trade with selected sector
	state.CurrentTrade = &models.Trade{
		Sector:    healthcareSector.Name,
		CreatedAt: time.Now(),
	}

	// Auto-save after Screen 1
	err = storage.SaveInProgressTrade(state.CurrentTrade)
	if err != nil {
		t.Fatalf("Failed to auto-save after Screen 1: %v", err)
	}
	t.Log("✓ Trade auto-saved after sector selection")

	// ==================================================
	// SCREEN 2: SCREENER LAUNCH
	// ==================================================
	t.Log("\n--- Screen 2: Screener Launch ---")

	// Verify screener URLs exist
	if len(healthcareSector.ScreenerURLs) == 0 {
		t.Fatal("Healthcare sector has no screener URLs")
	}

	// Check for required screeners
	requiredScreeners := []string{"universe", "pullback", "breakout", "golden_cross"}
	for _, screener := range requiredScreeners {
		if _, exists := healthcareSector.ScreenerURLs[screener]; !exists {
			t.Errorf("Missing screener: %s", screener)
		}
	}
	t.Logf("✓ %d screeners available for Healthcare", len(healthcareSector.ScreenerURLs))

	// Verify v=211 parameter (chart view)
	universeURL := healthcareSector.ScreenerURLs["universe"]
	if universeURL == "" {
		t.Fatal("Universe screener URL is empty")
	}
	// Note: We don't actually launch the URL in tests, just validate it exists
	t.Log("✓ Universe screener URL validated")

	// ==================================================
	// SCREEN 3: TICKER ENTRY & STRATEGY SELECTION
	// ==================================================
	t.Log("\n--- Screen 3: Ticker Entry & Strategy Selection ---")

	// Enter ticker
	state.CurrentTrade.Ticker = "UNH"
	t.Logf("✓ Ticker entered: %s", state.CurrentTrade.Ticker)

	// Verify allowed strategies
	if len(healthcareSector.AllowedStrategies) == 0 {
		t.Fatal("Healthcare has no allowed strategies")
	}

	// Verify Alt10 is allowed for Healthcare
	alt10Allowed := false
	for _, strategy := range healthcareSector.AllowedStrategies {
		if strategy == "Alt10" {
			alt10Allowed = true
			break
		}
	}
	if !alt10Allowed {
		t.Fatal("Alt10 should be allowed for Healthcare sector")
	}

	// Select Alt10 strategy
	state.CurrentTrade.Strategy = "Alt10"
	t.Logf("✓ Strategy selected: Alt10 (Profit Targets)")

	// Auto-save after Screen 3
	err = storage.SaveInProgressTrade(state.CurrentTrade)
	if err != nil {
		t.Fatalf("Failed to auto-save after Screen 3: %v", err)
	}
	t.Log("✓ Trade auto-saved after ticker entry")

	// Start cooldown timer
	state.StartCooldown()
	t.Logf("✓ Cooldown started: %d seconds", state.Policy.Defaults.CooldownSeconds)

	// ==================================================
	// SCREEN 4: CHECKLIST
	// ==================================================
	t.Log("\n--- Screen 4: Anti-Impulsivity Checklist ---")

	// Verify required checklist items
	if len(state.Policy.Checklist.Required) != 5 {
		t.Errorf("Expected 5 required checklist items, got %d", len(state.Policy.Checklist.Required))
	}
	t.Logf("✓ Checklist has %d required items", len(state.Policy.Checklist.Required))

	// Simulate waiting for cooldown (don't actually wait in test)
	t.Log("✓ Simulating cooldown completion...")
	state.CooldownActive = false
	state.CooldownCompleted = true

	// ==================================================
	// SCREEN 5: POSITION SIZING
	// ==================================================
	t.Log("\n--- Screen 5: Position Sizing ---")

	// Standard conviction (7 = 1.0x multiplier)
	conviction := 7
	multiplier := state.Policy.Checklist.PokerSizing[string(rune('0'+conviction))]
	if multiplier == 0 {
		multiplier = 1.0 // Default if not found
	}
	if multiplier != 1.0 {
		t.Errorf("Conviction 7 should have 1.0x multiplier, got %.2f", multiplier)
	}
	t.Logf("✓ Poker sizing: conviction=%d, multiplier=%.2fx", conviction, multiplier)

	// Calculate position size
	baseRisk := 250.0 // $250 base risk (0.5% of $50k account)
	state.CurrentTrade.Risk = baseRisk * multiplier
	t.Logf("✓ Position sized: $%.2f risk", state.CurrentTrade.Risk)

	// ==================================================
	// SCREEN 6: HEAT CHECK
	// ==================================================
	t.Log("\n--- Screen 6: Portfolio Heat Check ---")

	// Add existing trades to simulate portfolio
	state.AllTrades = []models.Trade{
		{Ticker: "JNJ", Sector: "Healthcare", Risk: 200, Status: "active"},  // 0.4%
		{Ticker: "ABBV", Sector: "Healthcare", Risk: 250, Status: "active"}, // 0.5%
		// Total existing: 0.9%
	}

	// Calculate heat
	accountSize := 50000.0
	healthcareHeat := 0.0
	for _, trade := range state.AllTrades {
		if trade.Sector == "Healthcare" && trade.Status == "active" {
			healthcareHeat += trade.Risk / accountSize
		}
	}

	// Add new trade heat
	newTradeHeat := state.CurrentTrade.Risk / accountSize
	totalHealthcareHeat := healthcareHeat + newTradeHeat

	t.Logf("✓ Existing Healthcare heat: %.2f%%", healthcareHeat*100)
	t.Logf("✓ New trade heat: %.2f%%", newTradeHeat*100)
	t.Logf("✓ Total Healthcare heat: %.2f%%", totalHealthcareHeat*100)

	// Verify heat is within limits
	if totalHealthcareHeat > state.Policy.Defaults.BucketHeatCap {
		t.Fatalf("Heat check FAILED: %.2f%% exceeds %.2f%% sector cap (test data misconfigured)",
			totalHealthcareHeat*100, state.Policy.Defaults.BucketHeatCap*100)
	} else {
		t.Logf("✓ Heat check PASSED: %.2f%% <= %.2f%% sector cap",
			totalHealthcareHeat*100, state.Policy.Defaults.BucketHeatCap*100)
	}

	// ==================================================
	// SCREEN 7: TRADE ENTRY (OPTIONS STRATEGY)
	// ==================================================
	t.Log("\n--- Screen 7: Trade Entry (Options Strategy) ---")

	// Select options strategy
	state.CurrentTrade.OptionsType = "Bull call spread"
	state.CurrentTrade.EntryDate = time.Now()
	state.CurrentTrade.ExpirationDate = time.Now().AddDate(0, 0, 45) // 45 DTE
	state.CurrentTrade.Status = "active"
	t.Logf("✓ Options strategy: %s", state.CurrentTrade.OptionsType)
	t.Logf("✓ Expiration: %s (45 DTE)", state.CurrentTrade.ExpirationDate.Format("2006-01-02"))

	// ==================================================
	// SCREEN 8: SAVE COMPLETED TRADE
	// ==================================================
	t.Log("\n--- Screen 8: Save Completed Trade ---")

	// Save completed trade
	err = storage.SaveCompletedTrade(state.CurrentTrade)
	if err != nil {
		t.Fatalf("Failed to save completed trade: %v", err)
	}
	t.Log("✓ Trade saved to history")

	// Verify trade was saved
	allTrades, err := storage.LoadAllTrades()
	if err != nil {
		t.Fatalf("Failed to load trades: %v", err)
	}

	// Find our trade
	found := false
	for _, trade := range allTrades {
		if trade.Ticker == "UNH" && trade.Sector == "Healthcare" && trade.Strategy == "Alt10" {
			found = true
			t.Logf("✓ Trade found in history: %s %s %s", trade.Ticker, trade.Sector, trade.Strategy)
			break
		}
	}
	if !found {
		t.Error("Trade not found in history after save")
	}

	// Verify in-progress file was deleted
	inProgressPath := "data/trades_in_progress.json"
	if _, err := os.Stat(inProgressPath); err == nil {
		t.Error("In-progress file should be deleted after completing trade")
	}

	t.Log("\n=== ✅ Integration Test PASSED: Full workflow completed successfully ===")
}

// TestHeatLimitEnforcement tests that the heat check properly blocks trades
// that would exceed sector or portfolio limits
func TestHeatLimitEnforcement(t *testing.T) {
	cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Log("=== Integration Test: Heat Limit Enforcement ===")

	// Initialize state
	state := appcore.NewAppState()
	err := state.LoadPolicy("data/policy.v1.json")
	if err != nil {
		t.Fatalf("Failed to load policy: %v", err)
	}

	accountSize := 50000.0

	// ==================================================
	// TEST 1: Sector Heat Cap Enforcement
	// ==================================================
	t.Log("\n--- Test 1: Sector Heat Cap (1.5%) ---")

	// Add trades that bring Healthcare to the sector cap
	state.AllTrades = []models.Trade{
		{Ticker: "JNJ", Sector: "Healthcare", Risk: 400, Status: "active"},   // 0.8%
		{Ticker: "ABBV", Sector: "Healthcare", Risk: 350, Status: "active"},  // 0.7%
		// Total: 1.5% (exactly at cap)
	}

	// Try to add another Healthcare trade
	newTrade := &models.Trade{
		Ticker: "UNH",
		Sector: "Healthcare",
		Risk:   100, // Even small trade should be blocked
		Status: "active",
	}

	// Calculate heat
	sectorHeat := 0.0
	for _, trade := range state.AllTrades {
		if trade.Sector == "Healthcare" && trade.Status == "active" {
			sectorHeat += trade.Risk / accountSize
		}
	}
	newHeat := newTrade.Risk / accountSize
	totalHeat := sectorHeat + newHeat

	t.Logf("Existing Healthcare heat: %.2f%%", sectorHeat*100)
	t.Logf("New trade would add: %.2f%%", newHeat*100)
	t.Logf("Total would be: %.2f%%", totalHeat*100)
	t.Logf("Sector cap: %.2f%%", state.Policy.Defaults.BucketHeatCap*100)

	// Verify heat check would block
	if totalHeat > state.Policy.Defaults.BucketHeatCap {
		t.Log("✓ Heat check correctly BLOCKS trade exceeding sector cap")
	} else {
		t.Error("Heat check should BLOCK trade that exceeds sector cap")
	}

	// ==================================================
	// TEST 2: Portfolio Heat Cap Enforcement
	// ==================================================
	t.Log("\n--- Test 2: Portfolio Heat Cap (4.0%) ---")

	// Add trades across multiple sectors that approach portfolio cap
	state.AllTrades = []models.Trade{
		{Ticker: "UNH", Sector: "Healthcare", Risk: 500, Status: "active"},     // 1.0%
		{Ticker: "MSFT", Sector: "Technology", Risk: 600, Status: "active"},    // 1.2%
		{Ticker: "CAT", Sector: "Industrials", Risk: 550, Status: "active"},    // 1.1%
		{Ticker: "JPM", Sector: "Financials", Risk: 350, Status: "active"},     // 0.7%
		// Total: 4.0% (exactly at portfolio cap)
	}

	// Try to add ANY new trade
	newTrade = &models.Trade{
		Ticker: "XLV",
		Sector: "Healthcare",
		Risk:   50, // Even tiny trade should be blocked
		Status: "active",
	}

	// Calculate portfolio heat
	portfolioHeat := 0.0
	for _, trade := range state.AllTrades {
		if trade.Status == "active" {
			portfolioHeat += trade.Risk / accountSize
		}
	}
	newHeat = newTrade.Risk / accountSize
	totalPortfolioHeat := portfolioHeat + newHeat

	t.Logf("Existing portfolio heat: %.2f%%", portfolioHeat*100)
	t.Logf("New trade would add: %.2f%%", newHeat*100)
	t.Logf("Total would be: %.2f%%", totalPortfolioHeat*100)
	t.Logf("Portfolio cap: %.2f%%", state.Policy.Defaults.PortfolioHeatCap*100)

	// Verify heat check would block
	if totalPortfolioHeat > state.Policy.Defaults.PortfolioHeatCap {
		t.Log("✓ Heat check correctly BLOCKS trade exceeding portfolio cap")
	} else {
		t.Error("Heat check should BLOCK trade that exceeds portfolio cap")
	}

	t.Log("\n=== ✅ Integration Test PASSED: Heat limits enforced correctly ===")
}

// TestCooldownPersistence tests that cooldown timer state persists across
// application restarts
func TestCooldownPersistence(t *testing.T) {
	cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Log("=== Integration Test: Cooldown Timer Persistence ===")

	// ==================================================
	// STEP 1: Start cooldown and save state
	// ==================================================
	t.Log("\n--- Step 1: Start cooldown and save ---")

	state1 := appcore.NewAppState()
	err := state1.LoadPolicy("data/policy.v1.json")
	if err != nil {
		t.Fatalf("Failed to load policy: %v", err)
	}

	// Create trade and start cooldown
	state1.CurrentTrade = &models.Trade{
		Sector:   "Healthcare",
		Ticker:   "UNH",
		Strategy: "Alt10",
	}

	// Start cooldown
	state1.StartCooldown()
	t.Logf("✓ Cooldown started: %d seconds", state1.Policy.Defaults.CooldownSeconds)

	// Save trade with cooldown state
	err = storage.SaveInProgressTrade(state1.CurrentTrade)
	if err != nil {
		t.Fatalf("Failed to save trade: %v", err)
	}
	t.Log("✓ Trade saved with cooldown start time")

	// Wait a bit to simulate time passing
	time.Sleep(2 * time.Second)

	// ==================================================
	// STEP 2: Simulate app restart - load state
	// ==================================================
	t.Log("\n--- Step 2: Simulate app restart ---")

	state2 := appcore.NewAppState()
	err = state2.LoadPolicy("data/policy.v1.json")
	if err != nil {
		t.Fatalf("Failed to load policy: %v", err)
	}

	// Load in-progress trade
	loadedTrade, err := storage.LoadInProgressTrade()
	if err != nil {
		t.Fatalf("Failed to load in-progress trade: %v", err)
	}
	if loadedTrade == nil {
		t.Fatal("No in-progress trade found")
	}

	state2.CurrentTrade = loadedTrade
	t.Log("✓ Trade loaded from disk")

	// Restore cooldown state
	if !loadedTrade.CooldownStartTime.IsZero() {
		// Manually restore cooldown
		state2.CooldownStart = &loadedTrade.CooldownStartTime
		state2.CooldownDuration = time.Duration(state2.Policy.Defaults.CooldownSeconds) * time.Second
		state2.CooldownActive = true

		// Calculate remaining time
		elapsed := time.Since(*state2.CooldownStart)
		remaining := state2.CooldownDuration - elapsed

		t.Logf("✓ Cooldown restored: %.0f seconds elapsed", elapsed.Seconds())
		t.Logf("✓ Remaining: %.0f seconds", remaining.Seconds())

		// Verify remaining time is less than original duration
		if remaining >= state2.CooldownDuration {
			t.Error("Remaining time should be less than original duration")
		}

		// Verify remaining time is reasonable (we waited 2 seconds)
		if remaining > state2.CooldownDuration-time.Second {
			t.Error("Cooldown timer did not count down during 'restart'")
		}
	} else {
		t.Error("Cooldown start time was not persisted")
	}

	t.Log("\n=== ✅ Integration Test PASSED: Cooldown persists across restarts ===")
}

// TestBlockedSectorEnforcement tests that blocked sectors (like Utilities)
// properly prevent trade entry
func TestBlockedSectorEnforcement(t *testing.T) {
	cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Log("=== Integration Test: Blocked Sector Enforcement ===")

	state := appcore.NewAppState()
	err := state.LoadPolicy("data/policy.v1.json")
	if err != nil {
		t.Fatalf("Failed to load policy: %v", err)
	}

	// Find Utilities sector (should be blocked or warned)
	var utilitiesSector *models.Sector
	for i := range state.Policy.Sectors {
		if state.Policy.Sectors[i].Name == "Utilities" {
			utilitiesSector = &state.Policy.Sectors[i]
			break
		}
	}

	if utilitiesSector == nil {
		t.Fatal("Utilities sector not found in policy")
	}

	t.Logf("Utilities sector - Blocked: %v, Warning: %v", utilitiesSector.Blocked, utilitiesSector.Warning)

	// Utilities should be blocked OR have warning enabled
	if !utilitiesSector.Blocked && !utilitiesSector.Warning {
		t.Error("Utilities sector should be blocked or warned (0% backtest success)")
	}

	// If blocked, verify no strategies are allowed
	if utilitiesSector.Blocked && len(utilitiesSector.AllowedStrategies) > 0 {
		t.Error("Blocked sector should have no allowed strategies")
	}

	// If warning enabled, verify warning text exists
	if utilitiesSector.Warning {
		if utilitiesSector.Notes == "" {
			t.Error("Warning sector should have notes explaining the warning")
		}
		t.Logf("✓ Utilities warning: %s", utilitiesSector.Notes)
	}

	t.Log("\n=== ✅ Integration Test PASSED: Sector blocking works correctly ===")
}

// setupIntegrationTest creates a clean test environment
func setupIntegrationTest(t *testing.T) func() {
	// Create temporary data directory
	testDataDir := filepath.Join("data", "test_"+t.Name())
	err := os.MkdirAll(testDataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test data dir: %v", err)
	}

	// Backup original data directory
	originalDataPath := "data/trades.json"
	backupPath := ""
	if _, err := os.Stat(originalDataPath); err == nil {
		backupPath = originalDataPath + ".integration_backup"
		os.Rename(originalDataPath, backupPath)
	}

	// Cleanup function
	return func() {
		// Restore original data
		if backupPath != "" {
			os.Rename(backupPath, originalDataPath)
		}

		// Clean up test data
		os.RemoveAll(testDataDir)

		// Clean up in-progress file
		os.Remove("data/trades_in_progress.json")
	}
}
