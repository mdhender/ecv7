---
title: Engine Determinism Reference
linkTitle: Determinism
weight: 20
---

Technical description of the EC engine's determinism machinery, implemented by
the `internal/prng` package. Given the same master seeds, a game produces
identical outcomes on any machine, independent of the order draws are made and
of Go map iteration order.

---

## 1. Master Seeds

Every game has exactly two master seeds.

| Property | Value                                             |
|----------|---------------------------------------------------|
| Name     | `seed_1`, `seed_2`                                |
| Type     | unsigned 64-bit integer (`uint64`)                |
| Storage  | `game` table, columns `seed_1` and `seed_2`       |
| Scope    | One pair per game, fixed at game creation         |

The master seeds are the root of all randomness in a game. They are never
consumed or advanced; every random value in the game is a pure function of the
seeds and a draw address.

---

## 2. Draw Addresses (Key Paths)

Every random draw has an address: a path of `Key` values.

| Property     | Value                                                        |
|--------------|--------------------------------------------------------------|
| Element type | `Key` (signed 64-bit integer, `int64`)                       |
| First element| A domain tag (§4)                                            |
| Remainder    | Instance keys identifying the specific object (§4, §5)       |
| Length       | Part of the address: `[K, q]` and `[K, q, r]` are distinct   |

The same address always yields the same stream. When the address is computed,
and in what order relative to other addresses, has no effect on its output.

**Instance keys must be intrinsic to the game.** Map objects are addressed by
their canonical coordinates ([§1.5 of the Canonical
Reference]({{< ref "ec-canonical-reference.md#15-identity-and-coordinate-display" >}})).
Never address a draw by a SQLite autoincrement row ID: row IDs depend on
insertion order, which is not part of the game's definition.

---

## 3. Stream Derivation

A stream is derived from the seeds and an address by one hash construction:

```
stream = PCG( SHA-256( seed1 ‖ seed2 ‖ len(path) ‖ path[0] ‖ path[1] ‖ … ) )
```

| Step          | Detail                                                        |
|---------------|---------------------------------------------------------------|
| Encoding      | Every element written as 8 bytes, big-endian                  |
| Length prefix | `len(path)` as `int64`, written before the path elements      |
| Coercion      | `Key` (`int64`) values are bit-cast to `uint64` for encoding  |
| Digest use    | Digest bytes 0–7 and 8–15 become the two PCG seed words       |
| Generator     | `math/rand/v2` PCG                                            |

The identical hash serves both operations in §6; only the destination of the
first 128 bits differs.

---

## 4. Domain Tags

The first element of every path is a domain tag. The registry is append-only
and defined in `internal/prng/tags.go`.

| Tag           | Value | Instance keys after the tag              |
|---------------|-------|------------------------------------------|
| `TagCluster`  | 1     | *(none — cluster generation)*            |
| `TagStellium` | 2     | `x, y, z`                                |
| `TagSystem`   | 3     | `x, y, z, seq`                           |
| `TagPlanet`   | 4     | `x, y, z, seq, orbit`                    |
| `TagDeposit`  | 5     | `x, y, z, seq, orbit, deposit_no`        |
| `TagPlayer`   | 6     | `player_no`                              |

Tag value 0 is invalid and never used.

---

## 5. Instance Key Domains

| Key          | Domain                                                        |
|--------------|---------------------------------------------------------------|
| `x`, `y`, `z`| Stellium coordinates, integers −15 to 15                      |
| `seq`        | System sequence letter, coerced `A` = 1, `B` = 2, …           |
| `orbit`      | Planet orbit, 1 to 10                                         |
| `deposit_no` | Deposit number on the planet, sequential 1 to 40; assigned at generation, never reused or renumbered |
| `player_no`  | Game-assigned player number; never a database row ID          |

---

## 6. Operations

Defined on `Seeds` in `internal/prng`.

| Operation          | Returns   | Behavior                                                         |
|--------------------|-----------|------------------------------------------------------------------|
| `Stream(path...)`  | `*Stream` | Draw source at the address; implements `math/rand/v2.Source`     |
| `Derive(path...)`  | `Seeds`   | Child seed pair at the address, for a subsystem's own randomness |
| `Roller(path...)`  | `*Roller` | Dice-rolling front end over the stream at the address            |

`Roller` methods, all drawing from the one underlying stream in order:

| Method                | Result                                                     |
|-----------------------|------------------------------------------------------------|
| `RollN(n, sides)`     | Sum of `n` uniform dice in `[1, sides]`; result in `[n, n × sides]` |
| `RollRange(lo, hi)`   | One uniform integer in `[lo, hi]` inclusive; requires `lo < hi` |
| `Shuffle(n, swap)`    | Pseudo-random reordering of `n` elements                   |
| `Perm(n)`             | Pseudo-random permutation of `[0, n)`                      |

Build a `Roller` once per address and reuse it. Building a fresh `Roller` per
die restarts the stream and correlates rolls.

---

## 7. Frozen Surfaces

The following are compatibility surfaces. Once any game exists, changing them
rewrites that game's outcomes; they must never change.

| Surface        | Frozen elements                                              |
|----------------|--------------------------------------------------------------|
| Path encoding  | Element order, `int64`/`uint64` coercions, big-endian layout, length prefix |
| Tag registry   | Numbering is append-only; never insert, remove, or reorder   |
| Generator      | SHA-256 digest layout and `math/rand/v2` PCG                 |
| Golden vectors | `internal/prng/testdata/golden.json` pins `(seed1, seed2, path) → output`; tests fail on any drift |

Regenerate golden vectors (`go test ./internal/prng/ -update`) only when
intentionally establishing a new surface before any game exists — never to
make a failing test pass.

---

## 8. Prohibitions

- Never seed a path element from a SQLite autoincrement row ID.
- Never draw from ambient sources: wall-clock time, package-level
  `math/rand/v2` functions, map iteration order, or goroutine scheduling.
- Never use the legacy `math/rand` package; the engine uses `math/rand/v2`
  exclusively.
- Never renumber deposits or players after generation.
- Never insert or reorder domain tags; append only.
