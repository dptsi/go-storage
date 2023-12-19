package s3_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/dptsi/go-storage/s3"
	"github.com/stretchr/testify/assert"
)

func getS3(ctx context.Context) *s3.S3 {
	s3, err := s3.NewS3(ctx, s3.Config{
		Region:          os.Getenv("S3_REGION"),
		Bucket:          os.Getenv("S3_BUCKET"),
		AccessKeyId:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
	})
	if err != nil {
		panic(err)
	}
	return s3
}

const sampleFileUrl = "https://placehold.co/100x100/jpg"

var preuploadedFileInfo = s3.FileInfo{
	FileID:       "13c91aa0-94f2-4e37-8167-5d6297a99646",
	FileExt:      "",
	FileMimetype: "image/jpeg",
	FileSize:     750,
	ETag:         "\"cf4cb127768fbc6ba1484fa6270d5c54\"",
	Timestamp:    "2023-12-19T03:23:15Z",
}

const preuploadedFileSha256 = "263cef32b161ff9365635244b899dd917b3f93b3d2ee322f0094c2b7f336a824"

func downloadSampleFile(ctx context.Context) (*os.File, error) {
	file, err := os.CreateTemp("", "sample")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(sampleFileUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func cleanupFile(t *testing.T, file *os.File) {
	name := file.Name()
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close file (%s): %v", name, err)
	}
	if err := os.Remove(name); err != nil {
		t.Fatalf("failed to remove file (%s): %v", name, err)
	}
}

func TestUploadFile(t *testing.T) {
	ctx := context.Background()
	s3 := getS3(ctx)

	file, err := downloadSampleFile(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFile(t, file)

	fileName := file.Name()
	fileExt := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)
	info, err := s3.Upload(ctx, file, fileName, fileExt)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(info)
	assert.NotEmpty(t, info.FileID)
	assert.Equal(t, fileExt, info.FileExt)
	assert.NotEmpty(t, info.FileMimetype)
	assert.NotEmpty(t, info.FileSize)
	assert.NotEmpty(t, info.ETag)
	assert.NotEmpty(t, info.Timestamp)
}

func TestUploadFileFromBase64(t *testing.T) {
	ctx := context.Background()
	s3 := getS3(ctx)

	file, err := downloadSampleFile(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFile(t, file)

	fileName := file.Name()
	fileExt := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	info, err := s3.UploadFromBase64(
		ctx,
		base64.StdEncoding.Strict().EncodeToString(fileBytes),
		fileName,
		fileExt,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(info)
	assert.NotEmpty(t, info.FileID)
	assert.Equal(t, fileExt, info.FileExt)
	assert.NotEmpty(t, info.FileMimetype)
	assert.NotEmpty(t, info.FileSize)
	assert.NotEmpty(t, info.ETag)
	assert.NotEmpty(t, info.Timestamp)
}

func TestFileInfo(t *testing.T) {
	ctx := context.Background()
	s3 := getS3(ctx)

	info, err := s3.FileInfo(ctx, preuploadedFileInfo.FileID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(info)
	assert.Equal(t, preuploadedFileInfo, info)
}

func TestDownloadFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sample")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	s3 := getS3(ctx)

	filePath := path.Join(tmpDir, preuploadedFileInfo.FileID)
	file, err := s3.Download(ctx, preuploadedFileInfo.FileID, filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFile(t, file)

	t.Log(file.Name())
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, preuploadedFileSha256, fmt.Sprintf("%x", sha256.Sum256(fileBytes)))
}

func TestDownloadFileAsBase64(t *testing.T) {
	ctx := context.Background()
	s3 := getS3(ctx)

	b64, err := s3.DownloadAsBase64(ctx, preuploadedFileInfo.FileID)
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := base64.StdEncoding.Strict().DecodeString(b64)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, preuploadedFileSha256, fmt.Sprintf("%x", sha256.Sum256(fileBytes)))
}

func TestPublicLink(t *testing.T) {
	ctx := context.Background()

	link, err := getS3(ctx).PublicLink(ctx, preuploadedFileInfo.FileID, s3.DefaultPublicLinkExpiration)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(link)
	assert.NotEmpty(t, link)
	assert.NotEmpty(t, link.ExpiredAt)
}

func TestDeleteFile(t *testing.T) {
	ctx := context.Background()
	s3 := getS3(ctx)

	file, err := downloadSampleFile(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFile(t, file)

	fileName := file.Name()
	fileExt := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)
	info, err := s3.Upload(ctx, file, fileName, fileExt)
	if err != nil {
		t.Fatal(err)
	}

	if err := s3.Delete(ctx, info.FileID); err != nil {
		t.Fatal(err)
	}
}
