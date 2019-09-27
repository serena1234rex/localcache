package benchmark

import (
	"fmt"
	"localcache"
	"testing"
)

func TestFlight(t *testing.T) {
	r := localcache.CreateRegister()

	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		OpenFlight(&r).
		Build()

	cache.Set("a", "aa")

	cache.Get("a")
	cache.Get("b")

	fmt.Println(r.MissCount())
}
