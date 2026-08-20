package main

import (
	"context"
	"database/sql"
	"time"
)

type SqlConfig struct {
	maxConns    int
	maxIdle     int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

func NewSqlConfig(
	maxConn int,
	maxIdle int,
	maxLifetime time.Duration,
	maxIdleTime time.Duration,
) SqlConfig {
	return SqlConfig{
		maxConns:    maxConn,
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
		maxIdleTime: maxIdleTime,
	}
}

func GetDBConnection(ctx context.Context, url string, config SqlConfig) (*sql.DB, error) {
	db, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.maxConns)
	db.SetMaxIdleConns(config.maxIdle)
	db.SetConnMaxLifetime(config.maxLifetime)
	db.SetConnMaxIdleTime(config.maxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
