package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"context"
	"encoding/json"
	"errors"
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
func loadEnv(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Loading env")

	for env := range envVars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			app.Logger.ErrorContext(ctx, "Unset env var: "+env)
			os.Exit(1)
		}
		envVars[env] = _env
	}
	setSettings(ctx, app)

	app.Logger.InfoContext(ctx, "Loaded env")
}

// Loads the type unsafe map into the type safe config
func setSettings(ctx context.Context, app *a.App) {
	app.Config.APIVersion = envVars["API_VERSION"]
	app.Config.Port = envVars["BACKEND_PORT"]

	app.Config.DBURL = envVars["DB_URL"]

	app.Config.FireBaseCredentials = []byte(envVars["FIREBASE_CREDENTIALS"])
	app.Config.FrontendSaveDir = envVars["FRONTEND_SAVE_DIR"]
	app.Config.ObjectStoreConfig = []byte(envVars["OBJECT_STORE_CONFIG"])

	firebaseClientConfig, err := jsObjectFromJSON(envVars["FIREBASE_CLIENT_CONFIG"])
	if err != nil {
		app.Logger.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
	app.Config.FireBaseClientConfig = firebaseClientConfig
}

// jsObjectFromJSON parses a JSON object string so that templ can render it
func jsObjectFromJSON(s string) (map[string]any, error) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return nil, errors.New("Invalid FIREBASE_CLIENT_CONFIG: " + err.Error())
	}
	return obj, nil
}
