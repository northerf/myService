package schema

import "time"

type PurchaseOutput struct {
	ID          int64      `json:"id"`
	ItemName    string     `json:"item_name"`
	Amount      float64    `json:"amount"`
	PurchasedAt *time.Time `json:"purchased_at"`
}
