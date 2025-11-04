package models

// Settings represents user preferences
type Settings struct {
	ThemeMode        string  `json:"theme_mode"`
	AccountEquity    float64 `json:"account_equity"`
	RiskPerTrade     float64 `json:"risk_per_trade"`
	PortfolioHeatCap float64 `json:"portfolio_heat_cap"`
	BucketHeatCap    float64 `json:"bucket_heat_cap"`
	VimiumEnabled    bool    `json:"vimium_enabled"`
	SampleDataMode   bool    `json:"sample_data_mode"`
}

// DefaultSettings returns default user settings
func DefaultSettings() *Settings {
	return &Settings{
		ThemeMode:        "day",
		AccountEquity:    100000.00,
		RiskPerTrade:     0.0075,
		PortfolioHeatCap: 0.04,
		BucketHeatCap:    0.015,
		VimiumEnabled:    false,
		SampleDataMode:   false,
	}
}
