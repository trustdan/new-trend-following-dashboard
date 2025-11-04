package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"tf-engine/internal/models"
)

func setupTestDataDir(t *testing.T) func() {
	// Create temp test directory
	testDir := filepath.Join("testdata", fmt.Sprintf("test_%d", time.Now().UnixNano()))
	os.MkdirAll(testDir, 0755)

	// Override constants for testing
	oldTradesFile := TradesFile
	oldInProgressFile := InProgressFile
	oldBackupDir := BackupDir

	// Update to use test directory (via reflection or exported setters)
	// For now, we'll use a simpler approach with relative paths
	os.MkdirAll("data", 0755)
	os.MkdirAll("data/backups", 0755)

	// Return cleanup function
	return func() {
		os.RemoveAll(InProgressFile)
		os.RemoveAll(TradesFile)
		os.RemoveAll(BackupDir)
		os.RemoveAll(testDir)

		// Restore original values
		_ = oldTradesFile
		_ = oldInProgressFile
		_ = oldBackupDir
	}
}

func TestSaveInProgressTrade_CreatesFile(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trade := &models.Trade{
		ID:     "test-123",
		Sector: "Healthcare",
		Ticker: "UNH",
	}

	err := SaveInProgressTrade(trade)
	if err != nil {
		t.Fatalf("SaveInProgressTrade failed: %v", err)
	}

	if !fileExists(InProgressFile) {
		t.Error("InProgressFile should exist after save")
	}

	// Verify JSON is valid
	data, _ := os.ReadFile(InProgressFile)
	var loaded models.Trade
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Errorf("JSON unmarshal failed: %v", err)
	}

	if loaded.Sector != "Healthcare" {
		t.Errorf("Expected sector Healthcare, got %s", loaded.Sector)
	}
}

func TestLoadInProgressTrade_RestoresData(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	originalTrade := &models.Trade{
		ID:       "test-456",
		Sector:   "Technology",
		Ticker:   "MSFT",
		Strategy: "Alt26",
	}

	err := SaveInProgressTrade(originalTrade)
	if err != nil {
		t.Fatalf("SaveInProgressTrade failed: %v", err)
	}

	loadedTrade, err := LoadInProgressTrade()
	if err != nil {
		t.Fatalf("LoadInProgressTrade failed: %v", err)
	}

	if loadedTrade == nil {
		t.Fatal("LoadInProgressTrade returned nil")
	}

	if loadedTrade.Sector != "Technology" {
		t.Errorf("Expected sector Technology, got %s", loadedTrade.Sector)
	}
	if loadedTrade.Ticker != "MSFT" {
		t.Errorf("Expected ticker MSFT, got %s", loadedTrade.Ticker)
	}
	if loadedTrade.Strategy != "Alt26" {
		t.Errorf("Expected strategy Alt26, got %s", loadedTrade.Strategy)
	}
}

func TestLoadInProgressTrade_NoFile_ReturnsNil(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	loadedTrade, err := LoadInProgressTrade()
	if err != nil {
		t.Errorf("LoadInProgressTrade should not error on missing file: %v", err)
	}

	if loadedTrade != nil {
		t.Error("LoadInProgressTrade should return nil when no file exists")
	}
}

func TestSaveCompletedTrade_AppendsToHistory(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Save first trade
	trade1 := &models.Trade{
		ID:     "trade-1",
		Ticker: "UNH",
		Sector: "Healthcare",
	}
	err := SaveCompletedTrade(trade1)
	if err != nil {
		t.Fatalf("SaveCompletedTrade failed: %v", err)
	}

	// Save second trade
	trade2 := &models.Trade{
		ID:     "trade-2",
		Ticker: "MSFT",
		Sector: "Technology",
	}
	err = SaveCompletedTrade(trade2)
	if err != nil {
		t.Fatalf("SaveCompletedTrade failed: %v", err)
	}

	// Load and verify
	trades, err := LoadAllTrades()
	if err != nil {
		t.Fatalf("LoadAllTrades failed: %v", err)
	}

	if len(trades) != 2 {
		t.Fatalf("Expected 2 trades, got %d", len(trades))
	}

	if trades[0].Ticker != "UNH" {
		t.Errorf("Expected first trade ticker UNH, got %s", trades[0].Ticker)
	}
	if trades[1].Ticker != "MSFT" {
		t.Errorf("Expected second trade ticker MSFT, got %s", trades[1].Ticker)
	}
}

func TestSaveCompletedTrade_CreatesBackup(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Save first trade
	trade1 := &models.Trade{ID: "trade-1", Ticker: "UNH"}
	SaveCompletedTrade(trade1)

	// Save second trade (should trigger backup of first)
	trade2 := &models.Trade{ID: "trade-2", Ticker: "MSFT"}
	err := SaveCompletedTrade(trade2)
	if err != nil {
		t.Fatalf("SaveCompletedTrade failed: %v", err)
	}

	// Check if backup directory exists and contains files
	if !fileExists(BackupDir) {
		t.Error("Backup directory should exist")
	}

	backups, err := filepath.Glob(filepath.Join(BackupDir, "trades_*.json"))
	if err != nil {
		t.Fatalf("Failed to glob backups: %v", err)
	}

	if len(backups) == 0 {
		t.Error("Expected at least one backup file")
	}
}

func TestSaveCompletedTrade_ClearsInProgress(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Save in-progress trade
	trade := &models.Trade{ID: "trade-1", Ticker: "UNH"}
	SaveInProgressTrade(trade)

	if !fileExists(InProgressFile) {
		t.Fatal("InProgressFile should exist")
	}

	// Complete the trade
	err := SaveCompletedTrade(trade)
	if err != nil {
		t.Fatalf("SaveCompletedTrade failed: %v", err)
	}

	// Verify in-progress file is deleted
	if fileExists(InProgressFile) {
		t.Error("InProgressFile should be deleted after completing trade")
	}
}

func TestConcurrentSaves_NoCorruption(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent saves to in-progress file
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			trade := &models.Trade{
				ID:     fmt.Sprintf("trade-%d", id),
				Ticker: fmt.Sprintf("TEST%d", id),
			}
			SaveInProgressTrade(trade)
		}(i)
	}

	wg.Wait()

	// Verify file is not corrupted
	trade, err := LoadInProgressTrade()
	if err != nil {
		t.Errorf("LoadInProgressTrade failed after concurrent writes: %v", err)
	}
	if trade == nil {
		t.Error("Trade should not be nil")
	}

	// Verify JSON is valid
	data, _ := os.ReadFile(InProgressFile)
	var loaded models.Trade
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Errorf("JSON corrupted after concurrent writes: %v", err)
	}
}

func TestLoadAllTrades_EmptyFile_ReturnsEmptySlice(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trades, err := LoadAllTrades()
	if err != nil {
		t.Errorf("LoadAllTrades should not error on missing file: %v", err)
	}

	if trades == nil {
		t.Error("LoadAllTrades should return empty slice, not nil")
	}

	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
	}
}

func TestDeleteInProgressTrade_RemovesFile(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Create in-progress trade
	trade := &models.Trade{ID: "test", Ticker: "SPY"}
	SaveInProgressTrade(trade)

	if !fileExists(InProgressFile) {
		t.Fatal("InProgressFile should exist")
	}

	// Delete it
	err := DeleteInProgressTrade()
	if err != nil {
		t.Errorf("DeleteInProgressTrade failed: %v", err)
	}

	if fileExists(InProgressFile) {
		t.Error("InProgressFile should be deleted")
	}
}

func TestDeleteInProgressTrade_NoFile_NoError(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	// Delete non-existent file
	err := DeleteInProgressTrade()
	if err != nil {
		t.Errorf("DeleteInProgressTrade should not error on missing file: %v", err)
	}
}

func TestAtomicWrites_NoPartialData(t *testing.T) {
	cleanup := setupTestDataDir(t)
	defer cleanup()

	trade := &models.Trade{
		ID:     "test-atomic",
		Ticker: "SPY",
		Sector: "Broad Market",
	}

	// Save trade
	err := SaveInProgressTrade(trade)
	if err != nil {
		t.Fatalf("SaveInProgressTrade failed: %v", err)
	}

	// Verify no .tmp file exists (should be renamed)
	tmpFile := InProgressFile + ".tmp"
	if fileExists(tmpFile) {
		t.Error("Temporary file should not exist after save completes")
	}

	// Verify main file exists and is valid
	if !fileExists(InProgressFile) {
		t.Error("InProgressFile should exist")
	}

	// Verify content is complete and valid JSON
	data, _ := os.ReadFile(InProgressFile)
	var loaded models.Trade
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Errorf("File should contain valid JSON: %v", err)
	}
}
