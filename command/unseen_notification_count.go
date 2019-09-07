package command

import (
	"../app"
	"../service"
	"../websocket"
)

// UnseenNotificationCountHandler ..
type UnseenNotificationCountHandler struct{}

// unseenNotificationCountOptions ...
type unseenNotificationCountOptions struct {
	Types []string `json:"types" validate:"required"`
}

// NewUnseenNotificationCountHandler ..
func NewUnseenNotificationCountHandler() *UnseenNotificationCountHandler {
	return &UnseenNotificationCountHandler{}
}

// Handle implementation of app.CommandHandler interface
func (h *UnseenNotificationCountHandler) Handle(s *app.Session, r *websocket.Request) error {
	unseenCountOptions := unseenNotificationCountOptions{}
	r.GetPayload(&unseenCountOptions)

	return service.RequestNotificationsCount(r.ID,
		s,
		r.Command,
		unseenCountOptions.Types,
		"unseen",
	)
}

// Validate implementation of app.CommandHandler interface
func (h *UnseenNotificationCountHandler) Validate(r *websocket.Request) error {
	unseenCountOptions := unseenNotificationCountOptions{}
	r.GetPayload(&unseenCountOptions)

	return validatePayload(unseenCountOptions)
}
