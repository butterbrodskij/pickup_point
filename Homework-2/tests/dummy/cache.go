package dummy

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

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
