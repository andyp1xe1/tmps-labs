package ai

import (
	"os"
	"strings"
	"testing"
	"tmps-go-labs/lab3/internal/character"
)

// TestMockAIClient tests the mock implementation
func TestMockAIClient(t *testing.T) {
	client := NewMockAIClient()

	// Test IsConfigured
	if client.IsConfigured() {
		t.Error("Mock client should not be configured")
	}

	// Test SendMessage
	response, err := client.SendMessage("Test message")
	if err != nil {
		t.Errorf("Mock client should not return error: %v", err)
	}
	if response == "" {
		t.Error("Mock client should return a response")
	}
	if !strings.Contains(response, "mock response") {
		t.Error("Mock response should indicate it's a mock")
	}

	// Test TestConnection
	err = client.TestConnection()
	if err == nil {
		t.Error("Mock client should return error for TestConnection")
	}
}

// TestGeminiClientConfiguration tests client configuration
func TestGeminiClientConfiguration(t *testing.T) {
	// Test with empty API key
	client := NewGeminiClient("")
	if client.IsConfigured() {
		t.Error("Client with empty API key should not be configured")
	}

	// Test TestConnection without API key
	err := client.TestConnection()
	if err == nil {
		t.Error("TestConnection should fail without API key")
	}

	// Test with non-empty API key (but don't actually call API)
	client = NewGeminiClient("test-key")
	if !client.IsConfigured() {
		t.Error("Client with API key should be configured")
	}
}

// TestCharacterAIAdapter tests the adapter functionality
func TestCharacterAIAdapter(t *testing.T) {
	// Test with mock client (no API key)
	adapter := NewCharacterAIAdapter()

	if adapter.IsConfigured() {
		t.Error("Adapter without API key should not be configured")
	}

	// Test TestConnection
	err := adapter.TestConnection()
	if err == nil {
		t.Error("Adapter without API key should fail TestConnection")
	}
}

// TestGeminiAPIConnectivity tests real API connectivity (requires GEMINI_API_KEY env var)
func TestGeminiAPIConnectivity(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping Gemini API connectivity test - GEMINI_API_KEY not set")
	}

	// Create client with real API key
	client := NewGeminiClient(apiKey)

	if !client.IsConfigured() {
		t.Error("Client with API key should be configured")
	}

	t.Log("Testing Gemini 2.5 Flash API connectivity...")

	// Test connectivity
	err := client.TestConnection()
	if err != nil {
		t.Fatalf("Gemini API connectivity test failed: %v", err)
	}

	t.Log("✅ Gemini 2.5 Flash API connectivity successful")

	// Test actual message sending
	response, err := client.SendMessage("Hello! Please respond with just 'API Test Successful'")
	if err != nil {
		t.Fatalf("Failed to send message to Gemini API: %v", err)
	}

	if response == "" {
		t.Error("Gemini API returned empty response")
	}

	t.Logf("✅ Gemini API response received: %s", response)
}

// TestCharacterChatWithGemini tests character chat functionality (requires GEMINI_API_KEY env var)
func TestCharacterChatWithGemini(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping Gemini character chat test - GEMINI_API_KEY not set")
	}

	// Create a test character
	baseChar := &character.BaseCharacter{
		Name:        "Gandalf",
		Race:        character.RaceHuman,
		Class:       character.ClassWizard,
		Background:  character.BackgroundScholar,
		Stats:       character.Stats{Strength: 10, Dexterity: 12, Constitution: 14, Intelligence: 18, Wisdom: 16, Charisma: 15},
		Abilities:   []string{character.AbilitySpellcasting, character.AbilityArcaneRecovery},
		Equipment:   []string{character.EquipmentSpellbook, character.EquipmentQuarterstaff},
		Personality: []string{character.PersonalityCurious, "wise"},
	}

	// Apply decorator to get full character
	testChar := character.NewHumanDecorator(baseChar)

	// Create adapter and test
	adapter := NewCharacterAIAdapter()
	adapter.SetAPIKey(apiKey)

	t.Log("Testing character chat with Gemini 2.5 Flash...")

	response, err := adapter.ChatWithCharacter(testChar, "Hello, who are you?")
	if err != nil {
		t.Fatalf("Character chat failed: %v", err)
	}

	if response == "" {
		t.Error("Character chat returned empty response")
	}

	// Check if response seems character-appropriate (should mention character name or role)
	if !strings.Contains(strings.ToLower(response), "gandalf") &&
		!strings.Contains(strings.ToLower(response), "wizard") {
		t.Logf("Warning: Response may not be character-specific: %s", response)
	}

	t.Logf("✅ Character chat successful: %s", response)
}

// TestBackstoryGeneration tests AI backstory generation (requires GEMINI_API_KEY env var)
func TestBackstoryGeneration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping backstory generation test - GEMINI_API_KEY not set")
	}

	// Create a test character
	baseChar := &character.BaseCharacter{
		Name:        "Thorin",
		Race:        character.RaceDwarf,
		Class:       character.ClassFighter,
		Background:  character.BackgroundCriminal,
		Stats:       character.Stats{Strength: 16, Dexterity: 10, Constitution: 15, Intelligence: 12, Wisdom: 13, Charisma: 8},
		Personality: []string{character.PersonalitySecretive, "gruff"},
	}

	testChar := character.NewDwarfDecorator(baseChar)

	adapter := NewCharacterAIAdapter()
	adapter.SetAPIKey(apiKey)

	t.Log("Testing backstory generation with Gemini 2.5 Flash...")

	backstory, err := adapter.GenerateCharacterBackstory(testChar)
	if err != nil {
		t.Fatalf("Backstory generation failed: %v", err)
	}

	if backstory == "" {
		t.Error("Generated backstory is empty")
	}

	if len(backstory) < 100 {
		t.Errorf("Generated backstory seems too short: %d characters", len(backstory))
	}

	t.Logf("✅ Backstory generation successful (%d chars): %s...", len(backstory), backstory[:100])
}

// TestAdviceGeneration tests character advice functionality (requires GEMINI_API_KEY env var)
func TestAdviceGeneration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping advice generation test - GEMINI_API_KEY not set")
	}

	// Create a test character
	baseChar := &character.BaseCharacter{
		Name:       "Legolas",
		Race:       character.RaceElf,
		Class:      character.ClassRogue,
		Background: character.BackgroundCriminal,
		Stats:      character.Stats{Strength: 12, Dexterity: 18, Constitution: 14, Intelligence: 14, Wisdom: 15, Charisma: 13},
	}

	testChar := character.NewElfDecorator(baseChar)

	adapter := NewCharacterAIAdapter()
	adapter.SetAPIKey(apiKey)

	t.Log("Testing character advice generation with Gemini 2.5 Flash...")

	advice, err := adapter.GetCharacterAdvice(testChar, "I'm starting a new campaign")
	if err != nil {
		t.Fatalf("Advice generation failed: %v", err)
	}

	if advice == "" {
		t.Error("Generated advice is empty")
	}

	t.Logf("✅ Advice generation successful: %s...", advice[:min(100, len(advice))])
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BenchmarkGeminiAPI benchmarks API response time (requires GEMINI_API_KEY env var)
func BenchmarkGeminiAPI(b *testing.B) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		b.Skip("Skipping Gemini API benchmark - GEMINI_API_KEY not set")
	}

	client := NewGeminiClient(apiKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.SendMessage("Hello")
		if err != nil {
			b.Fatalf("API call failed: %v", err)
		}
	}
}
