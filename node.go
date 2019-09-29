package localcache

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type Node struct {
	Key   interface{}
	Value interface{}
	Hash  int
	mu    sync.RWMutex
}

func (n *Node) HashCode() int {
	if n.Hash > 0 {
		return n.Hash
	}

	n.Hash = hash(n.Key)
	return n.Hash
}

func hash(raw interface{}) int {
	data, ok := raw.([]byte)
	if !ok {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(raw)
		if err != nil {
			return 0
		}

		data = buf.Bytes()
	}

	var h = 0
	for d := range data {
		h = 31*h + (d & 0xff)
	}
	return h
}

// 创造一个 node 节点
func CreateNode(key, value interface{}) *Node {
	n := &Node{
		Key:   key,
		Value: value,
	}
	n.HashCode()
	return n
}
