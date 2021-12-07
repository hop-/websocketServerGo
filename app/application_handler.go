package app

import (
	"wss/libs/websocket"

	log "github.com/hop-/golog"
)

// ApplicationHandler ..
type ApplicationHandler struct {
	commandHandlers map[string]CommandHandler
}

// NewApplicationHandler create a pointer to new ApplicationHandler object
func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		commandHandlers: make(map[string]CommandHandler),
	}
}

// AddCommandHandler ..
func (h *ApplicationHandler) AddCommandHandler(command string, ch CommandHandler) {
	h.commandHandlers[command] = ch
}

// handleWebsockets handle function for websocket server
func (h *ApplicationHandler) handleWebsockets(conn *websocket.Connection) {
	// create a new session
	s := NewSession(conn)

	// read session requests
	h.readRequests(s)

	// unsubscribe session from everything
	Unsubscribe(s, nil)
	// close session
	s.Close()
}

func (h *ApplicationHandler) readRequests(s *Session) {
	for {
		// Waiting for new request
		req, err := s.Read()

		if err != nil {
			// Check if connection fails
			if _, ok := err.(*websocket.CloseError); ok {
				log.Warning(err.Error())
				break
			}
			log.Error(err.Error())
			if err := sendErrorResponse(err.Error(), -1, s); err != nil {
				log.Error("Unable to send response, breaking connection")
				break
			}
			continue
		}

		// Get command handler
		commandHandler, ok := h.commandHandlers[req.Command]
		if !ok {
			errorMsg := "Unsupported command: " + req.Command
			log.Error(errorMsg)
			sendErrorResponse(errorMsg, req.ID, s)
			continue
		}

		log.Info("New command: ", req.Command)

		// Validate request/payload
		if err := commandHandler.Validate(req); err != nil {
			log.Error(err.Error())
			sendErrorResponse(err.Error(), req.ID, s)
			continue
		}

		// Handle command
		response, err := commandHandler.Handle(s, req)
		if err != nil {
			log.Error(err.Error())
			sendErrorResponse(err.Error(), req.ID, s)
			continue
		}
		s.SendResponse(response)
	}
}

// sendErrorResponse send response with error status for given session and request id
func sendErrorResponse(msg string, rID int64, s *Session) error {
	errorObject := NewError(msg)
	errorResponse := websocket.ErrorResponse(rID, errorObject)

	return s.SendResponse(errorResponse)
}
