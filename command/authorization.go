package command

import (
	"../app"
	"../service"
	"../websocket"
)

// AuthorizationHandler handler struct for authorize command
type AuthorizationHandler struct{}

// authorizationOptions payload object of authorize request/command
type authorizationOptions struct {
	Update    bool   `json:"update"`
	AuthToken string `json:"authToken"`
}

// NewAuthorizationHandler create a pointer to new AutorizationHandler object
func NewAuthorizationHandler() *AuthorizationHandler {
	return &AuthorizationHandler{}
}

// Handle implementation of app.CommandHandler interface
func (h *AuthorizationHandler) Handle(s *app.Session, r *websocket.Request) error {
	authOptions := authorizationOptions{}
	r.GetPayload(&authOptions)

	// Login user as guest if payload is empty
	if authOptions == (authorizationOptions{}) {
		service.Logout(r.ID, s, r.Command)
		return nil
	}

	err := service.RequestAuthorization(r.ID, s, r.Command, authOptions.Update, authOptions.AuthToken)
	if err != nil {
		return err
	}
	return nil
}

// Validate implementation of app.CommandHandler interface
func (h *AuthorizationHandler) Validate(r *websocket.Request) error {
	authOptions := authorizationOptions{}
	r.GetPayload(&authOptions)
	return validatePayload(authOptions)
}
