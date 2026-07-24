# AGENTS.md

## Project

EC v7 implements Epimethean Challenge, a 4X space game played via email. We
are building command-line tools for operating and playing the game.

## GitHub workflow

- Assign every new GitHub issue and pull request to `@mdhender`.
- Work on a feature branch and push it for review and safety. Do not commit
  directly to `main`.
- Before the final commit and push to upstream, bump the version in
  `version.go`:
  - Bump the minor version when adding, updating, or changing a feature.
  - Bump the patch version when fixing a bug, cleaning up code, or changing
    documentation under the `doc/` or `docs/` paths.

## Repository layout

- `version.go` is the shared version source for every CLI binary in this
  repository.
- Keep `version.go` as the only Go source file in the repository root. Put CLI
  entry points under `cmd/<tool>/` and supporting packages under `internal/`.
- Repository metadata and project directories such as `docs/`, `games/`, and
  configuration files may remain at the root.
- `docs/` is the authoritative source for game documentation.
- `doc/` is a working directory for maintainer documentation and is not
  authoritative.

## Command-line and database packages

- Use `github.com/peterbourgon/ff/v4` for command-line parsing and utilities.
  Do not introduce Cobra.
- Use ZombieZen's SQLite packages for the database driver and migrations
  (`zombiezen.com/go/sqlite` and `zombiezen.com/go/sqlite/sqlitemigration`).
- Keep the project CGO-free. Do not introduce SQLite drivers or other
  dependencies that require CGO.

## Randomness

- Use `math/rand/v2`, never the legacy `math/rand` package.
- Game randomness must derive from `internal/prng` addressed streams. Do not
  draw from ambient sources such as wall-clock time, package-level rand
  functions, or map iteration order.

## Test data

- Prefer temporary directories (for example, `t.TempDir()`) for test data.
- When test data must persist in the repository, keep it in the directory for
  the agent running the test:
  - Amp: `games/amp/`
  - Claude: `games/claude/`
  - Codex: `games/codex/`
- `games/alpha/` belongs to the user. Do not modify, delete, overwrite, or rely
  on its contents in automated tests. Avoid changes that would break the
  user's testing workflow there.
