package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Open connects to the SQLite database at path — creating the file if it
// doesn't exist, which the driver does automatically on first connection —
// then ensures all tables exist before returning.
func Open(ctx context.Context, path string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	sqlDB.SetMaxOpenConns(1) // SQLite: avoid concurrent-writer lock errors

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	if err := migrate(ctx, sqlDB); err != nil {
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	return sqlDB, nil
}
