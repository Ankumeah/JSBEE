package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

type v2 struct{}

func init() {
	migrations = append(migrations, v2{})
}

func (v2) Version() uint64 {
	return 2
}

func (v2) Apply(
	ctx context.Context,
	tx *sqlx.Tx,
) error {
	query := `
    CREATE TABLE issues (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      number INTEGER NOT NULL,
      volume INTEGER NOT NULL,
      filename TEXT NOT NULL
    );

    INSERT INTO migrations (version, applied_at)
    VALUES (2, ?);
  `

	if _, err := tx.ExecContext(
		ctx, query, time.Now().Unix(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
