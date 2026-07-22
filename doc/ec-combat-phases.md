# Combat Phases

---

## Pre-Combat (Round 1 Only)

1. **Troop Deployment** — Committed soldiers are armed and loaded into assault craft; overflow goes into transports with assault weapons; remainder stays behind as uncommitted defense.

---

## Each Round

2. **Weapons Fire** — All attacks execute simultaneously:
   - Energy beams fire at missiles, transports, and assault craft
   - Anti-missiles fire at incoming missiles, transports, and assault craft
   - Missiles are launched at ships/colonies
   - Defender fires back automatically (energy beams, anti-missiles, missiles)

3. **Intercept Resolution** — Determine what gets shot down before reaching the target:
   - Anti-missiles vs. incoming missiles
   - Energy beams vs. transports and assault craft (troop transport casualties)

4. **Casualty Calculation** — For raids/invasions: calculate combat factors and apply losses to both attacker and defender; split into KIA (70%) and wounded (30%).

5. **Damage Calculation** — For bombardment: apply damage from missiles and energy beams that weren't intercepted; distribute to weapons/drives (75%) and other units (25%).

6. **Surrender Check** — Check if any side faces 6:1 odds → auto-surrender.

7. **Ship Movement** — Ships with bombard orders move closer; ships under attack without bombard orders move away.

---

## End of Round

8. **Continuation Check** — Combat continues if: mission not complete AND soldiers/fuel/military supplies/missiles not exhausted; otherwise combat ends.

---

## Post-Combat

9. **Capture Resolution** — If all defenders are destroyed or surrendered, attacker takes control; for raids, stolen units are transferred.

---

## Notes

- **Raids are single-round only** — they skip the continuation check and end after round 1.
- **Troop deployment (step 1) only happens in round 1**, but feeds into every subsequent round's state.
- **All combat orders execute simultaneously** — there is no attacker-goes-first ordering within a round.
- **Intercept resolution (step 3)** is implicit in the source text rather than an explicitly named phase — may be modelled as a sub-step of weapons fire or kept separate depending on damage sequencing needs.
