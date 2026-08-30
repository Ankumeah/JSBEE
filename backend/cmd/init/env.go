package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"log"
	"os"
)

// Temprory map to store type unsafe env vars
var envVars = map[string]string{
	"DB_URL":            "",
	"FRONTEND_SAVE_DIR": "",
}

// Loads all needed env vars into passsed config.
// Exits program on first unset env var
func loadEnv(s *a.Config) {
	log.Println("Loading env")

	for env := range envVars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}
		envVars[env] = _env
	}
	setSettings(s)

	log.Println("Loaded env")
}

// Loads the type unsafe map into the type safe config
func setSettings(s *a.Config) {
	s.DBURL = envVars["DB_URL"]
	s.FrontendSaveDir = envVars["FRONTEND_SAVE_DIR"]
}
