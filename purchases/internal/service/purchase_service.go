package service

import (
	"awesomeProject1/purchases/internal/entity"
	"awesomeProject1/purchases/internal/repository"
	"context"
)

type PurchasesServices struct {
	repo repository.PurchaseRepository
}

func NewPurchasesServices(repo repository.PurchaseRepository) *PurchasesServices {
	return &PurchasesServices{repo: repo}
}

func (s *PurchasesServices) CreatePurchase(ctx context.Context, userID int64, itemName string, amount float64) (int64, error) {
	return s.repo.Create(ctx, &entity.Purchase{
		UserID:   userID,
		ItemName: itemName,
		Amount:   amount,
	})
}

func (s *PurchasesServices) ListPurchases(ctx context.Context, userID int64) ([]entity.Purchase, error) {
	purchases, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}
