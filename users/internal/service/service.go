package service

import (
	"awesomeProject1/internal/repository"
	"context"
	"time"
)

type Auth interface {
	CreateUser(ctx context.Context, username, password, email string) (int64, error)
	GenerateToken(ctx context.Context, email string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Auth
}

func NewService(repos *repository.Repository, secretkey []byte, ttl time.Duration) *Service {
	return &Service{
		Auth: NewAuthService(repos.UserRepository, secretkey, ttl),
	}
}
