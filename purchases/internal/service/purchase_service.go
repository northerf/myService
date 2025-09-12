package service

import (
	"awesomeProject1/purchases/internal/entity"
	"awesomeProject1/purchases/internal/kafka"
	"awesomeProject1/purchases/internal/repository"
	"context"
	"fmt"
)

type PurchasesServices struct {
	repo     repository.PurchaseRepository
	producer *kafka.Producer
}

func NewPurchasesServices(repo repository.PurchaseRepository, producer *kafka.Producer) *PurchasesServices {
	return &PurchasesServices{
		repo:     repo,
		producer: producer,
	}
}

func (s *PurchasesServices) CreatePurchase(ctx context.Context, userID int64, itemName string, amount float64) (int64, error) {
	fmt.Printf("Creating purchase: userID=%d, item=%s, amount=%.2f\n", userID, itemName, amount)

	purchaseID, err := s.repo.Create(ctx, &entity.Purchase{
		UserID:   userID,
		ItemName: itemName,
		Amount:   amount,
	})
	if err != nil {
		fmt.Printf("Failed to create purchase: %v\n", err)
		return 0, err
	}

	fmt.Printf("Purchase created successfully with ID: %d\n", purchaseID)

	event := kafka.PurchaseEvent{
		UserID:   userID,
		ItemName: itemName,
		Amount:   amount,
	}

	if err := s.producer.PublishPurchaseEvent(ctx, event); err != nil {
		fmt.Printf("Failed to publish purchase event to Kafka: %v\n", err)
	}

	return purchaseID, nil
}

func (s *PurchasesServices) ListPurchases(ctx context.Context, userID int64) ([]entity.Purchase, error) {
	purchases, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
