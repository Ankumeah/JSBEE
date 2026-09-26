package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore"

	"context"
	"log/slog"
	"os"
)

// Initalises the logger
// Should be initalised before anything
// else as almost everything can have errors
func initLogger(app *a.App) {
	app.Logger = slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)
}

// Initalises connection with firebase
// Exits program on connection failure
func initFirebase(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Initializeing firebase client")
	var err error
	app.FireBaseClient, err = middlewares.InitFirebase(
		ctx, app.Config.FireBaseCredentials,
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting firebase client: "+err.Error(),
		)
		os.Exit(1)
	}
	app.Logger.InfoContext(ctx, "Firebase client initialized")
}

// Creates the ComponentUpdater struct
// Creates the save dir if it dosent exist
// Exits program on any errors with mkdir
func getComponentUpdater(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Getting component updater")
	var err error
	app.ComponentUpdater, err = frontend.GetComponentUpdater(
		app.Config.FrontendSaveDir,
		app.Config.FireBaseClientConfig,
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting component updater: "+err.Error(),
		)
		os.Exit(1)
	}
	app.Logger.InfoContext(ctx, "Got component updater")
}

func connectObjectStore(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Connecting to object store")

	config, err := objectstore.NewS3StaticConfigFromJSON(
		app.Config.ObjectStoreConfig,
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting s3 config: "+err.Error(),
		)
		os.Exit(1)
	}

	app.ObjectStore, err = objectstore.GetStaticS3Client(
		ctx, config,
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting s3 client: "+err.Error(),
		)
		os.Exit(1)
	}

	app.Logger.InfoContext(ctx, "Connected to object store")
}
