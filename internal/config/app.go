package config

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"melvadTestTask/internal/interfaces"
)

type App struct {
	DB          interfaces.Postgres
	RedisClient interfaces.Redis
}

func GetAppConfig(db *sqlx.DB, redisClient *redis.Client) *App {
	return &App{
		DB:          db,
		RedisClient: redisClient,
	}
}
