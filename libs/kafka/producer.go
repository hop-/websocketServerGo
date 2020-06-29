package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

// Producer ..
type Producer struct {
	producer *kafkago.Writer
}

// NewProducer ..
func NewProducer(brokers []string, groupID string, topic string) *Producer {
	p := Producer{}

	p.producer = kafkago.NewWriter(kafkago.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	})

	return &p
}

// Close ..
func (p *Producer) Close() {
	if p.producer != nil {
		p.producer.Close()
		p.producer = nil
	}
}

// Send ..
func (p *Producer) Send(msg interface{}) error {
	if p.producer == nil {
		return fmt.Errorf("Producer not initialized")
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.producer.WriteMessages(context.Background(),
		kafkago.Message{
			Key:   []byte(""),
			Value: message,
		})
}
