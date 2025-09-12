package service

import (
	"awesomeProject1/posts/internal/entity"
	"awesomeProject1/posts/internal/repository"
	"context"
)

type Posts interface {
	CreatePost(ctx context.Context, userID int64, content string) (int64, error)
	ListRecent(ctx context.Context, limit int) ([]entity.Post, error)
}

type Service struct {
	Posts
}

func NewService(repos repository.PostsRepository) *Service {
	return &Service{
		Posts: NewPostsService(repos),
	}
}
