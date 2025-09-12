package entity

import "time"

type Post struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	Content   string     `db:"content"`
	CreatedAt *time.Time `db:"created_at"`
}
