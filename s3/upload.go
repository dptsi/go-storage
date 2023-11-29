package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cenkalti/backoff/v4"
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

func (s *S3) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (UploadResponse, error) {
	fileExt := filepath.Ext(fileHeader.Filename)

	fileName := strings.TrimSuffix(fileHeader.Filename, fileExt)
	fileName = strings.ReplaceAll(fileName, "/[^a-zA-Z0-9]+/", "_")
	if fileName == "" {
		fileName = fmt.Sprintf("undefined_%d", time.Now().Unix())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	fileExtWithoutDot := strings.TrimPrefix(fileExt, ".")
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
			FileExt:      fileExtWithoutDot,
			FileID:       fileId,
			FileMimetype: mime,
			FileName:     fileName,
			FileSize:     int(fileHeader.Size),
			PublicLink:   result.Location,
			ETag:         *result.ETag,
			Timestamp:    time.Now().Format(time.RFC3339),
		},
		Status: statusOk,
	}, nil
}
