package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Port        string
	BaseURL     string
	FrontendURL string
	LogLevel    string
	Env         string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	URL      string
	// SSLMode         string
	// MaxConns        int
	// MinConns        int
	// MaxConnLifetime time.Duration
	// MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	URL      string
}

type SlackConfig struct {
	ClientID      string
	ClientSecret  string
	SigningSecret string
	RedirectURL   string
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Slack    SlackConfig
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port:        getEnv("APP_PORT", "8080"),
			BaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
			FrontendURL: getEnv("APP_FRONTEND_URL", "http://localhost:4200"),
			LogLevel:    getEnv("APP_LOG_LEVEL", "info"),
			Env:         getEnv("APP_ENV", "development"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "admin"),
			Password: getEnv("POSTGRES_PASSWORD", "zenreply"),
			DB:       getEnv("POSTGRES_DB", "zenreply"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "zenreply"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: getEnvAsDuration("JWT_EXPIRY", 3600*time.Second),
		},
		Slack: SlackConfig{
			ClientID:      getEnv("SLACK_CLIENT_ID", ""),
			ClientSecret:  getEnv("SLACK_CLIENT_SECRET", ""),
			SigningSecret: getEnv("SLACK_SIGNING_SECRET", ""),
			RedirectURL:   fmt.Sprintf("%s/api/v1/slack/callback", getEnv("APP_BASE_URL", "https://localhost:8080")),
		},
	}
}

func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return fallback
}
