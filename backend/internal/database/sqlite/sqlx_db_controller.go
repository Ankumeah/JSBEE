package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"github.com/jmoiron/sqlx"
)

type SqlxDBController struct{ db *sqlx.DB }

func GetSqlxDBController(db *sqlx.DB) database.DBController {
	return &SqlxDBController{db}
}

// This exposes the raw underlying DB object.
// This is only to be used to execute migrations
func (s *SqlxDBController) DB() *sqlx.DB {
	return s.db
}
