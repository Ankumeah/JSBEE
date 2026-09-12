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
	savePath             string
	firebaseClientConfig string
}

func GetComponentUpdater(savePath string, firebaseClientConfig string) (*ComponentUpdater, error) {
	return &ComponentUpdater{
		savePath:             savePath,
		firebaseClientConfig: firebaseClientConfig,
	}, os.MkdirAll(savePath, 0o755)
}

func (u *ComponentUpdater) UpdateIndex(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, indexFile),
		components.IndexPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateVolumes(
	ctx context.Context,
	volumes []database.Volume,
) error {
	return updateComponent(
		ctx, path.Join(u.savePath, volumeFile),
		components.VolumesPage(volumes, u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateNotFound(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, notFoundFile),
		components.NotFoundPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateProfile(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, profileFile),
		components.ProfilePage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdatePaper(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, paperFile),
		components.PaperPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateReview(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, reviewFile),
		components.ReviewDashboardPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateAdmin(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, adminFile),
		components.AdminDashboardPage(u.firebaseClientConfig),
	)
}

// This function updates all frontend files and
// is to be called at application startup
// to make sure the static files always exist
func (u *ComponentUpdater) UpdateAll(
	ctx context.Context,
	volumes []database.Volume,
) error {
	for _, f := range []func() error{
		func() error { return u.UpdateIndex(ctx) },
		func() error { return u.UpdateVolumes(ctx, volumes) },
		func() error { return u.UpdateNotFound(ctx) },
		func() error { return u.UpdateProfile(ctx) },
		func() error { return u.UpdatePaper(ctx) },
		func() error { return u.UpdateReview(ctx) },
		func() error { return u.UpdateAdmin(ctx) },
	} {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
