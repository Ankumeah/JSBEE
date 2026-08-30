package main

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"context"
	"log"
)

// Initalises connection with firebase
// Exits program on connection failure
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

// Creates the ComponentUpdater struct
// Creates the save dir if it dosent exist
// Exits program on any errors with mkdir
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
