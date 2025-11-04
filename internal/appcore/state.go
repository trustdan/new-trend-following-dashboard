package appcore

import (
	"tf-engine/internal/config"
	"tf-engine/internal/models"
	"time"
)

// AppState holds the global application state
type AppState struct {
	Policy         *models.Policy
	FeatureFlags   *config.FeatureFlags
	Settings       *models.Settings
	CurrentTrade   *models.Trade
	CurrentScreen  string
	AllTrades      []models.Trade
	CooldownActive bool
	CooldownStart  *time.Time
	SafeModeActive bool
}

// NewAppState creates a new application state
func NewAppState() *AppState {
	return &AppState{
		Settings:      models.DefaultSettings(),
		AllTrades:     []models.Trade{},
		CurrentScreen: "dashboard",
	}
}

// LoadPolicy loads the policy file from disk
func (s *AppState) LoadPolicy(path string) error {
	policy, err := models.LoadPolicy(path)
	if err != nil {
		return err
	}
	s.Policy = policy
	return nil
}

// UseSafeMode activates safe mode with minimal policy
func (s *AppState) UseSafeMode() {
	s.Policy = models.SafeModePolicy()
	s.SafeModeActive = true
}

// StartCooldown begins the 120-second anti-impulsivity timer
func (s *AppState) StartCooldown() {
	now := time.Now()
	s.CooldownStart = &now
	s.CooldownActive = true

	// Also set cooldown start time in current trade for persistence
	if s.CurrentTrade != nil {
		s.CurrentTrade.CooldownStartTime = now
	}
}

// IsCooldownComplete checks if cooldown has expired
func (s *AppState) IsCooldownComplete() bool {
	if s.CooldownStart == nil {
		return false
	}
	elapsed := time.Since(*s.CooldownStart)
	return elapsed >= 120*time.Second
}

// GetCooldownRemaining returns seconds remaining in cooldown
func (s *AppState) GetCooldownRemaining() int {
	if s.CooldownStart == nil {
		return 0
	}
	elapsed := time.Since(*s.CooldownStart)
	remaining := 120*time.Second - elapsed
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}
