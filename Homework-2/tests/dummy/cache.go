package dummy

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type Cache struct {
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (c *Cache) Get(ctx context.Context, key string, value interface{}) error {
	return model.ErrorCacheMissed
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return nil
}
