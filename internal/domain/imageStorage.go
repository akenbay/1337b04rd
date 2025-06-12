package domain

type ImageStorageAPI interface {
	Store(imageData []byte) (string, error)
}
