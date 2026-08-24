package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"time"
)

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
    CREATE TABLE papers (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      approved INTEGER NOT NULL DEFAULT 0,
      filename TEXT NOT NULL,
      owner_id INTEGER,
      issue_id INTEGER,
      FOREIGN KEY (owner_id) REFERENCES users(id),
      FOREIGN KEY (issue_id) REFERENCES issues(id)
    );

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
