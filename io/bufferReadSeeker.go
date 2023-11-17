package ioutils

import (
	"io"
	"sync"

	"github.com/foxxorcat/library-go/pool"
	lru "github.com/hashicorp/golang-lru/v2"
)

// NewReaderAtBuffer return ReadSeekCloserAt
// 基于lru为io.ReadSeeker提供缓存支持
// io.ErrUnexpectedEOF 转换为 io.EOF
// @param blockSize 缓存块大小。
// @param blockNum 缓存块数量.
func NewBufferReadSeeker(r io.ReadSeeker, blockSize int, blockNum int) *bufferReadSeeker {
	pool := pool.NewPoolCap(blockNum, func() []byte {
		return make([]byte, blockSize)
	})
	cache, err := lru.NewWithEvict(blockNum, func(key int, value []byte) {
		pool.Put(value)
	})
	if err != nil {
		panic(err)
	}

	br := &bufferReadSeeker{
		r:           r,
		pool:        pool,
		blockSize:   blockSize,
		cacheBlocks: cache,
	}

	if c, ok := r.(io.Closer); ok {
		br.c = c
	}
	return br
}

type bufferReadSeeker struct {
	r   io.ReadSeeker
	c   io.Closer
	off int64

	lock sync.Mutex

	pool        *pool.PoolChan[[]byte]
	blockSize   int                     // 缓存块大小
	cacheBlocks *lru.Cache[int, []byte] // 块缓存
}

func (r *bufferReadSeeker) Read(p []byte) (n int, err error) {
	n, err = r.ReadAt(p, r.off)
	r.off += int64(n)
	return
}

func (r *bufferReadSeeker) Seek(offset int64, whence int) (n int64, err error) {
	n, err = r.r.Seek(offset, whence)
	r.off = n
	return
}

// 加载块到缓存
func (r *bufferReadSeeker) loadBlock(index int) ([]byte, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if buf, ok := r.cacheBlocks.Get(index); ok {
		return buf, nil
	}

	buf := r.pool.Get()
	_, err := r.r.Seek(int64(index)*int64(r.blockSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	n, err := io.ReadFull(r.r, buf[:r.blockSize])
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = buf[:n]
	r.cacheBlocks.Add(index, buf)
	return buf, nil
}

func (r *bufferReadSeeker) ReadAt(p []byte, off int64) (rn int, err error) {
	if off < 0 {
		return 0, ErrNegativeOffset
	}

	index := int(off / int64(r.blockSize))  // 缓存块编号
	offset := int(off % int64(r.blockSize)) //  缓存块偏移
	for len(p) > 0 {
		block, err := r.loadBlock(index)
		if err != nil {
			return rn, err
		}

		// 读取范围超过 block（仅在读取末端时触发）
		if offset >= len(block) {
			return rn, io.EOF
		}

		// 读取
		n := copy(p, block[offset:])
		p = p[n:]
		rn += n
		offset += n

		// 读取下一个块
		if offset >= r.blockSize {
			index++
			offset = 0
		}
	}
	return rn, nil
}

func (r *bufferReadSeeker) Close() error {
	if r.c != nil {
		return r.c.Close()
	}
	return nil
}
