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
func (kai *Kai) Respond(jsonStr string, branchCounts ...int) {
	// Set branchCount to 1 by default
	branchCount := 1
	if len(branchCounts) > 0 {
		branchCount = branchCounts[0]
	}

	// TODO: Testing
	// Print the current branch count for testing
	// fmt.Printf("Processing branch: %d\n", branchCount)

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
			processScript(kai, item.Data, branchCount)
		case "command":
			processCommand(kai, item.Data, branchCount)
		default:
			fmt.Println("Unknown type:", item.Type)
		}
	}
}

// Process a script response item
func processScript(kai *Kai, data json.RawMessage, branchCount int) {
	// Unmarshal the script data
	var scriptData struct {
		Message string `json:"message"`
		Role    string `json:"role"`
	}
	err := json.Unmarshal(data, &scriptData)
	if err != nil {
		log.Printf("Failed to parse script data: %v", err)
		return
	}
	// Skip unwanted conclusion scripts in higher branches
	if scriptData.Role == "conclusion" && branchCount > 1 {
		return
	}

	// TODO: Print the role for context and debugging
	// fmt.Printf("Speaking (%s): %s\n", scriptData.Role, scriptData.Message)

	// Have Kai speak the script
	if err := kai.Speak(scriptData.Message); err != nil {
		log.Printf("Failed to speak: %v", err)
	}
}

// Process a command response item
func processCommand(kai *Kai, data json.RawMessage, branchCount int) {
	// Unmarshal the command data
	var commandData struct {
		Command string `json:"command"`
	}
	err := json.Unmarshal(data, &commandData)
	if err != nil {
		log.Printf("Failed to parse command data: %v", err)
		return
	}

	// TODO: Testing
	// fmt.Println("Executing command:", commandData.Command)

	// Execute command
	sanitizedCommand := strings.ReplaceAll(commandData.Command, ",", "")
	output, err := kai.executeCommand(sanitizedCommand)
	if err != nil { // Error occurred
		errorMessage := fmt.Sprintf(
			"Command failed: %v. " + 
			"Please analyze the error and generate a new solution.", 
			err,
		)
		// Append command output if it exists
		if output != "" {
			errorMessage += fmt.Sprintf(" Command output: %s.", output)
		}
		// Handle the AI response for errors
		kai.handleAIResponse(errorMessage, branchCount)
	} else { // No error occurred
		if output != "" { // Handle the AI response for success with output
			successMessage := fmt.Sprintf(
				"Command executed successfully: %s. " +
				"Please analyze the output and provide a suitable response.",
				// + branchPrompt,
				output,
			)
			kai.handleAIResponse(successMessage, branchCount)
		}
		// Handle the AI response for success without output
	}
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

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
	// log.Printf("Raw JSON response: %s", jsonStr)

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

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method to handle AI response generation and chat history update
func (kai *Kai) handleAIResponse(message string, branchCount int) {
	// Feed the message back into the system to generate a new response
	if newResponse, err := kai.Reason(message); err != nil {
		log.Fatalf("Failed to process new response: %v", err)
	} else {
		kai.Respond(newResponse, branchCount + 1)
	}
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