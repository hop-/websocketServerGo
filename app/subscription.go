package app

import (
	"../libs/hasher"
)

// Subscription ..
type Subscription struct {
	Command           string      `json:"command"`
	Resource          string      `json:"resource"`
	ID                string      `json:"hash"`
	Where             interface{} `json:"where"`
	recallhandlerFunc func(s *Session, where interface{}) (interface{}, error)
}

// NewSubscription create a new Subscription object
func NewSubscription(resource string, where interface{}, recallhandlerFunc func(s *Session, where interface{}) (interface{}, error)) *Subscription {
	s := &Subscription{
		Resource:          resource,
		Where:             where,
		recallhandlerFunc: recallhandlerFunc,
	}
	s.ID, _ = hasher.HashJSONObject(s)

	return s
}
