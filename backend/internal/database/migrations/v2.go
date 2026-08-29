package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

// This creats the `papers` table
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
      uuid TEXT NOT NULL UNIQUE CHECK (uuid != ''),
      title TEXT NOT NULL CHECK (title != ''),
      number INTEGER,
      filename TEXT NOT NULL UNIQUE CHECK (filename != ''),
      volume INTEGER CHECK (volume > 0),
      issue INTEGER CHECK (issue > 0),

      owner_uuid INTEGER REFERENCES users(uuid) ON DELETE SET NULL,

      UNIQUE(title, owner_uuid),
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
