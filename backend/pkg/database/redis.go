package database

import (
	"context"
	"fmt"

	"github.com/kietle/zenreply/config"
	"github.com/redis/go-redis/v9"
)

// NewRedis initializes and returns a Redis client.
func NewRedis(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	// Enable keyspace notifications for Redis
	if err := rdb.ConfigSet(ctx, "notify-keyspace-events", "K$e").Err(); err != nil {
		return nil, fmt.Errorf("failed to set Redis keyspace events: %w", err)
	}

	return rdb, nil
}
