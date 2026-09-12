package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

// This adds the `reviewed` column to the `papers` table
type v4 struct{}

func init() {
	migrations = append(migrations, v4{})
}

func (v4) Version() uint64 {
	return 4
}

func (v4) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := `
    ALTER TABLE papers
    ADD COLUMN reviewed INTEGER NOT NULL DEFAULT 0
    CHECK (reviewed IN (0, 1));

    INSERT INTO migrations (version, applied_at)
    VALUES (4, ?);
  `

	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
