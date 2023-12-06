package storageapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/dptsi/go-storage/v2/contracts"
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

type UploadBody struct {
	FileName      string `json:"file_name"`
	FileExt       string `json:"file_ext"`
	FileMimetype  string `json:"mime_type"`
	BinaryDataB64 string `json:"binary_data_b64"`
}

func (s *StorageApi) Upload(ctx context.Context, file contracts.File, name, ext string) (UploadResponse, error) {
	mime, err := s.detectMimeType(file)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to detect mime type: %w", err)
	}

	// Convert file to base64 string.
	file.Seek(0, 0)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to read file: %w", err)
	}

	uploadBody := UploadBody{
		FileName:      name,
		FileExt:       ext,
		FileMimetype:  mime,
		BinaryDataB64: base64.StdEncoding.EncodeToString(fileBytes),
	}
	uploadBodyJson, err := json.Marshal(uploadBody)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to marshal upload body: %w", err)
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/d/files", s.storageApiUrl)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(uploadBodyJson))
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	if err := s.setAuthorizationHeader(ctx, req); err != nil {
		return UploadResponse{}, fmt.Errorf("failed to set authorization header: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := backoff.RetryWithData[*http.Response](func() (*http.Response, error) {
		return client.Do(req)
	}, s.backoff)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return UploadResponse{}, fmt.Errorf("failed to read response body: %w", err)
		}
		return UploadResponse{}, fmt.Errorf("failed to upload file: %s", string(body))
	}

	var uploadResponse UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return UploadResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if !uploadResponse.IsOk() {
		return UploadResponse{}, fmt.Errorf("failed to upload file: %s", uploadResponse.Message)
	}

	return uploadResponse, nil
}
