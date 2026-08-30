package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

// This creats the `state` table
type v3 struct{}

func init() {
	migrations = append(migrations, v3{})
}

func (v3) Version() uint64 {
	return 3
}

func (v3) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := `
    CREATE TABLE state (
      id INT PRIMARY KEY,
      volume INTEGER NOT NULL,
      issue INTEGER NOT NULL,

      CHECK (id = 1)
    );

    INSERT INTO state (volume, issue)
    VALUES (1, 1)
    ON CONFLICT DO NOTHING;

    INSERT INTO migrations (version, applied_at)
    VALUES (3, ?);
  `

	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
