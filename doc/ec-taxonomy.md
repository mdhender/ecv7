# Epimethean Challenge — Taxonomy

## 1. Universe Structure

The game universe is a fixed, fully bounded space called **The Cluster**.

```
The Cluster
└── Constellations (exactly 100)
    └── Systems (1–6 per Constellation)
        └── Planets (1–10 per System)
```

### 1.1 The Cluster
The entirety of accessible, explorable space. Contains exactly 100 Constellations.

### 1.2 Constellation
A gravitationally associated group of Systems. Contains between 1 and 6 Systems.

### 1.3 System
A star and its orbiting Planets. Contains between 1 and 10 Planets. A System is the Home System of any Faction whose Home Planet orbits it.

### 1.4 Planet
A body orbiting a System's star. Every Planet has a **Type** and **Resources**.

#### Planet Types
| Type | Notes |
|---|---|
| Rocky Terrestrial | Supports surface and orbital colonies |
| Asteroid Belt | Supports orbital colonies only |
| Gas Giant | Supports orbital colonies only |

#### Planet Resources
| Resource | Description |
|---|---|
| **Habitability** | Integer 0–25. Limits the number of FARM units installable on the planet's surface. |
| **Deposits** | A collection of Natural Resource deposits, each with a Quantity and Yield Percentage. |

---

## 2. Factions & Species

A **Faction** is the primary actor in the game, controlled by a human or agentic player. Every Faction originates from a **Home Planet**.

A **Species** is the group of all Factions sharing a common Home Planet. A Faction with a unique Home Planet constitutes a Species of one.

A Faction's **Home System** is the System containing their Home Planet.

Factions own and direct **Entities** (Ships and Colonies). Entities are composed of **Units**.

---

## 3. Entities

**Entities** are the two top-level game objects that may receive Orders. Each Entity is owned by a Faction, has a Tech Level (TL), and is composed of installed Units.

### 3.1 Ship (SHIP-TL)

A mobile Entity capable of movement within and between Systems. Ships are differentiated by their TL and their installed Units.

### 3.2 Colonies

A Colony is a fixed Entity established on or orbiting a Planet. A Faction is limited to one of each Colony type per Planet. Multiple Factions may each build their own Colonies on the same Planet.

#### Colony Types
| Code | Name | Location |
|---|---|---|
| COPN-TL | Open Surface Colony | Planet surface |
| CENC-TL | Enclosed Surface Colony | Planet surface |
| CORB-TL | Orbital Colony | Orbit |

#### Starting Condition
Every Faction begins the game with a COPN on their Home Planet. This starting COPN comes pre-installed with Factories, Mines, Farms, and Population Units, resolving the bootstrap dependency on Factories to produce Factories.

---

## 4. Units

Units are installed inside Entities. There are five categories of Units: **Natural Resources**, **Population**, **Cadre**, **Weapons**, and **Production**.

> **TL Notation:** All Production Units carry a `-TL` suffix (0–10), e.g. `FACT-3`, `MINE-7`. By convention, TL 0 units may omit the suffix (`FACT` and `FACT-0` are both valid).

---

### 4.1 Natural Resources

Natural Resources are extracted from planetary Deposits by MINE-TL units and stored as inventory within an Entity. They are consumed by Factories, Farms, and drives.

| Code | Name | Description |
|---|---|---|
| FUEL | Fuel | Used for all production and transportation |
| METL | Metallics | All metallic substances other than gold |
| NMTL | Non-Metallics | All non-metallic substances |

---

### 4.2 Population Units

Population Units represent the people living aboard an Entity. They are not produced by Factories; instead they grow and shrink each Turn through demographic and environmental forces.

| Code | Name | Description |
|---|---|---|
| USK | Unskilled | Workers not requiring long training; provide Labor alongside PRO |
| PRO | Professionals | Workers requiring long apprenticeships or extensive training; provide Labor alongside USK |
| SOL | Soldiers | All military personnel; required for combat |
| UEM | Unemployables | Non-working population: children, the elderly, pregnant females, combat casualties |

#### Population Dynamics
Each Turn the game engine adjusts Population Unit quantities via:

| Factor | Effect |
|---|---|
| Birth Rate | Increases population; derived from Entity SOL attribute |
| Death Rate | Decreases population; derived from Entity SOL attribute |
| Combat Damage | Decreases population |
| Starvation | Decreases population (insufficient FOOD) |
| Life Support Failure | Decreases population (SHIP-TL, CENC-TL, CORB-TL only) |

#### Standard of Living (SOL)
SOL is a per-Entity attribute calculated each Turn by the game engine. It drives Birth Rate and Death Rate and reflects a broader measure of population wellbeing than simple morale.

#### Labor
FACT-TL, FARM-TL, and MINE-TL units require **Labor** to operate. Labor is drawn from the **WRKR** Cadre, which is populated by allocating USK and PRO Population Units. Labor is allocated, not consumed — WRKR units are reserved for the duration of the Turn but the underlying Population Units are not destroyed.

*Example: an Entity with 10 USK/PRO assigned to WRKR, one FACT requiring 8 Labor, and one FACT requiring 3 Labor can only fully staff one Factory per Turn.*

---

### 4.3 Cadre Units

Cadre units are derived counts — they represent Population Units temporarily assigned to a specific role. Assigning Population Units to a Cadre allocates those units, making them unavailable for other uses that Turn.

**RBL is the sole exception:** it counts population willing to rebel but does not allocate those units. Population Units counted in RBL remain available for other assignments in the same Turn.

| Code | Name | Description |
|---|---|---|
| CNST | Construction | Execute assembly and disassembly orders |
| TRNE | Trainees | Unskilled workers in training |
| WRKR | Workers | PRO and USK units allocated to a FACT, FARM, or MINE |
| SPY | Spies | Report on other Factions and incite rebellion |
| RBL | Rebels | Count of population willing to rebel; does not allocate the underlying units |

---

### 4.4 Weapons Units

Weapons units are used in combat. All weapons units are Production Units (manufactured by FACT-TL) unless noted otherwise.

| Code | Name | Description |
|---|---|---|
| EWPN | Energy Weapons | Concentrated energy beam; used in all combat except surface-to-surface |
| ESHD | Energy Shields | Deflect energy beams |
| MSSL | Missiles | Used in any combat; less accurate than energy weapons |
| MLNC | Missile Launchers | Launch MSSL; accuracy depends on launcher TL |
| AMSL | Anti-Missiles | Launched by MLNC to destroy incoming MSSL |
| ACFT | Assault Craft | Land/space vehicles used to invade Colonies or Ships |
| AWPN | Assault Weapons | Used by SOL on a planet's surface |
| MRBT | Military Robots | Replace SOL units (TL × 2 soldiers per MRBT unit) |
| MSUP | Military Supplies | Ammunition and medicines; consumed in combat (1 unit per SOL per combat round) |

---

### 4.5 Production Units

Production Units are manufactured by FACT-TL (all except FOOD) or FARM-TL (FOOD only), using Natural Resources and Labor.

#### Agriculture
| Code | Name | Description |
|---|---|---|
| FARM-TL | Farm | Produces FOOD. TL 1: open-colony; TL 2–5: hydroponic; TL 6–10: artificial-sunlight hydroponic. Requires Habitability (surface only). |
| FOOD | Food | Produced by FARM-TL; consumed by all Population Units (¼ FOOD unit per population unit per Turn) |

#### Extraction
| Code | Name | Description |
|---|---|---|
| MINE-TL | Mine | Extracts and refines Natural Resources from Deposits |

#### Industry
| Code | Name | Description |
|---|---|---|
| FACT-TL | Factory | Manufactures all units except Natural Resources, FOOD, and Population |
| AUTO-TL | Automation | Replaces USK workers in FACT, FARM, or MINE (unit × TL = worker units replaced) |
| CSGD-TL | Consumer Goods | Produced by FACT; used to pay the population |

#### Research
| Code | Name | Description |
|---|---|---|
| LAB-TL | Laboratory | Generates Research Points (RPs) each Turn |
| RP | Research Points | Currency for TL advancement; not a physical unit |
| PRTO-TL | Prototypes | Used to transfer Tech Levels between Entities |

#### Movement
| Code | Name | Description |
|---|---|---|
| SDRV-TL | Space Drive | Maintains orbit and manoeuvres in combat; cannot be used for interplanetary or interstellar travel |
| HDRV-TL | Hyper Drive | Propels Ships through hyperspace; jump range = 1 light-year × TL |

#### Support & Infrastructure
| Code | Name | Description |
|---|---|---|
| LSU-TL | Life Support | Recycles air and water; required on SHIP-TL, CENC-TL, and CORB-TL (supports TL² population units per LSU) |
| SENS-TL | Sensors | Detect Ships and Colonies in orbit; conduct probes (1 × TL probes per Turn) |
| PROB-TL | Probes | Executed by SENS; report on Ships, Colonies, and Deposits in the same System |
| STRC-TL | Structural | Required to build Ships and Colonies; STRC-2 variant is Light Structural |
| TRNS-TL | Transports | Transfer Population and materials between Ships/Colonies at the same Planet |

---

## 5. Research & Tech Level

**Tech Level (TL)** is an integer from 0 to 10 tracked per Entity (Ship or Colony). A Factory can only manufacture units up to its Colony's or Ship's current TL.

**Research Points (RPs)** are generated each Turn by LAB-TL units. RPs may be expended to increase an Entity's TL per the following schedule:

| Target TL | RPs Required |
|---|---|
| 2 | 100,000 |
| 3 | 200,000 |
| 4 | 400,000 |
| 5 | 800,000 |
| 6 | 1,600,000 |
| 7 | 3,200,000 |
| 8 | 6,400,000 |
| 9 | 12,800,000 |
| 10 | 25,600,000 |

TL may also be advanced by transferring technology from another Entity (via PRTO-TL) or purchasing it at a market or trade station.

---

## 6. Victory Conditions

### 6.1 Control

**Colony Control:** A Colony is controlled by a Faction if it contains at least 1 SOL or PRO Population Unit belonging to that Faction.

**Planet Control (Faction):** A Faction controls a Planet when it controls every Colony on or orbiting that Planet. At least one Colony must exist on the Planet for it to be considered controlled.

**Planet Control (Species):** A Species controls a Planet if it controls the majority of Colonies on or orbiting that Planet (at least ⌈C × 0.5⌉ + 1, where C = total Colonies on that Planet).

### 6.2 Solo Victory

Victory is evaluated against **habitable Planets** — those with a Habitability value greater than 0.

Let **H** = total number of habitable Planets in the Cluster.

A single Faction wins when:
- It controls at least ⌈H × 0.5⌉ + 1 habitable Planets, **and**
- No other single Faction controls more than ⌈H × 0.1⌉ + 1 habitable Planets.

### 6.3 Species Victory

Let **H** = total number of habitable Planets in the Cluster.

A Species wins when:
- It controls at least ⌈H × 0.5⌉ + 1 habitable Planets, **and**
- No other single Species controls more than ⌈H × 0.1⌉ + 1 habitable Planets.

> A Species controls a Colony if any Faction within the Species controls it.

---

## 7. Game Mechanics

### 7.1 Turns
The game progresses in discrete **Turns**. Each Turn:
1. Players submit **Orders** for their Entities.
2. The game engine processes all Orders.
3. The engine returns a **Turn Report** to each player.

### 7.2 Orders
Orders are instructions submitted by players, targeting Ships or Colonies. Specific order vocabulary to be documented.

### 7.3 Turn Report
The game engine's response to a Turn, returned to each player. Specific report structure to be documented.

---

## 8. Open Questions / To Be Documented

- [ ] Complete Order vocabulary
- [ ] Turn Report structure and format
- [ ] Ship and Colony construction rules (STRC requirements, mass)
- [ ] Combat mechanics (round structure, damage resolution)
- [ ] Market and trade station mechanics
- [ ] SPY and rebellion mechanics
