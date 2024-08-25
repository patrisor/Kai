package core

import "encoding/json"

// ResponseItem represents a single response from the AI.
type ResponseItem struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}