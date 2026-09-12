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
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
				return
			}
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			size, err := io.Copy(buf, io.LimitReader(file, maxPaperSize+1))
			if !handleError(c, err) {
				return
			}

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

			// Published papers are served from the public bucket by
			// `GET /files/:filename`, so make the file public right away.
			// The filename is an unguessable uuid, paper metadata is public.
			if err := app.ObjectStore.PublicFile(
				ctx, filename,
			); !handleError(c, err) {
				app.ObjectStore.DeleteFile(ctx, filename)
				return
			}

			if err := app.DBController.AddPaper(ctx, database.Paper{
				UUID:      paperUUID,
				Title:     title,
				Filename:  filename,
				OwnerUUID: &userUUID,
			}); !handleError(c, err) {
				app.ObjectStore.DeleteFile(ctx, filename)
				return
			}

			c.JSON(http.StatusCreated, gin.H{"uuid": paperUUID})
		},
	)

	// This route returns the details of a paper
	group.GET("/:paperUUID", func(c *gin.Context) {
		ctx := c.Request.Context()

		paperUUID, err := uuid.Parse(c.Param("paperUUID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
			return
		}

		paper, err := app.DBController.GetPaper(ctx, paperUUID)
		if !handleError(c, err) {
			return
		}

		// Direct public-bucket URL so browsers fetch the PDF
		// straight from object storage, not through the backend
		fileURL := app.ObjectStore.PublicBaseURL() + "/" + paper.Filename

		c.JSON(http.StatusOK, gin.H{"paper": paper, "file_url": fileURL})
	})
}
