package core

import (
	"os"
	"fmt"
	"encoding/json"
	// Local utilities
	"kai/source/utils"
	// Gemini API
	"github.com/google/generative-ai-go/genai"
)

// Method primes the AI with the provided primer and history file.
func (kai *Kai) PrimeAI(primer, historyFile string) {
	// Initialize the chat session
	kai.Chat = kai.Model.StartChat()
	// Check if a history file is provided
	if historyFile != "" {
		// Load history if it exists, else use the primer
		if _, err := os.Stat(historyFile); err == nil {
			data, err := os.ReadFile(historyFile)
			if err != nil {
				fmt.Println("Failed to read chat history file:", err)
			} else {
				var rawHistory []map[string]interface{}
				if err := json.Unmarshal(data, &rawHistory); err != nil {
					fmt.Println("Failed to unmarshal chat history:", err)
				} else {
					var history []*genai.Content
					for _, rawContent := range rawHistory {
						content := &genai.Content{
							Role: rawContent["Role"].(string),
						}
						parts, ok := rawContent["Parts"].([]interface{})
						if !ok {
							fmt.Println("Failed to convert parts to []interface{}")
							continue
						}
						for _, part := range parts {
							textPart, ok := part.(string)
							if ok {
								content.Parts = append(
									content.Parts, genai.Text(textPart),
								)
							}
						}
						history = append(history, content)
					}
					kai.Chat.History = history
				}
			}
		} else { // History file does not exist, using primer
			if primer != "" {
				// Append system information to the primer
				systemInfo := "\n\nSystem Information:\n" + utils.GetSystemInfo()
				primerWithInfo := primer + systemInfo
				// Prime the AI
				kai.Chat.History = []*genai.Content{
					{
						Parts: []genai.Part{
							genai.Text(primerWithInfo),
						},
						Role: "model",
					},
				}
			}
		}
	}
}