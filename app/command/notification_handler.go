package command

import (
	"errors"
	"fmt"

	"wss/app"
	"wss/libs/websocket"
)

// NotificationHandler app.CommandHandler struct for notification command
type NotificationHandler struct{}

type notificationOptionsWhere struct {
	Types []string `json:"types"  validate:"required"`
}

// notificationOptions options object of notification request/command
type notificationOptions struct {
	Params struct {
		Limit int    `json:"limit"`
		Skip  int    `json:"skip"`
		Sort  string `json:"sort"`
	} `json:"params"`
	Where notificationOptionsWhere `json:"where" validate:"required"`
}

// NewNotificationHandler create a pointer to new Notificationhandler object
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// Handle method implementation
func (h *NotificationHandler) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	notificationOptions := notificationOptions{}
	r.GetOptions(&notificationOptions)

	payload, status, err := NotificationService.GetNotifications(
		s.AuthToken,
		notificationOptions.Params.Limit,
		notificationOptions.Params.Skip,
		notificationOptions.Params.Sort,
		notificationOptions.Where.Types,
	)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)

	if status == "ok" {
		where := struct {
			UserID string   `json:"userId"`
			Types  []string `json:"types"`
		}{s.UserID, notificationOptions.Where.Types}

		subID, err := app.Subscribe(s, "notification", where, notificationRecallHandler)
		if err != nil {
			return nil, errors.New("Failed to subscribe")
		}
		response.Event = subID
	}

	return response, nil
}

func notificationRecallHandler(s *app.Session, where interface{}) (interface{}, error) {
	whereObject := where.(notificationOptionsWhere)

	payload, status, err := NotificationService.GetNotifications(s.AuthToken,
		0,
		0,
		"",
		whereObject.Types,
	)
	if err != nil {
		return nil, err
	}

	if status != "ok" {
		return nil, fmt.Errorf("Some internal error on recall with status: %s", status)
	}

	return payload, nil
}

// Validate method implementation
func (h *NotificationHandler) Validate(r *websocket.Request) error {
	options := notificationOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
