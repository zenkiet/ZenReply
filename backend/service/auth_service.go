package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/pkg/slack"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	cfg        *config.Config
	slackOAuth *slack.SlackOAuth
	rdb        *redis.Client
}

func NewOAuthService(cfg *config.Config, rdb *redis.Client, slackOAuth *slack.SlackOAuth) *AuthService {
	return &AuthService{
		cfg:        cfg,
		rdb:        rdb,
		slackOAuth: slackOAuth,
	}
}

func (s *AuthService) ValidateOAuthState(ctx context.Context, state string) (bool, error) {
	value, err := s.rdb.Get(ctx, fmt.Sprintf("oauth:state:%s", state)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to get OAuth state: %w", err)
	}
	return value == state, nil
}

func (s *AuthService) BuildAuthURL(ctx context.Context) (string, string, error) {
	// Generate a new OAuth state
	state := uuid.New().String()
	s.rdb.Set(ctx, fmt.Sprintf("oauth:state:%s", state), state, time.Hour*24)

	// Build the Slack OAuth URL
	url := s.slackOAuth.BuildAuthURL(state)
	return url, state, nil
}
