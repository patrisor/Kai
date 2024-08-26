package core

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// Sends a message to the chat and processes the response.
func (kai *Kai) Reason(userInput string) (string, error) {
	if kai.Chat == nil {
		return "", fmt.Errorf("Kai's Chat is not initialized")
	}
	// Send a message to the chat
	resp, err := kai.Chat.SendMessage(kai.Context, genai.Text(userInput))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}
	// Parse the response
	candidates := resp.Candidates
	if len(candidates) <= 0 || len(candidates[0].Content.Parts) <= 0 {
		return "", fmt.Errorf("no content generated")
	}
	responseJSON := partToString(candidates[0].Content.Parts[0])
	return responseJSON, nil
}

// Helper function to convert genai.Part to a string
func partToString(part genai.Part) string {
	if v, ok := part.(genai.Text); ok {
		return string(v)
	}
	return ""
}
