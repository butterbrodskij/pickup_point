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

func (c *Cache) SetPickPoint(id int64, point model.PickPoint) {
}

func (c *Cache) GetPickPoint(id int64) (model.PickPoint, error) {
	return model.PickPoint{}, model.ErrorCacheMissed
}

func (c *Cache) DeletePickPoint(id int64) {
}

func (c *Cache) UpdatePickPoint(id int64, point model.PickPoint) {
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return "", model.ErrorCacheMissed
}

func (c *Cache) Delete(ctx context.Context, key ...string) error {
	return nil
}
