//go:build sqlite

package main

import (
	"context"
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/sqlite"
	"log"
)

// Initalises connection with DB with sqlite
// Exits program on connection failure
func getDBConnection(ctx context.Context, app *a.App) {
	log.Println("Getting DB connection")
	db, err := database.GetDBConnection(
		ctx,
		app.Config.DBURL,
		sqlite.DriverName,
	)
	if err != nil {
		log.Fatalf("Error while getting db connection: %v\n", err.Error())
	}

	app.DBController = sqlite.GetSqlxDBController(db)
	log.Println("DB connected")
}
