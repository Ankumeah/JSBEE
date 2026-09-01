//go:build init

package assets

import (
	"golang.org/x/sys/unix"

	"os"
	"path"
)

// This function saves all assets needed by the
// frontend and is to be called at application
// startup to make sure all assets always exist
func SaveAssets(baseDir string) error {
	for asset, data := range assets {
		if err := writeStaticFile(
			path.Join(baseDir, asset),
			data,
		); err != nil {
			return err
		}
	}

	return nil
}

// This function writes the given data at the given
// save path. As we cannot rely on `os.CreateTemp` due
// to running this on a scratch docker image
// we make use of the underlying kernals FSLock
// and then later call `os.Rename` to ensure atomic behaviour
func writeStaticFile(savePath string, data []byte) error {
	// Open or create the temp file
	file, err := os.OpenFile(
		savePath+".temp",
		os.O_WRONLY|os.O_CREATE, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// Aqquire the FSLock
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

	// Write the data
	if _, err := file.Write(data); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	// Atomic rename to final destination
	return os.Rename(savePath+".temp", savePath)
}
