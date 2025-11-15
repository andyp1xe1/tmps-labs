package character

import (
	"fmt"
	"slices"
)

// TemplateRegistry holds all game templates
type TemplateRegistry struct {
	RaceTemplates       map[string]*RaceTemplate
	ClassTemplates      map[string]*ClassTemplate
	BackgroundTemplates map[string]*BackgroundTemplate
	ValidCombinations   map[string][]string
}

// NewTemplateRegistry creates a registry with default templates
func NewTemplateRegistry() *TemplateRegistry {
	registry := &TemplateRegistry{
		RaceTemplates:       make(map[string]*RaceTemplate),
		ClassTemplates:      make(map[string]*ClassTemplate),
		BackgroundTemplates: make(map[string]*BackgroundTemplate),
		ValidCombinations:   make(map[string][]string),
	}
	registry.initializeDefaults()
	return registry
}

// CharacterFactory creates base characters with race/class combinations
type CharacterFactory interface {
	CreateCharacter(race, class string) (*BaseCharacter, error)
	GetAvailableRaces() []string
	GetAvailableClasses() []string
	IsValidCombination(race, class string) bool
}

// ConcreteCharacterFactory implements CharacterFactory
type ConcreteCharacterFactory struct {
	templates *TemplateRegistry
}

// NewCharacterFactory creates a new character factory
func NewCharacterFactory() CharacterFactory {
	return &ConcreteCharacterFactory{
		templates: NewTemplateRegistry(),
	}
}

// CreateCharacter creates a base character with the specified race and class
func (f *ConcreteCharacterFactory) CreateCharacter(race, class string) (*BaseCharacter, error) {
	// Validate race-class combination
	if !f.IsValidCombination(race, class) {
		return nil, fmt.Errorf("invalid combination: %s %s", race, class)
	}

	// Get race template
	raceTemplate, exists := f.templates.RaceTemplates[race]
	if !exists {
		return nil, fmt.Errorf("race %s not found", race)
	}

	// Get class template
	classTemplate, exists := f.templates.ClassTemplates[class]
	if !exists {
		return nil, fmt.Errorf("class %s not found", class)
	}

	// Create base character with race stats and class abilities
	character := &BaseCharacter{
		Name:        "", // Will be set by builder
		Race:        race,
		Class:       class,
		Stats:       raceTemplate.BaseStats,
		Background:  "", // Will be set by builder
		Personality: []string{},
		Equipment:   make([]string, len(classTemplate.Equipment)),
		Abilities:   make([]string, len(raceTemplate.Abilities)+len(classTemplate.Abilities)),
	}

	// Copy class equipment
	copy(character.Equipment, classTemplate.Equipment)

	// Combine race and class abilities
	copy(character.Abilities, raceTemplate.Abilities)
	copy(character.Abilities[len(raceTemplate.Abilities):], classTemplate.Abilities)

	return character, nil
}

// GetAvailableRaces returns all available races
func (f *ConcreteCharacterFactory) GetAvailableRaces() []string {
	races := make([]string, 0, len(f.templates.RaceTemplates))
	for race := range f.templates.RaceTemplates {
		races = append(races, race)
	}
	return races
}

// GetAvailableClasses returns all available classes
func (f *ConcreteCharacterFactory) GetAvailableClasses() []string {
	classes := make([]string, 0, len(f.templates.ClassTemplates))
	for class := range f.templates.ClassTemplates {
		classes = append(classes, class)
	}
	return classes
}

// IsValidCombination checks if race-class combination is valid
func (f *ConcreteCharacterFactory) IsValidCombination(race, class string) bool {
	validClasses, exists := f.templates.ValidCombinations[race]
	if !exists {
		return false
	}
	return slices.Contains(validClasses, class)
}

// GetBackgroundTemplate returns a background template
func (f *ConcreteCharacterFactory) GetBackgroundTemplate(background string) (*BackgroundTemplate, bool) {
	template, exists := f.templates.BackgroundTemplates[background]
	return template, exists
}

// initializeDefaults sets up default templates
func (tr *TemplateRegistry) initializeDefaults() {
	// Initialize Race Templates
	tr.RaceTemplates[RaceHuman] = &RaceTemplate{
		Name: RaceHuman,
		BaseStats: Stats{
			Strength: 10, Dexterity: 10, Constitution: 10,
			Intelligence: 10, Wisdom: 10, Charisma: 10,
		},
		Abilities:   []string{AbilityExtraSkill, "Versatile"},
		Description: "Adaptable and ambitious",
	}

	tr.RaceTemplates[RaceElf] = &RaceTemplate{
		Name: RaceElf,
		BaseStats: Stats{
			Strength: 8, Dexterity: 12, Constitution: 9,
			Intelligence: 11, Wisdom: 11, Charisma: 10,
		},
		Abilities:   []string{AbilityDarkvision, AbilityKeenSenses, AbilityFeyAncestry},
		Description: "Graceful and magical",
	}

	tr.RaceTemplates[RaceDwarf] = &RaceTemplate{
		Name: RaceDwarf,
		BaseStats: Stats{
			Strength: 11, Dexterity: 8, Constitution: 12,
			Intelligence: 10, Wisdom: 11, Charisma: 9,
		},
		Abilities:   []string{AbilityDarkvision, AbilityDwarvenResilience, AbilityStonecunning},
		Description: "Hardy and determined",
	}

	// Initialize Class Templates
	tr.ClassTemplates[ClassFighter] = &ClassTemplate{
		Name:         ClassFighter,
		HitDie:       10,
		PrimaryStats: []string{StatStrength, StatDexterity},
		Abilities:    []string{AbilityFightingStyle, AbilitySecondWind},
		Equipment:    []string{EquipmentChainMail, EquipmentShield, EquipmentSword},
		Description:  "Master of martial combat",
	}

	tr.ClassTemplates[ClassWizard] = &ClassTemplate{
		Name:         ClassWizard,
		HitDie:       6,
		PrimaryStats: []string{StatIntelligence},
		Abilities:    []string{AbilitySpellcasting, AbilityArcaneRecovery},
		Equipment:    []string{EquipmentSpellbook, EquipmentQuarterstaff, EquipmentDagger},
		Description:  "Scholar of arcane magic",
	}

	tr.ClassTemplates[ClassRogue] = &ClassTemplate{
		Name:         ClassRogue,
		HitDie:       8,
		PrimaryStats: []string{StatDexterity},
		Abilities:    []string{AbilityExpertise, AbilitySneakAttack},
		Equipment:    []string{EquipmentLeatherArmor, EquipmentShortsword, EquipmentThievesTools},
		Description:  "Skilled in stealth and precision",
	}

	// Initialize Background Templates
	tr.BackgroundTemplates[BackgroundAcolyte] = &BackgroundTemplate{
		Name:        BackgroundAcolyte,
		Skills:      []string{SkillInsight, SkillReligion},
		Equipment:   []string{EquipmentHolySymbol, EquipmentPrayerBook},
		Personality: []string{PersonalityDevout, PersonalityHelpful},
		Description: "Served in a temple",
	}

	tr.BackgroundTemplates[BackgroundCriminal] = &BackgroundTemplate{
		Name:        BackgroundCriminal,
		Skills:      []string{SkillDeception, SkillStealth},
		Equipment:   []string{EquipmentCrowbar, EquipmentDarkClothes},
		Personality: []string{PersonalitySecretive, PersonalityStreetSmart},
		Description: "Lived outside the law",
	}

	tr.BackgroundTemplates[BackgroundScholar] = &BackgroundTemplate{
		Name:        BackgroundScholar,
		Skills:      []string{SkillInvestigation, SkillHistory},
		Equipment:   []string{EquipmentBooks, EquipmentInkAndQuill},
		Personality: []string{PersonalityCurious, PersonalityMethodical},
		Description: "Spent years learning and researching",
	}

	// Initialize Valid Race-Class Combinations
	tr.ValidCombinations[RaceHuman] = []string{ClassFighter, ClassWizard, ClassRogue}
	tr.ValidCombinations[RaceElf] = []string{ClassWizard, ClassRogue, ClassFighter}
	tr.ValidCombinations[RaceDwarf] = []string{ClassFighter, ClassRogue}
}
