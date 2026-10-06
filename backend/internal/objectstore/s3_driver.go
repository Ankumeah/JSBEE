//go:build s3

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/s3"

	"context"
	"encoding/json"
)

func GetObjectStore(
	ctx context.Context,
	env map[string]string,
) (ObjectStore, error) {
	var raw struct {
		AccessKey      string `json:"access_key"`
		SecretKey      string `json:"secret_key"`
		Region         string `json:"region"`
		MaxDBSnapshots uint   `json:"max_db_snapshots"`
		Url            string `json:"url"`
		Secure         bool   `json:"secure"`
	}
	if err := json.Unmarshal([]byte(env["OBJECT_STORE_CONFIG"]), &raw); err != nil {
		return nil, err
	}

	return s3.New(ctx, s3.NewConfig(
		raw.AccessKey,
		raw.SecretKey,
		raw.Region,
		raw.MaxDBSnapshots,
		raw.Url,
		raw.Secure,
	))
}
