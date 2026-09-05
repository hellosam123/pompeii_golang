package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:embed database.db
var embeddedDB []byte

var (
	db   *sql.DB
	once sync.Once
)

func GetDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		var tmpDir string
		tmpDir, err = os.MkdirTemp("", "pompeii-db-*")
		if err != nil {
			err = fmt.Errorf("failed to create temp dir: %w", err)
			return
		}

		dbPath := filepath.Join(tmpDir, "database.db")
		if err = os.WriteFile(dbPath, embeddedDB, 0o644); err != nil {
			err = fmt.Errorf("failed to create temp dir: %w", err)
			return
		}

		dbPath = fmt.Sprintf("file:%s?mode=ro&cache=shared&immutable=1", dbPath)

		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			err = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err = db.Ping(); err != nil {
			db.Close()
			err = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
