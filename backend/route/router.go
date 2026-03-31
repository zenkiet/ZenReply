package route

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/handler"
	"github.com/kietle/zenreply/pkg/middleware"
	"github.com/kietle/zenreply/pkg/response"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *config.Config, h *handler.Handler, logger *slog.Logger) *gin.Engine {
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// --- Middleware ---
	r.Use(middleware.Logger(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS([]string{cfg.App.BaseURL, cfg.App.FrontendURL}))

	// --- Swagger UI ---
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/scalar", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, scalarHTML(cfg.App.BaseURL))
	})
	r.GET("/openapi.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// --- Systems ---
	r.GET("/health", h.HealthCheck)

	// --- API Routes ---
	v1 := r.Group("/api/v1")
	auth := v1.Group("/slack")
	{
		auth.GET("/auth", h.SlackAuthURL)
		auth.GET("/callback", h.SlackAuthCallback)
	}

	//--- Not Found Handler ---
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "the requested endpoint does not exist")
	})

	return r
}

func scalarHTML(baseURL string) string {
	return `<!doctype html>
<html>
  <head>
    <title>ZenReply API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body { margin: 0; padding: 0; }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="` + baseURL + `/openapi.json"
      data-configuration='{"theme":"purple","layout":"modern","defaultHttpClient":{"targetKey":"javascript","clientKey":"fetch"},"hideDownloadButton":false}'
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
}
