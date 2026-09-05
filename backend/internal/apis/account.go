package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"uuid"
)

// This route deals with account logic
func account(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/account", middlewares.FireBaseAuthMiddleware(app))

	// This route handles creates a user
	// Techinally it just sets the subscribed
	// value in the db as the `FireBaseAuthMiddleware`
	// creates the user
	group.POST("", func(c *gin.Context) {
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
		// `FireBaseAuthMiddleware` alreday created the user
		// and subscribed is the only unset var, and false is
		// the default, if user alreday wanted it to be false
		// no need to further waste compute
		if !subscribed {
			c.Status(http.StatusCreated)
		}

		ctx := c.Request.Context()

		userUUID, err := uuid.Parse(c.GetString(middlewares.UUIDFeild))
		if !handleError(c, err) {
			return
		}

		// Set subscribed of add as `FireBaseAuthMiddleware`
		// user exists in db but subscribed isent set
		if err := app.DBController.SetSubscription(
			ctx, subscribed, userUUID,
		); !handleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, gin.H{"uuid": userUUID})
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

		c.Status(http.StatusNoContent)
	})

	// TODO: Add PATCH sometime later
}
