package fileUtils

import (
	"1337b04rd/internal/domain"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type FileUtils struct{}

var _ domain.FileUtils = (*FileUtils)(nil)

func (f *FileUtils) ValidateImage(fileHeader *multipart.FileHeader) error {
	// 1. Check file size (without reading content)
	if fileHeader.Size > 5<<20 { // 5MB limit
		return fmt.Errorf("file too large (max 5MB)")
	}

	// 2. Check MIME type (reads only first 512 bytes)
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("invalid file content")
	}

	mimeType := http.DetectContentType(buf[:n])
	if !strings.HasPrefix(mimeType, "image/") {
		return fmt.Errorf("only images allowed")
	}

	return nil
}

func (f *FileUtils) FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read entire file (after validation)
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return buf.Bytes(), nil
}
