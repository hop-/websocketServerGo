package rest

import (
	"encoding/json"
)

// TournamentService rest implementation for app.TournamentService
type TournamentService struct {
	url string
}

// NewTournamentService create a pointer to new TournamentService object
func NewTournamentService(url string) *TournamentService {
	return &TournamentService{url}
}

// GetTournament method implementation
func (s *TournamentService) GetTournament(id string) (interface{}, string, error) {
	url := s.url + "/" + id

	result, err := request("GET", url, nil, nil, nil)
	if err != nil {
		return nil, "error", err
	}

	var payload interface{}
	if err = json.Unmarshal(result.Body, &payload); err != nil {
		return nil, "error", err
	}

	return payload, result.Status, nil
}
