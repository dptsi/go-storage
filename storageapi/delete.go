package storageapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cenkalti/backoff/v4"
)

type DeleteResponse struct {
	FileID  string         `json:"file_id"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message,omitempty"`
	Status  responseStatus `json:"status"`
}

func (u DeleteResponse) IsOk() bool {
	return u.Status == statusOk
}

func (s *StorageApi) Delete(ctx context.Context, fileId string) (DeleteResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/d/files/%s", s.storageApiUrl, fileId)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return DeleteResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	if err := s.setAuthorizationHeader(ctx, req); err != nil {
		return DeleteResponse{}, fmt.Errorf("failed to set authorization header: %w", err)
	}

	resp, err := backoff.RetryWithData[*http.Response](func() (*http.Response, error) {
		return client.Do(req)
	}, s.backoff)
	if err != nil {
		return DeleteResponse{}, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return DeleteResponse{}, fmt.Errorf("failed to read response body: %w", err)
		}
		return DeleteResponse{}, fmt.Errorf("failed to delete file: %s", string(body))
	}

	var deleteResponse DeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteResponse); err != nil {
		return DeleteResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if !deleteResponse.IsOk() {
		return DeleteResponse{}, fmt.Errorf("failed to delete file: %s", deleteResponse.Message)
	}

	return deleteResponse, nil
}
