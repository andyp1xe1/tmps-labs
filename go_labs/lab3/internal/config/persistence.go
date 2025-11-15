package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tmps-go-labs/lab3/internal/character"
)

// PersistentConfig handles saving/loading of app config and characters
type PersistentConfig struct {
	GeminiAPIKey    string           `json:"gemini_api_key"`
	AppName         string           `json:"app_name"`
	Version         string           `json:"version"`
	SavedCharacters []SavedCharacter `json:"saved_characters"`
	LastCharacterID string           `json:"last_character_id,omitempty"`
}

// SavedCharacter represents a serializable character
type SavedCharacter struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Race         string          `json:"race"`
	Class        string          `json:"class"`
	Background   string          `json:"background"`
	Stats        character.Stats `json:"stats"`
	Personality  []string        `json:"personality"`
	Equipment    []string        `json:"equipment"`
	Abilities    []string        `json:"abilities"`
	Backstory    string          `json:"backstory"`
	SystemPrompt string          `json:"system_prompt"`
	CreatedAt    string          `json:"created_at"`
}

// configDirFunc allows overriding the config directory for testing
var configDirFunc = getConfigDirDefault

// getConfigDir returns the application config directory
func getConfigDir() (string, error) {
	return configDirFunc()
}

// getConfigDirDefault returns the default config directory
func getConfigDirDefault() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "character-builder")
	return configDir, nil
}

// getConfigPath returns the full path to the config file
func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads the persistent configuration from disk
func LoadConfig() (*PersistentConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &PersistentConfig{
			AppName:         "Character Builder",
			Version:         "1.0.0",
			SavedCharacters: make([]SavedCharacter, 0),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config PersistentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the persistent configuration to disk
func (pc *PersistentConfig) SaveConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveCharacter adds or updates a character in the saved characters list
func (pc *PersistentConfig) SaveCharacter(char character.Character, id string) error {
	savedChar := SavedCharacter{
		ID:           id,
		Name:         char.GetName(),
		Race:         char.GetRace(),
		Class:        char.GetClass(),
		Background:   char.GetBackground(),
		Stats:        char.GetStats(),
		Personality:  char.GetPersonality(),
		Equipment:    char.GetEquipment(),
		Abilities:    char.GetAbilities(),
		Backstory:    char.GetBackstory(),
		SystemPrompt: char.GetSystemPrompt(),
		CreatedAt:    fmt.Sprintf("%d", getCurrentTimestamp()),
	}

	// Check if character already exists (update)
	for i, existing := range pc.SavedCharacters {
		if existing.ID == id {
			pc.SavedCharacters[i] = savedChar
			return pc.SaveConfig()
		}
	}

	// Add new character
	pc.SavedCharacters = append(pc.SavedCharacters, savedChar)
	pc.LastCharacterID = id
	return pc.SaveConfig()
}

// LoadCharacter loads a character by ID and returns a Character interface
func (pc *PersistentConfig) LoadCharacter(id string) (character.Character, error) {
	for _, savedChar := range pc.SavedCharacters {
		if savedChar.ID == id {
			baseChar := &character.BaseCharacter{
				Name:         savedChar.Name,
				Race:         savedChar.Race,
				Class:        savedChar.Class,
				Background:   savedChar.Background,
				Stats:        savedChar.Stats,
				Personality:  savedChar.Personality,
				Equipment:    savedChar.Equipment,
				Abilities:    savedChar.Abilities,
				Backstory:    savedChar.Backstory,
				SystemPrompt: savedChar.SystemPrompt,
			}

			// Apply appropriate racial decorator based on race
			switch savedChar.Race {
			case character.RaceHuman:
				return character.NewHumanDecorator(baseChar), nil
			case character.RaceElf:
				return character.NewElfDecorator(baseChar), nil
			case character.RaceDwarf:
				return character.NewDwarfDecorator(baseChar), nil
			default:
				return baseChar, nil
			}
		}
	}

	return nil, fmt.Errorf("character with ID %s not found", id)
}

// DeleteCharacter removes a character from the saved list
func (pc *PersistentConfig) DeleteCharacter(id string) error {
	for i, savedChar := range pc.SavedCharacters {
		if savedChar.ID == id {
			// Remove character from slice
			pc.SavedCharacters = append(pc.SavedCharacters[:i], pc.SavedCharacters[i+1:]...)

			// Clear last character ID if this was the last character
			if pc.LastCharacterID == id {
				pc.LastCharacterID = ""
			}

			return pc.SaveConfig()
		}
	}

	return fmt.Errorf("character with ID %s not found", id)
}

// ListCharacters returns a list of all saved characters
func (pc *PersistentConfig) ListCharacters() []SavedCharacter {
	return pc.SavedCharacters
}

// GetLastCharacterID returns the ID of the most recently used character
func (pc *PersistentConfig) GetLastCharacterID() string {
	return pc.LastCharacterID
}

// SetLastCharacterID sets the most recently used character
func (pc *PersistentConfig) SetLastCharacterID(id string) error {
	pc.LastCharacterID = id
	return pc.SaveConfig()
}

// Helper function to get current timestamp
func getCurrentTimestamp() int64 {
	return 1732486800 // Mock timestamp for now - in real app would use time.Now().Unix()
}
