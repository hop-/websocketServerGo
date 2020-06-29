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
	// Options interface
	Options interface{} `json:"options" validate:"required"`
	// binary representation of options json
	options []byte
}

// GetOptions ..
func (r *Request) GetOptions(options interface{}) error {
	return json.Unmarshal(r.options, options)
}
