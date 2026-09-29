package objectstore

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrNoSuchUpload = errors.New("No such upload")

type ObjectStore interface {
	// This is used for any init logic needed by the
	// object store sych as creating buckets, etc
	// This should be called at project init
	Init(ctx context.Context) error

	// This adds a file to storage, by default this adds the
	// file to private bucket if the driver supports it
	AddFile(ctx context.Context, filename string, content io.Reader, size int64) error

	// Moves a file from the private bucket to the
	// public bucket if the driver supports it
	PublicFile(ctx context.Context, filename string) error

	// Gets a file from the public bucket
	GetFile(ctx context.Context, filename string) (io.ReadCloser, error)

	// Returns the public base URL under which files in the
	// public bucket can be fetched directly, without going
	// through the backend (e.g. `https://store.example/public`)
	PublicBaseURL() string

	// Deletes a file from the public and private buckets
	// if the driver supports it
	DeleteFile(ctx context.Context, filename string) error

	// Returns a presigned URL the client can PUT the file to
	// directly. The upload lands in the private bucket
	PresignedUploadURL(ctx context.Context, filename string, expiry time.Duration) (string, error)

	// Returns the size of a file in the private bucket.
	// Returns ErrNoSuchUpload when the file does not exist
	PrivateFileSize(ctx context.Context, filename string) (int64, error)

	// Store a db snapshot and remove older snapshots
	StoreDBBackup(ctx context.Context, backupPath string) error
}
