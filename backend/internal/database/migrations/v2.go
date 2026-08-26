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
    CREATE TABLE papers (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      approved INTEGER NOT NULL DEFAULT 0 CHECK (approved IN (0, 1)),
      number INTEGER NOT NULL,
      filename TEXT NOT NULL UNIQUE,
      volume INTEGER NOT NULL CHECK (volume > 0),
      issue INTEGER NOT NULL CHECK (issue > 0),

      owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,

      UNIQUE(title, owner_id),
      UNIQUE(volume, number, issue)
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
