package core

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"strings"
)

// Starts the event loop for shell-based interactions.
func (kai *Kai) RunShell() {
	// Ensure history is saved when the function exits
	defer kai.SaveHistory()
	// Create a new buffered reader for user input
	reader := bufio.NewReader(os.Stdin)
	for {
		// Prompt user for input
		fmt.Print("Kai> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			continue
		}
		// Send the message and process the response
		responseJSON, err := kai.Reason(userInput)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		// Process the JSON response
		kai.Respond(responseJSON)

		// Testing response (readable)
		fmt.Print(responseJSON)

	}
}