package app

// ScheduleService ..
type ScheduleService interface {
	GetSchedules(from, to uint64, timezone string) (interface{}, string, error)
}
