---
title: About Transport Tempo
linkTitle: Transport Tempo
weight: 10
---

Transports (TRNS) are described in the [Canonical Reference]({{< ref "/references/ec-canonical-reference.md#7-transports-trns" >}})
by a rate — `TL² × 200 MU per turn` — rather than by how much they hold. That
framing surprises people. Every other cargo idea in the game is a *capacity*: a
hold has a volume, a colony has a footprint, a life-support unit sustains so many
people. A transport, uniquely, is defined by a *throughput*. This article is about
why, and about what that difference means once you start planning turns around it.

## A turn is a span of time, not a moment

The key to transports is that a turn is *long*. It is not a single action; it is a
whole period of activity — long enough to fight a multi-round battle, retool
factories, run surveys, and move ships between planets. So when a transport "moves
800 MU in a turn," it is not lifting 800 MU at once. A TRNS-2's physical volume is
only about 12 VU. It moves 800 MU the way a ferry moves ten thousand cars a day:
by making the same short trip over and over.

Run the arithmetic and the picture is unmistakable. To shift 800 units of cargo in
12-VU loads, a TRNS-2 makes roughly 67 round trips — all inside one turn. A TRNS-10
hauling 20,000 units makes over 300. Nobody tracks those trips individually; the
game abstracts the whole shuttle campaign into a single number. But the number
*is* a count of trips, not the size of a hold, and that is the mental model to
carry: a transport is a rate of work, sustained across the length of a turn.

This is why the reference is careful to call it throughput and not capacity. Treat
it as a hold and you will badly overestimate what a single lift can do — and you
will misread what happens the moment a turn stops being a leisurely span, which is
exactly what combat is.

## Why the rate scales with the square of tech level

If a transport were just a bigger box at higher tech levels, its rate would grow
in step with its volume. It doesn't. Physical volume grows *linearly* with tech
level (`6 × TL`), but throughput grows with the *square* (`TL² × 200`). A TRNS-10
is only five times the size of a TRNS-2, yet it moves twenty-five times the cargo
per turn.

That gap is the interesting part, because it tells you what the tech level is
actually buying. It is not buying a bigger hold — that grows modestly. It is
buying *tempo*: faster drives, better navigation, more automation, quicker
turnaround at each end. A higher-TL transport doesn't carry much more per trip; it
simply completes far more trips in the same turn. One of the two factors of TL is
the slightly larger load; the other is the far greater number of runs. Multiply
them and you get the square.

You can read this as a deliberate design choice, and a defensible one. Linear
scaling would have made tech level a dull lever — twice the level, twice the
freight. Square scaling makes investment in transport technology compound, so a
logistics-focused faction pulls away from a neglectful one much faster than a
straight ratio would suggest. It rewards specialization. Whether that is the
*right* balance is a matter of taste, but the intent is clear: transport tech is
meant to feel like a force multiplier, not a linear upgrade.

## Combat collapses the tempo

The throughput model depends entirely on the turn being long. Combat is not. A
combat round is a compressed slice of time, and inside it a transport gets exactly
one trip: `3 × TL MU` per round. No shuttling, no dozens of runs — one load, one
run, and the round is over.

This is the same vessel obeying the same physics; only the available time has
changed. Out of combat, the turn is generous and the trips pile up. In combat, the
clock is brutal and the transport manages a single dash under fire. A TRNS-2 that
moved 800 MU across a peaceful turn moves just 6 MU in a combat round. That is not
a different rule bolted on — it is the throughput idea taken to its limit, where
the number of trips falls to one. The consistency is the point: the same underlying
model, evaluated over a long span or a compressed one, produces both figures.

It also explains why transports are combat-fragile in a way their peacetime numbers
disguise. A faction used to thinking of a high-TL transport as an ocean of capacity
will find it can barely move a soldier or two per round when the shooting starts —
which is why committed troops end up riding assault craft, and why the transport's
own crewing rules change under fire (the soldiers aboard operate it themselves).

## What this means when you plan a turn

A few threads worth pulling together, because the tempo idea touches more than the
cargo number:

- **Fuel tracks how hard the transport works, not how big it is.** Fuel burn is
  proportional to the fraction of throughput used — a half-loaded transport burns
  half the fuel. That only makes sense once you see the rate as trips made: fewer
  trips, less fuel. It is the throughput model showing up in a second system.
- **Don't plan a single heroic lift.** There is no "fill it up and go" — the whole
  point is many small trips. If you need cargo moved *now*, in a single moment, a
  transport is the wrong tool; its strength is sustained movement across the whole
  turn.
- **Buy tech level for tempo, not for the hold.** Because throughput scales with
  the square, upgrading transports pays back faster than almost any linear-scaling
  unit. If logistics is your bottleneck, TL is where the leverage is.

For the exact formulas, worked examples, and the crew and fuel figures, see
[§7 of the Canonical Reference]({{< ref "/references/ec-canonical-reference.md#7-transports-trns" >}}).
This article is only about *why* those numbers take the shape they do — a
transport is a rate of work stretched across the length of a turn, and everything
else follows from how much turn there is to stretch it across.
