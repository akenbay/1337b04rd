package triples

import (
	"1337b04rd/internal/domain"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type Triples struct {
	port int
}

var _ domain.ImageStorageAPI = (*Triples)(nil)

func NewTriples(bucketName string, port int) *Triples {
	return &Triples{
		port: port,
	}
}

func (t *Triples) Store(imageData []byte, bucketName string) (string, error) {
	image_key, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	apiReq := "https://localhost:" + fmt.Sprint(t.port) + "/" + bucketName + "/" + image_key
	body := map[string]string{
		"content": string(imageData),
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// GenerateRandomToken creates a secure URL-safe token
func generateRandomToken() (string, error) {
	bytes := make([]byte, 22)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
