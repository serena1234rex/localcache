package benchmark

import (
	"container/heap"
	"fmt"
	"localcache"
	"testing"
)

func TestFre(t *testing.T) {
	var list localcache.FreqArr
	heap.Init(&list)

	a := &localcache.LFUItem{
		Key: "aa",
		Value: "aaa",
		Weight: 2,
	}
	heap.Push(&list, a)

	b := &localcache.LFUItem{
		Key: "bb",
		Value: "bbb",
		Weight: 3,
	}
	heap.Push(&list, b)

	c := &localcache.LFUItem{
		Key: "cc",
		Value: "ccc",
		Weight: 4,
	}
	heap.Push(&list, c)

	a.Weight = 5
	heap.Fix(&list, a.Index)

	k := heap.Pop(&list).(*localcache.LFUItem)
	fmt.Println(k.Key, k.Index, k.Weight)

	k = heap.Pop(&list).(*localcache.LFUItem)
	fmt.Println(k.Key, k.Index, k.Weight)

	k = heap.Pop(&list).(*localcache.LFUItem)
	fmt.Println(k.Key, k.Index, k.Weight)
}
