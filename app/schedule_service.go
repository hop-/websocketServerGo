package app

// ScheduleService ..
type ScheduleService interface {
	GetSchedules(from, to uint64, timezone string, gameID []string, substageTier []int, tournamentID []string) (interface{}, string, error)
}
