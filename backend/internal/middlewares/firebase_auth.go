package middlewares

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/Ankumeah/JSBEE/backend/internal/provider"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"context"
	"log"
	"net/http"
	"strings"
	"uuid"
)

// This initalises the connection to firebase
// and must be called before any request can be served
// with FireBaseAuthMiddleware
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

// This middleware check every request for a
// firebase JWT within the "Authorization" header
//
// Abort if:
//   - header dosent have "Bearer " prefix or
//   - token is invalid or
//   - email is not provided or unverified or
//   - name is not provided
//
// This middleware also sets:
//   - user's email in `EmailField`
//   - user's name in `NameField`
//   - user's uuid in `UUIDFeild` (will be issued if user does not exist)
//
// InitFirebase must be called before this middleware can be used
func FireBaseAuthMiddleware(app *a.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")

		// Reject on unexpected prefix

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "No auth token provided"},
			)
			return
		}

		// Verify the JWT

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := app.FireBaseClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid token"},
			)
			return
		}

		// Extracct needed feilds

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
		name, ok := token.Claims["name"].(string)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": "No name claim in token"},
			)
			return
		}

		// Check for uuid and issue if it dosent exist

		var userUUID uuid.UUID
		uuidInToken, ok := token.Claims[provider.SiteName+"-uuid"]
		if !ok {
			userUUID = uuid.New()
			err := app.FireBaseClient.SetCustomUserClaims(
				ctx,
				token.UID,
				map[string]interface{}{
					provider.SiteName + "-uuid": userUUID,
				},
			)
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"error": "Internal server error"},
				)
				log.Printf("Error while issueing custom user claim: %v\n", err.Error())
				return
			}
		} else {
			uuidInTokenString, ok := uuidInToken.(string)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Invalid user uuid in auth token"},
				)
				return
			}

			userUUID, err = uuid.Parse(uuidInTokenString)
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Invalid uuid in auth token"},
				)
				return
			}
		}

		c.Set(UUIDFeild, userUUID.String())
		c.Set(NameField, name)
		c.Set(EmailField, email)
		c.Next()
	}
}
