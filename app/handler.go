package app

import (
	"fmt"

	"../log"
	"../websocket"
)

var (
	sessions = make(map[string]map[*Session]bool)
	requests = make(map[int64]map[int64]*Session)
)

// handleWebsockets handle function for websocket server
func handleWebsockets(conn *websocket.Connection) {
	s := newSession(conn)

	readRequests(s)
	removeSession(s)
}

func readRequests(s *Session) {
	for {
		// Waiting for new request
		req, err := s.Read()

		if err != nil {
			// Check if connection fails
			if _, ok := err.(*websocket.CloseError); ok {
				log.Error(err.Error())
				break
			}
			log.Error(err.Error())
			SendErrorResponse(err.Error(), -1, s)
			continue
		}

		// Get command handler
		commandHandler, ok := commandHandlers[req.Command]
		if !ok {
			errorMsg := "Unsupported command: " + req.Command
			log.Error(errorMsg)
			SendErrorResponse(errorMsg, req.ID, s)
			continue
		}

		log.Info("New command: ", req.Command)

		// Validate request/payload
		if err := commandHandler.Validate(req); err != nil {
			log.Error(err.Error())
			SendErrorResponse(err.Error(), req.ID, s)
			continue
		}

		// Submit reqeust in request list
		if err := SubmitRequest(s, req.ID); err != nil {
			log.Error(err.Error())
			SendErrorResponse(err.Error(), -1, s)
		}

		// Handle command
		if err := commandHandler.Handle(s, req); err != nil {
			log.Error(err.Error())
			SendErrorResponse(err.Error(), req.ID, s)
			continue
		}
	}
}

// newSession create new session in session list
func newSession(conn *websocket.Connection) *Session {
	s := NewSession(conn)

	addSession(s)
	ClearSessionRequests(s)

	return s
}

// addSession add session into the smecified user group
func addSession(s *Session) {
	if len(sessions[s.UserID]) == 0 {
		sessions[s.UserID] = make(map[*Session]bool)
	}

	sessions[s.UserID][s] = true
}

// removeSession removes session from user group
func removeSession(s *Session) {

	delete(sessions[s.UserID], s)
	if len(sessions[s.UserID]) == 0 {
		delete(sessions, s.UserID)
	}

	ClearSessionRequests(s)
}

// ChangeSessionUserID changes session user group
func ChangeSessionUserID(s *Session, userID, authToken string) {
	log.Infof("Change UserID from '%s' to '%s' session %d", s.UserID, userID, s.ID)
	removeSession(s)

	s.ClearSubscriptions()

	s.UserID = userID
	ChangeSessionAuthToken(s, authToken)

	addSession(s)
}

// ChangeSessionAuthToken changes session authentication token
func ChangeSessionAuthToken(s *Session, authToken string) {
	log.Infof("Change authToken for session %d", s.ID)
	s.AuthToken = authToken
}

// SubmitRequest submit new request for the given session
func SubmitRequest(s *Session, requestID int64) error {
	if _, ok := requests[s.ID][requestID]; ok {
		return fmt.Errorf("Duplicate rid: %d", requestID)
	}
	requests[s.ID][requestID] = s
	return nil
}

// ClearSessionRequests clear all submitted requests for given session
func ClearSessionRequests(s *Session) {
	requests[s.ID] = make(map[int64]*Session)
}

// PushNotification push response
func PushNotification(userID string, subscription SubscriptionID, payload interface{}) {
	r := websocket.NewResponse(-1, subscription, "ok", payload)

	if group, ok := sessions[userID]; ok {
		log.Info("Push notification of type:", subscription)

		for s := range group {
			if _, ok := s.Subscriptions[subscription]; !ok {
				continue
			}

			log.Info("Pushing for session:", s.ID)

			if err := s.Write(r); err != nil {
				log.Error(err.Error())
			}
		}
	} else {
		log.Info("Unable to find userGroup for userID:", userID)
	}
}

// SendResponse send response for given session and request id
func SendResponse(sessionID int64, requestID int64, r websocket.Response) {
	log.Infof("Sending response for %d request", requestID)

	s := requests[sessionID][requestID]
	if s == nil {
		log.Error("Unknown sessionID requestID: ", sessionID, requestID)
		return
	}

	r.RequestID = requestID
	if err := s.Write(r); err != nil {
		log.Error(err.Error())
	}

	log.Info("Sent successfully")

	delete(requests[sessionID], requestID)
}

// SendErrorResponse send response with error status for given session and request id
func SendErrorResponse(msg string, rID int64, s *Session) {
	errorObject := NewError(msg)
	errorResponse := websocket.ErrorResponse(rID, errorObject)

	if err := s.Write(errorResponse); err != nil {
		log.Error(err.Error())
	}
}
