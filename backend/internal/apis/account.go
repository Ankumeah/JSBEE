package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"net/http"
)

// This route deals with account logic
func account(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/account", middlewares.FireBaseAuthMiddleware(app))

	// This route handles creating of a user
	group.POST("", func(c *gin.Context) {
		ctx := c.Request.Context()
		name := c.GetString(middlewares.NameField)
		email := c.GetString(middlewares.EmailField)

		if err := app.DBController.AddUser(ctx, database.User{
			Name:  name,
			Email: email,
		}); !handleError(c, err) {
			return
		}

		c.Status(http.StatusCreated)
	})

	// This route handles deletion of a user
	group.DELETE("", func(c *gin.Context) {
		ctx := c.Request.Context()
		email := c.GetString(middlewares.EmailField)

		if err := app.DBController.DeleteUser(ctx, email); !handleError(c, err) {
			return
		}
	})

	// TODO: Add PATCH sometime later
}
