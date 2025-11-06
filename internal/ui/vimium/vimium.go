package vimium

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"tf-engine/internal/config"
)

// VimiumManager manages Vimium mode state and UI
type VimiumManager struct {
	handler      *ShortcutHandler
	overlay      *ShortcutOverlay
	linkHintMode *LinkHintMode
	toggleButton *widget.Button
	featureFlags *config.FeatureFlags
	window       fyne.Window
	currentContent fyne.CanvasObject // Track current content for link hints
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
	} else {
		vm.overlay.Hide()
		vm.toggleButton.SetText("⌨️ Vim Mode")
	}
}

// SetCallbacks sets the navigation callbacks for keyboard shortcuts
func (vm *VimiumManager) SetCallbacks(next, prev, home, help func()) {
	// Create link hints callback that activates link hint mode
	linkHintsCallback := func() {
		vm.ActivateLinkHints()
	}
	vm.handler.SetCallbacks(next, prev, home, help, linkHintsCallback)
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

	// Build layers: content, shortcuts overlay, hints overlay
	layers := []fyne.CanvasObject{content}

	if vm.overlay.IsVisible() {
		layers = append(layers, vm.overlay)
	}

	if vm.linkHintMode.IsActive() {
		layers = append(layers, vm.linkHintMode.GetOverlay())
	}

	if len(layers) == 1 {
		return content
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
