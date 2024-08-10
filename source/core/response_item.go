package core

// ResponseItem represents a single response from the AI.
type ResponseItem struct {
    Type string `json:"type"`
    Data string `json:"data"`
}
