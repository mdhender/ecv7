# Unit Specification Decision

This maintainer reference records the unit-code, Tech-Level, and combined
unit-specification invariants for the EC v7 domain. It closes acceptance
criterion 1 of GitHub issue #25 ("Unit code, Tech Level, and combined unit
specification invariants are defined.").

Parsing and formatting rules, the canonical-reference/grammar reconciliation,
inventory-key details, and the remaining deferred questions are tracked under
the other acceptance criteria of #25 and are out of scope here.

The concrete set of codes and their Tech-Level domains is enumerated in
[`doc/ec-unit-tables.md`](ec-unit-tables.md): the Units and Classes table lists
the codes, and the Unit-Definition Catalog lists every valid `(code, TL)` pair.
`TechLevel` is the numeric type defined in the [numeric-types
decision](ec-numeric-types.md).

## 1. Unit code

A **unit code** identifies a kind of unit. Codes are short uppercase mnemonics
drawn from a fixed, closed set.

- **UC-1** A unit code is a member of the defined set (the Units and Classes
  table). Any other value is not a unit code.
- **UC-2** Each code belongs to exactly one Class (for example `Ammunition`,
  `Infrastructure`, `Resource`, `Structural`).
- **UC-3** Codes are compared exactly. They are case-sensitive uppercase
  mnemonics with no aliases.

## 2. Tech Level

A **Tech Level** qualifies a unit code. Each code has a **TL domain**: the set of
Tech Levels at which the code names a real unit.

- **TL-1** A Tech Level is an integer in the inclusive range `0`–`10`.
- **TL-2** Each code's TL domain is exactly one of:
  - `{0}` — a *TL-0 unit* (Natural Resources, `FOOD`, and the flat commodities
    `CSGD` and `CSUP`);
  - a non-empty subset of `{1..10}` — a *TL-bearing unit* (production, weapons,
    structural, and support units); or
  - empty — the code does not form a unit specification (see §3).
- **TL-3 (exclusivity)** Within a ruleset, no code has both a TL-0 member and a
  positive-TL member. A code is either TL-0-only or TL-bearing, never both.

## 3. Combined unit specification

A **unit specification** pairs a unit code with a Tech Level.

```go
type UnitSpec struct {
    Code UnitCode  // a defined unit code (§1)
    TL   TechLevel // a Tech Level in the code's TL domain (§2)
}
```

- **US-1 (validity)** A `UnitSpec` is valid if and only if `Code` is a defined
  code and `TL` is a member of that code's TL domain — equivalently, if and only
  if the `(Code, TL)` pair appears in the Unit-Definition Catalog. Unknown codes
  and out-of-domain Tech Levels (for example `TRNS-0`, `TRNS-11`, `STRC-0`) are
  invalid.
- **US-2 (equality)** Two `UnitSpec` values are equal if and only if their
  `Code` and `TL` are both equal. A TL-0 unit written bare and written with an
  explicit `-0` denote the same spec `(Code, 0)`. Inventory-key semantics beyond
  value equality are specified separately.
- **US-3 (non-spec inventory)** Kinds that never carry a Tech Level as a
  specification are not `UnitSpec` values:
  - Population and Cadre units are inventory of `PopulationClass`, not
    `UnitSpec`.
  - Research Points (`RP`) are a non-physical bookkeeping balance, not a
    `UnitSpec`.

## 4. Parsing and formatting

Text is parsed to a `UnitSpec` (or to a non-spec kind), and specs are formatted
back to a single canonical text form. Parsing validates against the authoritative
unit set — the codes and TL domains referenced in the introduction.

**Case.** Codes are accepted in any case and matched case-insensitively. The
canonical form emitted is always uppercase (UC-3).

### 4.1 Parsing (text → spec or kind)

- **P-1 (suffixed)** `CODE-N`, with `N` an integer `0`–`10`, resolves by exact
  `(CODE, N)` membership in the catalog and is valid if and only if the pair
  exists. Unknown codes and out-of-domain levels (`TRNS-0`, `TRNS-11`, `STRC-0`)
  are rejected (US-1).
- **P-2 (bare TL-0)** A bare `CODE` for a TL-0 unit resolves to `(CODE, 0)`.
- **P-3 (bare non-TL)** A bare `CODE` for a non-spec kind (Population, Cadre,
  `RP`) resolves to that kind, not a `UnitSpec` (US-3).
- **P-4 (bare TL-bearing rejected)** A bare `CODE` for a TL-bearing unit does not
  resolve: the Tech Level is required. In particular, `FACT`, `FARM`, and `MINE`
  require the `-TL` suffix in every position; forms that write them bare with a
  separate integer Tech Level are out of date and are reconciled under the
  grammar cleanup (AC 4 of #25).

### 4.2 Formatting (spec or kind → canonical text)

- **F-1** A TL-bearing spec formats as `CODE-N` (for example `FACT-3`).
- **F-2** A TL-0 spec formats as bare `CODE`, omitting the `-0`. `CODE` and
  `CODE-0` therefore parse to the same spec `(CODE, 0)` (US-2), and the canonical
  output is the bare form.
- **F-3** A non-spec kind formats as bare `CODE`.

> The concrete authoritative code set — including the `MSUP` → `CSUP` rename and
> the other maintainer updates from the consolidated tables — is being folded
> into the canonical reference. Until that lands, these rules stand but the
> enumerated examples (AC 3 of #25) are deferred.
