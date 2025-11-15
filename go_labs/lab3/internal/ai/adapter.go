package ai

import (
	"fmt"
	"log"
	"strings"
	"time"
	"tmps-go-labs/lab3/internal/character"
	"tmps-go-labs/lab3/internal/config"
)

// CharacterAIAdapter adapts between Character interface and AI communication (Adapter Pattern)
type CharacterAIAdapter struct {
	aiClient    AIProvider
	chatClient  ChatProvider // Optional: for chat session support
	config      *config.GameConfig
	chatSession ChatSession // Optional: for persistent conversations
}

// NewCharacterAIAdapter creates a new adapter for character-AI communication
func NewCharacterAIAdapter() *CharacterAIAdapter {
	cfg := config.GetGameConfig()
	var aiClient AIProvider
	var chatClient ChatProvider

	if cfg.GetGeminiAPIKey() != "" {
		geminiClient := NewGeminiClient(cfg.GetGeminiAPIKey())
		aiClient = geminiClient
		chatClient = geminiClient // GeminiClient implements both interfaces
	} else {
		mockClient := NewMockAIClient() // Fallback for testing
		aiClient = mockClient
		chatClient = mockClient // MockAIClient now also implements ChatProvider
	}

	return &CharacterAIAdapter{
		aiClient:   aiClient,
		chatClient: chatClient,
		config:     cfg,
	}
}

// SetAPIKey configures the AI client with an API key
func (caa *CharacterAIAdapter) SetAPIKey(apiKey string) {
	if err := caa.config.SetGeminiAPIKey(apiKey); err != nil {
		// Log the error but don't fail - the API key will still work in memory
		log.Printf("Warning: Failed to persist API key: %v", err)
	}
	geminiClient := NewGeminiClient(apiKey)
	caa.aiClient = geminiClient
	caa.chatClient = geminiClient
	// Close existing chat session when changing API key
	if caa.chatSession != nil {
		if err := caa.chatSession.Close(); err != nil {
			log.Printf("Warning: Failed to close existing chat session: %v", err)
		}
		caa.chatSession = nil
	}
}

// ChatWithCharacter converts character data to AI prompt and gets response (original method - backward compatible)
func (caa *CharacterAIAdapter) ChatWithCharacter(char character.Character, userMessage string) (string, error) {
	// Convert character to AI prompt context
	characterPrompt := caa.characterToPrompt(char)

	// Combine character context with user message
	fullPrompt := fmt.Sprintf("%s\n\nUser says: \"%s\"\n\nRespond as your character:", characterPrompt, userMessage)

	// Send to AI and get response (single message, no history)
	response, err := caa.aiClient.SendMessage(fullPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get AI response: %w", err)
	}

	return response, nil
}

// StartCharacterChatSession creates a persistent chat session with character context
func (caa *CharacterAIAdapter) StartCharacterChatSession(char character.Character) (ChatSession, error) {
	startTime := time.Now()
	log.Printf("⏱️  CHAT SESSION START: Beginning initialization for %s", char.GetName())

	if caa.chatClient == nil {
		log.Printf("❌ CHAT CLIENT ERROR: No chat client available")
		return nil, fmt.Errorf("chat sessions not supported with current AI client")
	}

	// Step 1: Create a new chat session
	log.Printf("📞 STEP 1/3: Creating AI chat session...")
	step1Start := time.Now()
	session, err := caa.chatClient.StartChatSession()
	step1Duration := time.Since(step1Start)
	if err != nil {
		log.Printf("❌ STEP 1 FAILED: Chat session creation failed after %v - %v", step1Duration, err)
		return nil, fmt.Errorf("failed to start chat session: %w", err)
	}
	log.Printf("✅ STEP 1 SUCCESS: Chat session created in %v", step1Duration)

	// Step 2: Generate character context (this might be slow if generating backstory/prompt)
	log.Printf("🎭 STEP 2/3: Generating character context...")
	step2Start := time.Now()
	characterPrompt := caa.characterToPrompt(char)
	step2Duration := time.Since(step2Start)
	log.Printf("✅ STEP 2 SUCCESS: Character context ready in %v (prompt length: %d chars)", step2Duration, len(characterPrompt))

	systemMessage := fmt.Sprintf("%s\n\nYou will now engage in a conversation. Stay in character at all times.", characterPrompt)

	// Step 3: Send the character context as the first message to AI
	log.Printf("💬 STEP 3/3: Initializing AI with character context...")
	step3Start := time.Now()
	_, err = session.SendMessage(systemMessage)
	step3Duration := time.Since(step3Start)
	if err != nil {
		log.Printf("❌ STEP 3 FAILED: Character context initialization failed after %v - %v", step3Duration, err)
		if closeErr := session.Close(); closeErr != nil {
			log.Printf("Warning: Failed to close session after initialization error: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to initialize character context: %w", err)
	}

	totalDuration := time.Since(startTime)
	log.Printf("🎉 CHAT SESSION READY: Total initialization took %v", totalDuration)
	log.Printf("   ├─ Step 1 (Create session): %v", step1Duration)
	log.Printf("   ├─ Step 2 (Generate context): %v", step2Duration)
	log.Printf("   └─ Step 3 (Initialize AI): %v", step3Duration)

	return session, nil
}

// ChatWithCharacterSession sends a message in an existing chat session (maintains history)
func (caa *CharacterAIAdapter) ChatWithCharacterSession(session ChatSession, userMessage string) (string, error) {
	if session == nil {
		return "", fmt.Errorf("chat session is nil")
	}

	response, err := session.SendMessage(userMessage)
	if err != nil {
		return "", fmt.Errorf("failed to send message in chat session: %w", err)
	}

	return response, nil
}

// GetChatHistory returns the conversation history from a chat session
func (caa *CharacterAIAdapter) GetChatHistory(session ChatSession) ([]ChatMessage, error) {
	if session == nil {
		return nil, fmt.Errorf("chat session is nil")
	}

	return session.GetHistory(), nil
}

// GenerateCharacterBackstory creates a backstory for the character
func (caa *CharacterAIAdapter) GenerateCharacterBackstory(char character.Character) (string, error) {
	log.Printf("🤖 AI BACKSTORY REQUEST: Starting generation for %s", char.GetName())
	startTime := time.Now()

	prompt := fmt.Sprintf(`Generate a detailed backstory for this D&D character:
Name: %s
Race: %s
Class: %s
Background: %s
Personality: %s
Stats: STR:%d DEX:%d CON:%d INT:%d WIS:%d CHA:%d

Create an engaging 2-3 paragraph backstory that explains how they became who they are.`,
		char.GetName(), char.GetRace(), char.GetClass(), char.GetBackground(),
		strings.Join(char.GetPersonality(), ", "),
		char.GetStats().Strength, char.GetStats().Dexterity, char.GetStats().Constitution,
		char.GetStats().Intelligence, char.GetStats().Wisdom, char.GetStats().Charisma)

	log.Printf("📤 AI REQUEST: Sending backstory prompt (%d chars) to API", len(prompt))
	result, err := caa.aiClient.SendMessage(prompt)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("❌ AI BACKSTORY FAILED: Request took %v, error: %v", duration, err)
		return "", err
	}

	log.Printf("✅ AI BACKSTORY SUCCESS: Generated %d characters in %v", len(result), duration)
	return result, nil
}

// GetCharacterAdvice asks the AI for character optimization advice
func (caa *CharacterAIAdapter) GetCharacterAdvice(char character.Character, context string) (string, error) {
	prompt := fmt.Sprintf(`As an experienced D&D DM, analyze this character and provide advice:
Character: %s
Context: %s

Character Details:
- Race: %s, Class: %s, Background: %s  
- Stats: STR:%d DEX:%d CON:%d INT:%d WIS:%d CHA:%d
- Abilities: %s
- Equipment: %s
- Personality: %s

Provide strategic advice for character development and gameplay.`,
		char.GetName(), context, char.GetRace(), char.GetClass(), char.GetBackground(),
		char.GetStats().Strength, char.GetStats().Dexterity, char.GetStats().Constitution,
		char.GetStats().Intelligence, char.GetStats().Wisdom, char.GetStats().Charisma,
		strings.Join(char.GetAbilities(), ", "),
		strings.Join(char.GetEquipment(), ", "),
		strings.Join(char.GetPersonality(), ", "))

	return caa.aiClient.SendMessage(prompt)
}

// IsConfigured checks if the adapter is properly configured
func (caa *CharacterAIAdapter) IsConfigured() bool {
	return caa.aiClient.IsConfigured()
}

// SupportsChatSessions checks if the current AI client supports chat sessions
func (caa *CharacterAIAdapter) SupportsChatSessions() bool {
	return caa.chatClient != nil
}

// TestConnection tests the AI service connectivity
func (caa *CharacterAIAdapter) TestConnection() error {
	return caa.aiClient.TestConnection()
}

// characterToPrompt converts Character data to AI prompt context
// Auto-generates backstory and system prompt if they're empty
func (caa *CharacterAIAdapter) characterToPrompt(char character.Character) string {
	// Check if we already have a system prompt - if so, return it immediately
	if existingPrompt := char.GetSystemPrompt(); existingPrompt != "" {
		log.Printf("🎯 CACHE HIT: Using existing system prompt for %s (length: %d chars)", char.GetName(), len(existingPrompt))
		return existingPrompt
	}

	log.Printf("🔄 CACHE MISS: Generating new system prompt for %s", char.GetName())

	// Check backstory status
	if char.GetBackstory() == "" {
		if caa.IsConfigured() {
			log.Printf("📝 GENERATING: New backstory for %s (AI call required)", char.GetName())
			// Only generate backstory once - this is the expensive operation
			backstory, err := caa.GenerateCharacterBackstory(char)
			if err == nil && backstory != "" {
				char.SetBackstory(backstory)
				log.Printf("✅ BACKSTORY CREATED: %d characters for %s", len(backstory), char.GetName())
			} else {
				// If backstory generation fails, set a placeholder to avoid retrying
				char.SetBackstory("A mysterious adventurer with an unknown past.")
				log.Printf("❌ BACKSTORY FAILED: Using placeholder for %s, error: %v", char.GetName(), err)
			}
		} else {
			log.Printf("⚠️  NO API KEY: Cannot generate backstory for %s", char.GetName())
		}
	} else {
		log.Printf("💾 BACKSTORY LOADED: Using existing backstory for %s (length: %d chars)", char.GetName(), len(char.GetBackstory()))
	}

	// Generate comprehensive system prompt
	stats := char.GetStats()
	backstory := char.GetBackstory()

	var systemPrompt strings.Builder

	// Core character identity
	systemPrompt.WriteString(fmt.Sprintf("You are %s, a %s %s with a %s background.\n\n",
		char.GetName(), char.GetRace(), char.GetClass(), char.GetBackground()))

	// Add backstory if available
	if backstory != "" {
		systemPrompt.WriteString(fmt.Sprintf("Your backstory: %s\n\n", backstory))
	}

	// Character traits and capabilities
	systemPrompt.WriteString(fmt.Sprintf(`Your personality traits: %s
Your abilities: %s
Your equipment: %s

Your stats are:
- Strength: %d
- Dexterity: %d  
- Constitution: %d
- Intelligence: %d
- Wisdom: %d
- Charisma: %d

`,
		strings.Join(char.GetPersonality(), ", "),
		strings.Join(char.GetAbilities(), ", "),
		strings.Join(char.GetEquipment(), ", "),
		stats.Strength, stats.Dexterity, stats.Constitution,
		stats.Intelligence, stats.Wisdom, stats.Charisma))

	// Behavioral instructions
	systemPrompt.WriteString("Respond to messages as this character would, staying true to their personality, background, and backstory. ")
	systemPrompt.WriteString("Keep responses engaging and in-character. ")
	systemPrompt.WriteString("Draw from your backstory and personality when appropriate.")

	finalPrompt := systemPrompt.String()

	// Cache the generated system prompt
	char.SetSystemPrompt(finalPrompt)
	log.Printf("🎭 SYSTEM PROMPT CREATED: Generated and cached %d character prompt for %s", len(finalPrompt), char.GetName())

	return finalPrompt
}
