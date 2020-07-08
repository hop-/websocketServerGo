package rest

import (
	"encoding/json"
	"strconv"
)

// ScheduleService rest implementation for app.ScheduleService
type ScheduleService struct {
	url string
}

// NewScheduleService create a pointer to new ScheduleService object
func NewScheduleService(url string) *ScheduleService {
	return &ScheduleService{url}
}

// GetSchedules method implementation
func (s *ScheduleService) GetSchedules(from, to uint64, timezone string, gameID []string, substageTier []int, tournamentID []string) (interface{}, string, error) {
	headers := make(map[string]string)
	qParams := map[string][]string{}

	qParams["from"] = []string{strconv.FormatUint(from, 10)}
	qParams["to"] = []string{strconv.FormatUint(to, 10)}
	qParams["tz"] = []string{timezone}

	if gameID != nil {
		qParams["game_id"] = gameID
	}
	if substageTier != nil {
		tiers := make([]string, len(substageTier))

		for i := range substageTier {
			tiers[i] = strconv.Itoa(substageTier[i])
		}

		qParams["substage.tier"] = tiers
	}
	if tournamentID != nil {
		qParams["tournament.id"] = tournamentID
	}

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
