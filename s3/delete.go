package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DeleteResponse struct {
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}

func (u DeleteResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *S3) Delete(ctx context.Context, fileId string) (DeleteResponse, error) {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &fileId,
	})
	if err != nil {
		return DeleteResponse{}, fmt.Errorf("failed to delete file, %v", err)
	}

	return DeleteResponse{
		Status:  statusOk,
		Message: fmt.Sprintf("DELETE %s", fileId),
	}, nil
}
