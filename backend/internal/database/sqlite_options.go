//go:build sqlite

package database

import (
	"modernc.org/sqlite"

	"errors"
)

const driverName = "sqlite"

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == 1555 {
		return true
	}
	return false
}
