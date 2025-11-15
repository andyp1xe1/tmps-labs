package facade

import (
	"fmt"
	"tmps-go-labs/lab3/internal/character"
	"tmps-go-labs/lab3/internal/config"
)

// CharacterCreationFacade provides a simplified interface for character creation
type CharacterCreationFacade struct {
	factory character.CharacterFactory
	config  *config.GameConfig
}

// NewCharacterCreationFacade creates a new facade instance
func NewCharacterCreationFacade() *CharacterCreationFacade {
	return &CharacterCreationFacade{
		factory: character.NewCharacterFactory(),
		config:  config.GetGameConfig(),
	}
}

// CreateBasicCharacter creates a character with standard array stats (15,14,13,12,10,8)
func (f *CharacterCreationFacade) CreateBasicCharacter(name, race, class, background string) (character.Character, error) {
	// Step 1: Use factory to create base character
	baseChar, err := f.factory.CreateCharacter(race, class)
	if err != nil {
		return nil, fmt.Errorf("failed to create base character: %w", err)
	}

	// Step 2: Use builder to add details with standard array stats
	builder := character.NewCharacterBuilder(baseChar)
	finalChar, err := builder.
		SetName(name).
		SetBackground(background).
		SetStats(15, 14, 13, 12, 10, 8). // Standard D&D point buy equivalent
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build character: %w", err)
	}

	// Step 3: Apply racial decorators
	decoratedChar := f.applyRacialDecorators(finalChar, race)

	return decoratedChar, nil
}

// CreateAdvancedCharacter creates a character with custom stats and equipment
func (f *CharacterCreationFacade) CreateAdvancedCharacter(name, race, class, background string, stats character.Stats, personality []string, equipment []string) (character.Character, error) {
	// Step 1: Use factory to create base character
	baseChar, err := f.factory.CreateCharacter(race, class)
	if err != nil {
		return nil, fmt.Errorf("failed to create base character: %w", err)
	}

	// Step 2: Use builder with custom parameters
	builder := character.NewCharacterBuilder(baseChar)
	finalChar, err := builder.
		SetName(name).
		SetStats(stats.Strength, stats.Dexterity, stats.Constitution,
			stats.Intelligence, stats.Wisdom, stats.Charisma).
		SetBackground(background).
		SetPersonality(personality).
		AddEquipment(equipment).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build character: %w", err)
	}

	// Step 3: Apply decorators
	decoratedChar := f.applyRacialDecorators(finalChar, race)

	return decoratedChar, nil
}

// CreateExperiencedCharacter creates a character with level-based bonuses
func (f *CharacterCreationFacade) CreateExperiencedCharacter(name, race, class, background string, level int) (character.Character, error) {
	// Create basic character first
	basicChar, err := f.CreateBasicCharacter(name, race, class, background)
	if err != nil {
		return nil, err
	}

	// Apply experience decorators
	levelBonuses := f.calculateLevelBonuses(level)
	levelAbilities := f.getLevelAbilities(class, level)

	experiencedChar := character.NewExperienceDecorator(basicChar, level, levelAbilities, levelBonuses)

	// Add some equipment based on level
	if level >= 3 {
		equipment := []string{character.EquipmentMagicWeapon, character.EquipmentMagicArmor}
		equipBonuses := character.Stats{Strength: 1, Dexterity: 1}
		experiencedChar = character.NewEquipmentDecorator(experiencedChar, equipment, equipBonuses, "Equipped with magical gear")
	}

	return experiencedChar, nil
}

// GetAvailableOptions returns all available creation options
func (f *CharacterCreationFacade) GetAvailableOptions() (races []string, classes []string, backgrounds []string) {
	races = f.factory.GetAvailableRaces()
	classes = f.factory.GetAvailableClasses()

	// Get backgrounds from factory (we need to cast to access this method)
	if concreteFactory, ok := f.factory.(*character.ConcreteCharacterFactory); ok {
		backgrounds = []string{character.BackgroundAcolyte, character.BackgroundCriminal, character.BackgroundScholar} // Using constants
		_ = concreteFactory                                                                                            // Suppress unused variable warning
	}

	return races, classes, backgrounds
}

// ValidateCharacterCombination checks if a race-class combination is valid
func (f *CharacterCreationFacade) ValidateCharacterCombination(race, class string) error {
	if !f.factory.IsValidCombination(race, class) {
		return fmt.Errorf("invalid combination: %s cannot be a %s", race, class)
	}
	return nil
}

// GenerateRandomCharacter creates a completely random character with random stats
func (f *CharacterCreationFacade) GenerateRandomCharacter(name string) (character.Character, error) {
	races, classes, backgrounds := f.GetAvailableOptions()

	// Pick random options (simplified random selection)
	race := races[len(name)%len(races)] // Use name length as pseudo-random
	class := classes[len(name)%len(classes)]
	background := backgrounds[len(name)%len(backgrounds)]

	// Validate combination, try alternatives if invalid
	if !f.factory.IsValidCombination(race, class) {
		// Try first valid combination for this race
		for _, validClass := range classes {
			if f.factory.IsValidCombination(race, validClass) {
				class = validClass
				break
			}
		}
	}

	// Step 1: Use factory to create base character
	baseChar, err := f.factory.CreateCharacter(race, class)
	if err != nil {
		return nil, fmt.Errorf("failed to create base character: %w", err)
	}

	// Step 2: Use builder with random stats
	builder := character.NewCharacterBuilder(baseChar)
	finalChar, err := builder.
		SetName(name).
		SetBackground(background).
		RandomizeStats(). // Use random stats for this method
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build character: %w", err)
	}

	// Step 3: Apply racial decorators
	decoratedChar := f.applyRacialDecorators(finalChar, race)

	return decoratedChar, nil
}

// applyRacialDecorators applies appropriate racial bonuses based on race
func (f *CharacterCreationFacade) applyRacialDecorators(char character.Character, race string) character.Character {
	switch race {
	case character.RaceElf:
		return character.NewElfDecorator(char)
	case character.RaceDwarf:
		return character.NewDwarfDecorator(char)
	case character.RaceHuman:
		return character.NewHumanDecorator(char)
	default:
		return char // No racial bonuses for unknown races
	}
}

// calculateLevelBonuses determines stat bonuses based on level
func (f *CharacterCreationFacade) calculateLevelBonuses(level int) character.Stats {
	// Simple level-based bonuses (every 4 levels = +1 to two stats)
	bonusPoints := level / 4
	return character.Stats{
		Strength:  bonusPoints,
		Dexterity: bonusPoints,
		// Other stats remain 0 for simplicity
	}
}

// getLevelAbilities determines abilities gained by level and class
func (f *CharacterCreationFacade) getLevelAbilities(class string, level int) []string {
	abilities := []string{}

	if level >= 2 {
		switch class {
		case character.ClassFighter:
			abilities = append(abilities, character.AbilityActionSurge)
		case character.ClassWizard:
			abilities = append(abilities, character.AbilityArcaneTradition)
		case character.ClassRogue:
			abilities = append(abilities, character.AbilityCunningAction)
		}
	}

	if level >= 5 {
		abilities = append(abilities, character.AbilityExtraAttack, character.AbilityImprovedAbilities)
	}

	return abilities
}
