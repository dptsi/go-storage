package storageapi

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/oauth2/clientcredentials"
)

type responseStatus string

const statusOk responseStatus = "OK"

type FileInfo struct {
	FileExt      string `json:"file_ext"`
	FileID       string `json:"file_id"`
	FileMimetype string `json:"file_mimetype"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	PublicLink   string `json:"public_link"`
	Tag          string `json:"tag"`
	Timestamp    string `json:"timestamp"`
}

type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// OidcProviderURL is the OpenID Connect Provider's
	// URL. This is a constant specific to each server.
	OidcProviderURL string

	// StorageApiURL is the Storage API's URL.
	StorageApiURL string
}

type StorageApi struct {
	oauth2Config  clientcredentials.Config
	storageApiUrl string
	backoff       backoff.BackOff
}

func NewStorageApi(ctx context.Context, config Config) (*StorageApi, error) {
	tokenUrl, err := getOidcTokenEndpoint(config.OidcProviderURL)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate storage api: %w", err)
	}
	return &StorageApi{
		oauth2Config: clientcredentials.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			TokenURL:     tokenUrl,
		},
		storageApiUrl: config.StorageApiURL,
		backoff:       backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3),
	}, nil
}

type UploadResponse struct {
	FileID  string         `json:"file_id"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}
