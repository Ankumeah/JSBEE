//go:build netlifyblob

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/netlifyblob"

	"context"
)

func GetObjectStore(
	ctx context.Context,
	rawJSON string,
) (ObjectStore, error) {
	config, err := netlifyblob.ParseConfig([]byte(rawJSON))
	if err != nil {
		return nil, err
	}

	return netlifyblob.New(ctx, config)
}
