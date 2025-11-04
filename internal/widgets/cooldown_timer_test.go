package widgets

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
)

func init() {
	// Initialize test app for all tests
	test.NewApp()
}

func TestCooldownTimer_StartsAtFullDuration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 120 * time.Second
	called := false
	onComplete := func() { called = true }

	timer := NewCooldownTimer(duration, onComplete)
	defer timer.Stop()

	remaining := timer.GetRemaining()
	// Allow 1 second tolerance for test execution time
	if remaining < 119*time.Second || remaining > 120*time.Second {
		t.Errorf("Expected remaining time around 120s, got %.0fs", remaining.Seconds())
	}

	if timer.IsComplete() {
		t.Error("Timer should not be complete immediately")
	}

	if called {
		t.Error("onComplete should not be called immediately")
	}
}

func TestCooldownTimer_CountsDown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 3 * time.Second
	timer := NewCooldownTimer(duration, nil)
	defer timer.Stop()

	initialRemaining := timer.GetRemaining()

	// Wait 1.5 seconds
	time.Sleep(1500 * time.Millisecond)

	remaining := timer.GetRemaining()

	// Remaining should be less than initial
	if remaining >= initialRemaining {
		t.Error("Timer should count down")
	}

	// Should have about 1.5 seconds remaining (allow 0.5s tolerance)
	if remaining < 1*time.Second || remaining > 2*time.Second {
		t.Errorf("Expected about 1.5s remaining, got %.2fs", remaining.Seconds())
	}
}

func TestCooldownTimer_CallsOnComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 1 * time.Second
	called := false
	onComplete := func() { called = true }

	timer := NewCooldownTimer(duration, onComplete)
	defer timer.Stop()

	// Wait for timer to complete (with buffer)
	time.Sleep(1500 * time.Millisecond)

	if !called {
		t.Error("onComplete should have been called")
	}

	if !timer.IsComplete() {
		t.Error("Timer should be complete")
	}

	remaining := timer.GetRemaining()
	if remaining != 0 {
		t.Errorf("Remaining time should be 0, got %.2fs", remaining.Seconds())
	}
}

func TestCooldownTimer_Stop(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 10 * time.Second
	timer := NewCooldownTimer(duration, nil)

	// Let it run briefly
	time.Sleep(100 * time.Millisecond)

	// Stop it
	timer.Stop()

	initialRemaining := timer.GetRemaining()

	// Wait a bit longer
	time.Sleep(1500 * time.Millisecond)

	// Remaining should be close to initial (allow max 1 second drift due to ticker)
	remaining := timer.GetRemaining()
	drift := initialRemaining - remaining
	if drift > 1*time.Second || drift < -1*time.Second {
		t.Errorf("Timer should not count down significantly after Stop(), drift: %.2fs", drift.Seconds())
	}
}

func TestCooldownTimer_Reset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 5 * time.Second
	timer := NewCooldownTimer(duration, nil)
	defer timer.Stop()

	// Wait a bit
	time.Sleep(1 * time.Second)

	// Reset timer
	timer.Reset()

	// Should be back at full duration
	remaining := timer.GetRemaining()
	if remaining < 4*time.Second || remaining > 5*time.Second {
		t.Errorf("Expected remaining time around 5s after reset, got %.0fs", remaining.Seconds())
	}
}

func TestNewCooldownTimerFromTime_AlreadyComplete(t *testing.T) {
	duration := 2 * time.Second
	startTime := time.Now().Add(-3 * time.Second) // Started 3 seconds ago

	called := false
	onComplete := func() { called = true }

	timer := NewCooldownTimerFromTime(duration, startTime, onComplete)
	defer timer.Stop()

	if !timer.IsComplete() {
		t.Error("Timer should be complete (started 3s ago with 2s duration)")
	}

	if !called {
		t.Error("onComplete should have been called for already-complete timer")
	}
}

func TestNewCooldownTimerFromTime_PartiallyComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	duration := 5 * time.Second
	startTime := time.Now().Add(-3 * time.Second) // Started 3 seconds ago

	timer := NewCooldownTimerFromTime(duration, startTime, nil)
	defer timer.Stop()

	if timer.IsComplete() {
		t.Error("Timer should not be complete yet")
	}

	remaining := timer.GetRemaining()
	// Should have about 2 seconds remaining (allow 1s tolerance)
	if remaining < 1*time.Second || remaining > 3*time.Second {
		t.Errorf("Expected about 2s remaining, got %.2fs", remaining.Seconds())
	}
}

func TestCooldownTimer_MultipleStopCalls_NoError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	timer := NewCooldownTimer(5*time.Second, nil)

	// Multiple stop calls should not cause issues
	timer.Stop()
	timer.Stop()
	timer.Stop()

	// Should not panic or hang
}

func TestCooldownTimer_ZeroDuration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	called := false
	onComplete := func() { called = true }

	timer := NewCooldownTimer(0, onComplete)
	defer timer.Stop()

	// Should complete immediately
	time.Sleep(100 * time.Millisecond)

	if !timer.IsComplete() {
		t.Error("Zero duration timer should complete immediately")
	}

	if !called {
		t.Error("onComplete should be called for zero duration timer")
	}
}

func TestCooldownTimer_NegativeDuration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing-sensitive cooldown timer test in short mode")
	}

	called := false
	onComplete := func() { called = true }

	timer := NewCooldownTimer(-5*time.Second, onComplete)
	defer timer.Stop()

	// Should treat as already complete
	time.Sleep(100 * time.Millisecond)

	if !timer.IsComplete() {
		t.Error("Negative duration timer should complete immediately")
	}

	if !called {
		t.Error("onComplete should be called for negative duration timer")
	}
}

func TestCooldownTimer_GetRemaining_BeforeStart(t *testing.T) {
	duration := 10 * time.Second
	timer := NewCooldownTimer(duration, nil)
	defer timer.Stop()

	// Should return reasonable value immediately
	remaining := timer.GetRemaining()
	if remaining <= 0 || remaining > duration {
		t.Errorf("GetRemaining should return valid value immediately, got %.2fs", remaining.Seconds())
	}
}
