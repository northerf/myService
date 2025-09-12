package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type PurchaseEvent struct {
	UserID   int64   `json:"user_id"`
	ItemName string  `json:"item_name"`
	Amount   float64 `json:"amount"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &Producer{
		writer: writer,
	}
}

func (p *Producer) PublishPurchaseEvent(ctx context.Context, event PurchaseEvent) error {
	msgValue, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: msgValue,
	})
	if err == nil {
		log.Printf("Message sent successfully to topic %s", p.writer.Topic)
	}
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
