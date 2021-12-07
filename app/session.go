package app

import (
	"fmt"
	"sync"
	"time"

	"wss/libs/websocket"

	log "github.com/hop-/golog"
)

// Session stucture to save session info and connection
type Session struct {
	// ID is a session id
	ID int64
	// conn is a websocket connection
	conn *websocket.Connection
	// subscriptions is a list/set of subscriptions for the session
	subscriptions map[string]bool
	// UserID an id of user group
	UserID string
	// AuthToken is a current authentication token which is neccessary for backand requests
	AuthToken string
	requests  map[int64]bool
	mutex     sync.RWMutex
}

// NewSession create pointer of new session object
func NewSession(conn *websocket.Connection) *Session {
	return &Session{
		ID:            time.Now().UnixNano(), // TODO check if this unique
		conn:          conn,
		subscriptions: make(map[string]bool),
		requests:      make(map[int64]bool),
	}
}

// Close session connection
func (s *Session) Close() {
	s.conn.Close()
}

// Read a request from the websocket connection
func (s *Session) Read() (*websocket.Request, error) {
	r, err := s.conn.Read()
	if err == nil {
		s.submitRequest(r.ID)
	}

	return r, err
}

// Write a responce into the websocket
func (s *Session) Write(response *websocket.Response) error {
	return s.conn.Write(response)
}

// submitRequest ..
func (s *Session) submitRequest(requestID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.requests[requestID]; ok {
		return fmt.Errorf("Duplicate rid: %d", requestID)
	}

	s.requests[requestID] = true
	return nil
}

// ClearRequests clear all submitted requests for given session
func (s *Session) ClearRequests() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requests = make(map[int64]bool)
}

// SendResponse send response for given session and request id
func (s *Session) SendResponse(r *websocket.Response) error {
	// TODO check if logs can be moved somwhere else
	if r.RequestID != -1 {
		log.Infof("Sending response for %d request", r.RequestID)
		s.mutex.RLock()
		_, ok := s.requests[r.RequestID]
		s.mutex.RUnlock()
		if !ok {
			log.Error("Unknown sessionID requestID: ", s.ID, r.RequestID)
			return nil
		}
	} else {
		log.Info("Sending update response without request id")
	}

	var err error
	if err = s.Write(r); err != nil {
		log.Error(err.Error())
	} else {
		log.Info("Sent successfully")
	}

	if r.RequestID != -1 {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		delete(s.requests, r.RequestID)
	}
	return err
}

// Subscriptions getter function for subscription list
func (s *Session) Subscriptions() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	subscriptions := make([]string, 0, len(s.subscriptions))
	for s := range s.subscriptions {
		subscriptions = append(subscriptions, s)
	}

	return subscriptions
}

// Subscribe add subscriptions
func (s *Session) Subscribe(ids []string) {
	// if ids is nil do nothing
	if ids == nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, id := range ids {
		s.subscriptions[id] = true
	}
}

// Unsubscribe remove subscriptions
func (s *Session) Unsubscribe(ids []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// if ids is nil remove all
	if ids == nil {
		s.subscriptions = make(map[string]bool)
		return
	}

	for _, id := range ids {
		delete(s.subscriptions, id)
	}
}

// ChangeUser ..
func (s *Session) ChangeUser(userID string, token string) {
	s.UserID = userID
	s.AuthToken = token
	s.Unsubscribe(nil)
}
