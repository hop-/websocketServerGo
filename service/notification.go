package service

import (
	"encoding/json"
	"strconv"

	"../app"
	"../websocket"
)

var (
	apiNotificationURL string
)

// RequestNotifications ..
func RequestNotifications(rID int64,
	s *app.Session,
	eventName string,
	limit int,
	skip int,
	sort string,
	types []string,
) error {
	headers := map[string]string{"Authorization": s.AuthToken, "rid": strconv.FormatInt(rID, 10)}
	qParams := map[string][]string{}

	if limit != 0 {
		qParams["limit"] = []string{strconv.Itoa(limit)}
	}
	if skip != 0 {
		qParams["skip"] = []string{strconv.Itoa(skip)}
	}
	if sort != "" {
		qParams["sort"] = []string{sort}
	}
	qParams["type"] = types

	result, err := makeRequest("GET", apiNotificationURL, headers, qParams, nil)
	if err != nil {
		return err
	}

	resp := websocket.NewResponse(rID, eventName, result.Status, nil)
	if err = json.Unmarshal(result.Body, &resp.Payload); err != nil {
		return err
	}

	app.SendResponse(s.ID, rID, resp)
	return nil
}

// RequestNotificationsCount ..
func RequestNotificationsCount(rID int64, s *app.Session, eventName string, notificationTypes []string, notifiacationStatus string) error {
	headers := map[string]string{"Authorization": s.AuthToken, "rid": strconv.FormatInt(rID, 10)}
	qParams := map[string][]string{}

	qParams["type"] = notificationTypes
	qParams["status"] = []string{notifiacationStatus}

	result, err := makeRequest("GET", apiNotificationURL, headers, qParams, nil)
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

	respPayload := struct {
		Total int `json:"total"`
	}{}
	if err = json.Unmarshal(result.Body, &respPayload); err != nil {
		return err
	}
	resp.Payload = respPayload

	app.SendResponse(s.ID, rID, resp)
	return nil
}

// UpdateNotifications ..
func UpdateNotifications(rID int64, s *app.Session, eventName string, ids []string, status string) error {
	headers := map[string]string{"Authorization": s.AuthToken, "rid": strconv.FormatInt(rID, 10)}
	body := struct {
		Ids    []string `json:"ids,omitempty"`
		Status string   `json:"status"`
	}{ids, status}

	result, err := makeRequest("PUT", apiNotificationURL, headers, nil, body)
	if err != nil {
		return err
	}

	resp := websocket.NewResponse(rID, eventName, result.Status, nil)
	if result.Status != "ok" {
		if err = json.Unmarshal(result.Body, &resp.Payload); err != nil {
			return err
		}
	}
	app.SendResponse(s.ID, rID, resp)

	return nil
}
