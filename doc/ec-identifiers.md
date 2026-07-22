# Identity and Identifier Decision

This maintainer reference records the identity model chosen for the EC v7 domain,
engine, persistence adapters, orders, and reports. It resolves GitHub issue #23.

## Identity classes

EC uses two identity classes:

1. **Domain keys** are intrinsic game values used by deterministic generators and
   the game engine.
2. **Store IDs** are opaque SQLite integer primary keys used by persistence,
   orders, and reports.

The domain model does not contain Store IDs. Persistence and reporting models
contain both forms when they must translate between them or present them to a
player.

Every store contains at most one game. Domain keys and Store IDs are therefore
meaningful only within that game store. Importing orders created for one game
into another game is invalid and unsupported.

## Domain keys

Domain keys are immutable. Equality compares all components of the same key
type. Values of different key types are never interchangeable.

| Object   | Domain key                          | Scope                               |
|----------|-------------------------------------|-------------------------------------|
| Game     | The current engine game             | One game per engine/store           |
| Player   | `PlayerNo`                          | Game                                |
| Faction  | `FactionNo`                         | Game                                |
| Stellium | `(x, y, z)`                         | Game                                |
| System   | `(x, y, z, seq)`                    | Game                                |
| Planet   | `(x, y, z, seq, orbit)`             | Game                                |
| Deposit  | `(x, y, z, seq, orbit, deposit_no)` | Game; `deposit_no` is Planet-scoped |
| Entity   | `(FactionNo, EntityNo)`             | Game; `EntityNo` is Faction-scoped  |

Hierarchical Go values should embed their parent key rather than repeat its
individual fields:

```go
type StelliumKey struct {
	X, Y, Z Coordinate
}

type SystemKey struct {
	Stellium StelliumKey
	Seq      SystemSeq
}

type PlanetKey struct {
	System SystemKey
	Orbit  Orbit
}

type DepositKey struct {
	Planet PlanetKey
	No     DepositNo
}

type EntityKey struct {
	Faction FactionNo
	Entity  EntityNo
}
```

This nesting is a domain representation choice. SQLite may flatten the
components into columns to enforce composite uniqueness.

Coordinates and sequence values are part of world-object identity, not merely
display labels. Mutable object names never participate in identity.

## Number allocation

`PlayerNo`, `FactionNo`, `EntityNo`, and `DepositNo` are positive domain values.
Zero is invalid or unset.

### Players and factions

Adding a player assigns both a unique `PlayerNo` and a unique `FactionNo`.
They are independent concepts and are not required to have the same numeric
value. Both are game-scoped, stable, immutable, and never reused. A retired
player or faction does not return its number to either allocator.

Player and Faction numbers are domain values, not aliases for `PlayerID` or
`FactionID`. Their allocation state must survive persistence and must not be
reconstructed from SQLite primary keys.

### Entities

`EntityNo` is scoped to its owning Faction. The single-threaded game engine
allocates Entity numbers sequentially from persistent per-Faction allocation
state. It never deletes, renumbers, or reuses an allocated Entity number.

Entity creation and the updated next-number state are committed atomically with
the turn. A rolled-back turn does not create an Entity and therefore does not
consume a committed Entity number. No concurrent Entity-number allocator is
required because one engine serially owns game-state mutation.

### Deposits

`DepositNo` ranges from 1 through 40 and is scoped to its parent Planet. A
Planet's first Deposit is number 1. Deposit numbers are assigned during
deterministic Planet generation and are never renumbered or reused.

## Store IDs

Persisted first-slice objects have distinct Store ID types:

- `GameID`
- `PlayerID`
- `FactionID`
- `StelliumID`
- `SystemID`
- `PlanetID`
- `DepositID`
- `EntityID`

Each is represented as a defined Go type with an `int64` underlying value and
as an SQLite `INTEGER PRIMARY KEY AUTOINCREMENT`. The exact SQLite declaration
is required: `AUTOINCREMENT` prevents reuse of a ROWID that belonged to a
previously committed row. Migrations and maintenance tools must not reset or
rewrite `sqlite_sequence`.

A Store ID is positive. Zero is invalid or unset. Equality is typed integer
equality; IDs belonging to different types cannot be compared as identities.
An ID is unique within its table and store, not across tables or independent
game stores.

The canonical textual form of a Store ID is an unsigned base-10 integer without
leading zeroes. Parsers reject zero, negative values, malformed text, and values
outside SQLite's signed 64-bit range.

Store IDs do not seed game generation, determine game outcomes, or establish
domain equality. Their values can depend on persistence insertion order and are
therefore not intrinsic game state.

## Persistence constraints

The SQLite adapter maintains a one-to-one mapping between each Store ID and its
domain key. In addition to primary keys and foreign keys, the schema enforces
these uniqueness constraints:

| Object   | Unique columns                               |
|----------|----------------------------------------------|
| Player   | `(game_id, player_no)`                       |
| Faction  | `(game_id, faction_no)`                      |
| Stellium | `(game_id, x, y, z)`                         |
| System   | `(game_id, x, y, z, seq)`                    |
| Planet   | `(game_id, x, y, z, seq, orbit)`             |
| Deposit  | `(game_id, x, y, z, seq, orbit, deposit_no)` |
| Entity   | `(game_id, faction_no, entity_no)`           |

The schema also constrains `deposit_no` to 1 through 40 and all domain ordinal
numbers to positive values. The application rejects a Store ID that does not
belong to the current store's game.

## Generator and engine boundary

Generators and the engine use only domain keys:

- deterministic world generation addresses Stellium, Systems, Planets, and
  Deposits by their coordinate tuples;
- player generation uses `PlayerNo` and, when Faction-specific behavior is
  needed, `FactionNo`;
- engine state and changes address Entities by `(FactionNo, EntityNo)`; and
- game-affecting traversal uses an explicit stable order and never relies on Go
  map iteration, SQL row order, or Store ID allocation order.

The persistence adapter translates domain keys to Store IDs after generation or
engine execution. It translates Store IDs to domain keys and loaded domain
values before invoking the engine. Neither generators nor the engine query the
store to allocate or resolve Store IDs.

## Orders

Players use Store IDs from their reports when writing orders. An inbound adapter
resolves every Store ID against the current game before invoking the engine.

For example, given this mapping:

```text
EntityID 481 <-> (FactionNo 7, EntityNo 12)
```

an order whose subject is `ENTITY 481` is loaded and passed to the engine as
Entity key `(7, 12)`. Authorization checks the submitting player's Faction
against that key. A `NAME` mutation and both endpoints of a `TRANSFER` are
represented inside the engine by domain keys, not by SQLite IDs.

An ID from another game is not portable. Orders are valid only for the game
store whose report supplied their IDs.

## Reports

Report generation uses a read model assembled by the persistence/reporting
adapter. The read model includes both the Store ID needed for future orders and
the domain key needed to identify the object in game terms.

Examples:

```text
Stellium 41  (-4, 7, 2)
System 137   (-4, 7, 2, B)
Planet 1937  (-4, 7, 2, B, 3)
Deposit 8861 (-4, 7, 2, B, 3, 1)
Entity 481   (Faction 7, Entity 12) at (-4, 7, 2, B, 3)
```

The engine does not receive Store IDs merely because reports display them.

## Determinism

Domain tuples and numbers are valid deterministic generator-address components
because they are intrinsic, immutable game values. Store IDs are forbidden in
generator addresses because they depend on persistence behavior.

Stable keys do not make unordered iteration deterministic. Code that consumes a
collection in a way that can affect state, random draw assignment, or output
must sort by an explicit domain-key order before iterating. This complements the
addressed PRNG streams and the repository rule forbidding game behavior from
depending on map iteration order.
