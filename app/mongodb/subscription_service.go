package mongodb

import (
	"context"

	"wss/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// SubscriptionService mongodb implementation for app.SubscriptionService
type SubscriptionService struct {
	url        string
	dbName     string
	ctx        context.Context
	clf        context.CancelFunc
	client     *mongo.Client
	collection *mongo.Collection
}

// NewSubscriptionService create a pointer to new SubscriptionduleService object
func NewSubscriptionService(url string, dbName string) *SubscriptionService {
	ctx, clf := context.WithCancel(context.Background())

	return &SubscriptionService{
		url:    url,
		dbName: dbName,
		ctx:    ctx,
		clf:    clf,
	}
}

// Connect ...
func (s *SubscriptionService) Connect() error {
	var err error
	// Connecting in background
	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(s.url))
	if err != nil {
		return err
	}

	s.collection = s.client.Database(s.dbName).Collection("subscriptions")

	return nil
}

// Save method implementation
func (s *SubscriptionService) Save(data *app.Subscription) (interface{}, string, error) {
	subscription := NewSubscription(data)

	filter := bson.M{"subscriptionId": subscription.ID, "clientId": subscription.ClientID}
	update := bson.M{"$set": subscription}

	updateOptions := &options.UpdateOptions{}
	updateOptions.SetUpsert(true)

	result, err := s.collection.UpdateOne(s.ctx, filter, update, updateOptions)
	if err != nil {
		return nil, "error", err
	}

	return result, "ok", nil
}

// Delete method implementation
func (s *SubscriptionService) Delete(subsciptionID string) (interface{}, string, error) {
	filter := bson.M{"subscriptionId": subsciptionID, "clientId": app.ID}

	result, err := s.collection.DeleteOne(s.ctx, filter)
	if err != nil {
		return nil, "error", err
	}

	return result, "ok", nil
}

// Disconnect disconnect mongo client
func (s *SubscriptionService) Disconnect() {
	s.client.Disconnect(s.ctx)
	s.clf()
}

// Ping check the connection status
func (s *SubscriptionService) Ping() {
	s.client.Ping(s.ctx, readpref.Primary())
}
