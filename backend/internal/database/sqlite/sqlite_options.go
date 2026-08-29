//go:build sqlite

package sqlite

import (
	"modernc.org/sqlite"
	"modernc.org/sqlite/lib"

	"errors"
)

const DriverName = "sqlite"

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return true
	}
	return false
}

func isForeignKeyViolation(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY {
		return true
	}
	return false
}
