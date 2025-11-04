@echo off
REM Build script for TF-Engine 2.0
REM This script ensures policy files stay in sync and builds the application

echo ========================================
echo TF-Engine 2.0 Build Script
echo ========================================
echo.

REM Step 1: Sync policy file to dist directory
echo [1/5] Syncing policy.v1.json to dist directory...
copy /Y data\policy.v1.json dist\policy.v1.json
if errorlevel 1 (
    echo ERROR: Failed to copy policy file
    exit /b 1
)
echo   ✓ Policy file synced
echo.

REM Step 2: Run tests
echo [2/5] Running tests...
go test ./... -v
if errorlevel 1 (
    echo WARNING: Some tests failed, but continuing build...
) else (
    echo   ✓ All tests passed
)
echo.

REM Step 3: Format code
echo [3/5] Formatting code...
go fmt ./...
echo   ✓ Code formatted
echo.

REM Step 4: Build application
echo [4/5] Building tf-engine.exe...
go build -o dist\tf-engine.exe .
if errorlevel 1 (
    echo ERROR: Build failed
    exit /b 1
)
echo   ✓ Build successful
echo.

REM Step 5: Display build info
echo [5/5] Build complete!
echo.
echo ========================================
echo Build Information
echo ========================================
dir dist\tf-engine.exe | findstr "tf-engine.exe"
echo.
echo Policy file hash:
certutil -hashfile dist\policy.v1.json SHA256 | findstr /v "hash CertUtil"
echo.
echo ========================================
echo Ready to test! Run: dist\tf-engine.exe
echo ========================================