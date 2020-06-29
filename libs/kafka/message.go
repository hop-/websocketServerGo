package kafka

import "encoding/json"

// Message ..
type Message struct {
	// Type of message
	Type string `json:"type"`
	// Document contains all data
	Document interface{}
	// Binary representation of json document
	document []byte
}

// GetDocument get document as given interface
func (m *Message) GetDocument(i interface{}) error {
	return json.Unmarshal(m.document, i)
}
