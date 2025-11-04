package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"tf-engine/internal/models"
)

// Storage paths
const (
	TradesFile     = "data/trades.json"
	InProgressFile = "data/trades_in_progress.json"
	BackupDir      = "data/backups/"
)

// TradeStorage provides thread-safe trade persistence
type TradeStorage struct {
	mu sync.RWMutex
}

var globalStorage = &TradeStorage{}

// SaveInProgressTrade saves current trade state atomically
func SaveInProgressTrade(trade *models.Trade) error {
	globalStorage.mu.Lock()
	defer globalStorage.mu.Unlock()

	trade.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(trade, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
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
		return nil, nil // No in-progress trade
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

// SaveCompletedTrade saves trade to history and creates backup
func SaveCompletedTrade(trade *models.Trade) error {
	globalStorage.mu.Lock()
	defer globalStorage.mu.Unlock()

	trade.UpdatedAt = time.Now()

	// Load existing trades
	trades, err := loadAllTradesUnsafe()
	if err != nil {
		trades = []models.Trade{}
	}

	// Append new trade
	trades = append(trades, *trade)

	// Backup existing file before overwriting
	if fileExists(TradesFile) {
		if err := os.MkdirAll(BackupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		backupPath := filepath.Join(BackupDir,
			fmt.Sprintf("trades_%s.json", time.Now().Format("20060102_150405")))

		if err := copyFile(TradesFile, backupPath); err != nil {
			// Log warning but don't fail
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
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

// loadAllTradesUnsafe loads trades without locking (internal use)
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

// SaveAllTrades saves the entire trade history (used for edit/delete operations)
func SaveAllTrades(trades []models.Trade) error {
	globalStorage.mu.Lock()
	defer globalStorage.mu.Unlock()

	// Backup existing file before overwriting
	if fileExists(TradesFile) {
		if err := os.MkdirAll(BackupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		backupPath := filepath.Join(BackupDir,
			fmt.Sprintf("trades_%s.json", time.Now().Format("20060102_150405")))

		if err := copyFile(TradesFile, backupPath); err != nil {
			// Log warning but don't fail
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
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

	return nil
}

// DeleteInProgressTrade removes the in-progress trade file
func DeleteInProgressTrade() error {
	globalStorage.mu.Lock()
	defer globalStorage.mu.Unlock()

	if fileExists(InProgressFile) {
		return os.Remove(InProgressFile)
	}
	return nil
}

// Helper functions

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
