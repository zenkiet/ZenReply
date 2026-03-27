package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/model"
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

func (s *AuthService) HandleCallback(ctx context.Context, code, state string) (*model.User, string, error) {
	result, err := s.rdb.GetDel(ctx, fmt.Sprintf("oauth:state:%s", state)).Result()
	if errors.Is(err, redis.Nil) || result == "" {
		return nil, "", fmt.Errorf("invalid or expired state")
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to get OAuth state: %w", err)
	}

	slackAuthResult, err := s.slackOAuth.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("authService.HandleCallback: %w", err)
	}

	user := &model.User{
		SlackUserID: slackAuthResult.SlackUserID,
		SlackTeamID: slackAuthResult.SlackTeamID,
		SlackName:   slackAuthResult.SlackName,
		Email:       slackAuthResult.Email,
		AvatarURL:   slackAuthResult.AvatarURL,
		AccessToken: slackAuthResult.AccessToken,
		TokenScope:  slackAuthResult.TokenScope,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// TODO: Create settings & JWT

	return user, slackAuthResult.AccessToken, nil
}
