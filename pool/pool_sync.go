package pool

import (
	"sync"
)

// 包装 sync.Pool
type Pool[T any] struct {
	pool sync.Pool
	New  func() T
}

func (p *Pool[T]) Get() T {
	if p.pool.New == nil {
		p.pool.New = func() any { return p.New() }
	}
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
