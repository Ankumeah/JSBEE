package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"encoding/json"
	"log"
	"os"
)

// Temprory map to store type unsafe env vars
var envVars = map[string]string{
	"API_VERSION":  "",
	"BACKEND_PORT": "",

	"DB_URL": "",

	"OBJECT_STORE_CONFIG":    "",
	"FIREBASE_CREDENTIALS":   "",
	"FIREBASE_CLIENT_CONFIG": "",

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
	s.APIVersion = envVars["API_VERSION"]
	s.Port = envVars["BACKEND_PORT"]

	s.DBURL = envVars["DB_URL"]

	s.FireBaseCredentials = []byte(envVars["FIREBASE_CREDENTIALS"])
	s.FireBaseClientConfig = jsObjectFromJSON(envVars["FIREBASE_CLIENT_CONFIG"])
	s.FrontendSaveDir = envVars["FRONTEND_SAVE_DIR"]
	s.ObjectStoreConfig = []byte(envVars["OBJECT_STORE_CONFIG"])
}

// jsObjectFromJSON parses a JSON object string so that templ can render it
func jsObjectFromJSON(s string) map[string]any {
	var obj map[string]any
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		log.Fatalf("Invalid FIREBASE_CLIENT_CONFIG: %v", err)
	}
	return obj
}
