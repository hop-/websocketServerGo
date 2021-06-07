package command

import (
	"errors"
	"fmt"

	"wss/app"
	"wss/libs/websocket"
)

// ScheduleHandlerMobile app.CommandHandler struct for schedule command
type ScheduleHandlerMobile struct{}

type scheduleOptionsWhere struct {
	From     uint64   `json:"from" validate:"required"`
	To       uint64   `json:"to" validate:"required"`
	Timezone string   `json:"timezone" validate:"required"`
	GameIDs  []string `json:"gameIds,omitempty" bson:"gameIds,omitempty"`
	Test     *bool    `json:"test,omitempty" bson:"test,omitempty"`
	Substage *struct {
		Tiers []int `json:"tiers,omitempty" bson:"tiers,omitempty"`
	} `json:"substage,omitempty" bson:"substage,omitempty"`
	Tournament *struct {
		IDs []string `json:"ids,omitempty" bson:"ids,omitempty"`
	} `json:"tournament,omitempty" bson:"tournament,omitempty"`
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

	var substageTiers []int
	if options.Where.Substage != nil {
		substageTiers = options.Where.Substage.Tiers
	}

	var tournamentIDs []string
	if options.Where.Tournament != nil {
		tournamentIDs = options.Where.Tournament.IDs
	}

	payload, status, err := ScheduleServiceMobile.GetSchedules(options.Where.From, options.Where.To, options.Where.Timezone, options.Where.GameIDs, substageTiers, tournamentIDs)
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

	var substageTiers []int
	if whereObject.Substage != nil {
		substageTiers = whereObject.Substage.Tiers
	}

	var tournamentIDs []string
	if whereObject.Tournament != nil {
		tournamentIDs = whereObject.Tournament.IDs
	}

	payload, status, err := ScheduleServiceMobile.GetSchedules(whereObject.From, whereObject.To, whereObject.Timezone, whereObject.GameIDs, substageTiers, tournamentIDs)
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
