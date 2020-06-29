package rest

import (
	"encoding/json"
)

// UserService rest implementation for app.UserService
type UserService struct {
	url string
}

// NewUserService create a pointer to new UserService object
func NewUserService(url string) *UserService {
	return &UserService{url}
}

// Login method implementation
func (s *UserService) Login() (interface{}, string, error) {
	// TODO
	return nil, "ok", nil
}

// Logout method implementation
func (s *UserService) Logout() (interface{}, string, error) {
	// TODO
	return nil, "ok", nil
}

// Authenticate method implementation
func (s *UserService) Authenticate(token string) (interface{}, string, error) {
	headers := map[string]string{"Authorization": token}

	result, err := request("GET", s.url, headers, nil, nil)
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
	info := make(map[string]interface{})
	if err = json.Unmarshal(result.Body, &info); err != nil {
		return nil, "error", err
	}

	return info, result.Status, nil
}
