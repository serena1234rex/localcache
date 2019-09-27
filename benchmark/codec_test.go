package benchmark

import (
	"fmt"
	"localcache"
	"testing"
)

func TestCodec(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		SerializeFunc(localcache.DefaultSerializeFunc).
		DeserializeFunc(localcache.DefaultDeserializeFunc).
		Build()

	cache.Set("try", "do")
	value,_ := cache.Get("try")
	fmt.Println(value)
}