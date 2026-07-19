package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mdhender/ecv7"
	"github.com/mdhender/ecv7/internal/sqlite"
	"github.com/peterbourgon/ff/v4"
	zsqlite "zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestCreateDatabase(t *testing.T) {
	dir := t.TempDir()
	if err := run(t.Context(), discardLogger, []string{"database", "create", "--path", dir}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, sqlite.DatabaseName)); err != nil {
		t.Fatalf("stat database: %v", err)
	}
}

func TestCreateDatabaseDefaultPath(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("db", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := run(t.Context(), discardLogger, []string{"database", "create"}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join("db", sqlite.DatabaseName)); err != nil {
		t.Fatalf("stat database: %v", err)
	}
}

func TestCreateDatabaseDoesNotCreatePath(t *testing.T) {
	parent := t.TempDir()
	missing := filepath.Join(parent, "missing")
	err := run(t.Context(), discardLogger, []string{"database", "create", "--path", missing}, &bytes.Buffer{})
	if !errors.Is(err, sqlite.ErrInvalidDirectory) {
		t.Fatalf("run error = %v, want ErrInvalidDirectory", err)
	}
	if _, err := os.Stat(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing path was created: %v", err)
	}
}

func TestBackupDatabase(t *testing.T) {
	sourceDir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), sourceDir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	outputDir := t.TempDir()

	if err := run(t.Context(), discardLogger, []string{"database", "backup", "--path", sourceDir, "--output-path", outputDir, "--version"}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(outputDir, "ec.db.*Z-"+strconv.Itoa(sqlite.ExpectedSchemaVersion)))
	if err != nil {
		t.Fatalf("glob backup: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("backup matches = %v, want one versioned backup", matches)
	}
}

func TestBackupDatabaseDefaultsOutputPath(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err := run(t.Context(), discardLogger, []string{"database", "backup", "--path", dir}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "ec.db.*Z"))
	if err != nil {
		t.Fatalf("glob backup: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("backup matches = %v, want one backup", matches)
	}
}

func TestBackupDatabaseRequiresPath(t *testing.T) {
	err := run(t.Context(), discardLogger, []string{"database", "backup"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--path is required") {
		t.Fatalf("run error = %v, want required path error", err)
	}
}

func TestCompactDatabase(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err := run(t.Context(), discardLogger, []string{"database", "compact", "--path", dir}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestCompactDatabaseRequiresPath(t *testing.T) {
	err := run(t.Context(), discardLogger, []string{"database", "compact"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--path is required") {
		t.Fatalf("run error = %v, want required path error", err)
	}
}

func TestUpgradeDatabase(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	conn, err := db.Get(t.Context())
	if err != nil {
		t.Fatalf("get connection: %v", err)
	}
	for _, query := range []string{"DROP TABLE metadata;", "PRAGMA user_version = 0;"} {
		if err := sqlitex.ExecuteTransient(conn, query, nil); err != nil {
			t.Fatalf("reset schema: %v", err)
		}
	}
	db.Put(conn)
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	output, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"database", "upgrade", "--path", dir}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := fmt.Sprintf("migrations applied to %s (version %d)\n", dir, sqlite.ExpectedSchemaVersion)
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestUpgradeDatabaseCurrentAndQuiet(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	output, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"database", "upgrade", "--path", dir}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := fmt.Sprintf("no migrations applied to %s (version %d)\n", dir, sqlite.ExpectedSchemaVersion)
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}

	output, err = captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"--quiet", "database", "upgrade", "--path", dir}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("quiet run: %v", err)
	}
	if output != "" {
		t.Fatalf("quiet output = %q, want empty", output)
	}
}

func TestUpgradeDatabaseFailures(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T) []string
		wantErr error
	}{
		{
			name:    "path required",
			prepare: func(*testing.T) []string { return []string{"database", "upgrade"} },
			wantErr: errors.New("--path is required"),
		},
		{
			name: "missing directory",
			prepare: func(t *testing.T) []string {
				return []string{"database", "upgrade", "--path", filepath.Join(t.TempDir(), "missing")}
			},
			wantErr: sqlite.ErrInvalidDirectory,
		},
		{
			name: "missing database",
			prepare: func(t *testing.T) []string {
				return []string{"database", "upgrade", "--path", t.TempDir()}
			},
			wantErr: sqlite.ErrDatabaseNotFound,
		},
		{
			name: "wrong application",
			prepare: func(t *testing.T) []string {
				dir := t.TempDir()
				conn, err := sqlitex.Open(filepath.Join(dir, sqlite.DatabaseName), 0, 1)
				if err != nil {
					t.Fatalf("create raw database: %v", err)
				}
				if err := conn.Close(); err != nil {
					t.Fatalf("close raw database: %v", err)
				}
				return []string{"database", "upgrade", "--path", dir}
			},
			wantErr: sqlite.ErrInvalidDatabase,
		},
		{
			name: "newer schema",
			prepare: func(t *testing.T) []string {
				dir := t.TempDir()
				db, err := sqlite.CreatePermanent(t.Context(), dir)
				if err != nil {
					t.Fatalf("CreatePermanent: %v", err)
				}
				conn, err := db.Get(t.Context())
				if err != nil {
					t.Fatalf("get connection: %v", err)
				}
				query := fmt.Sprintf("PRAGMA user_version = %d;", sqlite.ExpectedSchemaVersion+1)
				if err := sqlitex.ExecuteTransient(conn, query, nil); err != nil {
					t.Fatalf("set schema version: %v", err)
				}
				db.Put(conn)
				if err := db.Close(); err != nil {
					t.Fatalf("close: %v", err)
				}
				return []string{"database", "upgrade", "--path", dir}
			},
			wantErr: sqlite.ErrNewerSchemaVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captureStdout(t, func() error {
				return run(t.Context(), discardLogger, tt.prepare(t), &bytes.Buffer{})
			})
			if tt.name == "path required" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("run error = %v, want %v", err, tt.wantErr)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("run error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunPropagatesLoggerToUpgrade(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := run(t.Context(), log, []string{"--quiet", "database", "upgrade", "--path", dir}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got := logs.String(); !strings.Contains(got, "msg=\"database upgrade: starting\"") || !strings.Contains(got, "path="+filepath.Join(dir, sqlite.DatabaseName)) {
		t.Fatalf("log = %q, want start message and database path", got)
	}
}

func TestDatabaseVersion(t *testing.T) {
	for _, want := range []int{0, sqlite.ExpectedSchemaVersion, sqlite.ExpectedSchemaVersion + 1} {
		t.Run(strconv.Itoa(want), func(t *testing.T) {
			dir := t.TempDir()
			db, err := sqlite.CreatePermanent(t.Context(), dir)
			if err != nil {
				t.Fatalf("CreatePermanent: %v", err)
			}
			conn, err := db.Get(t.Context())
			if err != nil {
				t.Fatalf("get connection: %v", err)
			}
			if err := sqlitex.ExecuteTransient(conn, fmt.Sprintf("PRAGMA user_version = %d;", want), nil); err != nil {
				t.Fatalf("set schema version: %v", err)
			}
			db.Put(conn)
			if err := db.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}

			output, err := captureStdout(t, func() error {
				return run(t.Context(), discardLogger, []string{"database", "version", "--path", dir}, &bytes.Buffer{})
			})
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			if got := strings.TrimSpace(output); got != strconv.Itoa(want) {
				t.Fatalf("output = %q, want %q", got, strconv.Itoa(want))
			}

			readOnly, err := sqlite.OpenPermanentReadOnly(t.Context(), dir)
			if err != nil {
				t.Fatalf("OpenPermanentReadOnly: %v", err)
			}
			got, err := readOnly.SchemaVersion(t.Context())
			if closeErr := readOnly.Close(); err == nil {
				err = closeErr
			}
			if err != nil || got != want {
				t.Fatalf("schema version after command = %d, %v; want %d", got, err, want)
			}
		})
	}
}

func TestDatabaseVersionRequiresPath(t *testing.T) {
	_, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"database", "version"}, &bytes.Buffer{})
	})
	if err == nil || !strings.Contains(err.Error(), "--path is required") {
		t.Fatalf("run error = %v, want required path error", err)
	}
}

func TestDatabaseVersionRequiresExistingDatabase(t *testing.T) {
	dir := t.TempDir()
	_, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"database", "version", "--path", dir}, &bytes.Buffer{})
	})
	if !errors.Is(err, sqlite.ErrDatabaseNotFound) {
		t.Fatalf("run error = %v, want ErrDatabaseNotFound", err)
	}
}

func TestDatabaseVersionRequiresExistingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")
	_, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"database", "version", "--path", dir}, &bytes.Buffer{})
	})
	if !errors.Is(err, sqlite.ErrInvalidDirectory) {
		t.Fatalf("run error = %v, want ErrInvalidDirectory", err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing directory was created: %v", err)
	}
}

func TestVerifyDatabase(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	var stderr bytes.Buffer
	if err := run(t.Context(), discardLogger, []string{"database", "verify", "--path", dir}, &stderr); err != nil {
		t.Fatalf("run: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want no output", stderr.String())
	}
}

func TestVerifyDatabaseFailures(t *testing.T) {
	tests := []struct {
		name    string
		args    func(*testing.T) []string
		wantErr error
	}{
		{
			name:    "path required",
			args:    func(*testing.T) []string { return []string{"database", "verify"} },
			wantErr: errors.New("--path is required"),
		},
		{
			name: "missing directory",
			args: func(t *testing.T) []string {
				return []string{"database", "verify", "--path", filepath.Join(t.TempDir(), "missing")}
			},
			wantErr: sqlite.ErrInvalidDirectory,
		},
		{
			name: "missing database",
			args: func(t *testing.T) []string {
				return []string{"database", "verify", "--path", t.TempDir()}
			},
			wantErr: sqlite.ErrDatabaseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args(t)
			var stderr bytes.Buffer
			err := run(t.Context(), discardLogger, args, &stderr)
			if tt.name == "path required" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("run error = %v, want %v", err, tt.wantErr)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("run error = %v, want %v", err, tt.wantErr)
			}
			if got := stderr.String(); !strings.Contains(got, tt.wantErr.Error()) {
				t.Fatalf("stderr = %q, want underlying error %q", got, tt.wantErr)
			}
			for i, arg := range args {
				if arg == "--path" && !strings.Contains(stderr.String(), args[i+1]) {
					t.Fatalf("stderr = %q, want database path %q", stderr.String(), args[i+1])
				}
			}
		})
	}
}

func TestVerifyDatabaseQuietFailure(t *testing.T) {
	tests := [][]string{
		{"--quiet", "database", "verify", "--path", ""},
		{"database", "verify", "--path", "", "--quiet"},
	}
	for _, args := range tests {
		dir := filepath.Join(t.TempDir(), "missing")
		for i, arg := range args {
			if arg == "--path" {
				args[i+1] = dir
			}
		}
		var stderr bytes.Buffer
		err := run(t.Context(), discardLogger, args, &stderr)
		if !errors.Is(err, sqlite.ErrInvalidDirectory) {
			t.Fatalf("run(%v) error = %v, want ErrInvalidDirectory", args, err)
		}
		if stderr.Len() != 0 {
			t.Fatalf("run(%v) stderr = %q, want no output", args, stderr.String())
		}
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "core", args: []string{"version"}, want: ecv7.Version().Core()},
		{name: "build", args: []string{"version", "--build"}, want: ecv7.Version().Short()},
		{name: "long", args: []string{"version", "--long"}, want: ecv7.Version().String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := captureStdout(t, func() error {
				return run(t.Context(), discardLogger, tt.args, &bytes.Buffer{})
			})
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			if got := strings.TrimSpace(output); got != tt.want {
				t.Fatalf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVersionFlagsAreMutuallyExclusive(t *testing.T) {
	_, err := captureStdout(t, func() error {
		return run(t.Context(), discardLogger, []string{"version", "--build", "--long"}, &bytes.Buffer{})
	})
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("run error = %v, want mutually exclusive error", err)
	}
}

func TestDatabaseHelpShowsUpgradeDirectly(t *testing.T) {
	var stderr bytes.Buffer
	err := run(t.Context(), discardLogger, []string{"database", "--help"}, &stderr)
	if !errors.Is(err, ff.ErrHelp) {
		t.Fatalf("run error = %v, want ErrHelp", err)
	}
	help := stderr.String()
	if !strings.Contains(help, "upgrade   apply missing database migrations") {
		t.Fatalf("help = %q, want upgrade subcommand", help)
	}
	if strings.Contains(help, "migrate   ") {
		t.Fatalf("help = %q, must not contain migrate subcommand", help)
	}
}

func TestQuietHelpAndRemovedVerboseFlag(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"database", "verify", "--help"}} {
		var stderr bytes.Buffer
		err := run(t.Context(), discardLogger, args, &stderr)
		if !errors.Is(err, ff.ErrHelp) {
			t.Fatalf("run(%v) error = %v, want ErrHelp", args, err)
		}
		help := stderr.String()
		if !strings.Contains(help, "--quiet") {
			t.Fatalf("run(%v) help = %q, want quiet flag", args, help)
		}
		if strings.Contains(help, "--verbose") {
			t.Fatalf("run(%v) help = %q, must not contain verbose flag", args, help)
		}
	}

	var stderr bytes.Buffer
	err := run(t.Context(), discardLogger, []string{"database", "verify", "--path", t.TempDir(), "--verbose"}, &stderr)
	if err == nil {
		t.Fatal("run with --verbose succeeded, want parse error")
	}
}

func TestCommandErrors(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  string
		wantHelp string
	}{
		{name: "missing root command", wantErr: "no command specified", wantHelp: "ecdb <SUBCOMMAND>"},
		{name: "unknown root command", args: []string{"databse", "verify"}, wantErr: `unknown command "databse"`, wantHelp: "ecdb <SUBCOMMAND>"},
		{name: "missing database command", args: []string{"database"}, wantErr: "database: no command specified", wantHelp: "ecdb database <SUBCOMMAND>"},
		{name: "unknown database command", args: []string{"database", "varify"}, wantErr: `database: unknown command "varify"`, wantHelp: "ecdb database <SUBCOMMAND>"},
		{name: "removed migrate command", args: []string{"database", "migrate", "up"}, wantErr: `database: unknown command "migrate"`, wantHelp: "ecdb database <SUBCOMMAND>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			err := run(t.Context(), discardLogger, tt.args, &stderr)
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("run error = %v, want %q", err, tt.wantErr)
			}
			if got := stderr.String(); !strings.Contains(got, tt.wantHelp) {
				t.Fatalf("stderr = %q, want help containing %q", got, tt.wantHelp)
			}
		})
	}
}

func TestRunPropagatesCanceledContext(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err = run(ctx, discardLogger, []string{"database", "verify", "--path", dir}, &bytes.Buffer{})
	if code := zsqlite.ErrCode(err); !errors.Is(err, context.Canceled) && code != zsqlite.ResultInterrupt {
		t.Fatalf("run error = %v (code %v), want context cancellation or SQLite interrupt", err, code)
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = original })

	runErr := fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = original
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(output), runErr
}
