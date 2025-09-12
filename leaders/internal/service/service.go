package service

import (
	"awesomeProject1/leaders/internal/entity"
	"awesomeProject1/leaders/internal/redis-repository"
	"context"
)

type Leaders interface {
	ProcessPurchaseEvent(ctx context.Context, userID int64, amount float64) error
	GetTopLeaders(ctx context.Context, limit int) ([]entity.Leader, error)
}

type UserProvider interface {
	GetUserName(ctx context.Context, userID int64) (string, error)
}

type Service struct {
	repo  redis_repository.LeadersRepository
	users UserProvider
}

func NewLeadersService(repo redis_repository.LeadersRepository, users UserProvider) *Service {
	return &Service{
		repo:  repo,
		users: users,
	}
}
