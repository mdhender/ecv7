# Epimethean Challenge — Unit Mass & Volume

**Definitions:**
- **Mass (MU):** Physical mass in mass units (1 MU ≈ 17,000 lbs / 100 men). Sourced from document tables where stated; otherwise derived as the sum of METL + NMTL inputs.
- **Volume (VU):** Storage volume in volume units. Defaults to Mass (MU) unless otherwise specified.
- **Operational:** Unit must be assembled before it can be used; dis-assembly required before transfer.
- **TL:** Technological level of the unit (1–10).

---

## Natural Resources & Food

These units are not produced by Factories from METL/NMTL inputs. Natural resources are extracted by Mines; Food is produced by Farms. No input costs apply.

| Code | Unit        | Source  | Mass (MU) | Volume (VU) | Operational |
|------|-------------|---------|-----------|-------------|-------------|
| FUEL | Fuel          | Mined  | 1   | 1   | No |
| METL | Metallics     | Mined  | 1   | 0.5 | No |
| NMTL | Non-Metallics | Mined  | 1   | 0.5 | No |
| FOOD | Food          | Farmed | 1   | 6   | No |

---

## Population Units

Population units are not manufactured by Factories. They grow and shrink each turn through demographic forces. We are arbitrarily defining 1 MU ≈ the mass of 1 Population Unit.

| Code | Unit              | Mass (MU) | Volume (VU) | Operational |
|------|-------------------|-----------|-------------|-------------|
| UEM  | Unemployables     | 1         | 1           | No          |
| USK  | Unskilled Workers | 1         | 1           | No          |
| PRO  | Professionals     | 1         | 1           | No          |
| SOL  | Soldiers          | 1         | 1           | No          |

---

## Weapons

| Code | Unit             | METL Input    | NMTL Input    | Mass (MU)          | Volume (VU)        | Operational |
|------|------------------|---------------|---------------|--------------------|--------------------|-------------|
| AWPN | Assault Weapons  | `1 × TL`      | `1 × TL`      | `2 × TL`           | 20                 | No          |
| ACFT | Assault Craft    | `3 × TL`      | `2 × TL`      | `5 × TL`           | `5 × TL`           | No          |
| AMSL | Anti-Missiles    | `2 × TL`      | `2 × TL`      | `4 × TL`           | `4 × TL`           | No          |
| ESHD | Energy Shields   | `25 × TL`     | `25 × TL`     | `50 × TL`          | `50 × TL`          | Yes         |
| EWPN | Energy Weapons   | `5 × TL`      | `5 × TL`      | `10 × TL`          | `10 × TL`          | Yes         |
| MRBT | Military Robots  | `10 + TL`     | `10 + TL`     | `(2 × TL) + 20`    | `(2 × TL) + 20`    | No          |
| MSUP | Military Supplies| 0.02          | 0.02          | 0.04               | 0.04               | No          |
| MSSL | Missiles         | `2 × TL`      | `2 × TL`      | `4 × TL`           | `4 × TL`           | No          |
| MLNC | Missile Launchers| `15 × TL`     | `10 × TL`     | `25 × TL`          | `25 × TL`          | Yes         |

## Production

| Code | Unit      | METL Input | NMTL Input | Mass (MU)          | Volume (VU)        | Operational |
|------|-----------|------------|------------|--------------------|--------------------|-------------|
| FACT | Factories | `8 + TL`   | `4 + TL`   | `12 + (2 × TL)`    | `12 + (2 × TL)`    | Yes         |
| FARM | Farms     | `4 + TL`   | `2 + TL`   | `6 + TL`           | `6 + TL`           | Yes         |
| MINE | Mines     | `5 + TL`   | `5 + TL`   | `10 + (2 × TL)`    | `10 + (2 × TL)`    | Yes         |

## Miscellaneous

| Code   | Unit             | METL Input  | NMTL Input  | Mass (MU)    | Volume (VU)  | Operational |
|--------|------------------|-------------|-------------|--------------|--------------|-------------|
| AUTO   | Automation       | `2 × TL`    | `2 × TL`    | `4 × TL`     | `4 × TL`     | Yes         |
| CSGD   | Consumer Goods   | 0.2         | 0.4         | 0.6          | 1.0          | No          |
| HDRV   | Hyper Engines    | `25 × TL`   | `20 × TL`   | `45 × TL`    | `60 × TL`    | Yes         |
| LSU    | Life Support     | `3 × TL`    | `5 × TL`    | `8 × TL`     | `12 × TL`    | Yes         |
| SENS   | Sensors          | `10 × TL`   | `20 × TL`   | `30 × TL`    | `40 × TL`    | Yes         |
| SDRV   | Space Drives     | `15 × TL`   | `10 × TL`   | `25 × TL`    | `33 × TL`    | Yes         |
| STRC   | Structural       | 0.1         | 0.4         | 0.5          | 0.5          | Yes         |
| TRNS   | Transports       | `3 × TL`    | `1 × TL`    | `4 × TL`     | `6 × TL`     | No          |

---

## Notes

**Mass derivation rule:** Where mass is not explicitly stated in the document, it is derived as the sum of METL + NMTL inputs.

**Volume derivation rule:** Volume defaults to Mass unless a distinct value is specified.

**Operational units**: Space Drives, Sensors, Automation, Life Support, Energy Weapons, Energy Shields, Mining Units, Factories, Farms, Hyper Engines, Structural, Missile Launchers. These must be assembled after being taken out of storage to function.

---

## Addendum: Volume — Storage vs. Operational

The Volume (VU) values in this document currently conflate two distinct concepts that require further research and correction:

- **Storage volume:** The space a unit occupies when dis-assembled and held in storage.
- **Operational volume:** The space a unit occupies when assembled and in active use.

The upstream document notes (§5.1) that units in storage count as only ½ their mass for structural housing purposes, suggesting storage and operational volume are not equivalent. The current VU column does not yet distinguish between these two states. This requires further research and will be updated in a future revision.
