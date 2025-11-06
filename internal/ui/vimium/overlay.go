package vimium

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ShortcutOverlay displays keyboard shortcuts when Vimium mode is active
type ShortcutOverlay struct {
	widget.BaseWidget
	shortcuts map[string]string
	visible   bool
}

// NewShortcutOverlay creates a new keyboard shortcut overlay
func NewShortcutOverlay() *ShortcutOverlay {
	overlay := &ShortcutOverlay{
		shortcuts: GetShortcutHelp(),
		visible:   false,
	}
	overlay.ExtendBaseWidget(overlay)
	return overlay
}

// Show shows the keyboard shortcuts overlay
func (so *ShortcutOverlay) Show() {
	so.visible = true
	so.Refresh()
}

// Hide hides the keyboard shortcuts overlay
func (so *ShortcutOverlay) Hide() {
	so.visible = false
	so.Refresh()
}

// Toggle toggles the visibility of the overlay
func (so *ShortcutOverlay) Toggle() {
	so.visible = !so.visible
	so.Refresh()
}

// IsVisible returns whether the overlay is visible
func (so *ShortcutOverlay) IsVisible() bool {
	return so.visible
}

// CreateRenderer creates the renderer for the overlay
func (so *ShortcutOverlay) CreateRenderer() fyne.WidgetRenderer {
	if !so.visible {
		return widget.NewLabel("").CreateRenderer()
	}

	// Create shortcuts list
	items := []fyne.CanvasObject{}

	title := widget.NewLabelWithStyle("⌨️ Keyboard Shortcuts (Vimium Mode)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})
	items = append(items, title)

	items = append(items, widget.NewSeparator())

	// Add each shortcut
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"f", "Link hints - click anything"},
		{"j/k", "Navigate up/down in lists"},
		{"h / ←", "Previous screen"},
		{"l / →", "Next screen"},
		{"g", "Go to dashboard"},
		{"/ or ?", "Show help"},
		{"Enter", "Select/Continue"},
		{"Esc", "Cancel/Back"},
		{"Ctrl+V", "Toggle Vimium mode"},
	}

	for _, sc := range shortcuts {
		keyLabel := widget.NewLabelWithStyle(sc.key,
			fyne.TextAlignLeading,
			fyne.TextStyle{Monospace: true, Bold: true})
		keyLabel.Resize(fyne.NewSize(100, 0))

		descLabel := widget.NewLabel(sc.desc)

		row := container.NewHBox(
			keyLabel,
			widget.NewLabel(" - "),
			descLabel,
		)
		items = append(items, row)
	}

	content := container.NewVBox(items...)

	// Create background with padding
	bg := canvas.NewRectangle(theme.OverlayBackgroundColor())
	bg.Resize(fyne.NewSize(400, 250))

	return &overlayRenderer{
		overlay:    so,
		background: bg,
		content:    content,
		objects:    []fyne.CanvasObject{bg, content},
	}
}

// overlayRenderer handles rendering the overlay
type overlayRenderer struct {
	overlay    *ShortcutOverlay
	background *canvas.Rectangle
	content    fyne.CanvasObject
	objects    []fyne.CanvasObject
}

func (r *overlayRenderer) Layout(size fyne.Size) {
	// Position in bottom-right corner
	contentSize := fyne.NewSize(400, 250)
	r.background.Resize(contentSize)
	r.content.Resize(contentSize)

	// Position in bottom-right with padding
	padding := float32(20)
	pos := fyne.NewPos(size.Width-contentSize.Width-padding, size.Height-contentSize.Height-padding)
	r.background.Move(pos)
	r.content.Move(pos)
}

func (r *overlayRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 250)
}

func (r *overlayRenderer) Refresh() {
	r.Layout(r.overlay.Size())
	canvas.Refresh(r.overlay)
}

func (r *overlayRenderer) Objects() []fyne.CanvasObject {
	if !r.overlay.visible {
		return []fyne.CanvasObject{}
	}
	return r.objects
}

func (r *overlayRenderer) Destroy() {}
