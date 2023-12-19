package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const defaultPublicLinkExpiration = 30 * time.Minute

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

func (s *S3) Upload(ctx context.Context, file io.ReadSeeker, name, ext string) (FileInfo, error) {
	mime, err := s.detectMimeType(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to detect mime type: %w", err)
	}

	file.Seek(0, 0)
	fileId := uuid.NewString()
	output, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &fileId,
		Body:   file,
		Metadata: map[string]string{
			"ext": ext,
		},
		ContentType: aws.String(mime),
	})
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to put object to s3: %w", err)
	}

	return FileInfo{
		FileID:       fileId,
		FileExt:      ext,
		FileMimetype: mime,
		FileSize:     s.getSize(file),
		ETag:         *output.ETag,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *S3) DownloadAsFile(ctx context.Context, fileId, path string) (*os.File, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from s3: %w", err)
	}
	defer output.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file to path %s: %w", path, err)
	}

	if _, err := io.Copy(file, output.Body); err != nil {
		return nil, fmt.Errorf("failed to copy file to path %s: %w", path, err)
	}

	return file, nil
}

func (s *S3) Get(ctx context.Context, fileId string) (FileInfo, error) {
	output, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to head object from s3: %w", err)
	}

	metadata := output.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	return FileInfo{
		FileID:       fileId,
		FileExt:      metadata["ext"],
		FileMimetype: *output.ContentType,
		FileSize:     int(*output.ContentLength),
		ETag:         *output.ETag,
		Timestamp:    output.LastModified.UTC().Format(time.RFC3339),
	}, nil
}

func (s *S3) PublicLink(
	ctx context.Context,
	fileId string,
	publicLinkExpiration time.Duration,
) (PublicLinkResponse, error) {
	if publicLinkExpiration <= 0 {
		publicLinkExpiration = defaultPublicLinkExpiration
	}
	request, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = publicLinkExpiration
	})
	if err != nil {
		return PublicLinkResponse{}, fmt.Errorf("failed to presign object from s3: %w", err)
	}
	return PublicLinkResponse{
		Url:       request.URL,
		ExpiredAt: time.Now().Add(publicLinkExpiration).UTC().Format(time.RFC3339),
	}, nil
}

func (s *S3) Delete(ctx context.Context, fileId string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &fileId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from s3: %w", err)
	}

	return nil
}
