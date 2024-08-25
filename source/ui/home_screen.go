package ui

import (
	"fmt"
	"log"
	"time"
	"image/color"
	// Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	instructionText := createGreetingText()
	textEntryContainer := createTextEntryContainer(state)
	// Set the content of the window
	window.SetContent(
		container.NewStack(
			background,
			container.NewVBox(
				layout.NewSpacer(),
				container.NewCenter(instructionText),
				layout.NewSpacer(),
				textEntryContainer,
			),
		),
	)
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method creates the greeting text label.
func createGreetingText() *canvas.Text {
	greetingText := canvas.NewText(
		"How can I assist you today? Hold the button and speak, " + 
		"or type a message and press Enter.",
		color.Black,
	)
	greetingText.Alignment = fyne.TextAlignCenter
	greetingText.TextStyle = fyne.TextStyle{Bold: true}
	greetingText.TextSize = 16
	return greetingText
}

// Method creates the text entry field and its container.
func createTextEntryContainer(state *core.AppState) *fyne.Container {
	// Create the text entry
	textEntry := widget.NewEntry()
	textEntry.SetPlaceHolder("Type your message here...")
	textEntry.OnSubmitted = func(input string) {
		// Clear the text field after processing
		updateTextEntry(textEntry, "")
		// Process input
		processUserInput(state, input, textEntry)
	}
	// Set the size of the text entry to be shorter
	textEntryContainer := container.NewVBox(
		container.NewPadded(container.NewStack(textEntry)),
	)
	// Create the Listen button
	button := createListenButton(state, textEntry)
	// Combine the text entry and button in an HBox layout with padding
	content := container.NewBorder(
		nil, nil, nil, button,
		textEntryContainer,
	)
	// Align the container to the bottom with padding
	finalContainer := container.NewVBox(
		layout.NewSpacer(),
		container.NewPadded(content),
	)
	return finalContainer
}

// Method creates the transcription text label.
func createTranscriptionText() *canvas.Text {
	transcriptionText := canvas.NewText("", color.Black)
	transcriptionText.Alignment = fyne.TextAlignCenter
	transcriptionText.TextSize = 12
	return transcriptionText
}

// Method creates the Listen button with mouse event handling.
func createListenButton(
	state *core.AppState, 
	textEntry *widget.Entry,
) *gui_elements.HoldableImageButton {
	imageResource := loadImageResource("resources/assets/microphone.svg")
	var pressStartTime time.Time
	var stopChan chan struct{}
	var audioData []byte
	button := gui_elements.NewHoldableImageButton(
		imageResource, fyne.NewSize(40, 40),
		func() { // Button press event
			fmt.Println("Button pressed, starting recording...")
			handleListenButtonPress(
				state, &pressStartTime, 
				&stopChan, &audioData, 
				textEntry,
			)
		},
		func() { // Button release event
			fmt.Println("Button released, stopping recording...")
			handleListenButtonRelease(pressStartTime, stopChan)
		},
	)
	return button
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Loads the image resource from the given path.
func loadImageResource(imageLocation string) fyne.Resource {
	imageResource, err := fyne.LoadResourceFromPath(imageLocation)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}
	return imageResource
}

// Method processes the user input text.
func processUserInput(
	state *core.AppState, 
	input string, 
	textEntry *widget.Entry,
) {
	// Only process if the input is non-empty
	if input != "" {
		// Send transcription to Gemini API
		responseJSON, err := state.Kai.Reason(input)
		if err != nil {
			log.Fatalf("Failed to send command scan message: %v", err)
		}
		// Process the JSON response
		state.Kai.Respond(responseJSON)
	}
}

// Method updates the text entry with the given transcript.
func updateTextEntry(textEntry *widget.Entry, transcript string) {
	textEntry.SetText(transcript)
	textEntry.Refresh()
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method handles the button press event.
func handleListenButtonPress(
	state *core.AppState, 
	pressStartTime *time.Time, 
	stopChan *chan struct{},
	audioData *[]byte,
	textEntry *widget.Entry,
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

		// Clear the text field after processing
		updateTextEntry(textEntry, "")
		// Transcribe the recorded audio
		transcript, err := state.Kai.Recognize(*audioData)
		if err != nil {
			log.Println("Failed to transcribe audio:", err)
			// TODO: You might want to add some fallback behavior here, 
			// such as retrying or notifying the user
			return
		}
		// Update the transcription label
		updateTextEntry(textEntry, transcript)
		// Process the transcription
		processUserInput(state, transcript, textEntry)
		// Clear the text field after processing
		updateTextEntry(textEntry, "")
	}()
}

// Method handles the button release event.
func handleListenButtonRelease(
	pressStartTime time.Time, 
	stopChan chan struct{},
) {

	// TODO: Testing
	// pressDuration := time.Since(pressStartTime)
	// fmt.Printf("Button held for %v\n", pressDuration)

	// Signal to stop recording when the button is released
	if stopChan != nil {
		close(stopChan)
	}
}