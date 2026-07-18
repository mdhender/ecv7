// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mdhender/ecv7"
	"github.com/mdhender/ecv7/internal/dotenv"
	"github.com/mdhender/ecv7/internal/sqlite"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

const (
	envVarPrefix = "EC"
)

func command() *ff.Command {
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

	database.Subcommands = append(database.Subcommands, backup, compact, create, databaseVersion)

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

func run(ctx context.Context, args []string, stderr io.Writer) error {
	cmd := command()
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
		fmt.Fprintf(os.Stderr, "ecdb: %v\n", err)
		os.Exit(1)
	}
}
