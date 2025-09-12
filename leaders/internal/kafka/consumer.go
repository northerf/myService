package kafka

import (
	"awesomeProject1/leaders/internal/service"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type PurchaseEvent struct {
	UserID   int64  `json:"user_id"`
	ItemName string `json:"item_name"`
	Amount   int    `json:"amount"`
}

type Consumer struct {
	reader   *kafka.Reader
	services *service.Service
}

func NewConsumer(brokers []string, topic string, services *service.Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "leaders-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: reader, services: services}
}

func (c *Consumer) Run(ctx context.Context) {
	log.Println("Consumer started")
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("could not fetch message: %v\n", err)
			continue
		}
		var event PurchaseEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("could not unmarshal event: %v\n", err)
			c.reader.CommitMessages(ctx, msg)
			continue
		}
		log.Printf("Processing event for user %d with amount %.2f\n", event.UserID, event.Amount)
		if err := c.services.ProcessPurchaseEvent(ctx, event.UserID, float64(event.Amount)); err != nil {
			log.Printf("failed to process event: %v\n", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %v\n", err)
		}
	}
	log.Println("Kafka consumer stopped.")
}
