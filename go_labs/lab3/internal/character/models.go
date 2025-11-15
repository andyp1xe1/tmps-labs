package character

import "fmt"

// Stats represents the six core ability scores
type Stats struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// Character represents the core character interface
type Character interface {
	GetName() string
	GetRace() string
	GetClass() string
	GetStats() Stats
	GetBackground() string
	GetPersonality() []string
	GetEquipment() []string
	GetAbilities() []string
	GetBackstory() string
	GetSystemPrompt() string
	SetBackstory(string)
	SetSystemPrompt(string)
	Describe() string
	Clone() Character
}

// BaseCharacter is the concrete implementation
type BaseCharacter struct {
	Name         string
	Race         string
	Class        string
	Stats        Stats
	Background   string
	Personality  []string
	Equipment    []string
	Abilities    []string
	Backstory    string // AI-generated character backstory
	SystemPrompt string // AI-generated system prompt for chat
}

func (c *BaseCharacter) GetName() string               { return c.Name }
func (c *BaseCharacter) GetRace() string               { return c.Race }
func (c *BaseCharacter) GetClass() string              { return c.Class }
func (c *BaseCharacter) GetStats() Stats               { return c.Stats }
func (c *BaseCharacter) GetBackground() string         { return c.Background }
func (c *BaseCharacter) GetPersonality() []string      { return c.Personality }
func (c *BaseCharacter) GetEquipment() []string        { return c.Equipment }
func (c *BaseCharacter) GetAbilities() []string        { return c.Abilities }
func (c *BaseCharacter) GetBackstory() string          { return c.Backstory }
func (c *BaseCharacter) GetSystemPrompt() string       { return c.SystemPrompt }
func (c *BaseCharacter) SetBackstory(backstory string) { c.Backstory = backstory }
func (c *BaseCharacter) SetSystemPrompt(prompt string) { c.SystemPrompt = prompt }

func (c *BaseCharacter) Describe() string {
	return fmt.Sprintf("%s is a %s %s with %s background",
		c.Name, c.Race, c.Class, c.Background)
}

func (c *BaseCharacter) Clone() Character {
	clone := &BaseCharacter{
		Name:         c.Name,
		Race:         c.Race,
		Class:        c.Class,
		Stats:        c.Stats,
		Background:   c.Background,
		Personality:  make([]string, len(c.Personality)),
		Equipment:    make([]string, len(c.Equipment)),
		Abilities:    make([]string, len(c.Abilities)),
		Backstory:    c.Backstory,
		SystemPrompt: c.SystemPrompt,
	}
	copy(clone.Personality, c.Personality)
	copy(clone.Equipment, c.Equipment)
	copy(clone.Abilities, c.Abilities)
	return clone
}

// RaceTemplate defines base stats and traits for races
type RaceTemplate struct {
	Name        string
	BaseStats   Stats
	Abilities   []string
	Description string
}

// ClassTemplate defines base features for classes
type ClassTemplate struct {
	Name         string
	HitDie       int
	PrimaryStats []string
	Abilities    []string
	Equipment    []string
	Description  string
}

// BackgroundTemplate defines skills and features for backgrounds
type BackgroundTemplate struct {
	Name        string
	Skills      []string
	Equipment   []string
	Personality []string
	Description string
}
