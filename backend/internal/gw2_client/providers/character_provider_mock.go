package providers

import (
	"encoding/json"
	"fmt"
	"os"

	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type CharacterProviderMock struct{}

func (characterProvider *CharacterProviderMock) GetAllCharacters(apiKey string) ([]gw2models.GW2Character, error) {
	apiCharacters, err := characterProvider.ReadCharactersFromFile(testDataPath("character_test_data.txt"))
	if err != nil {
		return nil, fmt.Errorf("error reading from test data file: %s", err)
	}
	return apiCharacters, nil
}

func (characterProvider *CharacterProviderMock) ReadCharactersFromFile(filepath string) ([]gw2models.GW2Character, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var characters []gw2models.GW2Character
	err = json.Unmarshal(content, &characters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return characters, nil
}
