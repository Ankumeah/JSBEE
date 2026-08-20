package database

import (
	"github.com/jmoiron/sqlx"
)

type sqlxDBController struct { db *sqlx.DB }
func SqlxDBController(db *sqlx.DB) *DBController {
  return &DBController{
    &sqlxDBController{db},
  }
}
