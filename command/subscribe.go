package command

import (
	"../app"
	"../service"
	"../websocket"
)

type subscribePayload struct {
	Subscribe bool     `json:"subscribe" validate:"required"`
	List      []string `json:"list" validate:"dive,oneof=notifications seriesScores"`
}

// SubscribeHandler ..
type SubscribeHandler struct{}

// NewSubscribeHandler ..
func NewSubscribeHandler() *SubscribeHandler {
	return &SubscribeHandler{}
}

// Handle implementation of app.CommandHandler interface
func (h *SubscribeHandler) Handle(s *app.Session, r *websocket.Request) error {
	payload := subscribePayload{}
	r.GetPayload(&payload)

	service.Subscribe(s, r.ID, r.Command, payload.Subscribe, payload.List)
	return nil
}

// Validate implementation of app.CommandHandler interface
func (h *SubscribeHandler) Validate(r *websocket.Request) error {
	payload := subscribePayload{}
	r.GetPayload(&payload)

	return validatePayload(payload)
}
