package config

import (
	"sync"
)

// GameConfig holds application configuration (Singleton) and manages persistence
type GameConfig struct {
	persistentConfig *PersistentConfig
	mu               sync.RWMutex
}

var (
	instance *GameConfig
	once     sync.Once
)

// GetGameConfig returns the singleton instance of GameConfig
func GetGameConfig() *GameConfig {
	once.Do(func() {
		instance = &GameConfig{}
		instance.loadConfig()
	})
	return instance
}

// loadConfig loads the persistent configuration
func (gc *GameConfig) loadConfig() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	config, err := LoadConfig()
	if err != nil {
		// If loading fails, create default config
		config = &PersistentConfig{
			AppName:         "Character Builder",
			Version:         "1.0.0",
			SavedCharacters: make([]SavedCharacter, 0),
		}
	}
	gc.persistentConfig = config
}

// GetPersistentConfig returns the persistent config for direct access
func (gc *GameConfig) GetPersistentConfig() *PersistentConfig {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.persistentConfig
}

// SetGeminiAPIKey sets the API key for Gemini and persists it
func (gc *GameConfig) SetGeminiAPIKey(key string) error {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.persistentConfig.GeminiAPIKey = key
	return gc.persistentConfig.SaveConfig()
}

// GetGeminiAPIKey returns the Gemini API key
func (gc *GameConfig) GetGeminiAPIKey() string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.persistentConfig.GeminiAPIKey
}

// GetAppName returns the application name
func (gc *GameConfig) GetAppName() string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.persistentConfig.AppName
}

// GetVersion returns the application version
func (gc *GameConfig) GetVersion() string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.persistentConfig.Version
}
