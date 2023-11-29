package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cenkalti/backoff/v4"
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
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	PublicLink   string `json:"public_link"`
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
	bucket   string
	session  *session.Session
	backoff  backoff.BackOff
	uploader *s3manager.Uploader
}

func NewS3(ctx context.Context, config Config) (*S3, error) {
	s3Config := &aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyId,
			config.SecretAccessKey,
			"",
		),
	}
	sess, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}
	return &S3{
		bucket:   config.Bucket,
		session:  sess,
		backoff:  backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3),
		uploader: s3manager.NewUploader(sess),
	}, nil
}
