package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type responseStatus string

const (
	statusOk  responseStatus = "OK"
	statusErr responseStatus = "ERR"
)

type FileInfo struct {
	FileExt      string `json:"file_ext"`
	FileID       string `json:"file_id"`
	FileMimetype string `json:"file_mimetype"`
	FileSize     int    `json:"file_size"`
	ETag         string `json:"etag"`
	Timestamp    string `json:"timestamp"`
}

type Config struct {
	Region          string
	Bucket          string
	AccessKeyId     string
	SecretAccessKey string
}

type S3 struct {
	bucket        string
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewS3(ctx context.Context, cfg Config) (*S3, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyId,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg)

	return &S3{
		bucket:        cfg.Bucket,
		client:        client,
		presignClient: s3.NewPresignClient(client),
	}, nil
}
