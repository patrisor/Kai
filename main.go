package main

import (
    "os"
	"fmt"
	"log"
    "path/filepath"
	// GUI
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
    // Local utilities
    "kai/source/core"
    "kai/source/ui"
)

// Method sets the app's environment variables programmatically
func init() {
    configDir := filepath.Join("config", "service-account-file.json")
    if _, err := os.Stat(configDir); err == nil {
        os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", configDir)
    } else {
        log.Fatalf("Service account file not found: %v", err)
    }
}

// Method initializes an application window
func initializeWindow(
    app fyne.App, 
    title string, 
    width, height float32,
) fyne.Window {
	window := app.NewWindow(title)
	window.Resize(fyne.NewSize(width, height))
	window.CenterOnScreen()
	window.SetMaster()
	return window
}

// Method initializes the application state
func InitializeAppState() (*core.AppState, error) {
    configFile := "config/config.json"
    historyFile := "data/history.json"
    promptsFile := "resources/prompts.json"
    // Load the app configuration
    config, err := core.LoadConfig(configFile)
    if err != nil {
        return nil, err
    }
    // Load the prompts
    prompts, err := core.LoadPrompts(promptsFile)
    if err != nil {
        return nil, err
    }
    // Initialize AppState
    state := &core.AppState{
        Config:      config,
        ConfigFile:  configFile,
        HistoryFile: historyFile,
        Prompts:     prompts,
    }
    return state, nil
}

func main() {
    // Initialize AppState
    state, err := InitializeAppState()
    if err != nil {
        log.Fatalf("Failed to initialize AppState: %v", err)
    }
    // Initialize Fyne app
    fyneApp := app.New()
    // Initialize window
    window := initializeWindow(fyneApp, "Kai", 640, 480)
    // Attempt to initialize Kai with the API key if it exists
    if state.Config.APIKey != "" {
        state.Kai, err = core.InitializeKai(
            state.Config.APIKey,
            state.HistoryFile,
        )
        if err == nil {
            defer state.Kai.Client.Close()
            ui.ShowHomeScreen(window, state)
        } else {
            ui.ShowAuthScreen(window, state)
        }
    } else {
        ui.ShowAuthScreen(window, state)
    }
    // Show the window and run the application
    window.ShowAndRun()
    // Application exitting
    fmt.Println("Exiting...")
}