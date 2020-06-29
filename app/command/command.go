package command

import (
	"../../app"
)

var (
	// UserService an app.UserService used in command handlers
	UserService app.UserService

	// NotificationService an app.NotificationService used in command handlers
	NotificationService app.NotificationService

	// ScheduleServiceMobile an app.ScheduleService used in command handlers
	ScheduleServiceMobile app.ScheduleService

	// TournamentServiceMobile an app.TournamentService used in command handlers
	TournamentServiceMobile app.TournamentService
)
