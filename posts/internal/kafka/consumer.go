package kafka

import (
	"awesomeProject1/posts/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type PurchaseEvent struct {
	UserID   int64   `json:"user_id"`
	ItemName string  `json:"item_name"`
	Amount   float64 `json:"amount"`
}

type Consumer struct {
	reader   *kafka.Reader
	services *service.PostsService
}

func NewConsumer(brokers []string, topic string, services *service.PostsService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "posts-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: reader, services: services}
}

func (c *Consumer) Run(ctx context.Context) {
	log.Println("Posts Kafka consumer started...")
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
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("could not unmarshal message: %v\n", err)
			continue
		}

		log.Printf("Processing purchase event: user %d, item %s, amount %.2f\n", event.UserID, event.ItemName, event.Amount)

		if event.Amount > 1000 {
			content := fmt.Sprintf("Пользователь %d совершил крупную покупку '%s' на сумму %.2f", event.UserID, event.ItemName, event.Amount)

			if _, err := c.services.CreatePost(ctx, event.UserID, content); err != nil {
				log.Printf("failed to create post: %v\n", err)
				continue
			}

			log.Printf("Post created successfully for large purchase\n")
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %v\n", err)
		}
	}
	log.Println("Posts Kafka consumer stopped.")
}
