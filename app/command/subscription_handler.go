package command

import (
	"../../app"
	"../../libs/websocket"
)

// SubscriptionHandler app.CommandHandler struct for subscription command
type SubscriptionHandler struct{}

// subscriptionOptions options object of notification request/command
type subscriptionOptions struct {
	IDs    []string `json:"ids" validate:"required"`
	Action string   `json:"action" validate:"required,oneof=unsubscribe"`
}

// NewSubscriptionHandler create a pointer to new SubscriptionHandler object
func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{}
}

// Handle method implementation
func (h *SubscriptionHandler) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	options := subscriptionOptions{}
	r.GetOptions(&options)

	// TODO support options.Action if needed
	app.Unsubscribe(s, options.IDs)

	response := websocket.NewResponse(r.ID, r.Command, "ok", nil)

	return response, nil
}

// Validate method implementation
func (h *SubscriptionHandler) Validate(r *websocket.Request) error {
	options := subscriptionOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
