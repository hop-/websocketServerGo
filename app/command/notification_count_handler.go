package command

import (
	"wss/app"
	"wss/libs/websocket"
)

// NotificationCountHandler app.CommandHandler struct for notificationCount command
type NotificationCountHandler struct{}

// notificationCountOptions options object of notification request/command
type notificationCountOptions struct {
	Where struct {
		Types []string `json:"types" validate:"required"`
	} `json:"where" validate:"required"`
	Status string `json:"status" validate:"required,oneof=unseen unread read"`
}

// NewNotificationCountHandler create a pointer to new NotificationCountHandler object
func NewNotificationCountHandler() *NotificationCountHandler {
	return &NotificationCountHandler{}
}

// Handle method implementation
func (h *NotificationCountHandler) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	options := notificationCountOptions{}
	r.GetOptions(&options)

	payload, status, err := NotificationService.GetNotificationCount(
		s.AuthToken,
		options.Where.Types,
		options.Status,
	)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)

	return response, nil
}

// Validate method implementation
func (h *NotificationCountHandler) Validate(r *websocket.Request) error {
	options := notificationCountOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
