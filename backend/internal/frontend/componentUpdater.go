package frontend

import (
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/components"

	"context"
	"path"
)

type ComponentUpdater struct {
	savePath string
}

func GetComponentUpdater(savePath string) *ComponentUpdater {
	return &ComponentUpdater{
		savePath: savePath,
	}
}

func (u *ComponentUpdater) UpdateIndex(ctx context.Context) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, indexFile),
		components.Index(),
	)
}
