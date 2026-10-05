package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	storeerrors "github.com/Ankumeah/JSBEE/backend/internal/objectstore/errors"

	"context"
	"errors"
	"io"
	"net/http"
	"time"
	"uuid"
)

type Client struct {
	client         *minio.Client
	region         string
	maxDBSnapshots uint
	publicBaseURL  string
}

func New(
	ctx context.Context,
	config s3StaticConfig,
) (*Client, error) {
	client, err := minio.New(config.Url, &minio.Options{
		Creds: credentials.NewStaticV4(
			config.AccessKey,
			config.SecretKey,
			"",
		),
		Secure:       config.Secure,
		BucketLookup: minio.BucketLookupPath,
		Region:       config.Region,
	})

	return &Client{
		client:         client,
		region:         config.Region,
		maxDBSnapshots: config.MaxDBSnapshots,
		publicBaseURL:  publicBaseURLFromConfig(config),
	}, err
}

func (s *Client) Init(
	ctx context.Context,
) error {
	for _, bucket := range buckets {
		if err := s.client.MakeBucket(ctx, bucket,
			minio.MakeBucketOptions{Region: s.region},
		); err != nil {
			if exists, existsErr := s.client.BucketExists(
				ctx, bucket,
			); existsErr != nil {
				return err
			} else if !exists {
				return errors.New("Can't make bucket: " + bucket)
			}
		}
	}

	return nil
}

func (s *Client) AddFile(
	ctx context.Context,
	filename string,
	content io.Reader,
	size int64,
) error {
	_, err := s.client.PutObject(
		ctx, privateBucket, filename, content, size,
		minio.PutObjectOptions{
			ContentType: contentTypeFor(filename),
		},
	)

	return err
}

func (s *Client) PublicFile(
	ctx context.Context,
	filename string,
) error {
	src := minio.CopySrcOptions{
		Bucket: privateBucket,
		Object: filename,
	}
	dest := minio.CopyDestOptions{
		Bucket: publicBucket,
		Object: filename,
	}

	if _, err := s.client.CopyObject(
		ctx, dest, src,
	); err != nil {
		return err
	}

	return s.client.RemoveObject(
		ctx, privateBucket, filename,
		minio.RemoveObjectOptions{},
	)
}

func (s *Client) PublicBaseURL() string {
	return s.publicBaseURL
}

func (s *Client) GetFile(
	ctx context.Context,
	filename string,
) (io.ReadCloser, error) {
	object, err := s.client.GetObject(
		ctx, publicBucket, filename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	if _, err := object.Stat(); err != nil {
		return nil, err
	}

	return object, nil
}

func (s *Client) DeleteFile(
	ctx context.Context,
	filename string,
) error {
	if err := s.client.RemoveObject(
		ctx, publicBucket, filename,
		minio.RemoveObjectOptions{},
	); err != nil {
		return err
	}

	return s.client.RemoveObject(
		ctx, privateBucket, filename,
		minio.RemoveObjectOptions{},
	)
}

func (s *Client) PresignedUploadURL(
	ctx context.Context,
	filename string,
	expiry time.Duration,
) (string, error) {
	url, err := s.client.PresignedPutObject(
		ctx, privateBucket, filename, expiry,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *Client) PrivateFileSize(
	ctx context.Context,
	filename string,
) (int64, error) {
	info, err := s.client.StatObject(
		ctx, privateBucket, filename,
		minio.StatObjectOptions{},
	)
	if err != nil {
		var errResp minio.ErrorResponse
		if errors.As(err, &errResp) &&
			errResp.StatusCode == http.StatusNotFound {
			return 0, storeerrors.ErrNotFound
		}
		return 0, err
	}

	return info.Size, nil
}

func (s *Client) StoreDBBackup(
	ctx context.Context,
	backupPath string,
) error {
	if s.maxDBSnapshots < 1 {
		return nil
	}

	var minTime int64
	var minFile string
	var totalBackups uint
	first := true
	for backup := range s.client.ListObjects(
		ctx, backupBucket, minio.ListObjectsOptions{
			Prefix: backupBaseName,
		},
	) {
		if backup.Err != nil {
			continue
		}
		if first || backup.LastModified.Unix() < minTime {
			minTime = backup.LastModified.Unix()
			minFile = backup.Key
			first = false
		}
		totalBackups += 1
	}

	savePath := backupBaseName + uuid.New().String()
	if totalBackups >= s.maxDBSnapshots {
		savePath = minFile
	}

	_, err := s.client.FPutObject(
		ctx, backupBucket, savePath, backupPath,
		minio.PutObjectOptions{},
	)

	return err
}
