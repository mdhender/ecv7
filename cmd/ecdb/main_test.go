package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mdhender/ecv7"
	"github.com/mdhender/ecv7/internal/sqlite"
)

func TestCreateDatabase(t *testing.T) {
	dir := t.TempDir()
	if err := run(t.Context(), []string{"create", "database", "--path", dir}, &bytes.Buffer{}); err != nil {
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
	if err := run(t.Context(), []string{"create", "database"}, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join("db", sqlite.DatabaseName)); err != nil {
		t.Fatalf("stat database: %v", err)
	}
}

func TestCreateDatabaseDoesNotCreatePath(t *testing.T) {
	parent := t.TempDir()
	missing := filepath.Join(parent, "missing")
	err := run(t.Context(), []string{"create", "database", "--path", missing}, &bytes.Buffer{})
	if !errors.Is(err, sqlite.ErrInvalidDirectory) {
		t.Fatalf("run error = %v, want ErrInvalidDirectory", err)
	}
	if _, err := os.Stat(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing path was created: %v", err)
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
