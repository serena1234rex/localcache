package localcache

import (
	"sync"
	"sync/atomic"
	"time"
)

type SimpleCache struct {
	basicCache
	items map[interface{}]*Item
}

type Item struct {
	value      interface{}
	mu         sync.RWMutex
	expiration *time.Time
}

func (c *SimpleCache) Set(key, value interface{}) error {
	var err error
	if c.serializeFunc != nil {
		value, err = c.serializeFunc(value)
		if err != nil {
			return err
		}
	}
	err = c.setValue(key, value)
	return err
}

func (c *SimpleCache) setValue(key, value interface{}) error {
	item, ok := c.items[key];
	if !ok {
		item = &Item{}
		c.items[key] = item
		// 数量 + 1
		atomic.AddInt32(&c.size, 1)
	}

	// 针对资源上锁，不是针对表
	item.mu.Lock()
	defer item.mu.Unlock()

	item.value = value
	if c.basicCache.expiration != nil {
		item.SetExpire(*c.expiration)
	}

	return nil
}

func (c *Item) SetExpire(duration time.Duration) {
	t := time.Now().Add(duration)
	c.expiration = &t
}

func (c *SimpleCache) Get(key interface{}) (interface{}, error) {
	value, err := c.getValue(key)
	if err != nil {
		return nil, err
	}

	if c.deserializeFunc != nil {
		value, _ = c.basicCache.deserializeFunc(value)
	}

	if c.flight {
		if value != nil {
			(*c.register).IncrHicCount()
		} else {
			(*c.register).IncrMissCount()
		}
	}

	return value, nil
}

// 获取数据的私有方法
func (c *SimpleCache) getValue(key interface{}) (interface{}, error) {
	item, ok := c.items[key];
	if ok {
		item.mu.Lock()
		defer item.mu.Unlock()

		var value interface{}

		// 校验是否已经过期
		if item.IsExpire(time.Now()) {
			delete(c.items, key)
			atomic.AddInt32(&c.size, -1)
			value = nil
			// 执行过期策略
			c.expireFunc()
		} else {
			value = item.value
		}

		return value, nil
	} else {
		return nil, nil
	}
}

func (c *SimpleCache) Remove(key interface{}) error {
	item, ok := c.items[key]
	if ok {
		item.mu.Lock()
		delete(c.items, key)
		item.mu.Unlock()
		item = nil
	}
	return nil
}

func (c *SimpleCache) GetAll() map[interface{}]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make(map[interface{}]interface{}, len(c.items))

	for k, item := range c.items {
		if c.Has(k) {
			items[k] = item.value
		}
	}
	return items
}

func (c *SimpleCache) KeyCount() int32 {
	return c.size
}

func (c *SimpleCache) Has(key interface{}) bool {
	item, ok := c.items[key]
	if !ok {
		return false
	}
	return !item.IsExpire(time.Now())
}

func (item *Item) IsExpire(now time.Time) bool {
	if item.expiration == nil {
		return false
	}

	if &now == nil {
		return false
	}
	return item.expiration.Before(now)
}

func newSimpleCache(builder *CacheBuilder) *SimpleCache {
	cache := &SimpleCache{}
	cache.items = make(map[interface{}]*Item)

	buildCache(&cache.basicCache, builder)
	return cache
}
