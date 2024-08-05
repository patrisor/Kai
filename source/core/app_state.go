package core

// AppState holds the application state.
type AppState struct {
	Kai			*Kai
	Prompts     *Prompts
	Config		*Config
	ConfigFile	string
	HistoryFile	string
}