package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"log"
	"os"
)

var envVars = map[string]string{
	"API_VERSION": "",
	"PORT":        "",

	"CACHE_URL":            "",
	"DB_URL":               "",
	"FIREBASE_CREDENTIALS": "",

	"FRONTEND_SAVE_DIR": "",
}

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

func setSettings(s *a.Config) {
	s.APIVersion = envVars["API_VERSION"]
	s.Port = envVars["PORT"]

	s.CacheURL = envVars["CACHE_URL"]
	s.DBURL = envVars["DB_URL"]
	s.FireBaseCredentials = []byte(envVars["FIREBASE_CREDENTIALS"])
	s.FrontendSaveDir = envVars["FRONTEND_SAVE_DIR"]
}
