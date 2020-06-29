package command

import (
	"../../app"
	"../../libs/websocket"
)

// UpdateNotificationHandler app.CommandHandler struct for updateNotification command
type UpdateNotificationHandler struct{}

// updateNotificationOptions options object of notification request/command
type updateNotificationOptions struct {
	Where struct {
		Notifications []string `json:"notifications" validate:"required"`
	} `json:"where" validate:"required"`
	Status string `json:"status" validate:"required,oneof=unseen unread read"`
}

// NewUpdateNotificationHandler create a pointer to new UpdateNotificationHandler object
func NewUpdateNotificationHandler() *UpdateNotificationHandler {
	return &UpdateNotificationHandler{}
}

// Handle method implementation
func (h *UpdateNotificationHandler) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	options := updateNotificationOptions{}
	r.GetOptions(&options)

	payload, status, err := NotificationService.UpdateNotifications(
		s.AuthToken,
		options.Where.Notifications,
		options.Status,
	)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)

	return response, nil
}

// Validate method implementation
func (h *UpdateNotificationHandler) Validate(r *websocket.Request) error {
	options := updateNotificationOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
