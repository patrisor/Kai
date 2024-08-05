package core

import (
	"log"
	"os"
	"encoding/json"
)

// Saves the chat history to a file.
func (kai *Kai) SaveHistory() {
    if kai.HistoryFile == "" {
        log.Println("No history file specified.")
        return
    }
    // Check if there's anything in the history to save
    if len(kai.Chat.History) == 0 {
        log.Println("No chat history to save.")
        return
    }
    // Marshal the chat history into a formatted JSON byte slice
    data, err := json.MarshalIndent(kai.Chat.History, "", "  ")
    if err != nil {
        log.Println("Failed to marshal chat history:", err)
        return
    }
    // Write the JSON byte slice to the specified history file
    if err := os.WriteFile(kai.HistoryFile, data, 0o644); err != nil {
        log.Println("Failed to save chat history:", err)
        return
    }
}