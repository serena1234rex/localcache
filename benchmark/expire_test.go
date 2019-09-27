package benchmark

import (
	"fmt"
	"localcache"
	"testing"
	"time"
)

func TestExpire(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		SetDuration(time.Millisecond * 10).
		Build()

	cache.Set("boy", "yes")

	time.Sleep(time.Duration(time.Millisecond * 2))

	value, _ := cache.Get("boy")
	fmt.Println(value)
}
