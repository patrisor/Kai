package core

import (
	"os"
	"fmt"
	"encoding/json"
)

// Primer struct to hold individual primer data
type Prompts struct {
    Primers   map[string]string `json:"primers"`
}

// Method reads all prompts from the specified file into a Prompts struct
func LoadPrompts(file string) (*Prompts, error) {
    var prompts Prompts
    // Open the prompts file
    promptsFile, err := os.Open(file)
    if err != nil {
        return nil, fmt.Errorf("failed to open prompts file: %v", err)
    }
    defer promptsFile.Close()
    // Decode the JSON data into the Prompts struct
    decoder := json.NewDecoder(promptsFile)
    if err := decoder.Decode(&prompts); err != nil {
        return nil, fmt.Errorf("failed to decode prompts file: %v", err)
    }
    return &prompts, nil
}
