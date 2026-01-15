package minio

import (
	"context"
	"io"

	minioGo "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minioGo.PutObjectOptions) (info minioGo.UploadInfo, err error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	if args.Get(0) == nil {
		return minioGo.UploadInfo{}, args.Error(1)
	}
	return args.Get(0).(minioGo.UploadInfo), args.Error(1)
}

func (m *Mock) GetObject(ctx context.Context, bucketName, objectName string, opts minioGo.GetObjectOptions) (*minioGo.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Get(0).(*minioGo.Object), args.Error(1)
}

func (m *Mock) ListObjects(ctx context.Context, bucketName string, opts minioGo.ListObjectsOptions) <-chan minioGo.ObjectInfo {
	args := m.Called(ctx, bucketName, opts)
	if ch, ok := args.Get(0).(<-chan minioGo.ObjectInfo); ok {
		return ch
	}
	if ch, ok := args.Get(0).(chan minioGo.ObjectInfo); ok {
		return ch
	}
	return nil
}

func (m *Mock) StatObject(ctx context.Context, bucketName, objectName string, opts minioGo.StatObjectOptions) (minioGo.ObjectInfo, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	if info, ok := args.Get(0).(minioGo.ObjectInfo); ok {
		return info, args.Error(1)
	}
	if info, ok := args.Get(0).(*minioGo.ObjectInfo); ok {
		return *info, args.Error(1)
	}
	return minioGo.ObjectInfo{}, args.Error(1)
}

func (m *Mock) RemoveObject(ctx context.Context, bucketName, objectName string, opts minioGo.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

func (m *Mock) RemoveObjectsByPrefix(ctx context.Context, bucketName, prefix string) error {
	args := m.Called(ctx, bucketName, prefix)
	return args.Error(0)
}

func (m *Mock) RemoveObjects(ctx context.Context, bucketName string, objectsCh <-chan minioGo.ObjectInfo, opts minioGo.RemoveObjectsOptions) <-chan minioGo.RemoveObjectError {
	args := m.Called(ctx, bucketName, objectsCh, opts)
	if ch, ok := args.Get(0).(<-chan minioGo.RemoveObjectError); ok {
		return ch
	}
	if ch, ok := args.Get(0).(chan minioGo.RemoveObjectError); ok {
		return ch
	}
	return nil
}
