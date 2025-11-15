package ai

import (
	"strings"
	"testing"
	"tmps-go-labs/lab3/internal/character"
)

// TestChatHistoryFunctionality tests that chat sessions maintain conversation history
func TestChatHistoryFunctionality(t *testing.T) {
	// Create a test character
	baseChar := &character.BaseCharacter{
		Name:        "ChatBot",
		Race:        character.RaceElf,
		Class:       character.ClassWizard,
		Background:  character.BackgroundScholar,
		Stats:       character.Stats{Strength: 8, Dexterity: 14, Constitution: 12, Intelligence: 18, Wisdom: 16, Charisma: 13},
		Abilities:   []string{character.AbilitySpellcasting, character.AbilityArcaneRecovery},
		Equipment:   []string{character.EquipmentSpellbook, character.EquipmentQuarterstaff},
		Personality: []string{character.PersonalityCurious, character.PersonalitySecretive},
	}

	testChar := character.NewElfDecorator(baseChar)
	adapter := NewCharacterAIAdapter()

	// Test that the adapter supports chat sessions
	if !adapter.SupportsChatSessions() {
		t.Skip("Adapter doesn't support chat sessions - using mock client")
	}

	// Start a chat session
	session, err := adapter.StartCharacterChatSession(testChar)
	if err != nil {
		t.Fatalf("Failed to start chat session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			t.Logf("Warning: Failed to close session: %v", err)
		}
	}()

	// First message
	firstMessage := "Hello, what is your name?"
	firstResponse, err := adapter.ChatWithCharacterSession(session, firstMessage)
	if err != nil {
		t.Fatalf("Failed to send first message: %v", err)
	}

	if firstResponse == "" {
		t.Error("First response should not be empty")
	}

	// Second message that references the first
	secondMessage := "What did you just tell me about yourself?"
	secondResponse, err := adapter.ChatWithCharacterSession(session, secondMessage)
	if err != nil {
		t.Fatalf("Failed to send second message: %v", err)
	}

	if secondResponse == "" {
		t.Error("Second response should not be empty")
	}

	// Get chat history
	history, err := adapter.GetChatHistory(session)
	if err != nil {
		t.Fatalf("Failed to get chat history: %v", err)
	}

	// Verify history contains our messages
	if len(history) < 4 { // 2 user messages + 2 model responses (minimum)
		t.Errorf("Expected at least 4 messages in history, got %d", len(history))
	}

	// Check that our messages are in the history
	foundFirstMessage := false
	foundSecondMessage := false
	for _, msg := range history {
		if msg.Role == "user" {
			if strings.Contains(msg.Content, "what is your name") {
				foundFirstMessage = true
			}
			if strings.Contains(msg.Content, "What did you just tell me") {
				foundSecondMessage = true
			}
		}
	}

	if !foundFirstMessage {
		t.Error("First user message not found in history")
	}
	if !foundSecondMessage {
		t.Error("Second user message not found in history")
	}

	// Verify the character's system prompt was generated and cached
	systemPrompt := testChar.GetSystemPrompt()
	if systemPrompt == "" {
		t.Error("System prompt should be generated and cached")
	}

	if !strings.Contains(systemPrompt, testChar.GetName()) {
		t.Error("System prompt should contain character name")
	}

	t.Logf("✅ Chat history test passed with %d messages", len(history))
	t.Logf("✅ System prompt generated (%d chars): %s...", len(systemPrompt), systemPrompt[:100])
}
