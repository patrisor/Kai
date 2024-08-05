package ui

import (
	"fmt"
	"log"
	"time"
	"image/color"
	// Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	// Local imports
	"kai/source/ui/gui_elements"
	"kai/source/core"
)

// Method displays the home screen.
func ShowHomeScreen(window fyne.Window, state *core.AppState) {
	// Prime the AI with the default primer
	defaultPrimer, exists := state.Prompts.Primers["Default"]
	if !exists {
		log.Fatalf("Default primer not found")
	}
	state.Kai.PrimeAI(defaultPrimer, state.HistoryFile)
	// Set background color
	backgroundColor := color.NRGBA{R: 253, G: 252, B: 251, A: 255}
	background := canvas.NewRectangle(backgroundColor)
	// Create components
	instructionText := createInstructionText()
	transcriptionText := createTranscriptionText()
	button := createRecordButton(state, transcriptionText)
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

// Method creates the instruction text label.
func createInstructionText() *canvas.Text {
	instructionText := canvas.NewText(
		"How can I assist you today? Hold the button and speak.",
		color.Black,
	)
	instructionText.Alignment = fyne.TextAlignCenter
	instructionText.TextStyle = fyne.TextStyle{Bold: true}
	instructionText.TextSize = 16
	return instructionText
}

// Method creates the transcription text label.
func createTranscriptionText() *canvas.Text {
	transcriptionText := canvas.NewText("", color.Black)
	transcriptionText.Alignment = fyne.TextAlignCenter
	transcriptionText.TextSize = 12
	return transcriptionText
}

// Method creates the record button with mouse event handling.
func createRecordButton(
	state *core.AppState, 
	transcriptionText *canvas.Text,
) *gui_elements.HoldableImageButton {
	imageLocation := "resources/assets/logo.png"
	imageResource, err := fyne.LoadResourceFromPath(imageLocation)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}
	var pressStartTime time.Time
	var stopChan chan struct{}
	var audioData []byte
	button := gui_elements.NewHoldableImageButton(
		imageResource, fyne.NewSize(80, 80),
		func() { // Button press event

			// TODO: Testing
			fmt.Println("Button pressed, starting recording...")

			pressStartTime = time.Now()
			// Create a new stop channel for each recording
			stopChan = make(chan struct{})
			// Start recording audio when the button is pressed
			go func() {
				var err error
				audioData, err = state.Kai.Listen(stopChan)
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
				transcript, err := state.Kai.Recognize(audioData)
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
				_, err = state.Kai.Reason(transcript)
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
	return button
}