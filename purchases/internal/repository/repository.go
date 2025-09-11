package repository

import (
	"awesomeProject1/purchases/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entity.Purchase) (int64, error)
	GetByUserID(ctx context.Context, id int64) ([]entity.Purchase, error)
}

type Repository struct {
	PurchaseRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PurchaseRepository: NewPurchasesRepo(db),
	}
}
