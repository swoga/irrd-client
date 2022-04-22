package cache

import (
	"sync"
)

type Cache[TKey comparable, TValue any] interface {
	Get(TKey) (TValue, bool)
	Set(TKey, TValue)
}

func New[TKey comparable, TValue any]() Cache[TKey, TValue] {
	return &cache[TKey, TValue]{
		data: make(map[TKey]TValue),
	}
}

type cache[TKey comparable, TValue any] struct {
	mu   sync.RWMutex
	data map[TKey]TValue
}

// get cache key, safe for concurrent use
func (c *cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	i, f := c.data[key]
	return i, f
}

// set cache key, safe for concurrent use
func (c *cache[TKey, TValue]) Set(key TKey, value TValue) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}
