package main

import (
	"os"
    "fmt"
    "log"
	"time"
	"bufio"
	"strings"
	"os/exec"
	"runtime"
    "context"
	"encoding/json"
	// 
	"google.golang.org/api/option"
	// Gemini
    "github.com/google/generative-ai-go/genai"

	// 
	"github.com/gordonklaus/portaudio"

	// Text-to-Speech
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"

)

type ResponseItem struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

type Kai struct {
    apiKey		string
	historyFile	string
    client		*genai.Client
    model		*genai.GenerativeModel
    chat		*genai.ChatSession
    ctx			context.Context
	SampleRate	int
}

// Creates a new instance of the Kai struct.
func NewKai(apiKey, primer, historyFile string) (*Kai, error) {
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
	// Initialize the chat session
	chat := model.StartChat()
	primeAI(chat, primer, historyFile)
	// Initialize the AI
    return &Kai{
        apiKey: apiKey,
		historyFile: historyFile,
        client: client,
        model:  model,
        chat:   chat,
        ctx:    ctx,
		SampleRate: 44100, // CD quality
    }, nil
}

// Primes the AI.
func primeAI(chat *genai.ChatSession, primer, historyFile string) {
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
									content.Parts, 
									genai.Text(textPart),
								)
                            }
                        }
                        history = append(history, content)
                    }
                    chat.History = history
                }
            }
        } else { // History file does not exist, using primer
            if primer != "" {
				// Append system information to the primer
				systemInfo := "\n\nSystem Information:\n" + getSystemInfo()
				primerWithInfo := primer + systemInfo
				// Prime the AI
                chat.History = []*genai.Content{
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

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Retrieves system information and returns it as a string.
func getSystemInfo() string {
    systemInfo := map[string]string{
        "OS":               runtime.GOOS,
        "Architecture":     runtime.GOARCH,
        "CPU Count":        fmt.Sprintf("%d", runtime.NumCPU()),
        "Go Version":       runtime.Version(),
        "Hostname":         getHostname(),
        "Current User":     getCurrentUser(),
        "Home Directory":   getHomeDirectory(),
        "Environment Vars": getEnvironmentVariables(),
    }
    info, _ := json.MarshalIndent(systemInfo, "", "  ")
    return string(info)
}

// Helper function to get the hostname.
func getHostname() string {
    hostname, err := os.Hostname()
    if err != nil {
        return "unknown"
    }
    return hostname
}

// Helper function to get the current user.
func getCurrentUser() string {
    user := os.Getenv("USER")
    if user == "" {
        return "unknown"
    }
    return user
}

// Helper function to get the home directory.
func getHomeDirectory() string {
    home := os.Getenv("HOME")
    if home == "" {
        return "unknown"
    }
    return home
}

// Helper function to get the environment variables.
func getEnvironmentVariables() string {
    envVars := os.Environ()
    return strings.Join(envVars, "\n")
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Saves the chat history to a file.
func (kai *Kai) SaveHistory() {
    if kai.historyFile != "" {
		// Marshal the chat history into a formatted JSON byte slice
        data, err := json.MarshalIndent(kai.chat.History, "", "  ")
        if err != nil {
            log.Println("Failed to marshal chat history:", err)
            return
        }
		// Write the JSON byte slice to the specified history file
        if err := os.WriteFile(kai.historyFile, data, 0644); err != nil {
            log.Println("Failed to save chat history:", err)
			return
        } 
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

			// TODO: Add code to handle speaking the script
            if err := kai.speak(item.Data); err != nil {
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

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Synthesizes speech from the input text and plays it.
func (kai *Kai) speak(text string) error {
    ctx := context.Background()
    client, err := texttospeech.NewClient(ctx)
    if err != nil {
        return fmt.Errorf("failed to create text-to-speech client: %v", err)
    }
    defer client.Close()
    // Perform the text-to-speech request
    req := &texttospeechpb.SynthesizeSpeechRequest{
        Input: &texttospeechpb.SynthesisInput{
            InputSource: &texttospeechpb.SynthesisInput_Text{
                Text: text,
            },
        },
        Voice: &texttospeechpb.VoiceSelectionParams{
            LanguageCode: "en-US",
            Name:         "en-US-Polyglot-1",
            SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
        },
        AudioConfig: &texttospeechpb.AudioConfig{
            AudioEncoding:    texttospeechpb.AudioEncoding_LINEAR16,
            SpeakingRate:     1.0,
            Pitch:            0.0,
            SampleRateHertz:  int32(kai.SampleRate),
        },
    }
    resp, err := client.SynthesizeSpeech(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to synthesize speech: %v", err)
    }
    // Play the audio (you can use any audio playing library)
    if err := playAudio(resp.AudioContent, kai.SampleRate); err != nil {
        return fmt.Errorf("failed to play audio: %v", err)
    }
    return nil
}

// Helper function to play audio using PortAudio
func playAudio(audioData []byte, sampleRate int) error {
    // Initialize PortAudio
    err := portaudio.Initialize()
    if err != nil {
        return fmt.Errorf("failed to initialize PortAudio: %v", err)
    }
    defer portaudio.Terminate()
    // Define a buffer size and create a buffer
    bufferSize := 1024
    dataSize := len(audioData) / 2
    currentIndex := 0
    // Open a stream for audio playback
    stream, err := portaudio.OpenDefaultStream(0, 1, float64(sampleRate), bufferSize, func(out []int16) {
        for i := range out {
            if currentIndex < dataSize {
                out[i] = int16(audioData[2*currentIndex]) | int16(audioData[2*currentIndex+1])<<8
                currentIndex++
            } else {
                out[i] = 0 // Fill with silence if the data is finished
            }
        }
    })
    if err != nil {
        return fmt.Errorf("failed to open stream: %v", err)
    }
    defer stream.Close()
    // Start the audio stream
    if err := stream.Start(); err != nil {
        return fmt.Errorf("failed to start stream: %v", err)
    }
    // Wait for the audio to finish playing
    for currentIndex < dataSize {
        time.Sleep(10 * time.Millisecond)
    }
    // Stop the audio stream
    if err := stream.Stop(); err != nil {
        return fmt.Errorf("failed to stop stream: %v", err)
    }
    return nil
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Sends a message to the chat and processes the response.
func (kai *Kai) SendMessage(userInput string) (string, error) {
	// Send a message to the chat
    resp, err := kai.chat.SendMessage(kai.ctx, genai.Text(userInput))
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

// Starts the event loop for shell-based interactions.
func (kai *Kai) RunShell() {
	// Ensure history is saved when the function exits
	defer kai.SaveHistory()
	// Create a new buffered reader for user input
	reader := bufio.NewReader(os.Stdin)
	for {
		// Prompt user for input
		fmt.Print("Kai> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			continue
		}
		// Send the message and process the response
		responseJSON, err := kai.SendMessage(userInput)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		// TODO: Testing response (readable)
		fmt.Print(responseJSON)

	}
}