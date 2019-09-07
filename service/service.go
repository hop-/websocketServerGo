package service

// Init services
func Init(requestPoolSize int, notificationURL string, userURL string) {
	semaphore = make(chan struct{}, requestPoolSize)
	apiNotificationURL = notificationURL
	apiUserURL = userURL
}
