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
