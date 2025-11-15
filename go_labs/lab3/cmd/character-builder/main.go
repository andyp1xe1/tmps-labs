package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"tmps-go-labs/lab3/internal/ai"
	"tmps-go-labs/lab3/internal/character"
	"tmps-go-labs/lab3/internal/config"
	"tmps-go-labs/lab3/internal/facade"
)

// Application states
type state int

const (
	mainMenuState state = iota
	characterCreationState
	characterDisplayState
	chatState
	configState
)

// UI configuration constants
const (
	defaultWidth = 80  // Default terminal width for text wrapping
	maxWidth     = 120 // Maximum width before forced wrapping
)

// Styles with text wrapping support
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Width(defaultWidth).
			Align(lipgloss.Center)

	menuStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Width(defaultWidth)

	characterStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F25D94")).
			Width(defaultWidth)

	chatStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFA500")).
			Width(defaultWidth)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Width(defaultWidth)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(true).
			Width(defaultWidth).
			Align(lipgloss.Center)
)

type model struct {
	state         state
	facade        *facade.CharacterCreationFacade
	aiAdapter     *ai.CharacterAIAdapter
	config        *config.GameConfig
	currentChar   character.Character
	chatSession   ai.ChatSession // Persistent chat session with AI
	form          *huh.Form
	message       string
	error         string
	chatMessages  []string
	chatInput     string
	quitting      bool
	loading       bool   // General loading state
	loadingMsg    string // Loading message to display
	aiLoading     bool   // Specifically for AI operations
	terminalWidth int    // Terminal width for responsive design

	// UI Components
	loadingSpinner spinner.Model  // Proper spinner component
	chatViewport   viewport.Model // Scrollable chat viewport
	chatOffset     int            // Manual scroll offset for chat
}

func main() {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))

	// Initialize viewport for chat
	vp := viewport.New(defaultWidth-4, 10) // Initial size, will be updated
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFA500"))

	m := model{
		state:          mainMenuState,
		facade:         facade.NewCharacterCreationFacade(),
		aiAdapter:      ai.NewCharacterAIAdapter(),
		config:         config.GetGameConfig(),
		chatMessages:   make([]string, 0),
		terminalWidth:  defaultWidth,
		loadingSpinner: s,
		chatViewport:   vp,
		chatOffset:     0,
	}

	// Load saved character and API key on startup
	m.loadSavedDataOnStartup()

	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.loadingSpinner.Tick,
		tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

// tickMsg is used for spinner animation updates
type tickMsg struct{}

// characterCreatedMsg is sent when character creation completes
type characterCreatedMsg struct {
	character character.Character
	err       error
}

// chatResponseMsg is sent when AI chat response is received
type chatResponseMsg struct {
	response string
	err      error
}

// navigationCompleteMsg is sent when navigation is complete
type navigationCompleteMsg struct {
	targetState state
	err         error
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Continue spinner animation if loading or AI is working
		if m.loading || m.aiLoading {
			return m, tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd
	case characterCreatedMsg:
		// Handle async character creation completion (not currently used)
		m.loading = false
		if msg.err != nil {
			m.error = msg.err.Error()
			m.state = mainMenuState
		} else {
			m.currentChar = msg.character
			m.state = characterDisplayState
			m.message = "Character created successfully! 🎉"
		}
		return m, nil
	case chatResponseMsg:
		// Handle async chat response
		m.aiLoading = false

		if msg.err != nil {
			m.chatMessages = append(m.chatMessages, fmt.Sprintf("❌ Error: %s", msg.err.Error()))
			// If session fails, try to restart it
			if strings.Contains(msg.err.Error(), "session") {
				if closeErr := m.chatSession.Close(); closeErr != nil {
					log.Printf("Warning: Failed to close failed chat session: %v", closeErr)
				}
				m.chatSession = nil
				m.chatMessages = append(m.chatMessages, "🔄 Session lost. Please restart chat from main menu.")
			}
		} else {
			charName := m.currentChar.GetName()
			m.chatMessages = append(m.chatMessages, fmt.Sprintf("🎭 %s: %s", charName, msg.response))
		}

		// Keep last 20 messages for display (increased for better history)
		if len(m.chatMessages) > 20 {
			m.chatMessages = m.chatMessages[len(m.chatMessages)-20:]
		}

		// Reset chat scroll to bottom when new message arrives
		m.chatOffset = 0

		return m, nil
	case navigationCompleteMsg:
		// Handle navigation completion
		m.loading = false
		if msg.err != nil {
			m.error = msg.err.Error()
			m.state = mainMenuState
		} else {
			m.state = msg.targetState
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.updateTerminalSize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == mainMenuState {
				m.quitting = true
				// Clean up chat session before quitting
				if m.chatSession != nil {
					if err := m.chatSession.Close(); err != nil {
						// Log error but don't block quitting
						log.Printf("Warning: Failed to close chat session: %v", err)
					}
					m.chatSession = nil
				}
				return m, tea.Quit
			}
			// Clean up chat session when leaving chat state
			if m.state == chatState && m.chatSession != nil {
				if err := m.chatSession.Close(); err != nil {
					log.Printf("Warning: Failed to close chat session: %v", err)
				}
				m.chatSession = nil
			}
			m.state = mainMenuState
			m.error = ""
			return m, nil
		case "1":
			if m.state == mainMenuState {
				return m.startCharacterCreation()
			}
		case "2":
			if m.state == mainMenuState && m.currentChar != nil {
				m.loading = true
				m.loadingMsg = "Loading character details..."
				return m, m.navigateToCharacterDisplayCmd()
			}
		case "3":
			if m.state == mainMenuState && m.currentChar != nil {
				m.loading = true
				m.loadingMsg = "Starting chat session..."
				return m, m.startChatSessionAsyncCmd()
			}
		case "4":
			if m.state == mainMenuState {
				return m.startConfig()
			}
		case "enter":
			if m.state == chatState && m.chatInput != "" {
				return m.sendChatMessage()
			}
		case "backspace":
			if m.state == chatState && len(m.chatInput) > 0 {
				m.chatInput = m.chatInput[:len(m.chatInput)-1]
			}
		case "up", "k":
			if m.state == chatState {
				return m.scrollChatUp()
			}
		case "down", "j":
			if m.state == chatState {
				return m.scrollChatDown()
			}
		case "pgup":
			if m.state == chatState {
				return m.scrollChatPageUp()
			}
		case "pgdown":
			if m.state == chatState {
				return m.scrollChatPageDown()
			}
		default:
			if m.state == chatState {
				m.chatInput += msg.String()
			}
		}
	}

	if m.form != nil {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if m.form.State == huh.StateCompleted {
				return m.handleFormCompletion()
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	if m.quitting {
		return "Thanks for using Character Builder! 👋\n"
	}

	// Update styles with current terminal width
	m.updateStylesForWidth()

	var content string

	// Show loading spinner if we're in a loading state
	if m.loading {
		spinnerView := m.loadingSpinner.View()
		loadingText := fmt.Sprintf("%s %s", spinnerView, m.loadingMsg)
		content = loadingStyle.Render(loadingText)
	} else {
		switch m.state {
		case mainMenuState:
			content = m.renderMainMenu()
		case characterCreationState:
			if m.form != nil {
				content = m.form.View()
			}
		case characterDisplayState:
			content = m.renderCharacterDisplay()
		case chatState:
			content = m.renderChat()
		case configState:
			if m.form != nil {
				content = m.form.View()
			}
		}
	}

	title := titleStyle.Render("🎲 D&D Character Builder 🎲")

	// Handle error display with wrapping
	if m.error != "" {
		wrappedError := wrapText("❌ "+m.error, m.terminalWidth-4)
		errorMsg := errorStyle.Render(wrappedError)
		return fmt.Sprintf("%s\n\n%s\n\n%s", title, errorMsg, content)
	}

	// Handle success message display
	if m.message != "" {
		wrappedMessage := wrapText("✅ "+m.message, m.terminalWidth-4)
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Width(m.terminalWidth - 4)
		successMsg := messageStyle.Render(wrappedMessage)
		return fmt.Sprintf("%s\n\n%s\n\n%s", title, successMsg, content)
	}

	return fmt.Sprintf("%s\n\n%s", title, content)
}

func (m *model) renderMainMenu() string {
	options := []string{
		"1️⃣  Create New Character",
		"2️⃣  View Current Character" + m.getCharacterStatus(),
		"3️⃣  Chat with Character" + m.getChatStatus(),
		"4️⃣  Configure API Key",
		"",
		"Press number to select • q to quit",
	}
	return menuStyle.Render(strings.Join(options, "\n"))
}

func (m *model) getCharacterStatus() string {
	if m.currentChar == nil {
		return " (None)"
	}
	return fmt.Sprintf(" (%s)", m.currentChar.GetName())
}

func (m *model) getChatStatus() string {
	if m.currentChar == nil {
		return " (Create character first)"
	}
	if !m.aiAdapter.IsConfigured() {
		return " (Configure API key first)"
	}
	return ""
}

func (m *model) renderCharacterDisplay() string {
	if m.currentChar == nil {
		return "No character created yet."
	}

	stats := m.currentChar.GetStats()
	contentWidth := m.terminalWidth - 8 // Account for padding and borders

	lines := []string{
		fmt.Sprintf("🧙 Name: %s", m.currentChar.GetName()),
		fmt.Sprintf("🏃 Race: %s", m.currentChar.GetRace()),
		fmt.Sprintf("⚔️  Class: %s", m.currentChar.GetClass()),
		fmt.Sprintf("📜 Background: %s", m.currentChar.GetBackground()),
		"",
		"📊 Stats:",
		fmt.Sprintf("  💪 STR: %d  🏃 DEX: %d  ❤️  CON: %d", stats.Strength, stats.Dexterity, stats.Constitution),
		fmt.Sprintf("  🧠 INT: %d  🦉 WIS: %d  😊 CHA: %d", stats.Intelligence, stats.Wisdom, stats.Charisma),
		"",
	}

	// Wrap longer content sections
	abilities := "✨ Abilities: " + strings.Join(m.currentChar.GetAbilities(), ", ")
	lines = append(lines, wrapText(abilities, contentWidth))

	equipment := "🎒 Equipment: " + strings.Join(m.currentChar.GetEquipment(), ", ")
	lines = append(lines, wrapText(equipment, contentWidth))

	personality := "🎭 Personality: " + strings.Join(m.currentChar.GetPersonality(), ", ")
	lines = append(lines, wrapText(personality, contentWidth))

	lines = append(lines, "")

	description := "📖 " + m.currentChar.Describe()
	lines = append(lines, wrapText(description, contentWidth))

	lines = append(lines, "")
	lines = append(lines, "Press 'q' to go to previous screen")

	return characterStyle.Render(strings.Join(lines, "\n"))
}

func (m *model) renderChat() string {
	if m.currentChar == nil {
		return "No character to chat with."
	}

	contentWidth := m.terminalWidth - 8 // Account for padding and borders

	// Header
	header := fmt.Sprintf("💬 Chatting with %s", m.currentChar.GetName())

	// Show all chat messages without complex viewport logic
	var chatLines []string

	// Display all messages (simplified approach)
	if len(m.chatMessages) == 0 {
		chatLines = append(chatLines, "No messages yet. Start the conversation!")
	} else {
		// Show last N messages based on available space
		maxVisible := 10
		startIdx := 0
		if len(m.chatMessages) > maxVisible {
			startIdx = len(m.chatMessages) - maxVisible + m.chatOffset
			if startIdx < 0 {
				startIdx = 0
			}
			if startIdx > len(m.chatMessages)-maxVisible {
				startIdx = len(m.chatMessages) - maxVisible
			}
		}

		endIdx := startIdx + maxVisible
		if endIdx > len(m.chatMessages) {
			endIdx = len(m.chatMessages)
		}

		// Add scroll indicator at top if there are more messages above
		if startIdx > 0 {
			chatLines = append(chatLines, fmt.Sprintf("  ▲ %d more messages above (use ↑ to scroll) ▲", startIdx))
			chatLines = append(chatLines, "")
		}

		// Add visible messages
		for i := startIdx; i < endIdx; i++ {
			wrappedMsg := wrapText(m.chatMessages[i], contentWidth-2)
			chatLines = append(chatLines, wrappedMsg)
		}

		// Add scroll indicator at bottom if there are more messages below
		remainingBelow := len(m.chatMessages) - endIdx
		if remainingBelow > 0 {
			chatLines = append(chatLines, "")
			chatLines = append(chatLines, fmt.Sprintf("  ▼ %d more messages below (use ↓ to scroll) ▼", remainingBelow))
		}
	}

	// Show AI loading indicator if waiting for response
	if m.aiLoading {
		spinnerView := m.loadingSpinner.View()
		chatLines = append(chatLines, "")
		chatLines = append(chatLines, fmt.Sprintf("🎭 %s is thinking... %s", m.currentChar.GetName(), spinnerView))
	}

	// Create the chat display area
	chatDisplay := strings.Join(chatLines, "\n")

	// User input section
	userInput := fmt.Sprintf("You: %s_", m.chatInput)
	if len(userInput) > contentWidth {
		userInput = wrapText(userInput, contentWidth)
	}

	instructions := "Type your message and press Enter • ↑↓/k,j to scroll • PageUp/PageDown for fast scroll • q to return"

	// Combine all parts
	content := []string{
		header,
		"",
		chatDisplay,
		"",
		"─────────────────────────", // Separator line
		userInput,
		"",
		instructions,
	}

	return chatStyle.Render(strings.Join(content, "\n"))
}

func (m *model) startCharacterCreation() (tea.Model, tea.Cmd) {
	m.state = characterCreationState
	m.error = ""

	races, classes, backgrounds := m.facade.GetAvailableOptions()

	basicGroup := huh.NewGroup(
		huh.NewInput().
			Title("Character Name").
			Placeholder("Enter character name...").
			Key("name"),
		huh.NewSelect[string]().
			Title("Race").
			Options(huh.NewOptions(races...)...).
			Key("race"),
		huh.NewSelect[string]().
			Title("Class").
			Options(huh.NewOptions(classes...)...).
			Key("class"),
		huh.NewSelect[string]().
			Title("Background").
			Options(huh.NewOptions(backgrounds...)...).
			Key("background"),
		huh.NewSelect[string]().
			Title("Creation Method").
			Options(
				huh.NewOption("Standard Build (15,14,13,12,10,8 stats)", character.CreationMethodStandard),
				huh.NewOption("Random Build (4d6 drop lowest)", character.CreationMethodRandom),
				huh.NewOption("Custom Build (choose your own stats)", character.CreationMethodCustom),
			).
			Key(character.FormKeyMethod),
	).Title("🎲 Character Basics")

	statsGroup := huh.NewGroup(
		huh.NewInput().Title("Strength (3-18)").Placeholder("15").Key(character.FormKeyStrength),
		huh.NewInput().Title("Dexterity (3-18)").Placeholder("14").Key(character.FormKeyDexterity),
		huh.NewInput().Title("Constitution (3-18)").Placeholder("13").Key(character.FormKeyConstitution),
		huh.NewInput().Title("Intelligence (3-18)").Placeholder("12").Key(character.FormKeyIntelligence),
		huh.NewInput().Title("Wisdom (3-18)").Placeholder("10").Key(character.FormKeyWisdom),
		huh.NewInput().Title("Charisma (3-18)").Placeholder("8").Key(character.FormKeyCharisma),
	).Title("📊 Custom Stats").WithHideFunc(func() bool {
		return m.form != nil && m.form.GetString(character.FormKeyMethod) != character.CreationMethodCustom
	})

	m.form = huh.NewForm(basicGroup, statsGroup)
	return m, m.form.Init()
}

func (m *model) startConfig() (tea.Model, tea.Cmd) {
	m.state = configState
	m.error = ""

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Gemini API Key").
				Placeholder("Enter your Gemini API key...").
				EchoMode(huh.EchoModePassword).
				Key("apikey"),
		),
	)
	return m, m.form.Init()
}

func (m *model) handleFormCompletion() (tea.Model, tea.Cmd) {
	switch m.state {
	case characterCreationState:
		return m.createCharacter()
	case configState:
		return m.saveConfig()
	}
	return m, nil
}

func (m *model) createCharacter() (tea.Model, tea.Cmd) {
	name := m.form.GetString("name")
	race := m.form.GetString("race")
	class := m.form.GetString("class")
	background := m.form.GetString("background")
	method := m.form.GetString("method")

	if name == "" {
		m.error = "Character name is required"
		m.state = mainMenuState
		return m, nil
	}

	// Show loading state
	m.loading = true
	m.loadingMsg = "Creating your character..."
	m.error = ""
	m.message = ""

	var char character.Character
	var err error

	switch method {
	case character.CreationMethodRandom:
		char, err = m.facade.GenerateRandomCharacter(name)
	case character.CreationMethodCustom:
		stats := character.Stats{}
		if stats.Strength, err = m.parseStatInput(m.form.GetString(character.FormKeyStrength), character.StatStrength); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		if stats.Dexterity, err = m.parseStatInput(m.form.GetString(character.FormKeyDexterity), character.StatDexterity); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		if stats.Constitution, err = m.parseStatInput(m.form.GetString(character.FormKeyConstitution), character.StatConstitution); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		if stats.Intelligence, err = m.parseStatInput(m.form.GetString(character.FormKeyIntelligence), character.StatIntelligence); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		if stats.Wisdom, err = m.parseStatInput(m.form.GetString(character.FormKeyWisdom), character.StatWisdom); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		if stats.Charisma, err = m.parseStatInput(m.form.GetString(character.FormKeyCharisma), character.StatCharisma); err != nil {
			m.error = err.Error()
			m.loading = false
			m.state = mainMenuState
			return m, nil
		}
		char, err = m.facade.CreateAdvancedCharacter(name, race, class, background, stats, []string{}, []string{})
	default: // standard (character.CreationMethodStandard)
		char, err = m.facade.CreateBasicCharacter(name, race, class, background)
	}

	if err != nil {
		m.error = err.Error()
		m.loading = false
		m.state = mainMenuState
		return m, nil
	}

	m.currentChar = char
	m.loading = false

	// Hydrate character with AI content and auto-save
	if saveErr := m.hydrateAndSaveCharacter(char); saveErr != nil {
		m.error = fmt.Sprintf("Character created but save failed: %v", saveErr)
	} else {
		m.message = "Character created and saved successfully! 🎉💾"
	}

	m.state = characterDisplayState
	return m, nil
}

// hydrateAndSaveCharacter triggers AI content generation and saves the character
func (m *model) hydrateAndSaveCharacter(char character.Character) error {
	log.Printf("🧪 HYDRATING CHARACTER: %s", char.GetName())

	// Trigger AI hydration - the characterToPrompt method will auto-generate content
	// We can trigger this by attempting to chat, which calls characterToPrompt internally
	if m.aiAdapter.IsConfigured() {
		log.Printf("🔑 API KEY CONFIGURED: Generating AI content for %s", char.GetName())

		// Generate backstory if empty
		if char.GetBackstory() == "" {
			log.Printf("📝 GENERATING BACKSTORY: Character has no backstory yet")
			if backstory, err := m.aiAdapter.GenerateCharacterBackstory(char); err == nil && backstory != "" {
				char.SetBackstory(backstory)
				log.Printf("✅ BACKSTORY GENERATED: %d characters", len(backstory))
			} else {
				log.Printf("❌ BACKSTORY GENERATION FAILED: %v", err)
			}
		} else {
			log.Printf("💾 BACKSTORY EXISTS: Skipping generation (length: %d)", len(char.GetBackstory()))
		}

		// Trigger system prompt generation by doing a dummy chat
		// This will call characterToPrompt internally and hydrate the character
		log.Printf("🎭 TRIGGERING SYSTEM PROMPT: Dummy chat to hydrate character")
		_, _ = m.aiAdapter.ChatWithCharacter(char, "Hello") // Dummy call to trigger hydration
	} else {
		log.Printf("⚠️  NO API KEY: Skipping AI hydration for %s", char.GetName())
	}

	// Save to persistent storage (single character slot)
	persistentConfig := m.config.GetPersistentConfig()
	characterID := "main-character" // Single slot

	log.Printf("💾 SAVING CHARACTER: %s to persistent config", char.GetName())
	err := persistentConfig.SaveCharacter(char, characterID)
	if err != nil {
		log.Printf("❌ SAVE FAILED: %v", err)
	} else {
		log.Printf("✅ SAVE SUCCESS: Character %s saved with backstory (%d chars) and system prompt (%d chars)",
			char.GetName(), len(char.GetBackstory()), len(char.GetSystemPrompt()))
	}

	return err
}

func (m *model) saveConfig() (tea.Model, tea.Cmd) {
	apiKey := m.form.GetString("apikey")
	if apiKey == "" {
		m.error = "API key cannot be empty"
		m.state = mainMenuState
		return m, nil
	}

	// Set API key in memory
	m.aiAdapter.SetAPIKey(apiKey)

	// Persist API key to disk
	persistentConfig := m.config.GetPersistentConfig()
	persistentConfig.GeminiAPIKey = apiKey
	if err := persistentConfig.SaveConfig(); err != nil {
		m.error = fmt.Sprintf("API key set but failed to save: %v", err)
	} else {
		m.message = "API key configured and saved successfully! 🔑💾"
	}

	m.state = mainMenuState
	return m, nil
}

func (m *model) startChatSession() (tea.Model, tea.Cmd) {
	if !m.aiAdapter.IsConfigured() {
		m.error = "Please configure your Gemini API key first"
		m.state = mainMenuState
		return m, nil
	}

	if !m.aiAdapter.SupportsChatSessions() {
		m.error = "Chat sessions not supported with current AI client"
		m.state = mainMenuState
		return m, nil
	}

	// Close existing session if any
	if m.chatSession != nil {
		if err := m.chatSession.Close(); err != nil {
			log.Printf("Warning: Failed to close existing chat session: %v", err)
		}
		m.chatSession = nil
	}

	// Start new chat session with character context
	session, err := m.aiAdapter.StartCharacterChatSession(m.currentChar)
	if err != nil {
		m.error = fmt.Sprintf("Failed to start chat session: %v", err)
		m.state = mainMenuState
		return m, nil
	}

	m.chatSession = session
	m.state = chatState
	m.chatMessages = []string{
		"💬 Chat session started! Your character remembers the conversation.",
		fmt.Sprintf("🎭 You're now talking to %s. Type your message and press Enter.", m.currentChar.GetName()),
	}
	m.message = ""
	m.error = ""

	return m, nil
}

func (m *model) sendChatMessage() (tea.Model, tea.Cmd) {
	// Ensure we have a chat session
	if m.chatSession == nil {
		m.error = "No active chat session. Please restart the chat."
		m.state = mainMenuState
		return m, nil
	}

	userMsg := m.chatInput
	m.chatInput = ""

	m.chatMessages = append(m.chatMessages, fmt.Sprintf("🧑 You: %s", userMsg))

	// Show AI loading state and start spinner
	m.aiLoading = true

	// Return command to handle chat asynchronously
	return m, m.sendChatMessageCmd(userMsg)
}

// sendChatMessageCmd creates an async command for sending chat messages
func (m *model) sendChatMessageCmd(userMsg string) tea.Cmd {
	return func() tea.Msg {
		response, err := m.aiAdapter.ChatWithCharacterSession(m.chatSession, userMsg)
		return chatResponseMsg{
			response: response,
			err:      err,
		}
	}
}

func (m *model) parseStatInput(input string, statName string) (int, error) {
	if input == "" {
		return 0, fmt.Errorf("%s cannot be empty", statName)
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", statName)
	}

	if value < 3 || value > 18 {
		return 0, fmt.Errorf("%s must be between 3 and 18", statName)
	}

	return value, nil
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		width = defaultWidth
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if len(currentLine)+len(word)+1 > width && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// updateTerminalSize updates the terminal width for responsive design
func (m *model) updateTerminalSize(msg tea.WindowSizeMsg) {
	m.terminalWidth = msg.Width
	if m.terminalWidth > maxWidth {
		m.terminalWidth = maxWidth
	}
	if m.terminalWidth < 40 {
		m.terminalWidth = 40 // Minimum width for readability
	}

	// Update viewport size for chat
	viewportWidth := m.terminalWidth - 8
	viewportHeight := 8 // Fixed height for chat area
	if viewportWidth < 30 {
		viewportWidth = 30
	}

	m.chatViewport.Width = viewportWidth
	m.chatViewport.Height = viewportHeight
}

// updateStylesForWidth updates all styles to use current terminal width
func (m *model) updateStylesForWidth() {
	width := m.terminalWidth - 4 // Account for padding and borders
	if width < 30 {
		width = 30
	}

	titleStyle = titleStyle.Width(width)
	menuStyle = menuStyle.Width(width)
	characterStyle = characterStyle.Width(width)
	chatStyle = chatStyle.Width(width)
	errorStyle = errorStyle.Width(width)
	loadingStyle = loadingStyle.Width(width)
}

// loadSavedDataOnStartup loads saved character and API key when the app starts
func (m *model) loadSavedDataOnStartup() {
	persistentConfig := m.config.GetPersistentConfig()

	// Load saved API key
	if persistentConfig.GeminiAPIKey != "" {
		m.aiAdapter.SetAPIKey(persistentConfig.GeminiAPIKey)
	}

	// Load the last used character if available
	lastCharacterID := persistentConfig.GetLastCharacterID()
	if lastCharacterID != "" {
		if char, err := persistentConfig.LoadCharacter(lastCharacterID); err == nil {
			m.currentChar = char
			log.Printf("📂 CHARACTER LOADED: %s from persistent config", char.GetName())
			log.Printf("   └─ Backstory: %t (length: %d)", char.GetBackstory() != "", len(char.GetBackstory()))
			log.Printf("   └─ SystemPrompt: %t (length: %d)", char.GetSystemPrompt() != "", len(char.GetSystemPrompt()))
			m.message = fmt.Sprintf("Welcome back! Loaded character: %s 🎭", char.GetName())
		} else {
			log.Printf("❌ CHARACTER LOAD FAILED: %v", err)
		}
	} else {
		log.Printf("🆕 NO SAVED CHARACTER: Starting fresh")
	}
}

// scrollChatUp scrolls the chat up by one line
func (m *model) scrollChatUp() (tea.Model, tea.Cmd) {
	if len(m.chatMessages) > 8 { // Only scroll if there are more messages than visible
		if m.chatOffset < len(m.chatMessages)-8 {
			m.chatOffset++
		}
	}
	return m, nil
}

// scrollChatDown scrolls the chat down by one line
func (m *model) scrollChatDown() (tea.Model, tea.Cmd) {
	if m.chatOffset > 0 {
		m.chatOffset--
	}
	return m, nil
}

// scrollChatPageUp scrolls the chat up by a page
func (m *model) scrollChatPageUp() (tea.Model, tea.Cmd) {
	pageSize := 5
	if len(m.chatMessages) > 8 {
		m.chatOffset += pageSize
		maxOffset := len(m.chatMessages) - 8
		if m.chatOffset > maxOffset {
			m.chatOffset = maxOffset
		}
	}
	return m, nil
}

// scrollChatPageDown scrolls the chat down by a page
func (m *model) scrollChatPageDown() (tea.Model, tea.Cmd) {
	pageSize := 5
	m.chatOffset -= pageSize
	if m.chatOffset < 0 {
		m.chatOffset = 0
	}
	return m, nil
}

// navigateToCharacterDisplayCmd creates async command for character display navigation
func (m *model) navigateToCharacterDisplayCmd() tea.Cmd {
	return func() tea.Msg {
		// Simulate some loading time for demonstration
		time.Sleep(200 * time.Millisecond)
		return navigationCompleteMsg{
			targetState: characterDisplayState,
			err:         nil,
		}
	}
}

// navigateToChatCmd creates async command for chat navigation
func (m *model) navigateToChatCmd() tea.Cmd {
	return func() tea.Msg {
		// This is where chat session setup happens - this might be slow
		time.Sleep(300 * time.Millisecond) // Simulate the slow operation
		return navigationCompleteMsg{
			targetState: chatState,
			err:         nil,
		}
	}
}

// startChatSessionAsyncCmd creates async command for starting chat session
func (m *model) startChatSessionAsyncCmd() tea.Cmd {
	return func() tea.Msg {
		log.Printf("🚀 STARTING CHAT SESSION: For %s", m.currentChar.GetName())

		// Perform the actual chat session start (this is the slow part)
		if !m.aiAdapter.IsConfigured() {
			log.Printf("❌ CHAT SESSION FAILED: No API key configured")
			return navigationCompleteMsg{
				targetState: mainMenuState,
				err:         fmt.Errorf("please configure your Gemini API key first"),
			}
		}

		if !m.aiAdapter.SupportsChatSessions() {
			log.Printf("❌ CHAT SESSION FAILED: Chat sessions not supported")
			return navigationCompleteMsg{
				targetState: mainMenuState,
				err:         fmt.Errorf("chat sessions not supported with current AI client"),
			}
		}

		// Close existing session if any
		if m.chatSession != nil {
			log.Printf("🔄 CLOSING EXISTING SESSION: Cleaning up previous chat")
			if err := m.chatSession.Close(); err != nil {
				log.Printf("Warning: Failed to close existing chat session: %v", err)
			}
			m.chatSession = nil
		}

		// Start new chat session with character context (this is the slow part!)
		log.Printf("🎭 INITIALIZING CHARACTER CONTEXT: This will call characterToPrompt()")
		sessionStartTime := time.Now()
		session, err := m.aiAdapter.StartCharacterChatSession(m.currentChar)
		sessionDuration := time.Since(sessionStartTime)
		if err != nil {
			log.Printf("❌ CHAT SESSION START FAILED after %v: %v", sessionDuration, err)
			return navigationCompleteMsg{
				targetState: mainMenuState,
				err:         fmt.Errorf("failed to start chat session: %v", err),
			}
		}
		log.Printf("✅ CHAT SESSION COMPLETE: Total time %v", sessionDuration)

		// Store the session
		m.chatSession = session

		// Initialize chat messages
		m.chatMessages = []string{
			"💬 Chat session started! Your character remembers the conversation.",
			fmt.Sprintf("🎭 You're now talking to %s. Type your message and press Enter.", m.currentChar.GetName()),
		}
		m.chatOffset = 0

		return navigationCompleteMsg{
			targetState: chatState,
			err:         nil,
		}
	}
}
