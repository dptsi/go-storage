package its

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cenkalti/backoff/v4"
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

func (s *StorageApi) Get(ctx context.Context, fileId string) (GetResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/d/files/%s", s.storageApiUrl, fileId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return GetResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	if err := s.setAuthorizationHeader(ctx, req); err != nil {
		return GetResponse{}, fmt.Errorf("failed to set authorization header: %w", err)
	}

	resp, err := backoff.RetryWithData[*http.Response](func() (*http.Response, error) {
		return client.Do(req)
	}, s.backoff)
	if err != nil {
		return GetResponse{}, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return GetResponse{}, fmt.Errorf("failed to read response body: %w", err)
		}
		return GetResponse{}, fmt.Errorf("failed to get file by id: %s", string(body))
	}

	var getFileByIdResponse GetResponse
	if err := json.NewDecoder(resp.Body).Decode(&getFileByIdResponse); err != nil {
		return GetResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if !getFileByIdResponse.IsOk() {
		return GetResponse{}, fmt.Errorf("failed to get file by id: %s", getFileByIdResponse.Message)
	}

	return getFileByIdResponse, nil
}
