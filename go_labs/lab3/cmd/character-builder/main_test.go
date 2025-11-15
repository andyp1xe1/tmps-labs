package main

import (
	"testing"
	"tmps-go-labs/lab3/internal/character"
	"tmps-go-labs/lab3/internal/config"
	"tmps-go-labs/lab3/internal/facade"
)

func TestCharacterCreation(t *testing.T) {
	// Test Singleton pattern
	config1 := config.GetGameConfig()
	config2 := config.GetGameConfig()

	if config1 != config2 {
		t.Error("Singleton pattern failed: got different instances")
	}

	// Test Factory + Builder + Decorator patterns via Facade
	facade := facade.NewCharacterCreationFacade()

	char, err := facade.CreateBasicCharacter("TestHero", character.RaceHuman, character.ClassFighter, character.BackgroundAcolyte)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Verify character properties
	if char.GetName() != "TestHero" {
		t.Errorf("Expected name 'TestHero', got %s", char.GetName())
	}

	if char.GetRace() != character.RaceHuman {
		t.Errorf("Expected race '%s', got %s", character.RaceHuman, char.GetRace())
	}

	if char.GetClass() != character.ClassFighter {
		t.Errorf("Expected class '%s', got %s", character.ClassFighter, char.GetClass())
	}

	// Test that decorators are applied by checking abilities
	abilities := char.GetAbilities()
	foundHumanAbility := false
	for _, ability := range abilities {
		if ability == character.AbilityHumanDetermination {
			foundHumanAbility = true
			break
		}
	}
	if !foundHumanAbility {
		t.Error("Human racial abilities not applied correctly")
	}

	// Test abilities are present
	if len(abilities) == 0 {
		t.Error("Character should have abilities")
	}
}

func TestFactoryValidation(t *testing.T) {
	factory := character.NewCharacterFactory()

	// Test valid combination
	char, err := factory.CreateCharacter(character.RaceHuman, character.ClassFighter)
	if err != nil {
		t.Errorf("Valid combination should work: %v", err)
	}
	if char == nil {
		t.Error("Character should not be nil for valid combination")
	}

	// Test invalid combination (Dwarf Wizard is not in our valid combinations)
	_, err = factory.CreateCharacter(character.RaceDwarf, character.ClassWizard)
	if err == nil {
		t.Error("Invalid combination should return error")
	}
}

func TestBuilder(t *testing.T) {
	factory := character.NewCharacterFactory()
	baseChar, err := factory.CreateCharacter(character.RaceElf, character.ClassWizard)
	if err != nil {
		t.Fatalf("Failed to create base character: %v", err)
	}

	builder := character.NewCharacterBuilder(baseChar)

	// Test builder validation
	_, err = builder.SetName("").Build()
	if err == nil {
		t.Error("Builder should reject empty name")
	}

	// Reset builder and test successful build
	builder = character.NewCharacterBuilder(baseChar)
	char, err := builder.
		SetName("Elaria").
		SetBackground(character.BackgroundScholar).
		SetStats(10, 14, 12, 16, 13, 11).
		Build()

	if err != nil {
		t.Errorf("Valid build should succeed: %v", err)
	}

	if char != nil && char.GetName() != "Elaria" {
		t.Errorf("Expected name 'Elaria', got %s", char.GetName())
	}
}

func TestDecorators(t *testing.T) {
	// Create base character
	baseChar := &character.BaseCharacter{
		Name:  "TestChar",
		Race:  character.RaceHuman,
		Class: character.ClassFighter,
		Stats: character.Stats{Strength: 10, Dexterity: 10, Constitution: 10, Intelligence: 10, Wisdom: 10, Charisma: 10},
	}

	// Apply decorators
	decoratedChar := character.NewHumanDecorator(baseChar)

	// Test that bonuses are applied
	stats := decoratedChar.GetStats()
	if stats.Strength != 11 { // Base 10 + Human bonus 1
		t.Errorf("Expected Strength 11, got %d", stats.Strength)
	}

	// Test abilities are added
	abilities := decoratedChar.GetAbilities()
	found := false
	for _, ability := range abilities {
		if ability == character.AbilityHumanDetermination {
			found = true
			break
		}
	}
	if !found {
		t.Error("Human racial ability not found")
	}
}
