package api

import (
	"context"
	"log/slog"

	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/pkg/database"
	"github.com/kietle/zenreply/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.LogLevel)

	slog.SetDefault(log)

	log.Info("starting ZenReply API",
		slog.String("version", "1.0.0"),
		slog.String("port", cfg.App.Port),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB)
	db, err := database.NewPostgres(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb, err := database.NewRedis(ctx, redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer rdb.Close()
}
