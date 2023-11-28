package storageapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

type UploadResponse struct {
	FileID  string         `json:"file_id"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message"`
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

type GetResponse struct {
	Data   string         `json:"data"`
	Info   FileInfo       `json:"info"`
	Status responseStatus `json:"status"`
}

func (u GetResponse) IsOk() bool {
	return u.Status == statusOk
}

type DeleteResponse struct {
	FileID  string         `json:"file_id"`
	Info    FileInfo       `json:"info"`
	Message string         `json:"message"`
	Status  responseStatus `json:"status"`
}

func (u DeleteResponse) IsOk() bool {
	return u.Status == statusOk
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
	}, nil
}

func (s *StorageApi) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (UploadResponse, error) {
	fileExt := filepath.Ext(fileHeader.Filename)

	fileName := strings.TrimSuffix(fileHeader.Filename, fileExt)
	fileName = strings.ReplaceAll(fileName, "/[^a-zA-Z0-9]+/", "_")
	if fileName == "" {
		fileName = fmt.Sprintf("undefined_%d", time.Now().Unix())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return UploadResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	fileExtWithoutDot := strings.TrimPrefix(fileExt, ".")
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
		FileName:      fileName,
		FileExt:       fileExtWithoutDot,
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

	resp, err := client.Do(req)
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

	return uploadResponse, nil
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

	resp, err := client.Do(req)
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

	return getFileByIdResponse, nil
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

	resp, err := client.Do(req)
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

	return deleteResponse, nil
}

func (s *StorageApi) detectMimeType(file multipart.File) (string, error) {
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)

	// Copy the headers into the FileHeader buffer
	if _, err := file.Read(fileHeader); err != nil {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	// set position back to start.
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	return http.DetectContentType(fileHeader), nil
}

func (s *StorageApi) setAuthorizationHeader(ctx context.Context, req *http.Request) error {
	token, err := s.oauth2Config.TokenSource(ctx).Token()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	req.Header.Set("x-client-id", s.oauth2Config.ClientID)
	req.Header.Set("x-code", token.AccessToken)

	return nil
}
