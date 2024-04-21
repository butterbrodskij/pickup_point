package inmemorycache

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type cacheModel struct {
	str []byte
	time.Time
}

type InMemoryCache struct {
	pickPoints map[string]cacheModel
	mx         sync.RWMutex
	ttl        time.Duration
	ticker     *ticker
	wg         *sync.WaitGroup
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		pickPoints: make(map[string]cacheModel),
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

func (c *InMemoryCache) Set(_ context.Context, key string, value interface{}) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.pickPoints[key] = cacheModel{
		bytes,
		time.Now(),
	}
	return nil
}

func (c *InMemoryCache) updateUnsafe(_ context.Context, key string, value []byte) error {
	c.pickPoints[key] = cacheModel{
		value,
		time.Now(),
	}
	return nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mx.RLock()
	defer c.mx.RUnlock()
	el, ok := c.pickPoints[key]
	if !ok {
		return model.ErrorCacheMissed
	}
	err := json.Unmarshal(el.str, value)
	if err != nil {
		return err
	}
	return c.updateUnsafe(ctx, key, el.str)
}

func (c *InMemoryCache) Delete(_ context.Context, keys ...string) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	for _, key := range keys {
		delete(c.pickPoints, key)
	}
	return nil
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

func (r *InMemoryCache) Ping(ctx context.Context) error {
	return nil
}
