package core

import (
	"fmt"
	"log"
	"strings"
	"os/exec"
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
)

// Process JSON response
func (kai *Kai) Respond(jsonStr string) {
	// Sanitize and unmarshal the JSON response
	responseItems, err := sanitizeAndUnmarshal(jsonStr)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return
	}
	// Iterate through the slice and process each item
	for _, item := range responseItems {
		switch item.Type {
		case "script":
			processScript(kai, item.Data)
		case "command":
			processCommand(kai, item.Data)
		default:
			fmt.Println("Unknown type:", item.Type)
		}
	}
}

// Process a script response item
func processScript(kai *Kai, script string) {
	fmt.Println("Speaking:", script)
	// Have Kai speak the script
	if err := kai.Speak(script); err != nil {
		log.Printf("Failed to speak: %v", err)
	}
}

// Process a command response item
func processCommand(kai *Kai, command string) {
	fmt.Println("Executing command:", command)
	// Sanitize command string
	sanitizedCommand := strings.ReplaceAll(command, ",", "")
	// Execute command
	output, err := kai.executeCommand(sanitizedCommand)

	// TODO: Testing
	// fmt.Printf("Error: %v\n", err)
	// fmt.Printf("Output: %s\n", output)

	if err != nil {
		errorMessage := fmt.Sprintf(
			"Command failed: %v. " + 
			"Please analyze the error and generate a new solution.",
			err,
		)

		// TODO: Testing
		fmt.Println(errorMessage)

		// TODO: Recurse

		// Feed the error back into the system to generate a new solution
		kai.appendToChatHistory("model", errorMessage)
		newResponse, aiErr := kai.Reason(errorMessage)
		if aiErr != nil {
			log.Fatalf("Failed to process new response: %v", aiErr)
		}
		kai.Respond(newResponse)
		// Exit the current iteration to prevent further processing
		return 
	} 
	// Append command ouput to chat history
	if output != "" {
		fmt.Println("Command output:", output)
		kai.appendToChatHistory("model", output)
		// Save the appended content to disk
		kai.SaveHistory()
	}
}

// Method removes unwanted characters from the input JSON string to ensure the 
// input JSON string is properly formatted for unmarshalling.
func sanitizeJSONString(jsonStr string) string {
	// Remove unwanted markers
	jsonStr = strings.ReplaceAll(jsonStr, "```json", "")
	jsonStr = strings.ReplaceAll(jsonStr, "```", "")
	// Trim leading and trailing spaces
	return strings.TrimSpace(jsonStr)
}

// Sanitize and unmarshal the JSON string
func sanitizeAndUnmarshal(jsonStr string) ([]ResponseItem, error) {

	// TODO: Testing
	log.Printf("Raw JSON response: %s", jsonStr)

	// Sanitize the JSON string
	sanitizedJSON := sanitizeJSONString(jsonStr)
	// Ensure that opening and closing brackets match
	if strings.Count(sanitizedJSON, "[") != strings.Count(sanitizedJSON, "]") {
		return nil, fmt.Errorf("json format error: unmatched brackets")
	}
	// Unmarshal into a slice of ResponseItem
	var responseItems []ResponseItem
	err := json.Unmarshal([]byte(sanitizedJSON), &responseItems)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return responseItems, nil
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

// Appends the output to the chat history
func (kai *Kai) appendToChatHistory(role string, content string) {
	kai.Chat.History = append(kai.Chat.History, &genai.Content{
		Parts: []genai.Part{
			genai.Text(content),
		},
		Role: role,
	})
}