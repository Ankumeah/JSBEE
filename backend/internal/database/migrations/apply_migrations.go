package migrations

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
)

// This applies checks current DB migration version
// and all applies all migrations having a higher version.
// All migrations are run within a DB transaction
//
// # This function is safe to run on a new DB
//
// This function will error on any DB error or
// if any two migrations have the same number
func ApplyMigrations(ctx context.Context, db *sqlx.DB) error {
	var matchingVersions = false

	// Sort the migrations as they may be out of order as they
	// are registered via init functions
	slices.SortFunc(migrations,
		func(a migration, b migration) int {
			if a.Version() > b.Version() {
				return 1
			} else if a.Version() < b.Version() {
				return -1
			} else { // Just put the fries in the bag bro
				// We use this to exit the function later
				matchingVersions = true
				log.Printf(
					"Matching migration numbers: %v and %v\n",
					a.Version(), b.Version(),
				)
				return 0
			}
		},
	)
	if matchingVersions {
		return errors.New("Two migrations have matching version numbers")
	}

	// v0 is applied no matter what so that
	// `getCurrentMigration` does not fail
	err := v0{}.Apply(ctx, db)
	if err != nil {
		return fmt.Errorf("Migration 0: %w", err)
	}

	currentVersion, err := getCurrentMigration(ctx, db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.Version() <= currentVersion {
			continue
		}

		if err := func() error {
			tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback()

			err = m.Apply(ctx, tx)
			if err != nil {
				return fmt.Errorf("Migration %v: %w", m.Version(), err)
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}
