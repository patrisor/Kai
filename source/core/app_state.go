package core

import (
	"context"
)

// AppState holds the application state.
type AppState struct {
	Kai					*Kai
	Prompts     		*Prompts
	Config				*Config
	ConfigFile			string
	HistoryFile			string

	// TODO: Delete
	ProcessContext   	context.Context
	CancelProcessFunc 	context.CancelFunc 
	
}