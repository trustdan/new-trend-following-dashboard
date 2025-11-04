package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"

	"tf-engine/internal/appcore"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
)

// MockScreen implements the Screen interface for testing
type MockScreen struct {
	name          string
	isValid       bool
	renderCalled  int
	validateCalled int
}

func (m *MockScreen) Render() fyne.CanvasObject {
	m.renderCalled++
	return nil
}

func (m *MockScreen) Validate() bool {
	m.validateCalled++
	return m.isValid
}

func (m *MockScreen) GetName() string {
	return m.name
}

// MockWindow implements the minimal fyne.Window interface for testing
type MockWindow struct {
	content fyne.CanvasObject
	setContentCalled int
}

func (m *MockWindow) SetContent(content fyne.CanvasObject) {
	m.content = content
	m.setContentCalled++
}

// Stub implementations for other fyne.Window methods
func (m *MockWindow) Canvas() fyne.Canvas { return nil }
func (m *MockWindow) CenterOnScreen() {}
func (m *MockWindow) Clipboard() fyne.Clipboard { return nil }
func (m *MockWindow) Close() {}
func (m *MockWindow) Content() fyne.CanvasObject { return m.content }
func (m *MockWindow) Hide() {}
func (m *MockWindow) Icon() fyne.Resource { return nil }
func (m *MockWindow) MainMenu() *fyne.MainMenu { return nil }
func (m *MockWindow) RequestFocus() {}
func (m *MockWindow) Resize(fyne.Size) {}
func (m *MockWindow) SetCloseIntercept(func()) {}
func (m *MockWindow) SetFixedSize(bool) {}
func (m *MockWindow) SetFullScreen(bool) {}
func (m *MockWindow) SetIcon(fyne.Resource) {}
func (m *MockWindow) SetMainMenu(*fyne.MainMenu) {}
func (m *MockWindow) SetMaster() {}
func (m *MockWindow) SetOnClosed(func()) {}
func (m *MockWindow) SetPadded(bool) {}
func (m *MockWindow) SetTitle(string) {}
func (m *MockWindow) SetOnDropped(func(fyne.Position, []fyne.URI)) {}
func (m *MockWindow) RunWithContext(func()) {}
func (m *MockWindow) Show() {}
func (m *MockWindow) ShowAndRun() {}
func (m *MockWindow) Title() string { return "" }
func (m *MockWindow) FullScreen() bool { return false }
func (m *MockWindow) FixedSize() bool { return false }
func (m *MockWindow) Padded() bool { return false }

func setupTestNavigator(t *testing.T) (*Navigator, *appcore.AppState, *MockWindow) {
	// Setup test environment
	state := appcore.NewAppState()
	state.Policy = &models.Policy{
		Defaults: models.PolicyDefaults{
			CooldownSeconds: 120,
		},
	}
	state.CurrentTrade = &models.Trade{
		ID:        "test-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockWindow := &MockWindow{}

	nav := &Navigator{
		screens:      []Screen{},
		currentIndex: -1,
		history:      []int{},
		state:        state,
		window:       mockWindow,
	}

	// Add mock screens
	for i := 0; i < 3; i++ {
		screen := &MockScreen{
			name:    []string{"screen1", "screen2", "screen3"}[i],
			isValid: true,
		}
		nav.screens = append(nav.screens, screen)
	}

	return nav, state, mockWindow
}

func TestNavigator_Next_ValidData(t *testing.T) {
	nav, state, mockWindow := setupTestNavigator(t)

	// Move to first screen
	err := nav.Next()
	if err != nil {
		t.Errorf("Next() failed: %v", err)
	}

	if nav.currentIndex != 0 {
		t.Errorf("Expected currentIndex 0, got %d", nav.currentIndex)
	}

	if state.CurrentScreen != "screen1" {
		t.Errorf("Expected current screen 'screen1', got '%s'", state.CurrentScreen)
	}

	if mockWindow.setContentCalled != 1 {
		t.Errorf("Expected SetContent called 1 time, got %d", mockWindow.setContentCalled)
	}
}

func TestNavigator_Next_InvalidData_Fails(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Move to first screen
	nav.Next()

	// Make current screen invalid
	nav.screens[0].(*MockScreen).isValid = false

	// Try to move to next screen
	err := nav.Next()
	if err == nil {
		t.Error("Next() should fail with invalid data")
	}

	if nav.currentIndex != 0 {
		t.Errorf("CurrentIndex should remain 0, got %d", nav.currentIndex)
	}
}

func TestNavigator_Back_PreservesData(t *testing.T) {
	nav, state, _ := setupTestNavigator(t)

	// Navigate forward twice
	nav.Next()
	state.CurrentTrade.Sector = "Healthcare"
	nav.Next()
	state.CurrentTrade.Ticker = "UNH"

	// Navigate back
	err := nav.Back()
	if err != nil {
		t.Errorf("Back() failed: %v", err)
	}

	if nav.currentIndex != 0 {
		t.Errorf("Expected currentIndex 0 after Back(), got %d", nav.currentIndex)
	}

	// Verify data preserved
	if state.CurrentTrade.Sector != "Healthcare" {
		t.Error("Data should be preserved after Back()")
	}
	if state.CurrentTrade.Ticker != "UNH" {
		t.Error("Data should be preserved after Back()")
	}
}

func TestNavigator_Back_NoHistory_Fails(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Try to go back without any history
	err := nav.Back()
	if err == nil {
		t.Error("Back() should fail when history is empty")
	}
}

func TestNavigator_HistoryStack(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	if nav.CanGoBack() {
		t.Error("Should not be able to go back initially")
	}

	// Navigate forward
	nav.Next()
	if nav.GetHistoryDepth() != 1 {
		t.Errorf("Expected history depth 1, got %d", nav.GetHistoryDepth())
	}
	if !nav.CanGoBack() {
		t.Error("Should be able to go back after navigating forward")
	}

	// Navigate forward again
	nav.Next()
	if nav.GetHistoryDepth() != 2 {
		t.Errorf("Expected history depth 2, got %d", nav.GetHistoryDepth())
	}

	// Navigate back
	nav.Back()
	if nav.GetHistoryDepth() != 1 {
		t.Errorf("Expected history depth 1 after Back(), got %d", nav.GetHistoryDepth())
	}
}

func TestNavigator_AutoSave_CalledOnNavigation(t *testing.T) {
	nav, state, _ := setupTestNavigator(t)

	// Ensure clean state
	storage.DeleteInProgressTrade()

	// Set up a trade
	state.CurrentTrade = &models.Trade{
		ID:     "test-autosave",
		Sector: "Healthcare",
		Ticker: "UNH",
	}

	// Navigate forward (should trigger auto-save)
	err := nav.Next()
	if err != nil {
		t.Fatalf("Next() failed: %v", err)
	}

	// Verify trade was saved
	savedTrade, err := storage.LoadInProgressTrade()
	if err != nil {
		t.Fatalf("Failed to load saved trade: %v", err)
	}

	if savedTrade == nil {
		t.Fatal("Trade should have been auto-saved")
	}

	if savedTrade.Sector != "Healthcare" {
		t.Errorf("Expected saved sector 'Healthcare', got '%s'", savedTrade.Sector)
	}

	// Cleanup
	storage.DeleteInProgressTrade()
}

func TestNavigator_GetCurrentScreenName(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Dashboard
	name := nav.GetCurrentScreenName()
	if name != "dashboard" {
		t.Errorf("Expected 'dashboard', got '%s'", name)
	}

	// First screen (mock screen returns "screen1")
	nav.currentIndex = 0
	name = nav.GetCurrentScreenName()
	if name != "screen1" {
		t.Errorf("Expected 'screen1', got '%s'", name)
	}

	// Second screen (mock screen returns "screen2")
	nav.currentIndex = 1
	name = nav.GetCurrentScreenName()
	if name != "screen2" {
		t.Errorf("Expected 'screen2', got '%s'", name)
	}
}

func TestNavigator_NavigateToScreen(t *testing.T) {
	nav, _, mockWindow := setupTestNavigator(t)

	// Jump to screen 2
	err := nav.NavigateToScreen(2)
	if err != nil {
		t.Errorf("NavigateToScreen(2) failed: %v", err)
	}

	if nav.currentIndex != 2 {
		t.Errorf("Expected currentIndex 2, got %d", nav.currentIndex)
	}

	if mockWindow.setContentCalled != 1 {
		t.Errorf("Expected SetContent called 1 time, got %d", mockWindow.setContentCalled)
	}

	// History should be updated
	if len(nav.history) != 1 {
		t.Errorf("Expected history length 1, got %d", len(nav.history))
	}
}

func TestNavigator_NavigateToScreen_InvalidIndex(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Try invalid indices
	err := nav.NavigateToScreen(-1)
	if err == nil {
		t.Error("NavigateToScreen(-1) should fail")
	}

	err = nav.NavigateToScreen(100)
	if err == nil {
		t.Error("NavigateToScreen(100) should fail")
	}
}

func TestNavigator_ValidateCurrentScreen(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Dashboard - should always be valid
	if !nav.ValidateCurrentScreen() {
		t.Error("Dashboard should always validate")
	}

	// Move to first screen
	nav.Next()

	// Should be valid
	if !nav.ValidateCurrentScreen() {
		t.Error("Screen should be valid")
	}

	// Make screen invalid
	nav.screens[0].(*MockScreen).isValid = false

	if nav.ValidateCurrentScreen() {
		t.Error("Screen should be invalid")
	}
}

func TestNavigator_ClearHistory(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	// Build up history
	nav.Next()
	nav.Next()

	if nav.GetHistoryDepth() != 2 {
		t.Errorf("Expected history depth 2, got %d", nav.GetHistoryDepth())
	}

	// Clear history
	nav.ClearHistory()

	if nav.GetHistoryDepth() != 0 {
		t.Errorf("Expected history depth 0 after clear, got %d", nav.GetHistoryDepth())
	}

	if nav.CanGoBack() {
		t.Error("Should not be able to go back after clearing history")
	}
}

func TestNavigator_GetCurrentIndex(t *testing.T) {
	nav, _, _ := setupTestNavigator(t)

	if nav.GetCurrentIndex() != -1 {
		t.Errorf("Expected initial index -1, got %d", nav.GetCurrentIndex())
	}

	nav.Next()
	if nav.GetCurrentIndex() != 0 {
		t.Errorf("Expected index 0, got %d", nav.GetCurrentIndex())
	}

	nav.Next()
	if nav.GetCurrentIndex() != 1 {
		t.Errorf("Expected index 1, got %d", nav.GetCurrentIndex())
	}
}

func TestNavigator_AutoSave_NilTrade_NoError(t *testing.T) {
	nav, state, _ := setupTestNavigator(t)

	// Clear current trade
	state.CurrentTrade = nil

	// Auto-save should not error with nil trade
	err := nav.AutoSave()
	if err != nil {
		t.Errorf("AutoSave() should not error with nil trade: %v", err)
	}
}
