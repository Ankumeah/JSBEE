package objectstore

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"context"
	"encoding/json"
	"errors"
	"io"
	"uuid"
)

type s3StaticConfig struct {
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Region         string `json:"region"`
	MaxDBSnapshots uint   `json:"max_db_snapshots"`
	Url            string `json:"url"`
}

func NewS3StaticConfigFromJSON(
	configJSON []byte,
) (s3StaticConfig, error) {
	var config s3StaticConfig
	err := json.Unmarshal(configJSON, &config)

	return config, err
}

type s3StaticClient struct {
	client         *minio.Client
	region         string
	maxDBSnapshots uint
}

func GetStaticS3Client(
	ctx context.Context,
	config s3StaticConfig,
) (ObjectStore, error) {
	client, err := minio.New(config.Url, &minio.Options{
		Creds: credentials.NewStaticV4(
			config.AccessKey,
			config.SecretKey,
			"",
		),
		Secure:       true,
		BucketLookup: minio.BucketLookupPath,
		Region:       config.Region,
	})

	return s3StaticClient{
		client:         client,
		region:         config.Region,
		maxDBSnapshots: config.MaxDBSnapshots,
	}, err
}

func (s s3StaticClient) Init(
	ctx context.Context,
) error {
	for _, bucket := range buckets {
		if err := s.client.MakeBucket(ctx, bucket,
			minio.MakeBucketOptions{Region: s.region},
		); err != nil {
			if exists, existsErr := s.client.BucketExists(
				ctx, backupBucket,
			); existsErr != nil {
				return err
			} else if !exists {
				return errors.New("Can't make bucket: " + bucket)
			}
		}
	}

	return nil
}

func (s s3StaticClient) AddFile(
	ctx context.Context,
	filename string,
	content io.Reader,
	size int64,
) error {
	_, err := s.client.PutObject(
		ctx, privateBucket, filename, content, size,
		minio.PutObjectOptions{
			ContentType: "application/pdf",
		},
	)

	return err
}

func (s s3StaticClient) PublicFile(
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

func (s s3StaticClient) StoreDBBackup(
	ctx context.Context,
	baseFilename string,
	backupPath string,
) error {
	if s.maxDBSnapshots < 0 {
		return nil
	}

	var minTime int64
	var minFile string
	var totalBackups uint
	for backup := range s.client.ListObjects(
		ctx, backupBucket, minio.ListObjectsOptions{
			Prefix: backupBaseName,
		},
	) {
		if backup.LastModified.Unix() <= minTime {
			minTime = backup.LastModified.Unix()
			minFile = backup.Key
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
