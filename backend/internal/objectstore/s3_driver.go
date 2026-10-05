//go:build s3

package objectstore

import (
	"github.com/Ankumeah/JSBEE/backend/internal/objectstore/s3"

	"context"
)

func GetObjectStore(
	ctx context.Context,
	rawJSON string,
) (ObjectStore, error) {
	config, err := s3.NewS3StaticConfigFromJSON([]byte(rawJSON))
	if err != nil {
		return nil, err
	}

	return s3.New(ctx, config)
}
