package schema

import "time"

type PostsOutput struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
}

type CreatePostInput struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}
