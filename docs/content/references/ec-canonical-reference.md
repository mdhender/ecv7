---
title: Epimethean Challenge Canonical Reference
linkTitle: Canonical Reference
weight: 10
---

The single, reconciled reference for Epimethean Challenge (EC). It consolidates
the taxonomy, unit list, mass/volume tables, transport rules, combat phases,
turn-processing sequence, and order syntax into one authoritative source.

This reference records the **facts** of the game — its structure, taxonomy,
units, and fixed tables — so that a solid model can be built on them. Behavioral
**game rules** that operate on those facts (transfer mechanics, production
output rates, population dynamics, wages) are defined separately once the model
is well defined; until then they are listed in
[§12]({{< ref "#12-areas-not-yet-specified" >}}).

---

## 1. Universe Structure

The game universe is a fixed, fully bounded space called **The Cluster**.

```
The Cluster
└── Stellium (exactly 100)
    └── Systems (1–6 per Stellium)
        └── Planets (1–10 per System)
```

### 1.1 The Cluster
The entirety of accessible, explorable space. Contains exactly 100 Stellium.

### 1.2 Stellium
A gravitationally associated group of Systems. Contains between 1 and 6 Systems.
Every Stellium has a unique integer identifier and `(x, y, z)` map coordinates,
where x, y, and z are integers from -15 to 15 and the origin `(0, 0, 0)` is the
center of the Cluster. A Stellium is identified by its integer ID, not its
coordinates; the coordinates are displayed in reports (see
[§1.5]({{< ref "#15-identity-and-coordinate-display" >}})).

### 1.3 System
A star and its orbiting Planets. Contains between 1 and 10 Planets. A System is
the Home System of any Faction whose Home Planet orbits it.

Systems within a Stellium are ordered by a sequence letter (`seq_no`), starting
at `A` and progressing to `B`, `C`, `D`, and so on. Every System has a unique
integer identifier; reports display its coordinates as `(x, y, z, seq)` — its
Stellium's coordinates plus the sequence letter.

### 1.4 Planet
A body orbiting a System's star. Every Planet has a **Type** and **Resources**.
Every Planet has a unique integer identifier; reports display its coordinates
as `(x, y, z, seq, orbit)`.

#### Planet Types
| Type              | Notes                                 |
|-------------------|---------------------------------------|
| Rocky Terrestrial | Supports surface and orbital colonies |
| Asteroid Belt     | Supports orbital colonies only        |
| Gas Giant         | Supports orbital colonies only        |

#### Planet Resources
| Resource         | Description                                                                                                                       |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| **Habitability** | Integer 0–25. Limits the number of FARM-1 units installable on the planet's surface.                                              |
| **Deposits**     | A collection of Natural Resource deposits, each with a Quantity and Yield Percentage. Located by survey, extracted by MINE units. |

A Planet holds at most 40 Deposits, numbered sequentially (`deposit_no` 1–40)
when the Planet is generated; deposit numbers are never reused or renumbered. A
Deposit's coordinate is its Planet's coordinate plus its number:
`(x, y, z, seq, orbit, deposit_no)`.

### 1.5 Identity and Coordinate Display

Stellium, Systems, and Planets are identified by their unique integer IDs, not
by their coordinates. Coordinates are a display form used in reports:

| Object   | Identified by | Displayed as          |
|----------|---------------|-----------------------|
| Stellium | integer ID    | `(x, y, z)`           |
| System   | integer ID    | `(x, y, z, seq)`      |
| Planet   | integer ID    | `(x, y, z, seq, orbit)` |

Turn Reports display both the integer ID and the coordinate form for Stellium,
Systems, Planets, Ships, and Colonies, to help players write orders.

---

## 2. Factions & Species

A **Faction** is the primary actor in the game, controlled by a human or agentic
player. Every Faction originates from a **Home Planet**.

A **Species** is the group of all Factions sharing a common Home Planet. A Faction
with a unique Home Planet constitutes a Species of one.

A Faction's **Home System** is the System containing their Home Planet.

Factions own and direct **Entities** (Ships and Colonies). Entities are composed
of **Units**.

---

## 3. Entities

**Entities** are the two top-level game objects that may receive Orders. Each
Entity is owned by a Faction, has a Tech Level (TL), and is composed of installed
Units.

### 3.1 Ship (SHIP-TL)

A mobile Entity capable of movement within and between Systems. Ships are
differentiated by their TL and their installed Units.

### 3.2 Colonies

A Colony is a fixed Entity established on or orbiting a Planet. A Faction is
limited to one of each Colony type per Planet. Multiple Factions may each build
their own Colonies on the same Planet.

#### Colony Types
| Code    | Name                    | Location       |
|---------|-------------------------|----------------|
| COPN-TL | Open Surface Colony     | Planet surface |
| CENC-TL | Enclosed Surface Colony | Planet surface |
| CORB-TL | Orbital Colony          | Orbit          |

#### Starting Condition
Every Faction begins the game with a COPN on their Home Planet. This starting
COPN comes pre-installed with Factories, Mines, Farms, and Population Units,
resolving the bootstrap dependency on Factories to produce Factories.

---

## 4. Units

Units are installed inside Entities. Each unit belongs to one **Class**; the
Classes are listed below.

> **TL Notation:** All Production Units carry a `-TL` suffix (0–10), e.g. `FACT-3`,
> `MINE-7`. By convention, TL 0 units may omit the suffix (`CSUP` and `CSUP-0` are
> both valid). `FACT`, `FARM`, and `MINE` always require the `-TL` suffix.

### 4.1 Resource

| Code | Name          | Description                                |
|------|---------------|--------------------------------------------|
| FUEL | Fuel          | Used for all production and transportation |
| METL | Metallics     | All metallic substances                    |
| NMTL | Non-Metallics | All non-metallic substances                |

### 4.2 Commodity

| Code | Name           | Description                                                                                   |
|------|----------------|-----------------------------------------------------------------------------------------------|
| CSGD | Consumer Goods | Produced by FACT; used to pay the population                                                  |
| FOOD | Food           | Produced by FARM; consumed by all Population Units (¼ FOOD unit per population unit per Turn) |

### 4.3 Living

| Code | Name              | Description                                                                               |
|------|-------------------|-------------------------------------------------------------------------------------------|
| USK  | Unskilled Workers | Workers not requiring long training; provide Labor alongside PRO                          |
| PRO  | Professionals     | Workers requiring long apprenticeships or extensive training; provide Labor alongside USK |
| SOL  | Soldiers          | All military personnel; required for combat                                               |
| UEM  | Unemployables     | Non-working population: children, the elderly, pregnant females, combat casualties        |

#### Population Dynamics
Each Turn the game engine adjusts Population Unit quantities via Birth Rate and
Death Rate (both derived from the Entity's SOL attribute), Combat Damage,
Starvation (insufficient FOOD), and Life Support Failure (SHIP-TL, CENC-TL, and
CORB-TL only).

#### Standard of Living (SOL)
SOL is a per-Entity attribute calculated each Turn by the game engine. It drives
Birth Rate and Death Rate and reflects population wellbeing. Birth Rate increases
with unpopulated habitable land and inversely with SOL.

#### Labor
FACT, FARM, and MINE units require **Labor** to operate, drawn from the WRKR
Cadre (populated by allocating USK and PRO). Labor is allocated, not consumed —
WRKR units are reserved for the Turn but the underlying Population Units are not
destroyed.

### 4.4 Cadre

Cadre units are derived counts — Population Units temporarily assigned to a role.
Assigning Population Units to a Cadre allocates them for the Turn. **RBL is the
sole exception:** it counts population willing to rebel without allocating the
underlying units.

| Code | Name                | Description                                                                  |
|------|---------------------|------------------------------------------------------------------------------|
| CNST | Construction Worker | Execute assembly and disassembly orders                                      |
| POL  | Police              | Keep the peace and suppress rebellion by arresting rebels                    |
| RBL  | Rebel               | Count of population willing to rebel; does not allocate the underlying units |
| SPAG | Special Agent       | Infiltrate rebel population sectors and assist Police in locating rebels     |
| SPY  | Spy                 | Report on other Factions and incite rebellion                                |
| TRNE | Trainee             | Unskilled workers in training toward professional status                     |
| WRKR | Worker              | PRO and USK units allocated to a FACT, FARM, or MINE                         |

### 4.5 Ammunition

| Code | Name            | Description                                                                    |
|------|-----------------|--------------------------------------------------------------------------------|
| AMSL | Anti-Missile    | Launched by MLNC to destroy incoming MSSL                                      |
| CSUP | Combat Supplies | Ammunition and medicines; consumed in combat (1 unit per SOL per combat round) |
| MSSL | Missile         | Used in any combat; less accurate than energy weapons                          |

### 4.6 Weaponry

| Code | Name             | Description                                                            |
|------|------------------|------------------------------------------------------------------------|
| ACFT | Assault Craft    | Land/space vehicles used to invade Colonies or Ships                   |
| AWPN | Assault Weapons  | Used by SOL on a planet's surface                                      |
| ESHD | Energy Shield    | Deflect energy beams                                                   |
| EWPN | Energy Weapon    | Concentrated energy beam; used in all combat except surface-to-surface |
| MLNC | Missile Launcher | Launch MSSL; accuracy depends on launcher TL                           |
| MRBT | Military Robot   | Replace SOL units (TL × 2 soldiers per MRBT unit)                      |

### 4.7 Infrastructure

| Code | Name        | Description                                                                                   |
|------|-------------|-----------------------------------------------------------------------------------------------|
| AUTO | Automation  | Replaces USK workers in FACT, FARM, or MINE (unit × TL = worker units replaced)               |
| BMR  | Beamer      | Beams mass between Entities (5000 × TL² MU per use); 1 PRO per 25 Beamers                     |
| FACT | Factory     | Manufactures all units except Natural Resources, FOOD, and Population                         |
| FARM | Farm        | Produces FOOD. TL 1: open-colony; TL 2–5: hydroponic; TL 6–10: artificial-sunlight hydroponic |
| LAB  | Laboratory  | Generates Research Points (RPs) each Turn                                                     |
| MINE | Mine        | Extracts and refines Natural Resources from Deposits                                          |
| POWR | Power Plant | Produces TL Power per Turn (hydroelectric); Open Surface Colonies only                        |

Notes:

* FARM-1 may only be assembled in COPN on planets with Habitability > 0 and within orbits 1 through 5.
* FARM-2 through FARM-5 may only be assembled in CENC or CORB within orbits 1 through 5.
* FARM-6 through FARM-10 may be assembled on CENC, CORB, or SHIP in any orbit.
* MINE may only be assembled in COPN.
* FACT and FARM-2 through FARM-5 require no FUEL to operate when assembled in CORB within orbits 1 through 5.

### 4.8 Technology

| Code | Name           | Description                                      |
|------|----------------|--------------------------------------------------|
| PRTO | Prototype      | Used to transfer Tech Levels between Entities    |
| RP   | Research Point | Currency for TL advancement; not a physical unit |

### 4.9 Propulsion

| Code | Name        | Description                                                                                       |
|------|-------------|---------------------------------------------------------------------------------------------------|
| HDRV | Hyper Drive | Propels Ships through hyperspace; jump range = 1 light-year × TL                                  |
| SDRV | Space Drive | Maintains orbit and maneuvers in combat; cannot be used for interplanetary or interstellar travel |

### 4.10 Recon

| Code | Name   | Description                                                                  |
|------|--------|------------------------------------------------------------------------------|
| PROB | Probe  | Executed by SENS; report on Ships, Colonies, and Deposits in the same System |
| SENS | Sensor | Detect Ships and Colonies in orbit; conduct probes (1 × TL probes per Turn)  |

### 4.11 Static

| Code | Name         | Description                                                                                               |
|------|--------------|-----------------------------------------------------------------------------------------------------------|
| LSU  | Life Support | Recycles air and water; required on SHIP-TL, CENC-TL, and CORB-TL (supports TL² population units per LSU) |

### 4.12 Structural

| Code | Name            | Description                                                                                              |
|------|-----------------|----------------------------------------------------------------------------------------------------------|
| STRC | Structure       | Required to build Ships and Colonies; encloses (1 × TL²) ÷ structural requirement (based on Entity type) |
| STRL | Light Structure | Light structural variant of STRC at one-tenth the mass, volume, and cost                                 |

Notes:

* STRC may be built in COPN, CENC, or CORB.
* STRL may only be built in CORB.

### 4.13 Transportation

| Code | Name      | Description                                                                                                              |
|------|-----------|--------------------------------------------------------------------------------------------------------------------------|
| TRNS | Transport | Transfer Population and materials between Ships/Colonies at the same Planet (see [§7]({{< ref "#7-transports-trns" >}})) |

---

## 5. Mass, Volume & Inputs

**Definitions:** Mass and Volume are in mass units (MU) and volume units (VU).
**Volume** is the assembled/operational volume; **Volume Stowed** is the
dis-assembled storage volume. **METL/NMTL Input** is the build cost. An
**Operational** "Yes" means the unit must be assembled before use. `TBD` marks
values not yet specified.

Population and Cadre units are not manufactured: each is 1 MU and 1 VU and can
never be stowed. Research Points (`RP`) are non-physical: 0 MU, 0 VU.

| Class          | Code | Unit             | METL In         | NMTL In        | Mass            | Volume          | Stowed          | Op. |
|----------------|------|------------------|-----------------|----------------|-----------------|-----------------|-----------------|-----|
| Resource       | FUEL | Fuel             | —               | —              | 1               | 0.5             | 0.5             | No  |
| Resource       | METL | Metallics        | —               | —              | 1               | 0.5             | 0.5             | No  |
| Resource       | NMTL | Non-Metallics    | —               | —              | 1               | 0.5             | 0.5             | No  |
| Commodity      | CSGD | Consumer Goods   | 0.2             | 0.4            | 0.6             | 0.3             | 0.3             | No  |
| Commodity      | FOOD | Food             | —               | —              | 6               | 3               | 3               | No  |
| Ammunition     | AMSL | Anti-Missile     | 2 × TL          | 2 × TL         | 4 × TL          | 4 × TL          | 2 × TL          | No  |
| Ammunition     | CSUP | Combat Supplies  | 0.02            | 0.02           | 0.04            | 0.04            | 0.04            | No  |
| Ammunition     | MSSL | Missile          | 2 × TL          | 2 × TL         | 4 × TL          | 4 × TL          | 4 × TL          | No  |
| Weaponry       | ACFT | Assault Craft    | 3 × TL          | 2 × TL         | 5 × TL          | 5 × TL          | 2.5 × TL        | No  |
| Weaponry       | AWPN | Assault Weapons  | 1 × TL          | 1 × TL         | 2 × TL          | 2 × TL          | 1 × TL          | No  |
| Weaponry       | ESHD | Energy Shield    | 10 × TL         | 10 × TL        | 20 × TL         | 20 × TL         | 10 × TL         | Yes |
| Weaponry       | EWPN | Energy Weapon    | 5 × TL          | 5 × TL         | 10 × TL         | 10 × TL         | 5 × TL          | Yes |
| Weaponry       | MLNC | Missile Launcher | 15 × TL         | 10 × TL        | 25 × TL         | 25 × TL         | 12.5 × TL       | Yes |
| Weaponry       | MRBT | Military Robot   | TL + 10         | TL + 10        | 2 × (TL + 10)   | 2 × (TL + 10)   | TL + 10         | No  |
| Infrastructure | AUTO | Automation       | 2 × TL          | 2 × TL         | 4 × TL          | 4 × TL          | 2 × TL          | Yes |
| Infrastructure | BMR  | Beamer           | 10 × (TL + 210) | 30 × TL + 2500 | 40 × (TL + 115) | 40 × (TL + 115) | 20 × (TL + 115) | Yes |
| Infrastructure | FACT | Factory          | TL + 8          | 4 + TL         | 2 × (TL + 6)    | 2 × (TL + 6)    | TL + 6          | Yes |
| Infrastructure | FARM | Farm             | TL + 4          | 2 + TL         | 2 × (TL + 3)    | 2 × (TL + 3)    | TL + 3          | Yes |
| Infrastructure | LAB  | Laboratory       | TL + 5          | TL + 3         | 2 × TL + 8      | 2 × TL + 8      | TL + 4          | Yes |
| Infrastructure | MINE | Mine             | TL + 5          | TL + 5         | 2 × (TL + 5)    | 2 × (TL + 5)    | TL + 5          | Yes |
| Infrastructure | POWR | Power Plant      | TL + 5          | TL + 5         | 2 × (TL + 5)    | 2 × (TL + 5)    | TL + 5          | Yes |
| Technology     | PRTO | Prototype        | TBD             | TBD            | TL              | 3 × TL          | 3 × TL          | No  |
| Propulsion     | HDRV | Hyper Drive      | 25 × TL         | 20 × TL        | 45 × TL         | 45 × TL         | 22.5 × TL       | Yes |
| Propulsion     | SDRV | Space Drive      | 15 × TL         | 10 × TL        | 25 × TL         | 25 × TL         | 12.5 × TL       | Yes |
| Recon          | PROB | Probe            | 200/TL          | 300/TL         | 500/TL          | 500/TL          | 500/TL          | No  |
| Recon          | SENS | Sensor           | 1000 × TL       | 2000 × TL      | 3000 × TL       | 3000 × TL       | 1500 × TL       | Yes |
| Static         | LSU  | Life Support     | 3 × TL          | 5 × TL         | 8 × TL          | 8 × TL          | 4 × TL          | Yes |
| Structural     | STRC | Structure        | 0.7 × TL        | 0.3 × TL       | 1 × TL          | 1 × TL          | 1 × TL          | Yes |
| Structural     | STRL | Light Structure  | 0.07 × TL       | 0.03 × TL      | 0.1 × TL        | 0.1 × TL        | 0.1 × TL        | Yes |
| Transportation | TRNS | Transport        | 3 × TL          | 1 × TL         | 4 × TL          | 4 × TL          | 4 × TL          | No  |

### 5.1 Notes

* **Storage vs. operational volume:** the Volume Stowed column gives the
  dis-assembled storage volume; the Volume column gives the assembled volume.
* Natural Resources (FUEL, METL, NMTL) and FOOD have no METL/NMTL build cost;
  they are extracted or farmed.

---

## 6. Research & Tech Level

**Tech Level (TL)** is an integer from 0 to 10 tracked per Entity (Ship or Colony).
A Factory can only manufacture units up to its Colony's or Ship's current TL.

**Research Points (RPs)** are generated each Turn by LAB-TL units. RPs may be
expended to increase an Entity's TL per the following schedule:

| Target TL | RPs Required |
|-----------|--------------|
| 2         | 100,000      |
| 3         | 200,000      |
| 4         | 400,000      |
| 5         | 800,000      |
| 6         | 1,600,000    |
| 7         | 3,200,000    |
| 8         | 6,400,000    |
| 9         | 12,800,000   |
| 10        | 25,600,000   |

There is no research cost for reaching TL 1: no Entity ever enters the game
below TL 1, so the schedule begins at target TL 2.

TL may also be advanced by transferring technology from another Entity (via
PRTO-TL) or purchasing it at a market or trade station.

---

## 7. Transports (TRNS)

Transports transfer population and materials between Ships and Colonies at the
same Planet. They carry no armament.

### 7.1 Stats (per unit)

| Stat        | Formula  | TRNS-2 Example |
|-------------|----------|----------------|
| METL Input  | `3 × TL` | 6              |
| NMTL Input  | `1 × TL` | 2              |
| Mass (MU)   | `4 × TL` | 8              |
| Volume (VU) | `4 × TL` | 8              |
| Operational | No       | —              |

### 7.2 Throughput vs. Capacity

**Normal operations** — the transfer rate is a *throughput*, not a single-lift
capacity:

> **TL² × 200 MU per turn**

The throughput is achieved by repeated trips within the turn, each carrying up to
the transport's physical volume (8 VU for a TRNS-2), rather than a single lift.

**Combat operations** — only a single trip is possible per round:

> **3 × TL MU per combat round**

One load, one run.

### 7.3 Worked Examples

- **TRNS-2 moving 800 USK:** 800 VU of cargo ÷ 8 VU batch = **100 trips**, all
  within one turn → 800 MU/turn throughput. In combat: `3 × 2 = 6 MU` per round.
- **TRNS-10 moving 20,000 USK:** throughput `10² × 200 = 20,000 MU/turn`;
  20,000 VU ÷ 40 VU batch = **500 trips**. In combat: `3 × 10 = 30 MU` per round.

### 7.4 Crew & Fuel

- **Crew:** outside combat, 1 PRO unit per 10 TRNS units. In combat, the carried
  soldiers operate the transport themselves.
- **Fuel:** `TL² ÷ 10` per turn, proportional to the percentage of capacity used
  (a full TRNS-2 uses 0.4; half-loaded, 0.2). In combat, fuel use is `0.01 × TL²`
  per round trip.

---

## 8. Economy & Markets

Prices, wages, commissions, and fees are denominated in an abstract currency
(there is no gold or other named specie in the current game).

| Concept                | Description                                                                                                        |
|------------------------|--------------------------------------------------------------------------------------------------------------------|
| Market / Trade Station | Venues where units and technology are bought and sold. A Trade Station is an orbital colony dedicated to commerce. |
| Home Planet Market     | A fixed market on each Species' home planet.                                                                       |
| Commission             | 1% fee paid to a Trade Station on completed transactions.                                                          |
| Fee                    | Diplomatic payment for colonization permission or resource access.                                                 |
| Wage                   | Currency paid to population per turn; set by `PAY` orders per population class.                                    |
| Consumer Goods (CSGD)  | Manufactured goods used to pay/satisfy the population.                                                             |

Traded items include manufactured units and technology levels (e.g. STRC, SDRV,
and TL transfers are explicit examples).

---

## 9. Victory Conditions

### 9.1 Control

- **Colony Control:** A Colony is controlled by a Faction if it contains at least
  1 SOL or PRO Population Unit belonging to that Faction.
- **Planet Control (Faction):** A Faction controls a Planet when it controls every
  Colony on or orbiting it. At least one Colony must exist on the Planet.
- **Planet Control (Species):** A Species controls a Planet if it controls the
  majority of Colonies on or orbiting it (at least ⌈C × 0.5⌉ + 1, where C = total
  Colonies on that Planet). A Species controls a Colony if any Faction within the
  Species controls it.

### 9.2 Solo Victory

Victory is evaluated against **habitable Planets** (Habitability > 0). Let **H** =
total habitable Planets in the Cluster. A single Faction wins when:

- It controls at least ⌈H × 0.5⌉ + 1 habitable Planets, **and**
- No other single Faction controls more than ⌈H × 0.1⌉ + 1 habitable Planets.

### 9.3 Species Victory

Let **H** = total habitable Planets in the Cluster. A Species wins when:

- It controls at least ⌈H × 0.5⌉ + 1 habitable Planets, **and**
- No other single Species controls more than ⌈H × 0.1⌉ + 1 habitable Planets.

---

## 10. Game Mechanics

### 10.1 Turns
The game progresses in discrete **Turns**. Each Turn:
1. Players submit **Orders** for their Entities.
2. The game engine processes all Orders.
3. The engine returns a **Turn Report** to each player.

### 10.2 Orders
Orders are instructions submitted by players, targeting Ships, Colonies, or
(for a few order types) a Player directly. The canonical order syntax is defined
in [§11]({{< ref "#11-orders-reference-canonical" >}}).

### 10.3 Turn Report
The game engine's response to a Turn, returned to each player. Report structure is
not yet specified (see [§12]({{< ref "#12-areas-not-yet-specified" >}})). Whatever
its final format, a report displays both the integer ID and the coordinate form
for Stellium, Systems, Planets, Ships, and Colonies
([§1.5]({{< ref "#15-identity-and-coordinate-display" >}})), to help players
write orders.

### 10.4 Turn Processing Sequence

When the engine processes a Turn, it runs the following 21 stages in order. The
sequence is fixed: an order's effects are resolved at its stage, using the game
state left by all earlier stages.

| # | Stage | Description |
|---|---|---|
| 1 | Mining & Farming Production | MINE extraction and FARM food output are calculated. |
| 2 | Manufacturing Production | FACT manufacturing (including research production) is calculated. |
| 3 | Combat | All combat is resolved (see [§10.5]({{< ref "#105-combat" >}}) for the per-round phase sequence). |
| 4 | Set Up Orders | `SETUP` orders are processed (new Ships/Colonies established). |
| 5 | Dis-assembly Orders | `DISASSEMBLE` orders are processed. |
| 6 | Build Change Orders | `BUILDCHANGE` orders are entered. |
| 7 | Mining Change Orders | `MININGCHANGE` orders are entered. |
| 8 | Transfers | `TRANSFER` orders are processed. |
| 9 | Assembly Orders | `ASSEMBLE` orders are processed. |
| 10 | Market & Trade Station Activity | All `BUY`/`SELL`, `PERMIT`, and trade station commerce takes place. |
| 11 | Surveys | `SURVEY` orders are carried out. |
| 12 | Probe & Sensor Reports | PROB/SENS reports are compiled. |
| 13 | Espionage | `SPY` activity takes place. |
| 14 | Ship Movement | `MOVE` (jump) orders occur. |
| 15 | Draft Orders | `DRAFT`/`DISBAND` orders are processed. |
| 16 | Pay & Ration Orders | `PAY` and `RATION` orders are entered. |
| 17 | Rebellion | Rebellion occurs. |
| 18 | Rebel Increases | Rebel (RBL) counts increase. |
| 19 | Naming & Control Orders | `NAME`/`NAMEP`, `CONTROL`/`UNCONTROL`, and `COLONIZE` orders are processed. |
| 20 | Population Increases | Population growth (births) is calculated. |
| 21 | News Service Reports | `NEWS` service reports are compiled. |

> **Notes:**
> - Production (stages 1–2) precedes everything: newly produced resources and
>   units are available to later stages the same Turn.
> - Combat (stage 3) resolves before set-up, transfer, and assembly: only
>   survivors participate in construction and logistics.
> - Population increases (stage 20) are calculated near the end, after casualties,
>   drafts, rationing, and rebellion have adjusted the population base.

### 10.5 Combat

Combat is resolved in stage 3 of the Turn Processing Sequence ([§10.4]({{< ref "#104-turn-processing-sequence" >}})). The following
describes the per-round phase structure within a single combat.

#### Phase Sequence

**Pre-Combat (Round 1 only)**

1. **Troop Deployment** — Committed soldiers are armed and loaded into assault
   craft; overflow goes into transports with assault weapons; the remainder stays
   behind as uncommitted defense.

**Each Round**

2. **Weapons Fire** — All attacks execute simultaneously:
   - Energy beams fire at missiles, transports, and assault craft
   - Anti-missiles fire at incoming missiles, transports, and assault craft
   - Missiles are launched at ships/colonies
   - Defender fires back automatically (energy beams, anti-missiles, missiles)
3. **Intercept Resolution** — Determine what is shot down before reaching target:
   - Anti-missiles vs. incoming missiles
   - Energy beams vs. transports and assault craft (troop-transport casualties)
4. **Casualty Calculation** (raids/invasions) — Compute combat factors and apply
   losses to both sides; split into KIA (70%) and wounded (30%).
5. **Damage Calculation** (bombardment) — Apply un-intercepted missile/energy
   damage; distribute to weapons/drives (75%) and other units (25%).
6. **Surrender Check** — Any side facing 6:1 odds auto-surrenders.
7. **Ship Movement** — Ships with bombard orders move closer; ships under attack
   without bombard orders move away.

**End of Round**

8. **Continuation Check** — Combat continues if the mission is incomplete AND
   soldiers/fuel/military supplies/missiles are not exhausted; otherwise it ends.

**Post-Combat**

9. **Capture Resolution** — If all defenders are destroyed or surrendered, the
   attacker takes control; for raids, stolen units are transferred.

#### Combat Notes

- **Raids are single-round only** — they skip the continuation check and end after
  round 1.
- **Troop deployment (step 1) happens only in round 1** but feeds every subsequent
  round's state.
- **All combat orders execute simultaneously** — there is no attacker-goes-first
  ordering within a round.

---

## 11. Orders Reference (Canonical)

This is the modern keyword-tagged order syntax. It replaces the legacy positional
format.

### 11.1 Conventions

**General structure** — every order is one line; fields are space-separated;
string literals appear in `"double quotes"`; optional fields in `[brackets]`;
repeating fields with `...`.

```
ENTITY scID  KEYWORD  keyword  value  keyword  value  ...
```

The leading `ENTITY scID` (or `PLAYER playerID`) anchors every order to a subject
and gives the parser a reliable sync point at the start of each line.

**IDs**

| Token        | Description                                  |
|--------------|----------------------------------------------|
| `scID`       | Ship or colony ID (integer)                  |
| `playerID`   | Player ID (integer)                          |
| `systemID`   | System ID (integer)                          |
| `locationID` | Planet ID (integer) or System ID (integer)   |
| `groupNo`    | Factory or mine group number (integer)       |

**Values**

| Token      | Description                                                     |
|------------|-----------------------------------------------------------------|
| `qty`      | Integer quantity, e.g. `1000`                                   |
| `unitCode` | Unit code from the unit list, e.g. `SOL`, `FACT-3`, `MINE-2`    |
| `tlLevel`  | Tech level integer 0–10                                         |
| `pct`      | Percentage as integer 0–100 (no `%` symbol; `commit` labels it) |
| `price`    | Decimal price per unit, e.g. `0.25`                             |
| `wage`     | Decimal wage multiplier, e.g. `1.2`                             |
| `name`     | Quoted string, max 24 characters, e.g. `"Dragonfire"`           |
| `text`     | Quoted free text, e.g. `"Selling drives next turn"`             |

### 11.2 Combat Orders
One combat order per entity per turn.

```
ENTITY scID  BOMBARD   target scID  commit pct
ENTITY scID  INVADE    target scID  commit pct
ENTITY scID  RAID      target scID  commit pct  steal unitCode
ENTITY scID  SUPPORT   ally scID    commit pct
ENTITY scID  SUPPORT   ally scID    attacking scID  commit pct
```
- `SUPPORT` with only `ally` and `commit` is defensive support.
- `SUPPORT` with `attacking scID` specifies which enemy the ally is attacking.

### 11.3 Set Up Orders

```
SETUP  location locationID  type ("SHIP"|"COPN"|"CENC"|"CORB")  source scID
  TRANSFER  qty unitCode
  TRANSFER  qty unitCode
  ...
END SETUP
```
- `SETUP` / `END SETUP` are block delimiters — no `scID` prefix, as the new entity
  does not yet have one.
- Each `TRANSFER` line moves one unit type from the source entity.

### 11.4 Assembly Orders

```
ENTITY scID  ASSEMBLE  qty unitCode
ENTITY scID  ASSEMBLE  qty FACT-tlLevel  produce unitCode
ENTITY scID  ASSEMBLE  qty MINE-tlLevel  deposit locationID
```
- Plain `ASSEMBLE qty unitCode` covers all non-factory, non-mine assemblies.
- `produce` labels what the factory group will build.
- `deposit` labels the deposit the mine group will work.

### 11.5 Dis-assembly Orders
Same structure as Assembly, keyword changed to `DISASSEMBLE`.

```
ENTITY scID  DISASSEMBLE  qty unitCode
ENTITY scID  DISASSEMBLE  qty FACT-tlLevel  produce unitCode
ENTITY scID  DISASSEMBLE  qty MINE-tlLevel  deposit locationID
```

### 11.6 Build Change Orders

```
ENTITY scID  BUILDCHANGE  group groupNo  produce unitCode
ENTITY scID  BUILDCHANGE  group groupNo  retool
ENTITY scID  BUILDCHANGE  research
```
- `produce unitCode` sets the new output of the factory group.
- `retool` begins retooling without specifying a new product yet.
- `research` redirects the entity's factories to research production.

### 11.7 Transfer Orders

```
ENTITY scID  TRANSFER  qty unitCode  to scID
```

### 11.8 Mining Change Orders

```
ENTITY scID  MININGCHANGE  group groupNo  deposit locationID
```

### 11.9 Market Orders

```
ENTITY scID  SELL  unitCode  price price
ENTITY scID  BUY   qty unitCode  price price
```
- `SELL` lists a unit type at a price; quantity is not specified (sell all
  available at that price).
- `BUY` includes a quantity cap.

### 11.10 Survey Orders

```
ENTITY scID  SURVEY
```

### 11.11 Probe Orders

```
ENTITY scID  PROBE  orbit orbitNo  [orbit orbitNo  ...]
```
Multiple orbits may be listed, each prefixed with `orbit`.

### 11.12 Spy Orders

```
ENTITY scID  SPY  qty  CHECK REBELS
ENTITY scID  SPY  qty  CONVERT REBELS
ENTITY scID  SPY  qty  CHECK FOR SPIES
ENTITY scID  SPY  qty  ATTACK SPIES    from scID
ENTITY scID  SPY  qty  INCITE REBELS   at scID
ENTITY scID  SPY  qty  GATHER INFO     from scID
```
- `qty` is the number of spy units committed.
- `from scID` / `at scID` identifies the foreign entity targeted.

### 11.13 News Release

```
NEWS  at locationID  text  [sig text]
```
- `at locationID` is a planet or trade station ID.
- `text` is a quoted message body; `sig` is an optional quoted signature.

### 11.14 Move Orders

```
ENTITY scID  MOVE  to locationID
```
`locationID` is a Planet ID or a System ID (both integers).

### 11.15 Draft Orders

```
ENTITY scID  DRAFT    qty unitCode
ENTITY scID  DISBAND  qty unitCode
```

### 11.16 Pay Orders

```
ENTITY scID  PAY  wage wage  class ("USK"|"PRO"|"SOL")
```
Multiple `PAY` orders may be issued per entity, one per population class.

### 11.17 Ration Orders

```
ENTITY scID  RATION  pct pct
```

### 11.18 Control Orders

```
PLAYER playerID  CONTROL    system systemID  orbit orbitNo
PLAYER playerID  UNCONTROL  system systemID  orbit orbitNo
```
- `PLAYER playerID` is the subject, since control is asserted by a player rather
  than a specific ship or colony.

### 11.19 Naming Orders

```
PLAYER playerID  NAMEP  system systemID  planet planetNo  name name
ENTITY scID      NAME   name name
```
- `NAMEP` names a planet; `NAME` names a ship or colony.
- `name` is a quoted string, max 24 characters.

### 11.20 Trade Station Orders

```
ENTITY stationID  PERMIT  player playerID  ("GRANT"|"DENY")
```

### 11.21 Colonising Permission

```
PLAYER playerID  COLONIZE  system systemID  planet planetNo
```

### 11.22 Summary Table

| Keyword               | Subject            | Required fields                            |
|-----------------------|--------------------|--------------------------------------------|
| `BOMBARD`             | `ENTITY scID`      | `target`, `commit`                         |
| `INVADE`              | `ENTITY scID`      | `target`, `commit`                         |
| `RAID`                | `ENTITY scID`      | `target`, `commit`, `steal`                |
| `SUPPORT`             | `ENTITY scID`      | `ally`, `commit`, [`attacking`]            |
| `SETUP` … `END SETUP` | *(block)*          | `location`, `type`, `source`, `TRANSFER`s  |
| `ASSEMBLE`            | `ENTITY scID`      | `qty unitCode`, [`produce`\|`deposit`]     |
| `DISASSEMBLE`         | `ENTITY scID`      | `qty unitCode`, [`produce`\|`deposit`]     |
| `BUILDCHANGE`         | `ENTITY scID`      | `group`, (`produce`\|`retool`\|`research`) |
| `TRANSFER`            | `ENTITY scID`      | `qty unitCode`, `to`                       |
| `MININGCHANGE`        | `ENTITY scID`      | `group`, `deposit`                         |
| `SELL`                | `ENTITY scID`      | `unitCode`, `price`                        |
| `BUY`                 | `ENTITY scID`      | `qty unitCode`, `price`                    |
| `SURVEY`              | `ENTITY scID`      | *(none)*                                   |
| `PROBE`               | `ENTITY scID`      | `orbit` × 1+                               |
| `SPY`                 | `ENTITY scID`      | `qty`, spy-op keyword, [`from`\|`at`]      |
| `NEWS`                | *(global)*         | `at`, `text`, [`sig`]                      |
| `MOVE`                | `ENTITY scID`      | `to`                                       |
| `DRAFT`               | `ENTITY scID`      | `qty unitCode`                             |
| `DISBAND`             | `ENTITY scID`      | `qty unitCode`                             |
| `PAY`                 | `ENTITY scID`      | `wage`, `class`                            |
| `RATION`              | `ENTITY scID`      | `pct`                                      |
| `CONTROL`             | `PLAYER playerID`  | `system`, `orbit`                          |
| `UNCONTROL`           | `PLAYER playerID`  | `system`, `orbit`                          |
| `NAMEP`               | `PLAYER playerID`  | `system`, `planet`, `name`                 |
| `NAME`                | `ENTITY scID`      | `name`                                     |
| `PERMIT`              | `ENTITY stationID` | `player`, (`GRANT`\|`DENY`)                |
| `COLONIZE`            | `PLAYER playerID`  | `system`, `planet`                         |

---

## 12. Areas Not Yet Specified

The following aspects of EC are not defined in this reference:

- Turn Report structure and format
- Ship and Colony construction rules (STRC requirements, total mass)
- Storage vs. operational volume distinction ([§5.6]({{< ref "#56-derivation-rules--notes" >}}))
- Combat damage-resolution formulas (combat factors, accuracy)
- Market and trade station transaction mechanics
- SPY and rebellion mechanics in detail
- Transfer mechanics — how `TRANSFER` orders use TRNS throughput, crew, and
  fuel ([§7]({{< ref "#7-transports-trns" >}}), [§11.7]({{< ref "#117-transfer-orders" >}}))
- Production output rates (FACT manufacturing rates, LAB Research Point
  generation per Turn)
- Population dynamics formulas (SOL calculation, birth and death rates, wages,
  CSGD satisfaction)
- Order-language location tokens after the move to integer IDs: the legacy
  grammar's coordinate form (`TK_LOCATION`, e.g. `4-6-19`) predates ID-based
  identification and must be reconciled
  ([§1.5]({{< ref "#15-identity-and-coordinate-display" >}}),
  [§11.1]({{< ref "#111-conventions" >}}))
