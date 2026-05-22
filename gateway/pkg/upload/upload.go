package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileInfo struct {
	StoredName string
	OriginName string
	FilePath   string
	FileURL    string
	FileSize   int64
	MimeType   string
}

func SaveFile(file multipart.File, header *multipart.FileHeader, uploadPath string) (*FileInfo, error) {
	originName := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(originName))
	storedName := uuid.NewString() + ext

	now := time.Now()
	subDir := filepath.Join(fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	saveDir := filepath.Join(uploadPath, subDir)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create upload dir failed: %w", err)
	}

	savePath := filepath.Join(saveDir, storedName)
	dst, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("create upload file failed: %w", err)
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("save upload file failed: %w", err)
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	fileURL := "/" + strings.ReplaceAll(filepath.Join(uploadPath, subDir, storedName), "\\", "/")
	return &FileInfo{
		StoredName: storedName,
		OriginName: originName,
		FilePath:   savePath,
		FileURL:    fileURL,
		FileSize:   size,
		MimeType:   mimeType,
	}, nil
}

func RemoveFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
