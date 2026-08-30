package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"context"
	"log"
	"sync"
)

var Ctx = context.Background()
var app = &a.App{Config: &a.Config{}}

func main() {
	log.Println("Init started")

	loadEnv(app.Config)

	var wg sync.WaitGroup
	wg.Go(func() { getComponentUpdater(app) })
	wg.Go(func() {
		getDBConnection(Ctx, app)
		runDBMigrations(Ctx, app)

		generateInitalComponents(Ctx, app)
		saveAssets(Ctx, app)
	})
	wg.Wait()

	log.Println("Init completed")
}
