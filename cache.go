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
	KeyCount() int32                          // key 的数量
	Has(key interface{}) bool                 // 校验 key 是否存在
}

type (
	ADDCallback     func(interface{})                      // 加入元素后的回调
	SerializeFunc   func(interface{}) (interface{}, error) // 序列化
	DeserializeFunc func(interface{}) (interface{}, error) // 反序列化
	ExpireFunc      func()                                 // 超时函数
)

type basicCache struct {
	size       int32             // key 的数量
	expiration *time.Duration    // 过期时间
	register   *RegisterAccessor // 计数器
	flight     bool              // 是否启动飞行器
	mu         sync.RWMutex

	serializeFunc   SerializeFunc
	deserializeFunc DeserializeFunc
	expireFunc      ExpireFunc
}

// 组织器
type CacheBuilder struct {
	size            int32
	mu              sync.RWMutex
	serializeFunc   SerializeFunc
	deserializeFunc DeserializeFunc
	expireFunc      ExpireFunc
	tp              string
	expiration      *time.Duration
	flight          bool
	register        *RegisterAccessor // 计数器
}

// 创建一个构造器
func Create() *CacheBuilder {
	return &CacheBuilder{
		size: 0,
		tp:   SIMPLE,
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

func (builder *CacheBuilder) ExpireFunc(fc ExpireFunc) *CacheBuilder {
	builder.expireFunc = fc
	return builder
}

func (builder *CacheBuilder) SetDuration(duration time.Duration) *CacheBuilder {
	builder.expiration = &duration
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
	c.size = cb.size
	c.expiration = cb.expiration
	c.flight = cb.flight
	c.register = cb.register
}
