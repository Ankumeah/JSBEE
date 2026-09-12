package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"
	"github.com/Ankumeah/JSBEE/backend/internal/provider"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"github.com/gin-gonic/gin"

	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
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
			if err := app.ComponentUpdater.UpdateVolumes(
				ctx, volumes, app.ObjectStore.PublicBaseURL(),
			); err != nil {
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

	group.POST("/blog", func(c *gin.Context) {
		ctx := c.Request.Context()

		var body struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Must provide a title and content"},
			)
			return
		}
		if strings.TrimSpace(body.Title) == "" {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Title is required"},
			)
			return
		}
		if len(body.Content) == 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Content is required"},
			)
			return
		}
		if int64(len(body.Content)) > maxBlogSize {
			c.JSON(
				http.StatusRequestEntityTooLarge,
				gin.H{"error": fmt.Sprintf(
					"Content too large, max size is %d bytes", maxBlogSize,
				)},
			)
			return
		}

		blogUUID := uuid.New()
		now := time.Now().Unix()
		// Include the uuid so two posts in the same second never collide
		filename := fmt.Sprintf("blog-%d-%s.md", now, blogUUID.String())

		buf := bytes.NewBufferString(body.Content)
		if err := app.ObjectStore.AddFile(
			ctx, filename, buf, int64(buf.Len()),
		); !handleError(c, err) {
			return
		}
		if err := app.ObjectStore.PublicFile(
			ctx, filename,
		); !handleError(c, err) {
			app.ObjectStore.DeleteFile(ctx, filename)
			return
		}

		if err := app.DBController.AddBlog(ctx, database.Blog{
			UUID:      blogUUID,
			Title:     strings.TrimSpace(body.Title),
			Filename:  filename,
			CreatedAt: now,
			UpdatedAt: now,
		}); !handleError(c, err) {
			app.ObjectStore.DeleteFile(ctx, filename)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"uuid": blogUUID})
	})

	group.PATCH("/blog/:blogUUID", func(c *gin.Context) {
		ctx := c.Request.Context()

		blogUUID, err := uuid.Parse(c.Param("blogUUID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
			return
		}

		var body struct {
			Title   *string `json:"title"`
			Content *string `json:"content"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Must provide a title or content"},
			)
			return
		}
		if body.Title == nil && body.Content == nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Nothing to update"},
			)
			return
		}

		blog, err := app.DBController.GetBlog(ctx, blogUUID)
		if !handleError(c, err) {
			return
		}

		if body.Title != nil {
			if strings.TrimSpace(*body.Title) == "" {
				c.JSON(
					http.StatusBadRequest,
					gin.H{"error": "Title is required"},
				)
				return
			}
			blog.Title = strings.TrimSpace(*body.Title)
		}

		if body.Content != nil {
			if len(*body.Content) == 0 {
				c.JSON(
					http.StatusBadRequest,
					gin.H{"error": "Content is required"},
				)
				return
			}
			if int64(len(*body.Content)) > maxBlogSize {
				c.JSON(
					http.StatusRequestEntityTooLarge,
					gin.H{"error": fmt.Sprintf(
						"Content too large, max size is %d bytes", maxBlogSize,
					)},
				)
				return
			}

			buf := bytes.NewBufferString(*body.Content)
			if err := app.ObjectStore.AddFile(
				ctx, blog.Filename, buf, int64(buf.Len()),
			); !handleError(c, err) {
				return
			}
			if err := app.ObjectStore.PublicFile(
				ctx, blog.Filename,
			); !handleError(c, err) {
				return
			}
		}

		blog.UpdatedAt = time.Now().Unix()
		if err := app.DBController.UpdateBlog(ctx, blog); !handleError(c, err) {
			return
		}

		c.Status(http.StatusOK)
	})

	group.DELETE("/blog/:blogUUID", func(c *gin.Context) {
		ctx := c.Request.Context()

		blogUUID, err := uuid.Parse(c.Param("blogUUID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
			return
		}

		blog, err := app.DBController.GetBlog(ctx, blogUUID)
		if !handleError(c, err) {
			return
		}

		// Delete the DB row first so a DB failure never leaves
		// a row pointing at a missing file
		if err := app.DBController.DeleteBlog(
			ctx, blogUUID,
		); !handleError(c, err) {
			return
		}
		if err := app.ObjectStore.DeleteFile(
			ctx, blog.Filename,
		); err != nil {
			log.Printf("Error while deleting blog file: %v\n", err.Error())
		}

		c.Status(http.StatusNoContent)
	})

	group.PUT("/about", func(c *gin.Context) {
		ctx := c.Request.Context()

		var body struct {
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Must provide content"},
			)
			return
		}
		if len(body.Content) == 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Content is required"},
			)
			return
		}
		if int64(len(body.Content)) > maxBlogSize {
			c.JSON(
				http.StatusRequestEntityTooLarge,
				gin.H{"error": fmt.Sprintf(
					"Content too large, max size is %d bytes", maxBlogSize,
				)},
			)
			return
		}

		buf := bytes.NewBufferString(body.Content)
		if err := app.ObjectStore.AddFile(
			ctx, provider.AboutFilename, buf, int64(buf.Len()),
		); !handleError(c, err) {
			return
		}
		if err := app.ObjectStore.PublicFile(
			ctx, provider.AboutFilename,
		); !handleError(c, err) {
			return
		}

		c.Status(http.StatusOK)
	})
}
