package minio

import (
	"bytes"
	"fmt"
	"github.com/minio/minio-go/v7"
	"image-resize/internal/config"
	"io"

	"context"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	db *minio.Client
}

const bucketName = "image"

// New MinioConnection func for opening minio connection.
func New(config config.Minio) (*Storage, error) {
	const op = "storage.minio.New"

	ctx := context.Background()

	minioClient, errInit := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.Password, ""),
		Secure: false,
	})
	if errInit != nil {
		return nil, fmt.Errorf("%s: %w", op, errInit)
	}

	location := "us-east-1"

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if !exists || errBucketExists != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return &Storage{db: minioClient}, errInit
}

func (s *Storage) Upload(name string, b []byte, contentType string) (string, error) {
	const op = "storage.minio.Upload"
	uploadInfo, err := s.db.PutObject(
		context.Background(),
		bucketName,
		name,
		bytes.NewReader(b),
		int64(len(b)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return uploadInfo.VersionID, nil
}

func (s *Storage) Download(name string) ([]byte, string, error) {
	const op = "storage.minio.Download"

	object, err := s.db.GetObject(context.Background(), bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	defer object.Close()

	stat, _ := object.Stat()
	result, err := io.ReadAll(object)
	if err != nil || len(result) == 0 {
		return nil, "", fmt.Errorf("Not found %s: %w", op, err)
	}
	return result, stat.ContentType, nil
}
