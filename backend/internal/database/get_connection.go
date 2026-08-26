package database

import (
	"github.com/jmoiron/sqlx"

	"context"
)

func GetDBConnection(ctx context.Context, url string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
