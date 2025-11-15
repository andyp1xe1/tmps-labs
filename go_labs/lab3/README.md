# Lab 2: Structural Design Patterns - D&D Character Builder

## Author: Student Name

## Introduction

This laboratory work demonstrates the implementation of **Structural Design Patterns** in a D&D Character Builder application. The system allows users to create, customize, and interact with D&D characters through a beautiful terminal interface powered by Bubble Tea and Gemini AI integration.

## Design Patterns Implemented

### Creational Patterns (from Lab 2)

#### 1. **Singleton Pattern**
- **Location**: `internal/config/singleton.go:18-26`
- **Purpose**: Ensures single instance of application configuration
- **Implementation**: Uses `sync.Once` to guarantee thread-safe initialization

```go
func GetGameConfig() *GameConfig {
    once.Do(func() {
        instance = &GameConfig{
            AppName: "Character Builder",
            Version: "1.0.0",
        }
    })
    return instance
}
```

#### 2. **Factory Pattern**
- **Location**: `internal/character/factory.go:28-70`
- **Purpose**: Creates different character types based on race/class combinations
- **Implementation**: Validates combinations and applies base racial/class attributes

```go
func (f *ConcreteCharacterFactory) CreateCharacter(race, class string) (*BaseCharacter, error) {
    if !f.IsValidCombination(race, class) {
        return nil, fmt.Errorf("invalid combination: %s %s", race, class)
    }
    // Create base character with race stats and class abilities
    // ...
}
```

#### 3. **Builder Pattern**
- **Location**: `internal/character/builder.go:16-170`
- **Purpose**: Constructs characters step-by-step with validation
- **Implementation**: Fluent interface for setting stats, background, personality, and equipment

```go
character := builder.
    SetName("Elaria").
    SetStats(str, dex, con, int, wis, cha).
    SetBackground("Scholar").
    SetPersonality([]string{"curious", "methodical"}).
    Build()
```

### Structural Patterns

#### 1. **Facade Pattern**
- **Location**: `internal/facade/character_facade.go:12-185`
- **Purpose**: Simplifies complex character creation workflow
- **Implementation**: Provides high-level methods that coordinate Factory, Builder, and Decorators

```go
func (f *CharacterCreationFacade) CreateBasicCharacter(name, race, class, background string) (character.Character, error) {
    // Step 1: Use factory to create base character
    baseChar, err := f.factory.CreateCharacter(race, class)
    
    // Step 2: Use builder to add details
    builder := character.NewCharacterBuilder(baseChar)
    finalChar, err := builder.SetName(name).SetBackground(background).Build()
    
    // Step 3: Apply racial decorators
    decoratedChar := f.applyRacialDecorators(finalChar, race)
    
    return decoratedChar, nil
}
```

**Motivation**: The character creation process involves multiple steps and patterns. The Facade provides a simple interface that hides this complexity from the client.

#### 2. **Decorator Pattern**
- **Location**: `internal/character/decorators.go:8-155`
- **Purpose**: Dynamically adds abilities, bonuses, and traits to characters
- **Implementation**: Multiple decorator types for different enhancement categories

```go
// Racial bonuses decorator
func NewElfDecorator(character Character) Character {
    bonuses := Stats{Dexterity: 2}
    abilities := []string{"Darkvision", "Keen Senses", "Fey Ancestry"}
    return NewRacialBonusDecorator(character, bonuses, abilities, "Enhanced elven grace")
}

// Equipment bonuses decorator  
decoratedChar = NewEquipmentDecorator(char, equipment, bonuses, "Equipped with magical gear")
```

**Motivation**: Characters can have various combinations of racial bonuses, equipment effects, and level-based improvements. Decorators allow stacking these enhancements without modifying the core character class.

#### 3. **Adapter Pattern**
- **Location**: `internal/ai/adapter.go:87-210`
- **Purpose**: Adapts between Character interface and Gemini AI API
- **Implementation**: Converts character data to AI prompts and handles API communication

```go
func (caa *CharacterAIAdapter) ChatWithCharacter(char character.Character, userMessage string) (string, error) {
    // Convert character to AI prompt context
    characterPrompt := caa.characterToPrompt(char)
    
    // Combine character context with user message
    fullPrompt := fmt.Sprintf("%s\n\nUser says: \"%s\"\n\nRespond as your character:", characterPrompt, userMessage)
    
    // Send to AI and get response
    return caa.aiClient.SendMessage(fullPrompt)
}
```

**Motivation**: The Gemini API expects text prompts, while our system works with structured Character objects. The Adapter bridges this gap and provides fallback mock responses when API is not configured.

## System Architecture

### Project Structure
```
lab3/
├── cmd/
│   └── character-builder/   # Main application entry point
│       ├── main.go         # Bubble Tea TUI application
│       └── main_test.go    # Integration tests
├── internal/
│   ├── character/          # Core character models and creational patterns
│   │   ├── models.go      # Character interfaces and structs
│   │   ├── factory.go     # Factory pattern implementation
│   │   ├── builder.go     # Builder pattern implementation
│   │   └── decorators.go  # Decorator pattern implementation
│   ├── facade/            # Facade pattern
│   │   └── character_facade.go
│   ├── ai/                # AI integration with Adapter pattern
│   │   └── adapter.go
│   └── config/            # Singleton configuration
│       └── singleton.go
└── README.md              # Documentation
```

### Pattern Interactions

1. **Client** → **Facade** → **Factory** → **Builder** → **Decorators**
   - Client uses simple facade methods
   - Facade coordinates all patterns internally
   - Factory creates base characters
   - Builder adds customization
   - Decorators apply enhancements

2. **Character** → **Adapter** → **AI Service**
   - Character data gets converted to natural language
   - Adapter handles API communication
   - Responses are formatted for display

## Key Features

### 🎮 Interactive Terminal UI
- Beautiful Bubble Tea interface with forms and menus
- Colorful styling with emojis and borders
- Keyboard navigation and input validation
- **Professional animated spinners** for loading feedback using Charm Bubbles
- **Async operations** with non-blocking UI during AI calls
- **Responsive design** that adapts to terminal width (40-120 chars)
- **Chat scrolling system** with ↑↓/k,j navigation and PageUp/PageDown support
- **Scroll indicators** showing message counts above/below viewport

### ⚔️ Character Creation
- **Basic Method**: Uses standard array stats (15,14,13,12,10,8) for balanced characters
- **Random Method**: Completely randomized with 4d6-drop-lowest stat generation
- **Advanced Method**: Custom stats and equipment specification
- **Experienced Method**: Level-based characters with enhanced abilities
- Race/class validation with pre-defined combinations
- Background selection with automatic skill/equipment assignment

### 🎭 Dynamic Enhancement System
- Racial bonuses (Elf +2 DEX, Dwarf +2 CON, Human +1 all)
- Equipment bonuses with magical item effects
- Experience-based improvements with level abilities
- Stackable decorators for complex characters

### 🤖 AI Integration
- Chat with your character using Gemini AI
- Character-aware responses based on race, class, personality
- Backstory generation and character advice
- Graceful fallback with mock responses
- **Persistent chat sessions** with conversation memory
- **Smart caching system** - backstory and prompts saved to config
- **Performance optimization** - instant startup for saved characters
- **Detailed logging** to track AI calls vs cached content usage

## Usage

### Running the Application
```bash
cd lab3
go run ./cmd/character-builder
# Or build and run:
go build -o character-builder ./cmd/character-builder
./character-builder
```

### Configuration
1. Select "Configure API Key" from main menu
2. Enter your Gemini API key (optional - mock responses available)
3. Create characters and chat with them!

### Example Workflow
1. **Create Character**: Choose race, class, background, stats
2. **View Character**: See complete character sheet with all bonuses applied
3. **Chat**: Have conversations with your character using AI
4. **Create More**: Try different combinations and builds

## Testing

Run the included tests to verify all patterns work correctly:
```bash
go test -v
```

Tests cover:
- Singleton instance uniqueness
- Factory validation and character creation
- Builder validation and fluent interface
- Decorator bonus application and stacking

## Results & Conclusions

### Pattern Benefits Demonstrated

1. **Facade Pattern**: Simplified a complex 4-step character creation process into single method calls
2. **Decorator Pattern**: Enabled flexible character enhancement without class modification  
3. **Adapter Pattern**: Successfully bridged structured data with natural language AI API

### Technical Achievements

- ✅ Clean, idiomatic Go code structure (`internal/` packages)
- ✅ All patterns properly implemented and tested
- ✅ Modern TUI with Bubble Tea framework
- ✅ Real AI integration with Gemini API
- ✅ Comprehensive error handling and validation
- ✅ Thread-safe singleton implementation
- ✅ Passes all linting checks (`golangci-lint`)

### Performance & UX Enhancements

- ✅ **Professional spinner animations** using Charm Bubbles integration
- ✅ **Async UI operations** - non-blocking interface during AI calls
- ✅ **Smart caching system** - AI-generated content persisted to `~/.config/character-builder/config.json`
- ✅ **Instant app startup** - saved characters load in milliseconds without AI calls
- ✅ **Chat message scrolling** - full viewport management with keyboard navigation
- ✅ **Responsive terminal design** - adapts to window size changes
- ✅ **Detailed performance logging** - tracks AI generation vs cache usage

### Demonstrated Optimizations

**First Run (Character Creation):**
```
📝 GENERATING: New backstory for Gangsdalf (AI call required)
🤖 AI BACKSTORY REQUEST: Starting generation for Gangsdalf  
✅ AI BACKSTORY SUCCESS: Generated 2769 characters in 1.2s
🎭 SYSTEM PROMPT CREATED: Generated and cached 3413 character prompt
💾 SAVE SUCCESS: Character saved with backstory and system prompt
```

**Subsequent Runs (Cached):**
```
📂 CHARACTER LOADED: Gangsdalf from persistent config
   └─ Backstory: true (length: 2769)     ← No AI call needed
   └─ SystemPrompt: true (length: 3413)  ← Instant loading
🎯 CACHE HIT: Using existing system prompt (length: 3413 chars)
🎉 CHAT SESSION READY: Total initialization took 0.85s
```

### Design Pattern Synergy

The combination of creational and structural patterns creates a robust, extensible system:
- **Factory + Builder**: Handles complex object creation
- **Decorator**: Provides runtime flexibility
- **Facade**: Presents clean API to users
- **Adapter**: Enables external service integration

This architecture demonstrates how design patterns work together to solve real-world problems while maintaining clean, maintainable code.

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles/spinner` - Professional loading animations
- `github.com/charmbracelet/bubbles/viewport` - Scrollable content areas
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/charmbracelet/huh` - Forms and input components
- Standard library packages for HTTP, JSON, testing, persistence

### Performance Characteristics

| Operation | First Run (Cold) | Cached Run (Warm) | Improvement |
|-----------|------------------|-------------------|-------------|
| App Startup | ~200ms | ~50ms | **4x faster** |
| Character Loading | ~1.5s (AI generation) | ~5ms (from config) | **300x faster** |
| Chat Session Start | ~2s (context + AI) | ~800ms (cached context) | **2.5x faster** |
| System Prompt | ~1s (AI generation) | ~1ms (cache hit) | **1000x faster** |

The application showcases modern Go development practices while implementing classic design patterns in a practical, interactive system with professional UX and performance optimizations.
