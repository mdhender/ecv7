package sqlite

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestBackupPermanent(t *testing.T) {
	sourceDir := t.TempDir()
	db, err := CreatePermanent(t.Context(), sourceDir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	conn, err := db.Get(t.Context())
	if err != nil {
		t.Fatalf("get connection: %v", err)
	}
	if err := sqlitex.ExecuteTransient(conn, "INSERT INTO metadata (key, value) VALUES ('test', 'backup');", nil); err != nil {
		t.Fatalf("insert metadata: %v", err)
	}
	db.Put(conn)
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	outputDir := t.TempDir()
	stamp := time.Date(2026, 7, 8, 18, 32, 45, 0, time.FixedZone("test", -7*60*60))
	path, err := backupPermanent(t.Context(), sourceDir, outputDir, false, stamp)
	if err != nil {
		t.Fatalf("backupPermanent: %v", err)
	}
	wantPath := filepath.Join(outputDir, "ec.db.20260709T013245Z")
	if path != wantPath {
		t.Fatalf("backup path = %q, want %q", path, wantPath)
	}
	assertBackupMetadata(t, path)
}

func TestBackupPermanentIncludesVersion(t *testing.T) {
	sourceDir := t.TempDir()
	db, err := CreatePermanent(t.Context(), sourceDir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	stamp := time.Date(2026, 7, 8, 18, 32, 45, 0, time.UTC)
	path, err := backupPermanent(t.Context(), sourceDir, sourceDir, true, stamp)
	if err != nil {
		t.Fatalf("backupPermanent: %v", err)
	}
	wantPath := filepath.Join(sourceDir, fmt.Sprintf("ec.db.20260708T183245Z-%d", ExpectedSchemaVersion))
	if path != wantPath {
		t.Fatalf("backup path = %q, want %q", path, wantPath)
	}
}

func TestBackupPermanentRejectsExistingOutput(t *testing.T) {
	sourceDir := t.TempDir()
	db, err := CreatePermanent(t.Context(), sourceDir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	stamp := time.Date(2026, 7, 8, 18, 32, 45, 0, time.UTC)
	path := filepath.Join(sourceDir, "ec.db.20260708T183245Z")
	const original = "do not overwrite"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatalf("create existing output: %v", err)
	}
	if _, err := backupPermanent(t.Context(), sourceDir, sourceDir, false, stamp); err == nil {
		t.Fatal("backupPermanent succeeded with existing output")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read existing output: %v", err)
	}
	if string(got) != original {
		t.Fatalf("existing output = %q, want %q", got, original)
	}
}

func TestBackupPermanentValidatesSourceAndOutput(t *testing.T) {
	stamp := time.Date(2026, 7, 8, 18, 32, 45, 0, time.UTC)

	t.Run("missing source", func(t *testing.T) {
		if _, err := backupPermanent(t.Context(), t.TempDir(), t.TempDir(), false, stamp); !errors.Is(err, ErrDatabaseNotFound) {
			t.Fatalf("backupPermanent error = %v, want ErrDatabaseNotFound", err)
		}
	})

	t.Run("invalid database", func(t *testing.T) {
		sourceDir := t.TempDir()
		createRawDatabase(t, sourceDir, 0, ExpectedSchemaVersion)
		if _, err := backupPermanent(t.Context(), sourceDir, t.TempDir(), false, stamp); !errors.Is(err, ErrInvalidDatabase) {
			t.Fatalf("backupPermanent error = %v, want ErrInvalidDatabase", err)
		}
	})

	for _, version := range []int{0, ExpectedSchemaVersion + 1} {
		t.Run(fmt.Sprintf("schema version %d", version), func(t *testing.T) {
			sourceDir := t.TempDir()
			createRawDatabase(t, sourceDir, applicationID, version)
			if _, err := backupPermanent(t.Context(), sourceDir, t.TempDir(), false, stamp); !errors.Is(err, ErrUnexpectedSchemaVersion) {
				t.Fatalf("backupPermanent error = %v, want ErrUnexpectedSchemaVersion", err)
			}
		})
	}

	t.Run("missing output directory", func(t *testing.T) {
		sourceDir := t.TempDir()
		db, err := CreatePermanent(t.Context(), sourceDir)
		if err != nil {
			t.Fatalf("CreatePermanent: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
		outputDir := filepath.Join(t.TempDir(), "missing")
		if _, err := backupPermanent(t.Context(), sourceDir, outputDir, false, stamp); !errors.Is(err, ErrInvalidDirectory) {
			t.Fatalf("backupPermanent error = %v, want ErrInvalidDirectory", err)
		}
		if _, err := os.Stat(outputDir); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("output directory was created: %v", err)
		}
	})

	t.Run("output path is a file", func(t *testing.T) {
		sourceDir := t.TempDir()
		db, err := CreatePermanent(t.Context(), sourceDir)
		if err != nil {
			t.Fatalf("CreatePermanent: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
		outputPath := filepath.Join(t.TempDir(), "file")
		if err := os.WriteFile(outputPath, nil, 0o600); err != nil {
			t.Fatalf("create output file: %v", err)
		}
		if _, err := backupPermanent(t.Context(), sourceDir, outputPath, false, stamp); !errors.Is(err, ErrInvalidDirectory) {
			t.Fatalf("backupPermanent error = %v, want ErrInvalidDirectory", err)
		}
	})
}

func TestCompactPermanent(t *testing.T) {
	dir := t.TempDir()
	db, err := CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	conn, err := db.Get(t.Context())
	if err != nil {
		t.Fatalf("get connection: %v", err)
	}
	if err := sqlitex.ExecuteTransient(conn, "CREATE TABLE compact_test (value BLOB);", nil); err != nil {
		t.Fatalf("create compact test table: %v", err)
	}
	for range 256 {
		if err := sqlitex.ExecuteTransient(conn, "INSERT INTO compact_test (value) VALUES (zeroblob(4096));", nil); err != nil {
			t.Fatalf("insert compact test row: %v", err)
		}
	}
	if err := sqlitex.ExecuteTransient(conn, "DELETE FROM compact_test;", nil); err != nil {
		t.Fatalf("delete compact test rows: %v", err)
	}
	db.Put(conn)
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if before := rawPragmaInt(t, dir, "freelist_count"); before == 0 {
		t.Fatal("test database has no free pages before compaction")
	}

	if err := CompactPermanent(t.Context(), dir); err != nil {
		t.Fatalf("CompactPermanent: %v", err)
	}
	if after := rawPragmaInt(t, dir, "freelist_count"); after != 0 {
		t.Fatalf("freelist_count after compaction = %d, want 0", after)
	}

	opened, err := OpenPermanentReadOnly(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanentReadOnly: %v", err)
	}
	defer opened.Close()
	assertSchema(t, opened)
}

func TestCompactPermanentValidatesDatabase(t *testing.T) {
	t.Run("missing database", func(t *testing.T) {
		if err := CompactPermanent(t.Context(), t.TempDir()); !errors.Is(err, ErrDatabaseNotFound) {
			t.Fatalf("CompactPermanent error = %v, want ErrDatabaseNotFound", err)
		}
	})

	t.Run("invalid database", func(t *testing.T) {
		dir := t.TempDir()
		createRawDatabase(t, dir, 0, ExpectedSchemaVersion)
		if err := CompactPermanent(t.Context(), dir); !errors.Is(err, ErrInvalidDatabase) {
			t.Fatalf("CompactPermanent error = %v, want ErrInvalidDatabase", err)
		}
	})

	for _, version := range []int{0, ExpectedSchemaVersion + 1} {
		t.Run(fmt.Sprintf("schema version %d", version), func(t *testing.T) {
			dir := t.TempDir()
			createRawDatabase(t, dir, applicationID, version)
			if err := CompactPermanent(t.Context(), dir); !errors.Is(err, ErrUnexpectedSchemaVersion) {
				t.Fatalf("CompactPermanent error = %v, want ErrUnexpectedSchemaVersion", err)
			}
			if got := rawSchemaVersion(t, dir); got != version {
				t.Fatalf("schema version after rejection = %d, want %d", got, version)
			}
		})
	}
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

func TestOpenPermanentReadOnlyDoesNotMigrate(t *testing.T) {
	dir := t.TempDir()
	createRawDatabase(t, dir, applicationID, 0)
	beforeJournalMode := rawJournalMode(t, dir)

	db, err := OpenPermanentReadOnly(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanentReadOnly: %v", err)
	}
	if got, err := db.SchemaVersion(t.Context()); err != nil || got != 0 {
		t.Fatalf("SchemaVersion = %d, %v; want 0", got, err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if got := rawSchemaVersion(t, dir); got != 0 {
		t.Fatalf("database was migrated to version %d", got)
	}
	if got := rawJournalMode(t, dir); got != beforeJournalMode {
		t.Fatalf("journal mode changed from %q to %q", beforeJournalMode, got)
	}
}

func TestOpenPermanentReadOnlyAcceptsNewerSchema(t *testing.T) {
	dir := t.TempDir()
	want := ExpectedSchemaVersion + 1
	createRawDatabase(t, dir, applicationID, want)
	db, err := OpenPermanentReadOnly(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanentReadOnly: %v", err)
	}
	defer db.Close()
	if got, err := db.SchemaVersion(t.Context()); err != nil || got != want {
		t.Fatalf("SchemaVersion = %d, %v; want %d", got, err, want)
	}
}

func TestOpenPermanentReadOnlyReadsActiveWAL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DatabaseName)
	writer, err := zsqlite.OpenConn(path, zsqlite.OpenReadWrite|zsqlite.OpenCreate)
	if err != nil {
		t.Fatalf("create WAL database: %v", err)
	}
	defer writer.Close()
	want := ExpectedSchemaVersion + 1
	for _, query := range []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA wal_autocheckpoint = 0;",
		fmt.Sprintf("PRAGMA application_id = %d;", applicationID),
		fmt.Sprintf("PRAGMA user_version = %d;", want),
	} {
		if err := sqlitex.ExecuteTransient(writer, query, nil); err != nil {
			t.Fatalf("execute %q: %v", query, err)
		}
	}
	if _, err := os.Stat(path + "-wal"); err != nil {
		t.Fatalf("stat active WAL: %v", err)
	}

	db, err := OpenPermanentReadOnly(t.Context(), dir)
	if err != nil {
		t.Fatalf("OpenPermanentReadOnly: %v", err)
	}
	defer db.Close()
	if got, err := db.SchemaVersion(t.Context()); err != nil || got != want {
		t.Fatalf("SchemaVersion = %d, %v; want %d", got, err, want)
	}
}

func TestOpenPermanentReadOnlyRejectsWrongApplication(t *testing.T) {
	dir := t.TempDir()
	createRawDatabase(t, dir, 0, ExpectedSchemaVersion)
	if _, err := OpenPermanentReadOnly(t.Context(), dir); !errors.Is(err, ErrInvalidDatabase) {
		t.Fatalf("OpenPermanentReadOnly error = %v, want ErrInvalidDatabase", err)
	}
}

func TestVerifyPermanent(t *testing.T) {
	dir := t.TempDir()
	db, err := CreatePermanent(t.Context(), dir)
	if err != nil {
		t.Fatalf("CreatePermanent: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err := VerifyPermanent(t.Context(), dir); err != nil {
		t.Fatalf("VerifyPermanent: %v", err)
	}
}

func TestVerifyPermanentRejectsInvalidDatabase(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T) string
		wantErr error
	}{
		{
			name: "missing directory",
			prepare: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "missing")
			},
			wantErr: ErrInvalidDirectory,
		},
		{
			name:    "missing database",
			prepare: func(t *testing.T) string { return t.TempDir() },
			wantErr: ErrDatabaseNotFound,
		},
		{
			name: "wrong application",
			prepare: func(t *testing.T) string {
				dir := t.TempDir()
				createRawDatabase(t, dir, 0, ExpectedSchemaVersion)
				return dir
			},
			wantErr: ErrInvalidDatabase,
		},
		{
			name: "older schema",
			prepare: func(t *testing.T) string {
				dir := t.TempDir()
				createRawDatabase(t, dir, applicationID, ExpectedSchemaVersion-1)
				return dir
			},
			wantErr: ErrUnexpectedSchemaVersion,
		},
		{
			name: "newer schema",
			prepare: func(t *testing.T) string {
				dir := t.TempDir()
				createRawDatabase(t, dir, applicationID, ExpectedSchemaVersion+1)
				return dir
			},
			wantErr: ErrUnexpectedSchemaVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.prepare(t)
			if err := VerifyPermanent(t.Context(), dir); !errors.Is(err, tt.wantErr) {
				t.Fatalf("VerifyPermanent error = %v, want %v", err, tt.wantErr)
			}
		})
	}
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

func assertBackupMetadata(t *testing.T, path string) {
	t.Helper()
	conn, err := zsqlite.OpenConn(path, zsqlite.OpenReadOnly)
	if err != nil {
		t.Fatalf("open backup: %v", err)
	}
	defer conn.Close()
	var value string
	if err := sqlitex.ExecuteTransient(conn, "SELECT value FROM metadata WHERE key = 'test';", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *zsqlite.Stmt) error {
			value = stmt.ColumnText(0)
			return nil
		},
	}); err != nil {
		t.Fatalf("read backup metadata: %v", err)
	}
	if value != "backup" {
		t.Fatalf("backup metadata = %q, want %q", value, "backup")
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

func rawSchemaVersion(t *testing.T, dir string) int {
	return rawPragmaInt(t, dir, "user_version")
}

func rawPragmaInt(t *testing.T, dir, name string) int {
	t.Helper()
	conn, err := zsqlite.OpenConn(filepath.Join(dir, DatabaseName), zsqlite.OpenReadOnly)
	if err != nil {
		t.Fatalf("open raw database: %v", err)
	}
	defer conn.Close()
	value, err := pragmaInt(conn, name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return value
}
