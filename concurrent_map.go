package localcache

import (
	"sync"
)

type ConcurrentMap struct {
	Table []*Node // node 节点表
	mu 	sync.Mutex
}

func (m *ConcurrentMap) casTabAt(node *Node) bool {
	m.mu.Lock()
	m.mu.Unlock()

	return true
}

