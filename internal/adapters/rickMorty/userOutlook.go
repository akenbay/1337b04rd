package rickMorty

import (
	"1337b04rd/internal/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RickMortyAPI struct{}

var _ domain.UserOutlookAPI = (*RickMortyAPI)(nil)

func newRickMortyAPI() *RickMortyAPI {
	return &RickMortyAPI{}
}

func (r *RickMortyAPI) GenerateAvatarAndName(id int) (*domain.UserOutlook, error) {
	apiReq := "https://rickandmortyapi.com/api/character/" + fmt.Sprint(id)
	response, err := http.Get(apiReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
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
