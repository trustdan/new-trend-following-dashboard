package analytics

import (
	"testing"
	"time"

	"tf-engine/internal/models"
)

func TestCalculateTradeStats_EmptyTrades(t *testing.T) {
	trades := []models.Trade{}
	stats := CalculateTradeStats(trades)

	if stats.TotalTrades != 0 {
		t.Errorf("Expected 0 total trades, got %d", stats.TotalTrades)
	}
	if stats.WinRate != 0 {
		t.Errorf("Expected 0 win rate, got %.2f", stats.WinRate)
	}
}

func TestCalculateTradeStats_AllWinners(t *testing.T) {
	pnl1 := 100.0
	pnl2 := 150.0
	pnl3 := 200.0

	trades := []models.Trade{
		{ProfitLoss: &pnl1},
		{ProfitLoss: &pnl2},
		{ProfitLoss: &pnl3},
	}

	stats := CalculateTradeStats(trades)

	if stats.TotalTrades != 3 {
		t.Errorf("Expected 3 total trades, got %d", stats.TotalTrades)
	}
	if stats.WinningTrades != 3 {
		t.Errorf("Expected 3 winning trades, got %d", stats.WinningTrades)
	}
	if stats.WinRate != 100.0 {
		t.Errorf("Expected 100%% win rate, got %.2f%%", stats.WinRate)
	}
	if stats.TotalPnL != 450.0 {
		t.Errorf("Expected $450 total P&L, got $%.2f", stats.TotalPnL)
	}
	if stats.AveragePnL != 150.0 {
		t.Errorf("Expected $150 average P&L, got $%.2f", stats.AveragePnL)
	}
	if stats.LargestWin != 200.0 {
		t.Errorf("Expected $200 largest win, got $%.2f", stats.LargestWin)
	}
}

func TestCalculateTradeStats_MixedResults(t *testing.T) {
	pnl1 := 100.0
	pnl2 := -50.0
	pnl3 := 150.0
	pnl4 := -75.0

	trades := []models.Trade{
		{ProfitLoss: &pnl1},
		{ProfitLoss: &pnl2},
		{ProfitLoss: &pnl3},
		{ProfitLoss: &pnl4},
	}

	stats := CalculateTradeStats(trades)

	if stats.TotalTrades != 4 {
		t.Errorf("Expected 4 total trades, got %d", stats.TotalTrades)
	}
	if stats.WinningTrades != 2 {
		t.Errorf("Expected 2 winning trades, got %d", stats.WinningTrades)
	}
	if stats.LosingTrades != 2 {
		t.Errorf("Expected 2 losing trades, got %d", stats.LosingTrades)
	}
	if stats.WinRate != 50.0 {
		t.Errorf("Expected 50%% win rate, got %.2f%%", stats.WinRate)
	}

	expectedTotal := 125.0
	if stats.TotalPnL != expectedTotal {
		t.Errorf("Expected $%.2f total P&L, got $%.2f", expectedTotal, stats.TotalPnL)
	}

	if stats.LargestWin != 150.0 {
		t.Errorf("Expected $150 largest win, got $%.2f", stats.LargestWin)
	}
	if stats.LargestLoss != -75.0 {
		t.Errorf("Expected $-75 largest loss, got $%.2f", stats.LargestLoss)
	}

	expectedAvgWin := 125.0
	if stats.AverageWin != expectedAvgWin {
		t.Errorf("Expected $%.2f average win, got $%.2f", expectedAvgWin, stats.AverageWin)
	}

	expectedAvgLoss := -62.5
	if stats.AverageLoss != expectedAvgLoss {
		t.Errorf("Expected $%.2f average loss, got $%.2f", expectedAvgLoss, stats.AverageLoss)
	}

	// Profit factor = total wins / abs(total losses)
	// Total wins = 250, total losses = -125
	expectedProfitFactor := 2.0
	if stats.ProfitFactor != expectedProfitFactor {
		t.Errorf("Expected %.2f profit factor, got %.2f", expectedProfitFactor, stats.ProfitFactor)
	}
}

func TestCalculateTradeStats_Streaks(t *testing.T) {
	pnl1 := 100.0
	pnl2 := 150.0
	pnl3 := 200.0
	pnl4 := -50.0
	pnl5 := -75.0

	trades := []models.Trade{
		{ProfitLoss: &pnl1}, // Win
		{ProfitLoss: &pnl2}, // Win
		{ProfitLoss: &pnl3}, // Win (streak = 3)
		{ProfitLoss: &pnl4}, // Loss
		{ProfitLoss: &pnl5}, // Loss (streak = 2)
	}

	stats := CalculateTradeStats(trades)

	if stats.LongestWinStreak != 3 {
		t.Errorf("Expected longest win streak of 3, got %d", stats.LongestWinStreak)
	}
	if stats.LongestLossStreak != 2 {
		t.Errorf("Expected longest loss streak of 2, got %d", stats.LongestLossStreak)
	}
}

func TestCalculateSectorStats(t *testing.T) {
	pnl1 := 100.0
	pnl2 := -50.0
	pnl3 := 150.0
	pnl4 := 75.0

	trades := []models.Trade{
		{Sector: "Healthcare", ProfitLoss: &pnl1},
		{Sector: "Healthcare", ProfitLoss: &pnl2},
		{Sector: "Technology", ProfitLoss: &pnl3},
		{Sector: "Technology", ProfitLoss: &pnl4},
	}

	stats := CalculateSectorStats(trades)

	if len(stats) != 2 {
		t.Errorf("Expected 2 sectors, got %d", len(stats))
	}

	// Stats should be sorted by total PnL descending
	// Technology: 225 (150 + 75)
	// Healthcare: 50 (100 - 50)

	if stats[0].Sector != "Technology" {
		t.Errorf("Expected first sector to be Technology, got %s", stats[0].Sector)
	}
	if stats[0].TotalPnL != 225.0 {
		t.Errorf("Expected Technology total P&L of $225, got $%.2f", stats[0].TotalPnL)
	}
	if stats[0].WinRate != 100.0 {
		t.Errorf("Expected Technology win rate of 100%%, got %.2f%%", stats[0].WinRate)
	}

	if stats[1].Sector != "Healthcare" {
		t.Errorf("Expected second sector to be Healthcare, got %s", stats[1].Sector)
	}
	if stats[1].TotalPnL != 50.0 {
		t.Errorf("Expected Healthcare total P&L of $50, got $%.2f", stats[1].TotalPnL)
	}
	if stats[1].WinRate != 50.0 {
		t.Errorf("Expected Healthcare win rate of 50%%, got %.2f%%", stats[1].WinRate)
	}
}

func TestCalculateStrategyStats(t *testing.T) {
	pnl1 := 100.0
	pnl2 := -50.0
	pnl3 := 150.0
	pnl4 := 75.0

	trades := []models.Trade{
		{Strategy: "Alt10", ProfitLoss: &pnl1},
		{Strategy: "Alt10", ProfitLoss: &pnl2},
		{Strategy: "Alt26", ProfitLoss: &pnl3},
		{Strategy: "Alt26", ProfitLoss: &pnl4},
	}

	stats := CalculateStrategyStats(trades)

	if len(stats) != 2 {
		t.Errorf("Expected 2 strategies, got %d", len(stats))
	}

	// Stats should be sorted by total PnL descending
	// Alt26: 225 (150 + 75)
	// Alt10: 50 (100 - 50)

	if stats[0].Strategy != "Alt26" {
		t.Errorf("Expected first strategy to be Alt26, got %s", stats[0].Strategy)
	}
	if stats[0].TotalPnL != 225.0 {
		t.Errorf("Expected Alt26 total P&L of $225, got $%.2f", stats[0].TotalPnL)
	}

	if stats[1].Strategy != "Alt10" {
		t.Errorf("Expected second strategy to be Alt10, got %s", stats[1].Strategy)
	}
	if stats[1].TotalPnL != 50.0 {
		t.Errorf("Expected Alt10 total P&L of $50, got $%.2f", stats[1].TotalPnL)
	}
}

func TestCalculateEquityCurve(t *testing.T) {
	now := time.Now()
	pnl1 := 100.0
	pnl2 := -50.0
	pnl3 := 150.0

	trades := []models.Trade{
		{CreatedAt: now.AddDate(0, 0, -2), UpdatedAt: now.AddDate(0, 0, -2), ProfitLoss: &pnl1},
		{CreatedAt: now.AddDate(0, 0, -1), UpdatedAt: now.AddDate(0, 0, -1), ProfitLoss: &pnl2},
		{CreatedAt: now, UpdatedAt: now, ProfitLoss: &pnl3},
	}

	curve := CalculateEquityCurve(trades)

	// Should have 4 points: initial + 3 trades
	if len(curve) != 4 {
		t.Errorf("Expected 4 points in equity curve, got %d", len(curve))
	}

	// Initial point should be at 0
	if curve[0].Equity != 0.0 {
		t.Errorf("Expected initial equity of $0, got $%.2f", curve[0].Equity)
	}

	// After first trade: +100
	if curve[1].Equity != 100.0 {
		t.Errorf("Expected equity after first trade to be $100, got $%.2f", curve[1].Equity)
	}

	// After second trade: +100 - 50 = 50
	if curve[2].Equity != 50.0 {
		t.Errorf("Expected equity after second trade to be $50, got $%.2f", curve[2].Equity)
	}

	// After third trade: +50 + 150 = 200
	if curve[3].Equity != 200.0 {
		t.Errorf("Expected equity after third trade to be $200, got $%.2f", curve[3].Equity)
	}
}

func TestCalculateEquityCurve_Empty(t *testing.T) {
	trades := []models.Trade{}
	curve := CalculateEquityCurve(trades)

	if len(curve) != 0 {
		t.Errorf("Expected empty equity curve, got %d points", len(curve))
	}
}
