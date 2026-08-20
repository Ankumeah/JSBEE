package main

import (
	"github.com/Ankumeah/JSBEE/backend/internal/apis"
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"context"
	"log"
  "sync"
)

var Ctx = context.Background()
var app = &a.App{Config: &a.Config{}}

func main() {
	loadEnv(app.Config)

  var wg sync.WaitGroup
	wg.Go(func(){initFirebase(Ctx, app)})
  wg.Go(func(){getDBConnection(Ctx, app)})
  wg.Wait()

	log.Println("Starting http server")
	r := gin.Default()
	apiGroup := r.Group(
		"/api/"+app.Config.APIVersion+"/",
		middlewares.LogMiddleware(),
	)
	apis.Apis(apiGroup, app)

	log.Println("Running backend on port: " + app.Config.Port)
	log.Println("API_VERSION: " + app.Config.APIVersion)

	err := r.Run(":" + app.Config.Port)
	if err != nil {
		log.Fatalf("Error while running server: %v\n", err.Error())
	}
}
