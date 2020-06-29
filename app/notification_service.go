package app

// NotificationService ..
type NotificationService interface {
	GetNotifications(token string, limit, skip int, sort string, types []string) (interface{}, string, error)
	GetNotificationCount(token string, types []string, status string) (interface{}, string, error)
	UpdateNotifications(token string, ids []string, status string) (interface{}, string, error)
}
