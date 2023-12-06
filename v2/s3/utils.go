package s3

import (
	"fmt"
	"net/http"

	"github.com/dptsi/go-storage/contracts"
)

func (s *S3) detectMimeType(file contracts.File) (string, error) {
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
