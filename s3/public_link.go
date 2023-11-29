package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

const publicLinkExpiration = 30 * time.Minute

type PublicLinkResponse struct {
	Url       string         `json:"url"`
	ExpiredAt string         `json:"expired_at"`
	Status    responseStatus `json:"status"`
}

func (u PublicLinkResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *S3) PublicLink(ctx context.Context, fileId string) (PublicLinkResponse, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &fileId,
	})
	url, err := req.Presign(publicLinkExpiration)
	if err != nil {
		return PublicLinkResponse{}, fmt.Errorf("failed to get public link, %v", err)
	}

	return PublicLinkResponse{
		Url:       url,
		ExpiredAt: time.Now().Add(publicLinkExpiration).Format(time.RFC3339),
		Status:    statusOk,
	}, nil
}
