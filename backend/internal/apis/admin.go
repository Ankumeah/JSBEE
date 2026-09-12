package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
	"uuid"
)

// This route handles all admin oprations
func admin(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/admin",
		middlewares.FireBaseAuthMiddleware(app),
		middlewares.GetUserMiddleware(app),
		func(c *gin.Context) { // Checks if user can edit db
			value, ok := c.Get(middlewares.RoleField)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"error": "Internal server error"},
				)
				log.Println("Role field not set")
				return
			}

			role, ok := value.(roles.Role)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"error": "Internal server error"},
				)
				log.Printf("Invalid role: %v\n", value)
				return
			}

			if !role.CanEditDB {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					gin.H{"error": "You are not allowed to preform this action"},
				)
				return
			}

			c.Next()
		},
	)

	group.POST("/increment", func(c *gin.Context) {
		ctx := c.Request.Context()
		field := c.Query("field")

		var err error
		switch field {
		case "volume":
			err = app.DBController.IncrementVolume(ctx)
		case "issue":
			err = app.DBController.IncrementIssue(ctx)
		default:
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Must provide a field to increment"},
			)
			return
		}

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			log.Printf("Error while incrementing %v: %v\n", field, err.Error())
			return
		}

		c.Status(http.StatusNoContent)
	})

	group.GET("/reviewed", func(c *gin.Context) {
		ctx := c.Request.Context()

		papers, err := app.DBController.GetReviewedPapers(ctx)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			log.Printf("Error while getting reviewed papers: %v\n", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"papers": papers})
	})

	group.POST("/publish", func(c *gin.Context) {
		ctx := c.Request.Context()

		count, err := app.DBController.PublishPapers(ctx)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			log.Printf("Error while publishing papers: %v\n", err.Error())
			return
		}

		if count > 0 {
			volumes, err := app.DBController.GetVolumes(ctx)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": "Internal server error"},
				)
				log.Printf("Error while getting volumes: %v\n", err.Error())
				return
			}
			if err := app.ComponentUpdater.UpdateVolumes(ctx, volumes); err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": "Internal server error"},
				)
				log.Printf("Error while updating volumes: %v\n", err.Error())
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
	})

	group.PATCH("/role/:userUUID", func(c *gin.Context) {
		ctx := c.Request.Context()

		userUUID, err := uuid.Parse(c.Param("userUUID"))
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid user UUID"},
			)
			return
		}

		var newRole roles.Role
		if err = newRole.Scan(c.Query("role")); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid role"},
			)
			return
		} else if newRole.Role == roles.Owner.Role {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "Cannot change role to owner"},
			)
			return
		}

		user, err := app.DBController.GetUser(ctx, userUUID)
		if !handleError(c, err) {
			return
		}
		if user.Role.Role == roles.Owner.Role {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "Cannot change the role of owner"},
			)
			return
		}

		err = app.DBController.ChangeRole(ctx, userUUID, newRole)
		if !handleError(c, err) {
			return
		}

		c.Status(http.StatusOK)
	})
}
