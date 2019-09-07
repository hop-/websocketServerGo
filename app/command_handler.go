package app

import (
	"../websocket"
)

var (
	commandHandlers = make(map[string]CommandHandler)
)

// CommandHandler ..
type CommandHandler interface {
	Handle(s *Session, r *websocket.Request) error
	Validate(r *websocket.Request) error
}

// AddCommandHandler ..
func AddCommandHandler(name string, h CommandHandler) {
	commandHandlers[name] = h
}
