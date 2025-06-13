package domain

type ImageStorageAPI interface {
	Store(imageData []byte, bucketName string) (string, error)
}
