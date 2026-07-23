# Unit Specification Decision

This maintainer reference records the unit-code, Tech-Level, and combined
unit-specification invariants for the EC v7 domain. It closes acceptance
criterion 1 of GitHub issue #25 ("Unit code, Tech Level, and combined unit
specification invariants are defined.").

The canonical-reference/grammar reconciliation and the remaining deferred
questions are tracked under the other acceptance criteria of #25 and are out of
scope here.

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
  explicit `-0` denote the same spec `(Code, 0)`. Use as an inventory key is
  specified in §6.
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

## 5. Worked examples

These examples exercise the rules above against the five required kinds of unit
— production (`FACT`, `MINE`), population, natural resources, and `TRNS` — using
the TL domains from the Unit-Definition Catalog in
[`doc/ec-unit-tables.md`](ec-unit-tables.md). `FACT`, `MINE`, and `TRNS` are
TL-bearing with the domain `{1..10}`; `FUEL`, `METL`, and `NMTL` are TL-0; `PRO`,
`SOL`, `UEM`, and `USK` are `PopulationClass`, not `UnitSpec` (US-3).

### 5.1 Parsing

| Input     | Result                        | Rule  | Notes                                            |
|-----------|-------------------------------|-------|--------------------------------------------------|
| `FACT-3`  | `UnitSpec{FACT, 3}`           | P-1   | TL-bearing production unit; `(FACT, 3)` in domain |
| `fact-3`  | `UnitSpec{FACT, 3}`           | P-1   | Case-insensitive match; same spec as `FACT-3`     |
| `FACT`    | rejected                      | P-4   | TL required; `FACT`/`FARM`/`MINE` never bare      |
| `MINE-5`  | `UnitSpec{MINE, 5}`           | P-1   | TL-bearing production unit                        |
| `MINE-0`  | rejected                      | P-1   | `0` outside domain `{1..10}` (US-1)               |
| `PRO`     | `PopulationClass` (`PRO`)     | P-3   | Population is non-spec inventory (US-3)           |
| `SOL`     | `PopulationClass` (`SOL`)     | P-3   | Population is non-spec inventory (US-3)           |
| `METL`    | `UnitSpec{METL, 0}`           | P-2   | Bare TL-0 resource resolves to `(METL, 0)`        |
| `METL-0`  | `UnitSpec{METL, 0}`           | P-1   | Explicit `-0` denotes the same spec as bare `METL` |
| `FUEL`    | `UnitSpec{FUEL, 0}`           | P-2   | Bare TL-0 resource                                |
| `TRNS-4`  | `UnitSpec{TRNS, 4}`           | P-1   | TL-bearing transport unit                         |
| `TRNS-0`  | rejected                      | P-1   | `0` outside domain `{1..10}` (US-1)               |
| `TRNS-11` | rejected                      | P-1   | `11` outside the `0`–`10` Tech-Level range (TL-1) |

### 5.2 Formatting

| Spec / kind                | Output    | Rule | Notes                                    |
|----------------------------|-----------|------|------------------------------------------|
| `UnitSpec{FACT, 3}`        | `FACT-3`  | F-1  | TL-bearing spec                          |
| `UnitSpec{MINE, 5}`        | `MINE-5`  | F-1  | TL-bearing spec                          |
| `UnitSpec{TRNS, 4}`        | `TRNS-4`  | F-1  | TL-bearing spec                          |
| `UnitSpec{METL, 0}`        | `METL`    | F-2  | TL-0 spec emits bare form, omitting `-0` |
| `UnitSpec{FUEL, 0}`        | `FUEL`    | F-2  | TL-0 spec emits bare form                |
| `PopulationClass` (`PRO`)  | `PRO`     | F-3  | Non-spec kind                            |

## 6. Equality and inventory key

An inventory holds quantities against identities. This section specifies when two
identities are equal and how each kind of identity is keyed. It complements the
[numeric-types decision](ec-numeric-types.md), which defines `InventoryUnit`
(a `UnitSpec` plus a `Quantity`) and places per-unit mass, volume, and assembly
in `UnitDefinition` rather than in the identity.

- **EQ-1 (equality)** Two `UnitSpec` values are equal if and only if their `Code`
  and `TL` are both equal (US-2). `UnitCode` is a case-sensitive uppercase
  mnemonic (UC-3) and `TL` is an integer, so `UnitSpec` is a small comparable
  value and equality is exact field equality — no normalization at compare time.
- **IK-1 (spec key)** For TL-bearing and TL-0 units, the `UnitSpec` is itself the
  inventory key. It is directly usable as a map key; no separate derived key is
  introduced.
- **IK-2 (canonical key)** Keys are canonical specs. A TL-0 unit written bare and
  written with an explicit `-0` resolve to the same spec `(Code, 0)` (P-2, US-2)
  and therefore to a single key: `METL` and `METL-0` are one inventory line, not
  two. Text is normalized to a spec (§4) before it is used as a key.
- **IK-3 (non-spec keys)** Kinds that are not `UnitSpec` (US-3) are keyed by their
  own identity, never by a fabricated `(Code, TL)`:
  - `PopulationClass` units (Population and Cadre) are keyed by their class code
    (`PRO`, `SOL`, `WRKR`, …).
  - Research Points (`RP`) are a single scalar balance, not a keyed inventory
    line.
- **IK-4 (state is not identity)** An entity may hold both stored and assembled
  quantities of the same `UnitSpec`. The stored/assembled distinction is a
  separate dimension of the inventory line, not part of `UnitSpec` equality or
  the key; its final shape is owned by issue #26.

## 7. Decision status and deferred questions

This document is the recorded maintainer decision for the unit specification. The
decisions it fixes are:

- unit-code, Tech-Level, and combined-spec invariants (§1–§3);
- parsing and formatting, including case handling and the bare-form rules (§4);
- worked examples for production, population, resource, and transport units (§5);
  and
- equality and inventory-key semantics (§6).

The following questions are deferred. None of them blocks the invariants above;
each is owned elsewhere and is listed here so the open surface is recorded in one
place:

- **Bare-form grammar cleanup.** `FACT`, `FARM`, and `MINE` always require the
  `-TL` suffix (P-4); the Lemon grammar and `ec-modern-orders.md` still admit bare
  forms with a separate integer Tech Level and are the tracked follow-up (AC 4 of
  #25).
- **Stored vs. assembled inventory shape.** The stored/assembled dimension of an
  inventory line and its transition rules are owned by issue #26 (IK-4).
- **Stored vs. operational volume.** Whether a unit occupies different volume when
  stored versus assembled is open in the canonical facts and is isolated in issue
  #26; it does not affect unit identity.
- **`PROB` unit costs.** Production inputs, mass, and volume for `PROB` are not yet
  specified (`TBD` in [`ec-mass-volume.md`](ec-mass-volume.md) and the
  Unit-Definition Catalog in [`ec-unit-tables.md`](ec-unit-tables.md)).
- **`PRTO` production inputs.** Metal and non-metal inputs for `PRTO` are not yet
  specified (`TBD` in the same tables); its mass and volume are known.
