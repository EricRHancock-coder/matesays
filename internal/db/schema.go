// Package db manages db creation and schema creeation
package db

import (
	"context"
	"database/sql"
)

const createQuotesTable = `
	CREATE TABLE  if not exists quotes ( 
		id integer primary key autoincrement, 
		quote string, 
		created_time datetime default current_timestamp 
	);
`

// can construct more schema related constants here

// migrate creates every table if it doesn't already exist. Unexported —
// only Open calls it, so bootstrapping can never be skipped or duplicated
// by accident.
func migrate(ctx context.Context, sqlDB *sql.DB) error {
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// append schema constant after createQuotesTable constant ie ',createAuditTable'
	for _, stmt := range []string{createQuotesTable} {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return tx.Commit()
}
