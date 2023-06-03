package pool

import "sync"

// 包装 sync.Pool
type Pool[T any] struct {
	pool sync.Pool
	New  func() T
	flag bool
}

func (p *Pool[T]) Get() T {
	if !p.flag {
		p.pool.New = func() any {
			return p.New()
		}
		p.flag = true
	}
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
