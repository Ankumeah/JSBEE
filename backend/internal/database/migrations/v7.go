package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

// This adds the `location` and `category` columns to the `papers` table.
// Both are required for new submissions (enforced by the API); existing
// rows keep an empty string so the migration applies to old databases.
type v7 struct{}

func init() {
	migrations = append(migrations, v7{})
}

func (v7) Version() uint64 {
	return 7
}

func (v7) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	if _, err := tx.ExecContext(
		ctx, tx.Rebind(`
      ALTER TABLE papers
      ADD COLUMN location TEXT NOT NULL DEFAULT ''
      CHECK (length(location) <= 120);
    `),
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx, tx.Rebind(`
      ALTER TABLE papers
      ADD COLUMN category TEXT NOT NULL DEFAULT ''
      CHECK (length(category) <= 120);
    `),
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx, tx.Rebind(`
      INSERT INTO migrations (version, applied_at)
      VALUES (7, ?);
    `), time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
