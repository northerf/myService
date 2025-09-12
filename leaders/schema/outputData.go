package schema

type LeaderOutput struct {
	UserID   int64  `json:"user_id"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
	Username string `json:"username"`
}
