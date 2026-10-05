//go:build netlifyblob

package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	storeerrors "github.com/Ankumeah/JSBEE/backend/internal/objectstore/errors"
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/netlifyblob"
	"github.com/gin-gonic/gin"

	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const maxProxyUploadBytes = 100 << 20

func blobConfig(app *a.App) (netlifyblob.Config, error) {
	return netlifyblob.ParseConfig(app.Config.ObjectStoreConfig)
}

func blobAPIPrefix(app *a.App) string {
	return "/api/" + app.Config.APIVersion + "/blob"
}

func issueUploadURL(
	ctx context.Context,
	app *a.App,
	filename string,
) (string, error) {
	cfg, err := blobConfig(app)
	if err != nil {
		return "", err
	}
	exp, sig, err := netlifyblob.IssueUploadTicket(cfg, filename, uploadURLTTL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/upload/%s?exp=%d&sig=%s",
		blobAPIPrefix(app), filename, exp, sig), nil
}

func publicFileURL(app *a.App, filename string) string {
	return fmt.Sprintf("%s/file/%s", blobAPIPrefix(app), filename)
}

// This handles routes used for blob uploads
func blob(r *gin.RouterGroup, app *a.App) {
	g := r.Group("/blob")

	g.PUT("/upload/:filename", func(c *gin.Context) {
		ctx := c.Request.Context()
		filename := c.Param("filename")
		cfg, err := blobConfig(app)
		if !handleError(c, "blobConfig", err) {
			return
		}
		exp, err := strconv.ParseInt(c.Query("exp"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket"})
			return
		}
		if err := netlifyblob.VerifyUploadTicket(
			cfg, filename, exp, c.Query("sig"),
		); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid ticket"})
			return
		}

		tmp, err := os.CreateTemp("", "jsbee-upload-*")
		if !handleError(c, "CreateTemp", err) {
			return
		}
		tmpName := tmp.Name()
		defer os.Remove(tmpName)
		defer tmp.Close()

		c.Request.Body = http.MaxBytesReader(
			c.Writer, c.Request.Body, maxProxyUploadBytes,
		)
		size, err := io.Copy(tmp, c.Request.Body)
		if err != nil {
			if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
				c.JSON(http.StatusRequestEntityTooLarge,
					gin.H{"error": "File too large"})
				return
			}
			handleError(c, "read upload", err)
			return
		}
		if _, err := tmp.Seek(0, io.SeekStart); !handleError(c, "Seek", err) {
			return
		}
		if !handleError(c, "AddFile",
			app.ObjectStore.AddFile(ctx, filename, tmp, size)) {
			return
		}

		c.Status(http.StatusOK)
	})

	g.GET("/file/:filename", func(c *gin.Context) {
		ctx := c.Request.Context()
		filename := c.Param("filename")

		body, err := app.ObjectStore.GetFile(ctx, filename)
		if errors.Is(err, storeerrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		if !handleError(c, "GetFile", err) {
			return
		}
		defer body.Close()

		c.DataFromReader(
			http.StatusOK, -1, contentTypeFor(filename), body, nil,
		)
	})
}

func contentTypeFor(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown"):
		return "text/markdown; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}
