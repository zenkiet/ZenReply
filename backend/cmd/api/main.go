//	@title ZenReply API
//	@version 1.0.0
//	@description ZenReply is an intelligent Slack auto-reply system for deep work sessions.
//	@termsOfService https://reply.zenkiet.com/terms
//
//	@contact.name ZenKiet
//	@contact.url https://zenkiet.com
//	@contact.email kietgolx65234@gmail.com
//
//	@license.name Apache 2.0
//	@license.url https://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host localhost:8080
//	@basePath	/api/v1
//	@schemes http https
//
//	@securityDefinitions.apikey	BearerAuth
//	@in header
//	@name Authorization
//	@description Enter: Bearer <your-jwt-token>
//
//	@tag.name system
//	@tag.description System health and diagnostics
//	@tag.name auth
//	@tag.description Slack OAuth 2.0 authentication flow
//	@tag.name users
//	@tag.description User profile management
//	@tag.name deep-work
//	@tag.description Deep work session management
//	@tag.name settings
//	@tag.description User auto-reply configuration
//	@tag.name logs
//	@tag.description Auto-reply message history
//	@tag.name slack
//	@tag.description Slack Events API webhook

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/handler"
	"github.com/kietle/zenreply/pkg/database"
	"github.com/kietle/zenreply/pkg/logger"
	"github.com/kietle/zenreply/route"
)

func main() {
	//--- Config ---
	cfg := config.Load()
	log := logger.New(cfg.App.LogLevel)

	slog.SetDefault(log)

	log.Info("starting ZenReply API",
		slog.String("version", "1.0.0"),
		slog.String("port", cfg.App.Port),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//--- Database ---
	db, err := database.NewPostgres(ctx, &cfg.Postgres)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
	}
	defer db.Close()

	//--- Redis ---
	rdb, err := database.NewRedis(ctx, &cfg.Redis)
	if err != nil {
		log.Error("Failed to connect to redis", "error", err)
	}
	defer rdb.Close()

	//--- Migration ---
	if err := database.RunMigrations(ctx, db, log); err != nil {
		log.Error("Failed to run migrations", "error", err)
	}
	log.Info("migrations applied successfully")

	// --- Handler ---
	h := handler.New(cfg, db, rdb)

	// --- HTTP Server ---
	router := route.Setup(cfg, h)
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info("HTTP server is running on port " + cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down HTTP server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	log.Info("HTTP server shutdown successfully")
	os.Exit(0)
}
