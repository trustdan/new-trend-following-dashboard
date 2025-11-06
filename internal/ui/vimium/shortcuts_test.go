package vimium

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestShortcutHandler_SetEnabled(t *testing.T) {
	handler := NewShortcutHandler()

	if handler.IsEnabled() {
		t.Error("Handler should be disabled by default")
	}

	handler.SetEnabled(true)
	if !handler.IsEnabled() {
		t.Error("Handler should be enabled after SetEnabled(true)")
	}

	handler.SetEnabled(false)
	if handler.IsEnabled() {
		t.Error("Handler should be disabled after SetEnabled(false)")
	}
}

func TestShortcutHandler_HandleKeyboard_WhenDisabled(t *testing.T) {
	handler := NewShortcutHandler()
	handler.SetEnabled(false)

	nextCalled := false
	handler.SetCallbacks(func() { nextCalled = true }, nil, nil, nil, nil)

	key := &fyne.KeyEvent{Name: fyne.KeyL}
	handled := handler.HandleKeyboard(key)

	if handled {
		t.Error("Handler should not handle keys when disabled")
	}
	if nextCalled {
		t.Error("Next callback should not be called when handler is disabled")
	}
}

func TestShortcutHandler_HandleKeyboard_NextScreen(t *testing.T) {
	handler := NewShortcutHandler()
	handler.SetEnabled(true)

	nextCalled := false
	handler.SetCallbacks(func() { nextCalled = true }, nil, nil, nil, nil)

	// Test 'l' key
	key := &fyne.KeyEvent{Name: fyne.KeyL}
	handled := handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle 'l' key")
	}
	if !nextCalled {
		t.Error("Next callback should be called for 'l' key")
	}

	// Test right arrow
	nextCalled = false
	key = &fyne.KeyEvent{Name: fyne.KeyRight}
	handled = handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle right arrow key")
	}
	if !nextCalled {
		t.Error("Next callback should be called for right arrow")
	}
}

func TestShortcutHandler_HandleKeyboard_PrevScreen(t *testing.T) {
	handler := NewShortcutHandler()
	handler.SetEnabled(true)

	prevCalled := false
	handler.SetCallbacks(nil, func() { prevCalled = true }, nil, nil, nil)

	// Test 'h' key
	key := &fyne.KeyEvent{Name: fyne.KeyH}
	handled := handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle 'h' key")
	}
	if !prevCalled {
		t.Error("Prev callback should be called for 'h' key")
	}

	// Test left arrow
	prevCalled = false
	key = &fyne.KeyEvent{Name: fyne.KeyLeft}
	handled = handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle left arrow key")
	}
	if !prevCalled {
		t.Error("Prev callback should be called for left arrow")
	}

	// Test Escape
	prevCalled = false
	key = &fyne.KeyEvent{Name: fyne.KeyEscape}
	handled = handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle Escape key")
	}
	if !prevCalled {
		t.Error("Prev callback should be called for Escape")
	}
}

func TestShortcutHandler_HandleKeyboard_Home(t *testing.T) {
	handler := NewShortcutHandler()
	handler.SetEnabled(true)

	homeCalled := false
	handler.SetCallbacks(nil, nil, func() { homeCalled = true }, nil, nil)

	key := &fyne.KeyEvent{Name: fyne.KeyG}
	handled := handler.HandleKeyboard(key)

	if !handled {
		t.Error("Handler should handle 'g' key")
	}
	if !homeCalled {
		t.Error("Home callback should be called for 'g' key")
	}
}

func TestGetShortcutHelp(t *testing.T) {
	help := GetShortcutHelp()

	if len(help) == 0 {
		t.Error("GetShortcutHelp should return shortcuts")
	}

	expectedKeys := []string{"f", "j/k", "h/←", "l/→", "g", "/ or ?", "Enter", "Esc", "Ctrl+V"}
	for _, key := range expectedKeys {
		if _, exists := help[key]; !exists {
			t.Errorf("GetShortcutHelp should include '%s' key", key)
		}
	}
}
