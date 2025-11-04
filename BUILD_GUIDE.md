# TF-Engine 2.0 - Build Guide

Complete guide to building, testing, and distributing TF-Engine 2.0.

---

## Quick Reference

| Script | Purpose | Use When |
|--------|---------|----------|
| `build.bat` | Full build with tests | Before commits, releases |
| `quick-build.bat` | Fast build without tests | During development |
| `run.bat` | Build and run in one command | Quick testing |
| `build-installer.bat` | Create Windows installer | Creating releases |
| `clean.bat` | Remove build artifacts | Fresh start needed |

---

## Prerequisites

### Required:
- **Go 1.21+** - Install from: https://go.dev/dl/
- **Git** - For version control

### Optional (for installer):
- **NSIS** - Install from: https://nsis.sourceforge.io/Download
  - Default path: `C:\Program Files (x86)\NSIS\`

---

## Build Scripts Explained

### 1. `build.bat` - Full Build (Recommended)

**What it does:**
1. Syncs policy.v1.json to dist/
2. Runs all tests (`go test`)
3. Formats code (`go fmt`)
4. Builds `dist\tf-engine.exe`
5. Displays build info and SHA256 hash

**When to use:**
- Before committing code
- Before creating a release
- When you want to verify tests pass

**Usage:**
```cmd
build.bat
```

**Time:** ~10-30 seconds (depending on tests)

---

### 2. `quick-build.bat` - Fast Build (Development)

**What it does:**
1. Syncs policy.v1.json to dist/
2. Builds `dist\tf-engine.exe`
3. **Skips tests and formatting** (fast!)

**When to use:**
- During rapid development iteration
- When you just want to test UI changes
- When you know tests pass

**Usage:**
```cmd
quick-build.bat
```

**Time:** ~3-5 seconds

---

### 3. `run.bat` - Build and Run

**What it does:**
1. Calls `quick-build.bat`
2. Launches `dist\tf-engine.exe`

**When to use:**
- Quick testing during development
- One-command build-and-run workflow

**Usage:**
```cmd
run.bat
```

**Time:** ~3-5 seconds + app startup

---

### 4. `build-installer.bat` - Create Windows Installer

**What it does:**
1. Calls `build.bat` (full build with tests)
2. Verifies required files (LICENSE.txt, README.md, etc.)
3. Compiles NSIS installer script
4. Creates `TFEngine-Setup-1.0.0.exe`
5. Generates SHA256 hash for distribution

**When to use:**
- Creating release packages
- Distributing to users
- Before uploading to GitHub releases

**Prerequisites:**
- NSIS installed at `C:\Program Files (x86)\NSIS\`
- LICENSE.txt exists
- README.md exists

**Usage:**
```cmd
build-installer.bat
```

**Output:**
- `TFEngine-Setup-1.0.0.exe` (installer)
- `TFEngine-Setup-1.0.0.exe.sha256` (hash for verification)

**Time:** ~30-60 seconds

---

### 5. `clean.bat` - Clean Build Artifacts

**What it does:**
1. Removes `dist\tf-engine.exe`
2. Removes old installers
3. **Optionally** removes logs
4. **Optionally** removes trade data (asks confirmation)

**When to use:**
- Before fresh builds
- When disk space is low
- After testing is complete

**Usage:**
```cmd
clean.bat
```

**Interactive prompts:**
- "Delete logs? (y/N)"
- "Delete trade data? [WARNING] (y/N)"

**Time:** Instant

---

## Typical Workflows

### Development Workflow

**Rapid iteration (making UI changes):**
```cmd
# Edit code...
quick-build.bat
dist\tf-engine.exe

# Or even faster:
run.bat
```

**Before committing code:**
```cmd
build.bat
# Review test results
# If all pass, commit
```

---

### Release Workflow

**Creating a new release:**

1. **Update version numbers:**
   - Update `APP_VERSION` in `tf-engine-installer.nsi`
   - Update version in `README.md`

2. **Full build and test:**
   ```cmd
   build.bat
   ```
   Verify all tests pass.

3. **Create installer:**
   ```cmd
   build-installer.bat
   ```

4. **Test installer:**
   - Run `TFEngine-Setup-1.0.0.exe`
   - Install to test location
   - Verify application runs
   - Test core features
   - Uninstall

5. **Create GitHub release:**
   - Tag version: `git tag v1.0.0`
   - Push tag: `git push origin v1.0.0`
   - Upload:
     - `TFEngine-Setup-1.0.0.exe`
     - `TFEngine-Setup-1.0.0.exe.sha256`

---

## Build Outputs

### After `build.bat` or `quick-build.bat`:

```
dist/
├── tf-engine.exe          # Main executable (40-45 MB)
├── policy.v1.json         # Policy configuration (synced from data/)
├── data/                  # Runtime data directory
│   └── ui/
│       └── trades_in_progress.json
├── logs/                  # Application logs
│   └── tf-engine_2025-*.log
└── README.md              # User documentation (copied for installer)
```

### After `build-installer.bat`:

```
project-root/
├── TFEngine-Setup-1.0.0.exe        # Windows installer (~45 MB)
└── TFEngine-Setup-1.0.0.exe.sha256 # SHA256 hash for verification
```

---

## Troubleshooting

### "go: command not found"

**Problem:** Go is not installed or not in PATH

**Fix:**
1. Install Go from https://go.dev/dl/
2. Restart terminal
3. Verify: `go version`

---

### "NSIS not found"

**Problem:** NSIS is not installed or installed in non-default location

**Fix Option 1 (Install NSIS):**
1. Download from https://nsis.sourceforge.io/Download
2. Install to default location: `C:\Program Files (x86)\NSIS\`
3. Run `build-installer.bat`

**Fix Option 2 (Custom NSIS path):**
Edit `build-installer.bat` line 15:
```batch
set "NSIS_PATH=C:\Your\Custom\Path\makensis.exe"
```

---

### "Tests failed"

**Problem:** Some tests are failing

**What `build.bat` does:**
- Shows WARNING but continues building
- Use this to identify failing tests

**What to do:**
1. Review test output
2. Fix failing tests
3. Re-run `build.bat`

**If you need to build urgently:**
Use `quick-build.bat` (skips tests)

---

### "Missing required files" (installer)

**Problem:** LICENSE.txt or README.md not found

**Fix:**
1. Ensure `LICENSE.txt` exists in project root
2. Ensure `README.md` exists in project root
3. Run `build-installer.bat` again

---

### Large executable size (40+ MB)

**Why?** Go includes runtime and all dependencies in single binary.

**This is normal and expected.**

**Benefits:**
- No external dependencies needed
- Easy distribution
- Works on any Windows machine

**To reduce size (advanced):**
```cmd
go build -ldflags="-s -w" -o dist\tf-engine.exe .
```
This strips debug symbols (~10-15% reduction).

---

## Build Configuration

### Policy File Sync

All build scripts automatically sync `data/policy.v1.json` to `dist/policy.v1.json`.

**Why?** The application reads policy from:
1. Current directory
2. `dist/` directory
3. Executable directory

**Important:** Always edit `data/policy.v1.json` (source), not `dist/policy.v1.json` (copy).

---

### Feature Flags

The installer creates `feature.flags.json` with Phase 2+ features disabled:

```json
{
  "trade_management": false,
  "sample_data_generator": false,
  "vimium_mode": false,
  "advanced_analytics": false
}
```

Users can edit this file to enable experimental features.

---

## Advanced: Manual Build Commands

If you prefer manual control:

### Build executable:
```cmd
go build -o dist\tf-engine.exe .
```

### Run tests:
```cmd
go test ./...
```

### Run tests with verbose output:
```cmd
go test ./... -v
```

### Format code:
```cmd
go fmt ./...
```

### Build installer (manual):
```cmd
"C:\Program Files (x86)\NSIS\makensis.exe" tf-engine-installer.nsi
```

### Generate SHA256 hash:
```cmd
certutil -hashfile TFEngine-Setup-1.0.0.exe SHA256
```

---

## Continuous Integration

For GitHub Actions or other CI/CD:

### Basic CI build:
```yaml
- name: Build
  run: go build -o dist/tf-engine.exe .

- name: Test
  run: go test ./... -v
```

### Create release:
```yaml
- name: Build Installer (Windows)
  run: |
    choco install nsis
    ./build-installer.bat
```

---

## File Locations Reference

### Source Files:
- `data/policy.v1.json` - Policy source (edit this)
- `main.go` - Application entry point
- `internal/` - Go packages

### Build Outputs:
- `dist/tf-engine.exe` - Built executable
- `dist/policy.v1.json` - Policy copy (auto-synced)

### Installer Files:
- `tf-engine-installer.nsi` - NSIS installer script
- `LICENSE.txt` - Required for installer
- `README.md` - Required for installer

### Generated Files:
- `TFEngine-Setup-*.exe` - Windows installer
- `TFEngine-Setup-*.exe.sha256` - Hash for verification
- `logs/*.log` - Application logs

---

## Best Practices

### During Development:
1. Use `quick-build.bat` or `run.bat` for speed
2. Run `build.bat` before commits (verify tests)
3. Use `clean.bat` if builds behave strangely

### Before Release:
1. Update version in `tf-engine-installer.nsi`
2. Run `build.bat` and verify all tests pass
3. Run `build-installer.bat`
4. Test installer on clean VM or test machine
5. Generate release notes
6. Create GitHub release with installer + hash

### After Major Changes:
1. Run `clean.bat` for fresh start
2. Run `build.bat` to rebuild everything
3. Test thoroughly before distribution

---

## FAQ

### Q: Do I need to run `build.bat` every time?

**A:** No! Use `quick-build.bat` during development. Only use `build.bat` before commits/releases.

---

### Q: Can I skip tests completely?

**A:** Yes, `quick-build.bat` skips tests. But run `build.bat` before committing to verify tests pass.

---

### Q: How do I distribute the application?

**A:** Use `build-installer.bat` to create `TFEngine-Setup-1.0.0.exe`, then share that file + the .sha256 hash.

---

### Q: Why is the installer so large (40+ MB)?

**A:** Go includes the entire runtime in the binary (no dependencies needed). This is normal and makes distribution easier.

---

### Q: Can I build on Mac/Linux?

**A:** The Go code will build, but the installer scripts are Windows-specific (NSIS). For Mac/Linux, use:
```bash
go build -o tf-engine .
```

---

### Q: What if I modify policy.v1.json?

**A:** Just rebuild. All build scripts automatically sync `data/policy.v1.json` to `dist/policy.v1.json`.

---

### Q: How do I test the installer?

**A:** Run `TFEngine-Setup-1.0.0.exe`, install to a test location, verify the app works, then uninstall.

---

## Summary

**For development:**
```cmd
run.bat              # Fastest: build and run
quick-build.bat      # Fast: build only
```

**Before committing:**
```cmd
build.bat            # Full build with tests
```

**For releases:**
```cmd
build-installer.bat  # Create installer
```

**For cleanup:**
```cmd
clean.bat            # Remove build artifacts
```

---

**Last Updated:** November 4, 2025
**TF-Engine Version:** 2.0.0
**Build System Version:** 1.0.0
