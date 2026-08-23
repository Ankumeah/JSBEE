package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"context"
	"log"
)

func initFirebase(ctx context.Context, app *a.App) {
	log.Println("Initializeing firebase client")
	var err error
	app.FireBaseClient, err = middlewares.InitFirebase(
		ctx, app.Config.FireBaseCredentials,
	)
	if err != nil {
		log.Fatalf("Error while getting firebase client: %v\n", err.Error())
	}
	log.Println("Firebase client initialized")
}

func getDBConnection(ctx context.Context, app *a.App) {
	log.Println("Getting DB connection")
	db, err := database.GetDBConnection(ctx, app.Config.DBURL,
		database.NewSqlConfig(
			app.Config.DBMaxConn,
			app.Config.DBMaxIdleConn,
			app.Config.DBMaxLifetime,
			app.Config.DBMaxIdleTime,
		),
	)
	if err != nil {
		log.Fatalf("Error while getting db connection: %v\n", err.Error())
	}

	app.DBController = database.GetSqlxDBController(db)
	log.Println("DB connected")
}

func runDBMigrations(ctx context.Context, app *a.App) {
	log.Println("Running DB migrations")

	if err := migrations.ApplyMigrations(
		ctx, app.DBController.DB(),
	); err != nil {
		log.Fatalf("Error while running DB migrations: %v\n", err)
	}

	log.Println("DB migrations completed")
}

func getComponentUpdater(app *a.App) {
	log.Println("Getting component updater")
	app.ComponentUpdater = frontend.GetComponentUpdater(
		app.Config.FrontendSaveDir,
	)
	log.Println("Got component updater")
}
