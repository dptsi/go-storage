package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cenkalti/backoff/v4"
	"github.com/dptsi/go-storage/v2/contracts"
	"github.com/google/uuid"
)

type UploadResponse struct {
	FileID  string         `json:"file_id"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}

func (u UploadResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *S3) Upload(ctx context.Context, file contracts.File, name, ext string) (UploadResponse, error) {
	mime, err := s.detectMimeType(file)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to detect mime type: %w", err)
	}

	file.Seek(0, 0)
	fileId := uuid.NewString()
	result, err := backoff.RetryWithData[*s3manager.UploadOutput](func() (*s3manager.UploadOutput, error) {
		return s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: &s.bucket,
			Key:    &fileId,
			Body:   file,
		})
	}, s.backoff)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to upload file: %w", err)
	}

	return UploadResponse{
		FileID: fileId,
		Info: FileInfo{
			FileExt:      ext,
			FileID:       fileId,
			FileMimetype: mime,
			FileSize:     s.getSize(file),
			PublicLink:   result.Location,
			ETag:         *result.ETag,
			Timestamp:    time.Now().Format(time.RFC3339),
		},
		Status: statusOk,
	}, nil
}

func (*S3) getSize(stream io.Reader) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Len()
}
