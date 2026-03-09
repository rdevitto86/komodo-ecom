package config

import (
	"os"
	"sync"
)

type Config struct {
	mu   sync.RWMutex
	data map[string]string
}

var (
	once     sync.Once
	instance *Config
)

func init() {
	once.Do(func() {
		instance = &Config{
			data: make(map[string]string),
		}
	})
}

// Checks local in-memory config first, then falls back to environment variable
func GetConfigValue(key string) string {
	if key == "" || instance == nil { return "" }
	instance.mu.RLock()
	val := instance.data[key]
	instance.mu.RUnlock()
	if val != "" { return val }
	return os.Getenv(key)
}

// Sets value in local in-memory config
func SetConfigValue(key, value string) {
	if value == "" || key == "" || instance == nil { return }
	instance.mu.Lock()
	instance.data[key] = value
	instance.mu.Unlock()
}

// Removes value from local in-memory config only
func DeleteConfigValue(key string) {
	if key == "" || instance == nil { return }
	instance.mu.Lock()
	delete(instance.data, key)
	instance.mu.Unlock()
}

// Returns all keys currently stored in the config
func GetAllKeys() []string {
	if instance == nil { return []string{} }

	instance.mu.RLock()
	defer instance.mu.RUnlock()

	keys := make([]string, 0, len(instance.data))
	for k := range instance.data {
		keys = append(keys, k)
	}
	return keys
}
