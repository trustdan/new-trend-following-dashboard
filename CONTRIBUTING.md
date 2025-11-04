# Contributing to TF-Engine 2.0

## Golden Rule: No Feature Creep

**DO NOT add features unless:**
1. Explicitly requested by the project architect, OR
2. Approved in writing before implementation

This rule exists because the previous version failed due to excessive features.

---

## Pull Request Checklist

Before submitting a PR:

- [ ] Feature is in the approved roadmap (`plans/roadmap.md`)
- [ ] If Phase 2 feature, it's behind a feature flag (default OFF)
- [ ] Unit tests added (>80% coverage for new code)
- [ ] Policy signature still validates (`go run scripts/verify_policy_hash.go`)
- [ ] No hardcoded business logic (check `data/policy.v1.json`)
- [ ] Behavioral guardrails preserved (cooldowns, checklist, heat checks)
- [ ] Documentation updated if adding public API
- [ ] All tests pass (`go test ./...`)

---

## Prohibited Changes

❌ Removing cooldown timer
❌ Skipping checklist gates
❌ Bypassing heat caps
❌ Hardcoding sector/strategy mappings
❌ Building custom screener (use FINVIZ)
❌ Adding analytics dashboards without approval

---

## Phase 2 Features

All Phase 2 features MUST:
1. Be in `internal/screens/_post_mvp/` or similar segregation
2. Have a feature flag in `feature.flags.json` (default: false)
3. Not execute if flag is disabled
4. Be documented as "Phase 2" in code comments

---

## Testing Requirements

### Unit Tests
```bash
go test ./...
```
All unit tests must pass. Aim for >80% coverage on new code.

### Integration Tests
```bash
go test -tags=integration ./...
```
Integration tests verify end-to-end workflows.

### Policy Validation
```bash
go run scripts/verify_policy_hash.go
```
This must pass to ensure policy integrity.

### Manual Testing
Complete at least 1 full trade entry workflow manually before submitting PR.

---

## Code Style Guidelines

### Go Code
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) conventions
- Add comments for exported functions and types
- Keep functions focused and small (<50 lines when possible)

### Policy-Driven Design
Business logic belongs in `data/policy.v1.json`, not in Go code:

**❌ Bad:**
```go
if sector == "Healthcare" {
    return []string{"Alt10", "Alt46", "Alt43"}
}
```

**✅ Good:**
```go
sector := policy.GetSector(sectorName)
return sector.AllowedStrategies
```

---

## Commit Message Format

Use clear, descriptive commit messages:

```
[Component] Brief description

Detailed explanation of what changed and why.

- Bullet points for multiple changes
- Reference issue numbers if applicable
```

**Examples:**
- `[Policy] Add Real Estate sector with warning flag`
- `[Screens] Implement Screen 4 checklist validation`
- `[Tests] Add unit tests for heat calculation logic`

---

## Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes
- Follow the coding guidelines above
- Add tests for new functionality
- Update documentation if needed

### 3. Test Locally
```bash
# Run all tests
go test ./...

# Verify policy signature
go run scripts/verify_policy_hash.go

# Build to check for compilation errors
go build
```

### 4. Commit Changes
```bash
git add .
git commit -m "[Component] Your descriptive message"
```

### 5. Push and Create PR
```bash
git push origin feature/your-feature-name
```
Then create a Pull Request on GitHub.

---

## Adding a New Screen

If adding a new screen to the 9-screen workflow:

1. Create screen file in `internal/screens/`
2. Implement the `Screen` interface:
```go
type Screen interface {
    Render() fyne.CanvasObject
    Validate() bool
    OnContinue() error
    OnBack() error
}
```
3. Add screen to Navigator in `internal/ui/navigator.go`
4. Write unit tests in `internal/screens/screenname_test.go`
5. Update `plans/roadmap.md` to mark screen as complete

---

## Adding a New Feature Flag

If adding a Phase 2 feature:

1. Add flag to `feature.flags.json`:
```json
{
  "new_feature": {
    "enabled": false,
    "description": "Brief description",
    "phase": 2,
    "since_version": "2.x.0"
  }
}
```

2. Check flag in code:
```go
flags, _ := config.LoadFeatureFlags("feature.flags.json")
if flags.IsEnabled("new_feature") {
    // Feature code here
}
```

3. Document feature in roadmap under Phase 2 section

---

## Policy File Updates

If updating `data/policy.v1.json`:

1. Make your changes to the policy file
2. **Temporarily set signature to "PLACEHOLDER"**
3. Run verification script to get new hash:
```bash
go run scripts/verify_policy_hash.go
```
4. Copy the calculated hash
5. Update the `signature` field in policy file
6. Run verification again to confirm:
```bash
go run scripts/verify_policy_hash.go
# Should output: ✅ Policy signature valid
```

**Never commit a policy file with mismatched signature!**

---

## Common Issues

### "Policy signature mismatch" error
- You modified `data/policy.v1.json` without updating the signature
- Run `go run scripts/verify_policy_hash.go` to get correct hash
- Update signature field and try again

### "Feature flag not working"
- Ensure `feature.flags.json` is in the root directory
- Check that flag name matches exactly (case-sensitive)
- Verify flag is set to `"enabled": true` if testing

### Tests failing after policy change
- Update `internal/testdata/` fixtures with new policy structure
- Regenerate test data if necessary
- Check if policy changes broke assumptions in test code

---

## Getting Help

- Read the [Roadmap](plans/roadmap.md) for implementation details
- Check [CLAUDE.md](CLAUDE.md) for architectural rules
- Review [DISCOVERIES_AND_LEARNINGS.md](DISCOVERIES_AND_LEARNINGS.md) for research context
- Ask questions in GitHub Issues before implementing uncertain features

---

## Code of Conduct

### Be Respectful
- Treat all contributors with respect
- Provide constructive feedback
- Assume good intentions

### Stay Focused
- This project has a narrow scope by design
- Resist the urge to add "just one more feature"
- Remember: previous version failed due to feature bloat

### Quality Over Speed
- Take time to write tests
- Don't rush changes that affect behavioral guardrails
- Security and correctness are non-negotiable

---

**Thank you for contributing to TF-Engine 2.0!**

This project succeeds by staying disciplined and focused. Your adherence to these guidelines helps maintain that discipline.
