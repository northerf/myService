package redis_repository

import (
	"awesomeProject1/leaders/internal/entity"
	"context"

	"github.com/redis/go-redis/v9"
)

type LeadersRepository interface {
	UpdateScore(ctx context.Context, leaderboardKey string, userID int64, increment float64) error
	GetTop(ctx context.Context, leaderboardKey string, limit int64) ([]entity.Leader, error)
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}
