package core

import (
	"fmt"
	"log"
	"strings"
    "os/exec"
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
)

// ResponseItem represents a single response from the AI.
type ResponseItem struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

// Sends a message to the chat and processes the response.
func (kai *Kai) Reason(userInput string) (string, error) {
	// Send a message to the chat
    resp, err := kai.Chat.SendMessage(kai.Context, genai.Text(userInput))
    if err != nil {
        return "", fmt.Errorf("error sending message: %v", err)
    }
	// Save the history after sending the message
	kai.SaveHistory()
	// Parse the response
	candidates := resp.Candidates
	if len(candidates) <= 0 || len(candidates[0].Content.Parts) <= 0 {
		return "", fmt.Errorf("no content generated")
	}
	responseJSON := partToString(candidates[0].Content.Parts[0])
	// Process the JSON response
	kai.processJSONResponse(responseJSON)
	return responseJSON, nil
}

// Executes a given shell command and returns the output.
func (kai *Kai) executeCommand(command string) (string, error) {
    cmd := exec.Command("sh", "-c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("failed to execute command: %w", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// Process JSON response
func (kai *Kai) processJSONResponse(jsonStr string) {
	// Unmarshal the JSON string into a slice of ResponseItem
	var responseItems []ResponseItem
	err := json.Unmarshal([]byte(jsonStr), &responseItems)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	// Iterate through the slice and process each item
	for _, item := range responseItems {
		switch item.Type {
		case "script":
			fmt.Println("Speaking:", item.Data)
			// Have Kai speak the script
            if err := kai.Speak(item.Data); err != nil {
                log.Printf("Failed to speak: %v", err)
            }
		case "command":
			fmt.Println("Executing command:", item.Data)
			output, err := kai.executeCommand(item.Data)
			if err != nil {
				fmt.Println("Command execution failed:", err)
				// TODO: Recurse
				return
			} 
			fmt.Println("Command output:", output)
		default:
			fmt.Println("Unknown type:", item.Type)
		}
	}
}

// Helper function to convert genai.Part to a string
func partToString(part genai.Part) string {
    if v, ok := part.(genai.Text); ok {
        return string(v)
    }
	return ""
}