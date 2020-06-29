package websocket

// Response structure
type Response struct {
	// RequestID request id
	RequestID int64 `json:"rid"`
	// Event name string
	Event string `json:"event"`
	// Status response status
	Status string `json:"status"`
	// Payload response payload object
	Payload interface{} `json:"payload"`
}

// NewResponse create new response objcet
func NewResponse(rID int64, eventName string, status string, payload interface{}) *Response {
	return &Response{
		RequestID: rID,
		Event:     eventName,
		Status:    status,
		Payload:   payload,
	}
}

// ErrorResponse create new error response
func ErrorResponse(rID int64, payload interface{}) *Response {
	return &Response{
		RequestID: rID,
		Status:    "error",
		Payload:   payload,
	}
}
