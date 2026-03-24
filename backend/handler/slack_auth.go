package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kietle/zenreply/pkg/response"
)

type SlackAuthURLResponse struct {
	URL   string `json:"url" example:"https://slack.com/oauth/v2/authorize?client_id=1234567890&scope=chat:write,im:write,channels:read,groups:read,users:read,reactions:read&redirect_uri=https://localhost:8080/api/v1/slack/auth/callback"`
	State string `json:"state" example:"123456"`
}

type SlackAuthCallbackRequest struct {
	Token   string `json:"token" example:"xoxb-1234567890"`
	UserID  string `json:"user_id" example:"U0123456789"`
	SlackID string `json:"slack_id" example:"U0123456789"`
	Name    string `json:"name" example:"Zen Kiet"`
	Email   string `json:"email" example:"zenkiet0906@gmail.com"`
	Avatar  string `json:"avatar" example:"https://avatars.slack-edge.com/2026-03-24/1234567890_1234567890_1234567890.png"`
}

// SlackAuthURL godoc
// @Summary Generate Slack OAuth URL
// @Description Generate a Slack OAuth URL for the user login
// @Tags slack
// @Accept json
// @Produce json
//
// @Success 200 {object} response.Response{data=SlackAuthURLResponse}
// @Failure 500 {object} response.Response
// @Router /slack/auth [get]
func (h *Handler) SlackAuthURL(c *gin.Context) {
	ctx := c.Request.Context()

	url, state, err := h.authSvc.BuildAuthURL(ctx)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "Slack authentication URL generated successfully", SlackAuthURLResponse{
		URL:   url,
		State: state,
	})
}
