---
title: Unit Reference
linkTitle: Units
weight: 30
---

Technical description of every unit in Epimethean Challenge: its category,
mass, volume, build inputs, operational requirements, and per-turn output.
Units are the components installed inside Entities (Ships and Colonies). For
the universe structure and the broader unit taxonomy, see the
[Canonical Reference]({{< ref "ec-canonical-reference.md" >}}).

> **⚠️ Work in progress.** This reference is under active development. Some
> values are provisional (`TBD`), and unit statistics are being reconciled
> against the current engine. Figures may change.

---

## 1. Categories

Every unit belongs to one **Category** (the unit's `Class`). The table below
summarizes the categories; the [full unit table](#3-units) lists each unit's
statistics.

| Category       | Item Codes                             | Description                                                                                                                               |
|----------------|----------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| Ammunition     | AMSL, CSUP, MSSL                       | Expendable munitions consumed in combat.                                                                                                  |
| Cadre          | CNST, POL, RBL, SPAG, SPY, TRNE, WRKR  | Roles filled by population units on temporary assignment; their food and consumer-goods needs are counted with the underlying population. |
| Commodity      | CSGD, FOOD                             | Consumable goods that feed the population and set its standard of living.                                                                 |
| Infrastructure | AUTO, BMR, FACT, FARM, LAB, MINE, POWR | Assembled installations that produce output each turn (production, power, research).                                                      |
| Living         | PRO, SOL, UEM, USK                     | Population units whose counts change each turn through demographics.                                                                      |
| Propulsion     | HDRV, SDRV                             | Assembled drives that move or maneuver an entity.                                                                                         |
| Recon          | PROB, SENS                             | Sensor and probe equipment used to gather information.                                                                                    |
| Resource       | FUEL, METL, NMTL                       | Raw materials extracted from planetary deposits and consumed in production.                                                               |
| Static         | LSU                                    | Assembled support installation (life support).                                                                                            |
| Structural     | STRC, STRL                             | Material assembled to enclose volume for ships and colonies.                                                                              |
| Technology     | PRTO, RP                               | Units used to advance or transfer Tech Level; may be non-physical.                                                                        |
| Transportation | TRNS                                   | Units that move population and materials between entities at a planet.                                                                    |
| Weaponry       | ACFT, AWPN, ESHD, EWPN, MLNC, MRBT     | Combat systems that inflict or deflect damage; most require assembly and crew.                                                            |

---

## 2. Conventions

Values that scale with a unit's Tech Level are written as expressions in `TL`
(for example `4 × TL`, `TL + 6`, `10 × TL^2`).

| Column | Meaning |
|--------|---------|
| **Item Code** | Canonical short code used in orders and reports. |
| **Category** | The unit's `Class` (see [§1](#1-categories)). |
| **Assembled?** | `yes` if the unit must be assembled before it operates. |
| **Mass (MU)** | Mass in Mass Units. |
| **Volume (VU)** | Assembled (operational) volume in Volume Units. |
| **Volume Stowed (VU)** | Dis-assembled (storage) volume; `—` for units that cannot be stowed (population and cadre). |
| **Metals / Non-Metals to Build** | Metallics (`METL`) and Non-Metallics (`NMTL`) consumed to manufacture the unit. `N/A` for units that are not manufactured (resources, food, population, cadre); `TBD` where not yet specified. |
| **Operational Requirements** | Labor, fuel, power, or other units needed to operate the unit each turn; `N/A` if none. |

Abbreviations: **TL** = Tech Level (integer 0–10) · **MU** = Mass Unit ·
**VU** = Volume Unit · **TBD** = not yet specified.

**Stowability.** Every unit except Population and Cadre can be stowed; each such
unit lists its dis-assembled storage volume in the **Volume Stowed** column.
Population and Cadre can never be stowed — they show `—` and are always carried at
their full 1 VU. Units that must be assembled to operate (**Assembled?** = `yes`)
are stowed while unassembled and, like all stowable units, must be in the stowed
(unassembled) state to be transferred.

---

## 3. Units

Sorted by Item Code.

| Item Code | Item Name           | Category       | Assembled? | Mass (MU)       | Volume (VU)     | Volume Stowed (VU) | Metals to Build | Non-Metals to Build | Operational Requirements                                               | Output / Notes                                                                                |
|-----------|---------------------|----------------|------------|-----------------|-----------------|--------------------|-----------------|---------------------|------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| ACFT      | Assault Craft       | Weaponry       | no         | 5 × TL          | 5 × TL          | 2.5 × TL           | 3 × TL          | 2 × TL              | 1 soldier or military robot equivalent + 0.1 fuel in combat            | Provides 10 × TL combat factors does not require transports to attack                         |
| AMSL      | Anti-Missile        | Ammunition     | no         | 4 × TL          | 4 × TL          | 2 × TL             | 2 × TL          | 2 × TL              | Missile Launcher of same TL                                            | Destroys Missiles see combat                                                                  |
| AUTO      | Automation          | Infrastructure | yes        | 4 × TL          | 4 × TL          | 2 × TL             | 2 × TL          | 2 × TL              | N/A                                                                    | Replaces 1 × TL^2 Unskilled see Automation in Production chapter                              |
| AWPN      | Assault Weapons     | Weaponry       | no         | 2 × TL          | 2 × TL          | 1 × TL             | 1 × TL          | 1 × TL              | 1 soldier or military robot equivalent                                 | Provides 2 × TL^2 combat factors                                                              |
| BMR       | Beamer              | Infrastructure | yes        | 40 × (TL + 115) | 40 × (TL + 115) | 20 × (TL + 115)    | 10 × (TL + 210) | 30 × TL + 2500      | 1 PRO / 25 Beamers and 20 × TL^2 fuel per Beamer                       | Beams 5000 × TL^2 Mass Units                                                                  |
| CNST      | Construction Worker | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | Provide manpower to assemble, disassemble, add on and scrap items to/from SCs                 |
| CSGD      | Consumer Goods      | Commodity      | no         | 0.6             | 0.3             | 0.3                | 0.2             | 0.4                 | N/A                                                                    | Consumption determines S.O.L.                                                                 |
| CSUP      | Combat Supplies     | Ammunition     | no         | 0.04            | 0.04            | 0.04               | 0.02            | 0.02                | N/A                                                                    | Consumed in combat - 1 per soldier per combat round                                           |
| ESHD      | Energy Shield       | Weaponry       | yes        | 20 × TL         | 20 × TL         | 10 × TL            | 10 × TL         | 10 × TL             | 1 soldier / 100 and uses 10 × TL fuel only                             | Deflects 10 × TL^2 energy units per use                                                       |
| EWPN      | Energy Weapon       | Weaponry       | yes        | 10 × TL         | 10 × TL         | 5 × TL             | 5 × TL          | 5 × TL              | 1 soldier / 100 and uses 4 × TL fuel only                              | Destroys 10 × TL^2 mass per hit                                                               |
| FACT      | Factory             | Infrastructure | yes        | 2 × (TL + 6)    | 2 × (TL + 6)    | TL + 6             | TL + 8          | 4 + TL              | 1 professional 3 unskilled and 0.5 × TL fuel or power                  | Produces 20 × TL mass per turn see Manufacturing                                              |
| FARM      | Farm                | Infrastructure | yes        | 2 × (TL + 3)    | 2 × (TL + 3)    | TL + 3             | TL + 4          | 2 + TL              | 1 professional 3 unskilled fuel or power varies                        | Production varies via TL see Farming                                                          |
| FOOD      | Food                | Commodity      | no         | 6               | 3               | 3                  | N/A             | N/A                 | N/A                                                                    | Feeds 4 to 16 population each turn see Basic Needs                                            |
| FUEL      | Fuel                | Resource       | no         | 1               | 0.5             | 0.5                | N/A             | N/A                 | N/A                                                                    | Raw material used by many units                                                               |
| HDRV      | Hyper Drive         | Propulsion     | yes        | 45 × TL         | 45 × TL         | 22.5 × TL          | 25 × TL         | 20 × TL             | 1 professional / 100 and 40 fuel only per light year                   | Lift capacity 1045 × TL range is (Square root of TL) + 4                                      |
| LAB       | Laboratory          | Infrastructure | yes        | 2 × TL + 8      | 2 × TL + 8      | TL + 4             | TL + 5          | TL + 3              | 3 professional 1 unskilled, 0.5 × TL fuel or power                     | Produces 0.25 × TL research per turn                                                          |
| LSU       | Life Support        | Static         | yes        | 8 × TL          | 8 × TL          | 4 × TL             | 3 × TL          | 5 × TL              | 1 × TL fuel or power                                                   | Supports 1 × TL^2 population                                                                  |
| METL      | Metallics           | Resource       | no         | 1               | 0.5             | 0.5                | N/A             | N/A                 | N/A                                                                    | Raw material used in production                                                               |
| MINE      | Mine                | Infrastructure | yes        | 2 × (TL + 5)    | 2 × (TL + 5)    | TL + 5             | TL + 5          | TL + 5              | 1 professional 3 unskilled, 0.5 × TL fuel or power                     | Mines 25 × TL per turn in raw ore. Actual net depends on yield of deposit see Mining          |
| MLNC      | Missile Launcher    | Weaponry       | yes        | 25 × TL         | 25 × TL         | 12.5 × TL          | 15 × TL         | 10 × TL             | 1 soldier / 100                                                        | Launches 1 missile per attack see Combat                                                      |
| MRBT      | Military Robot      | Weaponry       | no         | 2 × (TL + 10)   | 2 × (TL + 10)   | TL + 10            | TL + 10         | TL + 10             | N/A                                                                    | Same as 2 × TL soldiers                                                                       |
| MSSL      | Missile             | Ammunition     | no         | 4 × TL          | 4 × TL          | 4 × TL             | 2 × TL          | 2 × TL              | Missile Launcher of same TL                                            | Destroys 100 × TL^2 Mass                                                                      |
| NMTL      | Non-Metallics       | Resource       | no         | 1               | 0.5             | 0.5                | N/A             | N/A                 | N/A                                                                    | Raw material used in production                                                               |
| POL       | Police              | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | Drafted to keep the peace and suppress rebellion by arresting rebels                          |
| POWR      | Power Plant         | Infrastructure | yes        | 2 × (TL + 5)    | 2 × (TL + 5)    | TL + 5             | TL + 5          | TL + 5              | Open Surface Colonies ONLY                                             | Produces TL Power per turn (think hydro electric)                                             |
| PRO       | Professionals       | Living         | no         | 1               | 1               | —                  | N/A             | N/A                 | Food and Consumer Goods                                                | Provide administrative and managerial services to operate certain assembled items             |
| PROB      | Probe               | Recon          | no         | 500/TL          | 500/TL          | 500/TL             | 200/TL          | 300/TL              | N/A                                                                    | Obtains probe data - expended when used                                                       |
| PRTO      | Prototype           | Technology     | no         | TL              | 3 × TL          | 3 × TL             | TBD             | TBD                 | N/A                                                                    | Transfer Tech Levels between Entities                                                         |
| RBL       | Rebel               | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | Tally of population willing to rebel; does not allocate the underlying units                  |
| RP        | Research Point      | Technology     | no         | 0               | 0               | 0                  | N/A             | N/A                 | N/A                                                                    | Expended to increase TLs                                                                      |
| SDRV      | Space Drive         | Propulsion     | yes        | 25 × TL         | 25 × TL         | 12.5 × TL          | 15 × TL         | 10 × TL             | 1 professional / 100 and 1 × TL fuel only                              | Produces 3000 × TL^2 thrust                                                                   |
| SENS      | Sensor              | Recon          | yes        | 3000 × TL       | 3000 × TL       | 1500 × TL          | 1000 × TL       | 2000 × TL           | 0.05 × TL fuel or power                                                | Used to obtain probe information                                                              |
| SOL       | Soldiers            | Living         | no         | 1               | 1               | —                  | N/A             | N/A                 | Food and Consumer Goods as well as 1 × military supplies during combat | Provide combat factors and also staffing for certain weapons systems                          |
| SPAG      | Special Agent       | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | Used to infiltrate rebel population sectors and assist police in locating rebels              |
| SPY       | Spy                 | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | Used to perform espionage                                                                     |
| STRC      | Structure           | Structural     | yes        | 1 × TL          | 1 × TL          | 1 × TL             | 0.7 × TL        | 0.3 × TL            | Must be assembled to enclose volume                                    | Encloses (1 × TL^2) divided by structural requirement (based on Entity type)                  |
| STRL      | Light Structure     | Structural     | yes        | 0.1 × TL        | 0.1 × TL        | 0.1 × TL           | 0.07 × TL       | 0.03 × TL           | Must be assembled to enclose volume                                    | Light structure; encloses (1 × TL^2) divided by structural requirement (based on Entity type) |
| TRNE      | Trainee             | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | Professionals to perform Training                                      | Trainees are in training to become Professionals                                              |
| TRNS      | Transport           | Transportation | no         | 4 × TL          | 4 × TL          | 4 × TL             | 3 × TL          | 1 × TL              | 1 professional / 10 and 0.1 × TL^2 fuel                                | Transports 20 × TL^2 Mass; In Combat = 3 × TL^2 Mass                                          |
| UEM       | Unemployables       | Living         | no         | 1               | 1               | —                  | N/A             | N/A                 | Food and Consumer Goods                                                | Population growth occurs in this category                                                     |
| USK       | Unskilled Workers   | Living         | no         | 1               | 1               | —                  | N/A             | N/A                 | Food and Consumer Goods                                                | Provide labor to operate certain assembled items                                              |
| WRKR      | Worker              | Cadre          | no         | 1               | 1               | —                  | N/A             | N/A                 | N/A                                                                    | PRO and USK units allocated to a FACT, FARM, or MINE                                          |
