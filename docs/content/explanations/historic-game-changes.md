---
title: Changes from the Historic Game
linkTitle: Changes from 1978
weight: 20
---

Epimethean Challenge began as a play-by-mail game first published in 1978. EC v7
is a modern re-implementation of that design, and most of its machinery is
faithful to the original: the same resources are mined, the same units are
manufactured, combat resolves in broadly the same phases. But a handful of things
were deliberately changed, and this article is about those changes and the
reasoning behind them.

It is not a specification — the [Canonical Reference]({{< ref "/references/ec-canonical-reference.md" >}})
describes the current game as it stands — and it is not a guide to migrating old
games. It exists to answer a single question: *why isn't EC v7 a byte-for-byte
reproduction of the 1978 rules?* The changes fall into three loose groups: things
that were **removed** because they no longer earned their keep, things that were
**renamed or restructured** for consistency and tooling, and things that were
**added** because they made the game clearer.

## Gold was removed

In the historic game, **gold** did double duty: it was a natural resource you
mined from deposits *and* the currency you paid wages and bought units with. That
coupling is the problem. Because the medium of exchange was itself a mined,
transportable, hoardable commodity, every economic decision became a logistics
decision. Players spent real effort shipping specie around, and a faction that
happened to be short on gold deposits found its whole development and expansion
throttled by something that had nothing to do with its industrial capacity. The
currency was cumbersome to manage and imposed a severe constraint on growth.

EC v7 drops gold entirely. Prices, wages, commissions, and fees are denominated in
an abstract currency with no mined counterpart (see
[§8, Economy & Markets]({{< ref "/references/ec-canonical-reference.md#8-economy--markets" >}})),
so metallics no longer have to be described as "all metallic substances *other
than gold*," and economic strength is no longer gated by where your gold happens
to be. This is not a novel departure, either: later versions of the historic game
reached the same conclusion and removed gold as well. We are following a path the
original design eventually took itself.

## Light structural units were removed

The historic game offered two structural products: an ordinary structural unit and
a cheaper, lighter variant that could only be manufactured in orbiting colonies.
The light variant was a special case — a unit whose cost and availability depended
on *where* it was built, in a way that no other unit's did. Special cases are
exactly the kind of thing that make a system harder to learn and harder to
maintain.

EC v7 removes the light variant and keeps a single structural family (STRC) that
scales by tech level like every other production unit, with the ordinary
location rules that govern the rest of the game — STRC-1 on the surface, STRC-2
through STRC-10 in orbit (see
[§4.5, Production Units]({{< ref "/references/ec-canonical-reference.md#45-production-units" >}})).
The goal was to make STRC behave consistently with everything else, so that
"structural" is one idea with one cost curve rather than two units with different
rules.

## Unit codes and terminology were standardized

The original rules used long English names, a scattering of ad-hoc abbreviations,
and vocabulary — *race*, *year* — inherited from a 1970s postal game. EC v7
regularizes all of it: every unit has a terse, fixed code (`FUEL`, `METL`,
`FACT-3`), production units carry a uniform `-TL` suffix, a *race* is now a
[**Faction** (and Species)]({{< ref "/references/ec-canonical-reference.md#2-factions--species" >}}),
and a *year* is now a **Turn**.

None of this changes how the game plays. It is an ergonomics decision, made to make
the maintainer's life easier: consistent codes are easier to hold in your head,
easier to grep for, and — importantly — far easier to parse reliably than a pile of
irregular names. That last point is not hypothetical; the order language now has a
formal grammar (the Lemon grammar, published in the
[References section]({{< ref "/references/_index.md" >}})), and a
regular vocabulary is what makes such a grammar tractable to write and validate.

## The order syntax was modernized

The historic game accepted orders as **positional, comma-separated** lines: the
meaning of each field came from its place in the sequence. That format is terse but
brittle. A misplaced comma silently shifts every field after it, the format resists
extension (where do you add a new optional argument?), and a parser has little to
check against.

EC v7 replaces it with a **keyword-tagged** format, where each value is introduced
by a label — `target`, `commit`, `produce` — and every line begins with its subject
(see [§11, Orders Reference]({{< ref "/references/ec-canonical-reference.md#11-orders-reference-canonical" >}})).
The change was made because the tagged form is easier for players to write and
easier for parsers to validate: orders are self-describing, new fields can be added
without disturbing old ones, and the leading subject gives the parser a reliable
sync point at the start of each line. The trade-off is verbosity — a keyword order
is longer than its positional equivalent — but that is a price well worth paying for
orders that fail loudly and legibly instead of silently doing the wrong thing.

## Cadres were introduced

The historic game had spies, construction workers, trainees, workers assigned to
factories, and a running tally of rebels — but it treated them as a loose
assortment of population-related entries with no unifying concept. In a turn report,
that scatters closely related information across unrelated categories.

EC v7 gathers them under a single idea: the **Cadre** (see
[§4.3, Cadre Units]({{< ref "/references/ec-canonical-reference.md#43-cadre-units" >}})).
A cadre is a *derived count* — population temporarily assigned to a role — rather
than a distinct kind of unit. The reason for naming the concept was report
clarity: "these are the people currently acting as workers / spies / construction"
reads as one coherent section instead of several disconnected ones. The framing
also makes the one true exception explicit — rebels (RBL) are *counted* but not
*allocated*, so they remain available for other assignments — which is much easier
to state cleanly once "cadre" is a defined term.

## Mass and volume were split

The original design leaned on a single notion of size. EC v7 tracks two distinct
quantities: **Mass (MU)** and **Volume (VU)**, where volume defaults to mass but is
allowed to differ (see
[§5, Mass & Volume]({{< ref "/references/ec-canonical-reference.md#5-mass--volume" >}})).
The split matters wherever a unit is light but bulky or dense but compact — food,
for instance, is 1 MU but 6 VU, so it weighs little yet fills a hold quickly. That
distinction is invisible until the two measures are separate, and once they are, it
feeds directly into things like how much a
[transport]({{< ref "transport-tempo.md" >}}) can shift per trip.

The two were separated to stay aligned with where later versions of the game take
the distinction, so that EC v7's data model does not have to be reworked to
accommodate it down the line.

## The Faction / Species distinction was added

Where the historic game identified a player by their *race*, EC v7 draws a
distinction the original did not: a **Faction** is the actor that gives orders,
while a **Species** is the group of all Factions that share a common home planet
(see [§2, Factions & Species]({{< ref "/references/ec-canonical-reference.md#2-factions--species" >}})).
This is more than a rename — it introduces a second, higher level at which control
and victory can be evaluated, so several factions of the same species can win
together.

The split was a concession to database modeling. Separating the actor (Faction)
from the group that shares a home planet (Species) gave the two concepts clean,
independent identities in the schema, rather than overloading a single "race"
record with both roles. The gameplay consequence — species-level control and
victory — followed naturally once the data model drew that line.

---

This list is not closed. As more of the historic design is reconciled against the
implementation, further deltas may surface and be added here. For the current state
of any mechanic — as opposed to the story of how it got that way — the
[Canonical Reference]({{< ref "/references/ec-canonical-reference.md" >}}) remains
the authority.
