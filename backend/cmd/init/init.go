package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/assets"
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore"
	"github.com/Ankumeah/JSBEE/backend/internal/provider"

	"bytes"
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

	if err := app.ComponentUpdater.UpdateAll(
		ctx, volumes, app.ObjectStore.PublicBaseURL(),
	); err != nil {
		log.Fatalf("Error while generating inital components: %v\n", err.Error())
	}

	log.Println("Generated inital components")
}

func connectObjectStore(ctx context.Context, app *a.App) {
	log.Println("Connecting to object store")

	config, err := objectstore.NewS3StaticConfigFromJSON(
		app.Config.ObjectStoreConfig,
	)
	if err != nil {
		log.Fatalf("Error while getting s3 config: %v\n", err.Error())
	}

	app.ObjectStore, err = objectstore.GetStaticS3Client(
		ctx, config,
	)
	if err != nil {
		log.Fatalf("Error while getting s3 client: %v\n", err.Error())
	}

	log.Println("Connected to object store")
}

func initObjectStore(ctx context.Context, app *a.App) {
	log.Println("Initalizing object store")

	if err := app.ObjectStore.Init(ctx); err != nil {
		log.Fatalf("Error while initalizing object store: %v\n", err.Error())
	}

	log.Println("Object store initalized")
}

// Creates the about us page in the object store
// if it dosent exist yet, exits on any error
func seedAboutPage(ctx context.Context, app *a.App) {
	log.Println("Seeding about page")

	content, err := app.ObjectStore.GetFile(ctx, provider.AboutFilename)
	if err == nil {
		content.Close()
		log.Println("About page already exists")
		return
	}

	buf := bytes.NewBufferString(provider.AboutSeedContent)
	if err := app.ObjectStore.AddFile(
		ctx, provider.AboutFilename, buf, int64(buf.Len()),
	); err != nil {
		log.Fatalf("Error while seeding about page: %v\n", err.Error())
	}
	if err := app.ObjectStore.PublicFile(ctx, provider.AboutFilename); err != nil {
		log.Fatalf("Error while publishing about page: %v\n", err.Error())
	}

	log.Println("About page seeded")
}
