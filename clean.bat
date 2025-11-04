@echo off
REM ============================================================================
REM TF-Engine 2.0 - Clean Build Artifacts
REM ============================================================================
REM Removes all build artifacts and temporary files
REM ============================================================================

echo.
echo ========================================
echo TF-Engine 2.0 Clean Script
echo ========================================
echo.

set "CLEANED=0"

REM Clean dist directory executables
if exist "dist\tf-engine.exe" (
    echo [1] Removing dist\tf-engine.exe...
    del /Q dist\tf-engine.exe
    set /a CLEANED+=1
)

if exist "dist\tf-engine.exe~" (
    echo [2] Removing dist\tf-engine.exe~ (backup)...
    del /Q dist\tf-engine.exe~
    set /a CLEANED+=1
)

REM Clean installers
if exist "TFEngine-Setup-*.exe" (
    echo [3] Removing installer(s)...
    del /Q TFEngine-Setup-*.exe
    set /a CLEANED+=1
)

if exist "TFEngine-Setup-*.exe.sha256" (
    echo [4] Removing installer hash(es)...
    del /Q TFEngine-Setup-*.exe.sha256
    set /a CLEANED+=1
)

REM Clean test executables
if exist "tf-engine-test.exe" (
    echo [5] Removing test executable...
    del /Q tf-engine-test.exe
    set /a CLEANED+=1
)

REM Clean logs (optional)
echo.
set /p CLEAN_LOGS="Delete logs? (y/N): "
if /i "%CLEAN_LOGS%"=="y" (
    if exist "logs\*.log" (
        echo [6] Removing logs...
        del /Q logs\*.log
        set /a CLEANED+=1
    )
)

REM Clean user data (DANGEROUS - confirm)
echo.
set /p CLEAN_DATA="Delete trade data? [WARNING: Cannot be undone] (y/N): "
if /i "%CLEAN_DATA%"=="y" (
    if exist "internal\ui\data\trades_in_progress.json" (
        echo [7] Removing trade data...
        del /Q internal\ui\data\trades_in_progress.json
        set /a CLEANED+=1
    )
    if exist "dist\data\*.json" (
        echo [8] Removing dist data...
        del /Q dist\data\*.json
        set /a CLEANED+=1
    )
)

echo.
echo ========================================
echo Clean Complete!
echo ========================================
echo Files cleaned: %CLEANED%
echo.
echo Ready for fresh build: build.bat
echo ========================================
echo.
