package main

import (
	"log"
	"os"
)

func loadEnv() {
	for env, _ := range envVars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Env var unset: %v\n", env)
		}
		envVars[env] = _env
	}
}
