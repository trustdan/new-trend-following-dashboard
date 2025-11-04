Below is a single, detailed **markdown** document you can drop into your repo (or the `Readme` sheet) to guide how the dashboard enforces *“fewer gates, stronger rules”*—i.e., curb impulsivity without blocking legitimate trend trades. I’ve also included **Gherkin scenarios** so you can treat the workbook like a spec you can validate by inspection.

------

# Trend‑Following Without Impulsivity

**An Excel dashboard for Seykota/Turtle‑style options trading**

> **Intent:** trade obvious trends (Donchian breakout, add by N, exit mechanically) while preventing impulse entries. We keep only a handful of *required* gates; everything else is an *optional* “quality” cushion that nudges, not blocks.

------

## 1) Design principles

- **Trade the tide, not the splash.** A breakout and a mechanical exit are the core.
- **Friction where it matters.** Hard gates for *signal*, *risk/size*, *liquidity/DTE*, *exit*, and *behavior*.
- **Nudge for better trades.** Optional quality items (regime check, no chase, earnings blackout) affect a **score**, not permission.
- **Immediate feedback.** A large **3‑state banner** (RED / YELLOW / GREEN) updates as you tick boxes or change inputs.
- **Journal while you decide.** One click logs the decision and the plan (N, risk%, units, add prices, contracts).
- **Calendar awareness.** Added trades appear across a **rolling 10‑week** sector calendar (2 back + 8 forward) so you can see basket diversification and crowding.

------

## 2) What’s in the workbook

**Sheets**

- **Trades**: `Symbol, Sector, StartDate, EndDate, Active, Strategy, Notes`.
- **Calendar**: sector × week grid (Mon‑Sun ranges) for 2 weeks back + 8 forward.
- **Summary**: active‑count and sector‑coverage charts per week.
- **Checklist**: ticket, required/optional items, banner, contract helper, presets, buttons.
- **Decisions**: log created by “Save Decision”.
- **Lists** (hidden): sector values for validation (edit to customize).

**Macros (by module)**

- **CalendarModule**:
   `ForceRepairAndRefresh`, `RefreshCalendar`, `StartOfWeek`.
- **ChecklistModule**:
  - Build & layout: `BuildInteractiveChecklist`, `ApplyChecklistFixes`, `FixChecklistUI`
  - Actions: `ResetChecklist`, `AddToTrades`, `RecalcChecklist`, `SaveDecision`, `ApplyPreset`
  - Live UI: `Checklist_CheckBoxChanged`, `RefreshBanner`
- **Worksheet (Checklist)** events:
   `Worksheet_Change` and `Worksheet_Calculate` → automatically recolor banner on any checkbox/input change.

> **No Python required.** All interactivity is VBA + worksheet formulas.

------

## 3) Mapping from strategy to checklist

| Strategy dial (PineScript)                    | Excel input / behavior                                       |
| --------------------------------------------- | ------------------------------------------------------------ |
| `entryLen = 55`, `exitLen = 10`, `nLen = 20`  | Ticket fields record **Preset** (20/10 or 55/10) as a journal lens; checklist requires *breakout* and *10‑bar or 2×N exit*. |
| `stopN = 2`, `addStepN = 0.5`, `maxUnits = 4` | Inputs: **Stop multiple (N)**, **Add step (N)**, **Max units**. Helper computes **Add1/Add2/Add3 prices** = `Entry + k × AddStep × N`. |
| `riskPct`                                     | **Risk % per unit**; per‑unit risk $ = `Account × Risk%`.    |
| `N = ATR(20)`                                 | **N (ATR) at entry** field (you can paste N or compute elsewhere and type it). |
| Market regime toggle                          | Optional quality item (score only, does not block).          |
| Time exit (Close/Roll)                        | Optional quality/journal note; exits enforced by 10‑bar/2×N rule. |
| Liquidity / slippage                          | **Required**: bid‑ask < 10% of mid; OI > 100.                |

------

## 4) Minimal “required” gates (hard stops)

- **Signal (SIG_REQ):** 55‑bar breakout (long > 55‑high / short < 55‑low).
- **Risk/Size (RISK_REQ):** per‑unit risk = `% of equity` using **2×N** stop; pyramids **every 0.5×N** to **Max units**.
- **Options (OPT_REQ):** **60–90 DTE**, **roll/close ~21 DTE**, liquidity **required** (bid/ask < 10% mid; OI > 100).
- **Exits (EXIT_REQ):** exit by **10‑bar opposite Donchian OR closer of 2×N**.
- **Behavior (BEHAV_REQ):** 2‑minute **cool‑off** + **no intraday overrides**.

If any required gate is unchecked → **RED (Do Not Trade)**.

------

## 5) Optional “quality” nudges (score, not permission)

- **Regime OK** (e.g., SPY > 200SMA for longs).
- **No chase** (> 2N above 20‑EMA at entry).
- **Earnings blackout** when long premium.
- **Journal note**: *why this unit/now*, *profit‑take plan for debit verticals*, etc.

Each optional box adds **1 point**. **Quality threshold** defaults to **3.0** (cell `I6`).
 Required pass + theta OK + score ≥ threshold → **GREEN (OK TO TRADE)**.
 Required pass + (score < threshold or theta fails) → **YELLOW (Caution)**.

------

## 6) Contracts helper (for quick sizing)

Inputs: **Account size**, **N**, **StopN**, **Add step N**, **Risk %/unit**, **Max units**, **Underlier @ entry**, **Option delta**, **Debit per vertical**, **Structure**.
 Computed:

- **Per‑unit risk $** = `Account × Risk%`.
- **Single‑call contracts** ≈ `unitRisk$ / (delta × (2×N) × 100)`.
- **Debit‑vertical spreads** ≈ `unitRisk$ / debit`.
- **Add prices**: `Entry + (k × AddStep × N)` for k = 1..(MaxUnits − 1)`.

> These are **planning** numbers; execution is still discretionary within the rules.

------

## 7) Workflow (one pass)

1. Fill **ticket** (Symbol, Sector, DTE, Strategy, N, risk %, etc.).
2. Tick required gates as you verify them.
3. Watch **banner** move RED → YELLOW → GREEN as inputs/boxes change.
4. Click **Add To Trades** to push the plan to the Calendar’s 10‑week view.
5. Click **Save Decision** to log everything to **Decisions** (time‑stamped).

------

## 8) Gherkin scenarios (spec you can check by inspection)

### Feature: Banner reflects readiness to trade

```gherkin
Rule: RED if any required group fails
  Scenario: Missing liquidity
    Given "Liquidity OK (bid-ask <10% mid; OI >100)" is unchecked
    When the user checks all other required items
    Then cell J5 shows text "DO NOT TRADE"
    And J5 background is red

Rule: YELLOW if required pass but quality or theta fails
  Scenario: All required pass, score below threshold
    Given all required items are checked
    And I5 (Quality Score) = 2.0
    And I6 (Threshold) = 3.0
    Then J5 shows "CAUTION"
    And J5 background is yellow

  Scenario: Theta budget missing or failing
    Given all required items are checked
    And C4 (Account Size) is blank
    Then I24 (Theta budget OK) evaluates to TRUE
    And J5 may still be "CAUTION" if score < threshold

Rule: GREEN only if all gates pass
  Scenario: Good quality and theta
    Given all required items are checked
    And I5 >= I6
    And I24 = TRUE
    Then J5 shows "OK TO TRADE"
    And J5 background is green
```

### Feature: Contracts helper computes counts

```gherkin
Background:
  Given C4 = 100000 (Account size)
  And C11 = 0.0075 (Risk % per unit = 0.75%)
  And C8 = 3.0 (N at entry)
  And C9 = 2.0 (Stop multiple N)
  And C14 = 0.65 (Delta single call)
  And C15 = 5.00 (Debit per vertical)

Scenario: Single-call contracts
  When the sheet recalculates
  Then I27 (unit risk $) = 750.00
  And I30 (call contracts) = floor( 750 / (0.65 * (2*3.0) * 100) ) >= 1

Scenario: Debit vertical spreads
  When the sheet recalculates
  Then I31 (vertical spreads) = floor( 750 / 5.00 ) = 150
```

### Feature: Add‑on levels displayed

```gherkin
Scenario: Compute add levels from N and AddStep
  Given C13 = 100.00 (Underlier @ entry)
  And C8 = 3.0 (N)
  And C10 = 0.5 (Add every 0.5N)
  And C12 = 4 (Max units)
  Then I32 (Add1) = 100 + (0.5*3) = 101.5
  And  I33 (Add2) = 100 + (1.0*3) = 103.0
  And  I34 (Add3) = 100 + (1.5*3) = 104.5
```

### Feature: Preset journaling

```gherkin
Scenario: System-2 selected and applied
  Given C18 = "System-2 (55/10)"
  When the user clicks Apply Preset
  Then I36 = "System-2"
  And  I37 = "55/10"
  And Save Decision will include "Preset" = "System-2 (55/10)"
```

### Feature: Add To Trades integrates with Calendar

```gherkin
Scenario: Row appears in calendar window
  Given Checklist banner J5 is "OK TO TRADE"
  And user clicks "Add To Trades"
  Then a new row is appended to the Trades sheet with StartDate = Today and EndDate = Today + DTE
  And the symbol appears in the sector/week cell for any week overlapping [StartDate, EndDate]
```

### Feature: Anti‑impulsivity without paralysis

```gherkin
Scenario: A valid breakout with minimal quality passes
  Given Signal, Risk/Size, Liquidity/DTE, Exit, and Behavior are all checked
  And only two optional quality items are checked
  And I6 (Threshold) = 3.0
  Then the banner is "CAUTION"
  When the user checks one more optional item or raises the plan quality
  Then the banner turns "OK TO TRADE"
```

### Feature: Reset and Save Decision

```gherkin
Scenario: Reset clears all Done checkboxes
  When the user clicks Reset
  Then all linked cells in column F become FALSE
  And J5 becomes "DO NOT TRADE" (red) unless required are re-checked

Scenario: Save Decision logs everything
  When the user clicks Save Decision
  Then a row is appended to "Decisions" with timestamp and fields including:
       Symbol, Sector, Strategy, N, StopN, AddStepN, Risk%/unit, MaxUnits, DTE,
       Structure, Delta, Debit, Entry price, CallContracts, VerticalSpreads, Preset, Banner Status
```

------

## 9) Guardrails and failure states

- **Liquidity rule** is *required*—prevents most poor fills.
- **Theta rule** is **lenient**: if Account Size or Max Theta are blank, it evaluates to OK (keeps friction low while you’re configuring).
- **Button placement** uses a dedicated UI column (**L**) to avoid overlapping the grid.
- **Events** are re‑enabled by `RecalcChecklist` if something feels “stuck”.

------

## 10) Customization knobs

- **Threshold** (`I6`): raise from 3 → 4–6 to demand more quality boxes for GREEN.
- **Risk % per unit** (`C11`): 0.5–1.0% typical for swing.
- **Add step N** (`C10`): 0.5N standard; smaller → more frequent adds.
- **Max units** (`C12`): cap to 3 if you want a gentler book.
- **Sectors**: edit the **Lists** sheet.

------

## 11) What not to add (to avoid paralysis)

- No subjective pattern checks in required gates.
- No multi‑factor overlays that invert the breakout logic.
- No rigid earnings ban for all strategies; keep it as an optional quality item.

------

## 12) Definition of Done (for a **trade ticket**)

- Required gates checked **and** banner **GREEN**.
- Contracts helper reviewed (counts make sense vs liquidity).
- **Add To Trades** pressed (shows on calendar), and **Save Decision** logged.

------

### Appendix: Column/cell cheat‑sheet (Checklist)

- **Inputs:** `C2:C7` (ticket text), `C8:N` (N/risk/contract params, presets).
- **Helpers:** `H/I 27–34` (risk $, contracts, add prices), `I36–I37` (preset journal).
- **Quality threshold:** `I6`; **Score:** `I5`.
- **Banner text:** `J5` (color set in code).

------

If you’d like, I can also deliver a **printable one‑pager** (A4/Letter) with the required gates and a short checklist you can keep at your desk; it will mirror the workbook’s logic and wording.