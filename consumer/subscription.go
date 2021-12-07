package consumer

import (
	"encoding/json"

	"wss/app"
	"wss/libs/kafka"

	log "github.com/hop-/golog"
)

// SubscriptionsDocument is a golang representation of notification document object in kafka message
type SubscriptionsDocument struct {
	Subscriptions []string `json:"subscriptions"`
	Type          string   `json:"type"`
	Event         string   `json:"evnet"`
	Options       struct {
		Reload bool `json:"reload"`
	} `json:"options"`
	// Data structure may very depending on the Type of Document
	// json will unmarshall to map[string]interface{}
	Data interface{} `json:"data"`
}

// GetData structured data object
func (d *SubscriptionsDocument) GetData(i interface{}) error {
	data, err := json.Marshal(d.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, i)
}

// HandleSubscriptions handle subscriptions messages in kafka consumer
func HandleSubscriptions(msg *kafka.Message) {
	var document SubscriptionsDocument
	if err := msg.GetDocument(&document); err != nil {
		log.Error(err.Error())
		return
	}

	// Send update response to each subscription id
	for _, subID := range document.Subscriptions {
		app.SendUpdateMessage(subID, document.Data, document.Options.Reload)
	}
}
