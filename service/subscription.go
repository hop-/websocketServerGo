package service

import (
	"../app"
	"../log"
	"../websocket"
)

// Subscribe ..
func Subscribe(s *app.Session, rID int64, event string, subscribe bool, list []app.SubscriptionID) {
	if subscribe {
		log.Info("Subscribe for", list)
		s.Subscribe(list)
	} else {
		log.Info("Unsubscribe from", list)
		s.Unsubscribe(list)
	}

	response := websocket.NewResponse(rID, event, "ok", nil)
	app.SendResponse(s.ID, rID, response)
}
