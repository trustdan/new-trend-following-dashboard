# TF-Engine Windows Installer Guide

## Overview

The TF-Engine installer (`TFEngine-Setup-1.0.0.exe`) is a professional Windows installer built with NSIS (Nullsoft Scriptable Install System). It provides a smooth installation experience with proper registry integration, Start Menu shortcuts, and clean uninstallation.

## Build Information

**Installer Version:** 1.0.0
**Build Date:** November 3, 2025
**Installer Size:** 1.8 MB (compressed from 3.3 MB)
**Compression:** zlib (57.4% compression ratio)
**Target Platform:** Windows 10/11 (64-bit)

## What's Included

The installer packages the following components:

### Core Application
- `tf-engine.exe` (3.0 MB) - Main application executable
- `policy.v1.json` (9.2 KB) - Strategy and sector configuration
- `README.md` (6.8 KB) - Project documentation
- `LICENSE.txt` (1.1 KB) - MIT License

### Configuration
- `feature.flags.json` - Feature toggles (Phase 2 features OFF by default)
- `data/` directory - Created for trade history storage
- `data/ui/` directory - Created for UI state persistence

### Registry Entries
- `HKLM\Software\TF-Engine\InstallPath` - Installation directory
- `HKLM\Software\TF-Engine\Version` - Installed version number
- Uninstall registry keys for Add/Remove Programs integration

### Shortcuts
- Start Menu folder: `TF-Engine`
  - TF-Engine.lnk - Launch application
  - Uninstall.lnk - Remove application
- Desktop shortcut (optional, user-selectable during install)

## Installation Process

### Prerequisites
- Windows 10 or Windows 11 (64-bit)
- Administrator privileges (required for Program Files installation)
- 10 MB free disk space

### Installation Steps

1. **Download** the installer: `TFEngine-Setup-1.0.0.exe`

2. **Run the installer** (right-click → Run as Administrator if needed)

3. **Welcome Screen**
   - Click "Next" to continue

4. **License Agreement**
   - Review the MIT License
   - Click "I Agree" to accept

5. **Choose Install Location**
   - Default: `C:\Program Files\TF-Engine`
   - Can be customized to any directory
   - Click "Install" to begin

6. **Installation Progress**
   - Files are extracted and copied
   - Registry entries are created
   - Shortcuts are generated
   - Takes 5-10 seconds

7. **Desktop Shortcut Prompt**
   - Dialog: "Create desktop shortcut?"
   - Choose Yes or No

8. **Completion Screen**
   - Option: "Launch TF-Engine" checkbox (checked by default)
   - Click "Finish"

### Post-Installation

After successful installation:
- Application is registered in Windows Add/Remove Programs
- Start Menu folder is created with launch icon
- Desktop shortcut is created (if selected)
- Application can be launched immediately

## Testing the Installer

### Test Checklist

#### Pre-Installation Tests
- [ ] Verify installer size (1.8 MB)
- [ ] Check digital signature (if code-signed)
- [ ] Scan with antivirus (should be clean)

#### Installation Tests
- [ ] **Fresh Install** - Install on clean system
  - Verify all files are copied to `C:\Program Files\TF-Engine`
  - Check that `data/` and `data/ui/` directories exist
  - Confirm feature.flags.json contains Phase 2 flags set to false
  - Test Start Menu shortcut launches application
  - Test Desktop shortcut launches application (if created)

- [ ] **Version Check** - Install when already installed
  - Installer should detect existing version
  - Dialog: "TF-Engine is already installed (version 1.0.0). Continue?"
  - Can proceed or abort

- [ ] **Custom Directory** - Install to non-default location
  - Choose custom path (e.g., `D:\Apps\TF-Engine`)
  - Verify installation completes successfully
  - Check registry points to custom location

#### Application Launch Tests
- [ ] Launch from Start Menu
  - Navigate to Start → TF-Engine → TF-Engine
  - Application should open to Sector Selection screen

- [ ] Launch from Desktop shortcut (if created)
  - Double-click desktop icon
  - Application should open normally

- [ ] Launch from executable
  - Navigate to install directory
  - Double-click tf-engine.exe
  - Application should open normally

#### Registry Verification
- [ ] Open Registry Editor (regedit)
- [ ] Check `HKLM\Software\TF-Engine`:
  - InstallPath = installation directory
  - Version = "1.0.0"
- [ ] Check `HKLM\Software\Microsoft\Windows\CurrentVersion\Uninstall\TF-Engine`:
  - DisplayName = "TF-Engine"
  - DisplayVersion = "1.0.0"
  - Publisher = "TF Systems"
  - UninstallString = path to Uninstall.exe

#### Uninstallation Tests
- [ ] **Uninstall via Start Menu**
  - Start → TF-Engine → Uninstall
  - Confirm uninstall dialog appears
  - Dialog: "Delete trade history and settings?"
  - Test both Yes and No options

- [ ] **Uninstall via Add/Remove Programs**
  - Settings → Apps → Installed apps
  - Find "TF-Engine"
  - Click Uninstall
  - Follow prompts

- [ ] **Verify Clean Removal**
  - Check that `C:\Program Files\TF-Engine` is removed (or empty if data preserved)
  - Check that Start Menu folder is removed
  - Check that Desktop shortcut is removed
  - Check that registry keys are removed
  - If "Delete trade history" = No, verify `data/` folder still exists

### Expected Behaviors

#### Installation
- **Admin Required**: Installer requests elevation (UAC prompt)
- **Duplicate Detection**: Warns if already installed
- **Progress Display**: Shows file copying progress
- **Desktop Shortcut**: Optional, user-selectable
- **Auto-Launch**: Can launch immediately after install

#### Uninstallation
- **Data Preservation**: Prompts before deleting trade history
- **Clean Removal**: Removes all files except preserved data
- **Shortcut Cleanup**: Removes all shortcuts
- **Registry Cleanup**: Removes all registry entries

## Feature Flags Configuration

The installer creates `feature.flags.json` with Phase 2 features disabled by default:

```json
{
  "trade_management": false,
  "sample_data_generator": false,
  "vimium_mode": false,
  "advanced_analytics": false
}
```

To enable Phase 2 features after installation:
1. Navigate to installation directory
2. Edit `feature.flags.json`
3. Set desired features to `true`
4. Save and restart TF-Engine

## Troubleshooting

### Installation Issues

**Problem:** "Access Denied" during installation
**Solution:** Run installer as Administrator (right-click → Run as Administrator)

**Problem:** "Cannot write to C:\Program Files"
**Solution:** Choose a user-writable directory (e.g., `C:\Users\YourName\TF-Engine`)

**Problem:** Installer hangs or freezes
**Solution:** Close installer, restart computer, run installer again

**Problem:** Antivirus blocks installer
**Solution:** Add exception for TFEngine-Setup-1.0.0.exe (false positive)

### Launch Issues

**Problem:** Application doesn't launch from shortcut
**Solution:**
- Right-click shortcut → Properties → Verify Target path
- Navigate to install directory and launch exe directly
- Reinstall application

**Problem:** "Missing DLL" error on launch
**Solution:**
- Application is self-contained Go binary (no DLLs needed)
- If error persists, check Windows Event Viewer for details

**Problem:** Policy file not found
**Solution:**
- Verify `policy.v1.json` exists in installation directory
- Re-run installer to restore missing files

### Uninstallation Issues

**Problem:** Uninstaller doesn't remove all files
**Solution:**
- Manually delete installation directory
- Run CCleaner or similar to clean registry

**Problem:** Start Menu entries remain after uninstall
**Solution:**
- Right-click Start Menu → Open file location
- Manually delete TF-Engine folder

## Building the Installer from Source

If you need to rebuild the installer:

### Requirements
- NSIS 3.x installed (`choco install nsis` or download from nsis.sourceforge.io)
- Go 1.21+ for building the executable
- Git for version control

### Build Steps

```bash
# 1. Build the Go executable
mkdir dist
go build -o dist\tf-engine.exe .

# 2. Copy required files to dist
cp data\policy.v1.json dist\
cp README.md dist\

# 3. Ensure LICENSE.txt exists in root
# (already created by installer build process)

# 4. Build the installer
makensis tf-engine-installer.nsi

# Output: TFEngine-Setup-1.0.0.exe (1.8 MB)
```

### Optional: Add Icon

To add a custom icon to the installer:

1. Create or obtain `icon.ico` (256x256, multiple sizes)
2. Copy to `assets\icon.ico`
3. Uncomment these lines in `tf-engine-installer.nsi`:
   ```nsi
   !define MUI_ICON "assets\icon.ico"
   !define MUI_UNICON "assets\icon.ico"
   ```
4. Rebuild: `makensis tf-engine-installer.nsi`

### Versioning

To create a new version (e.g., 1.1.0):

1. Edit `tf-engine-installer.nsi`
2. Change `!define APP_VERSION "1.1.0"`
3. Rebuild installer
4. Output: `TFEngine-Setup-1.1.0.exe`

## Distribution

### Recommended Distribution Methods

1. **GitHub Releases**
   - Tag release: `git tag v1.0.0`
   - Push tag: `git push origin v1.0.0`
   - Create release on GitHub
   - Attach `TFEngine-Setup-1.0.0.exe`

2. **Direct Download**
   - Host on website
   - Provide SHA256 checksum for verification
   - Include installation instructions

3. **Code Signing (Recommended for Production)**
   - Obtain code signing certificate (Sectigo, DigiCert)
   - Sign installer with `signtool.exe`
   - Prevents Windows SmartScreen warnings

### Checksum Generation

```bash
# Generate SHA256 checksum
sha256sum TFEngine-Setup-1.0.0.exe > TFEngine-Setup-1.0.0.exe.sha256

# Verify checksum (users)
sha256sum -c TFEngine-Setup-1.0.0.exe.sha256
```

## Support

For installation support or issues:
- File issue on GitHub: https://github.com/[repo]/issues
- Email: support@tfsystems.com
- Documentation: README.md

## License

TF-Engine is released under the MIT License. See LICENSE.txt for details.

---

**Last Updated:** November 3, 2025
**Installer Version:** 1.0.0
**Build Tool:** NSIS 3.x
