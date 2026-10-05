package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
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

func connectObjectStore(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Connecting to object store")

	var err error
	app.ObjectStore, err = objectstore.GetObjectStore(
		ctx, string(app.Config.ObjectStoreConfig),
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting object store client: "+err.Error(),
		)
		os.Exit(1)
	}

	app.Logger.InfoContext(ctx, "Connected to object store")
}
