package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"log"
	"os"
	"strconv"
	"time"
)

var envVars = map[string]string{
	"API_VERSION": "",
	"PORT":        "",

	"CACHE_URL":            "",
	"DB_URL":               "",
	"FIREBASE_CREDENTIALS": "",

	"DB_MAX_CONN":      "",
	"DB_MAX_IDLE_CONN": "",
	"DB_MAX_LIFETIME":  "",
	"DB_MAX_IDLE_TIME": "",

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

	var err error
	s.DBMaxConn, err = strconv.Atoi(envVars["DB_MAX_CONN"])
	if err != nil {
		log.Fatalf("Error while parseing DB_MAX_CONN: %v\n", err.Error())
	}
	s.DBMaxIdleConn, err = strconv.Atoi(envVars["DB_MAX_IDLE_CONN"])
	if err != nil {
		log.Fatalf("Error while parseing DB_MAX_IDLE_CONN: %v\n", err.Error())
	}
	s.DBMaxLifetime, err = time.ParseDuration(envVars["DB_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parseing DB_MAX_LIFETIME: %v\n", err.Error())
	}
	s.DBMaxIdleTime, err = time.ParseDuration(envVars["DB_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parseing DB_MAX_IDLE_TIME: %v\n", err.Error())
	}

	s.FrontendSaveDir = envVars["FRONTEND_SAVE_DIR"]
}
