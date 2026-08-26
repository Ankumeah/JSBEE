package frontend

import (
	"github.com/a-h/templ"
	"golang.org/x/sys/unix"

	"context"
	"os"
)

// This function generates the given component at the given
// save path. As we cannot rely on `os.CreateTemp` due
// to running this on a scratch docker image
// we make use of the underlying kernals FSLock
// and then laster call `os.Rename` to ensure atomic behaviour
func updateComponent(
	ctx context.Context,
	savePath string,
	component templ.Component,
) error {
	// Open or create the temp file
	file, err := os.OpenFile(
		savePath+".temp",
		os.O_WRONLY|os.O_CREATE, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// Accquire the FSLock
	if err = unix.Flock(int(file.Fd()), unix.LOCK_EX); err != nil {
		return err
	}
	defer unix.Flock(int(file.Fd()), unix.LOCK_UN)

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// Render the component inot the temp file
	if err = component.Render(ctx, file); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	// Atomic rename to final desination
	return os.Rename(savePath+".temp", savePath)
}
