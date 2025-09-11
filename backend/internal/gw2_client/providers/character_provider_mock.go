package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	gw2models "github.com/zoehay/gw2-armory/backend/internal/gw2_client/models"
)

type CharacterProviderMock struct{}

func (characterProvider *CharacterProviderMock) GetAllCharacters(apiKey string) ([]gw2models.GW2Character, error) {
	wd, _ := os.Getwd()
	isTesting := strings.Contains(wd, "test")
	leadingFilepath := ""

	if isTesting {
		leadingFilepath = "../."
	}

	filepath := fmt.Sprintf("%s./test_data/character_test_data.txt", leadingFilepath)
	apiCharacters, err := characterProvider.ReadCharactersFromFile(filepath)

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
