package pool_test

import (
	"sync"
	"testing"

	"github.com/foxxorcat/library-go/pool"
)

func BenchmarkPoolCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool := pool.NewPoolCap(1000, func() int {
			return 0
		})

		for i := 0; i < 1000; i++ {
			pool.Put(i)
		}

		for i := 0; i < 1000; i++ {
			pool.Get()
		}
	}
}

func BenchmarkSyncPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool := sync.Pool{
			New: func() any {
				return 0
			},
		}
		for i := 0; i < 1000; i++ {
			//buf := pool.Get()
			pool.Put(i)
		}
		for i := 0; i < 1000; i++ {
			pool.Get()
		}
	}
}
