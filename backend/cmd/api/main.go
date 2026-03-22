package api

import (
	"context"
	"log/slog"

	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/pkg/database"
	"github.com/kietle/zenreply/pkg/logger"
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
}
