package s3

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type GetResponse struct {
	Data    string         `json:"data"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}

func (u GetResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *S3) Get(ctx context.Context, fileId string) (GetResponse, error) {
	result, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return GetResponse{}, fmt.Errorf("failed to get file, %v", err)
	}
	defer result.Body.Close()

	content, err := io.ReadAll(result.Body)
	if err != nil {
		return GetResponse{}, fmt.Errorf("failed to read file, %v", err)
	}
	// Convert content to Base64
	encoded := base64.StdEncoding.EncodeToString(content)
	effectiveUri, ok := result.Metadata["effectiveUri"]
	if !ok {
		return GetResponse{}, fmt.Errorf("failed to get effectiveUri")
	}

	return GetResponse{
		Data: encoded,
		Info: FileInfo{
			FileID:       fileId,
			FileExt:      path.Ext(*effectiveUri),
			FileMimetype: *result.ContentType,
			FileSize:     int(*result.ContentLength),
			ETag:         *result.ETag,
			Timestamp:    time.Now().Format(time.RFC3339),
		},
		Status: statusOk,
	}, nil

}
