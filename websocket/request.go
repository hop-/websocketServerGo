package websocket

import (
	"encoding/json"
)

// Request ...
type Request struct {
	// ID reqeust id
	ID int64 `json:"rid" validate:"required"`
	// Command request command string
	Command string `json:"command" validate:"required"`
	// Payload interface
	Payload interface{} `json:"payload" validate:"required"`
	// binary representation of payload json
	payload []byte
}

// GetPayload ..
func (r *Request) GetPayload(payload interface{}) error {
	return json.Unmarshal(r.payload, payload)
}
