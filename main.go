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
    configDir := filepath.Join(".config", "service-account-file.json")
    if _, err := os.Stat(configDir); err == nil {
        os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", configDir)
    } else {
        log.Fatalf("Service account file not found: %v", err)
    }
}

// Method initializes the application state
func InitializeAppState() (*core.AppState, error) {
    configFile := ".config/config.json"
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

// Method initializes Kai with the provided API key and manages the flow of the 
// app.
//
// Parameters:
//  - state: A pointer to the current application state (AppState)
//  - window: A reference to the Fyne window to display the UI
//
// Returns:
//  - An error if the initialization fails, otherwise nil
func initializeKai(state *core.AppState, window fyne.Window) error { 
    if state.Config.APIKey == "" {
        // Show Auth Screen if API Key is not present
        ui.ShowAuthScreen(window, state)
        return nil
    }
    // Attempt to initialize Kai with the API key
    kai, err := core.InitializeKai(state.Config.APIKey, state.HistoryFile)
    if err != nil {
        // If initialization fails, show the Auth Screen
        ui.ShowAuthScreen(window, state)
        return err
    }
    // Assign Kai instance to the application state
    state.Kai = kai
    defer state.Kai.Client.Close()
    // Prime the AI with the default primer
    defaultPrimer, exists := state.Prompts.Primers["Default"]
    if !exists {
        return fmt.Errorf("default primer not found")
    }
    // Prime AI and greet user
    state.Kai.PrimeAI(defaultPrimer, state.HistoryFile)
    go greetUser(state)
    // Show Home Screen after successful initialization
    ui.ShowHomeScreen(window, state)
    return nil
}

/* ************************************************************************* */
/* ************************************************************************* */
/* ************************************************************************* */

// Method requests Kai to send a greeting message to the user.
//
// Parameters:
//  - state: A pointer to the current application state (AppState)
func greetUser(state *core.AppState) {
    greetingMessage := "Greet the user by their name if it is available, " + 
                       "otherwise just greet the user."
    responseJSON, err := state.Kai.Reason(greetingMessage)
    if err != nil {
        log.Fatalf("Failed to send greeting message: %v", err)
    }
    state.Kai.Respond(responseJSON)
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
    window := initializeWindow(fyneApp, "Kai", 1280, 720)

    // TODO: Convert block of code to reusable method, called `initializeKai`
    // if err := initializeKai(state, window); err != nil {
    //     log.Fatalf("Failed to initialize Kai: %v", err)
    // }
    // Attempt to initialize Kai with the API key if it exists
    if state.Config.APIKey != "" {
        state.Kai, err = core.InitializeKai(
            state.Config.APIKey,
            state.HistoryFile,
        )
        if err == nil {
            defer state.Kai.Client.Close()
            // Prime the AI with the default primer
            defaultPrimer, exists := state.Prompts.Primers["Default"]
            if !exists {
                log.Fatalf("Default primer not found")
            }
            state.Kai.PrimeAI(defaultPrimer, state.HistoryFile)
            // Greet the user on a separate goroutine
            go greetUser(state)
            // Show Home Screen
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