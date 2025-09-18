package signaling

import "encoding/json"

const (
	MessageTypeClientRequestPing uint32 = iota
	MessageTypeClientSendSignal
)

// Message ..
type Message struct {
	Type uint32      `json:"Type"`
	From string      `json:"From,omitempty"`
	To   json.Number `json:"To,omitempty"`
	Data string      `json:"Message,omitempty"`
}
