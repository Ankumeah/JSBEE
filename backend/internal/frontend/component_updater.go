package frontend

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/components"

	"context"
	"os"
	"path"
)

// This struct is responsible for updating fronted components
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

// This function is to be called at application startup
// to make sure the static files always exist
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
