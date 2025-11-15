package ai

import "fmt"

// MockAIClient provides a fallback for when Gemini API is not configured
// It implements both AIProvider and ChatProvider interfaces
type MockAIClient struct{}

// NewMockAIClient creates a mock AI client for testing
func NewMockAIClient() *MockAIClient {
	return &MockAIClient{}
}

// SendMessage returns a mock response (implements AIProvider)
func (mac *MockAIClient) SendMessage(prompt string) (string, error) {
	return "I am a brave adventurer ready to face any challenge! (Note: This is a mock response - configure Gemini API for real AI interaction)", nil
}

// IsConfigured always returns false for mock client (implements AIProvider)
func (mac *MockAIClient) IsConfigured() bool {
	return false
}

// TestConnection always returns an error for mock client (implements AIProvider)
func (mac *MockAIClient) TestConnection() error {
	return fmt.Errorf("mock client cannot test real API connection - configure Gemini API key")
}

// StartChatSession creates a mock chat session (implements ChatProvider)
func (mac *MockAIClient) StartChatSession() (ChatSession, error) {
	return NewMockChatSession(), nil
}

// MockChatSession implements ChatSession interface for testing
type MockChatSession struct {
	history []ChatMessage
}

// NewMockChatSession creates a new mock chat session
func NewMockChatSession() *MockChatSession {
	return &MockChatSession{
		history: make([]ChatMessage, 0),
	}
}

// SendMessage sends a message in the mock chat session
func (mcs *MockChatSession) SendMessage(message string) (string, error) {
	// Add user message to history
	mcs.history = append(mcs.history, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// Generate mock response
	response := fmt.Sprintf("Mock character response to: '%s' (This is a mock response - configure Gemini API for real AI interaction)", message)

	// Add mock response to history
	mcs.history = append(mcs.history, ChatMessage{
		Role:    "model",
		Content: response,
	})

	return response, nil
}

// GetHistory returns the chat history
func (mcs *MockChatSession) GetHistory() []ChatMessage {
	// Return a copy to prevent external modification
	history := make([]ChatMessage, len(mcs.history))
	copy(history, mcs.history)
	return history
}

// Close closes the mock chat session
func (mcs *MockChatSession) Close() error {
	mcs.history = nil
	return nil
}
