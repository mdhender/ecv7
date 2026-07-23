# Numeric Types and Unit-Value Ownership Decision

This maintainer reference records the numeric representation and formatting
decisions chosen for the first EC v7 `NAME` and `TRANSFER` slice. It resolves
GitHub issue #24. Rule-specific rounding, allocation, and converted-value
validation remain with the game rules that perform those operations.

## Domain numeric types

Semantically different dimensions use distinct named Go types at domain
boundaries. Whole-unit quantities and aggregates use explicit `int64` backing
rather than platform-dependent `int`.

```go
type Quantity int64
type MassUnit int64
type VolumeUnit int64
type TechLevel int64
```

`Quantity` represents a count of inventory units. It is non-negative in game
state, although a specific operation may impose a higher minimum. In particular,
a `TRANSFER` order requests at least one whole unit.

`MassUnit` and `VolumeUnit` are scaled integer values with a scale of 1,000:

| Canonical value | Stored value |
|-----------------|-------------:|
| `0.02`          |           20 |
| `0.04`          |           40 |
| `0.1`           |          100 |
| `0.25`          |          250 |
| `0.5`           |          500 |
| `1`             |        1,000 |

The factory production input table expresses its fractional mass and volume
facts in these scaled units. The domain does not introduce separate
`MassFactor` or `VolumeFactor` types.

`TechLevel` is a named type because it has domain semantics and validation
distinct from other integers. Its canonical range is 0 through 10. The rules
for which unit codes require or permit a particular Tech Level belong to the
unit-specification decision in issue #25.

The theoretical maximum FUEL quantity in a cluster containing 100 stelliums,
six systems per stellium, ten planets per system, and 40 maximum-size deposits
per planet is:

```text
99,999,999 × 40 × 10 × 6 × 100 = 23,999,999,760,000
```

This exceeds `int32` but is represented exactly by `int64`. It is also below
the largest consecutive integer represented exactly by `float64`, `2^53`.

## Inventory units and labor

FUEL is not a separate numeric type. Like every inventory unit, it is identified
by a typed unit code and has a `Quantity`. Unit code plus Tech Level forms the
unit identity; issue #25 owns its final invariants and textual notation.

Conceptually:

```go
type UnitCode string

type UnitSpec struct {
	Code UnitCode
	TL   TechLevel
}

type InventoryUnit struct {
	Spec     UnitSpec
	Quantity Quantity
}
```

`UnitSpec` remains a small, comparable identity value suitable for use as an
inventory key. `InventoryUnit` carries the quantity held for that identity.
`UnitSpec` does not carry mass, volume, or current operational state.

Labor is a composition of population quantities rather than a scalar value:

```go
type Labor struct {
	PRO Quantity
	USK Quantity
	SOL Quantity
}
```

`SOL` means Soldiers. `S.O.L.` means Standard of Living.

## Unit definitions and inventory state

Static per-unit facts belong to the ruleset's unit definition, not to
`UnitSpec`. These include per-unit mass, stored volume, and whether the unit
requires assembly before operation. A representative shape is:

```go
type UnitDefinition struct {
	Spec             UnitSpec
	Mass             MassUnit
	StoredVolume     VolumeUnit
	RequiresAssembly bool
}
```

Current stored or assembled quantities belong to inventory state. An entity may
hold stored and assembled quantities of the same `UnitSpec`, so operational
state is not a boolean on either `UnitSpec` or the inventory as a whole. Issue
#26 owns the final inventory shape and stored/assembled transition invariants.

Occupied volume is derived from inventory quantity and its `UnitDefinition`; it
is not part of unit identity. Persistence may store source state needed to
derive the value but should not make a redundant cached total authoritative.

| Concern                              | Owner                         |
|--------------------------------------|-------------------------------|
| Unit code and Tech Level identity    | `UnitSpec`                    |
| Per-unit mass and stored volume      | `UnitDefinition`              |
| Whether assembly applies             | `UnitDefinition`              |
| Stored and assembled quantities      | Inventory state               |
| Total occupied mass or volume        | Derived domain calculation    |

The distinction between stored and operational volume remains open in the
canonical facts and is isolated in issue #26. This decision does not invent an
operational-volume rule.

## Other canonical numeric values

- S.O.L. and Rations are non-negative `float64` values with four canonical
  decimal places and no defined upper bound. Practical S.O.L. values range from
  `0.0625` through `1.2500`; practical Rations values range from `0.0625`
  through `1.0000`.
- `Yield_Pct` is an integer percentage point in the inclusive range 1 through
  99 and is displayed as an integer.
- Prices and currency are whole units.
- The marketplace purchase premium is an integer percentage point in the
  inclusive range 1 through 5.
- Production retains only quarter-unit remainders: `0`, `0.25`, `0.5`, and
  `0.75`, represented internally as 0, 1, 2, and 3. Detailed production
  behavior is outside the first slice.

Additional named types should be introduced with the game systems that require
them rather than speculated into the first domain slice.

## Calculation representation

Authoritative calculations use Go `float64` where fractional intermediates are
required. The engine typically rounds a calculated result before storing it or
applying it to other units. The applicable game rule defines the exact rounding
boundary, direction, converted-value range validation, and failure behavior.

Non-finite values and invalid conversions must not be committed accidentally to
game state. Detailed checks are defined with the rules that can produce them.

Known arithmetic boundaries are:

- MINE output rounds down to a whole unit after applying `Yield_Pct`.
- A production pipeline emits whole units and may retain only its encoded
  quarter-unit remainder.
- Combat accumulates fractional damage during a round and rounds before applying
  it.
- Transfer quantities and transport limits are whole units.
- Deterministic engine allocation may partially fulfill or reject an
  over-extended order according to the applicable game rule.

The following recommendations are retained for later game-rule discussions;
they are not adopted rules:

| Calculation                                | Recommendation |
|--------------------------------------------|----------------|
| Mine output                                | Down           |
| Factory inputs consumed                    | Retain the permitted quarter-unit pipeline remainder |
| Factory output exposed to inventory        | Down to whole units and retain the permitted remainder |
| Births added                               | Down           |
| Deaths applied                             | Up             |
| Player-caused damage                       | Down           |
| Damage against the player                  | Up             |
| Defender-favored tie or ambiguous result   | Define explicitly in the combat rules |
| Transfer capacity                          | Down to a whole transferable quantity |
| Fuel charged                               | Define explicitly in the transport rules |

“Round against the player” is guidance rather than a complete algorithm.
Neutral and multi-player calculations require explicit rule-specific behavior.

## Parsing and report display

The order grammar's lexical picture for a non-negative value with up to four
decimal places is:

```text
[0-9]+(\.[0-9]{1,4})?
```

It accepts `0`, `0.0`, `0.0125`, and `1.25`. It rejects `0.00001`, `.5`, `1.`,
negative values, and exponent notation.

Reports display Rations and S.O.L. as `N.NNNN`, always including four digits
after the decimal point. No other currently known fractional value is displayed
on a report.

An individual stored item's volume is rounded up for display. A total stored
volume sums the individual scaled values before rounding the total for display;
it is not the sum of the separately rounded display values.

## Deferred rule decisions

The following decisions do not block the first `NAME` and `TRANSFER` slice:

- exact rounding boundaries and directions beyond the known boundaries above;
- range validation and failure behavior for each calculated conversion;
- resource allocation and partial-fulfillment behavior for each order type;
- production-pipeline behavior beyond the quarter-unit representation;
- combat damage formulas and application rules; and
- operational volume and stored/assembled transition rules from issue #26.
