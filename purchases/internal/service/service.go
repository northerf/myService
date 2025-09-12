package service

import (
	"awesomeProject1/purchases/internal/entity"
	"awesomeProject1/purchases/internal/kafka"
	"awesomeProject1/purchases/internal/repository"
	"context"
)

type Purchase interface {
	CreatePurchase(ctx context.Context, userID int64, itemName string, amount float64) (int64, error)
	ListPurchases(ctx context.Context, userID int64) ([]entity.Purchase, error)
}

type Service struct {
	Purchase
}

func NewService(repos *repository.Repository, producer *kafka.Producer) *Service {
	return &Service{
		Purchase: NewPurchasesServices(repos.PurchaseRepository, producer),
	}
}
