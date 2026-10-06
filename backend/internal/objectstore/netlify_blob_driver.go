//go:build netlifyblob

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/netlifyblob"

	"context"
	"strconv"
	"time"
)

func StoreEnvKeys() []string {
	return []string{
		"BLOB_SITE_ID",
		"BLOB_TOKEN",
		"BLOB_API_URL",
		"BLOB_STORE",
		"BLOB_TIMEOUT_MS",
		"BLOB_MAX_DB_SNAPSHOTS",
		"BLOB_UPLOAD_SECRET",
	}
}

func GetObjectStore(
	ctx context.Context,
	env map[string]string,
) (ObjectStore, error) {
	ms, _ := strconv.ParseInt(env["BLOB_TIMEOUT_MS"], 10, 64)
	n, _ := strconv.ParseUint(env["BLOB_MAX_DB_SNAPSHOTS"], 10, 32)
	return netlifyblob.New(ctx, netlifyblob.NewConfig(
		env["BLOB_SITE_ID"],
		env["BLOB_TOKEN"],
		env["BLOB_API_URL"],
		env["BLOB_STORE"],
		time.Duration(ms)*time.Millisecond,
		uint(n),
		env["BLOB_UPLOAD_SECRET"],
	))
}
