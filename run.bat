@echo off
REM ============================================================================
REM TF-Engine 2.0 - Quick Run Script
REM ============================================================================
REM Builds and runs the application in one command
REM ============================================================================

echo.
echo ========================================
echo TF-Engine 2.0 Quick Run
echo ========================================
echo.

REM Quick build
echo [Step 1/2] Building application...
call quick-build.bat
if errorlevel 1 (
    echo.
    echo ERROR: Build failed
    pause
    exit /b 1
)

REM Run application
echo [Step 2/2] Launching tf-engine.exe...
echo.
echo ========================================
echo.
start "" dist\tf-engine.exe

echo Application started!
echo.
echo Check the application window...
echo.
