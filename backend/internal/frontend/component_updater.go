package frontend

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/components"

	"context"
	"os"
	"path"
)

type ComponentUpdater struct {
	savePath string
}

func GetComponentUpdater(savePath string) (*ComponentUpdater, error) {
	return &ComponentUpdater{savePath}, os.MkdirAll(savePath, 0o755)
}

func (u *ComponentUpdater) UpdateIndex(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, indexFile),
		components.IndexPage(),
	)
}

func (u *ComponentUpdater) UpdateVolumes(
	ctx context.Context,
	volumes []database.Volume,
) error {
	return updateComponent(
		ctx, path.Join(u.savePath, volumeFile),
		components.VolumesPage(volumes),
	)
}

func (u *ComponentUpdater) UpdateAll(
	ctx context.Context,
	volumes []database.Volume,
) error {
	if err := u.UpdateIndex(ctx); err != nil {
		return err
	}
	if err := u.UpdateVolumes(ctx, volumes); err != nil {
		return err
	}

	return nil
}
