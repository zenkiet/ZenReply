package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kietle/zenreply/config"
	"github.com/kietle/zenreply/service"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	config  *config.Config
	db      *pgxpool.Pool
	rdb     *redis.Client
	authSvc *service.AuthService
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, authSvc *service.AuthService) *Handler {
	return &Handler{
		config:  cfg,
		db:      db,
		rdb:     rdb,
		authSvc: authSvc,
	}
}
