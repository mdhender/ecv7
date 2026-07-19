# `ec`

`ec` is the command-line application for playing EC.

```text
ec [--quiet] [--path PATH] <SUBCOMMAND>
```

Use `ec --help` to display command help.

## Environment variables and `.env` files

`ec` reads environment variables whose names start with `EC_`. A flag's
environment variable name is its long name in uppercase, with hyphens replaced
by underscores. For example, `EC_PATH` supplies `--path` and `EC_QUIET=true`
supplies `--quiet`. Explicit command-line flags take precedence over environment
variables.

Before parsing flags, `ec` loads dotenv files from the current working
directory. Set the exported `EC_ENV` variable to select `development`, `test`,
`production`, or `agent`; it defaults to `development`. Files are considered
from highest to lowest priority:

1. Variables already exported in the process environment.
2. `.env.<environment>.local`
3. `.env.local`
4. `.env.<environment>`
5. `.env`
6. Command defaults

Missing dotenv files are ignored. Files with `.local` in their names are
ignored by Git and may contain local settings or secrets.

## Root options

- `--path PATH` identifies the directory containing `ec.db`. It defaults to
  `db` and is available to every subcommand.
- `--quiet` suppresses status and diagnostic output from commands that support
  it. It does not suppress intentional version output.

## Print the application version

```text
ec version [--build | --long]
```

With no flag, `ec version` prints the core semantic version. `--build` also
includes pre-release information, while `--long` includes both pre-release and
build information. The two flags are mutually exclusive.

```sh
ec version
ec version --build
ec version --long
```
