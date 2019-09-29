package localcache

import (
	"container/list"
	"sync"
	"time"
)

type LRUCache struct {
	basicCache
	items     map[interface{}]*list.Element
	evictList *list.List
}

type LRUItem struct {
	key        interface{}
	value      interface{}
	expiration *time.Time
	mu         sync.RWMutex
}

func (c *LRUCache) Set(key, value interface{}) error {
	var err error
	if c.serializeFunc != nil {
		value, err = c.serializeFunc(value)
		if err != nil {
			return err
		}
	}
	err = c.setValue(key, value)

	if c.addCallback != nil {
		c.addCallback(key, value)
	}
	return err
}

func (c *LRUCache) setValue(key, value interface{}) error {
	c.basicCache.mu.Lock()
	item, ok := c.items[key]
	if !ok {
		newItem := &LRUItem{
			key:   key,
			value: value,
		}
		c.items[key] = c.evictList.PushFront(newItem)
		item = c.items[key]
		c.basicCache.mu.Unlock()

		if c.isEvict() {
			c.evictItems()
		}
	} else {
		c.basicCache.mu.Unlock()
		c.evictList.MoveToFront(item)
	}

	originItem := item.Value.(*LRUItem)

	originItem.mu.Lock()
	defer originItem.mu.Unlock()

	originItem.value = value
	if c.basicCache.duration != nil {
		originItem.SetExpire(*c.duration)
	}
	return nil
}

func (c *LRUCache) Get(key interface{}) (interface{}, error) {
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

func (c *LRUCache) getValue(key interface{}) (interface{}, error) {
	c.basicCache.mu.Lock()
	item, ok := c.items[key]
	c.basicCache.mu.Unlock()

	if !ok {
		return nil, KeyNotFoundError
	}

	originItem := item.Value.(*LRUItem)
	originItem.mu.RLock()
	defer originItem.mu.RUnlock()

	ret := originItem.value
	if originItem.IsExpire(time.Now()) {
		ret = nil
		delete(c.items, key)

		c.evictList.Remove(item)
		if c.expireFunc != nil {
			c.expireFunc()
		}
	} else {
		c.evictList.MoveToFront(item)
	}

	return ret, nil
}

func (c *LRUCache) Remove(key interface{}) error {
	c.basicCache.mu.Lock()
	item, ok := c.items[key]
	c.basicCache.mu.Unlock()

	if !ok {
		return KeyNotFoundError
	} else {
		c.removeValue(item)
	}
	return nil
}

func (c *LRUCache) removeValue(item *list.Element) error {
	originItem := item.Value.(*LRUItem)
	originItem.mu.Lock()
	defer originItem.mu.Unlock()

	delete(c.items, originItem.key)
	c.evictList.Remove(item)

	item = nil
	return nil
}

func (c *LRUCache) GetAll() map[interface{}]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make(map[interface{}]interface{}, len(c.items))

	for k, item := range c.items {
		if c.Has(k) {
			items[k] = item.Value.(*LRUItem).value
		}
	}
	return items
}

func (c *LRUCache) KeyCount() int {
	return len(c.items)
}

func (c *LRUCache) Has(key interface{}) bool {
	item, ok := c.items[key]
	if !ok {
		return false
	}
	originItem := item.Value.(*LRUItem)
	return !originItem.IsExpire(time.Now())
}

// 判断是否过载
func (c *LRUCache) isEvict() bool {
	return (c.basicCache.capacity > 0 && len(c.items) > c.basicCache.capacity)
}

// 剔除超过容量的数据
func (c *LRUCache) evictItems() {
	over := c.KeyCount() - c.basicCache.capacity
	if over > 0 {
		for i := 0; i < over; i++ {
			item := c.evictList.Back()
			c.removeValue(item)
		}
	}
}

func (it *LRUItem) SetExpire(duration time.Duration) {
	t := time.Now().Add(duration)
	it.expiration = &t
}

func (it *LRUItem) IsExpire(now time.Time) bool {
	if it.expiration == nil {
		return false
	}

	if &now == nil {
		return false
	}
	return it.expiration.Before(now)
}

// new a LRU cache
func newLRUCache(builder *CacheBuilder) *LRUCache {
	cache := &LRUCache{}
	buildCache(&cache.basicCache, builder)

	cache.init()
	return cache
}

// init this cache
func (c *LRUCache) init() {
	c.items = make(map[interface{}]*list.Element, c.capacity)
	c.evictList = list.New()
}
