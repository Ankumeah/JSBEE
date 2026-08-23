package app

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"

	"firebase.google.com/go/v4/auth"

	"time"
)

type Config struct {
	Port                string
	APIVersion          string
	DBURL               string
	CacheURL            string
	FireBaseCredentials []byte
	DBMaxConn           int
	DBMaxIdleConn       int
	DBMaxLifetime       time.Duration
	DBMaxIdleTime       time.Duration
	FrontendSaveDir     string
}

type App struct {
	Config           *Config
	FireBaseClient   *auth.Client
	DBController     database.DBController
	ComponentUpdater *frontend.ComponentUpdater
}
