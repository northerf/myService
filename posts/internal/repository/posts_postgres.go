package repository

import (
	"awesomeProject1/posts/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type PostsRepo struct {
	db *sqlx.DB
}

func NewPostsRepo(db *sqlx.DB) *PostsRepo {
	return &PostsRepo{
		db: db,
	}
}

func (r *PostsRepo) CreatePost(ctx context.Context, post *entity.Post) (postID int64, err error) {
	query := `INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id`
	var newPostID int64

	row := r.db.QueryRowContext(ctx, query, post.UserID, post.Content)
	err = row.Scan(&newPostID)
	if err != nil {
		return 0, err
	}

	return newPostID, nil
}

func (r *PostsRepo) ListRecent(ctx context.Context, limit int) ([]entity.Post, error) {
	var posts []entity.Post
	query := `SELECT id, user_id, content, created_at FROM posts LIMIT $1`
	err := r.db.SelectContext(ctx, &posts, query, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
