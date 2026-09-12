package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

// This adds the `blogs` table
type v5 struct{}

func init() {
	migrations = append(migrations, v5{})
}

func (v5) Version() uint64 {
	return 5
}

func (v5) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := `
    CREATE TABLE blogs (
      uuid TEXT PRIMARY KEY,
      title TEXT NOT NULL CHECK (length(title) > 0),
      filename TEXT NOT NULL UNIQUE,
      created_at INTEGER NOT NULL,
      updated_at INTEGER NOT NULL
    );

    INSERT INTO migrations (version, applied_at)
    VALUES (5, ?);
  `

	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
