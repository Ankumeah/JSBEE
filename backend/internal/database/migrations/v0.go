package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

type v0 struct{}

func (v0) Version() uint64 {
	return 0
}

func (v0) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := tx.Rebind(`
    CREATE TABLE IF NOT EXISTS migrations (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      version BIGINT UNIQUE NOT NULL,
      applied_at BIGINT NOT NULL
    );

    INSERT INTO migrations (version, applied_at)
    VALUES (0, ?)
    ON CONFLICT(version) DO NOTHING;
  `)
	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
