package widgets

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CooldownTimer is a reusable countdown timer widget
type CooldownTimer struct {
	widget.BaseWidget

	duration        time.Duration
	startTime       time.Time
	ticker          *time.Ticker
	done            chan bool
	onComplete      func()
	stopped         bool
	frozenRemaining time.Duration

	label       *widget.Label
	progressBar *widget.ProgressBar
	container   *fyne.Container
}

// NewCooldownTimer creates a new cooldown timer
func NewCooldownTimer(duration time.Duration, onComplete func()) *CooldownTimer {
	timer := &CooldownTimer{
		duration:    duration,
		startTime:   time.Now(),
		ticker:      time.NewTicker(1 * time.Second),
		done:        make(chan bool),
		onComplete:  onComplete,
		label:       widget.NewLabel(""),
		progressBar: widget.NewProgressBar(),
	}

	timer.ExtendBaseWidget(timer)
	timer.progressBar.Min = 0
	timer.progressBar.Max = 1.0

	// Create container with label and progress bar
	timer.container = container.NewVBox(
		timer.label,
		timer.progressBar,
	)

	// Start the countdown
	timer.start()

	return timer
}

// CreateRenderer returns the widget renderer
func (t *CooldownTimer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.container)
}

// start begins the countdown timer
func (t *CooldownTimer) start() {
	// Update immediately
	t.update()

	// Start update loop
	go func() {
		for {
			select {
			case <-t.done:
				return
			case <-t.ticker.C:
				t.update()
			}
		}
	}()
}

// update refreshes the timer display and checks for completion
func (t *CooldownTimer) update() {
	elapsed := time.Since(t.startTime)
	remaining := t.duration - elapsed

	if remaining <= 0 {
		t.complete()
		return
	}

	// Update label with MM:SS format
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	t.label.SetText(fmt.Sprintf("Cooldown: %d:%02d remaining", minutes, seconds))

	// Update progress bar (inverse - counts down)
	progress := float64(elapsed) / float64(t.duration)
	t.progressBar.SetValue(progress)

	t.Refresh()
}

// complete is called when the timer reaches zero
func (t *CooldownTimer) complete() {
	t.ticker.Stop()
	t.label.SetText("✓ Ready to continue")
	t.progressBar.SetValue(1.0)

	if t.onComplete != nil {
		t.onComplete()
	}

	t.Refresh()
}

// Stop stops the timer and freezes the remaining time
func (t *CooldownTimer) Stop() {
	if !t.stopped {
		// Calculate remaining before setting stopped flag
		elapsed := time.Since(t.startTime)
		remaining := t.duration - elapsed
		if remaining < 0 {
			remaining = 0
		}

		t.stopped = true
		t.frozenRemaining = remaining

		select {
		case t.done <- true:
		default:
		}
		t.ticker.Stop()
	}
}

// GetRemaining returns the remaining duration
func (t *CooldownTimer) GetRemaining() time.Duration {
	if t.stopped {
		return t.frozenRemaining
	}

	elapsed := time.Since(t.startTime)
	remaining := t.duration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsComplete returns true if the timer has finished
func (t *CooldownTimer) IsComplete() bool {
	return t.GetRemaining() <= 0
}

// Reset resets the timer to the original duration
func (t *CooldownTimer) Reset() {
	t.Stop()
	t.stopped = false
	t.frozenRemaining = 0
	t.startTime = time.Now()
	t.ticker = time.NewTicker(1 * time.Second)
	t.done = make(chan bool)
	t.start()
}

// NewCooldownTimerFromTime creates a timer that started at a specific time
// Used for resuming cooldowns after app restart
func NewCooldownTimerFromTime(duration time.Duration, startTime time.Time, onComplete func()) *CooldownTimer {
	timer := &CooldownTimer{
		duration:    duration,
		startTime:   startTime,
		ticker:      time.NewTicker(1 * time.Second),
		done:        make(chan bool),
		onComplete:  onComplete,
		label:       widget.NewLabel(""),
		progressBar: widget.NewProgressBar(),
	}

	timer.ExtendBaseWidget(timer)
	timer.progressBar.Min = 0
	timer.progressBar.Max = 1.0

	// Create container
	timer.container = container.NewVBox(
		timer.label,
		timer.progressBar,
	)

	// Check if already complete
	if time.Since(startTime) >= duration {
		timer.label.SetText("✓ Ready to continue")
		timer.progressBar.SetValue(1.0)
		if onComplete != nil {
			onComplete()
		}
	} else {
		// Start countdown from where it left off
		timer.start()
	}

	return timer
}
