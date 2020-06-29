package command

import (
	"errors"
	"fmt"

	"../../app"
	"../../libs/websocket"
)

// ScheduleHandlerMobile app.CommandHandler struct for schedule command
type ScheduleHandlerMobile struct{}

type scheduleOptionsWhere struct {
	From     uint64 `json:"from" validate:"required"`
	To       uint64 `json:"to" validate:"required"`
	Timezone string `json:"timezone" validate:"required"`
}

// scheduleOptions options object of schedule request/command
type scheduleOptions struct {
	Where scheduleOptionsWhere `json:"where" validate:"required"`
}

// NewScheduleHandlerMobile create a pointer to new Schedulendler object
func NewScheduleHandlerMobile() *ScheduleHandlerMobile {
	return &ScheduleHandlerMobile{}
}

// Handle method implementation
func (h *ScheduleHandlerMobile) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	options := scheduleOptions{}
	r.GetOptions(&options)

	payload, status, err := ScheduleServiceMobile.GetSchedules(options.Where.From, options.Where.To, options.Where.Timezone)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)
	if status == "ok" {
		subID, err := app.Subscribe(s, "schedule", options.Where, scheduleRecallHandler)
		if err != nil {
			return nil, errors.New("Failed to subscribe")
		}

		response.Event = subID
	}

	return response, nil
}

func scheduleRecallHandler(s *app.Session, where interface{}) (interface{}, error) {
	whereObject := where.(scheduleOptionsWhere)

	payload, status, err := ScheduleServiceMobile.GetSchedules(whereObject.From, whereObject.To, whereObject.Timezone)
	if err != nil {
		return nil, err
	}

	if status != "ok" {
		return nil, fmt.Errorf("Some internal error on recall with status: %s", status)
	}

	return payload, nil
}

// Validate method implementation
func (h *ScheduleHandlerMobile) Validate(r *websocket.Request) error {
	options := scheduleOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
