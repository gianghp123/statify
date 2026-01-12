package minio

import (
	"context"
	"io"

	minioGo "github.com/minio/minio-go/v7"
)

type Interface interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minioGo.PutObjectOptions) (info minioGo.UploadInfo, err error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minioGo.GetObjectOptions) (*minioGo.Object, error)
	StatObject(ctx context.Context, bucketName, objectName string, opts minioGo.StatObjectOptions) (minioGo.ObjectInfo, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minioGo.RemoveObjectOptions) error
	RemoveObjectsByPrefix(ctx context.Context, bucketName, prefix string) error
}
