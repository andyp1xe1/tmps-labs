package character

import "fmt"

// CharacterDecorator wraps a Character and adds additional functionality
type CharacterDecorator struct {
	Character Character
}

// Implement Character interface by delegating to wrapped character
func (cd *CharacterDecorator) GetName() string               { return cd.Character.GetName() }
func (cd *CharacterDecorator) GetRace() string               { return cd.Character.GetRace() }
func (cd *CharacterDecorator) GetClass() string              { return cd.Character.GetClass() }
func (cd *CharacterDecorator) GetStats() Stats               { return cd.Character.GetStats() }
func (cd *CharacterDecorator) GetBackground() string         { return cd.Character.GetBackground() }
func (cd *CharacterDecorator) GetPersonality() []string      { return cd.Character.GetPersonality() }
func (cd *CharacterDecorator) GetEquipment() []string        { return cd.Character.GetEquipment() }
func (cd *CharacterDecorator) GetAbilities() []string        { return cd.Character.GetAbilities() }
func (cd *CharacterDecorator) GetBackstory() string          { return cd.Character.GetBackstory() }
func (cd *CharacterDecorator) GetSystemPrompt() string       { return cd.Character.GetSystemPrompt() }
func (cd *CharacterDecorator) SetBackstory(backstory string) { cd.Character.SetBackstory(backstory) }
func (cd *CharacterDecorator) SetSystemPrompt(prompt string) { cd.Character.SetSystemPrompt(prompt) }
func (cd *CharacterDecorator) Describe() string              { return cd.Character.Describe() }
func (cd *CharacterDecorator) Clone() Character              { return cd.Character.Clone() }

// RacialBonusDecorator adds racial stat bonuses and abilities
type RacialBonusDecorator struct {
	*CharacterDecorator
	bonusStats     Stats
	bonusAbilities []string
	description    string
}

// NewRacialBonusDecorator creates a decorator for racial bonuses
func NewRacialBonusDecorator(character Character, bonusStats Stats, bonusAbilities []string, description string) Character {
	return &RacialBonusDecorator{
		CharacterDecorator: &CharacterDecorator{Character: character},
		bonusStats:         bonusStats,
		bonusAbilities:     bonusAbilities,
		description:        description,
	}
}

// GetStats returns base stats plus racial bonuses
func (rbd *RacialBonusDecorator) GetStats() Stats {
	baseStats := rbd.Character.GetStats()
	return Stats{
		Strength:     baseStats.Strength + rbd.bonusStats.Strength,
		Dexterity:    baseStats.Dexterity + rbd.bonusStats.Dexterity,
		Constitution: baseStats.Constitution + rbd.bonusStats.Constitution,
		Intelligence: baseStats.Intelligence + rbd.bonusStats.Intelligence,
		Wisdom:       baseStats.Wisdom + rbd.bonusStats.Wisdom,
		Charisma:     baseStats.Charisma + rbd.bonusStats.Charisma,
	}
}

// GetAbilities returns base abilities plus racial abilities
func (rbd *RacialBonusDecorator) GetAbilities() []string {
	baseAbilities := rbd.Character.GetAbilities()
	allAbilities := make([]string, len(baseAbilities)+len(rbd.bonusAbilities))
	copy(allAbilities, baseAbilities)
	copy(allAbilities[len(baseAbilities):], rbd.bonusAbilities)
	return allAbilities
}

// Describe returns enhanced description with racial features
func (rbd *RacialBonusDecorator) Describe() string {
	baseDesc := rbd.Character.Describe()
	return fmt.Sprintf("%s. %s", baseDesc, rbd.description)
}

// EquipmentDecorator adds magical or enhanced equipment
type EquipmentDecorator struct {
	*CharacterDecorator
	magicalEquipment []string
	equipmentBonuses Stats
	description      string
}

// NewEquipmentDecorator creates a decorator for magical equipment
func NewEquipmentDecorator(character Character, equipment []string, bonuses Stats, description string) Character {
	return &EquipmentDecorator{
		CharacterDecorator: &CharacterDecorator{Character: character},
		magicalEquipment:   equipment,
		equipmentBonuses:   bonuses,
		description:        description,
	}
}

// GetEquipment returns base equipment plus magical items
func (ed *EquipmentDecorator) GetEquipment() []string {
	baseEquipment := ed.Character.GetEquipment()
	allEquipment := make([]string, len(baseEquipment)+len(ed.magicalEquipment))
	copy(allEquipment, baseEquipment)
	copy(allEquipment[len(baseEquipment):], ed.magicalEquipment)
	return allEquipment
}

// GetStats returns base stats plus equipment bonuses
func (ed *EquipmentDecorator) GetStats() Stats {
	baseStats := ed.Character.GetStats()
	return Stats{
		Strength:     baseStats.Strength + ed.equipmentBonuses.Strength,
		Dexterity:    baseStats.Dexterity + ed.equipmentBonuses.Dexterity,
		Constitution: baseStats.Constitution + ed.equipmentBonuses.Constitution,
		Intelligence: baseStats.Intelligence + ed.equipmentBonuses.Intelligence,
		Wisdom:       baseStats.Wisdom + ed.equipmentBonuses.Wisdom,
		Charisma:     baseStats.Charisma + ed.equipmentBonuses.Charisma,
	}
}

// Describe returns enhanced description with equipment
func (ed *EquipmentDecorator) Describe() string {
	baseDesc := ed.Character.Describe()
	return fmt.Sprintf("%s. %s", baseDesc, ed.description)
}

// ExperienceDecorator adds level-based bonuses and abilities
type ExperienceDecorator struct {
	*CharacterDecorator
	level          int
	levelAbilities []string
	levelBonuses   Stats
}

// NewExperienceDecorator creates a decorator for experience-based improvements
func NewExperienceDecorator(character Character, level int, abilities []string, bonuses Stats) Character {
	return &ExperienceDecorator{
		CharacterDecorator: &CharacterDecorator{Character: character},
		level:              level,
		levelAbilities:     abilities,
		levelBonuses:       bonuses,
	}
}

// GetAbilities returns base abilities plus level abilities
func (ed *ExperienceDecorator) GetAbilities() []string {
	baseAbilities := ed.Character.GetAbilities()
	allAbilities := make([]string, len(baseAbilities)+len(ed.levelAbilities))
	copy(allAbilities, baseAbilities)
	copy(allAbilities[len(baseAbilities):], ed.levelAbilities)
	return allAbilities
}

// GetStats returns base stats plus level bonuses
func (ed *ExperienceDecorator) GetStats() Stats {
	baseStats := ed.Character.GetStats()
	return Stats{
		Strength:     baseStats.Strength + ed.levelBonuses.Strength,
		Dexterity:    baseStats.Dexterity + ed.levelBonuses.Dexterity,
		Constitution: baseStats.Constitution + ed.levelBonuses.Constitution,
		Intelligence: baseStats.Intelligence + ed.levelBonuses.Intelligence,
		Wisdom:       baseStats.Wisdom + ed.levelBonuses.Wisdom,
		Charisma:     baseStats.Charisma + ed.levelBonuses.Charisma,
	}
}

// Describe returns enhanced description with level information
func (ed *ExperienceDecorator) Describe() string {
	baseDesc := ed.Character.Describe()
	return fmt.Sprintf("%s (Level %d)", baseDesc, ed.level)
}

// Helper functions for common racial decorators
func NewElfDecorator(character Character) Character {
	bonuses := Stats{Dexterity: 2}
	abilities := []string{AbilityDarkvision, AbilityKeenSenses, AbilityFeyAncestry, AbilityTrance}
	return NewRacialBonusDecorator(character, bonuses, abilities, "Enhanced elven grace and magical heritage")
}

func NewDwarfDecorator(character Character) Character {
	bonuses := Stats{Constitution: 2}
	abilities := []string{AbilityDarkvision, AbilityDwarvenResilience, AbilityDwarvenCombatTraining, AbilityStonecunning}
	return NewRacialBonusDecorator(character, bonuses, abilities, "Dwarven toughness and mountain wisdom")
}

func NewHumanDecorator(character Character) Character {
	bonuses := Stats{Strength: 1, Dexterity: 1, Constitution: 1, Intelligence: 1, Wisdom: 1, Charisma: 1}
	abilities := []string{AbilityExtraLanguage, AbilityExtraSkill, AbilityHumanDetermination}
	return NewRacialBonusDecorator(character, bonuses, abilities, "Human versatility and adaptability")
}
