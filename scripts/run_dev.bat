@echo off
REM Development run script for TF-Engine 2.0

echo Starting TF-Engine 2.0 in development mode...
echo.

REM Download dependencies if needed
go mod tidy

REM Run the application
go run main.go

pause
