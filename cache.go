package localcache

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key, value interface{}) error         // 写入
	Get(key interface{}) (interface{}, error) // 抽取
	Remove(key interface{}) error             // 删除
	GetAll() map[interface{}]interface{}      // 获取所有
	KeyCount() int                            // key 的数量
	Has(key interface{}) bool                 // 校验 key 是否存在
}

type (
	ADDCallback     func(key, value interface{})                      // 加入元素后的回调
	SerializeFunc   func(interface{}) (interface{}, error) // 序列化
	DeserializeFunc func(interface{}) (interface{}, error) // 反序列化
	ExpireFunc      func()                                 // 超时函数
)

type basicCache struct {
	capacity   int               // 容量
	expiration *time.Duration    // 过期时间
	register   *RegisterAccessor // 计数器
	flight     bool              // 是否启动飞行器
	mu         sync.RWMutex

	serializeFunc   SerializeFunc
	deserializeFunc DeserializeFunc
	expireFunc      ExpireFunc
	addCallback 	ADDCallback
}

// 组织器
type CacheBuilder struct {
	capacity        int
	mu              sync.RWMutex
	serializeFunc   SerializeFunc
	deserializeFunc DeserializeFunc
	expireFunc      ExpireFunc
	tp              string
	expiration      *time.Duration
	flight          bool
	register        *RegisterAccessor // 计数器
	addCallback 	ADDCallback
}

// 创建一个构造器
func Create() *CacheBuilder {
	return &CacheBuilder{
		capacity: -1,
		tp:       SIMPLE,
	}
}

// 组织序列化
func (builder *CacheBuilder) SerializeFunc(fc SerializeFunc) *CacheBuilder {
	builder.serializeFunc = fc
	return builder
}

// 组织反序列化
func (builder *CacheBuilder) DeserializeFunc(fc DeserializeFunc) *CacheBuilder {
	builder.deserializeFunc = fc
	return builder
}

func (builder *CacheBuilder) Tp(tp string) *CacheBuilder {
	builder.tp = tp
	return builder
}

func (builder *CacheBuilder) AddCallback(fc ADDCallback) *CacheBuilder {
	builder.addCallback = fc
	return builder
}

func (builder *CacheBuilder) ExpireFunc(fc ExpireFunc) *CacheBuilder {
	builder.expireFunc = fc
	return builder
}

func (builder *CacheBuilder) SetDuration(duration time.Duration) *CacheBuilder {
	builder.expiration = &duration
	return builder
}

// 设置容量
func (builder *CacheBuilder) Capacity(capacity int) *CacheBuilder {
	builder.capacity = capacity
	return builder
}

// 启动飞行器
func (builder *CacheBuilder) OpenFlight(r *RegisterAccessor) *CacheBuilder {
	builder.flight = true
	builder.register = r
	return builder
}

func (builder *CacheBuilder) Build() Cache {
	if builder.tp == SIMPLE {
		return newSimpleCache(builder)
	}
	return nil
}

func buildCache(c *basicCache, cb *CacheBuilder) {
	c.deserializeFunc = cb.deserializeFunc
	c.serializeFunc = cb.serializeFunc
	c.expireFunc = cb.expireFunc
	c.serializeFunc = cb.serializeFunc
	c.capacity = cb.capacity
	c.expiration = cb.expiration
	c.flight = cb.flight
	c.register = cb.register
	c.addCallback = cb.addCallback
}
