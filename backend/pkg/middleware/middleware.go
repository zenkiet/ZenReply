package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kietle/zenreply/pkg/response"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		protocol := c.Request.Proto
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("request", "path", path, "latency", latency, "statusCode", statusCode, "clientIP", clientIP, "userAgent", userAgent, "referer", referer, "protocol", protocol, "method", method)
	}

}

func CORS(allowedOrigins []string) gin.HandlerFunc {
	originMap := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		originMap[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if originMap[origin] || originMap["*"] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Encoding")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("panic recovered",
			slog.Any("error", recovered),
			slog.String("path", c.Request.URL.Path),
		)
		response.InternalServerError(c, "internal server error")
	})
}
