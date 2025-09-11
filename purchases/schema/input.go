package schema

type CreatePuchaseInput struct {
	ItemName string  `json:"item_name" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
}
