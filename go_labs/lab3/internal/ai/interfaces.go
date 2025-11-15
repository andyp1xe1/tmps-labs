package ai

// AIProvider defines the interface for AI communication
type AIProvider interface {
	SendMessage(prompt string) (string, error)
	IsConfigured() bool
	TestConnection() error
}

// ChatProvider extends AIProvider with chat session capabilities
type ChatProvider interface {
	AIProvider
	StartChatSession() (ChatSession, error)
}

// ChatSession represents an ongoing conversation with context
type ChatSession interface {
	SendMessage(message string) (string, error)
	GetHistory() []ChatMessage
	Close() error
}

// ChatMessage represents a single message in a chat session
type ChatMessage struct {
	Role    string // "user" or "model"
	Content string
}

// APIError represents errors returned by AI APIs
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
