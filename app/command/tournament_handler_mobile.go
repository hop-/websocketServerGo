package command

import (
	"../../app"
	"../../libs/websocket"
)

// TournamentHandlerMobile app.CommandHandler struct for tournament command
type TournamentHandlerMobile struct{}

type tournamentOptionsWhere struct {
	ID string `json:"id"`
}

// tournamentOptions options object of tournament request/command
type tournamentOptions struct {
	Where tournamentOptionsWhere `json:"where"`
}

// NewTournamentHandlerMobile create a pointer to new Schedulendler object
func NewTournamentHandlerMobile() *TournamentHandlerMobile {
	return &TournamentHandlerMobile{}
}

// Handle method implementation
func (h *TournamentHandlerMobile) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	options := tournamentOptions{}
	r.GetOptions(&options)

	payload, status, err := TournamentServiceMobile.GetTournament(options.Where.ID)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)

	return response, nil
}

// Validate method implementation
func (h *TournamentHandlerMobile) Validate(r *websocket.Request) error {
	options := tournamentOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
