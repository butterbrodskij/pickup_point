package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
)

func NewRedisDB(cfg config.Config) *Redis {
	return NewRedis(&redis.Options{
		Addr:     generateAddr(cfg),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

func generateAddr(cfg config.Config) string {
	pattern := "%s:%d"
	return fmt.Sprintf(pattern, cfg.Redis.Host, cfg.Redis.Port)
}
