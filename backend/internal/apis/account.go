package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"uuid"
)

// This route deals with account logic
func account(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/account", middlewares.FireBaseAuthMiddleware(app))

	// This route handles creating of a user
	group.POST("", func(c *gin.Context) {
		ctx := c.Request.Context()
		name := c.GetString(middlewares.NameField)
		email := c.GetString(middlewares.EmailField)

		userUUID, err := uuid.Parse(c.GetString(middlewares.UUIDFeild))
		if !handleError(c, err) {
			return
		}

		subscribedString := c.Query("sub")
		subscribed, err := strconv.ParseBool(subscribedString)
		if subscribedString == "" {
			subscribed = false
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := app.DBController.AddUser(ctx, database.User{
			UUID:       userUUID,
			Name:       name,
			Email:      email,
			Subscribed: subscribed,
		}); !handleError(c, err) {
			return
		}

		c.Status(http.StatusCreated)
	})

	// This route handles deletion of a user
	group.DELETE("", func(c *gin.Context) {
		ctx := c.Request.Context()

		userUUID, err := uuid.Parse(c.GetString(middlewares.UUIDFeild))
		if !handleError(c, err) {
			return
		}

		if err := app.DBController.DeleteUser(ctx, userUUID); !handleError(c, err) {
			return
		}
	})

	// TODO: Add PATCH sometime later
}
