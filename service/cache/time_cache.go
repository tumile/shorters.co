package cache

import (
	"time"
)

type TimeCache interface {
	Put(key, value string, expiredTime int64)
	Get(key string) string
}

type timeCache struct {
	dict map[string]cacheEntry
}

type cacheEntry struct {
	value       string
	expiredTime int64
}

func NewTimeCache() TimeCache {
	c := &timeCache{dict: make(map[string]cacheEntry)}
	go func() {
		for range time.Tick(30 * time.Second) {
			c.evict()
		}
	}()
	return c
}

func (c *timeCache) Put(key, value string, expiredTime int64) {
	c.dict[key] = cacheEntry{value, expiredTime}
}

func (c *timeCache) Get(key string) string {
	e, ok := c.dict[key]
	if !ok || e.expiredTime <= time.Now().Unix() {
		return ""
	}
	return e.value
}

func (c *timeCache) evict() {
	now := time.Now().Unix()
	for k := range c.dict {
		if c.dict[k].expiredTime <= now {
			delete(c.dict, k)
		}
	}
}
