package app

// SubscriptionService ..
type SubscriptionService interface {
	Save(sub *Subscription) (interface{}, string, error)
	Delete(hash string) (interface{}, string, error)
}
