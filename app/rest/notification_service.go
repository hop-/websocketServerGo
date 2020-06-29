package rest

import (
	"encoding/json"
	"strconv"
)

// NotificationService rest implementation for app.NotificationService
type NotificationService struct {
	url string
}

// NewNotificationService create a pointer to new NotificationService object
func NewNotificationService(url string) *NotificationService {
	return &NotificationService{url}
}

// GetNotifications method implementation
func (s *NotificationService) GetNotifications(token string, limit, skip int, sort string, types []string) (interface{}, string, error) {
	headers := map[string]string{"Authorization": token}
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

	result, err := request("GET", s.url, headers, qParams, nil)
	if err != nil {
		return nil, "error", err
	}

	var payload interface{}
	if err = json.Unmarshal(result.Body, &payload); err != nil {
		return nil, "error", err
	}

	return payload, result.Status, nil
}

// GetNotificationCount method implementation
func (s *NotificationService) GetNotificationCount(token string, types []string, status string) (interface{}, string, error) {
	headers := map[string]string{"Authorization": token}

	qParams := map[string][]string{}
	qParams["type"] = types
	qParams["status"] = []string{status}

	result, err := request("GET", s.url, headers, qParams, nil)
	if err != nil {
		return nil, "error", err
	}

	if result.Status != "ok" {
		var errorPayload interface{}
		if err = json.Unmarshal(result.Body, &errorPayload); err != nil {
			return nil, "error", err
		}

		return errorPayload, result.Status, nil
	}

	payload := struct {
		Total int `json:"total"`
	}{}

	if err = json.Unmarshal(result.Body, &payload); err != nil {
		return nil, "error", err
	}

	return payload, "ok", nil
}

// UpdateNotifications method implementation
func (s *NotificationService) UpdateNotifications(token string, ids []string, status string) (interface{}, string, error) {
	headers := map[string]string{"Authorization": token}

	body := struct {
		Ids    []string `json:"ids,omitempty"`
		Status string   `json:"status"`
	}{ids, status}

	result, err := request("PUT", s.url, headers, nil, body)
	if err != nil {
		return nil, "error", err
	}

	if result.Status != "ok" {
		var errorPayload interface{}
		if err = json.Unmarshal(result.Body, &errorPayload); err != nil {
			return nil, "error", err
		}
		return errorPayload, result.Status, nil
	}

	return nil, result.Status, nil
}
