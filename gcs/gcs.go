package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type Config struct {
	Bucket string
}

type GCS struct {
	bucket string
	client *storage.Client
}

func NewGCS(ctx context.Context, cfg Config) (*GCS, error) {
	s, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCS{
		bucket: cfg.Bucket,
		client: s,
	}, nil
}

func (s *GCS) Upload(ctx context.Context, file io.Reader) (FileInfo, error) {
	fileId := uuid.NewString()
	object := s.client.Bucket(s.bucket).Object(fileId)
	w := object.NewWriter(ctx)

	if _, err := io.Copy(w, file); err != nil {
		return FileInfo{}, fmt.Errorf("failed to put object to GCS: %w", err)
	}
	if err := w.Close(); err != nil {
		return FileInfo{}, fmt.Errorf("failed to close writer: %w", err)
	}
	attrs := w.Attrs()

	return FileInfo{
		FileID:    fileId,
		FileSize:  int(attrs.Size),
		ETag:      attrs.Etag,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *GCS) Delete(ctx context.Context, fileId string) error {
	return s.client.Bucket(s.bucket).Object(fileId).Delete(ctx)
}

func (s *GCS) Stream(ctx context.Context, fileId string) (io.ReadCloser, error) {
	return s.client.Bucket(s.bucket).Object(fileId).NewReader(ctx)
}
