package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kietle/zenreply/config"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	config *config.Config
	db     *pgxpool.Pool
	rdb    *redis.Client
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{
		config: cfg,
		db:     db,
		rdb:    rdb,
	}
}
