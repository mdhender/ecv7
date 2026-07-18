// Copyright (c) 2026 Michael D Henderson. All rights reserved.

// Package sqlite manages EC's SQLite databases.
package sqlite

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mdhender/ecv7/internal/cerrs"
	zsqlite "zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
	"zombiezen.com/go/sqlite/sqlitex"
)

const (
	// DatabaseName is the name of an EC database within its directory.
	DatabaseName = "ec.db"

	// ExpectedSchemaVersion is the schema version supported by this build.
	ExpectedSchemaVersion = 1

	// applicationID identifies SQLite files created by this package. This value
	// is part of the on-disk format and must never change.
	applicationID int32 = 0x0EC7DB
)

const (
	ErrDatabaseExists     = cerrs.Error("database already exists")
	ErrDatabaseNotFound   = cerrs.Error("database not found")
	ErrInvalidDirectory   = cerrs.Error("invalid database directory")
	ErrInvalidDatabase    = cerrs.Error("not an EC database")
	ErrNewerSchemaVersion = cerrs.Error("database schema is newer than this binary")
)

var schema = sqlitemigration.Schema{
	AppID: applicationID,
	Migrations: []string{
		`CREATE TABLE metadata (
	key   TEXT NOT NULL,
	value TEXT NOT NULL
);`,
	},
}

// DB is a pool of connections to an initialized EC database.
type DB struct {
	pool *sqlitex.Pool
}

// CreatePermanent creates and migrates ec.db in dir. The directory must
// already exist, and the database must not.
func CreatePermanent(ctx context.Context, dir string) (*DB, error) {
	if err := validateDirectory(dir); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, DatabaseName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("%s: %w", path, ErrDatabaseExists)
		}
		return nil, fmt.Errorf("create %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return nil, fmt.Errorf("create %s: %w", path, err)
	}

	db, err := open(ctx, path, zsqlite.OpenReadWrite, true, false)
	if err != nil {
		_ = os.Remove(path)
		_ = os.Remove(path + "-shm")
		_ = os.Remove(path + "-wal")
		return nil, err
	}
	return db, nil
}

// CreateTemporary creates and migrates an isolated in-memory database.
func CreateTemporary(ctx context.Context) (*DB, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, fmt.Errorf("name temporary database: %w", err)
	}
	uri := "file:ecv7-" + hex.EncodeToString(id[:]) + "?mode=memory&cache=shared"
	return open(ctx, uri, zsqlite.OpenReadWrite|zsqlite.OpenURI|zsqlite.OpenSharedCache, false, false)
}

// OpenPermanent opens and migrates the existing ec.db in dir. It never
// creates a database. Databases created by another application and databases
// newer than this build are rejected.
func OpenPermanent(ctx context.Context, dir string) (*DB, error) {
	if err := validateDirectory(dir); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, DatabaseName)
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%s: %w", path, ErrDatabaseNotFound)
		}
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s: %w", path, ErrDatabaseNotFound)
	}

	return open(ctx, path, zsqlite.OpenReadWrite, true, true)
}

func open(ctx context.Context, uri string, flags zsqlite.OpenFlags, enableWAL, validateExisting bool) (*DB, error) {
	pool, err := sqlitex.NewPool(uri, sqlitex.PoolOptions{
		Flags:       flags,
		PrepareConn: enableForeignKeys,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	fail := func(err error) (*DB, error) {
		_ = pool.Close()
		return nil, err
	}

	conn, err := pool.Take(ctx)
	if err != nil {
		return fail(fmt.Errorf("open database: %w", err))
	}
	failWithConn := func(err error) (*DB, error) {
		pool.Put(conn)
		return fail(err)
	}

	if validateExisting {
		got, err := pragmaInt(conn, "application_id")
		if err != nil {
			return failWithConn(err)
		}
		if int32(got) != applicationID {
			return failWithConn(fmt.Errorf("%w: application ID is %#x", ErrInvalidDatabase, got))
		}
	}

	version, err := pragmaInt(conn, "user_version")
	if err != nil {
		return failWithConn(err)
	}
	if version > ExpectedSchemaVersion {
		return failWithConn(fmt.Errorf("%w: database is version %d, binary expects %d", ErrNewerSchemaVersion, version, ExpectedSchemaVersion))
	}

	if enableWAL {
		if err := sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode = WAL;", nil); err != nil {
			return failWithConn(fmt.Errorf("enable WAL: %w", err))
		}
	}
	if err := sqlitemigration.Migrate(ctx, conn, schema); err != nil {
		return failWithConn(fmt.Errorf("migrate database: %w", err))
	}
	version, err = pragmaInt(conn, "user_version")
	if err != nil {
		return failWithConn(err)
	}
	if version > ExpectedSchemaVersion {
		return failWithConn(fmt.Errorf("%w: database is version %d, binary expects %d", ErrNewerSchemaVersion, version, ExpectedSchemaVersion))
	}
	pool.Put(conn)
	return &DB{pool: pool}, nil
}

func validateDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", dir, ErrInvalidDirectory, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s: %w", dir, ErrInvalidDirectory)
	}
	return nil
}

func enableForeignKeys(conn *zsqlite.Conn) error {
	return sqlitex.ExecuteTransient(conn, "PRAGMA foreign_keys = ON;", nil)
}

func pragmaInt(conn *zsqlite.Conn, name string) (int, error) {
	var value int
	err := sqlitex.ExecuteTransient(conn, "PRAGMA "+name+";", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *zsqlite.Stmt) error {
			value = stmt.ColumnInt(0)
			return nil
		},
	})
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", name, err)
	}
	return value, nil
}

// Get obtains a connection from the database pool. The caller must return it
// with Put.
func (db *DB) Get(ctx context.Context) (*zsqlite.Conn, error) {
	return db.pool.Take(ctx)
}

// Put returns a connection obtained with Get.
func (db *DB) Put(conn *zsqlite.Conn) {
	db.pool.Put(conn)
}

// Close closes the database pool. A temporary database is discarded when its
// pool is closed.
func (db *DB) Close() error {
	return db.pool.Close()
}
