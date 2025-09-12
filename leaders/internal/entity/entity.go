package entity

type Leader struct {
	UserID   int64   `json:"user_id"`
	Score    float64 `json:"score"`
	Rank     int     `json:"rank"`
	Username string  `json:"username"`
}
