package ui

import (
	"fmt"
	"log"
	"time"
	"image/color"
	// Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// Local imports
	"kai/source/ui/gui_elements"
	"kai/source/core"
)

// Method displays the home screen.
func ShowHomeScreen(window fyne.Window, state *core.AppState) {
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
				container.NewCenter(transcriptionText),
				layout.NewSpacer(),
				container.NewCenter(button),
				layout.NewSpacer(),
			),
		),
	)

	// TODO: Make this run only once
	// Instruct Kai to scan available commands
	// primeCommandScan(state)

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
	imageResource := loadImageResource("resources/assets/logo.png")
	var pressStartTime time.Time
	var stopChan chan struct{}
	var audioData []byte
	button := gui_elements.NewHoldableImageButton(
		imageResource, fyne.NewSize(80, 80),
		func() { // Button press event
			fmt.Println("Button pressed, starting recording...")
			handleButtonPress(
				state, &pressStartTime, 
				&stopChan, &audioData, 
				transcriptionText,
			)
		},
		func() { // Button release event
			fmt.Println("Button released, stopping recording...")
			handleButtonRelease(pressStartTime, stopChan)
		},
	)
	return button
}

// Loads the image resource from the given path.
func loadImageResource(imageLocation string) fyne.Resource {
	imageResource, err := fyne.LoadResourceFromPath(imageLocation)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}
	return imageResource
}

// Method handles the button press event.
func handleButtonPress(
	state *core.AppState, 
	pressStartTime *time.Time, 
	stopChan *chan struct{},
	audioData *[]byte,
	transcriptionText *canvas.Text,
) {
	*pressStartTime = time.Now()
	// Create a new stop channel for each recording
	*stopChan = make(chan struct{})
	// Start recording audio when the button is pressed
	go func() {
		var err error
		*audioData, err = state.Kai.Listen(*stopChan)
		if err != nil {
			log.Fatalf("Failed to record audio: %v", err)
		}

		// TODO: Testing
		fmt.Println("Recording complete")
		fmt.Printf("Recorded %d bytes of audio data\n", len(*audioData))
		// Save the recorded audio to a WAV file
		// err = saveToWavFile("recording.wav", audioData)
		// if err != nil {
		//     log.Fatalf("Failed to save WAV file: %v", err)
		// }
		// fmt.Println("Saved recording to recording.wav")

		// Transcribe the recorded audio
		transcript, err := state.Kai.Recognize(*audioData)
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
		responseJSON, err := state.Kai.Reason(transcript)
		if err != nil {
			log.Fatalf("Failed to send command scan message: %v", err)
		}
		// Process the JSON response
		state.Kai.Respond(responseJSON)
	}()
}

// Method handles the button release event.
func handleButtonRelease(
	pressStartTime time.Time, 
	stopChan chan struct{},
) {
	pressDuration := time.Since(pressStartTime)

	// TODO: Testing
	fmt.Printf("Button held for %v\n", pressDuration)

	// Signal to stop recording when the button is released
	if stopChan != nil {
		close(stopChan)
	}
}