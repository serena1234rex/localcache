package localcache

import "sync/atomic"

type Register struct {
	hitCount  int32 // 命中数
	missCount int32 // miss 数
}

type RegisterAccessor interface {
	HitCount() int32
	MissCount() int32
	HitRate() float32
	TotalCount() int32
	IncrHicCount() int32
	IncrMissCount() int32
}

func CreateRegister() RegisterAccessor {
	var r RegisterAccessor
	r = &Register{}
	return r
}

func (r *Register) IncrHicCount() int32 {
	return atomic.AddInt32(&r.hitCount, 1)
}

func (r *Register) HitCount() int32 {
	return atomic.LoadInt32(&r.hitCount)
}

func (r *Register) IncrMissCount() int32 {
	return atomic.AddInt32(&r.missCount, 1)
}

func (r *Register) MissCount() int32 {
	return atomic.LoadInt32(&r.missCount)
}

func (r *Register) HitRate() float32 {
	hc, mc := r.HitCount(), r.MissCount()
	total := hc + mc
	if total == 0.0 {
		return 0.0
	}
	return float32(mc) / float32(total)
}

func (r *Register) TotalCount() int32 {
	hc, mc := r.HitCount(), r.MissCount()
	return hc + mc
}
