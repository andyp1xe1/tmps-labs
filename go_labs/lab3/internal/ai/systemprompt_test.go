package ai

import (
	"strings"
	"testing"
	"tmps-go-labs/lab3/internal/character"
)

// TestSystemPromptAutoGeneration tests the auto-generation and caching of system prompts
func TestSystemPromptAutoGeneration(t *testing.T) {
	// Create a test character
	baseChar := &character.BaseCharacter{
		Name:        "TestHero",
		Race:        character.RaceHuman,
		Class:       character.ClassFighter,
		Background:  character.BackgroundCriminal,
		Stats:       character.Stats{Strength: 16, Dexterity: 12, Constitution: 14, Intelligence: 10, Wisdom: 13, Charisma: 8},
		Abilities:   []string{character.AbilitySecondWind, character.AbilityActionSurge},
		Equipment:   []string{character.EquipmentChainMail, character.EquipmentSword},
		Personality: []string{character.PersonalityCurious, "disciplined"},
	}

	testChar := character.NewHumanDecorator(baseChar)
	adapter := NewCharacterAIAdapter()

	// Initially, character should have no system prompt
	if testChar.GetSystemPrompt() != "" {
		t.Error("Character should start with empty system prompt")
	}
	if testChar.GetBackstory() != "" {
		t.Error("Character should start with empty backstory")
	}

	// Generate system prompt (should auto-generate with mock client)
	systemPrompt := adapter.characterToPrompt(testChar)

	// Verify system prompt was generated and cached
	if systemPrompt == "" {
		t.Error("System prompt should be generated")
	}
	if testChar.GetSystemPrompt() != systemPrompt {
		t.Error("System prompt should be cached in character")
	}

	// Verify content includes character details
	if !strings.Contains(systemPrompt, testChar.GetName()) {
		t.Error("System prompt should include character name")
	}
	if !strings.Contains(systemPrompt, testChar.GetRace()) {
		t.Error("System prompt should include character race")
	}
	if !strings.Contains(systemPrompt, testChar.GetClass()) {
		t.Error("System prompt should include character class")
	}

	// Call again - should return cached version
	secondCall := adapter.characterToPrompt(testChar)
	if secondCall != systemPrompt {
		t.Error("Second call should return cached system prompt")
	}

	t.Logf("✅ Generated system prompt (%d chars): %s...", len(systemPrompt), systemPrompt[:100])
}
