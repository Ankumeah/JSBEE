package migrations

import (
  "database/sql"
	"context"
)

type v3 struct{}

func init() {
	migrations = append(migrations, v3{})
}

func (v3) Version() uint64 {
	return 3
}

func (v3) Apply(ctx context.Context, tx *sql.Tx) error {
  query := `
    CREATE TABLE issues (
      id SERIAL PRIMARY KEY,
      number INTEGER NOT NULL,
      volume INTEGER NOT NULL,
      filename TEXT NOT NULL,
    );
  `
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (v3) Revert(ctx context.Context, tx *sql.Tx) error {
  query := `
    ALTER TABLE issues
    RENAME TO _issues;
  `

  _, err := tx.ExecContext(ctx, query)
  if err != nil {
    return err
  }

  return tx.Commit()
}
