package apis

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

func ping(r *gin.RouterGroup) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "PONG")
	})
}
