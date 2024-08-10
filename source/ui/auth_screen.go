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
	// Initialize widgets
	apiKeyEntry, apiKeyEntryContainer := createAPIKeyEntry()
	submitButtonContainer := createSubmitButton(apiKeyEntry, window, state)
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

// Method creates the API key input field and its container.
func createAPIKeyEntry() (*widget.Entry, *fyne.Container) {
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
		if enteredAPIKey == "" {
			fmt.Println("API Key cannot be empty")
		} else {
			// Attempt to initialize Kai
			kai, err := core.InitializeKai(enteredAPIKey, state.HistoryFile)
			if err != nil {
				fmt.Println("Invalid API Key")
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
	})
	submitButtonContainer := container.NewHBox(
		layout.NewSpacer(),
		submitButton,
		layout.NewSpacer(),
	)
	submitButtonContainer.Resize(fyne.NewSize(400, 40))
	return submitButtonContainer
}