package schema

import "time"

type PostsOutput struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
}
