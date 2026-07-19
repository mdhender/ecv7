package main

import (
	"bytes"
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
	"zombiezen.com/go/sqlite/sqlitex"
)

func TestCreateDatabase(t *testing.T) {
	dir := t.TempDir()
	if err := run(t.Context(), []string{"database", "create", "--path", dir}, &bytes.Buffer{}); err != nil {
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
	if err := run(t.Context(), []string{"database", "create"}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join("db", sqlite.DatabaseName)); err != nil {
		t.Fatalf("stat database: %v", err)
	}
}

func TestCreateDatabaseDoesNotCreatePath(t *testing.T) {
	parent := t.TempDir()
	missing := filepath.Join(parent, "missing")
	err := run(t.Context(), []string{"database", "create", "--path", missing}, &bytes.Buffer{})
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

	if err := run(t.Context(), []string{"database", "backup", "--path", sourceDir, "--output-path", outputDir, "--version"}, &bytes.Buffer{}); err != nil {
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

	if err := run(t.Context(), []string{"database", "backup", "--path", dir}, &bytes.Buffer{}); err != nil {
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
	err := run(t.Context(), []string{"database", "backup"}, &bytes.Buffer{})
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

	if err := run(t.Context(), []string{"database", "compact", "--path", dir}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestCompactDatabaseRequiresPath(t *testing.T) {
	err := run(t.Context(), []string{"database", "compact"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--path is required") {
		t.Fatalf("run error = %v, want required path error", err)
	}
}

func TestMigrateUpDatabase(t *testing.T) {
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
		return run(t.Context(), []string{"database", "migrate", "up", "--path", dir}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := fmt.Sprintf("migrations applied to %s (version %d)\n", dir, sqlite.ExpectedSchemaVersion)
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestMigrateUpDatabaseCurrentAndQuiet(t *testing.T) {
	dir := t.TempDir()
	db, err := sqlite.CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	output, err := captureStdout(t, func() error {
		return run(t.Context(), []string{"database", "migrate", "up", "--path", dir}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := fmt.Sprintf("no migrations applied to %s (version %d)\n", dir, sqlite.ExpectedSchemaVersion)
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}

	output, err = captureStdout(t, func() error {
		return run(t.Context(), []string{"database", "migrate", "up", "--path", dir, "--quiet"}, &bytes.Buffer{})
	})
	if err != nil {
		t.Fatalf("quiet run: %v", err)
	}
	if output != "" {
		t.Fatalf("quiet output = %q, want empty", output)
	}
}

func TestMigrateUpDatabaseFailures(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T) []string
		wantErr error
	}{
		{
			name:    "path required",
			prepare: func(*testing.T) []string { return []string{"database", "migrate", "up"} },
			wantErr: errors.New("--path is required"),
		},
		{
			name: "missing directory",
			prepare: func(t *testing.T) []string {
				return []string{"database", "migrate", "up", "--path", filepath.Join(t.TempDir(), "missing")}
			},
			wantErr: sqlite.ErrInvalidDirectory,
		},
		{
			name: "missing database",
			prepare: func(t *testing.T) []string {
				return []string{"database", "migrate", "up", "--path", t.TempDir()}
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
				return []string{"database", "migrate", "up", "--path", dir}
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
				return []string{"database", "migrate", "up", "--path", dir}
			},
			wantErr: sqlite.ErrNewerSchemaVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captureStdout(t, func() error {
				return run(t.Context(), tt.prepare(t), &bytes.Buffer{})
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

func TestMigrateUpLogsDatabasePath(t *testing.T) {
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
	if err := migrateUp(t.Context(), log, dir, true); err != nil {
		t.Fatalf("migrateUp: %v", err)
	}
	if got := logs.String(); !strings.Contains(got, "msg=\"migration up: starting\"") || !strings.Contains(got, "path="+filepath.Join(dir, sqlite.DatabaseName)) {
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
				return run(t.Context(), []string{"database", "version", "--path", dir}, &bytes.Buffer{})
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
		return run(t.Context(), []string{"database", "version"}, &bytes.Buffer{})
	})
	if err == nil || !strings.Contains(err.Error(), "--path is required") {
		t.Fatalf("run error = %v, want required path error", err)
	}
}

func TestDatabaseVersionRequiresExistingDatabase(t *testing.T) {
	dir := t.TempDir()
	_, err := captureStdout(t, func() error {
		return run(t.Context(), []string{"database", "version", "--path", dir}, &bytes.Buffer{})
	})
	if !errors.Is(err, sqlite.ErrDatabaseNotFound) {
		t.Fatalf("run error = %v, want ErrDatabaseNotFound", err)
	}
}

func TestDatabaseVersionRequiresExistingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")
	_, err := captureStdout(t, func() error {
		return run(t.Context(), []string{"database", "version", "--path", dir}, &bytes.Buffer{})
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
	if err := run(t.Context(), []string{"database", "verify", "--path", dir}, &stderr); err != nil {
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
			var stderr bytes.Buffer
			err := run(t.Context(), tt.args(t), &stderr)
			if tt.name == "path required" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("run error = %v, want %v", err, tt.wantErr)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("run error = %v, want %v", err, tt.wantErr)
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q, want no output", stderr.String())
			}
		})
	}
}

func TestVerifyDatabaseVerboseFailure(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")
	var stderr bytes.Buffer
	err := run(t.Context(), []string{"database", "verify", "--path", dir, "--verbose"}, &stderr)
	if !errors.Is(err, sqlite.ErrInvalidDirectory) {
		t.Fatalf("run error = %v, want ErrInvalidDirectory", err)
	}
	if got := stderr.String(); !strings.Contains(got, sqlite.ErrInvalidDirectory.Error()) || !strings.Contains(got, dir) {
		t.Fatalf("stderr = %q, want path and invalid-directory error", got)
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
				return run(t.Context(), tt.args, &bytes.Buffer{})
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
		return run(t.Context(), []string{"version", "--build", "--long"}, &bytes.Buffer{})
	})
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("run error = %v, want mutually exclusive error", err)
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
