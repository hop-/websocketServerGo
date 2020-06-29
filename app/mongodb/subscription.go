package mongodb

import (
	"time"

	"../../app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscription is a bson ready structure for app.Subscription
type Subscription struct {
	Resource  string             `bson:"resource"`
	ID        string             `bson:"subscriptionId"`
	ClientID  string             `bson:"clientId"`
	Where     interface{}        `bson:"where"`
	UpdatedAt primitive.DateTime `bson:"updatedAt"`
}

// NewSubscription ..
func NewSubscription(s *app.Subscription) Subscription {
	return Subscription{
		Resource:  s.Resource,
		ID:        s.ID,
		Where:     s.Where,
		ClientID:  app.ID,
		UpdatedAt: primitive.DateTime(time.Now().Unix()),
	}
}
