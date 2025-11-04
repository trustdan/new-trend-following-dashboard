package generators

import (
	"fmt"
	"math/rand"
	"time"

	"tf-engine/internal/models"
)

// GenerateSampleTrades creates realistic sample trades for calendar view testing
func GenerateSampleTrades(count int) []models.Trade {
	sectors := []string{"Healthcare", "Technology", "Industrials", "Consumer", "Financials"}
	tickers := map[string][]string{
		"Healthcare":  {"UNH", "JNJ", "PFE", "ABBV", "TMO"},
		"Technology":  {"MSFT", "AAPL", "NVDA", "AMD", "CRM"},
		"Industrials": {"CAT", "BA", "UNP", "GE", "HON"},
		"Consumer":    {"AMZN", "HD", "NKE", "MCD", "SBUX"},
		"Financials":  {"JPM", "BAC", "GS", "MS", "BLK"},
	}
	strategies := []string{"Alt10", "Alt26", "Alt43", "Alt46"}
	optionsStrategies := []string{
		"Bull call spread",
		"Bear put spread",
		"Bull put credit spread",
		"Iron condor",
		"Long call",
		"Long put",
		"Calendar call spread",
		"Straddle",
	}

	trades := make([]models.Trade, count)
	now := time.Now()

	// Seed random number generator for reproducibility in tests
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		sector := sectors[rand.Intn(len(sectors))]
		ticker := tickers[sector][rand.Intn(len(tickers[sector]))]

		// Random entry date within past 14 days
		entryDate := now.AddDate(0, 0, -rand.Intn(14))

		// Random expiration date 2-12 weeks out from entry
		daysToExpiration := rand.Intn(70) + 14 // 14 to 84 days
		expirationDate := entryDate.AddDate(0, 0, daysToExpiration)

		// Random risk between $200-$700
		risk := float64(rand.Intn(500) + 200)

		// Random P&L between -$100 and +$300
		pnl := float64(rand.Intn(400) - 100)

		// Determine status based on expiration
		status := "active"
		if expirationDate.Before(now) {
			status = "expired"
		} else if rand.Float64() < 0.2 { // 20% chance of being closed early
			status = "closed"
		}

		trades[i] = models.Trade{
			ID:              generateTradeID(i),
			CreatedAt:       entryDate,
			UpdatedAt:       entryDate,
			Sector:          sector,
			Ticker:          ticker,
			Strategy:        strategies[rand.Intn(len(strategies))],
			OptionsStrategy: optionsStrategies[rand.Intn(len(optionsStrategies))],
			ExpirationDate:  expirationDate,
			MaxLoss:         risk,
			ProfitLoss:      &pnl,
			Status:          status,

			// Populate some additional realistic fields
			Conviction:       rand.Intn(4) + 5, // 5-8
			SizingMultiplier: []float64{0.5, 0.75, 1.0, 1.25}[rand.Intn(4)],
			PositionSize:     rand.Intn(10) + 1, // 1-10 contracts
			ChecklistPassed:  true,
			HeatCheckPassed:  true,
			CooldownComplete: true,
		}
	}

	return trades
}

// GenerateHeatCheckScenario creates trades that test heat limit enforcement
func GenerateHeatCheckScenario() []models.Trade {
	return []models.Trade{
		{
			ID:        "heat-1",
			Ticker:    "UNH",
			Sector:    "Healthcare",
			Strategy:  "Alt10",
			MaxLoss:   450,
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, 0, -5),
			UpdatedAt: time.Now().AddDate(0, 0, -5),
		},
		{
			ID:        "heat-2",
			Ticker:    "JNJ",
			Sector:    "Healthcare",
			Strategy:  "Alt46",
			MaxLoss:   300,
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, 0, -3),
			UpdatedAt: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:        "heat-3",
			Ticker:    "ABBV",
			Sector:    "Healthcare",
			Strategy:  "Alt43",
			MaxLoss:   400,
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		}, // Should trigger 1.5% sector limit (total $1150 = 2.3% of $50k)
		{
			ID:        "heat-4",
			Ticker:    "MSFT",
			Sector:    "Technology",
			Strategy:  "Alt10",
			MaxLoss:   500,
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, 0, -7),
			UpdatedAt: time.Now().AddDate(0, 0, -7),
		},
		{
			ID:        "heat-5",
			Ticker:    "AAPL",
			Sector:    "Technology",
			Strategy:  "Alt26",
			MaxLoss:   350,
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, 0, -2),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
		},
	}
}

// GenerateMixedStatusTrades creates trades with various statuses for testing
func GenerateMixedStatusTrades() []models.Trade {
	now := time.Now()

	pnl1 := 150.0
	pnl2 := -75.0
	pnl3 := 0.0
	pnl4 := 220.0

	return []models.Trade{
		{
			ID:              "mixed-1",
			Ticker:          "UNH",
			Sector:          "Healthcare",
			Strategy:        "Alt10",
			OptionsStrategy: "Bull call spread",
			MaxLoss:         350,
			ProfitLoss:      &pnl1,
			Status:          "active",
			ExpirationDate:  now.AddDate(0, 0, 30),
			CreatedAt:       now.AddDate(0, 0, -10),
			UpdatedAt:       now.AddDate(0, 0, -10),
		},
		{
			ID:              "mixed-2",
			Ticker:          "MSFT",
			Sector:          "Technology",
			Strategy:        "Alt26",
			OptionsStrategy: "Bear put spread",
			MaxLoss:         400,
			ProfitLoss:      &pnl2,
			Status:          "closed",
			ExpirationDate:  now.AddDate(0, 0, 15),
			CreatedAt:       now.AddDate(0, 0, -20),
			UpdatedAt:       now.AddDate(0, 0, -5),
		},
		{
			ID:              "mixed-3",
			Ticker:          "CAT",
			Sector:          "Industrials",
			Strategy:        "Alt43",
			OptionsStrategy: "Iron condor",
			MaxLoss:         300,
			ProfitLoss:      &pnl3,
			Status:          "expired",
			ExpirationDate:  now.AddDate(0, 0, -2),
			CreatedAt:       now.AddDate(0, 0, -45),
			UpdatedAt:       now.AddDate(0, 0, -2),
		},
		{
			ID:              "mixed-4",
			Ticker:          "JPM",
			Sector:          "Financials",
			Strategy:        "Alt10",
			OptionsStrategy: "Long call",
			MaxLoss:         250,
			ProfitLoss:      &pnl4,
			Status:          "closed",
			ExpirationDate:  now.AddDate(0, 0, 20),
			CreatedAt:       now.AddDate(0, 0, -15),
			UpdatedAt:       now.AddDate(0, 0, -3),
		},
	}
}

// generateTradeID creates a unique trade ID
func generateTradeID(index int) string {
	return fmt.Sprintf("sample-%d-%d", time.Now().Unix(), index)
}
