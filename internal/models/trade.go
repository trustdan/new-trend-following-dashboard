package models

import (
	"time"
)

// Trade represents a single options trade
type Trade struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Screen 1: Sector Selection
	Sector string `json:"sector"`

	// Screen 3: Ticker + Strategy
	Ticker            string    `json:"ticker"`
	Strategy          string    `json:"strategy"`
	CooldownStartTime time.Time `json:"cooldown_start_time"`
	CooldownComplete  bool      `json:"cooldown_complete"`

	// Screen 4: Checklist
	ChecklistPassed   bool            `json:"checklist_passed"`
	ChecklistRequired map[string]bool `json:"checklist_required"`
	ChecklistOptional map[string]bool `json:"checklist_optional"`

	// Screen 5: Position Sizing
	Conviction       int     `json:"conviction"` // 5-8 poker-bet sizing
	AccountEquity    float64 `json:"account_equity"`
	RiskPerTrade     float64 `json:"risk_per_trade"`
	SizingMultiplier float64 `json:"sizing_multiplier"` // From poker sizing
	PositionSize     int     `json:"position_size"`
	MaxLoss          float64 `json:"max_loss"`

	// Screen 6: Heat Check
	PortfolioHeat   float64 `json:"portfolio_heat"`
	BucketHeat      float64 `json:"bucket_heat"`
	HeatCheckPassed bool    `json:"heat_check_passed"`

	// Screen 7: Trade Entry
	OptionsStrategy string    `json:"options_strategy"`
	Strike1         float64   `json:"strike1"`
	Strike2         float64   `json:"strike2"`
	Strike3         float64   `json:"strike3"`
	Strike4         float64   `json:"strike4"`
	ExpirationDate  time.Time `json:"expiration_date"`
	Premium         float64   `json:"premium"`

	// Exit Information (filled later)
	ExitDate   *time.Time `json:"exit_date,omitempty"`
	ExitPrice  *float64   `json:"exit_price,omitempty"`
	ProfitLoss *float64   `json:"profit_loss,omitempty"`
	Status     string     `json:"status"` // "active", "closed", "expired"
}

// GetStatus returns the current status of the trade
func (t *Trade) GetStatus() string {
	if t.Status != "" {
		return t.Status
	}

	// Derive status if not explicitly set
	if t.ExitDate != nil {
		return "closed"
	}

	// Check if expired
	if time.Now().After(t.ExpirationDate) {
		return "expired"
	}

	return "active"
}

// GetPnL returns the profit/loss for the trade
func (t *Trade) GetPnL() float64 {
	if t.ProfitLoss != nil {
		return *t.ProfitLoss
	}
	return 0.0
}
