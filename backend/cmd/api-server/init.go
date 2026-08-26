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
	var err error
	app.ComponentUpdater, err = frontend.GetComponentUpdater(
		app.Config.FrontendSaveDir,
	)
	if err != nil {
		log.Fatalf("Error while getting component updater: %v\n", err.Error())
	}
	log.Println("Got component updater")
}

func generateInitalComponents(ctx context.Context, app *a.App) {
	log.Println("Generating inital components")

	volumes, err := app.DBController.GetVolumes(ctx)
	if err != nil {
		log.Fatalf("Error while getting volumes: %v\n", err.Error())
	}

	if err := app.ComponentUpdater.UpdateAll(ctx, volumes); err != nil {
		log.Fatalf("Error while generating inital components: %v\n", err.Error())
	}

	log.Println("Generated inital components")
}
