package providers

import (
	"encoding/json"
	"fmt"
	"io"

	gw2client "github.com/zoehay/gw2-armory/backend/internal/gw2_client"
	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type CharacterDataProvider interface {
	GetAllCharacters(apiKey string) ([]gw2models.GW2Character, error)
}

type CharacterProvider struct{}

func (characterProvider *CharacterProvider) GetAllCharacters(apiKey string) ([]gw2models.GW2Character, error) {
	res, err := gw2client.GetAllCharacters(apiKey)

	if err != nil {
		return nil, fmt.Errorf("provider get error: %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("provider io.ReadAll error: %s", err)
	}

	var result []gw2models.GW2Character

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("provider json.Unmarshal error: %s", err)
	}

	return result, nil
}
