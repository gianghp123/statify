package minio

import (
	"context"
	"io"
	"net/url"

	minioGo "github.com/minio/minio-go/v7"
)

type Interface interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minioGo.PutObjectOptions) (info minioGo.UploadInfo, err error)
	PresignedPostPolicy(ctx context.Context, post *minioGo.PostPolicy) (*url.URL, map[string]string, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minioGo.GetObjectOptions) (*minioGo.Object, error)
	ListObjects(ctx context.Context, bucketName string, opts minioGo.ListObjectsOptions) <-chan minioGo.ObjectInfo
	StatObject(ctx context.Context, bucketName, objectName string, opts minioGo.StatObjectOptions) (minioGo.ObjectInfo, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minioGo.RemoveObjectOptions) error
	RemoveObjectsByPrefix(ctx context.Context, bucketName, prefix string) error
	RemoveObjects(ctx context.Context, bucketName string, objectsCh <-chan minioGo.ObjectInfo, opts minioGo.RemoveObjectsOptions) <-chan minioGo.RemoveObjectError
}
