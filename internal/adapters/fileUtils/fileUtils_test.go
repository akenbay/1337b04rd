package fileUtils

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

// mockFileHeader creates a mock multipart.FileHeader for testing
func mockFileHeader(filename string, content []byte, size int64) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(content)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		panic(err)
	}

	// Override size if specified
	if size > 0 {
		fileHeader.Size = size
	}

	return fileHeader
}

func TestFileUtils_ValidateImage(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func() *multipart.FileHeader
		expectedErr error
	}{
		{
			name: "valid small JPEG",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("test.jpg", []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF"), 1<<20) // 1MB
			},
			expectedErr: nil,
		},
		{
			name: "valid PNG",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("test.png", []byte("\x89PNG\x0D\x0A\x1A\x0A"), 2<<20) // 2MB
			},
			expectedErr: nil,
		},
		{
			name: "file too large",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("large.jpg", []byte("\xFF\xD8\xFF"), 6<<20) // 6MB
			},
			expectedErr: errors.New("file too large (max 5MB)"),
		},
		{
			name: "invalid file type (text)",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("text.txt", []byte("just some text"), 1<<10) // 1KB
			},
			expectedErr: errors.New("only images allowed"),
		},
		{
			name: "empty file",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("empty.jpg", []byte{}, 0)
			},
			expectedErr: errors.New("invalid file content"),
		},
		{
			name: "corrupt image header",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("corrupt.jpg", []byte("not an image"), 1<<10) // 1KB
			},
			expectedErr: errors.New("only images allowed"),
		},
	}

	utils := NewFileUtils()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHeader := tt.setupFile()
			err := utils.ValidateImage(fileHeader)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.expectedErr.Error())
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %q, got %q", tt.expectedErr.Error(), err.Error())
				}
			}
		})
	}
}

func TestFileUtils_FileHeaderToBytes(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func() *multipart.FileHeader
		expected    []byte
		expectedErr error
	}{
		{
			name: "successful conversion",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("test.txt", []byte("file content"), 0)
			},
			expected:    []byte("file content"),
			expectedErr: nil,
		},
		{
			name: "empty file",
			setupFile: func() *multipart.FileHeader {
				return mockFileHeader("empty.txt", []byte{}, 0)
			},
			expected:    []byte{},
			expectedErr: nil,
		},
		{
			name: "large file",
			setupFile: func() *multipart.FileHeader {
				content := bytes.Repeat([]byte("a"), 2<<20) // 2MB
				return mockFileHeader("large.bin", content, 0)
			},
			expected:    bytes.Repeat([]byte("a"), 2<<20),
			expectedErr: nil,
		},
		// Note: It's difficult to simulate file open errors with the current implementation
		// since we're using httptest to create the file header. In a real scenario,
		// you might want to create a mock FileHeader that fails on Open().
	}

	utils := NewFileUtils()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHeader := tt.setupFile()
			result, err := utils.FileHeaderToBytes(fileHeader)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if !bytes.Equal(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.expectedErr.Error())
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %q, got %q", tt.expectedErr.Error(), err.Error())
				}
			}
		})
	}
}

// Additional test for error cases that are hard to simulate with mockFileHeader
func TestFileUtils_FileHeaderToBytes_ErrorCases(t *testing.T) {
	utils := NewFileUtils()

	t.Run("nil file header", func(t *testing.T) {
		_, err := utils.FileHeaderToBytes(nil)
		if err == nil {
			t.Error("expected error for nil file header")
		}
	})

	// To test file open errors, we'd need to create a custom FileHeader implementation
	// that fails when Open() is called. This would require more extensive mocking.
}
