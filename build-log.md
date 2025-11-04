PS C:\Users\Dan\new-trend-following-dashboard> .\build.bat
========================================
TF-Engine 2.0 Build Script
========================================

[1/5] Syncing policy.v1.json to dist directory...
        1 file(s) copied.
  Γ£ô Policy file synced

[2/5] Running tests...
# tf-engine/internal/ui/screens [tf-engine/internal/ui/screens.test]
internal\ui\screens\screener_launch_test.go:272:3: not enough arguments in call to screen.createScreenerCard
        have (string, string, string, string, string, string)
        want (string, string, string, string, string, string, string)
?       tf-engine       [no test files]
?       tf-engine/internal/appcore      [no test files]
=== RUN   TestLoadFeatureFlags
--- PASS: TestLoadFeatureFlags (0.00s)
=== RUN   TestIsEnabled
=== RUN   TestIsEnabled/Enabled_feature
=== RUN   TestIsEnabled/Disabled_feature
=== RUN   TestIsEnabled/Non-existent_feature
--- PASS: TestIsEnabled (0.00s)
    --- PASS: TestIsEnabled/Enabled_feature (0.00s)
    --- PASS: TestIsEnabled/Disabled_feature (0.00s)
    --- PASS: TestIsEnabled/Non-existent_feature (0.00s)
=== RUN   TestGetFlag
--- PASS: TestGetFlag (0.00s)
=== RUN   TestListEnabledFlags
--- PASS: TestListEnabledFlags (0.00s)
=== RUN   TestListPhase2Flags
--- PASS: TestListPhase2Flags (0.00s)
=== RUN   TestLoadFeatureFlags_FileNotFound
--- PASS: TestLoadFeatureFlags_FileNotFound (0.00s)
=== RUN   TestLoadFeatureFlags_InvalidJSON
--- PASS: TestLoadFeatureFlags_InvalidJSON (0.00s)
PASS
ok      tf-engine/internal/config       (cached)
?       tf-engine/internal/logging      [no test files]
=== RUN   TestLoadPolicyStrategyCount
    policy_integration_test.go:17: policy file not present at data\policy.v1.json: GetFileAttributesEx data\policy.v1.json: The system cannot find the path specified.
--- SKIP: TestLoadPolicyStrategyCount (0.00s)
PASS
ok      tf-engine/internal/models       1.244s
=== RUN   TestSaveInProgressTrade_CreatesFile
--- PASS: TestSaveInProgressTrade_CreatesFile (0.01s)
=== RUN   TestLoadInProgressTrade_RestoresData
--- PASS: TestLoadInProgressTrade_RestoresData (0.01s)
=== RUN   TestLoadInProgressTrade_NoFile_ReturnsNil
--- PASS: TestLoadInProgressTrade_NoFile_ReturnsNil (0.00s)
=== RUN   TestSaveCompletedTrade_AppendsToHistory
--- PASS: TestSaveCompletedTrade_AppendsToHistory (0.02s)
=== RUN   TestSaveCompletedTrade_CreatesBackup
--- PASS: TestSaveCompletedTrade_CreatesBackup (0.01s)
=== RUN   TestSaveCompletedTrade_ClearsInProgress
--- PASS: TestSaveCompletedTrade_ClearsInProgress (0.00s)
=== RUN   TestConcurrentSaves_NoCorruption
--- PASS: TestConcurrentSaves_NoCorruption (0.08s)
=== RUN   TestLoadAllTrades_EmptyFile_ReturnsEmptySlice
--- PASS: TestLoadAllTrades_EmptyFile_ReturnsEmptySlice (0.00s)
=== RUN   TestDeleteInProgressTrade_RemovesFile
--- PASS: TestDeleteInProgressTrade_RemovesFile (0.00s)
=== RUN   TestDeleteInProgressTrade_NoFile_NoError
--- PASS: TestDeleteInProgressTrade_NoFile_NoError (0.00s)
=== RUN   TestAtomicWrites_NoPartialData
--- PASS: TestAtomicWrites_NoPartialData (0.03s)
PASS
ok      tf-engine/internal/storage      1.388s
=== RUN   TestGenerateSampleTrades_CreatesCorrectCount
--- PASS: TestGenerateSampleTrades_CreatesCorrectCount (0.00s)
=== RUN   TestGenerateSampleTrades_AllFieldsPopulated
--- PASS: TestGenerateSampleTrades_AllFieldsPopulated (0.00s)
=== RUN   TestGenerateSampleTrades_ValidSectors
--- PASS: TestGenerateSampleTrades_ValidSectors (0.00s)
=== RUN   TestGenerateSampleTrades_ValidStrategies
--- PASS: TestGenerateSampleTrades_ValidStrategies (0.00s)
=== RUN   TestGenerateSampleTrades_DateRanges
--- PASS: TestGenerateSampleTrades_DateRanges (0.00s)
=== RUN   TestGenerateSampleTrades_RiskRange
--- PASS: TestGenerateSampleTrades_RiskRange (0.00s)
=== RUN   TestGenerateSampleTrades_PnLRange
--- PASS: TestGenerateSampleTrades_PnLRange (0.00s)
=== RUN   TestGenerateSampleTrades_ConvictionRange
--- PASS: TestGenerateSampleTrades_ConvictionRange (0.00s)
=== RUN   TestGenerateSampleTrades_StatusValidity
--- PASS: TestGenerateSampleTrades_StatusValidity (0.00s)
=== RUN   TestGenerateHeatCheckScenario_CreatesExpectedTrades
--- PASS: TestGenerateHeatCheckScenario_CreatesExpectedTrades (0.00s)
=== RUN   TestGenerateHeatCheckScenario_CalculatesTotalRisk
--- PASS: TestGenerateHeatCheckScenario_CalculatesTotalRisk (0.00s)
=== RUN   TestGenerateMixedStatusTrades_CreatesVariedStatuses
--- PASS: TestGenerateMixedStatusTrades_CreatesVariedStatuses (0.00s)
=== RUN   TestGenerateMixedStatusTrades_AllFieldsPopulated
--- PASS: TestGenerateMixedStatusTrades_AllFieldsPopulated (0.00s)
=== RUN   TestGenerateTradeID_CreatesUniqueIDs
--- PASS: TestGenerateTradeID_CreatesUniqueIDs (0.00s)
PASS
ok      tf-engine/internal/testing/generators   1.191s
=== RUN   TestNavigator_Next_ValidData
--- PASS: TestNavigator_Next_ValidData (0.00s)
=== RUN   TestNavigator_Next_InvalidData_Fails
--- PASS: TestNavigator_Next_InvalidData_Fails (0.00s)
=== RUN   TestNavigator_Back_PreservesData
--- PASS: TestNavigator_Back_PreservesData (0.00s)
=== RUN   TestNavigator_Back_NoHistory_Fails
--- PASS: TestNavigator_Back_NoHistory_Fails (0.00s)
=== RUN   TestNavigator_HistoryStack
--- PASS: TestNavigator_HistoryStack (0.00s)
=== RUN   TestNavigator_AutoSave_CalledOnNavigation
--- PASS: TestNavigator_AutoSave_CalledOnNavigation (0.01s)
=== RUN   TestNavigator_GetCurrentScreenName
--- PASS: TestNavigator_GetCurrentScreenName (0.00s)
=== RUN   TestNavigator_NavigateToScreen
--- PASS: TestNavigator_NavigateToScreen (0.00s)
=== RUN   TestNavigator_NavigateToScreen_InvalidIndex
--- PASS: TestNavigator_NavigateToScreen_InvalidIndex (0.00s)
=== RUN   TestNavigator_ValidateCurrentScreen
--- PASS: TestNavigator_ValidateCurrentScreen (0.00s)
=== RUN   TestNavigator_ClearHistory
--- PASS: TestNavigator_ClearHistory (0.00s)
=== RUN   TestNavigator_GetCurrentIndex
--- PASS: TestNavigator_GetCurrentIndex (0.00s)
=== RUN   TestNavigator_AutoSave_NilTrade_NoError
--- PASS: TestNavigator_AutoSave_NilTrade_NoError (0.00s)
PASS
ok      tf-engine/internal/ui   0.333s
=== RUN   TestGetHelpForScreen_SectorSelection
--- PASS: TestGetHelpForScreen_SectorSelection (0.00s)
=== RUN   TestGetHelpForScreen_AllScreens
=== RUN   TestGetHelpForScreen_AllScreens/sector_selection
=== RUN   TestGetHelpForScreen_AllScreens/screener_launch
=== RUN   TestGetHelpForScreen_AllScreens/ticker_entry
=== RUN   TestGetHelpForScreen_AllScreens/checklist
=== RUN   TestGetHelpForScreen_AllScreens/position_sizing
=== RUN   TestGetHelpForScreen_AllScreens/heat_check
=== RUN   TestGetHelpForScreen_AllScreens/trade_entry
=== RUN   TestGetHelpForScreen_AllScreens/calendar
=== RUN   TestGetHelpForScreen_AllScreens/trade_management
--- PASS: TestGetHelpForScreen_AllScreens (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/sector_selection (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/screener_launch (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/ticker_entry (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/checklist (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/position_sizing (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/heat_check (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/trade_entry (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/calendar (0.00s)
    --- PASS: TestGetHelpForScreen_AllScreens/trade_management (0.00s)
=== RUN   TestGetHelpForScreen_UnknownScreen
--- PASS: TestGetHelpForScreen_UnknownScreen (0.00s)
=== RUN   TestGetHelpForScreen_TickerEntry_HasCooldownInfo
--- PASS: TestGetHelpForScreen_TickerEntry_HasCooldownInfo (0.00s)
=== RUN   TestGetHelpForScreen_HeatCheck_HasLimits
--- PASS: TestGetHelpForScreen_HeatCheck_HasLimits (0.00s)
=== RUN   TestGetHelpForScreen_PositionSizing_HasConvictionScale
--- PASS: TestGetHelpForScreen_PositionSizing_HasConvictionScale (0.00s)
=== RUN   TestShowHelpDialog_DoesNotPanic
--- PASS: TestShowHelpDialog_DoesNotPanic (0.13s)
=== RUN   TestShowWelcomeScreen_DoesNotPanic
--- PASS: TestShowWelcomeScreen_DoesNotPanic (0.02s)
=== RUN   TestGetHelpForScreen_Checklist_Has5RequiredItems
--- PASS: TestGetHelpForScreen_Checklist_Has5RequiredItems (0.00s)
=== RUN   TestGetHelpForScreen_Calendar_HasColorCoding
--- PASS: TestGetHelpForScreen_Calendar_HasColorCoding (0.00s)
=== RUN   TestGetHelpForScreen_TradeManagement_MentionsFeatureFlag
--- PASS: TestGetHelpForScreen_TradeManagement_MentionsFeatureFlag (0.00s)
PASS
ok      tf-engine/internal/ui/help      0.457s
FAIL    tf-engine/internal/ui/screens [build failed]
=== RUN   TestCooldownTimer_StartsAtFullDuration
--- PASS: TestCooldownTimer_StartsAtFullDuration (0.06s)
=== RUN   TestCooldownTimer_CountsDown
--- PASS: TestCooldownTimer_CountsDown (1.50s)
=== RUN   TestCooldownTimer_CallsOnComplete
--- PASS: TestCooldownTimer_CallsOnComplete (1.50s)
=== RUN   TestCooldownTimer_Stop
--- PASS: TestCooldownTimer_Stop (1.60s)
=== RUN   TestCooldownTimer_Reset
--- PASS: TestCooldownTimer_Reset (1.00s)
=== RUN   TestNewCooldownTimerFromTime_AlreadyComplete
--- PASS: TestNewCooldownTimerFromTime_AlreadyComplete (0.00s)
=== RUN   TestNewCooldownTimerFromTime_PartiallyComplete
--- PASS: TestNewCooldownTimerFromTime_PartiallyComplete (0.00s)
=== RUN   TestCooldownTimer_MultipleStopCalls_NoError
--- PASS: TestCooldownTimer_MultipleStopCalls_NoError (0.00s)
=== RUN   TestCooldownTimer_ZeroDuration
--- PASS: TestCooldownTimer_ZeroDuration (0.10s)
=== RUN   TestCooldownTimer_NegativeDuration
--- PASS: TestCooldownTimer_NegativeDuration (0.10s)
=== RUN   TestCooldownTimer_GetRemaining_BeforeStart
--- PASS: TestCooldownTimer_GetRemaining_BeforeStart (0.00s)
PASS
ok      tf-engine/internal/widgets      6.170s
?       tf-engine/scripts       [no test files]
FAIL
WARNING: Some tests failed, but continuing build...

[3/5] Formatting code...
integration_test.go
internal\models\settings.go
internal\storage\settings.go
internal\ui\screens\settings.go
  Γ£ô Code formatted

[4/5] Building tf-engine.exe...
  Γ£ô Build successful

[5/5] Build complete!

========================================
Build Information
========================================
11/04/2025  03:07 PM        42,652,352 tf-engine.exe

Policy file hash:
6ae095e862638de9bb2f124e0482bddc3b50059b2189fd63c78f85f09e387c1b

========================================
Ready to test! Run: dist\tf-engine.exe
========================================