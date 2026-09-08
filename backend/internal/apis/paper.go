package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"bytes"
	"fmt"
	"io"
	"net/http"
	"uuid"
)

// This route deals with paper relasted features
func paper(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/paper")

	// This route adds a new paper
	group.POST("", middlewares.FireBaseAuthMiddleware(app),
		func(c *gin.Context) {
			ctx := c.Request.Context()

			userUUID, err := uuid.Parse(c.GetString(middlewares.UUIDField))
			if !handleError(c, err) {
				return
			}

			title := c.PostForm("title")
			if title == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
				return
			}

			file, _, err := c.Request.FormFile("file")
			if !handleError(c, err) {
				return
			}
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			_, err = io.Copy(buf, file)
			if !handleError(c, err) {
				return
			}

			size := int64(buf.Len())
			if size <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Empty file"})
				return
			}
			if size > maxPaperSize {
				c.JSON(
					http.StatusRequestEntityTooLarge,
					gin.H{"error": fmt.Sprintf(
						"File too large, max size is %d bytes", maxPaperSize,
					)},
				)
				return
			}

			paperUUID := uuid.New()
			filename := paperUUID.String() + ".pdf"

			if err := app.ObjectStore.AddFile(
				ctx, filename, buf, size,
			); !handleError(c, err) {
				return
			}

			if err := app.DBController.AddPaper(ctx, database.Paper{
				UUID:      paperUUID,
				Title:     title,
				Filename:  filename,
				OwnerUUID: &userUUID,
			}); !handleError(c, err) {
				return
			}

			c.JSON(http.StatusCreated, gin.H{"uuid": paperUUID})
		},
	)

	// This route returns the details of a paper
	group.GET("/:paperUUID", func(c *gin.Context) {
		ctx := c.Request.Context()

		paperUUID, err := uuid.Parse(c.Param("paperUUID"))
		if !handleError(c, err) {
			return
		}

		paper, err := app.DBController.GetPaper(ctx, paperUUID)
		if !handleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, gin.H{"paper": paper})
	})
}
