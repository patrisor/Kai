package core

import (
	"os"
	"encoding/json"
)

// Config structure to hold API key.
type Config struct {
	APIKey      string `json:"api_key"`
}

// SaveConfig writes the Config struct to the configuration file.
func SaveConfig(file string, config *Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, configData, 0644)
}

// LoadConfig reads a configuration file and returns a Config struct.
func LoadConfig(file string) (*Config, error) {
    var config Config
    // Open the configuration file
    configFile, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer configFile.Close()
    // Decode the JSON data into the Config struct
    decoder := json.NewDecoder(configFile)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}