package core

import (
	"fmt"
	"context"
	// Google Cloud
	"google.golang.org/api/option"
	// Gemini API
	"github.com/google/generative-ai-go/genai"
)

type Kai struct {
	ApiKey      string
	HistoryFile string
	Client      *genai.Client
	Model       *genai.GenerativeModel
	Chat        *genai.ChatSession
	Context     context.Context
	SampleRate  int
}

// Method initializes and validates a new Kai instance with the 
// provided API key.
func InitializeKai(apiKey, historyFile string) (*Kai, error) {
	// Initialize the Gemini API client with the API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}
	// Define and configure the model
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}
	// Create the Kai instance
	kai := &Kai{
		ApiKey:      apiKey,
		HistoryFile: historyFile,
		Client:      client,
		Model:       model,
		Chat:        model.StartChat(),
		Context:     ctx,
		SampleRate:  44100, // CD quality
	}
	// Validate the API key by making a lightweight request
	iter := kai.Client.ListModels(kai.Context)
	if iter == nil {
		return nil, fmt.Errorf("invalid API key")
	}
	// Check if we can fetch at least one model to confirm the API key is valid
	_, err = iter.Next()
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}
	return kai, nil
}