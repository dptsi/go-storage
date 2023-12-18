package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GetResponse struct {
	Data    FileInfo       `json:"data"`
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}

func (u GetResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *S3) Get(ctx context.Context, fileId string) (GetResponse, error) {
	output, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	})
	if err != nil {
		return GetResponse{}, err
	}

	metadata := output.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	return GetResponse{
		Data: FileInfo{
			FileID:       fileId,
			FileExt:      metadata["ext"],
			FileMimetype: *output.ContentType,
			FileSize:     int(*output.ContentLength),
			ETag:         *output.ETag,
			Timestamp:    time.Now().Format(time.RFC3339),
		},
		Status: statusOk,
	}, nil
}
