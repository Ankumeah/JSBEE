package migrations

import (
  "database/sql"
	"context"
)

type v1 struct{}

func init() {
	migrations = append(migrations, v1{})
}

func (v1) Version() uint64 {
	return 1
}

func (v1) Apply(ctx context.Context, tx *sql.Tx) error {
  query := `
    CREATE TABLE users (
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      role TEXT NOT NULL DEFAULT 'viewer' CHECK (role IN (
        'owner',
        'admin',
        'reviewer',
        'viewer'
      ))
    );

    CREATE INDEX idx_users
    ON users(name, role);
  `
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (v1) Revert(ctx context.Context, tx *sql.Tx) error {
  query := `
    ALTER TABLE users
    RENAME TO _users;
  `

  _, err := tx.ExecContext(ctx, query)
  if err != nil {
    return err
  }

  return tx.Commit()
}
