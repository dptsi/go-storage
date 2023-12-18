package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (s *S3) PublicLink(
	ctx context.Context,
	fileId string,
) (string, error) {
	request, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileId),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = publicLinkExpiration
	})
	if err != nil {
		return "", err
	}
	return request.URL, err
}
