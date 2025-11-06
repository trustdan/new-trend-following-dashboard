package vimium

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LinkHintMode manages the link hint overlay system
type LinkHintMode struct {
	active       bool
	hints        []Hint
	currentInput string
	overlay      *HintsOverlay
	onComplete   func()
}

// Hint represents a clickable element with its hint label
type Hint struct {
	Label    string
	Element  fyne.Tappable
	Position fyne.Position
	Size     fyne.Size
}

// NewLinkHintMode creates a new link hint mode manager
func NewLinkHintMode() *LinkHintMode {
	return &LinkHintMode{
		active:       false,
		hints:        []Hint{},
		currentInput: "",
		overlay:      NewHintsOverlay(),
	}
}

// IsActive returns whether link hint mode is currently active
func (lhm *LinkHintMode) IsActive() bool {
	return lhm.active
}

// Activate activates link hint mode and scans for clickable elements
func (lhm *LinkHintMode) Activate(content fyne.CanvasObject, window fyne.Window) {
	lhm.active = true
	lhm.currentInput = ""
	lhm.hints = lhm.scanClickableElements(content)
	lhm.overlay.SetHints(lhm.hints)
	lhm.overlay.Show()
}

// Deactivate deactivates link hint mode
func (lhm *LinkHintMode) Deactivate() {
	lhm.active = false
	lhm.currentInput = ""
	lhm.hints = []Hint{}
	lhm.overlay.Hide()
}

// HandleKeyPress processes a key press in link hint mode
// Returns true if the key was handled, false otherwise
func (lhm *LinkHintMode) HandleKeyPress(key *fyne.KeyEvent) bool {
	if !lhm.active {
		return false
	}

	// Escape cancels link hint mode
	if key.Name == fyne.KeyEscape {
		lhm.Deactivate()
		return true
	}

	// Backspace removes last character
	if key.Name == fyne.KeyBackspace {
		if len(lhm.currentInput) > 0 {
			lhm.currentInput = lhm.currentInput[:len(lhm.currentInput)-1]
			lhm.updateVisibleHints()
			return true
		}
		lhm.Deactivate()
		return true
	}

	// Only handle letter keys
	keyChar := lhm.keyToChar(key.Name)
	if keyChar == "" {
		return false
	}

	lhm.currentInput += keyChar
	lhm.updateVisibleHints()

	// Check if we have a complete match
	for _, hint := range lhm.hints {
		if strings.ToLower(hint.Label) == strings.ToLower(lhm.currentInput) {
			// Execute the click action
			if hint.Element != nil {
				hint.Element.Tapped(&fyne.PointEvent{
					Position: fyne.NewPos(hint.Position.X+hint.Size.Width/2,
						hint.Position.Y+hint.Size.Height/2),
				})
			}
			lhm.Deactivate()
			if lhm.onComplete != nil {
				lhm.onComplete()
			}
			return true
		}
	}

	return true
}

// SetCompleteCallback sets the callback to execute when a hint is selected
func (lhm *LinkHintMode) SetCompleteCallback(callback func()) {
	lhm.onComplete = callback
}

// GetOverlay returns the hints overlay widget
func (lhm *LinkHintMode) GetOverlay() *HintsOverlay {
	return lhm.overlay
}

// scanClickableElements recursively scans the UI for clickable elements
func (lhm *LinkHintMode) scanClickableElements(obj fyne.CanvasObject) []Hint {
	hints := []Hint{}
	counter := 0

	var scan func(fyne.CanvasObject, fyne.Position)
	scan = func(obj fyne.CanvasObject, offset fyne.Position) {
		if obj == nil {
			return
		}

		pos := obj.Position().Add(offset)

		// Check if object is tappable (button, widget, etc.)
		if tappable, ok := obj.(fyne.Tappable); ok {
			// Generate hint label
			label := generateHintLabel(counter)
			hints = append(hints, Hint{
				Label:    label,
				Element:  tappable,
				Position: pos,
				Size:     obj.Size(),
			})
			counter++
		}

		// Recursively scan containers
		if container, ok := obj.(*fyne.Container); ok {
			for _, child := range container.Objects {
				scan(child, pos)
			}
		}
	}

	scan(obj, fyne.NewPos(0, 0))
	return hints
}

// updateVisibleHints filters hints based on current input
func (lhm *LinkHintMode) updateVisibleHints() {
	if lhm.currentInput == "" {
		lhm.overlay.SetHints(lhm.hints)
		return
	}

	// Filter hints that start with current input
	filtered := []Hint{}
	for _, hint := range lhm.hints {
		if strings.HasPrefix(strings.ToLower(hint.Label), strings.ToLower(lhm.currentInput)) {
			filtered = append(filtered, hint)
		}
	}
	lhm.overlay.SetHints(filtered)
}

// keyToChar converts a key name to a character
func (lhm *LinkHintMode) keyToChar(keyName fyne.KeyName) string {
	// Map key names to characters
	keyMap := map[fyne.KeyName]string{
		fyne.KeyA: "a", fyne.KeyB: "b", fyne.KeyC: "c", fyne.KeyD: "d",
		fyne.KeyE: "e", fyne.KeyF: "f", fyne.KeyG: "g", fyne.KeyH: "h",
		fyne.KeyI: "i", fyne.KeyJ: "j", fyne.KeyK: "k", fyne.KeyL: "l",
		fyne.KeyM: "m", fyne.KeyN: "n", fyne.KeyO: "o", fyne.KeyP: "p",
		fyne.KeyQ: "q", fyne.KeyR: "r", fyne.KeyS: "s", fyne.KeyT: "t",
		fyne.KeyU: "u", fyne.KeyV: "v", fyne.KeyW: "w", fyne.KeyX: "x",
		fyne.KeyY: "y", fyne.KeyZ: "z",
	}

	if char, exists := keyMap[keyName]; exists {
		return char
	}
	return ""
}

// generateHintLabel generates a hint label (a, b, c, ... z, aa, ab, ...)
func generateHintLabel(index int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	if index < 26 {
		return string(chars[index])
	}

	// For indices >= 26, use two characters: aa, ab, ac, ..., ba, bb, ...
	first := index / 26
	second := index % 26
	return string(chars[first-1]) + string(chars[second])
}

// HintsOverlay displays hint labels over clickable elements
type HintsOverlay struct {
	widget.BaseWidget
	hints   []Hint
	visible bool
}

// NewHintsOverlay creates a new hints overlay
func NewHintsOverlay() *HintsOverlay {
	overlay := &HintsOverlay{
		hints:   []Hint{},
		visible: false,
	}
	overlay.ExtendBaseWidget(overlay)
	return overlay
}

// Show shows the hints overlay
func (ho *HintsOverlay) Show() {
	ho.visible = true
	ho.Refresh()
}

// Hide hides the hints overlay
func (ho *HintsOverlay) Hide() {
	ho.visible = false
	ho.Refresh()
}

// SetHints sets the hints to display
func (ho *HintsOverlay) SetHints(hints []Hint) {
	ho.hints = hints
	ho.Refresh()
}

// CreateRenderer creates the renderer for the hints overlay
func (ho *HintsOverlay) CreateRenderer() fyne.WidgetRenderer {
	if !ho.visible || len(ho.hints) == 0 {
		return &hintsRenderer{
			overlay: ho,
			objects: []fyne.CanvasObject{},
		}
	}

	objects := []fyne.CanvasObject{}

	for _, hint := range ho.hints {
		// Create hint label with yellow background
		label := widget.NewLabel(strings.ToUpper(hint.Label))
		label.TextStyle = fyne.TextStyle{Bold: true}

		// Create background rectangle
		bg := canvas.NewRectangle(theme.WarningColor())
		bg.CornerRadius = 4

		// Position at top-left of element
		labelSize := label.MinSize()
		bg.Resize(fyne.NewSize(labelSize.Width+8, labelSize.Height+4))
		bg.Move(hint.Position)

		label.Move(hint.Position.Add(fyne.NewPos(4, 2)))

		objects = append(objects, bg, label)
	}

	return &hintsRenderer{
		overlay: ho,
		objects: objects,
	}
}

// hintsRenderer renders the hints overlay
type hintsRenderer struct {
	overlay *HintsOverlay
	objects []fyne.CanvasObject
}

func (r *hintsRenderer) Layout(size fyne.Size) {
	// Hints are positioned absolutely, no layout needed
}

func (r *hintsRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *hintsRenderer) Refresh() {
	r.objects = r.overlay.CreateRenderer().(*hintsRenderer).objects
	canvas.Refresh(r.overlay)
}

func (r *hintsRenderer) Objects() []fyne.CanvasObject {
	if !r.overlay.visible {
		return []fyne.CanvasObject{}
	}
	return r.objects
}

func (r *hintsRenderer) Destroy() {}
