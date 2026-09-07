package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"github.com/gin-gonic/gin"

	"net/http"
	"uuid"
)

func paper(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/paper")

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
