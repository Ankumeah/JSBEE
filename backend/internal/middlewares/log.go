package middlewares

import "github.com/gin-gonic/gin"

// This middleware logs any request
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Acctually log to somewhere
		c.Next()
	}
}
