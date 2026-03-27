package model

import "time"

type User struct {
	ID          string `json:"id" db:"id"`
	SlackUserID string `json:"slack_user_id" db:"slack_user_id"`
	SlackTeamID string `json:"slack_team_id" db:"slack_team_id"`
	SlackName   string `json:"slack_name" db:"slack_name"`
	Email       string `json:"email" db:"email"`
	AvatarURL   string `json:"avatar_url" db:"avatar_url"`

	AccessToken string `json:"access_token" db:"access_token"`
	TokenScope  string `json:"token_scope" db:"token_scope"`
	IsActive    bool   `json:"is_active" db:"is_active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
