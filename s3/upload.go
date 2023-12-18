package s3

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (s *S3) Upload(ctx context.Context, file io.ReadSeeker, name, ext string) (UploadResponse, error) {
	mime, err := s.detectMimeType(file)
	if err != nil {
		return UploadResponse{}, err
	}

	file.Seek(0, 0)
	fileId := uuid.NewString()
	output, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &fileId,
		Body:   file,
		Metadata: map[string]string{
			"ext":  ext,
			"name": name,
		},
	})
	if err != nil {
		return UploadResponse{}, err
	}

	return UploadResponse{
		FileID: fileId,
		Info: FileInfo{
			FileExt:      ext,
			FileID:       fileId,
			FileMimetype: mime,
			FileSize:     s.getSize(file),
			ETag:         *output.ETag,
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
