package imageValidator

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func validateImage(fileHeader *multipart.FileHeader) error {
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
