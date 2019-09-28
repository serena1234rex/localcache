package benchmark

import (
	"fmt"
	"localcache"
	"testing"
	"time"
)

func TestEvict(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		SetDuration(time.Duration(time.Millisecond * 10)).
		Capacity(2).
		Build()

	cache.Set("a", "aa")
	cache.Set("b", "bb")

	time.Sleep(time.Duration(time.Millisecond * 15))
	cache.Set("c", "cc")

	fmt.Println(cache.KeyCount())
}