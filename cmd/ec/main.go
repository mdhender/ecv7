// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdhender/ecv7"
	"github.com/mdhender/ecv7/internal/dotenv"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

const envVarPrefix = "EC"

func command() *ff.Command {
	rootFlags := ff.NewFlagSet("ec")
	rootFlags.BoolLong("quiet", "suppress status and diagnostic output")
	rootFlags.StringLong("path", "db", "directory containing ec.db")
	root := &ff.Command{
		Name:      "ec",
		Usage:     "ec <SUBCOMMAND>",
		ShortHelp: "play EC",
		Flags:     rootFlags,
	}

	versionFlags := ff.NewFlagSet("version").SetParent(rootFlags)
	build := versionFlags.BoolLong("build", "include pre-release version information")
	long := versionFlags.BoolLong("long", "include pre-release and build version information")
	version := &ff.Command{
		Name:      "version",
		Usage:     "ec version [--build | --long]",
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

	root.Subcommands = append(root.Subcommands, version)
	return root
}

func run(ctx context.Context, args []string, stderr io.Writer) error {
	cmd := command()
	if err := cmd.Parse(args, ff.WithEnvVarPrefix(envVarPrefix)); err != nil {
		fmt.Fprint(stderr, ffhelp.Command(cmd.GetSelected()))
		return err
	}
	selected := cmd.GetSelected()
	if selected.Exec == nil && len(selected.Subcommands) != 0 {
		fmt.Fprint(stderr, ffhelp.Command(selected))
		if args := selected.Flags.GetArgs(); len(args) != 0 {
			return fmt.Errorf("unknown command %q", args[0])
		}
		return errors.New("no command specified")
	}
	return cmd.Run(ctx)
}

func main() {
	env, ok := os.LookupEnv(envVarPrefix + "_ENV")
	if !ok {
		env = "development"
	}
	if err := dotenv.Load(env); err != nil {
		fmt.Fprintf(os.Stderr, "ec: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, os.Args[1:], os.Stderr); err != nil {
		if errors.Is(err, ff.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "ec: %v\n", err)
		os.Exit(1)
	}
}
