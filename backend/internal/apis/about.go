package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/provider"

	"github.com/gin-gonic/gin"

	"io"
	"net/http"
)

// This route deals with reading the about us page,
// writing is done by admins in `admin.go`
func about(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/about")

	// This route returns the raw markdown content of the about us page
	group.GET("/content", func(c *gin.Context) {
		ctx := c.Request.Context()

		content, err := app.ObjectStore.GetFile(ctx, provider.AboutFilename)
		if err != nil {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "About page not found"},
			)
			return
		}
		defer content.Close()

		markdown, err := io.ReadAll(content)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			return
		}

		c.Data(http.StatusOK, "text/markdown; charset=utf-8", markdown)
	})
}
