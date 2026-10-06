//go:build s3

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/s3"

	"context"
	"strconv"
)

func StoreEnvKeys() []string {
	return []string{
		"S3_ACCESS_KEY",
		"S3_SECRET_KEY",
		"S3_REGION",
		"S3_MAX_DB_SNAPSHOTS",
		"S3_URL",
		"S3_SECURE",
	}
}

func GetObjectStore(
	ctx context.Context,
	env map[string]string,
) (ObjectStore, error) {
	n, _ := strconv.ParseUint(env["S3_MAX_DB_SNAPSHOTS"], 10, 32)
	secure, _ := strconv.ParseBool(env["S3_SECURE"])
	return s3.New(ctx, s3.NewConfig(
		env["S3_ACCESS_KEY"],
		env["S3_SECRET_KEY"],
		env["S3_REGION"],
		uint(n),
		env["S3_URL"],
		secure,
	))
}
