package character

import (
	"fmt"
	"math/rand"
	"time"
)

// CharacterBuilder builds detailed characters step by step
type CharacterBuilder interface {
	SetName(name string) CharacterBuilder
	SetStats(str, dex, con, intel, wis, cha int) CharacterBuilder
	SetBackground(background string) CharacterBuilder
	SetPersonality(traits []string) CharacterBuilder
	AddEquipment(equipment []string) CharacterBuilder
	AddAbilities(abilities []string) CharacterBuilder
	RandomizeStats() CharacterBuilder
	Build() (*BaseCharacter, error)
	Reset() CharacterBuilder
}

// ConcreteCharacterBuilder implements CharacterBuilder
type ConcreteCharacterBuilder struct {
	character *BaseCharacter
	factory   CharacterFactory
	errors    []string
}

// NewCharacterBuilder creates a new builder with a base character from factory
func NewCharacterBuilder(baseCharacter *BaseCharacter) CharacterBuilder {
	return &ConcreteCharacterBuilder{
		character: baseCharacter,
		factory:   NewCharacterFactory(), // For background templates
		errors:    make([]string, 0),
	}
}

// SetName sets the character's name
func (b *ConcreteCharacterBuilder) SetName(name string) CharacterBuilder {
	if name == "" {
		b.errors = append(b.errors, "name cannot be empty")
		return b
	}
	b.character.Name = name
	return b
}

// SetStats sets all ability scores
func (b *ConcreteCharacterBuilder) SetStats(str, dex, con, intel, wis, cha int) CharacterBuilder {
	// Validate stat ranges (3-18 for D&D)
	stats := []int{str, dex, con, intel, wis, cha}
	statNames := []string{StatStrength, StatDexterity, StatConstitution, StatIntelligence, StatWisdom, StatCharisma}

	for i, stat := range stats {
		if stat < 3 || stat > 18 {
			b.errors = append(b.errors, fmt.Sprintf("%s must be between 3 and 18, got %d", statNames[i], stat))
		}
	}

	b.character.Stats = Stats{
		Strength:     str,
		Dexterity:    dex,
		Constitution: con,
		Intelligence: intel,
		Wisdom:       wis,
		Charisma:     cha,
	}
	return b
}

// SetBackground sets the character's background and applies its benefits
func (b *ConcreteCharacterBuilder) SetBackground(background string) CharacterBuilder {
	// Get background template from factory
	if factory, ok := b.factory.(*ConcreteCharacterFactory); ok {
		if template, exists := factory.GetBackgroundTemplate(background); exists {
			b.character.Background = background

			// Add background equipment
			b.character.Equipment = append(b.character.Equipment, template.Equipment...)

			// Add background personality traits if not already set
			if len(b.character.Personality) == 0 {
				b.character.Personality = make([]string, len(template.Personality))
				copy(b.character.Personality, template.Personality)
			}
		} else {
			b.errors = append(b.errors, fmt.Sprintf("background %s not found", background))
		}
	}
	return b
}

// SetPersonality sets personality traits (overrides background traits)
func (b *ConcreteCharacterBuilder) SetPersonality(traits []string) CharacterBuilder {
	if len(traits) == 0 {
		b.errors = append(b.errors, "personality traits cannot be empty")
		return b
	}
	b.character.Personality = make([]string, len(traits))
	copy(b.character.Personality, traits)
	return b
}

// AddEquipment adds additional equipment to the character
func (b *ConcreteCharacterBuilder) AddEquipment(equipment []string) CharacterBuilder {
	b.character.Equipment = append(b.character.Equipment, equipment...)
	return b
}

// AddAbilities adds additional abilities to the character
func (b *ConcreteCharacterBuilder) AddAbilities(abilities []string) CharacterBuilder {
	b.character.Abilities = append(b.character.Abilities, abilities...)
	return b
}

// RandomizeStats generates random stats using 4d6 drop lowest method
func (b *ConcreteCharacterBuilder) RandomizeStats() CharacterBuilder {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	rollStat := func() int {
		rolls := make([]int, 4)
		for i := range 4 {
			rolls[i] = r.Intn(6) + 1
		}

		// Find and remove the lowest roll
		lowest := 0
		for i := range len(rolls) {
			if i > 0 && rolls[i] < rolls[lowest] {
				lowest = i
			}
		}

		sum := 0
		for i := range len(rolls) {
			if i != lowest {
				sum += rolls[i]
			}
		}
		return sum
	}

	b.character.Stats = Stats{
		Strength:     rollStat(),
		Dexterity:    rollStat(),
		Constitution: rollStat(),
		Intelligence: rollStat(),
		Wisdom:       rollStat(),
		Charisma:     rollStat(),
	}
	return b
}

// Build finalizes and validates the character
func (b *ConcreteCharacterBuilder) Build() (*BaseCharacter, error) {
	// Check for any accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("character build failed: %v", b.errors)
	}

	// Validate required fields
	if b.character.Name == "" {
		return nil, fmt.Errorf("character name is required")
	}

	if b.character.Background == "" {
		return nil, fmt.Errorf("character background is required")
	}

	// Return a copy to prevent external modification
	return b.character.Clone().(*BaseCharacter), nil
}

// Reset clears the current character and errors for reuse
func (b *ConcreteCharacterBuilder) Reset() CharacterBuilder {
	// Keep the original character structure but clear personalization
	if b.character != nil {
		b.character.Name = ""
		b.character.Background = ""
		b.character.Personality = make([]string, 0)
		// Don't clear race/class specific equipment and abilities
	}
	b.errors = make([]string, 0)
	return b
}
