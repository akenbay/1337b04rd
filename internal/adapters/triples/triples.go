package triples

import (
	"1337b04rd/internal/domain"
	"1337b04rd/pkg/logger"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Triples struct {
	port int
}

var _ domain.ImageStorageAPI = (*Triples)(nil)

func NewTriples(port int) *Triples {
	return &Triples{
		port: port,
	}
}

func (t *Triples) Store(imageData []byte, bucketName string) (string, error) {
	logger.Info("Fisrt 10 chars of image data:", "chars", imageData[:10])
	image_key, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	logger.Info("Storing new image:", "iamge key", image_key)

	saveImageURL := "http://triple-s:" + fmt.Sprint(t.port) + "/" + bucketName + "/" + image_key

	urlOfImage := "http://localhost:" + fmt.Sprint(t.port) + "/" + bucketName + "/" + image_key

	saveImageReq, err := http.NewRequest(http.MethodPut, saveImageURL, bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}

	// Send request
	client := &http.Client{}

	imageResp, err := client.Do(saveImageReq)
	if err != nil {
		logger.Error("Error when saving image", "error", err)
		return "", err
	}
	defer imageResp.Body.Close()

	// if imageResp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("image upload failed with status: %s", imageResp.Status)
	// }

	if imageResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(imageResp.Body)
		return "", fmt.Errorf("upload failed (status %d): %s", imageResp.StatusCode, string(body))
	}

	return urlOfImage, nil
}

func (t *Triples) CreateBucket(bucketName string) error {
	createBucketURL := "http://triple-s:" + fmt.Sprint(t.port) + "/" + bucketName

	createBucketReq, err := http.NewRequest(http.MethodPut, createBucketURL, nil)
	if err != nil {
		return err
	}

	// Send request
	client := &http.Client{}
	bucketResp, err := client.Do(createBucketReq)
	if err != nil {
		logger.Error("Error when creating bucket", "error", err)
		return err
	}
	defer bucketResp.Body.Close()

	if bucketResp.StatusCode != http.StatusOK && bucketResp.StatusCode != http.StatusConflict {
		return fmt.Errorf("bucket creation failed with status: %s", bucketResp.Status)
	}

	return nil
}

// GenerateRandomToken creates a secure URL-safe token
func generateRandomToken() (string, error) {
	bytes := make([]byte, 22)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}
	key := base64.URLEncoding.EncodeToString(bytes)
	key = strings.ReplaceAll(key, ".", "_") // Replace unsafe chars
	key = strings.ReplaceAll(key, "/", "_") // Replace unsafe chars
	key = strings.ReplaceAll(key, "=", "")  // Replace unsafe chars
	return key, nil
}
