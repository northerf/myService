package redis_repository

import (
	"awesomeProject1/leaders/internal/entity"
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (r *RedisRepository) UpdateScore(ctx context.Context, leaderboardKey string, userID int64, increment float64) error {
	member := fmt.Sprintf("user:%d", userID)

	err := r.client.ZIncrBy(ctx, leaderboardKey, increment, member).Err()
	if err != nil {
		return fmt.Errorf("failed to update score for user %d: %w", userID, err)
	}
	return nil
}

func (r *RedisRepository) GetTop(ctx context.Context, leaderboardKey string, limit int64) ([]entity.Leader, error) {

	results, err := r.client.ZRevRangeWithScores(ctx, leaderboardKey, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get top leaders from redis-repository: %w", err)
	}

	leaders := make([]entity.Leader, 0, len(results))

	for i, z := range results {
		memberStr, ok := z.Member.(string)
		if !ok {
			continue
		}

		userIDStr := strings.TrimPrefix(memberStr, "user:")

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return nil, err
		}

		leaders = append(leaders, entity.Leader{
			UserID: userID,
			Score:  z.Score,
			Rank:   i + 1,
		})
	}

	return leaders, nil
}
