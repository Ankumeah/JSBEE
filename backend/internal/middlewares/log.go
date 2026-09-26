package middlewares

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	"github.com/gin-gonic/gin"

	"context"
	"log/slog"
)

// This middleware logs any errors that
// may have appeared during the execution
// of a route
// This should be the top most middleware
// to make sure it catches all errors
// This takes its own context rather then
// the request context to make sure that
// logs are written even after the client disconnects
func LogMiddleware(app *a.App, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, err := range c.Errors {
			app.Logger.LogAttrs(
				ctx, slog.LevelError,
				err.Error(),
				slog.Any("error", err),
			)
		}

		c.Next()
	}
}
