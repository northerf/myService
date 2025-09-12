package repository

import (
	"awesomeProject1/posts/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)
	ListRecent(ctx context.Context, limit int) ([]entity.Post, error)
}

type Repository struct {
	PostsRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PostsRepository: NewPostsRepo(db),
	}
}
