package app

// TournamentService ..
type TournamentService interface {
	GetTournament(id string) (interface{}, string, error)
}
