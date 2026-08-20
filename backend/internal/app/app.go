package app

import (
  "github.com/Ankumeah/JSBEE/backend/internal/database"

	"firebase.google.com/go/v4/auth"

  "time"
)

type Config struct {
	Port                string
	APIVersion          string
	DBURL               string
	CacheURL            string
	FireBaseCredentials []byte
	DBMaxConn int
	DBMaxIdleConn int
	DBMaxLifetime time.Duration
	DBMaxIdleTime time.Duration
}

type App struct {
	Config         *Config
	FireBaseClient *auth.Client
  DBController *database.DBController
}
