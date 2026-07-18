package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

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
