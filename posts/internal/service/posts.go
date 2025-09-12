package service

import (
	"awesomeProject1/posts/internal/entity"
	"awesomeProject1/posts/internal/repository"
	"context"
)

type PostsService struct {
	repo repository.PostsRepository
}

func NewPostsService(repo repository.PostsRepository) *PostsService {
	return &PostsService{repo: repo}
}

func (s *PostsService) CreatePost(ctx context.Context, userID int64, content string) (int64, error) {
	post := &entity.Post{
		UserID:  userID,
		Content: content,
	}
	return s.repo.CreatePost(ctx, post)
}

func (s *PostsService) ListRecent(ctx context.Context, limit int) ([]entity.Post, error) {
	return s.repo.ListRecent(ctx, limit)
}
