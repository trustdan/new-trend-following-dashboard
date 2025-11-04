package models

import (
	"encoding/json"
	"os"
	"time"
)

// Policy represents the complete policy configuration
type Policy struct {
	PolicyID      string              `json:"policy_id"`
	Version       string              `json:"version"`
	GeneratedAt   time.Time           `json:"generated_at"`
	Sectors       []Sector            `json:"sectors"`
	Strategies    map[string]Strategy `json:"strategies"`
	Checklist     Checklist           `json:"checklist"`
	Defaults      PolicyDefaults      `json:"defaults"`
	Calendar      CalendarConfig      `json:"calendar"`
	FinvizHelpers map[string]string   `json:"finviz_helpers"`
}

// Sector represents a trading sector configuration
type Sector struct {
	Name              string            `json:"name"`
	Priority          int               `json:"priority"`
	Blocked           bool              `json:"blocked"`
	Warning           bool              `json:"warning"`
	HeatCapPercent    float64           `json:"heat_cap_percent"`
	Notes             string            `json:"notes"`
	AllowedStrategies []string          `json:"allowed_strategies"`
	ScreenerURLs      map[string]string `json:"screener_urls"`
}

// Strategy represents a Pine Script trading strategy
type Strategy struct {
	Label              string   `json:"label"`
	OptionsSuitability string   `json:"options_suitability"`
	HoldWeeks          string   `json:"hold_weeks"`
	BestExamples       []string `json:"best_examples"`
	Notes              string   `json:"notes"`
}

// Checklist defines the anti-impulsivity gates
type Checklist struct {
	Required     []string           `json:"required"`
	Optional     []string           `json:"optional"`
	PokerSizing  map[string]float64 `json:"poker_sizing"`
	MinContracts int                `json:"min_contracts"`
}

// PolicyDefaults contains default system values
type PolicyDefaults struct {
	PortfolioHeatCap float64 `json:"portfolio_heat_cap"`
	BucketHeatCap    float64 `json:"bucket_heat_cap"`
	RiskPerTrade     float64 `json:"risk_per_trade"`
	CooldownSeconds  int     `json:"cooldown_seconds"`
	ChartViewParam   string  `json:"chart_view_param"`
}

// CalendarConfig defines calendar view settings
type CalendarConfig struct {
	PastDays   int    `json:"past_days"`
	FutureDays int    `json:"future_days"`
	YAxis      string `json:"y_axis"`
	BarLabel   string `json:"bar_label"`
}

// LoadPolicy loads policy from JSON file
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// SafeModePolicy returns a minimal safe-mode policy
func SafeModePolicy() *Policy {
	return &Policy{
		Version:     "safe-mode",
		GeneratedAt: time.Now(),
		Sectors: []Sector{
			{
				Name:              "Healthcare",
				Priority:          1,
				AllowedStrategies: []string{"Alt10", "Alt43", "Alt46"},
				Blocked:           false,
			},
			{
				Name:              "Technology",
				Priority:          2,
				AllowedStrategies: []string{"Alt26", "Alt22", "Alt10"},
				Blocked:           false,
			},
			{
				Name:    "Utilities",
				Blocked: true,
			},
		},
		Strategies: map[string]Strategy{
			"Alt10": {Label: "Profit Targets (3N/6N/9N)", OptionsSuitability: "excellent"},
			"Alt26": {Label: "Fractional Pyramid", OptionsSuitability: "excellent"},
			"Alt43": {Label: "Volatility-Adaptive", OptionsSuitability: "excellent"},
		},
		Defaults: PolicyDefaults{
			PortfolioHeatCap: 0.04,
			BucketHeatCap:    0.015,
			CooldownSeconds:  120,
		},
	}
}
