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

	createFlags := ff.NewFlagSet("create").SetParent(rootFlags)
	create := &ff.Command{
		Name:      "create",
		Usage:     "ecdb create <SUBCOMMAND>",
		ShortHelp: "create database resources",
		Flags:     createFlags,
	}

	databaseFlags := ff.NewFlagSet("database").SetParent(createFlags)
	path := databaseFlags.StringLong("path", "db", "directory in which to create ec.db")
	database := &ff.Command{
		Name:      "database",
		Usage:     "ecdb create database [--path PATH]",
		ShortHelp: "create and migrate a persistent database",
		Flags:     databaseFlags,
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

	create.Subcommands = append(create.Subcommands, database)

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

	root.Subcommands = append(root.Subcommands, create, version)
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
