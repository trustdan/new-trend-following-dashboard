package vimium

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/config"
)

// VimiumManager manages Vimium mode state and UI
type VimiumManager struct {
	handler         *ShortcutHandler
	overlay         *ShortcutOverlay
	linkHintMode    *LinkHintMode
	toggleButton    *widget.Button
	featureFlags    *config.FeatureFlags
	window          fyne.Window
	currentContent  fyne.CanvasObject // Track current content for link hints
	onRefresh       func()            // Callback to refresh the entire UI
	scrollContainer *container.Scroll // Track scroll container for keyboard scrolling
}

// NewVimiumManager creates a new Vimium mode manager
func NewVimiumManager(featureFlags *config.FeatureFlags, window fyne.Window) *VimiumManager {
	vm := &VimiumManager{
		handler:      NewShortcutHandler(),
		overlay:      NewShortcutOverlay(),
		linkHintMode: NewLinkHintMode(),
		featureFlags: featureFlags,
		window:       window,
	}

	vm.toggleButton = widget.NewButton("⌨️ Vim Mode", func() {
		vm.Toggle()
	})

	// Check if feature is enabled
	if featureFlags != nil && !featureFlags.IsEnabled("vimium_mode") {
		vm.toggleButton.Disable()
	}

	return vm
}

// IsEnabled returns whether Vimium mode is currently enabled
func (vm *VimiumManager) IsEnabled() bool {
	return vm.handler.IsEnabled()
}

// Toggle toggles Vimium mode on/off
func (vm *VimiumManager) Toggle() {
	// Check feature flag
	if vm.featureFlags != nil && !vm.featureFlags.IsEnabled("vimium_mode") {
		return
	}

	enabled := !vm.handler.IsEnabled()
	vm.handler.SetEnabled(enabled)

	if enabled {
		vm.overlay.Show()
		vm.toggleButton.SetText("⌨️ Vim Mode ON")
		if vm.window != nil {
			vm.window.RequestFocus()
		}

		// Auto-hide overlay after 3 seconds so it doesn't cover UI
		go func() {
			time.Sleep(3 * time.Second)
			if vm.handler.IsEnabled() && vm.overlay.IsVisible() {
				vm.overlay.Hide()
				if vm.window != nil && vm.window.Canvas() != nil {
					vm.window.Canvas().Refresh(vm.overlay)
				} else {
					canvas.Refresh(vm.overlay)
				}
			}
		}()
	} else {
		vm.overlay.Hide()
		vm.toggleButton.SetText("⌨️ Vim Mode")
	}

	// Trigger UI refresh to show/hide overlay and keep content reference fresh
	if vm.onRefresh != nil {
		vm.onRefresh()
	}
}

// SetCallbacks sets the navigation callbacks for keyboard shortcuts
func (vm *VimiumManager) SetCallbacks(next, prev, home, help func()) {
	// Create link hints callback that activates link hint mode
	linkHintsCallback := func() {
		vm.ActivateLinkHints()
	}
	vm.handler.SetCallbacks(next, prev, home, help, linkHintsCallback)

	// Set up scroll callbacks
	vm.handler.SetScrollCallbacks(
		vm.scrollDown,
		vm.scrollUp,
		vm.pageDown,
		vm.pageUp,
	)
}

// scrollDown scrolls down by a small amount (j key)
func (vm *VimiumManager) scrollDown() {
	if vm.scrollContainer != nil {
		offset := vm.scrollContainer.Offset
		offset.Y += 30 // Scroll down 30 pixels
		vm.scrollContainer.Offset = offset
		vm.scrollContainer.Refresh()
	}
}

// scrollUp scrolls up by a small amount (k key)
func (vm *VimiumManager) scrollUp() {
	if vm.scrollContainer != nil {
		offset := vm.scrollContainer.Offset
		offset.Y -= 30 // Scroll up 30 pixels
		if offset.Y < 0 {
			offset.Y = 0
		}
		vm.scrollContainer.Offset = offset
		vm.scrollContainer.Refresh()
	}
}

// pageDown scrolls down by a page (d key)
func (vm *VimiumManager) pageDown() {
	if vm.scrollContainer != nil {
		offset := vm.scrollContainer.Offset
		offset.Y += float32(vm.scrollContainer.Size().Height) // Scroll by viewport height
		vm.scrollContainer.Offset = offset
		vm.scrollContainer.Refresh()
	}
}

// pageUp scrolls up by a page (u key)
func (vm *VimiumManager) pageUp() {
	if vm.scrollContainer != nil {
		offset := vm.scrollContainer.Offset
		offset.Y -= float32(vm.scrollContainer.Size().Height) // Scroll by viewport height
		if offset.Y < 0 {
			offset.Y = 0
		}
		vm.scrollContainer.Offset = offset
		vm.scrollContainer.Refresh()
	}
}

// HandleKeyboard handles keyboard events
func (vm *VimiumManager) HandleKeyboard(key *fyne.KeyEvent) bool {
	// If link hint mode is active, it gets first priority
	if vm.linkHintMode.IsActive() {
		return vm.linkHintMode.HandleKeyPress(key)
	}

	// Otherwise, handle normal shortcuts
	return vm.handler.HandleKeyboard(key)
}

// ActivateLinkHints activates link hint mode
func (vm *VimiumManager) ActivateLinkHints() {
	if !vm.handler.IsEnabled() {
		return // Vimium mode must be enabled
	}

	if vm.currentContent == nil {
		return // No content to scan
	}

	vm.linkHintMode.Activate(vm.currentContent, vm.window)
}

// SetCurrentContent updates the current content for link hints
func (vm *VimiumManager) SetCurrentContent(content fyne.CanvasObject) {
	vm.currentContent = content
}

// GetToggleButton returns the toggle button for the UI
func (vm *VimiumManager) GetToggleButton() *widget.Button {
	return vm.toggleButton
}

// GetOverlay returns the keyboard shortcuts overlay
func (vm *VimiumManager) GetOverlay() *ShortcutOverlay {
	return vm.overlay
}

// WrapContent wraps content with Vimium overlay support
func (vm *VimiumManager) WrapContent(content fyne.CanvasObject) fyne.CanvasObject {
	// Store current content for link hints
	vm.currentContent = content

	// Build layers: content, shortcuts overlay, hints overlay (always included)
	layers := []fyne.CanvasObject{
		content,
		vm.overlay,
		vm.linkHintMode.GetOverlay(),
	}

	return container.NewStack(layers...)
}

// CreateVimiumEnabledCanvas creates a canvas that listens for Vimium shortcuts
// This is a helper for attaching keyboard listeners to windows
func (vm *VimiumManager) CreateVimiumEnabledCanvas(window fyne.Window) {
	window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		vm.HandleKeyboard(key)
	})
}

// AttachKeyboardHandler attaches the keyboard event handler to the window
func (vm *VimiumManager) AttachKeyboardHandler() {
	if vm.window != nil {
		vm.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			vm.HandleKeyboard(key)
		})
	}
}

// SetRefreshCallback sets a callback to refresh the UI when Vim mode is toggled
func (vm *VimiumManager) SetRefreshCallback(callback func()) {
	vm.onRefresh = callback
}

// SetScrollContainer sets the scroll container for keyboard scrolling
func (vm *VimiumManager) SetScrollContainer(scroll *container.Scroll) {
	vm.scrollContainer = scroll
}
