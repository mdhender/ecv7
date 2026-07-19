// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mdhender/ecv7"
	"github.com/mdhender/ecv7/internal/dotenv"
	"github.com/mdhender/ecv7/internal/sqlite"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

const (
	envVarPrefix = "EC"
)

type quietError struct {
	err error
}

func (e *quietError) Error() string { return e.err.Error() }
func (e *quietError) Unwrap() error { return e.err }

func command(stderr io.Writer) *ff.Command {
	rootFlags := ff.NewFlagSet("ecdb")
	root := &ff.Command{
		Name:      "ecdb",
		Usage:     "ecdb <SUBCOMMAND>",
		ShortHelp: "manage EC databases",
		Flags:     rootFlags,
	}

	databaseFlags := ff.NewFlagSet("database").SetParent(rootFlags)
	database := &ff.Command{
		Name:      "database",
		Usage:     "ecdb database <SUBCOMMAND>",
		ShortHelp: "manage the persistent database",
		Flags:     databaseFlags,
	}

	createFlags := ff.NewFlagSet("create").SetParent(databaseFlags)
	path := createFlags.StringLong("path", "db", "directory in which to create ec.db")
	create := &ff.Command{
		Name:      "create",
		Usage:     "ecdb database create [--path PATH]",
		ShortHelp: "create and migrate a persistent database",
		Flags:     createFlags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			db, err := sqlite.CreatePermanent(ctx, *path)
			if err != nil {
				return err
			}
			return db.Close()
		},
	}

	backupFlags := ff.NewFlagSet("backup").SetParent(databaseFlags)
	backupPath := backupFlags.StringLong("path", "", "directory containing ec.db")
	backupOutputPath := backupFlags.StringLong("output-path", "", "directory in which to write the backup")
	backupVersion := backupFlags.BoolLong("version", "append the database schema version to the backup name")
	backup := &ff.Command{
		Name:      "backup",
		Usage:     "ecdb database backup --path PATH [--output-path OUTPUT_PATH] [--version]",
		ShortHelp: "write a consistent timestamped database backup",
		Flags:     backupFlags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			if *backupPath == "" {
				return errors.New("--path is required")
			}
			outputPath := *backupOutputPath
			if outputPath == "" {
				outputPath = *backupPath
			}
			_, err := sqlite.BackupPermanent(ctx, *backupPath, outputPath, *backupVersion)
			return err
		},
	}

	compactFlags := ff.NewFlagSet("compact").SetParent(databaseFlags)
	compactPath := compactFlags.StringLong("path", "", "directory containing ec.db")
	compact := &ff.Command{
		Name:      "compact",
		Usage:     "ecdb database compact --path PATH",
		ShortHelp: "reclaim unused database space",
		Flags:     compactFlags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			if *compactPath == "" {
				return errors.New("--path is required")
			}
			return sqlite.CompactPermanent(ctx, *compactPath)
		},
	}

	migrateFlags := ff.NewFlagSet("migrate").SetParent(databaseFlags)
	migrate := &ff.Command{
		Name:      "migrate",
		Usage:     "ecdb database migrate <SUBCOMMAND>",
		ShortHelp: "migrate the persistent database",
		Flags:     migrateFlags,
	}

	migrateUpFlags := ff.NewFlagSet("up").SetParent(migrateFlags)
	migrateUpPath := migrateUpFlags.StringLong("path", "", "directory containing ec.db")
	migrateUpQuiet := migrateUpFlags.BoolLong("quiet", "do not write migration status")
	migrateUpCommand := &ff.Command{
		Name:      "up",
		Usage:     "ecdb database migrate up --path PATH [--quiet]",
		ShortHelp: "apply missing database migrations",
		Flags:     migrateUpFlags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			if *migrateUpPath == "" {
				return errors.New("--path is required")
			}
			return migrateUp(ctx, slog.Default(), *migrateUpPath, *migrateUpQuiet)
		},
	}
	migrate.Subcommands = append(migrate.Subcommands, migrateUpCommand)

	databaseVersionFlags := ff.NewFlagSet("version").SetParent(databaseFlags)
	databaseVersionPath := databaseVersionFlags.StringLong("path", "", "directory containing ec.db")
	databaseVersion := &ff.Command{
		Name:      "version",
		Usage:     "ecdb database version --path PATH",
		ShortHelp: "print the database migration version",
		Flags:     databaseVersionFlags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			if *databaseVersionPath == "" {
				return errors.New("--path is required")
			}
			db, err := sqlite.OpenPermanentReadOnly(ctx, *databaseVersionPath)
			if err != nil {
				return err
			}
			migrationVersion, err := db.SchemaVersion(ctx)
			if closeErr := db.Close(); err == nil {
				err = closeErr
			}
			if err != nil {
				return err
			}
			fmt.Println(migrationVersion)
			return nil
		},
	}

	verifyFlags := ff.NewFlagSet("verify").SetParent(databaseFlags)
	verifyPath := verifyFlags.StringLong("path", "", "directory containing ec.db")
	verbose := verifyFlags.BoolLong("verbose", "write verification errors to standard error")
	verify := &ff.Command{
		Name:      "verify",
		Usage:     "ecdb database verify --path PATH [--verbose]",
		ShortHelp: "verify the database type and migration version",
		Flags:     verifyFlags,
		Exec: func(ctx context.Context, args []string) error {
			var err error
			switch {
			case len(args) != 0:
				err = fmt.Errorf("unexpected arguments: %v", args)
			case *verifyPath == "":
				err = errors.New("--path is required")
			default:
				err = sqlite.VerifyPermanent(ctx, *verifyPath)
			}
			if err == nil {
				return nil
			}
			if *verbose {
				fmt.Fprintf(stderr, "ecdb: %v\n", err)
			}
			return &quietError{err: err}
		},
	}

	database.Subcommands = append(database.Subcommands, backup, compact, create, migrate, databaseVersion, verify)

	versionFlags := ff.NewFlagSet("version").SetParent(rootFlags)
	build := versionFlags.BoolLong("build", "include pre-release version information")
	long := versionFlags.BoolLong("long", "include pre-release and build version information")
	version := &ff.Command{
		Name:      "version",
		Usage:     "ecdb version [--build | --long]",
		ShortHelp: "print version information",
		Flags:     versionFlags,
		Exec: func(_ context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			if *build && *long {
				return errors.New("--build and --long are mutually exclusive")
			}
			switch {
			case *build:
				fmt.Println(ecv7.Version().Short())
			case *long:
				fmt.Println(ecv7.Version().String())
			default:
				fmt.Println(ecv7.Version().Core())
			}
			return nil
		},
	}

	root.Subcommands = append(root.Subcommands, database, version)
	return root
}

func migrateUp(ctx context.Context, log *slog.Logger, path string, quiet bool) error {
	dbPath := filepath.Join(path, sqlite.DatabaseName)
	log.Debug("migration up: starting", "path", dbPath)
	applied, err := sqlite.MigratePermanent(ctx, path)
	if err != nil {
		return err
	}
	if quiet {
		return nil
	}
	if !applied {
		fmt.Printf("no migrations applied to %s (version %d)\n", path, sqlite.ExpectedSchemaVersion)
		return nil
	}
	fmt.Printf("migrations applied to %s (version %d)\n", path, sqlite.ExpectedSchemaVersion)
	return nil
}

func run(ctx context.Context, args []string, stderr io.Writer) error {
	cmd := command(stderr)
	if err := cmd.Parse(args, ff.WithEnvVarPrefix(envVarPrefix)); err != nil {
		fmt.Fprint(stderr, ffhelp.Command(cmd))
		return err
	}
	if err := cmd.Run(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	env, ok := os.LookupEnv(envVarPrefix + "_ENV")
	if !ok {
		env = "development"
	}
	if err := dotenv.Load(env); err != nil {
		fmt.Fprintf(os.Stderr, "ecdb: %v\n", err)
		os.Exit(1)
	}

	if err := run(context.Background(), os.Args[1:], os.Stderr); err != nil {
		if errors.Is(err, ff.ErrHelp) {
			return
		}
		var quiet *quietError
		if !errors.As(err, &quiet) {
			fmt.Fprintf(os.Stderr, "ecdb: %v\n", err)
		}
		os.Exit(1)
	}
}
