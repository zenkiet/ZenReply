package slack

import (
	"fmt"

	"github.com/kietle/zenreply/config"
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
