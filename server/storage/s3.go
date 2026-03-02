package storage

import (
	"bytes"
	"context"
	"fmt"
	"funtabs-server/config"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client   *s3.Client
	bucket   string
	endpoint string
	cdn      string
}

func newS3Storage() Storage {
	cfg := config.Cfg.Storage.S3

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("初始化 S3 配置失败: %v", err))
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			// 兼容非 AWS 的 S3 服务（如阿里云 OSS、MinIO）
			o.UsePathStyle = true
		}
	})

	return &s3Storage{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		cdn:      config.Cfg.Storage.CDN,
	}
}

func (s *s3Storage) Save(file multipart.File, header *multipart.FileHeader, dir string) (string, error) {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".bin"
	}
	key := fmt.Sprintf("%s/%d%s", dir, time.Now().UnixNano(), ext)

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("读取上传文件失败: %w", err)
	}

	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", fmt.Errorf("上传到 S3 失败: %w", err)
	}

	// 优先使用 CDN 前缀，否则拼接 endpoint/bucket/key
	if cdn := s.cdn; cdn != "" {
		return strings.TrimRight(cdn, "/") + "/" + key, nil
	}
	return strings.TrimRight(s.endpoint, "/") + "/" + s.bucket + "/" + key, nil
}

func (s *s3Storage) Delete(urlPath string) error {
	if urlPath == "" {
		return nil
	}
	// 从 URL 中提取 S3 key
	key := extractS3Key(urlPath, s.endpoint, s.bucket, s.cdn)
	if key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// extractS3Key 从完整 URL 中还原 S3 key
func extractS3Key(url, endpoint, bucket, cdn string) string {
	for _, prefix := range []string{
		strings.TrimRight(cdn, "/") + "/",
		strings.TrimRight(endpoint, "/") + "/" + bucket + "/",
	} {
		if prefix != "/" && strings.HasPrefix(url, prefix) {
			return strings.TrimPrefix(url, prefix)
		}
	}
	return ""
}
