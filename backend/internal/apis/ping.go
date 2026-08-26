package apis

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

// This is a simple health check route
func ping(r *gin.RouterGroup) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "PONG")
	})
}
