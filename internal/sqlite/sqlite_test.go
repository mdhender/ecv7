package sqlite

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	zsqlite "zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func TestCreatePermanent(t *testing.T) {
	dir := t.TempDir()
	db, err := CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	assertSchema(t, db)
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, DatabaseName)); err != nil {
		t.Fatalf("stat database: %v", err)
	}
	if _, err := CreatePermanent(t.Context(), dir); !errors.Is(err, ErrDatabaseExists) {
		t.Fatalf("second CreatePermanent error = %v, want ErrDatabaseExists", err)
	}
}

func TestCreatePermanentRequiresExistingDirectory(t *testing.T) {
	parent := t.TempDir()
	missing := filepath.Join(parent, "missing")
	if _, err := CreatePermanent(t.Context(), missing); !errors.Is(err, ErrInvalidDirectory) {
		t.Fatalf("CreatePermanent error = %v, want ErrInvalidDirectory", err)
	}
	if _, err := os.Stat(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing directory was created: %v", err)
	}

	file := filepath.Join(parent, "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := CreatePermanent(t.Context(), file); !errors.Is(err, ErrInvalidDirectory) {
		t.Fatalf("CreatePermanent with file error = %v, want ErrInvalidDirectory", err)
	}
}

func TestCreateTemporary(t *testing.T) {
	db, err := CreateTemporary(t.Context())
	if err != nil {
		t.Fatalf("CreateTemporary: %v", err)
	}
	defer db.Close()
	assertSchema(t, db)
}

func TestOpenPermanent(t *testing.T) {
	dir := t.TempDir()
	created, err := CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := created.Close(); err != nil {
		t.Fatalf("close created database: %v", err)
	}

	opened, err := OpenPermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanent: %v", err)
	}
	defer opened.Close()
	assertSchema(t, opened)
}

func TestOpenPermanentDoesNotCreate(t *testing.T) {
	dir := t.TempDir()
	if _, err := OpenPermanent(t.Context(), dir); !errors.Is(err, ErrDatabaseNotFound) {
		t.Fatalf("OpenPermanent error = %v, want ErrDatabaseNotFound", err)
	}
	if _, err := os.Stat(filepath.Join(dir, DatabaseName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("database was created: %v", err)
	}
}

func TestOpenPermanentRejectsWrongApplication(t *testing.T) {
	dir := t.TempDir()
	createRawDatabase(t, dir, 0, 0)
	before := rawJournalMode(t, dir)
	if _, err := OpenPermanent(t.Context(), dir); !errors.Is(err, ErrInvalidDatabase) {
		t.Fatalf("OpenPermanent error = %v, want ErrInvalidDatabase", err)
	}
	if after := rawJournalMode(t, dir); after != before {
		t.Fatalf("journal mode changed from %q to %q", before, after)
	}
}

func TestOpenPermanentRejectsNewerSchema(t *testing.T) {
	dir := t.TempDir()
	createRawDatabase(t, dir, applicationID, ExpectedSchemaVersion+1)
	if _, err := OpenPermanent(t.Context(), dir); !errors.Is(err, ErrNewerSchemaVersion) {
		t.Fatalf("OpenPermanent error = %v, want ErrNewerSchemaVersion", err)
	}
}

func TestOpenPermanentMigratesOlderSchema(t *testing.T) {
	dir := t.TempDir()
	createRawDatabase(t, dir, applicationID, 0)
	db, err := OpenPermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanent: %v", err)
	}
	defer db.Close()
	assertSchema(t, db)
}

func assertSchema(t *testing.T, db *DB) {
	t.Helper()
	conn, err := db.Get(t.Context())
	if err != nil {
		t.Fatalf("get connection: %v", err)
	}
	defer db.Put(conn)

	if got, err := pragmaInt(conn, "application_id"); err != nil || int32(got) != applicationID {
		t.Fatalf("application_id = %#x, %v; want %#x", got, err, applicationID)
	}
	if got, err := pragmaInt(conn, "user_version"); err != nil || got != ExpectedSchemaVersion {
		t.Fatalf("user_version = %d, %v; want %d", got, err, ExpectedSchemaVersion)
	}
	var columns int
	err = sqlitex.ExecuteTransient(conn, "SELECT count(*) FROM pragma_table_info('metadata') WHERE name IN ('key', 'value') AND type = 'TEXT' AND \"notnull\" = 1;", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *zsqlite.Stmt) error {
			columns = stmt.ColumnInt(0)
			return nil
		},
	})
	if err != nil || columns != 2 {
		t.Fatalf("metadata columns = %d, %v; want 2", columns, err)
	}
}

func createRawDatabase(t *testing.T, dir string, appID int32, version int) {
	t.Helper()
	conn, err := zsqlite.OpenConn(filepath.Join(dir, DatabaseName), zsqlite.OpenReadWrite|zsqlite.OpenCreate)
	if err != nil {
		t.Fatalf("create raw database: %v", err)
	}
	defer conn.Close()
	for _, query := range []string{
		fmt.Sprintf("PRAGMA application_id = %d;", appID),
		fmt.Sprintf("PRAGMA user_version = %d;", version),
	} {
		if err := sqlitex.ExecuteTransient(conn, query, nil); err != nil {
			t.Fatalf("execute %q: %v", query, err)
		}
	}
}

func rawJournalMode(t *testing.T, dir string) string {
	t.Helper()
	conn, err := zsqlite.OpenConn(filepath.Join(dir, DatabaseName), zsqlite.OpenReadOnly)
	if err != nil {
		t.Fatalf("open raw database: %v", err)
	}
	defer conn.Close()

	var mode string
	err = sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode;", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *zsqlite.Stmt) error {
			mode = stmt.ColumnText(0)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("read journal mode: %v", err)
	}
	return mode
}
