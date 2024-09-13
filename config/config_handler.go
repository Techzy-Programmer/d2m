package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/univ"
)

var configFile string;

func init() {
	configPath, err := univ.GetUserConfigPath("d2m");
	if err != nil {
		log.Fatalf("Failed to get config path: %v", err)
	}

	configFile = filepath.Join(configPath, "_config.json")
	ensureConfigExists()
}

// GetData retrieves data from the config file and returns it as the specified type.
func GetData[T any](key string, def ...T) T {
	var defValue T

	// If a default value is provided, use it
	if len(def) > 0 {
		defValue = def[0]
	}

	strData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(strData, &data); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Type assertion to convert `interface{}` to the desired type `T`
	value, _ := data[key].(T)

	if reflect.ValueOf(value).IsZero() {
		return defValue
	}
	return value
}

// SetData sets data in the config file with the specified key and value.
func SetData[T any](key string, value T) {
	// Read existing config data
	strData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(strData, &data); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Set the new value
	data[key] = value

	// Marshal updated data back to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal config data: %v", err)
	}

	// Write updated data to the file
	err = os.WriteFile(configFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write config data: %v", err)
	}
}

// ensureConfigExists checks if the config file exists and creates it if not.
func ensureConfigExists() {
	configDir := filepath.Dir(configFile)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create config directory: %v", err)
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		paint.InfoF("Config file not found, creating \"%s\"...", configFile)

		if err := os.WriteFile(configFile, []byte(`{}`), 0644); err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
	}
}
