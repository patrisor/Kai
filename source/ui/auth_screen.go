package ui

import (
	"fmt"
	"log"
	"image/color"
	// Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	// Local imports
	"kai/source/core"
)

// Method displays the authentication screen.
func ShowAuthScreen(window fyne.Window, state *core.AppState) {
	// Change background color
	backgroundColor := color.NRGBA{R: 253, G: 252, B: 251, A: 255}
	background := canvas.NewRectangle(backgroundColor)
	// Initialize error label component
	errorLabel := createErrorLabel()
	// Initialize logo component
	logo := createLogo("resources/assets/logo.png")
	// Initialize key entry component
	apiKeyEntry, apiKeyEntryContainer := createAPIKeyEntry(window, state)
	apiKeyEntry.OnSubmitted = func(input string) {
		initializeKai(input, window, state, errorLabel)
	}
	// Initialize submit button component
	submitButtonContainer := createSubmitButton(apiKeyEntry, window, state)
	// Set the content of the window
	window.SetContent(
		container.NewStack(
			background,
			container.NewVBox(
				layout.NewSpacer(),
				container.NewCenter(
					container.NewVBox(
						logo,
						container.NewVBox(
							layout.NewSpacer(),
							container.NewGridWrap(fyne.NewSize(0, 50)),
						),
						apiKeyEntryContainer,
						submitButtonContainer,
					),
				),
				container.NewCenter(errorLabel),
				layout.NewSpacer(),
			),
		),
	)	
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method to create the logo widget using a resource.
func createLogo(imageLocation string) *canvas.Image {
	// Load the image resource
	imageResource := loadImageResource(imageLocation)
	// Create an image from the resource
	logo := canvas.NewImageFromResource(imageResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(100, 100))
	return logo
}

// Method creates the API key input field and its container.
func createAPIKeyEntry(
	window fyne.Window, 
	state *core.AppState,
) (*widget.Entry, *fyne.Container) {
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enter API Key")
	apiKeyEntryContainer := container.NewGridWrap(
		fyne.NewSize(400, 40), apiKeyEntry,
	)
	return apiKeyEntry, apiKeyEntryContainer
}

// Method to creates the submit button and its container.
func createSubmitButton(
	apiKeyEntry *widget.Entry, 
	window fyne.Window, 
	state *core.AppState,
) *fyne.Container {
	submitButton := widget.NewButton("Submit", func() {
		enteredAPIKey := apiKeyEntry.Text
		initializeKai(enteredAPIKey, window, state, nil)
	})
	submitButtonContainer := container.NewHBox(
		layout.NewSpacer(),
		submitButton,
		layout.NewSpacer(),
	)
	submitButtonContainer.Resize(fyne.NewSize(400, 40))
	return submitButtonContainer
}

// Method to create a reusable error label.
func createErrorLabel() *canvas.Text {
	errorLabel := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	errorLabel.Alignment = fyne.TextAlignCenter
	errorLabel.Hide()
	return errorLabel
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method to display an error message using the error label.
func displayError(message string, errorLabel *canvas.Text) {
	if errorLabel != nil {
		errorLabel.Text = message
		errorLabel.Show()
		errorLabel.Refresh()
	}
	fmt.Println(message)
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method initializes Kai.
func initializeKai(
	enteredAPIKey string,
	window fyne.Window,
	state *core.AppState,
	errorLabel *canvas.Text,
) {
	// Clear previous error messages
	if errorLabel != nil {
		errorLabel.Hide()
		errorLabel.Text = ""
		errorLabel.Refresh()
	}
	// Validate API Key
	if enteredAPIKey == "" {
		displayError("API Key cannot be empty", errorLabel)
		return
	}
	// Attempt to initialize Kai
	kai, err := core.InitializeKai(enteredAPIKey, state.HistoryFile)
	if err != nil {
		displayError("Invalid API Key", errorLabel)
		return
	}
	// Successfully initialized Kai, update app state
	state.Kai = kai
	state.Config.APIKey = enteredAPIKey
	// Save configuration
	err = core.SaveConfig(state.ConfigFile, state.Config)
	if err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
	// Show Loading screen
	ShowLoadingScreen(window, state)
}