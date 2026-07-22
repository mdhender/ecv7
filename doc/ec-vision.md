# Epimethean Challenge Implementation Vision

The project should implement the game from the domain outward. The database
and command-line foundations are useful, but the next step should not be a
complete database schema or a complete game engine. Both depend on game concepts
and invariants that have not yet been represented in code, and parts of the
canonical reference remain intentionally unspecified.

The recommended sequence is:

```text
Resolve the minimum rules needed by the first vertical slice
    -> build a small pure-Go domain model
    -> prove it with a thin in-memory engine slice
    -> persist that slice in SQLite
    -> expand the engine one phase at a time
```

## Phase 1: Resolve decisions needed by the first slice

Do not try to complete the entire canonical specification before implementation.
Resolve only the foundational representation and behavior choices needed for the
first vertical slice:

1. Stable identifiers for games, players, factions, entities, systems, planets,
   and deposits.
2. Whether game quantities are always integral or may be fractional.
3. Fixed-point representations and scales for currency, mass, volume, yield,
   percentages, and fuel.
4. The representation and textual notation of unit code plus Tech Level.
5. The representation and rules for stored versus assembled units.
6. The representation of entity locations.
7. Whether invalid or unsatisfied orders are rejected completely or may be
   partially fulfilled.
8. Deterministic ordering when multiple orders compete for the same resources.
9. Whether canonical rules and unit definitions are versioned per game.
10. How random behavior is seeded, persisted, and reproduced.

Prefer these defaults unless investigation reveals a rule conflict:

- Use typed integer identifiers internally.
- Use scaled integers rather than binary floating point for exact game values.
- Represent a unit as a structured `UnitSpec` containing a code and Tech Level,
  rather than treating `FACT-3` as an opaque string.
- Record an immutable ruleset version on every game.
- Process a turn transactionally while retaining the result of each individual
  order.
- Require every random decision to derive from explicit, persisted state.

Each decision should identify the canonical rules it supports, state the chosen
invariants, include representative examples, and call out questions deferred to
later slices.

## Phase 2: Build a small pure-Go domain kernel

Build a domain package independent of SQLite, command-line parsing, and report
formatting. Initially model only the concepts needed to express and validate the
first vertical slice.

### Identity and world structure

- Game
- Turn number and status
- Player
- Faction
- Species
- Stellium
- System
- Planet
- Deposit
- Location

### Entities and inventory

- Ship
- Colony and colony type
- Entity owner, location, Tech Level, and name
- Unit code and unit Tech Level
- Stored and assembled quantities
- Population ownership where necessary

### Canonical value rules

- Valid Tech Level range
- Valid Habitability range
- Valid colony placement
- One colony of each type per faction per planet
- Unit mass and volume calculations
- Operational versus non-operational units
- Research cost lookup
- Transport throughput and fuel calculations

Prefer pure functions for calculations and validation. For example:

```go
func ResearchCost(targetTL TechLevel) (ResearchPoints, error)
func MassOf(spec UnitSpec, quantity Quantity) Mass
func TransportThroughput(tl TechLevel, quantity Quantity) Mass
func CanEstablishColony(faction Faction, planet Planet, kind ColonyType) error
```

Domain objects must not save themselves, contain SQLite row details, or hold
database connections. Persistence should translate between records and domain
state.

## Phase 3: Prove the model with one vertical engine slice

Implement a small end-to-end turn path rather than expanding the domain model
indefinitely:

1. Create a small deterministic test game.
2. Add two factions and a few entities.
3. Submit `NAME` and `TRANSFER` orders.
4. Parse the orders into typed commands.
5. Validate authorization and game invariants.
6. Execute the orders in their canonical turn phases.
7. Produce per-order results and a minimal turn report.
8. Advance the turn.

`NAME` proves basic subject lookup, authorization, mutation, and reporting.
`TRANSFER` exercises the important foundations: ownership, co-location,
inventory, unit identity and Tech Level, quantities, transport capacity, crew,
fuel, competing orders, and deterministic state changes.

Avoid beginning with production, population, markets, espionage, or combat;
their rules are substantially less complete.

## Phase 4: Persist the proven slice

After the domain and execution shape are demonstrated in memory, add only the
SQLite persistence required by that slice. Likely concepts include:

- Games and turns
- Players, factions, and species
- Stellium, systems, planets, and deposits
- Entities and inventories
- Order submissions, parsed orders, and order results
- Turn reports

The exact tables should follow demonstrated access patterns and transaction
boundaries rather than this conceptual list.

Keep immutable canonical catalogs, formulas, and validation rules in typed and
tested Go code initially. Store dynamic game state in SQLite, along with the
ruleset version required to reproduce that state.

A turn should conceptually:

1. Load its authoritative starting state.
2. Validate submitted orders.
3. Execute stages deterministically.
4. Store resulting state, order outcomes, and reports.
5. Mark the turn complete.

Execution and persistence of the resulting turn should be atomic. A failed turn
must leave the previous state authoritative.

## Phase 5: Expand the engine in dependency order

After the first persistent slice, expand approximately in this order:

1. Naming and basic administration
2. Transfers and inventory
3. Assembly, disassembly, and setup
4. Mining and farming
5. Manufacturing and research
6. Surveys, probes, and movement
7. Population, rationing, and life support
8. Control and victory
9. Markets
10. Espionage and rebellion
11. Combat

Combat should come late because it depends on most other subsystems and its core
resolution formulas remain unspecified.

Represent the canonical 21 processing stages explicitly. Each stage should take
deterministic state and produce state changes plus report events; avoid allowing
turn processing to grow into one monolithic function.

## Immediate implementation milestone

After Phase 1 decisions are recorded, implement the pure domain model for
identities, locations, entities, unit specifications, inventories, and
mass/volume calculations. Validate it with table-driven tests derived from the
canonical reference, including the TRNS-2 and TRNS-10 examples.

Do not include SQLite migrations, full universe generation, all 21 stages,
combat, markets, or a complete report format in that first domain milestone.
