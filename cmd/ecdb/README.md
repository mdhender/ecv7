# `ecdb`

`ecdb` creates and maintains the persistent SQLite database used by EC. The
database is always named `ec.db`; path arguments identify the directory that
contains the file, not the file itself.

```text
ecdb [--quiet] <SUBCOMMAND>
ecdb [--quiet] database <SUBCOMMAND>
```

Use `ecdb --help` or `ecdb database --help` to display command help.

## Environment variables and `.env` files

`ecdb` reads environment variables whose names start with `EC_`. A flag's
environment variable name is its long name in uppercase, with hyphens replaced
by underscores. For example:

- `EC_PATH` supplies `--path`.
- `EC_OUTPUT_PATH` supplies `--output-path`.
- `EC_QUIET=true` supplies `--quiet`.
- `EC_VERSION=true`, when running `database backup`, supplies `--version`.

An explicit command-line flag takes precedence over the corresponding
environment variable. Environment variables only apply to flags available for
the command being run. Thus, `EC_VERSION` controls `database backup`; it does
not affect the separate `ecdb version` command.

Before parsing flags, `ecdb` loads dotenv files from the current working
directory. Set the exported `EC_ENV` variable to select an environment:

```sh
export EC_ENV=production
ecdb database verify
```

The supported values are `development`, `test`, `production`, and `agent`.
When `EC_ENV` is not exported, it defaults to `development`. Because the
environment is selected before dotenv files are loaded, put `EC_ENV` in the
process environment rather than in a dotenv file.

For the selected environment, files are considered in this order, from
highest to lowest priority:

1. Variables already exported in the process environment.
2. `.env.<environment>.local`
3. `.env.local`
4. `.env.<environment>`
5. `.env`
6. Command defaults

Missing dotenv files are ignored. Files with `.local` in their names are
ignored by Git and may contain local settings or secrets; `.env` and
`.env.<environment>` are shared files and must not contain secrets.

For example, this `.env.development.local` supplies the database directory for
development commands:

```dotenv
EC_PATH=games/example/db
```

Then the path may be omitted:

```sh
ecdb database verify
```

An explicit flag still wins:

```sh
ecdb database verify --path games/other/db
```

## Important limitation of `verify`

**`ecdb database verify` does not verify database contents or database
integrity.** It checks only these two SQLite metadata values:

1. `application_id` identifies the file as an EC database.
2. `user_version` matches the schema version expected by this build of
   `ecdb`.

It does **not** inspect table contents or schema objects, check foreign-key
consistency, or run a SQLite integrity check. A successful result means only
that the two metadata values match.

## Database lifecycle

### Create a database

```text
ecdb database create [--path PATH]
```

Creates `ec.db` and applies all migrations. `PATH` must already be an existing
directory. It defaults to `db`, relative to the current directory. The command
fails rather than replacing an existing `ec.db` and produces no output on
success.

```sh
mkdir -p games/example/db
ecdb database create --path games/example/db
```

### Upgrade a database

```text
ecdb database upgrade --path PATH
```

Applies migrations missing from an existing `PATH/ec.db`. `--path` is
required, and the directory and database must already exist. The command
rejects a file that is not identified as an EC database and a database whose
schema is newer than this build supports.

On success, the command reports whether migrations were applied and prints
the resulting schema version:

```sh
ecdb database upgrade --path games/example/db
# migrations applied to games/example/db (version 1)
```

Use the root `--quiet` flag to suppress this status line:

```sh
ecdb --quiet database upgrade --path games/example/db
```

### Verify database metadata

```text
ecdb database verify --path PATH
```

Opens an existing `PATH/ec.db` read-only and checks only its SQLite
`application_id` and `user_version`, as described in the
[verification limitation](#important-limitation-of-verify). `--path` is
required. Success produces no output.

```sh
ecdb database verify --path games/example/db
```

On failure, `verify` exits nonzero and writes a diagnostic that includes the
path when one was supplied. `--quiet` suppresses that diagnostic but does not
change the exit status:

```sh
if ecdb --quiet database verify --path games/example/db; then
    echo "EC database metadata matches this build"
fi
```

## Supporting database commands

### Back up a database

```text
ecdb database backup --path PATH [--output-path OUTPUT_PATH] [--version]
```

Writes a consistent, compacted backup of `PATH/ec.db`. The required `PATH`
and the output directory must already exist. `--output-path` defaults to
`PATH`. The backup is named with a UTC timestamp, for example
`ec.db.20260719T142530Z`; `--version` appends the schema version, for example
`ec.db.20260719T142530Z-1`.

The source must be an EC database at the schema version expected by this
build. On success, the command writes the path of the backup file to standard
output.

```sh
mkdir -p backups
ecdb database backup \
    --path games/example/db \
    --output-path backups \
    --version
# backups/ec.db.20260719T142530Z-1
```

### Compact a database

```text
ecdb database compact --path PATH
```

Reclaims unused space in an existing `PATH/ec.db`. `--path` is required. The
file must be an EC database at the schema version expected by this build. The
command produces no output on success.

```sh
ecdb database compact --path games/example/db
```

### Print the database schema version

```text
ecdb database version --path PATH
```

Opens `PATH/ec.db` read-only and prints its numeric SQLite `user_version`.
`--path` is required. The file must be identified as an EC database, but its
schema version does not need to match the version expected by this build.

```sh
ecdb database version --path games/example/db
# 1
```

## Print the application version

```text
ecdb version [--build | --long]
```

Prints the version of the `ecdb` application. With no flag it prints the core
semantic version. `--build` also includes pre-release information, while
`--long` includes both pre-release and build information. The two flags are
mutually exclusive.

```sh
ecdb version
ecdb version --build
ecdb version --long
```

The application version and database schema version are independent:

- `ecdb version` identifies the `ecdb` executable.
- `ecdb database version --path PATH` reads SQLite `user_version` from a
  database.

Do not use the application version to determine whether a database needs an
upgrade. Use `database verify`, `database version`, or `database upgrade`, as
appropriate.

## Paths, output, and failures

- Every database path is a directory containing a file named `ec.db`.
- Commands do not create missing directories.
- Only `database create` has a default input path (`db`). All other database
  commands require `--path`.
- `database backup` defaults its output directory to its input path.
- `database backup` writes the created file's path to standard output.
  Application and database `version` commands also write their value to
  standard output. `database upgrade` writes a status line unless `--quiet`
  is set. Other successful database operations are silent.
- Command failures return a nonzero exit status and normally write a
  diagnostic to standard error. For `database verify`, `--quiet` suppresses
  the diagnostic while preserving the nonzero status.
- `--quiet` does not suppress intentional version output.
