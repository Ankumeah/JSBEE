//go:build s3

package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/gin-gonic/gin"

	"context"
)

// issueUploadURL hands the browser a direct-to-S3 presigned PUT URL
func issueUploadURL(
	ctx context.Context,
	app *a.App,
	filename string,
) (string, error) {
	return app.ObjectStore.PresignedUploadURL(ctx, filename, uploadURLTTL)
}

// publicFileURL is the direct S3 public URL for a published file.
func publicFileURL(app *a.App, filename string) string {
	return app.ObjectStore.PublicBaseURL() + "/" + filename
}

func blob(r *gin.RouterGroup, app *a.App) {}
