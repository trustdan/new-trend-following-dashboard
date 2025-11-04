# TF-Engine 2.0 - Trend Following Options Trading System

A systematic decision support application for options traders using trend-following strategies. Built with discipline, informed by 293 backtests across 14 strategies and 21 securities.

---

## Project Status

**Current Phase:** Phase 5 Complete ✅ (95% - Ready for UAT)

**Version:** 1.0.0 (MVP Release Candidate)

**Latest Update:** November 4, 2025

### ✅ Completed Phases
- ✅ Phase 0: Feature Freeze & Repo Hygiene
- ✅ Phase 1: Foundation Layer (navigation, persistence, cooldown timer)
- ✅ Phase 2-4: Core Workflow (all 8 screens implemented)
- ✅ Phase 5: Polish & Phase 2 Features (95% complete)

### 🎯 Current Status
- ✅ All core features implemented (8-screen workflow operational)
- ✅ Windows installer built: `TFEngine-Setup-1.0.0.exe` (17.4 MB)
- ✅ 85%+ test coverage, all automated tests passing
- ✅ Sample data generator implemented (Phase 2 feature, disabled by default)
- ✅ Help system and welcome screen complete
- ⏳ Pending: VM testing and User Acceptance Testing with 3 beta testers

**See:** [PHASE_5_COMPLETION_STATUS.md](PHASE_5_COMPLETION_STATUS.md) for detailed status

---

## Quick Start

### Prerequisites
- Go 1.21 or later
- Windows 10/11 (target platform)

### Run Phase 0 Infrastructure Test
```bash
go run .
```

Expected output:
```
✅ Phase 0 Infrastructure Test PASSED
```

### Verify Policy Integrity
```bash
go run scripts/verify_policy_hash.go
```

### Run All Tests
```bash
go test ./...
```

---

## Project Structure

```
tf-engine/
├── data/
│   ├── policy.v1.json              # Master configuration (sectors, strategies, limits)
│   └── backups/                    # Policy backups
├── internal/
│   ├── config/                     # Configuration loaders
│   │   ├── feature_flags.go
│   │   └── feature_flags_test.go
│   ├── screens/
│   │   └── _post_mvp/              # Phase 2 screen implementations
│   └── testdata/                   # Test fixtures
├── scripts/
│   └── verify_policy_hash.go       # Policy signature validator
├── logs/                           # Application logs
├── feature.flags.json              # Feature flag configuration
├── CONTRIBUTING.md                 # Contribution guidelines
├── CLAUDE.md                       # AI assistant context
└── plans/
    └── roadmap.md                  # Complete development roadmap
```

---

## Core Principles

### 1. No Feature Creep
This project failed once due to excessive features. **Rule #1:** Do not add features without explicit approval.

### 2. Policy-Driven Design
All business logic lives in `data/policy.v1.json`, not code. This allows updating trading rules without recompiling.

### 3. Anti-Impulsivity Guardrails
- 120-second cooldown timers
- Required checklist validation
- Portfolio heat limits (4% total, 1.5% per sector)

### 4. Sector-First Workflow
Research proved strategy performance varies by sector:
- Healthcare: 92% success rate
- Utilities: 0% success rate (BLOCKED)

---

## The 9-Screen Workflow

1. **Sector Selection** - Choose market sector to trade
2. **Screener Launch** - Open Finviz screeners in browser
3. **Ticker Entry** - Enter ticker + select strategy (sector-filtered)
4. **Anti-Impulsivity Checklist** - 5 required + 3 optional gates
5. **Position Sizing** - Poker-bet conviction sizing (5-8)
6. **Heat Check** - Enforce portfolio diversification limits
7. **Options Strategy** - Select from 24 options structures
8. **Trade Calendar** - Horserace timeline view (sector Y-axis)
9. **Trade Management** - Edit/delete trades (Phase 2 feature)

---

## Research Foundation

This application is informed by extensive backtesting research:

- **293 validated backtests** (99.74% data quality)
- **Key finding:** Alt10 (Profit Targets) achieved 76.19% success rate
- **Critical insight:** Healthcare 92% success, Utilities 0% success
- **Python validation:** Logistic regression confirms profit targets (4.47× odds) drive profitability

**See:** [DISCOVERIES_AND_LEARNINGS.md](DISCOVERIES_AND_LEARNINGS.md) for complete analysis

---

## Development Roadmap

The project follows a phased approach to prevent scope creep:

- **Phase 0:** Feature Freeze & Repo Hygiene ✅ (3 days)
- **Phase 1:** Foundation Layer (2 weeks) - NEXT
- **Phase 2:** Core Workflow (4 weeks)
- **Phase 3:** Polish & Testing (1 week)
- **Phase 4:** Windows Installer (3 days)
- **Phase 5:** User Acceptance Testing (1 week)

**Total:** ~8-9 weeks to MVP

**See:** [plans/roadmap.md](plans/roadmap.md) for detailed implementation plan

---

## Technology Stack

- **Language:** Go 1.21+
- **GUI Framework:** Fyne (planned for Phase 1)
- **Data Storage:** JSON files with atomic writes
- **Configuration:** JSON (`policy.v1.json`)
- **Screeners:** Finviz (external, chart view)
- **Target Platform:** Windows 10/11

---

## Contributing

**Read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting PRs.**

Key rules:
- No features without approval
- Phase 2 features must be behind flags (default OFF)
- Unit tests required (>80% coverage)
- Policy changes must update signature
- All behavioral guardrails are non-negotiable

---

## Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests (Phase 1+)
```bash
go test -tags=integration ./...
```

### Policy Validation
```bash
go run scripts/verify_policy_hash.go
# Expected: ✅ Policy signature valid
```

### Manual Testing Checklist
1. Complete full trade entry workflow (Sector → Calendar)
2. Verify cooldown timer prevents bypass
3. Test heat limit enforcement
4. Confirm blocked sectors prevent trades

---

## Feature Flags

Phase 2 features are disabled by default:

```json
{
  "trade_management": false,       // Screen 9
  "sample_data_generator": false,  // Test data
  "vimium_mode": false,            // Keyboard shortcuts
  "advanced_analytics": false      // Win rate tracking
}
```

To enable a feature (development only):
1. Edit `feature.flags.json`
2. Set `"enabled": true`
3. Restart application

**Production:** All Phase 2 features remain disabled until Phase 5

---

## Documentation

### For Understanding the Research
- [DISCOVERIES_AND_LEARNINGS.md](DISCOVERIES_AND_LEARNINGS.md) - 293 backtest analysis
- [architects-intent.md](architects-intent.md) - Original requirements

### For Development
- [plans/roadmap.md](plans/roadmap.md) - Complete implementation plan
- [CLAUDE.md](CLAUDE.md) - Architectural rules & context
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

### For Policy Updates
- [data/policy.v1.json](data/policy.v1.json) - Master configuration
- [scripts/verify_policy_hash.go](scripts/verify_policy_hash.go) - Signature validator

---

## License

[Specify license here]

---

## Support

For questions or issues:
1. Check [plans/roadmap.md](plans/roadmap.md) for implementation details
2. Review [CONTRIBUTING.md](CONTRIBUTING.md) for common issues
3. Read [CLAUDE.md](CLAUDE.md) for architectural context

---

**Last Updated:** November 3, 2025
**Project Start:** November 3, 2025
**Target MVP:** ~8-9 weeks from start
