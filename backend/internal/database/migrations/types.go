package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
)

type migration interface {
	Version() uint64
	Apply(ctx context.Context, tx *sqlx.Tx) error
}
