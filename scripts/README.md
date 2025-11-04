# Scripts Directory

Python scripts for policy management and data generation.

---

## Available Scripts

### Bearish Screener Scripts

#### `add_bearish_screeners_FIXED.py` ✅ **Use This**

**Purpose:** Generates bearish (short-side) Finviz screeners for all compatible sectors.

**What it does:**
1. Removes conflicting positive fundamental filters from bearish universe screeners
2. Skips bearish screeners for incompatible sectors (Utilities, Energy)
3. Adds volume confirmation to breakdown screeners
4. Creates 4 bearish screener types per compatible sector

**Usage:**
```bash
python scripts/add_bearish_screeners_FIXED.py
```

**Output:**
- Updates `data/policy.v1.json` with fixed bearish screeners
- Creates 32 bearish screeners (4 types × 8 sectors)

**Bearish Screener Types:**
- **universe_bearish**: Stocks in downtrends (below SMA200)
- **bounce_bearish**: Bear flag setups (bounces in downtrends, RSI > 60)
- **breakdown_bearish**: 52-week lows with volume confirmation
- **death_cross_bearish**: SMA50 crossed below SMA200

**Sectors with bearish screeners:**
- Healthcare, Technology, Consumer Discretionary, Industrials
- Communication Services, Consumer Defensive, Financials, Real Estate

**Sectors WITHOUT bearish screeners:**
- Utilities (0% trend-following success)
- Energy (mean-reverting, whipsaws)

**Documentation:**
- See [BEARISH_SCREENER_EVALUATION.md](../BEARISH_SCREENER_EVALUATION.md) for analysis
- See [BEARISH_SCREENERS_FIXED.md](../BEARISH_SCREENERS_FIXED.md) for fixes applied

---

#### `add_bearish_screeners.py` ❌ **Deprecated**

**Status:** Original version with bugs (DO NOT USE)

**Problems:**
- Includes conflicting positive fundamental filters in bearish universe
- Creates bearish screeners for incompatible sectors (Utilities, Energy)
- No volume confirmation on breakdowns

**Replaced by:** `add_bearish_screeners_FIXED.py`

---

## Future Scripts

As new policy management needs arise, add scripts here:

- `validate_policy.py` - Policy JSON validation
- `generate_sample_trades.py` - Sample data for testing
- `export_screeners.py` - Export screeners to CSV/Excel
- `backtest_data_parser.py` - Parse backtest results

---

## Script Guidelines

When adding new scripts:

1. **Document purpose** in this README
2. **Handle encoding properly** (avoid Windows emoji issues)
3. **Validate inputs** before modifying policy.v1.json
4. **Create backups** of policy.v1.json before making changes
5. **Use descriptive output** so users know what happened
6. **Test on sample data first** before production policy files

---

## Policy File Location

Scripts should read/write:
- **Source:** `data/policy.v1.json`
- **Build copy:** `dist/policy.v1.json` (synced automatically by build scripts)

**Important:** Always edit `data/policy.v1.json` (source), not `dist/policy.v1.json` (copy).

---

## Running Scripts

### From Project Root:
```bash
python scripts/script_name.py
```

### From Scripts Directory:
```bash
cd scripts
python script_name.py
```

Scripts use relative paths from project root, so running from root is recommended.

---

## Dependencies

Current scripts require:
- Python 3.7+
- No external dependencies (uses standard library only)

If future scripts need packages:
```bash
pip install -r requirements.txt
```

---

**Last Updated:** November 4, 2025
**Scripts:** 2 (1 active, 1 deprecated)
