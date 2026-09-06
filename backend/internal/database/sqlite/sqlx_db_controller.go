package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"github.com/jmoiron/sqlx"
)

type SqlxDBController struct{ db *sqlx.DB }

func GetSqlxDBController(db *sqlx.DB) database.DBController {
	return &SqlxDBController{db}
}

func (s *SqlxDBController) DB() *sqlx.DB {
	return s.db
}
