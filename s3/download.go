package s3

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *S3) DownloadAsFile(ctx context.Context, fileId, path string) (*os.File, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(file, output.Body); err != nil {
		return nil, err
	}

	return file, nil
}
