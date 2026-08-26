package app

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"

	"firebase.google.com/go/v4/auth"
)

// This stuct contains all needed env vars
// in a type safe manner. Any newly added
// env var should also be added here and all env vars are
// to be load with in env.go of the main package
//
// It should not be mutated except in its inital loading
type Config struct {
	Port                string
	APIVersion          string
	DBURL               string
	CacheURL            string
	FireBaseCredentials []byte
	FrontendSaveDir     string
}

// This struct containing shared structs
// needed by the app. Any newly created such struct
// are to be added here and receive their initalisation
// at init.go og the main package
//
// It should not be mutated except in its inital loading
type App struct {
	Config           *Config
	FireBaseClient   *auth.Client
	DBController     database.DBController
	ComponentUpdater *frontend.ComponentUpdater
}
