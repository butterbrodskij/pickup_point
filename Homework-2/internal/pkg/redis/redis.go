package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(opt *redis.Options) *Redis {
	return &Redis{
		redis.NewClient(opt),
	}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, time.Minute).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.Result()
}

func (r *Redis) Delete(ctx context.Context, key ...string) error {
	return r.client.Pipeline().Del(ctx, key...).Err()
}
