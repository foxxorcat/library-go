package pool

// 由chan实现的池
// 避免 sync.Pool GC 回收问题
type PoolChan[T any] struct {
	New func() T
	c   chan T
}

func NewPoolCap[T any](maxSize int, new func() T) (bp *PoolChan[T]) {
	return &PoolChan[T]{
		New: new,
		c:   make(chan T, maxSize),
	}
}

func (p *PoolChan[T]) Get() (b T) {
	select {
	case b = <-p.c:
	default:
		b = p.New()
	}
	return
}

// 放回池，等待复用
func (p *PoolChan[T]) Put(b T) {
	select {
	case p.c <- b:
	default:
	}
}

func (p *PoolChan[T]) Len() int {
	return len(p.c)
}

func (p *PoolChan[T]) Cap() int {
	return cap(p.c)
}

// //go:linkname runtime_LoadAcquintptr runtime/internal/atomic.LoadAcquintptr
// func runtime_LoadAcquintptr(ptr *uintptr) uintptr

// //go:linkname runtime_StoreReluintptr runtime/internal/atomic.StoreReluintptr
// func runtime_StoreReluintptr(ptr *uintptr, val uintptr) uintptr

// //go:linkname runtime_procPin sync.runtime_procPin
// func runtime_procPin() int

// //go:linkname runtime_procUnpin sync.runtime_procUnpin
// func runtime_procUnpin()
