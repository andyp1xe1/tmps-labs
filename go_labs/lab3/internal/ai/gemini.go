package ai

import (
	"context"
	"fmt"
	"google.golang.org/genai"
)

// GeminiClient implements both AIProvider and ChatProvider using official genai SDK
type GeminiClient struct {
	apiKey string
	client *genai.Client
	ctx    context.Context
}

// NewGeminiClient creates a new Gemini API client with 2.5 Flash model using official SDK
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		ctx:    context.Background(),
	}
}

// ensureClient initializes the genai client if not already done
func (gc *GeminiClient) ensureClient() error {
	if gc.client != nil {
		return nil
	}

	if gc.apiKey == "" {
		return fmt.Errorf("gemini API key not configured")
	}

	// Create client with API key via environment variable or option
	client, err := genai.NewClient(gc.ctx, nil) // API key should be in GEMINI_API_KEY env var
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	gc.client = client
	return nil
}

// SendMessage sends a single message to Gemini API using GenerateContent (implements AIProvider)
func (gc *GeminiClient) SendMessage(prompt string) (string, error) {
	if err := gc.ensureClient(); err != nil {
		return "", err
	}

	// Use GenerateContent for simple text generation
	result, err := gc.client.Models.GenerateContent(
		gc.ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates from Gemini")
	}

	if len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in Gemini response")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

// IsConfigured checks if the client is properly configured
func (gc *GeminiClient) IsConfigured() bool {
	return gc.apiKey != ""
}

// TestConnection tests the API connectivity with a simple request
func (gc *GeminiClient) TestConnection() error {
	if !gc.IsConfigured() {
		return fmt.Errorf("gemini API key not configured")
	}

	// Send a simple test message
	_, err := gc.SendMessage("Hello, are you working?")
	if err != nil {
		return fmt.Errorf("gemini API connection test failed: %w", err)
	}

	return nil
}

// StartChatSession creates a new chat session (implements ChatProvider)
func (gc *GeminiClient) StartChatSession() (ChatSession, error) {
	if err := gc.ensureClient(); err != nil {
		return nil, err
	}

	return NewGeminiChatSession(gc.client, gc.ctx), nil
}

// GetModelInfo returns information about the model being used
func (gc *GeminiClient) GetModelInfo() string {
	return "Gemini 2.5 Flash"
}

// GeminiChatSession implements the ChatSession interface using genai SDK
type GeminiChatSession struct {
	client  *genai.Client
	ctx     context.Context
	chat    *genai.Chat
	history []ChatMessage
}

// NewGeminiChatSession creates a new chat session
func NewGeminiChatSession(client *genai.Client, ctx context.Context) *GeminiChatSession {
	// Start with empty history - can be customized later
	initialHistory := []*genai.Content{}

	chat, err := client.Chats.Create(ctx, "gemini-2.5-flash", nil, initialHistory)
	if err != nil {
		// If chat creation fails, we'll handle it in SendMessage
		return &GeminiChatSession{
			client:  client,
			ctx:     ctx,
			chat:    nil,
			history: make([]ChatMessage, 0),
		}
	}

	return &GeminiChatSession{
		client:  client,
		ctx:     ctx,
		chat:    chat,
		history: make([]ChatMessage, 0),
	}
}

// NewGeminiChatSessionWithHistory creates a chat session with initial history
func NewGeminiChatSessionWithHistory(client *genai.Client, ctx context.Context, history []ChatMessage) (*GeminiChatSession, error) {
	// Convert our ChatMessage format to genai.Content format
	genaiHistory := make([]*genai.Content, 0, len(history))
	for _, msg := range history {
		var role genai.Role
		if msg.Role == "user" {
			role = genai.RoleUser
		} else {
			role = genai.RoleModel
		}
		genaiHistory = append(genaiHistory, genai.NewContentFromText(msg.Content, role))
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.5-flash", nil, genaiHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat with history: %w", err)
	}

	return &GeminiChatSession{
		client:  client,
		ctx:     ctx,
		chat:    chat,
		history: append([]ChatMessage{}, history...), // Copy the history
	}, nil
}

// SendMessage sends a message in the chat session
func (gcs *GeminiChatSession) SendMessage(message string) (string, error) {
	// Ensure chat is created
	if gcs.chat == nil {
		chat, err := gcs.client.Chats.Create(gcs.ctx, "gemini-2.5-flash", nil, []*genai.Content{})
		if err != nil {
			return "", fmt.Errorf("failed to create chat: %w", err)
		}
		gcs.chat = chat
	}

	// Add user message to history
	gcs.history = append(gcs.history, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// Send message to Gemini
	res, err := gcs.chat.SendMessage(gcs.ctx, genai.Part{Text: message})
	if err != nil {
		return "", fmt.Errorf("failed to send message to chat: %w", err)
	}

	// Extract response
	if len(res.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates from Gemini")
	}

	if len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in Gemini response")
	}

	response := res.Candidates[0].Content.Parts[0].Text

	// Add model response to history
	gcs.history = append(gcs.history, ChatMessage{
		Role:    "model",
		Content: response,
	})

	return response, nil
}

// GetHistory returns the chat history
func (gcs *GeminiChatSession) GetHistory() []ChatMessage {
	// Return a copy to prevent external modification
	history := make([]ChatMessage, len(gcs.history))
	copy(history, gcs.history)
	return history
}

// Close closes the chat session
func (gcs *GeminiChatSession) Close() error {
	// genai.Chat doesn't have a close method, so we just cleanup references
	gcs.chat = nil
	gcs.history = nil
	return nil
}
