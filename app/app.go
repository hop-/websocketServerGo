package app

import (
	"../log"
	"../websocket"
)

// Init ..
func Init() {
	websocket.SetHandlerFunc(handleWebsockets)
}

// Serve websocket
func Serve(port int) {
	log.Info("Starting listening on port", port)
	websocket.Serve(port)
}
