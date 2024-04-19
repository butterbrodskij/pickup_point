package inmemorycache

import (
	"sync"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type cacheModel struct {
	model.PickPoint
	time.Time
}

type InMemoryCache struct {
	pickPoints map[int64]cacheModel
	mx         sync.RWMutex
	ttl        time.Duration
	ticker     *ticker
	wg         *sync.WaitGroup
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		pickPoints: make(map[int64]cacheModel),
		mx:         sync.RWMutex{},
		ttl:        time.Minute,
		ticker:     newTicker(time.Minute),
		wg:         &sync.WaitGroup{},
	}
	cache.wg.Add(1)
	go func() {
		cache.ticker.Run(cache)
		cache.wg.Done()
	}()
	return cache
}

func (c *InMemoryCache) Close() {
	close(c.ticker.stop)
	c.wg.Wait()
}

func (c *InMemoryCache) SetPickPoint(id int64, point model.PickPoint) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.pickPoints[id] = cacheModel{
		point,
		time.Now(),
	}
}

func (c *InMemoryCache) UpdatePickPoint(id int64, point model.PickPoint) {
	c.mx.Lock()
	defer c.mx.Unlock()
	_, ok := c.pickPoints[id]
	if !ok {
		c.SetPickPoint(id, point)
	}
	c.pickPoints[id] = cacheModel{
		point,
		time.Now(),
	}
}

func (c *InMemoryCache) GetPickPoint(id int64) (model.PickPoint, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	pointCache, ok := c.pickPoints[id]
	if !ok {
		return model.PickPoint{}, model.ErrorCacheMissed
	}
	return pointCache.PickPoint, nil
}

func (c *InMemoryCache) DeletePickPoint(id int64) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.pickPoints, id)
}

func (c *InMemoryCache) DeleteExpired() {
	c.mx.Lock()
	defer c.mx.Unlock()
	for id, point := range c.pickPoints {
		if c.expired(point) {
			delete(c.pickPoints, id)
		}
	}
}

func (c *InMemoryCache) expired(el cacheModel) bool {
	return time.Since(el.Time) > c.ttl
}
