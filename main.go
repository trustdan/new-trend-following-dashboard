package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"tf-engine/internal/appcore"
	"tf-engine/internal/config"
	"tf-engine/internal/logging"
	"tf-engine/internal/models"
	"tf-engine/internal/storage"
	"tf-engine/internal/ui"
)

const (
	AppName    = "TF-Engine 2.0"
	AppID      = "com.tfsystems.tfengine"
	AppVersion = "1.0.0"
)

func main() {
	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			if logging.ErrorLogger != nil {
				logging.LogPanic(r)
				logging.ErrorLogger.Printf("Application crashed: %v", r)
			}
			fmt.Fprintf(os.Stderr, "FATAL: Application crashed: %v\n", r)
			os.Exit(1)
		}
	}()

	// Initialize logging FIRST (so we can see what's happening)
	if err := logging.InitializeLogging(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}
	defer logging.CloseLogging()

	// Log startup info
	logging.LogStartup()
	logging.InfoLogger.Printf("Starting %s version %s", AppName, AppVersion)

	// Create required directories
	if err := createRequiredDirectories(); err != nil {
		logging.ErrorLogger.Printf("Failed to create required directories: %v", err)
		// Continue anyway - app can still run
	}

	// Clean up old logs (ignore errors)
	logging.CleanupOldLogs()

	// Initialize application state
	logging.InfoLogger.Println("Initializing application state...")
	state := appcore.NewAppState()

	// Load policy file
	logging.InfoLogger.Println("Loading policy configuration...")
	policyPath := findPolicyFile()
	if err := state.LoadPolicy(policyPath); err != nil {
		logging.ErrorLogger.Printf("Failed to load policy: %v", err)
		logging.ErrorLogger.Println("Activating safe mode with minimal policy")
		state.UseSafeMode()
	} else {
		logging.InfoLogger.Printf("Policy loaded successfully from %s", policyPath)
	}

	// Load feature flags
	logging.InfoLogger.Println("Loading feature flags...")
	if _, err := config.LoadFeatureFlags("feature.flags.json"); err != nil {
		logging.ErrorLogger.Printf("Failed to load feature flags: %v", err)
		logging.InfoLogger.Println("Continuing with default feature flags (all Phase 2 features OFF)")
	} else {
		logging.InfoLogger.Println("Feature flags loaded successfully")
	}

	// Load existing trades
	logging.InfoLogger.Println("Loading existing trades...")
	trades, err := storage.LoadAllTrades()
	if err != nil {
		logging.ErrorLogger.Printf("Failed to load trades: %v", err)
		// Continue with empty trade list
		state.AllTrades = []models.Trade{}
	} else {
		state.AllTrades = trades
		logging.InfoLogger.Printf("Loaded %d existing trades", len(trades))
	}

	// Check for in-progress trade
	inProgressTrade, err := storage.LoadInProgressTrade()
	if err == nil && inProgressTrade != nil {
		logging.InfoLogger.Printf("Found in-progress trade: %s", inProgressTrade.Ticker)
		state.CurrentTrade = inProgressTrade
	}

	// Create Fyne application
	logging.InfoLogger.Println("Creating Fyne application...")
	fyneApp := app.NewWithID(AppID)
	fyneApp.Settings().SetTheme(ui.NewTFEngineTheme())

	// Create main window
	logging.InfoLogger.Println("Creating main window...")
	window := fyneApp.NewWindow(AppName)
	window.Resize(fyne.NewSize(1024, 768))
	window.CenterOnScreen()

	// Create navigator
	logging.InfoLogger.Println("Initializing navigator...")
	navigator := ui.NewNavigator(state, window)

	// Show welcome screen on first launch (if feature enabled)
	if shouldShowWelcome() {
		logging.InfoLogger.Println("Showing welcome screen")
		showWelcomeScreen(window, navigator)
	}

	// Start at dashboard
	logging.InfoLogger.Println("Navigating to dashboard...")
	navigator.NavigateToDashboard()

	// Show window and run
	logging.InfoLogger.Println("Application initialized successfully")
	logging.InfoLogger.Println("Showing main window...")
	window.ShowAndRun()

	// Cleanup on exit
	logging.InfoLogger.Println("Application shutting down...")
}

// createRequiredDirectories creates all directories the app needs
func createRequiredDirectories() error {
	dirs := []string{
		"data",
		"data/ui",
		"data/backups",
		"logs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logging.DebugLogger.Printf("Created directory: %s", dir)
	}

	return nil
}

// findPolicyFile locates the policy file in various possible locations
func findPolicyFile() string {
	// Check multiple locations
	locations := []string{
		"data/policy.v1.json",    // Development location
		"policy.v1.json",         // Installed location (same as exe)
		"../data/policy.v1.json", // One level up
		"./data/policy.v1.json",  // Explicit current dir
	}

	// Also check relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		locations = append(locations,
			filepath.Join(exeDir, "policy.v1.json"),
			filepath.Join(exeDir, "data", "policy.v1.json"),
		)
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			logging.DebugLogger.Printf("Found policy file at: %s", loc)
			return loc
		}
	}

	// Default to first location
	logging.DebugLogger.Printf("Policy file not found, using default path: data/policy.v1.json")
	return "data/policy.v1.json"
}

// shouldShowWelcome checks if welcome screen should be shown
func shouldShowWelcome() bool {
	// Check for marker file
	welcomePath := filepath.Join("data", "ui", ".welcome_shown")
	if _, err := os.Stat(welcomePath); err == nil {
		return false // Already shown
	}
	return true // First launch
}

// showWelcomeScreen displays the welcome screen dialog
func showWelcomeScreen(window fyne.Window, navigator *ui.Navigator) {
	welcomeContent := `Welcome to TF-Engine 2.0!

This application guides you through a systematic options trading workflow:

1. Sector Selection - Choose market sector
2. Screener Launch - Find trade candidates
3. Ticker Entry - Enter symbol & strategy
4. Checklist - Anti-impulsivity verification
5. Position Sizing - Calculate contracts
6. Heat Check - Verify portfolio limits
7. Trade Entry - Select options structure
8. Calendar View - Visualize all trades

Key Features:
• Policy-driven strategy selection
• 120-second cooldown timer
• Portfolio heat limits (4% max)
• Sector-based strategy filtering

Click OK to begin!`

	dialog.ShowInformation("Welcome to TF-Engine 2.0", welcomeContent, window)

	// Mark welcome as shown
	welcomePath := filepath.Join("data", "ui", ".welcome_shown")
	os.MkdirAll(filepath.Dir(welcomePath), 0755)
	os.WriteFile(welcomePath, []byte("shown"), 0644)
}
