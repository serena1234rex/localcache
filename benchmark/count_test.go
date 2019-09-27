package benchmark

import (
	"fmt"
	"localcache"
	"sync"
	"testing"
)

func TestCount(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		Build()

	var sg sync.WaitGroup

	sg.Add(1)
	go func() {
		defer sg.Done()
		cache.Set("a", "aa")
	}()

	sg.Add(1)
	go func() {
		defer sg.Done()
		cache.Set("b", "bb")
	}()

	sg.Add(1)
	go func() {
		defer sg.Done()
		cache.Set("c", "cc")
	}()

	sg.Wait()

	fmt.Println(cache.KeyCount())
}
