package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
)

func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	var matchingVersions = false
	slices.SortFunc(migrations,
		func(a migration, b migration) int {
			if a.Version() > b.Version() {
				return 1
			} else if a.Version() < b.Version() {
				return -1
			} else { // Just put the fries in the bag bro
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

	for _, m := range migrations {
		err := func() error {
			tx, err := db.BeginTx(ctx, &sql.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback()

			err = m.Apply(ctx, tx)
			if err != nil {
				return fmt.Errorf("Migration %v: %w", m.Version(), err)
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}

	return nil
}
