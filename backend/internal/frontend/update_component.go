package frontend

import (
	"github.com/a-h/templ"
	"golang.org/x/sys/unix"

	"context"
	"os"
)

func updateComponent(
	ctx context.Context,
	savePath string,
	component templ.Component,
) error {
	file, err := os.OpenFile(
		savePath+".temp",
		os.O_WRONLY|os.O_CREATE, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

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

	if err = component.Render(ctx, file); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return os.Rename(savePath+".temp", savePath)
}
