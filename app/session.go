package app

import (
	"time"

	"../websocket"
)

// SubscriptionID is an type alias to string
type SubscriptionID = string

const (
	// NotificationID SubscriptionID
	NotificationID SubscriptionID = "notifications"
	// SeriesScoreID SubscriptionID
	SeriesScoreID SubscriptionID = "seriesScores"
)

// Session stucture to save session info and connection
type Session struct {
	// ID is a session id
	ID int64
	// conn is a websocket connection
	conn *websocket.Connection
	// UserID an id of user group
	UserID string
	// Subscriptions is a map of subscriptions whom was subuscribed the session
	Subscriptions map[SubscriptionID]bool
	// AuthTocken is a current authentication token which is neccessary for backand requests
	AuthToken string
}

// NewSession create pointer of new session object
func NewSession(conn *websocket.Connection) *Session {
	return &Session{
		ID:            time.Now().UnixNano(),
		conn:          conn,
		Subscriptions: make(map[SubscriptionID]bool),
	}
}

// Read a request from the websocket connection
func (s *Session) Read() (*websocket.Request, error) {
	return s.conn.Read()
}

// Write a responce into the websocket
func (s *Session) Write(response websocket.Response) error {
	return s.conn.Write(response)
}

// Subscribe add subscriptions
func (s *Session) Subscribe(list []SubscriptionID) {
	// if list is nil do nothing
	if list == nil {
		return
	}
	for _, id := range list {
		s.Subscriptions[id] = true
	}
}

// Unsubscribe remove subscriptions
func (s *Session) Unsubscribe(list []SubscriptionID) {
	// if list is nil remove all
	if list == nil {
		s.ClearSubscriptions()
		return
	}
	for _, id := range list {
		delete(s.Subscriptions, id)
	}
}

// ClearSubscriptions clear all subscriptions
func (s *Session) ClearSubscriptions() {
	s.Subscriptions = make(map[SubscriptionID]bool)
}
