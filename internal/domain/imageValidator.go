package domain

import "mime/multipart"

type ImageValidator interface {
	Validate(fileHeader *multipart.FileHeader) error
}
