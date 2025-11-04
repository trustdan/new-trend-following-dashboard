@echo off
REM ============================================================================
REM TF-Engine 2.0 - Quick Build Script (No Tests)
REM ============================================================================
REM Use this for rapid iteration during development
REM Skips tests and formatting checks for speed
REM ============================================================================

echo.
echo ========================================
echo TF-Engine 2.0 Quick Build
echo ========================================
echo.

REM Step 1: Sync policy file
echo [1/2] Syncing policy.v1.json...
if not exist "dist" mkdir dist
copy /Y data\policy.v1.json dist\policy.v1.json >nul
if errorlevel 1 (
    echo   ERROR: Failed to copy policy file
    exit /b 1
)
echo   [OK] Policy synced
echo.

REM Step 2: Build application
echo [2/2] Building tf-engine.exe...
go build -o dist\tf-engine.exe .
if errorlevel 1 (
    echo   ERROR: Build failed
    echo.
    pause
    exit /b 1
)
echo   [OK] Build successful
echo.

REM Display build info
echo ========================================
echo Build Complete!
echo ========================================
dir dist\tf-engine.exe | findstr "tf-engine.exe"
echo.
echo Ready to test: dist\tf-engine.exe
echo.
echo TIP: Use "build.bat" for full build with tests
echo ========================================
echo.
