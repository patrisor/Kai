package ui

import (
	"log"
	"image/color"
	// Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// Local imports
	"kai/source/core"
)

// Method displays the loading screen.
func ShowLoadingScreen(window fyne.Window, state *core.AppState) {
	// Change background color
	backgroundColor := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	background := canvas.NewRectangle(backgroundColor)
	// Create a loading text label
	loadingText := canvas.NewText("Please wait while I scan your system...", color.Black)
	loadingText.Alignment = fyne.TextAlignCenter
	loadingText.TextStyle = fyne.TextStyle{Bold: true}
	loadingText.TextSize = 16
	// Set the content of the window
	window.SetContent(
		container.NewStack(
			background,
			container.NewCenter(loadingText),
		),
	)
	// Run the system scan in a separate goroutine
	go func() {
		// Prime the AI with the default primer
		defaultPrimer, exists := state.Prompts.Primers["Default"]
		if !exists {
			log.Fatalf("Default primer not found")
		}
		state.Kai.PrimeAI(defaultPrimer, state.HistoryFile)
		// Create a primer message that instructs Kai to scan for available commands
		systemScanPrimer, exists := state.Prompts.Primers["SystemScan"]
		if !exists {
			log.Fatalf("SystemScan primer not found")
		}
		// Send the message to Kai
		responseJSON, err := state.Kai.Reason(systemScanPrimer)
		if err != nil {
			log.Fatalf("Failed to send command scan message: %v", err)
		}
		// Process the JSON response
		state.Kai.Respond(responseJSON)
		// After the scan, transition to the Home screen
		ShowHomeScreen(window, state)
	}()
}
