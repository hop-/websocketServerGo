package consumer

import (
	"encoding/json"
	"time"

	"../app"
	"../kafka"
	"../log"
)

// Document is a golang representation of notification document object in kafka message
type Document struct {
	ID     string    `json:"_id"`
	Owner  string    `json:"owner"`
	Type   string    `json:"type"`
	Status string    `json:"status"`
	SentAt time.Time `json:"sentAt"`
	// Data structure may very depending on the Type of Document
	// json will unmarshall to map[string]interface{}
	Data interface{} `json:"data"`
}

// GetData structured data object
func (d *Document) GetData(i interface{}) error {
	data, err := json.Marshal(d.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, i)
}

// handleNotification handle notification messages in kafka consumer
func handleNotification(msg *kafka.Message) {
	var document Document
	if err := msg.GetDocument(&document); err != nil {
		log.Error(err.Error())
		return
	}

	user := document.Owner
	app.PushNotification(user, app.NotificationID, document)
}
