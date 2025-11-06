package vimium

import (
	"fyne.io/fyne/v2"
)

// ShortcutHandler manages keyboard shortcuts for Vimium mode
type ShortcutHandler struct {
	enabled      bool
	onNext       func() // Next screen (l or Right)
	onPrev       func() // Previous screen (h or Left)
	onHome       func() // Go to dashboard (g)
	onHelp       func() // Show help (/ or ?)
	onLinkHints  func() // Activate link hints (f)
}

// NewShortcutHandler creates a new keyboard shortcut handler
func NewShortcutHandler() *ShortcutHandler {
	return &ShortcutHandler{
		enabled: false,
	}
}

// SetEnabled enables or disables Vimium mode
func (sh *ShortcutHandler) SetEnabled(enabled bool) {
	sh.enabled = enabled
}

// IsEnabled returns whether Vimium mode is enabled
func (sh *ShortcutHandler) IsEnabled() bool {
	return sh.enabled
}

// SetCallbacks sets the callback functions for shortcuts
func (sh *ShortcutHandler) SetCallbacks(next, prev, home, help, linkHints func()) {
	sh.onNext = next
	sh.onPrev = prev
	sh.onHome = home
	sh.onHelp = help
	sh.onLinkHints = linkHints
}

// HandleKeyboard handles keyboard events when Vimium mode is enabled
func (sh *ShortcutHandler) HandleKeyboard(key *fyne.KeyEvent) bool {
	if !sh.enabled {
		return false
	}

	switch key.Name {
	case fyne.KeyJ: // Down/Next in lists
		return true
	case fyne.KeyK: // Up/Previous in lists
		return true
	case fyne.KeyH: // Previous screen
		if sh.onPrev != nil {
			sh.onPrev()
			return true
		}
	case fyne.KeyL: // Next screen
		if sh.onNext != nil {
			sh.onNext()
			return true
		}
	case fyne.KeyLeft: // Previous screen (arrow key)
		if sh.onPrev != nil {
			sh.onPrev()
			return true
		}
	case fyne.KeyRight: // Next screen (arrow key)
		if sh.onNext != nil {
			sh.onNext()
			return true
		}
	case fyne.KeySlash: // Show help with '/' or '?'
		if sh.onHelp != nil {
			sh.onHelp()
			return true
		}
	case fyne.KeyF: // Activate link hints mode
		if sh.onLinkHints != nil {
			sh.onLinkHints()
			return true
		}
	case fyne.KeyG: // 'g' for "go to" commands
		// For now, just go to dashboard
		if sh.onHome != nil {
			sh.onHome()
			return true
		}
	case fyne.KeyEscape: // Cancel/Back
		if sh.onPrev != nil {
			sh.onPrev()
			return true
		}
	}

	return false
}

// GetShortcutHelp returns a map of keyboard shortcuts and their descriptions
func GetShortcutHelp() map[string]string {
	return map[string]string{
		"f":          "Link hints - show clickable elements",
		"j/k":        "Navigate up/down in lists",
		"h/←":        "Previous screen",
		"l/→":        "Next screen",
		"g":          "Go to dashboard",
		"/ or ?":     "Show help",
		"Enter":      "Select/Continue",
		"Esc":        "Cancel/Back",
		"Ctrl+V":     "Toggle Vimium mode",
	}
}
