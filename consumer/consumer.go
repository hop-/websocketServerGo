package consumer

import (
	"../kafka"
)

var ()

// Init initiate kafka consumer
func Init() {
	kafka.AddConsumerHandler("create_notification", handleNotification)
}

// Consume starts consuming in background
func Consume(brokers []string, groupID string, topic string, partition int) {
	go kafka.Consume(brokers, groupID, topic, partition)
}
