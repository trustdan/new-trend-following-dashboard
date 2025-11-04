@echo off
REM ============================================================================
REM TF-Engine 2.0 - Installer Build Script
REM ============================================================================
REM This script builds the Windows installer using NSIS (Nullsoft Scriptable Install System)
REM Prerequisites: NSIS must be installed at C:\Program Files (x86)\NSIS\
REM ============================================================================

setlocal enabledelayedexpansion

echo.
echo ========================================
echo TF-Engine 2.0 Installer Build
echo ========================================
echo.

REM Check if NSIS is installed
set "NSIS_PATH=C:\Program Files (x86)\NSIS\makensis.exe"
if not exist "%NSIS_PATH%" (
    echo ERROR: NSIS not found at: %NSIS_PATH%
    echo.
    echo Please install NSIS from: https://nsis.sourceforge.io/Download
    echo.
    pause
    exit /b 1
)

REM Step 1: Build the application first
echo [Step 1/5] Building application...
echo.
call build.bat
if errorlevel 1 (
    echo.
    echo ERROR: Application build failed. Cannot create installer.
    echo.
    pause
    exit /b 1
)
echo.
echo   [OK] Application built successfully
echo.

REM Step 2: Verify required files exist
echo [Step 2/5] Verifying required files...
set "MISSING_FILES="

if not exist "dist\tf-engine.exe" set "MISSING_FILES=!MISSING_FILES! tf-engine.exe"
if not exist "dist\policy.v1.json" set "MISSING_FILES=!MISSING_FILES! policy.v1.json"
if not exist "LICENSE.txt" set "MISSING_FILES=!MISSING_FILES! LICENSE.txt"
if not exist "README.md" set "MISSING_FILES=!MISSING_FILES! README.md"

if not "!MISSING_FILES!"=="" (
    echo   ERROR: Missing required files:!MISSING_FILES!
    echo.
    pause
    exit /b 1
)

REM Copy README.md to dist for installer
copy /Y README.md dist\README.md >nul
echo   [OK] All required files present
echo.

REM Step 3: Clean old installer if exists
echo [Step 3/5] Cleaning old installer...
if exist "TFEngine-Setup-*.exe" (
    del /Q "TFEngine-Setup-*.exe"
    echo   [OK] Old installer removed
) else (
    echo   [OK] No old installer to remove
)
echo.

REM Step 4: Build installer with NSIS
echo [Step 4/5] Building installer with NSIS...
echo.
echo   Running: "%NSIS_PATH%" tf-engine-installer.nsi
echo.
"%NSIS_PATH%" tf-engine-installer.nsi
if errorlevel 1 (
    echo.
    echo ERROR: NSIS compilation failed
    echo.
    pause
    exit /b 1
)
echo.
echo   [OK] Installer compiled successfully
echo.

REM Step 5: Generate SHA256 hash
echo [Step 5/5] Generating SHA256 hash...
for %%F in (TFEngine-Setup-*.exe) do (
    set "INSTALLER_NAME=%%F"
    certutil -hashfile "%%F" SHA256 > "%%F.sha256"
    echo   [OK] Hash saved to: %%F.sha256
)
echo.

REM Display results
echo ========================================
echo Build Complete!
echo ========================================
echo.
echo Installer: %INSTALLER_NAME%
if exist "%INSTALLER_NAME%" (
    dir "%INSTALLER_NAME%" | findstr "TFEngine-Setup"
)
echo.
echo SHA256 Hash:
if exist "%INSTALLER_NAME%.sha256" (
    type "%INSTALLER_NAME%.sha256" | findstr /v "CertUtil"
)
echo.
echo ========================================
echo Next Steps:
echo ========================================
echo 1. Test installer: %INSTALLER_NAME%
echo 2. Distribute both:
echo    - %INSTALLER_NAME%
echo    - %INSTALLER_NAME%.sha256
echo ========================================
echo.

pause
