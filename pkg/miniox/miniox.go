// Package miniox 封装 MinIO(S3) 对象存储客户端。
package miniox

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Bucket          string
	PublicBaseURL   string
}

type Client struct {
	cli    *minio.Client
	bucket string
	base   string
}

// New 创建客户端并确保 bucket 存在（设为公共读，便于直接访问）。
func New(c Config) (*Client, error) {
	m, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := m.BucketExists(ctx, c.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := m.MakeBucket(ctx, c.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		// 允许匿名下载对象
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + c.Bucket + `/*"]}]}`
		_ = m.SetBucketPolicy(ctx, c.Bucket, policy)
	}

	return &Client{cli: m, bucket: c.Bucket, base: c.PublicBaseURL}, nil
}

// Put 上传对象，返回可访问 URL。
func (c *Client) Put(ctx context.Context, objectName, contentType string, r io.Reader, size int64) (string, error) {
	_, err := c.cli.PutObject(ctx, c.bucket, objectName, r, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return c.base + "/" + c.bucket + "/" + objectName, nil
}
