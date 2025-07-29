package domain

import "mime/multipart"

type FileUtils interface {
	ValidateImage(fileHeader *multipart.FileHeader) error
	FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error)
}
