package main

import (
	"github.com/Ankumeah/JSBEE/db_init/migrations"

	"context"
	"log"
	"os"
	"strconv"
	"time"
)

var Ctx = context.Background()
var Config SqlConfig
var DBURL string

func loadEnv() {
	var envVars = map[string]string{
		"DB_URL":            "",
		"DB_MAX_CONN":      "",
		"DB_MAX_IDLE_CONN": "",
		"DB_MAX_LIFETIME":   "",
		"DB_MAX_IDLE_TIME":  "",
	}

	for env := range envVars {
		Env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}

		envVars[env] = Env
	}

	DBURL = envVars["DB_URL"]
	maxConn, err := strconv.Atoi(envVars["DB_MAX_CONN"])
	if err != nil {
		log.Fatalf("Error while parsing DB_MAX_CONN: %v\n", err.Error())
	}
	maxIdle, err := strconv.Atoi(envVars["DB_MAX_IDLE_CONN"])
	if err != nil {
		log.Fatalf("Error while parsing DB_MAX_IDLE_CONN: %v\n", err.Error())
	}
	maxLifetime, err := time.ParseDuration(envVars["DB_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parsing DB_MAX_LIFETIME: %v\n", err.Error())
	}
	maxIdleTime, err := time.ParseDuration(envVars["DB_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parsing DB_MAX_IDLE_TIME: %v\n", err.Error())
	}

	Config = NewSqlConfig(int(maxConn), int(maxIdle), maxLifetime, maxIdleTime)
}

func main() {
	loadEnv()

	db, err := GetDBConnection(Ctx, DBURL, Config)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v\n", err.Error())
	}

	err = migrations.ApplyMigrations(Ctx, db)
	if err != nil {
		log.Fatalf("Error while applying migration: %v\n", err.Error())
	}
}
