package command

import (
	"../app"
	"../service"
	"../websocket"
)

// NotificationListHandler handler struct for notification command
type NotificationListHandler struct{}

// notification list is a golang representation of notification request payload object
type notificationOptions struct {
	Limit int      `json:"limit"`
	Skip  int      `json:"skip"`
	Sort  string   `json:"sort"`
	Types []string `json:"types"  validate:"required"`
}

// NewNotificationListHandler create a pointer to new notification object
func NewNotificationListHandler() *NotificationListHandler {
	return &NotificationListHandler{}
}

// Handle implementation of app.CommandHandler interface
func (h *NotificationListHandler) Handle(s *app.Session, r *websocket.Request) error {
	notificationOptions := notificationOptions{}
	r.GetPayload(&notificationOptions)

	return service.RequestNotifications(r.ID,
		s,
		r.Command,
		notificationOptions.Limit,
		notificationOptions.Skip,
		notificationOptions.Sort,
		notificationOptions.Types,
	)
}

// Validate implementation of app.CommandHandler interface
func (h *NotificationListHandler) Validate(r *websocket.Request) error {
	notificationOptions := notificationOptions{}
	r.GetPayload(&notificationOptions)

	return validatePayload(notificationOptions)
}
