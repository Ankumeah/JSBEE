package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"context"
	"log"
)

var Ctx = context.Background()
var app = &a.App{Config: &a.Config{}}

func main() {
	log.Println("Init started")

	loadEnv(app.Config)
	getComponentUpdater(app)
	getDBConnection(Ctx, app)
	runDBMigrations(Ctx, app)
	generateInitalComponents(Ctx, app)
	saveAssets(Ctx, app)

	log.Println("Init completed")
}
