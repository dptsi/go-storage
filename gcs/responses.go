package gcs

type FileInfo struct {
	FileID    string `json:"file_id"`
	FileSize  int    `json:"file_size"`
	ETag      string `json:"etag"`
	Timestamp string `json:"timestamp"`
}

type PublicLinkResponse struct {
	Url       string `json:"url"`
	ExpiredAt string `json:"expired_at"`
}
