package config

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/Techzy-Programmer/d2m/config/paint"
)

var configPath = "_config.json"

func init() {
	ensureConfigExists()
}

// GetData retrieves data from the config file and returns it as the specified type.
func GetData[T any](key string, def ...T) T {
	var defValue T

	// If a default value is provided, use it
	if len(def) > 0 {
		defValue = def[0]
	}

	strData, err := os.ReadFile(configPath)
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
	strData, err := os.ReadFile(configPath)
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
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write config data: %v", err)
	}
}

// ensureConfigExists checks if the config file exists and creates it if not.
func ensureConfigExists() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		paint.Info("Config file not found, creating a new one...")

		if err := os.WriteFile(configPath, []byte(`{}`), 0644); err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
	}
}
