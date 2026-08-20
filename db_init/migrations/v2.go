package migrations

import (
  "database/sql"
	"context"
)

type v2 struct{}

func init() {
	migrations = append(migrations, v2{})
}

func (v2) Version() uint64 {
	return 2
}

func (v2) Apply(ctx context.Context, tx *sql.Tx) error {
  query := `
    CREATE TABLE papers (
      id SERIAL PRIMARY KEY,
      title TEXT NOT NULL,
      approved BOOLEAN NOT NULL DEFAULT false,
      filename TEXT NOT NULL,
      owner_id INTEGER,
      issue_id INTEGER,
      FOREIGN KEY (owner_id) REFERENCES users(id),
      FOREIGN KEY (issue_id) REFERENCES issues(id)
    );
  `
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (v2) Revert(ctx context.Context, tx *sql.Tx) error {
  query := `
    ALTER TABLE papers
    RENAME TO _papers;
  `

  _, err := tx.ExecContext(ctx, query)
  if err != nil {
    return err
  }

  return tx.Commit()
}
