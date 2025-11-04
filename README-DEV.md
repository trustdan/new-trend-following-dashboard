# TF-Engine 2.0 - Development Guide

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Windows OS (primary target)
- Git

### Installation

1. **Install Go dependencies:**
```bash
go mod download
```

2. **Run the application:**
```bash
go run main.go
```

3. **Build executable:**
```bash
go build -o tf-engine.exe
```

## Project Structure

```
tf-engine/
├── main.go                          # Application entry point
├── internal/
│   ├── appcore/
│   │   └── state.go                 # Global application state
│   ├── models/
│   │   ├── policy.go                # Policy data structures
│   │   ├── trade.go                 # Trade data structures
│   │   └── settings.go              # Settings data structures
│   └── ui/
│       ├── theme.go                 # Day/night mode themes
│       ├── dashboard.go             # Main dashboard
│       └── screens/                 # Individual workflow screens
│           ├── sector_selection.go  # Screen 1
│           ├── screener_launch.go   # Screen 2
│           ├── ticker_entry.go      # Screen 3
│           ├── checklist.go         # Screen 4
│           ├── position_sizing.go   # Screen 5
│           ├── heat_check.go        # Screen 6
│           ├── trade_entry.go       # Screen 7
│           └── calendar.go          # Screen 8
├── data/
│   └── policy.v1.json               # Master configuration (READ ONLY)
└── go.mod                           # Go module definition
```

## Development Workflow

### Adding a New Screen

1. Create new file in `internal/ui/screens/`
2. Implement the screen struct with `Render()` method
3. Add navigation logic in the calling screen
4. Update auto-save logic if screen modifies trade state

### Modifying Policy

1. Edit `data/policy.v1.json`
2. Update corresponding structs in `internal/models/policy.go`
3. Restart application to load new policy

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/models

# Run with coverage
go test -cover ./...
```

## Key Architectural Principles

1. **Policy-Driven**: All business rules come from `policy.v1.json`
2. **Linear Workflow**: Users progress through screens 1→8 sequentially
3. **Anti-Impulsivity**: Cooldown timers and checklists are mandatory
4. **Auto-Save**: State persists after each screen transition
5. **No Feature Creep**: Only implement explicitly requested features

## Next Steps for Full Implementation

1. Complete screen implementations (currently stubs)
2. Add cooldown timer widget
3. Implement FINVIZ URL launcher (browser integration)
4. Add heat calculation logic
5. Build calendar timeline widget
6. Add trade persistence (JSON file I/O)
7. Implement Vimium keyboard navigation
8. Add sample data generator
9. Create help system
10. Polish day/night mode themes

## Common Issues

### "module not found"
Run `go mod tidy` to sync dependencies.

### "policy.v1.json not found"
Ensure you're running from the project root directory where `data/` exists.

### Theme not applying
Check that `fyneApp.Settings().SetTheme()` is called in `main.go`.

## Resources

- [Fyne Documentation](https://docs.fyne.io/)
- [Go Documentation](https://go.dev/doc/)
- [Project Architecture](CLAUDE.md)
