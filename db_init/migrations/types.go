package migrations

import (
	"context"
  "database/sql"
)

type migration interface {
	Version() uint64
	Apply(ctx context.Context, tx *sql.Tx) error
	Revert(ctx context.Context, tx *sql.Tx) error
}
