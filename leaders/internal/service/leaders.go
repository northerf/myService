package service

import (
	"awesomeProject1/leaders/internal/entity"
	"awesomeProject1/leaders/schema"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserServiceClient struct {
	BaseURL string
	client  *http.Client
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserServiceClient) GetUserName(ctx context.Context, userID int64) (string, error) {
	url := fmt.Sprintf("%s/api/users/%d", c.BaseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request to user service: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request to user service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user service returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from user service: %w", err)
	}

	var user schema.User

	if err := json.Unmarshal(body, &user); err != nil {
		return "", fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user.Name, nil
}

func (s *Service) ProcessPurchaseEvent(ctx context.Context, userID int64, amount float64) error {
	const leaderboardKey = "leaders"
	return s.repo.UpdateScore(ctx, leaderboardKey, userID, amount)
}

func (s *Service) GetTopLeaders(ctx context.Context, limit int64) ([]entity.Leader, error) {
	const leaderboardKey = "leaders"
	data, err := s.repo.GetTop(ctx, leaderboardKey, limit)
	if err != nil {
		return nil, fmt.Errorf("Failed")
	}
	for i := range data {
		name, err := s.users.GetUserName(ctx, data[i].UserID)
		if err == nil {
			data[i].Username = name
		}
	}
	return data, nil
}
