package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"net/http"
)

func account(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/account", middlewares.FireBaseAuthMiddleware(app))

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

	group.DELETE("", func(c *gin.Context) {
		ctx := c.Request.Context()
		email := c.GetString(middlewares.EmailField)

		if err := app.DBController.DeleteUser(ctx, email); !handleError(c, err) {
			return
		}
	})

	// TODO: Add PATCH sometime later
}
