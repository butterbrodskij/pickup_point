package inmemorycache

import "time"

type ticker struct {
	*time.Ticker
	interval time.Duration
	stop     chan struct{}
}

func newTicker(interval time.Duration) *ticker {
	return &ticker{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (t *ticker) Run(c *InMemoryCache) {
	t.Ticker = time.NewTicker(t.interval)
	for {
		select {
		case <-t.Ticker.C:
			c.DeleteExpired()
		case <-t.stop:
			t.Ticker.Stop()
			return
		}
	}
}
