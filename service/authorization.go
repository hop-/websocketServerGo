package service

import (
	"encoding/json"
	"errors"
	"strconv"

	"../app"
	"../websocket"
)

var (
	apiUserURL string
)

// userInfo ..
type userInfo struct {
	UserID string `json:"id"`
}

// Logout session
func Logout(rID int64, s *app.Session, eventName string) {
	response := websocket.NewResponse(rID, eventName, "ok", nil)
	app.SendResponse(s.ID, rID, response)
	app.ChangeSessionUserID(s, "", "")
}

// RequestAuthorization ..
func RequestAuthorization(rID int64, s *app.Session, eventName string, update bool, authToken string) error {
	headers := map[string]string{"Authorization": authToken, "rid": strconv.FormatInt(rID, 10)}

	result, err := makeRequest("GET", apiUserURL, headers, nil, nil)
	if err != nil {
		return err
	}

	resp := websocket.NewResponse(rID, eventName, result.Status, nil)
	if result.Status != "ok" {
		if err = json.Unmarshal(result.Body, &resp.Payload); err != nil {
			return err
		}

		app.SendResponse(s.ID, rID, resp)
		return nil
	}

	info := userInfo{}
	if err = json.Unmarshal(result.Body, &info); err != nil {
		return err
	}

	if update {
		if info.UserID != s.UserID {
			return errors.New("Bad access token")
		}

		app.SendResponse(s.ID, rID, resp)
		app.ChangeSessionAuthToken(s, authToken)
		return nil
	}

	app.SendResponse(s.ID, rID, resp)
	app.ChangeSessionUserID(s, info.UserID, authToken)
	return nil
}
