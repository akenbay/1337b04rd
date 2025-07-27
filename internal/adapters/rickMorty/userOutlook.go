package rickMorty

import (
	"1337b04rd/internal/domain"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RickMortyAPI struct{}

var _ domain.UserOutlookAPI = (*RickMortyAPI)(nil)

func NewRickMortyAPI() *RickMortyAPI {
	return &RickMortyAPI{}
}

func (r *RickMortyAPI) GenerateAvatarAndName(id int) (*domain.UserOutlook, error) {
	apiReq := "https://rickandmortyapi.com/api/character/" + fmt.Sprint(id)

	// Create a custom transport with TLS config that skips verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{Transport: tr}
	response, err := client.Get(apiReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var userOutlook domain.UserOutlook
	err = json.Unmarshal(body, &userOutlook)
	if err != nil {
		return nil, err
	}

	return &userOutlook, nil
}
