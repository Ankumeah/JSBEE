//go:build netlifyblob

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/netlifyblob"

	"context"
	"encoding/json"
	"time"
)

func GetObjectStore(
	ctx context.Context,
	env map[string]string,
) (ObjectStore, error) {
	var raw struct {
		SiteID         string `json:"site_id"`
		Token          string `json:"token"`
		APIURL         string `json:"api_url"`
		Store          string `json:"store"`
		TimeoutMS      int64  `json:"timeout_ms"`
		MaxDBSnapshots uint   `json:"max_db_snapshots"`
		UploadSecret   string `json:"upload_secret"`
	}
	if err := json.Unmarshal([]byte(env["OBJECT_STORE_CONFIG"]), &raw); err != nil {
		return nil, err
	}

	return netlifyblob.New(ctx, netlifyblob.NewConfig(
		raw.SiteID,
		raw.Token,
		raw.APIURL,
		raw.Store,
		time.Duration(raw.TimeoutMS)*time.Millisecond,
		raw.MaxDBSnapshots,
		raw.UploadSecret,
	))
}
