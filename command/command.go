package command

import (
	"../app"
)

// Setup handler routes
func Setup() {
	// Route command handlers to command strings
	app.AddCommandHandler("notificationList", NewNotificationListHandler())
	app.AddCommandHandler("authorize", NewAuthorizationHandler())
	app.AddCommandHandler("unseenNotificationCount", NewUnseenNotificationCountHandler())
	app.AddCommandHandler("updateNotification", NewNotificationUpdateHandler())
	app.AddCommandHandler("subscribe", NewSubscribeHandler())
}
