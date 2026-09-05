package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/assets"

	"context"
	"log"
)

// Runs DB migrations
// Exits program on any errors
func runDBMigrations(ctx context.Context, app *a.App) {
	log.Println("Running DB migrations")

	if err := migrations.ApplyMigrations(
		ctx, app.DBController.DB(),
	); err != nil {
		log.Fatalf("Error while running DB migrations: %v\n", err)
	}

	log.Println("DB migrations completed")
}

// Creates the ComponentUpdater struct
// Creates the save dir if it dosent exist
// Exits program on any errors with mkdir
func getComponentUpdater(app *a.App) {
	log.Println("Getting component updater")
	var err error
	app.ComponentUpdater, err = frontend.GetComponentUpdater(
		app.Config.FrontendSaveDir,
		app.Config.FireBaseClientConfig,
	)
	if err != nil {
		log.Fatalf("Error while getting component updater: %v\n", err.Error())
	}
	log.Println("Got component updater")
}

// Asset files may exist at launch, this
// writes the bundled asset files
// Exits on any error
func saveAssets(ctx context.Context, app *a.App) {
	log.Println("Saveing assets")

	if err := assets.SaveAssets(
		app.Config.FrontendSaveDir,
	); err != nil {
		log.Fatalf("Error while saveing assets: %v\n", err.Error())
	}

	log.Println("Assets saved")
}

// Component files may not exist at launch, this generates
// them and updates them if they alreday exist
// Exits on any error
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
