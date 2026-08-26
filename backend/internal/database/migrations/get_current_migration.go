package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
)

// This function returns the max migration
// present in the DB
func getCurrentMigration(
	ctx context.Context,
	db *sqlx.DB,
) (uint64, error) {
	query := `SELECT MAX(version) FROM migrations;`

	var version uint64
	if err := db.QueryRowContext(
		ctx, query,
	).Scan(&version); errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("No migration found")
	} else if err != nil {
		return 0, err
	}

	return version, nil
}
