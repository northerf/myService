package repository

import (
	"awesomeProject1/purchases/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type PurchasesRepo struct {
	db *sqlx.DB
}

func NewPurchasesRepo(db *sqlx.DB) *PurchasesRepo {
	return &PurchasesRepo{
		db: db,
	}
}

func (r *PurchasesRepo) Create(ctx context.Context, purchase *entity.Purchase) (int64, error) {
	query := `INSERT INTO purchases (item_name, amount, user_id) VALUES ($1, $2, $3) RETURNING id, purchased_at`
	row := r.db.QueryRowContext(ctx, query, purchase.ItemName, purchase.Amount, purchase.UserID)
	if err := row.Scan(&purchase.ID, &purchase.PurchasedAt); err != nil {
		return 0, err
	}
	return purchase.ID, nil
}

func (r *PurchasesRepo) GetByUserID(ctx context.Context, id int64) ([]entity.Purchase, error) {
	var purchase []entity.Purchase
	query := `SELECT id, user_id, item_name, amount, purchased_at FROM purchases WHERE user_id = $1 ORDER BY purchased_at DESC`
	err := r.db.SelectContext(ctx, &purchase, query, id)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}
