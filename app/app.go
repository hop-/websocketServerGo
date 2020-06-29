package app

import (
	"fmt"
	"strings"
	"sync"

	"../libs/log"
	"../libs/websocket"

	"github.com/google/uuid"
)

var (
	subscriptions       = make(map[string]map[*Session]bool)
	queryObjects        = make(map[string]*Subscription)
	subscriptionsMutex  sync.RWMutex
	subscriptionService SubscriptionService

	// ID clinet runtime unique id
	ID = strings.Replace(uuid.New().String(), "-", "", -1)
)

// Serve websocket
func Serve(port int) {
	log.Info("Starting listening on port", port)
	websocket.Serve(port)
}

// AddHandler ..
func AddHandler(endpoint string, h *ApplicationHandler) {
	websocket.SetHandlerFunc(endpoint, h.handleWebsockets)
}

// SetSubscriptionService ..
func SetSubscriptionService(s SubscriptionService) {
	subscriptionService = s
}

// Subscribe is a fucntion to generate/get sub-id and subscribe session
func Subscribe(s *Session, resource string, where interface{}, recallhandlerFunc func(s *Session, where interface{}) (interface{}, error)) (string, error) {
	subscription := NewSubscription(resource, where, recallhandlerFunc)

	_, status, err := subscriptionService.Save(subscription)
	if err != nil {
		log.Error("Failed to save subscription:", err.Error())
		return "", err
	}

	if status != "ok" {
		// TODO use payload and status
		return "", fmt.Errorf("Failed to save subscription")
	}

	subscribe(s, subscription)

	return subscription.ID, nil
}

// subscribe a session to subscription ids
func subscribe(s *Session, sub *Subscription) {
	subscriptionsMutex.Lock()

	if _, ok := subscriptions[sub.ID]; !ok {
		subscriptions[sub.ID] = make(map[*Session]bool)
	}
	subscriptions[sub.ID][s] = true

	// TODO: the if statment can be removed if queryObjects and subscriptions are sync
	if _, ok := queryObjects[sub.ID]; !ok {
		queryObjects[sub.ID] = sub
	}

	subscriptionsMutex.Unlock()

	s.Subscribe([]string{sub.ID})
}

// Unsubscribe session from given subscription ids, if ids is nil unsubscribe from all
func Unsubscribe(s *Session, ids []string) {
	subscriptionIDs := ids
	// if ids array is nil remove all
	if ids == nil {
		subscriptionIDs = s.Subscriptions()
	}

	subscriptionsMutex.Lock()
	defer subscriptionsMutex.Unlock()

	for _, id := range subscriptionIDs {
		delete(subscriptions[id], s)

		if len(subscriptions[id]) == 0 {
			log.Info("Last subscriber has been unsubscribed. Deleting record")
			delete(subscriptions, id)
			delete(queryObjects, id)

			_, _, err := subscriptionService.Delete(id)
			if err != nil {
				log.Error("Failed to delete subscription:", err.Error())
			}
		}
	}

	s.Unsubscribe(ids)
}

// SendUpdateMessage to subscribers
func SendUpdateMessage(id string, data interface{}, reload bool) {
	if reload {
		reloadAndSendMessages(id)
		return
	}

	sendUpdateMessages(id, data)
}

func reloadAndSendMessages(id string) {
	subscriptionsMutex.RLock()
	defer subscriptionsMutex.RUnlock()

	subscriptionObject, ok := queryObjects[id]
	if !ok {
		log.Infof("Subscription with %s id doesn't exist, skipping", id)
		return
	}
	sessions, ok := subscriptions[id]
	if !ok {
		log.Infof("Subscriptions with %s id doesn't exist, skipping", id)
		return
	}
	for s := range sessions {
		payload, err := subscriptionObject.recallhandlerFunc(s, subscriptionObject.Where)
		if err != nil {
			continue
		}

		response := websocket.NewResponse(-1, id, "reload", payload)

		s.SendResponse(response)
	}
}

func sendUpdateMessages(id string, data interface{}) {
	subscriptionsMutex.RLock()
	defer subscriptionsMutex.RUnlock()

	sessions, ok := subscriptions[id]
	if !ok {
		log.Infof("Subscriptions with %s id doesn't exist, skipping", id)
		return
	}

	response := websocket.NewResponse(-1, id, "update", data)

	for s := range sessions {
		s.SendResponse(response)
	}
}
