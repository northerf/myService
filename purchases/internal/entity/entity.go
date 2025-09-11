package entity

import "time"

type Purchase struct {
	ID          int64      `db:"id"`
	UserID      int64      `db:"user_id"`
	ItemName    string     `db:"item_name"`
	Amount      float64    `db:"amount"`
	PurchasedAt *time.Time `db:"purchased_at"`
}
