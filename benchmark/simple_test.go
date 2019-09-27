package benchmark

import (
	"fmt"
	"localcache"
	"testing"
)

func TestSimple(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		Build()

	cache.Set("key", "ddddddd")

	value, err := cache.Get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(value)
}

func TestGetAll(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		Build()

	cache.Set("a", "aa")
	cache.Set("b", "bb")
	cache.Set("c", "cc")

	fmt.Println(cache.GetAll())
}