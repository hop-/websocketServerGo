package command

import (
	"../app"
	"../service"
	"../websocket"
)

// NotificationUpdateHandler ..
type NotificationUpdateHandler struct{}

// notificationUpdateOptions (payload)
type notificationUpdateOptions struct {
	Notifications []string `json:"ids,omitempty"`
	Status        string   `json:"status" validate:"required,oneof=unseen unread read"`
}

// NewNotificationUpdateHandler ..
func NewNotificationUpdateHandler() *NotificationUpdateHandler {
	return &NotificationUpdateHandler{}
}

// Handle implementation of app.CommandHandler interface
func (h *NotificationUpdateHandler) Handle(s *app.Session, r *websocket.Request) error {
	notifUpdateOptions := notificationUpdateOptions{}
	r.GetPayload(&notifUpdateOptions)

	return service.UpdateNotifications(r.ID,
		s,
		r.Command,
		notifUpdateOptions.Notifications,
		notifUpdateOptions.Status,
	)
}

// Validate implementation of app.CommandHandler interface
func (h *NotificationUpdateHandler) Validate(r *websocket.Request) error {
	notifUpdateOptions := notificationUpdateOptions{}
	r.GetPayload(&notifUpdateOptions)

	return validatePayload(notifUpdateOptions)
}
