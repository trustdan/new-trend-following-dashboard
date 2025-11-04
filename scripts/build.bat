@echo off
REM Build script for TF-Engine 2.0

echo Building TF-Engine 2.0...

REM Download dependencies
go mod tidy

REM Build the executable
go build -o tf-engine.exe

if %ERRORLEVEL% == 0 (
    echo.
    echo Build successful! Run tf-engine.exe to start the application.
) else (
    echo.
    echo Build failed. Check errors above.
)

pause
