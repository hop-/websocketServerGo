package kafka

import (
	"context"
	"encoding/json"
	"time"

	"../log"

	kafkago "github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/gzip"
)

var (
	handlers = make(map[string]func(*Message))
)

// handleMessage handles kafka message in consumer
func handleMessage(msg *kafkago.Message) {
	message := Message{}
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Create document []byte from Document object
	message.document, err = json.Marshal(message.Document)
	if err != nil {
		log.Error(err.Error())
	}

	// Handling message by its type
	if handler, ok := handlers[message.Type]; ok {
		log.Info("Processing message of type:", message.Type)
		handler(&message)
	} else {
		log.Error("Unsupported message type:", message.Type)
	}
}

// AddConsumerHandler ..
func AddConsumerHandler(typeName string, handler func(*Message)) {
	handlers[typeName] = handler
}

// Consume ..
func Consume(brokers []string, groupID string, topic string, partition int) {
	// Creating reader with given config
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   brokers,
		GroupID:   groupID,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	// Setting consummer/reader offset to last message
	r.SetOffsetAt(context.Background(), time.Now())
	defer r.Close()

	// Consuming loop
	for {
		// Wait for new message
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Error("Kafka -", err.Error())
			continue
		}
		log.Info("Received new message from kafka")
		// Handle message
		handleMessage(&msg)
	}
}
