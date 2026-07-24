# Inventory State Decision

This maintainer reference fixes how an entity's inventory distinguishes **stored**
(unassembled) units from **assembled** operating units, and which transitions
between those states are legal. It records the decision for GitHub issue #26
("Define stored and assembled unit state") and closes its inventory-shape,
transition-invariant, and worked-example acceptance criteria.

The player-facing statements of these facts live in the canonical reference:
unassembled/assembled state and transfer eligibility in
[§4.14 Unit State](../docs/content/references/ec-canonical-reference.md#414-unit-state),
and the group-numbering rule in §4.15 Production Groups. This document is the
maintainer decision behind them.

It builds on two existing decisions and does not restate them:

- The [numeric-types decision](ec-numeric-types.md) defines `Quantity` (a
  non-negative count), the placeholder `InventoryUnit`, and `UnitDefinition`
  (per-unit mass, stored volume, and `RequiresAssembly`). It delegated the final
  inventory shape and the stored/assembled transition invariants to this issue.
- The [unit-spec decision](ec-unit-spec.md) defines `UnitSpec` and its equality
  and inventory-key rules. IK-4 states that the stored/assembled distinction is a
  separate dimension of an inventory line, not part of `UnitSpec` equality or the
  key; US-3/IK-3 state that Population, Cadre, and Research Points are not
  `UnitSpec` inventory.

## 1. Inventory line

An inventory holds quantities against unit identities. This decision fixes the
shape of the line held for a `UnitSpec`, refining the placeholder
single-`Quantity` `InventoryUnit` from the numeric-types decision into two
quantity dimensions:

```go
type InventoryUnit struct {
    Spec      UnitSpec // unit identity and inventory key (IK-1 of the unit-spec decision)
    Stored    Quantity // on hand and unassembled
    Assembled Quantity // assembled and operating this turn
}
```

`Stored` and `Assembled` are the two dimensions of one line for a single `Spec`
(IK-4). Both are `Quantity`, so both are non-negative. Per-unit mass, stored
volume, and `RequiresAssembly` remain in `UnitDefinition`, not on the line; the
line carries counts only.

A **spec is operational** when its `UnitDefinition.RequiresAssembly` is `true`
(equivalently, **Assembled?** = `yes`, the **Op.** column of the canonical
mass/volume table). "Operational" names the *requirement to be assembled before
it functions*, not activity in general: `TRNS` transports without ever being
assembled and is therefore **non-operational** here.

## 2. State invariants

- **IS-1 (one line per spec)** A `UnitSpec`-keyed inventory holds at most one
  `InventoryUnit` per `Spec`. That single line carries both `Stored` and
  `Assembled`; the same `Spec` is never split across two competing lines (IK-4).
  A spec absent from the inventory denotes `Stored = 0` and `Assembled = 0`.
- **IS-2 (assembled requires an operational spec)** `Assembled > 0` is
  representable only when `Spec` is operational
  (`UnitDefinition.RequiresAssembly == true`). For every non-operational spec —
  resources, commodities, ammunition, `PROB`, `PRTO`, `ACFT`, `AWPN`, `MRBT`,
  `TRNS`, and the like — `Assembled` is always `0`.
- **IS-3 (non-negative counts)** `Stored >= 0` and `Assembled >= 0` (guaranteed
  by `Quantity`). The line's **total on hand** is `Stored + Assembled`; there is
  no separate stored count for the assembled portion.
- **IS-4 (non-spec kinds have no assembled state)** Population, Cadre, and
  Research Points are not `InventoryUnit` lines (US-3, IK-3) and therefore have no
  `Assembled` dimension. Population and Cadre are `PopulationClass` quantities;
  `RP` is a scalar balance. "Assembled population" is not representable.

IS-2 is the core of "no impossible combinations": the assembled dimension exists
only where assembly applies. Construction should keep an illegal line
unrepresentable rather than merely rejected — for example, admit `Assembled` only
through operations gated on an operational spec (§3), so a non-operational spec
can never be handed a non-zero assembled count.

## 3. Assembly, disassembly, and transfer transitions

The `ASSEMBLE` and `DISASSEMBLE` *orders* are deferred (their stages appear in the
turn sequence, but their order syntax is owned elsewhere). The state transitions
they drive are fixed here so the inventory type is closed under them.

- **AT-1 (assemble)** Assembling `q` units of an operational `Spec` sets
  `Stored -= q` and `Assembled += q`. It requires `Spec` operational (IS-2) and
  `q <= Stored`.
- **AT-2 (disassemble)** Disassembling `q` units sets `Assembled -= q` and
  `Stored += q`. It requires `q <= Assembled`.
- **AT-3 (conservation)** AT-1 and AT-2 conserve the line total
  (`Stored + Assembled`); assembly state changes never create or destroy units.
  Only production, transfer, combat, and consumption change a line total.
- **AT-4 (transfer eligibility)** `TRANSFER` and the `SETUP … TRANSFER` block move
  units out of `Stored` only. An assembled unit is not transferable; it must be
  disassembled (AT-2) first. A transfer of `q` decrements the source's `Stored`
  by `q` and increments the destination's `Stored` by `q`, subject to `q <=`
  source `Stored`. Non-operational manufactured units (e.g. `TRNS`) are always in
  `Stored` and are transferable without any disassembly step.

## 4. Worked examples

These cover the four required cases — a transferable resource, population, `TRNS`,
and an assembled operational unit — plus the impossible combinations IS-1–IS-4
exclude. Domains are from the Unit-Definition Catalog in
[`ec-unit-tables.md`](ec-unit-tables.md).

### 4.1 Required cases

| Case                        | Representation                                          | Transferable?                       | Rules            |
|-----------------------------|--------------------------------------------------------|-------------------------------------|------------------|
| Resource `FUEL`             | `InventoryUnit{(FUEL,0), Stored: 5000, Assembled: 0}`  | Yes — all 5000 (from `Stored`)      | IS-2, AT-4       |
| Population `USK`            | `PopulationClass` quantity `20000` — not `InventoryUnit`| Yes, as population (via `TRNS`)     | IS-4 (US-3, IK-3)|
| Transport `TRNS-4`          | `InventoryUnit{(TRNS,4), Stored: 12, Assembled: 0}`    | Yes — all 12; never assembled       | IS-2, AT-4       |
| Factory `FACT-3` (operating)| `InventoryUnit{(FACT,3), Stored: 2, Assembled: 2}`     | Only the 2 in `Stored`              | IS-2, AT-1, AT-4 |

Walk-through for `FACT-3`. A colony manufactures 4 FACT-3, which arrive
unassembled: `{Stored: 4, Assembled: 0}` (factory output is unassembled). An
`ASSEMBLE 2` moves two into operation: `{Stored: 2, Assembled: 2}` (AT-1, total 4
conserved). At transfer time only the two in `Stored` may be moved (AT-4); moving
the operating pair first requires `DISASSEMBLE 2` back to `{Stored: 4,
Assembled: 0}` (AT-2).

### 4.2 Impossible combinations (excluded by construction)

| Attempted line / state                                   | Why it is invalid                                       | Rule |
|----------------------------------------------------------|---------------------------------------------------------|------|
| `InventoryUnit{(FUEL,0), Stored: 0, Assembled: 100}`     | `FUEL` is non-operational; `Assembled` must be `0`      | IS-2 |
| `InventoryUnit{(TRNS,4), Stored: 0, Assembled: 12}`      | `TRNS` never requires assembly; no assembled state      | IS-2 |
| Two lines both keyed `(FACT,3)`                          | at most one line per `Spec`                             | IS-1 |
| `InventoryUnit{(FACT,3), Stored: -1, Assembled: 0}`      | counts are non-negative                                 | IS-3 |
| An `InventoryUnit` keyed by `USK` / `SOL` / `RP`         | population, cadre, and `RP` are not `UnitSpec` inventory | IS-4 |
| `TRANSFER` drawn from `Assembled`                        | transfer moves `Stored` only                            | AT-4 |

## 5. Decision status and deferred questions

This document is the recorded maintainer decision for inventory state. The
decisions it fixes are:

- the two-dimension inventory line (`Stored`, `Assembled`) refining the
  placeholder `InventoryUnit` (§1);
- the state invariants that exclude impossible combinations (§2);
- the assembly, disassembly, and transfer transition invariants (§3); and
- worked examples for a resource, population, `TRNS`, and an assembled unit (§4).

The following remain deferred; none blocks the invariants above:

- **`ASSEMBLE` / `DISASSEMBLE` / `TRANSFER` order syntax and processing.** Only the
  *state transitions* are fixed here (§3); the order-language details and their
  stage processing are owned by the orders work.
- **Stored vs. operational volume values.** Whether an assembled unit occupies a
  different volume than the same unit stored, and by how much, is unverified and
  isolated in the canonical facts. This decision does not invent an
  operational-volume rule; occupied volume remains a derived calculation over
  `UnitDefinition` (numeric-types decision).
- **Transfer capacity mechanics.** How a `TRANSFER` consumes `TRNS` throughput,
  crew, and fuel is a separate rule; AT-4 fixes only *eligibility* (source
  `Stored`), not the per-turn capacity that bounds `q`.
