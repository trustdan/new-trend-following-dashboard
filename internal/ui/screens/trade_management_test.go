package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"tf-engine/internal/appcore"
	"tf-engine/internal/config"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
)

// setupTestDataDir creates a temporary test data directory
func setupTestDataDir(t *testing.T) func() {
	// Create temp test directory
	testDir := filepath.Join("testdata", fmt.Sprintf("test_%d", time.Now().UnixNano()))
	os.MkdirAll(testDir, 0755)

	// Ensure data directory exists for tests
	os.MkdirAll("data", 0755)
	os.MkdirAll("data/backups", 0755)

	// Return cleanup function
	return func() {
		os.RemoveAll("data/trades_in_progress.json")
		os.RemoveAll("data/trades.json")
		os.RemoveAll("data/backups")
		os.RemoveAll(testDir)
	}
}

func TestTradeManagement_RenderWithFeatureFlagEnabled(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true, Phase: 2},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render returned nil")
	}
}

func TestTradeManagement_RenderWithFeatureFlagDisabled(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Version: "1.0.0",
		Flags: map[string]config.FeatureFlag{
			"trade_management": {
				Enabled:      false,
				Phase:        2,
				Description:  "Screen 9: Edit/delete trades",
				SinceVersion: "2.1.0",
			},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	content := screen.Render()

	// Assert
	if content == nil {
		t.Fatal("Render returned nil even when disabled")
	}
	// Should show disabled state message
}

func TestTradeManagement_GetFilteredTrades_All(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Create sample trades
	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", Status: "active"},
		{ID: "2", Ticker: "MSFT", Status: "closed"},
		{ID: "3", Ticker: "GOOGL", Status: "active"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)
	screen.filterStatus = "all"

	// Act
	filtered := screen.getFilteredTrades()

	// Assert
	if len(filtered) != 3 {
		t.Errorf("Expected 3 trades, got %d", len(filtered))
	}
}

func TestTradeManagement_GetFilteredTrades_ActiveOnly(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", Status: "active"},
		{ID: "2", Ticker: "MSFT", Status: "closed"},
		{ID: "3", Ticker: "GOOGL", Status: "active"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)
	screen.filterStatus = "active"

	// Act
	filtered := screen.getFilteredTrades()

	// Assert
	if len(filtered) != 2 {
		t.Errorf("Expected 2 active trades, got %d", len(filtered))
	}

	// Verify all are active
	for _, trade := range filtered {
		if trade.Status != "active" {
			t.Errorf("Expected active status, got %s", trade.Status)
		}
	}
}

func TestTradeManagement_GetFilteredTrades_ClosedOnly(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", Status: "active"},
		{ID: "2", Ticker: "MSFT", Status: "closed"},
		{ID: "3", Ticker: "GOOGL", Status: "closed"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)
	screen.filterStatus = "closed"

	// Act
	filtered := screen.getFilteredTrades()

	// Assert
	if len(filtered) != 2 {
		t.Errorf("Expected 2 closed trades, got %d", len(filtered))
	}

	// Verify all are closed
	for _, trade := range filtered {
		if trade.Status != "closed" {
			t.Errorf("Expected closed status, got %s", trade.Status)
		}
	}
}

func TestTradeManagement_UpdateTrade(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	pnl1 := 100.0
	pnl2 := -50.0
	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", ProfitLoss: &pnl1, Status: "active"},
		{ID: "2", Ticker: "MSFT", ProfitLoss: &pnl2, Status: "active"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Modify first trade
	pnlUpdated := 250.0
	updatedTrade := &models.Trade{
		ID:         "1",
		Ticker:     "AAPL",
		ProfitLoss: &pnlUpdated,
		Status:     "closed",
	}

	// Act
	err := screen.updateTrade(updatedTrade)

	// Assert
	if err != nil {
		t.Fatalf("updateTrade failed: %v", err)
	}

	// Verify trade was updated
	allTrades, _ := storage.LoadAllTrades()
	found := false
	for _, trade := range allTrades {
		if trade.ID == "1" {
			found = true
			if trade.GetPnL() != 250.0 {
				t.Errorf("Expected P&L 250.0, got %.2f", trade.GetPnL())
			}
			if trade.Status != "closed" {
				t.Errorf("Expected status 'closed', got '%s'", trade.Status)
			}
		}
	}

	if !found {
		t.Error("Updated trade not found in storage")
	}
}

func TestTradeManagement_UpdateTrade_NotFound(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", Status: "active"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Try to update non-existent trade
	nonExistentTrade := &models.Trade{
		ID:     "999",
		Ticker: "INVALID",
	}

	// Act
	err := screen.updateTrade(nonExistentTrade)

	// Assert
	if err == nil {
		t.Error("Expected error when updating non-existent trade")
	}
}

func TestTradeManagement_DeleteTrade(t *testing.T) {
	// Arrange
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trades := []models.Trade{
		{ID: "1", Ticker: "AAPL", Status: "active"},
		{ID: "2", Ticker: "MSFT", Status: "active"},
		{ID: "3", Ticker: "GOOGL", Status: "active"},
	}
	storage.SaveAllTrades(trades)

	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	err := screen.deleteTrade(&trades[1]) // Delete MSFT

	// Assert
	if err != nil {
		t.Fatalf("deleteTrade failed: %v", err)
	}

	// Verify trade was deleted
	allTrades, _ := storage.LoadAllTrades()
	if len(allTrades) != 2 {
		t.Errorf("Expected 2 trades remaining, got %d", len(allTrades))
	}

	// Verify MSFT is gone
	for _, trade := range allTrades {
		if trade.ID == "2" {
			t.Error("Deleted trade still exists")
		}
	}
}

func TestTradeManagement_CreateTradesTable_Empty(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	table := screen.createTradesTable([]models.Trade{})

	// Assert
	if table == nil {
		t.Fatal("createTradesTable returned nil for empty trades")
	}
}

func TestTradeManagement_CreateTradesTable_WithTrades(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	pnl := 150.0
	trades := []models.Trade{
		{
			ID:              "1",
			Ticker:          "AAPL",
			Sector:          "Technology",
			Strategy:        "Alt10",
			OptionsStrategy: "Bull call spread",
			ProfitLoss:      &pnl,
			Status:          "active",
			CreatedAt:       time.Now(),
		},
	}

	// Act
	table := screen.createTradesTable(trades)

	// Assert
	if table == nil {
		t.Fatal("createTradesTable returned nil")
	}
}

func TestTradeManagement_Validate(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	valid := screen.Validate()

	// Assert
	if !valid {
		t.Error("Trade management screen should always be valid (read-only)")
	}
}

func TestTradeManagement_GetName(t *testing.T) {
	// Arrange
	state := appcore.NewAppState()
	window := test.NewWindow(nil)
	flags := &config.FeatureFlags{
		Flags: map[string]config.FeatureFlag{
			"trade_management": {Enabled: true},
		},
	}

	screen := NewTradeManagement(state, window, flags)

	// Act
	name := screen.GetName()

	// Assert
	expected := "trade_management"
	if name != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, name)
	}
}
