package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kietle/zenreply/pkg/response"
)

type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Service   string `json:"service" example:"ZenReply API"`
	Version   string `json:"version" example:"1.0.0"`
	Timestamp string `json:"timestamp" example:"2021-01-01T00:00:00Z"`
	Database  string `json:"database" example:"ok"`
	Redis     string `json:"redis" example:"ok"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Check the health of the service
// @Tags system
// @Accept json
// @Produce json
//
//	@Success		200	{object}	response.Response{data=HealthResponse}
//	@Failure		503	{object}	response.Response
//	@Router			/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	dbStatus := "ok"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "error"
	}

	redisStatus := "ok"
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		redisStatus = "error"
	}

	data := HealthResponse{
		Status:    "ok",
		Service:   "ZenReply API",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  dbStatus,
		Redis:     redisStatus,
	}

	if dbStatus != "ok" || redisStatus != "ok" {
		response.ServiceUnavailable(c, "service is unavailable")
		return
	}

	response.OK(c, "service is available", data)
}
