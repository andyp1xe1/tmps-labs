package config

import (
	"os"
	"path/filepath"
	"testing"
	"tmps-go-labs/lab3/internal/character"
)

func TestPersistentConfig(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()

	// Override configDirFunc for testing
	originalConfigDirFunc := configDirFunc
	configDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() { configDirFunc = originalConfigDirFunc }()

	// Test loading non-existent config (should create default)
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.AppName != "Character Builder" {
		t.Errorf("Expected app name 'Character Builder', got '%s'", config.AppName)
	}

	// Test saving configuration
	config.GeminiAPIKey = "test-api-key"
	if err := config.SaveConfig(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify config file was created
	configPath := filepath.Join(tempDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Test character persistence
	testChar := &character.BaseCharacter{
		Name:         "Testadin",
		Race:         character.RaceHuman,
		Class:        character.ClassFighter,
		Background:   character.BackgroundCriminal,
		Stats:        character.Stats{Strength: 16, Dexterity: 12, Constitution: 14, Intelligence: 10, Wisdom: 13, Charisma: 8},
		Abilities:    []string{character.AbilitySecondWind},
		Equipment:    []string{character.EquipmentSword},
		Personality:  []string{character.PersonalityCurious},
		Backstory:    "A test character with an interesting past.",
		SystemPrompt: "You are a test character.",
	}

	characterID := "test-char-001"

	// Save character
	if err := config.SaveCharacter(testChar, characterID); err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	// Reload config from disk
	reloadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	// Verify character was persisted
	if len(reloadedConfig.SavedCharacters) != 1 {
		t.Fatalf("Expected 1 saved character, got %d", len(reloadedConfig.SavedCharacters))
	}

	savedChar := reloadedConfig.SavedCharacters[0]
	if savedChar.Name != testChar.Name {
		t.Errorf("Expected character name '%s', got '%s'", testChar.Name, savedChar.Name)
	}
	if savedChar.Backstory != testChar.Backstory {
		t.Errorf("Expected backstory '%s', got '%s'", testChar.Backstory, savedChar.Backstory)
	}
	if savedChar.SystemPrompt != testChar.SystemPrompt {
		t.Errorf("Expected system prompt '%s', got '%s'", testChar.SystemPrompt, savedChar.SystemPrompt)
	}

	// Test character loading
	loadedChar, err := reloadedConfig.LoadCharacter(characterID)
	if err != nil {
		t.Fatalf("Failed to load character: %v", err)
	}

	if loadedChar.GetName() != testChar.Name {
		t.Errorf("Expected loaded character name '%s', got '%s'", testChar.Name, loadedChar.GetName())
	}
	if loadedChar.GetBackstory() != testChar.Backstory {
		t.Errorf("Expected loaded backstory '%s', got '%s'", testChar.Backstory, loadedChar.GetBackstory())
	}
	if loadedChar.GetSystemPrompt() != testChar.SystemPrompt {
		t.Errorf("Expected loaded system prompt '%s', got '%s'", testChar.SystemPrompt, loadedChar.GetSystemPrompt())
	}

	t.Log("✅ Character persistence test passed")
}
