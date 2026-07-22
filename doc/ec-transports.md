# Epimethean Challenge — Transports (TRNS)

## Overview

Transports (TRNS) transfer population and materials between ships and colonies at the same planet. They are not weapons and carry no armament.

---

## Stats (per unit)

| Stat         | Formula       | TRNS-2 Example |
|--------------|---------------|----------------|
| METL Input   | `3 × TL`      | 6              |
| NMTL Input   | `1 × TL`      | 2              |
| Mass (MU)    | `4 × TL`      | 8              |
| Volume (VU)  | `6 × TL`      | 12             |
| Operational  | No            | —              |

---

## Throughput vs. Capacity

### Normal Operations

The transfer rate of a TRNS unit is:

> **TL² × 200 MU per turn**

This is a **throughput rate**, not a cargo capacity. A TRNS-2 moves 800 MU per turn by making multiple trips within the turn, not by carrying 800 MU at once. Its physical volume is only 12 VU, so each trip carries a batch of up to 12 VU of cargo. The length of a turn is long enough to accommodate many such trips.

### Combat Operations

In combat, only a **single trip** is possible. The carrying capacity is:

> **3 × TL MU per combat round**

This represents one load, one run — consistent with the compressed timescale of a combat round.

---

## Worked Example: TRNS-2 moving USK

**Scenario:** A TRNS-2 is tasked with transferring 800 units of Unskilled Workers (USK) from a surface colony to an orbiting colony.

**Unit stats:**
- TRNS-2 Volume: `6 × 2 = 12 VU`
- USK Volume: 1 VU per unit (1 MU mass, 1 VU volume)
- 800 USK = 800 VU total

**Apparent contradiction:**  
800 VU of USK cannot fit inside a 12 VU transport. Yet the rules say a TRNS-2 moves 800 MU per turn.

**Resolution:**  
The 800 MU/turn figure is a rate, not a single-lift capacity. The TRNS-2 shuttles in batches:

- Batch size: ~12 VU (one full load)
- Trips required: 800 ÷ 12 ≈ **67 trips**
- All 67 trips occur within a single turn

The turn is long enough for the transport to complete all 67 round trips. The 800 MU/turn throughput is the aggregate result.

**In combat**, the same TRNS-2 can deliver only:

> `3 × 2 = 6 MU` of soldiers per combat round

One trip. No time for a second run.

---

## Worked Example: TRNS-10 moving USK

**Scenario:** A TRNS-10 is tasked with transferring 20,000 units of Unskilled Workers (USK) from a surface colony to an orbiting colony.

**Unit stats:**
- TRNS-10 Volume: `6 × 10 = 60 VU`
- USK Volume: 1 VU per unit
- 20,000 USK = 20,000 VU total
- Throughput: `10² × 200 = 20,000 MU/turn`

**Apparent contradiction:**  
20,000 VU of USK is 333× the volume of the transport itself. The ferrying explanation is the same as for TRNS-2 — the transport makes repeated trips within the turn.

**Resolution:**  
- Batch size: 60 VU (one full load)
- Trips required: 20,000 ÷ 60 ≈ **333 trips**

This is where the model is under stress. 333 round trips in a single turn is a significant demand on the abstraction. The turn must be long enough to accommodate not just the travel time per trip but also the loading and unloading of each batch. At this scale, the "ferrying within a turn" explanation is still technically valid, but it implies a very high operational tempo — the TRNS-10 is essentially running a continuous shuttle service for the entire turn.

The TL-10 rating justifies this: a TL-10 transport is a highly advanced craft, presumably faster and more automated than lower-TL equivalents. The throughput increase (TL² scaling) reflects this — a TRNS-10 moves 25× as much per turn as a TRNS-2, despite having only 5× the volume.

**In combat**, the TRNS-10 delivers:

> `3 × 10 = 30 MU` of soldiers per combat round

One trip. The same constraint applies regardless of TL.

---

## Crew

When not in combat, a TRNS unit is operated by professionals: **1 PRO unit per 10 TRNS units**.

In combat, the soldiers being carried operate the transport themselves.

---

## Fuel Use

Fuel consumed per turn:

> **TL² ÷ 10**, proportional to the percentage of capacity actually used.

A TRNS-2 running at full throughput uses `4 ÷ 10 = 0.4 fuel units` per turn. A half-loaded TRNS-2 uses 0.2.

In combat, fuel use is `0.01 × TL²` per round trip.
