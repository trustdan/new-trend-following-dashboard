package generators

import (
	"testing"
	"time"
)

func TestGenerateSampleTrades_CreatesCorrectCount(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(10)

	// Assert
	if len(trades) != 10 {
		t.Errorf("Expected 10 trades, got %d", len(trades))
	}
}

func TestGenerateSampleTrades_AllFieldsPopulated(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(5)

	// Assert
	for i, trade := range trades {
		if trade.ID == "" {
			t.Errorf("Trade %d has empty ID", i)
		}
		if trade.Ticker == "" {
			t.Errorf("Trade %d has empty Ticker", i)
		}
		if trade.Sector == "" {
			t.Errorf("Trade %d has empty Sector", i)
		}
		if trade.Strategy == "" {
			t.Errorf("Trade %d has empty Strategy", i)
		}
		if trade.OptionsStrategy == "" {
			t.Errorf("Trade %d has empty OptionsStrategy", i)
		}
		if trade.MaxLoss == 0 {
			t.Errorf("Trade %d has zero MaxLoss", i)
		}
		if trade.Status == "" {
			t.Errorf("Trade %d has empty Status", i)
		}
	}
}

func TestGenerateSampleTrades_ValidSectors(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Valid sectors
	validSectors := map[string]bool{
		"Healthcare":  true,
		"Technology":  true,
		"Industrials": true,
		"Consumer":    true,
		"Financials":  true,
	}

	// Assert
	for i, trade := range trades {
		if !validSectors[trade.Sector] {
			t.Errorf("Trade %d has invalid sector: %s", i, trade.Sector)
		}
	}
}

func TestGenerateSampleTrades_ValidStrategies(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Valid strategies
	validStrategies := map[string]bool{
		"Alt10": true,
		"Alt26": true,
		"Alt43": true,
		"Alt46": true,
	}

	// Assert
	for i, trade := range trades {
		if !validStrategies[trade.Strategy] {
			t.Errorf("Trade %d has invalid strategy: %s", i, trade.Strategy)
		}
	}
}

func TestGenerateSampleTrades_DateRanges(t *testing.T) {
	// Arrange
	now := time.Now()

	// Act
	trades := GenerateSampleTrades(10)

	// Assert
	for i, trade := range trades {
		// Entry date should be within past 14 days
		daysSinceEntry := int(now.Sub(trade.CreatedAt).Hours() / 24)
		if daysSinceEntry < 0 || daysSinceEntry > 14 {
			t.Errorf("Trade %d entry date %d days ago (expected 0-14)", i, daysSinceEntry)
		}

		// Expiration should be after entry date
		if trade.ExpirationDate.Before(trade.CreatedAt) {
			t.Errorf("Trade %d expiration before entry date", i)
		}

		// Expiration should be 14-84 days from entry
		daysToExpire := int(trade.ExpirationDate.Sub(trade.CreatedAt).Hours() / 24)
		if daysToExpire < 14 || daysToExpire > 84 {
			t.Errorf("Trade %d has %d days to expiration (expected 14-84)", i, daysToExpire)
		}
	}
}

func TestGenerateSampleTrades_RiskRange(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Assert
	for i, trade := range trades {
		if trade.MaxLoss < 200 || trade.MaxLoss > 700 {
			t.Errorf("Trade %d has risk $%.2f (expected $200-$700)", i, trade.MaxLoss)
		}
	}
}

func TestGenerateSampleTrades_PnLRange(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Assert
	for i, trade := range trades {
		if trade.ProfitLoss != nil {
			pnl := *trade.ProfitLoss
			if pnl < -100 || pnl > 300 {
				t.Errorf("Trade %d has P&L $%.2f (expected -$100 to +$300)", i, pnl)
			}
		}
	}
}

func TestGenerateSampleTrades_ConvictionRange(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Assert
	for i, trade := range trades {
		if trade.Conviction < 5 || trade.Conviction > 8 {
			t.Errorf("Trade %d has conviction %d (expected 5-8)", i, trade.Conviction)
		}
	}
}

func TestGenerateSampleTrades_StatusValidity(t *testing.T) {
	// Act
	trades := GenerateSampleTrades(20)

	// Valid statuses
	validStatuses := map[string]bool{
		"active":  true,
		"closed":  true,
		"expired": true,
	}

	// Assert
	for i, trade := range trades {
		if !validStatuses[trade.Status] {
			t.Errorf("Trade %d has invalid status: %s", i, trade.Status)
		}
	}
}

func TestGenerateHeatCheckScenario_CreatesExpectedTrades(t *testing.T) {
	// Act
	trades := GenerateHeatCheckScenario()

	// Assert
	if len(trades) != 5 {
		t.Errorf("Expected 5 trades for heat check scenario, got %d", len(trades))
	}

	// Verify Healthcare sector has 3 trades
	healthcareCount := 0
	for _, trade := range trades {
		if trade.Sector == "Healthcare" {
			healthcareCount++
		}
	}

	if healthcareCount != 3 {
		t.Errorf("Expected 3 Healthcare trades, got %d", healthcareCount)
	}

	// Verify Technology sector has 2 trades
	techCount := 0
	for _, trade := range trades {
		if trade.Sector == "Technology" {
			techCount++
		}
	}

	if techCount != 2 {
		t.Errorf("Expected 2 Technology trades, got %d", techCount)
	}
}

func TestGenerateHeatCheckScenario_CalculatesTotalRisk(t *testing.T) {
	// Act
	trades := GenerateHeatCheckScenario()

	// Calculate total risk
	totalRisk := 0.0
	for _, trade := range trades {
		if trade.Status == "active" {
			totalRisk += trade.MaxLoss
		}
	}

	// Assert - should be $2000 total ($1150 Healthcare + $850 Technology)
	expected := 2000.0
	if totalRisk != expected {
		t.Errorf("Expected total risk $%.2f, got $%.2f", expected, totalRisk)
	}
}

func TestGenerateMixedStatusTrades_CreatesVariedStatuses(t *testing.T) {
	// Act
	trades := GenerateMixedStatusTrades()

	// Assert
	if len(trades) != 4 {
		t.Errorf("Expected 4 mixed status trades, got %d", len(trades))
	}

	// Count statuses
	statusCounts := make(map[string]int)
	for _, trade := range trades {
		statusCounts[trade.Status]++
	}

	// Should have at least 2 different statuses
	if len(statusCounts) < 2 {
		t.Error("Expected multiple different trade statuses")
	}

	// Verify specific expected statuses
	if statusCounts["active"] == 0 {
		t.Error("Expected at least one active trade")
	}
	if statusCounts["closed"] == 0 {
		t.Error("Expected at least one closed trade")
	}
	if statusCounts["expired"] == 0 {
		t.Error("Expected at least one expired trade")
	}
}

func TestGenerateMixedStatusTrades_AllFieldsPopulated(t *testing.T) {
	// Act
	trades := GenerateMixedStatusTrades()

	// Assert
	for i, trade := range trades {
		if trade.ID == "" {
			t.Errorf("Trade %d has empty ID", i)
		}
		if trade.Ticker == "" {
			t.Errorf("Trade %d has empty Ticker", i)
		}
		if trade.Sector == "" {
			t.Errorf("Trade %d has empty Sector", i)
		}
		if trade.Strategy == "" {
			t.Errorf("Trade %d has empty Strategy", i)
		}
		if trade.OptionsStrategy == "" {
			t.Errorf("Trade %d has empty OptionsStrategy", i)
		}
		if trade.MaxLoss == 0 {
			t.Errorf("Trade %d has zero MaxLoss", i)
		}
		if trade.ProfitLoss == nil {
			t.Errorf("Trade %d has nil ProfitLoss", i)
		}
	}
}

func TestGenerateTradeID_CreatesUniqueIDs(t *testing.T) {
	// Act
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateTradeID(i)
		if ids[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		ids[id] = true
	}

	// Assert
	if len(ids) != 100 {
		t.Errorf("Expected 100 unique IDs, got %d", len(ids))
	}
}
