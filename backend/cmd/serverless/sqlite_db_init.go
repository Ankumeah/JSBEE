//go:build sqlite

package main

import (
	"context"
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/sqlite"

	"os"
)

// Initalises connection with DB with sqlite
// Exits program on connection failure
func getDBConnection(ctx context.Context, app *a.App) {
	app.Logger.InfoContext(ctx, "Getting DB connection")
	db, err := database.GetDBConnection(
		ctx,
		app.Config.DBURL,
		sqlite.DriverName,
	)
	if err != nil {
		app.Logger.ErrorContext(ctx,
			"Error while getting db connection: "+err.Error(),
		)
		os.Exit(1)
	}

	app.DBController = sqlite.GetSqlxDBController(db)
	app.Logger.InfoContext(ctx, "DB connected")
}
