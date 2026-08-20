package middlewares

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"context"
	"net/http"
	"strings"
)

func InitFirebase(ctx context.Context, creds []byte) (*auth.Client, error) {
	opts := []option.ClientOption{
		option.WithAuthCredentialsJSON(option.ServiceAccount, creds),
		option.WithTelemetryDisabled(),
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func FireBaseAuthMiddleware(app *a.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "No auth token provided"},
			)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := app.FireBaseClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid token"},
			)
			return
		}

		email, ok := token.Claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": "No email claim in token"},
			)
			return
		}
		verified, ok := token.Claims["email_verified"].(bool)
		if !ok || !verified {
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": "Email not verified"},
			)
			return
		}

		c.Set(EmailField, email)
		c.Next()
	}
}
