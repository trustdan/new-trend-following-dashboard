package analytics

import (
	"sort"
	"time"

	"tf-engine/internal/models"
)

// TradeStats holds performance statistics for trades
type TradeStats struct {
	TotalTrades     int
	WinningTrades   int
	LosingTrades    int
	WinRate         float64
	TotalPnL        float64
	AveragePnL      float64
	AverageWin      float64
	AverageLoss     float64
	LargestWin      float64
	LargestLoss     float64
	ProfitFactor    float64
	SharpeRatio     float64
	MaxDrawdown     float64
	MaxDrawdownPct  float64
	CurrentStreak   int
	LongestWinStreak  int
	LongestLossStreak int
}

// SectorStats holds performance statistics by sector
type SectorStats struct {
	Sector      string
	TotalTrades int
	WinRate     float64
	TotalPnL    float64
	AveragePnL  float64
}

// StrategyStats holds performance statistics by strategy
type StrategyStats struct {
	Strategy    string
	TotalTrades int
	WinRate     float64
	TotalPnL    float64
	AveragePnL  float64
}

// CalculateTradeStats computes overall performance statistics
func CalculateTradeStats(trades []models.Trade) TradeStats {
	stats := TradeStats{}

	if len(trades) == 0 {
		return stats
	}

	var totalWins, totalLosses float64
	var winCount, lossCount int
	var currentStreak int
	var longestWinStreak, longestLossStreak int
	var lastTradeWin bool

	// Track equity curve for drawdown calculation
	equity := 0.0
	maxEquity := 0.0
	maxDrawdown := 0.0

	for i, trade := range trades {
		// Skip active trades without PnL
		if trade.ProfitLoss == nil {
			continue
		}

		pnl := *trade.ProfitLoss
		stats.TotalTrades++
		stats.TotalPnL += pnl

		// Track wins/losses
		if pnl > 0 {
			stats.WinningTrades++
			totalWins += pnl
			winCount++

			if pnl > stats.LargestWin {
				stats.LargestWin = pnl
			}

			// Track win streak
			if i == 0 || lastTradeWin {
				currentStreak++
			} else {
				currentStreak = 1
			}
			lastTradeWin = true

			if currentStreak > longestWinStreak {
				longestWinStreak = currentStreak
			}
		} else if pnl < 0 {
			stats.LosingTrades++
			totalLosses += pnl
			lossCount++

			if pnl < stats.LargestLoss {
				stats.LargestLoss = pnl
			}

			// Track loss streak
			if i == 0 || !lastTradeWin {
				currentStreak++
			} else {
				currentStreak = 1
			}
			lastTradeWin = false

			if currentStreak > longestLossStreak {
				longestLossStreak = currentStreak
			}
		}

		// Update equity curve
		equity += pnl
		if equity > maxEquity {
			maxEquity = equity
		}

		// Calculate drawdown
		drawdown := maxEquity - equity
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	// Calculate derived statistics
	if stats.TotalTrades > 0 {
		stats.WinRate = float64(stats.WinningTrades) / float64(stats.TotalTrades) * 100
		stats.AveragePnL = stats.TotalPnL / float64(stats.TotalTrades)
	}

	if winCount > 0 {
		stats.AverageWin = totalWins / float64(winCount)
	}

	if lossCount > 0 {
		stats.AverageLoss = totalLosses / float64(lossCount)
	}

	// Profit factor = total wins / abs(total losses)
	if totalLosses < 0 {
		stats.ProfitFactor = totalWins / (-totalLosses)
	}

	// Max drawdown
	stats.MaxDrawdown = maxDrawdown
	if maxEquity > 0 {
		stats.MaxDrawdownPct = (maxDrawdown / maxEquity) * 100
	}

	// Current and longest streaks
	stats.CurrentStreak = currentStreak
	stats.LongestWinStreak = longestWinStreak
	stats.LongestLossStreak = longestLossStreak

	return stats
}

// CalculateSectorStats computes performance statistics by sector
func CalculateSectorStats(trades []models.Trade) []SectorStats {
	sectorMap := make(map[string]*SectorStats)

	for _, trade := range trades {
		if trade.ProfitLoss == nil {
			continue
		}

		pnl := *trade.ProfitLoss

		if _, exists := sectorMap[trade.Sector]; !exists {
			sectorMap[trade.Sector] = &SectorStats{
				Sector: trade.Sector,
			}
		}

		stats := sectorMap[trade.Sector]
		stats.TotalTrades++
		stats.TotalPnL += pnl

		if pnl > 0 {
			// Win
		}
	}

	// Calculate derived statistics and convert to slice
	result := []SectorStats{}
	for _, stats := range sectorMap {
		if stats.TotalTrades > 0 {
			stats.AveragePnL = stats.TotalPnL / float64(stats.TotalTrades)

			// Calculate win rate (need to count wins)
			winCount := 0
			for _, trade := range trades {
				if trade.Sector == stats.Sector && trade.ProfitLoss != nil && *trade.ProfitLoss > 0 {
					winCount++
				}
			}
			stats.WinRate = float64(winCount) / float64(stats.TotalTrades) * 100
		}
		result = append(result, *stats)
	}

	// Sort by total PnL descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalPnL > result[j].TotalPnL
	})

	return result
}

// CalculateStrategyStats computes performance statistics by strategy
func CalculateStrategyStats(trades []models.Trade) []StrategyStats {
	strategyMap := make(map[string]*StrategyStats)

	for _, trade := range trades {
		if trade.ProfitLoss == nil {
			continue
		}

		pnl := *trade.ProfitLoss

		if _, exists := strategyMap[trade.Strategy]; !exists {
			strategyMap[trade.Strategy] = &StrategyStats{
				Strategy: trade.Strategy,
			}
		}

		stats := strategyMap[trade.Strategy]
		stats.TotalTrades++
		stats.TotalPnL += pnl
	}

	// Calculate derived statistics and convert to slice
	result := []StrategyStats{}
	for _, stats := range strategyMap {
		if stats.TotalTrades > 0 {
			stats.AveragePnL = stats.TotalPnL / float64(stats.TotalTrades)

			// Calculate win rate (need to count wins)
			winCount := 0
			for _, trade := range trades {
				if trade.Strategy == stats.Strategy && trade.ProfitLoss != nil && *trade.ProfitLoss > 0 {
					winCount++
				}
			}
			stats.WinRate = float64(winCount) / float64(stats.TotalTrades) * 100
		}
		result = append(result, *stats)
	}

	// Sort by total PnL descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalPnL > result[j].TotalPnL
	})

	return result
}

// EquityCurvePoint represents a point on the equity curve
type EquityCurvePoint struct {
	Date   time.Time
	Equity float64
}

// CalculateEquityCurve computes the equity curve over time
func CalculateEquityCurve(trades []models.Trade) []EquityCurvePoint {
	if len(trades) == 0 {
		return []EquityCurvePoint{}
	}

	// Sort trades by created date
	sortedTrades := make([]models.Trade, len(trades))
	copy(sortedTrades, trades)
	sort.Slice(sortedTrades, func(i, j int) bool {
		return sortedTrades[i].CreatedAt.Before(sortedTrades[j].CreatedAt)
	})

	curve := []EquityCurvePoint{}
	equity := 0.0

	// Add initial point at zero
	if len(sortedTrades) > 0 {
		curve = append(curve, EquityCurvePoint{
			Date:   sortedTrades[0].CreatedAt,
			Equity: 0.0,
		})
	}

	// Add points for each trade with PnL
	for _, trade := range sortedTrades {
		if trade.ProfitLoss == nil {
			continue
		}

		equity += *trade.ProfitLoss
		curve = append(curve, EquityCurvePoint{
			Date:   trade.UpdatedAt, // Use UpdatedAt for when trade was closed
			Equity: equity,
		})
	}

	return curve
}
