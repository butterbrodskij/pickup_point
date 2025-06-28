package redis

import (
	"context"
	"encoding/json"
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
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, time.Minute).Err()
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return res.Err()
	}
	data, err := res.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Pipeline().Del(ctx, keys...).Err()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
