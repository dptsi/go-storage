package gcs

import (
	"fmt"
	"io"
	"net/http"
)

func (s *GCS) detectMimeType(file io.Reader) (string, error) {
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)

	// Copy the headers into the FileHeader buffer
	if _, err := file.Read(fileHeader); err != nil {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	return http.DetectContentType(fileHeader), nil
}

// func (*GCS) getSize(stream io.Reader) int {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(stream)
// 	return buf.Len()
// }
