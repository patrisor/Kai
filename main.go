package main

import (
    "os"
	"fmt"
	"log"
    "time"
    "image/color"
    "path/filepath"
    "encoding/json"
    "encoding/binary"
	// GUI
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/driver/desktop"
    "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

    // TODO: Transfer code to Kai
    // Audio Recording
    "github.com/gordonklaus/portaudio"
    // Speech-to-Text
    "context"
    speech "cloud.google.com/go/speech/apiv1"
    "cloud.google.com/go/speech/apiv1/speechpb"

)

// AppState holds the application state.
type AppState struct {
	kai     *Kai
	config  *Config
	configFile string
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Config structure to hold API key.
type Config struct {
    APIKey string `json:"api_key"`
    Primer string `json:"primer"`
    HistoryFile string `json:"history_file"`
}

// Method writes the Config struct to the configuration file.
func saveConfig(file string, config *Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, configData, 0644)
}

// Method reads a configuration file and returns a Config struct.
func loadConfig(file string) (*Config, error) {
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

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Custom holdable image button to handle mouse events
type holdableImageButton struct {
    widget.BaseWidget
    image     *canvas.Image
    onPress   func()
    onRelease func()
    size      fyne.Size
}

func newHoldableImageButton(
    image fyne.Resource, 
    size fyne.Size, 
    onPress, 
    onRelease func(),
) *holdableImageButton {
    button := &holdableImageButton{
        image:     canvas.NewImageFromResource(image),
        onPress:   onPress,
        onRelease: onRelease,
        size:      size,
    }
    button.image.SetMinSize(size)
    button.image.FillMode = canvas.ImageFillContain
    button.ExtendBaseWidget(button)
    return button
}

func (h *holdableImageButton) MouseDown(*desktop.MouseEvent) {
    if h.onPress != nil {
        h.onPress()
    }
}

func (h *holdableImageButton) MouseUp(*desktop.MouseEvent) {
    if h.onRelease != nil {
        h.onRelease()
    }
}

func (h *holdableImageButton) CreateRenderer() fyne.WidgetRenderer {
    return widget.NewSimpleRenderer(h.image)
}

func (h *holdableImageButton) MinSize() fyne.Size {
    return h.size
}

func (h *holdableImageButton) Resize(size fyne.Size) {
    h.size = size
    h.image.Resize(size)
    h.BaseWidget.Resize(size)
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method displays the authentication screen.
func showAuthScreen(window fyne.Window, state *AppState) {
	// Change background color
	backgroundColor := color.NRGBA{R: 253, G: 252, B: 251, A: 255}
	background := canvas.NewRectangle(backgroundColor)

	// Initialize widgets
	// API key input field
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enter API Key")
	apiKeyEntryContainer := container.NewGridWrap(
		fyne.NewSize(400, 40), apiKeyEntry,
	)
	// Submit button
	submitButton := widget.NewButton("Submit", func() {
		enteredAPIKey := apiKeyEntry.Text
		if enteredAPIKey == "" {
			fmt.Println("API Key cannot be empty")
		} else {
			kai, err := initializeKai(
                enteredAPIKey, 
                state.config.Primer, 
                state.config.HistoryFile,
            )
			if err != nil {
				fmt.Println("Invalid API Key")
				return
			}
			state.kai = kai
			state.config.APIKey = enteredAPIKey
			err = saveConfig(state.configFile, state.config)
			if err != nil {
				log.Fatalf("Failed to save config: %v", err)
			}
			showHomeScreen(window, state)
		}
	})
	submitButtonContainer := container.NewHBox(
		layout.NewSpacer(),
		submitButton,
		layout.NewSpacer(),
	)
	submitButtonContainer.Resize(fyne.NewSize(400, 40))

	// Set the content of the window
	window.SetContent(
		container.NewStack(
			background,
			container.NewCenter(
				container.NewVBox(
					apiKeyEntryContainer,
					submitButtonContainer,
				),
			),
		),
	)
}

// Method displays the home screen.
func showHomeScreen(window fyne.Window, state *AppState) {
    // Set background color
    backgroundColor := color.NRGBA{R: 253, G: 252, B: 251, A: 255}
    background := canvas.NewRectangle(backgroundColor)
    // Create label to instruct user
    instructionText := canvas.NewText(
        "How can I assist you today? Hold the button and speak.",
        color.Black,
    )
    instructionText.Alignment = fyne.TextAlignCenter
    instructionText.TextStyle = fyne.TextStyle{Bold: true}
    instructionText.TextSize = 16
    // Create label to display transcription result
    transcriptionText := canvas.NewText("", color.Black)
    transcriptionText.Alignment = fyne.TextAlignCenter
    transcriptionText.TextSize = 12
    // Load the image for the button
    imageResource, err := fyne.LoadResourceFromPath("assets/logo.png")
    if err != nil {
        log.Fatalf("Failed to load image: %v", err)
    }
    // Custom button with mouse event handling
    var pressStartTime time.Time
    var stopChan chan struct{}
    var audioData []byte
    button := newHoldableImageButton(imageResource, fyne.NewSize(80, 80),
        func() { // Button press event

            // TODO: Testing
            fmt.Println("Button pressed, starting recording...")

            pressStartTime = time.Now()
            // Create a new stop channel for each recording
            stopChan = make(chan struct{})
            // Start recording audio when the button is pressed
            go func() {
                var err error
                audioData, err = recordAudio(stopChan)
                if err != nil {
                    log.Fatalf("Failed to record audio: %v", err)
                }

                // TODO: Testing
                fmt.Println("Recording complete")
                fmt.Printf("Recorded %d bytes of audio data\n", len(audioData))
                // Save the recorded audio to a WAV file
                // err = saveToWavFile("recording.wav", audioData)
                // if err != nil {
                //     log.Fatalf("Failed to save WAV file: %v", err)
                // }
                // fmt.Println("Saved recording to recording.wav")

                // Transcribe the recorded audio
                transcript, err := transcribeAudio(audioData)
                if err != nil {
                    log.Println("Failed to transcribe audio:", err)
                    // You might want to add some fallback behavior here, 
                    // such as retrying or notifying the user
                    return
                }
                // Update the transcription label
                transcriptionText.Text = transcript
                transcriptionText.Refresh()
                // Send transcription to Gemini API
                _, err = state.kai.SendMessage(transcript)
                if err != nil {
                    log.Fatalf("Failed to send message: %v", err)
                }
            }()
        },
        func() { // Button release event
            pressDuration := time.Since(pressStartTime)

            // TODO: Testing
            fmt.Println("Button released, stopping recording...")
            fmt.Printf("Button held for %v\n", pressDuration)

            // Signal to stop recording when the button is released
            if stopChan != nil {
                close(stopChan)
            }
        },
    )
    // Set the content of the window
    window.SetContent(
        container.NewStack(
            background,
            container.NewVBox(
                layout.NewSpacer(),
                container.NewCenter(instructionText),
                layout.NewSpacer(),
                container.NewCenter(button),
                layout.NewSpacer(),
                container.NewCenter(transcriptionText),
                layout.NewSpacer(),
            ),
        ),
    )
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method captures audio input from the microphone using the portaudio library.
func recordAudio(stop <-chan struct{}) ([]byte, error) {
    // Initialize the PortAudio library
    err := portaudio.Initialize()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize PortAudio: %v", err)
    }
    defer portaudio.Terminate()
    // Create an input buffer to store audio samples
    in := make([]int16, 64)
    // Open a default stream for audio input
    stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
    if err != nil {
        return nil, fmt.Errorf("failed to open default stream: %v", err)
    }
    defer stream.Close()
    // Start the audio stream
    err = stream.Start()
    if err != nil {
        return nil, fmt.Errorf("failed to start the stream: %v", err)
    }
    defer stream.Stop()
    // Record audio data
    var audioData []int16
    for {
        select {
        case <-stop:
            // Stop recording when a signal is received on the stop channel
            return convertToBytes(audioData), nil
        default:
            // Read audio samples into the input buffer
            err := stream.Read()
            if err != nil {
                return nil, fmt.Errorf("failed to read from stream: %v", err)
            }
            // Append the samples to the audio data slice
            audioData = append(audioData, in...)
        }
    }
}

// convertToBytes converts audio data from int16 to a byte slice.
func convertToBytes(audioData []int16) []byte {
    audioBytes := make([]byte, len(audioData)*2)
    for i, sample := range audioData {
        audioBytes[i*2] = byte(sample)
        audioBytes[i*2+1] = byte(sample >> 8)
    }
    return audioBytes
}

// Method saves the recorded audio data to a WAV file.
func saveToWavFile(filename string, audioData []byte) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()
    // WAV file header
    var header = []byte{
        'R', 'I', 'F', 'F',
        0, 0, 0, 0, // ChunkSize (to be filled later)
        'W', 'A', 'V', 'E',
        'f', 'm', 't', ' ',
        16, 0, 0, 0, // Subchunk1Size (16 for PCM)
        1, 0, // AudioFormat (1 for PCM)
        1, 0, // NumChannels (1 for mono)
        0x44, 0xac, 0, 0, // SampleRate (44100 Hz)
        0x88, 0x58, 1, 0, // ByteRate (SampleRate * NumChannels * BitsPerSample/8)
        2, 0, // BlockAlign (NumChannels * BitsPerSample/8)
        16, 0, // BitsPerSample (16 bits)
        'd', 'a', 't', 'a',
        0, 0, 0, 0, // Subchunk2Size (to be filled later)
    }
    // Fill in the ChunkSize and Subchunk2Size
    chunkSize := 36 + len(audioData)
    subchunk2Size := len(audioData)
    binary.LittleEndian.PutUint32(header[4:], uint32(chunkSize))
    binary.LittleEndian.PutUint32(header[40:], uint32(subchunk2Size))
    // Write the header and audio data to the file
    _, err = file.Write(header)
    if err != nil {
        return fmt.Errorf("failed to write header: %v", err)
    }
    _, err = file.Write(audioData)
    if err != nil {
        return fmt.Errorf("failed to write audio data: %v", err)
    }
    return nil
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method sends recorded audio to the Google Cloud Speech-to-Text API for 
// transcription.
func transcribeAudio(audioData []byte) (string, error) {
    ctx := context.Background()
    client, err := speech.NewClient(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to create speech client: %v", err)
    }
    defer client.Close()
    // Configure the request with the correct audio encoding and sample rate
    req := &speechpb.RecognizeRequest{
        Config: &speechpb.RecognitionConfig{
            Encoding:        speechpb.RecognitionConfig_LINEAR16,
            SampleRateHertz: 44100,
            LanguageCode:    "en-US",
        },
        Audio: &speechpb.RecognitionAudio{
            AudioSource: &speechpb.RecognitionAudio_Content{
                Content: audioData,
            },
        },
    }
    // Send the request and get the response
    resp, err := client.Recognize(ctx, req)
    if err != nil {
        return "", fmt.Errorf("failed to recognize speech: %v", err)
    }
    // Process the response and extract the transcribed text
    if len(resp.Results) > 0 && len(resp.Results[0].Alternatives) > 0 {
        return resp.Results[0].Alternatives[0].Transcript, nil
    }
    return "", fmt.Errorf("no transcription results")
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method creates a new Kai instance and validates the API key.
func initializeKai(apiKey, primer, historyFile string) (*Kai, error) {
    // Initialize a new Kai instance
	kai, err := NewKai(apiKey, primer, historyFile)
	if err != nil {
		return nil, err
	}
    // Make a lightweight request to validate the API key
	iter := kai.client.ListModels(kai.ctx)
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

// Set the GOOGLE_APPLICATION_CREDENTIALS environment variable programmatically
func init() {
    configDir := filepath.Join("config", "service-account-file.json")
    if _, err := os.Stat(configDir); err == nil {
        os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", configDir)
    } else {
        log.Fatalf("Service account file not found: %v", err)
    }
}

func main() {
    configFile := "config/config.json"
    // Load the configuration from the config.json file
    config, err := loadConfig(configFile)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

	// Initialize AppState
	state := &AppState{
		config: config,
		configFile: configFile,
	}

    // Initialize Fyne app
    fyneApp := app.New()
    // Initialize window
    windowSize := fyne.NewSize(640, 480)
    window := fyneApp.NewWindow("Kai")
    window.Resize(windowSize)
    window.CenterOnScreen()
    window.SetMaster()
	// Check if API key exists
	if config.APIKey != "" {
		state.kai, err = initializeKai(
            config.APIKey, 
            config.Primer, 
            config.HistoryFile,
        )
		if err == nil {
            defer state.kai.client.Close()
			showHomeScreen(window, state)
		} else {
			showAuthScreen(window, state)
		}
	} else {
		showAuthScreen(window, state)
	}
    // Show the window and run the application
    window.ShowAndRun()
    // Application exitting
    fmt.Println("Exiting...")
}
