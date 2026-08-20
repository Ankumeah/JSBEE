package migrations

import (
	"context"
  "database/sql"
	"errors"
)

type v0 struct{}

func init() {
	migrations = append(migrations, v0{})
}

func (v0) Version() uint64 {
	return 0
}

func (v0) Apply(ctx context.Context, tx *sql.Tx) error {
  query := `
    CREATE TABLE migration (
      id SERIAL PRIMARY KEY,
      version BIGINT NOT NULL,
      applied_at BIGINT NOT NULL
    );
  `
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (v0) Revert(ctx context.Context, tx *sql.Tx) error {
	return errors.New("Cannot revert migration v0")
}
