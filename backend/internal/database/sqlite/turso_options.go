//go:build sqlite && turso

package sqlite

import (
	"errors"
	"strings"

	turso "turso.tech/database/tursogo-serverless"
)

const DriverName = "turso-serverless"

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var terr *turso.Error
	if errors.As(err, &terr) {
		if terr.Code == "SQLITE_CONSTRAINT_UNIQUE" ||
			terr.Code == "SQLITE_CONSTRAINT_PRIMARYKEY" ||
			terr.ExtendedCode == "SQLITE_CONSTRAINT_UNIQUE" ||
			terr.ExtendedCode == "SQLITE_CONSTRAINT_PRIMARYKEY" {
			return true
		}
		if terr.Code == "" && terr.ExtendedCode == "" {
			return strings.Contains(terr.Message, "UNIQUE constraint failed")
		}
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "SQLITE_CONSTRAINT_UNIQUE") ||
		strings.Contains(msg, "SQLITE_CONSTRAINT_PRIMARYKEY")
}

func isForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}
	var terr *turso.Error
	if errors.As(err, &terr) {
		if terr.Code == "SQLITE_CONSTRAINT_FOREIGNKEY" ||
			terr.ExtendedCode == "SQLITE_CONSTRAINT_FOREIGNKEY" {
			return true
		}
		if terr.Code == "" && terr.ExtendedCode == "" {
			return strings.Contains(terr.Message, "FOREIGN KEY constraint failed")
		}
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "FOREIGN KEY constraint failed") ||
		strings.Contains(msg, "SQLITE_CONSTRAINT_FOREIGNKEY")
}
