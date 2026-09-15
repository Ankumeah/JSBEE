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
	filesBaseURL string,
) error {
	return updateComponent(
		ctx, path.Join(u.savePath, volumeFile),
		components.VolumesPage(volumes, u.firebaseClientConfig, filesBaseURL),
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
	aboutURL string,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, adminFile),
		components.AdminDashboardPage(u.firebaseClientConfig, aboutURL),
	)
}

func (u *ComponentUpdater) UpdateBlog(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, blogFile),
		components.BlogPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateAbout(
	ctx context.Context,
	aboutURL string,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, aboutFile),
		components.AboutPage(u.firebaseClientConfig, aboutURL),
	)
}

func (u *ComponentUpdater) UpdateAuthor(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, authorFile),
		components.AuthorPage(u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateTeam(
	ctx context.Context,
	leaders []database.User,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, teamFile),
		components.TeamPage(leaders, u.firebaseClientConfig),
	)
}

func (u *ComponentUpdater) UpdateContact(
	ctx context.Context,
) error {
	return updateComponent(
		ctx,
		path.Join(u.savePath, contactFile),
		components.ContactPage(u.firebaseClientConfig),
	)
}

// This function updates all frontend files and
// is to be called at application startup
// to make sure the static files always exist
func (u *ComponentUpdater) UpdateAll(
	ctx context.Context,
	volumes []database.Volume,
	leaders []database.User,
	filesBaseURL string,
) error {
	aboutURL := filesBaseURL + "/" + AboutFilename

	for _, f := range []func() error{
		func() error { return u.UpdateIndex(ctx) },
		func() error { return u.UpdateVolumes(ctx, volumes, filesBaseURL) },
		func() error { return u.UpdateNotFound(ctx) },
		func() error { return u.UpdateProfile(ctx) },
		func() error { return u.UpdatePaper(ctx) },
		func() error { return u.UpdateReview(ctx) },
		func() error { return u.UpdateAdmin(ctx, aboutURL) },
		func() error { return u.UpdateBlog(ctx) },
		func() error { return u.UpdateAbout(ctx, aboutURL) },
		func() error { return u.UpdateAuthor(ctx) },
		func() error { return u.UpdateTeam(ctx, leaders) },
		func() error { return u.UpdateContact(ctx) },
	} {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
