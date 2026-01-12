package minio

import (
	"context"
	"io"

	minioGo "github.com/minio/minio-go/v7"
)

type Client struct {
	client *minioGo.Client
}

func NewClient(client *minioGo.Client) *Client {
	return &Client{client: client}
}

func (c *Client) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minioGo.PutObjectOptions) (info minioGo.UploadInfo, err error) {
	return c.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (c *Client) GetObject(ctx context.Context, bucketName, objectName string, opts minioGo.GetObjectOptions) (*minioGo.Object, error) {
	return c.client.GetObject(ctx, bucketName, objectName, opts)
}

func (c *Client) StatObject(ctx context.Context, bucketName, objectName string, opts minioGo.StatObjectOptions) (minioGo.ObjectInfo, error) {
	return c.client.StatObject(ctx, bucketName, objectName, opts)
}

func (c *Client) RemoveObject(ctx context.Context, bucketName, objectName string, opts minioGo.RemoveObjectOptions) error {
	return c.client.RemoveObject(ctx, bucketName, objectName, opts)
}

func (c *Client) RemoveObjectsByPrefix(ctx context.Context, bucketName, prefix string) error {
	objectsCh := make(chan minioGo.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range c.client.ListObjects(ctx, bucketName, minioGo.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			if object.Err != nil {
				return
			}
			objectsCh <- object
		}
	}()

	errorCh := c.client.RemoveObjects(ctx, bucketName, objectsCh, minioGo.RemoveObjectsOptions{
		GovernanceBypass: true,
	})

	for err := range errorCh {
		if err.Err != nil {
			return err.Err
		}
	}

	return nil
}
