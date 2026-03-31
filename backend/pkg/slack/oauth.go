package slack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kietle/zenreply/config"
	slacklib "github.com/slack-go/slack"
)

const userScopes = "chat:write,im:history,im:read,mpim:history,mpim:read,channels:history,groups:history,users:read,users:read.email"

type SlackOAuth struct {
	cfg *config.Config
}

type SlackAuthResult struct {
	AccessToken string `json:"access_token"`
	BotToken    string `json:"bot_token"`
	TokenScope  string `json:"token_scope"`
	SlackUserID string `json:"slack_user_id"`
	SlackTeamID string `json:"slack_team_id"`
	SlackName   string `json:"slack_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
}

func NewSlackOAuth(cfg *config.Config) *SlackOAuth {
	return &SlackOAuth{
		cfg: cfg,
	}
}

// https://docs.slack.dev/authentication/installing-with-oauth/
func (s *SlackOAuth) BuildAuthURL(state string) string {
	return fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&user_scope=%s&redirect_uri=%s&state=%s",
		s.cfg.Slack.ClientID,
		userScopes,
		s.cfg.Slack.RedirectURL,
		state,
	)
}

func (s *SlackOAuth) ExchangeCodeForToken(ctx context.Context, code string) (*SlackAuthResult, error) {
	resp, err := slacklib.GetOAuthV2ResponseContext(ctx, http.DefaultClient, s.cfg.Slack.ClientID, s.cfg.Slack.ClientSecret, code, s.cfg.Slack.RedirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	if !resp.Ok {
		return nil, fmt.Errorf("oauth error: %s", resp.Error)
	}

	userToken := resp.AuthedUser.AccessToken
	if userToken == "" {
		return nil, fmt.Errorf("slack oauth: no user access token is response (check user_scope in OAuth configuration)")
	}

	result := &SlackAuthResult{
		AccessToken: userToken,
		TokenScope:  resp.AuthedUser.Scope,
		SlackUserID: resp.AuthedUser.ID,
		SlackTeamID: resp.Team.ID,
	}

	userClient := slacklib.New(userToken)
	userInfo, err := userClient.GetUserInfoContext(ctx, resp.AuthedUser.ID)
	if err != nil && userInfo != nil {
		result.SlackName = userInfo.Name
		result.Email = userInfo.Profile.Email
		result.AvatarURL = userInfo.Profile.Image192
	}

	return result, nil
}
