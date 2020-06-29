package command

import (
	"errors"

	"../../app"
	"../../libs/websocket"
)

// AuthenticationHandler app.CommandHandler struct for authorize command
type AuthenticationHandler struct{}

// authenticationOptions options object of authorize request/command
type authenticationOptions struct {
	Update    bool   `json:"update"`
	AuthToken string `json:"authToken"`
}

// NewAuthenticationHandler create a pointer to new AutorizationHandler object
func NewAuthenticationHandler() *AuthenticationHandler {
	return &AuthenticationHandler{}
}

// Handle method implementation
func (h *AuthenticationHandler) Handle(s *app.Session, r *websocket.Request) (*websocket.Response, error) {
	authOptions := authenticationOptions{}
	r.GetOptions(&authOptions)

	// Login user as guest if options are empty
	if authOptions == (authenticationOptions{}) {
		response := websocket.NewResponse(r.ID, r.Command, "ok", nil)
		s.ChangeUser("", "")
		return response, nil
	}

	payload, status, err := UserService.Authenticate(authOptions.AuthToken)
	if err != nil {
		return nil, err
	}

	response := websocket.NewResponse(r.ID, r.Command, status, payload)

	if status != "ok" {
		return response, nil
	}

	info := payload.(map[string]interface{}) // todo ugly
	if authOptions.Update {
		if info["id"] != s.UserID {
			return nil, errors.New("Invalid access token")
		}

		s.AuthToken = authOptions.AuthToken
		return response, nil
	}

	s.ChangeUser(info["id"].(string), authOptions.AuthToken)

	return response, nil
}

// Validate method implementation
func (h *AuthenticationHandler) Validate(r *websocket.Request) error {
	options := authenticationOptions{}
	r.GetOptions(&options)

	return validateOptions(options)
}
