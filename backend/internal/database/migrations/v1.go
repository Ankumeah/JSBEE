package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

type v1 struct{}

func init() {
	migrations = append(migrations, v1{})
}

func (v1) Version() uint64 {
	return 1
}

func (v1) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := tx.Rebind(`
    CREATE TABLE users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      email TEXT NOT NULL,
      role TEXT NOT NULL DEFAULT 'viewer' CHECK (role IN (
        'owner',
        'admin',
        'reviewer',
        'viewer'
      ))
    );

    CREATE INDEX idx_users
    ON users(name, role);

    INSERT INTO migrations (version, applied_at)
    VALUES (1, ?);
  `)
	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
