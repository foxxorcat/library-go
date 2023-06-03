package ioutils

import (
	"io"

	"github.com/foxxorcat/library-go/pool"
	lru "github.com/hashicorp/golang-lru/v2"
)

// NewReaderAtBuffer return io.ReaderAt
// 基于lru为io.ReaderAt提供缓存支持
// io.ErrUnexpectedEOF 转换为 io.EOF
// @param blockSize 缓存块大小。
// @param blockNum 缓存块数量.
func NewReaderAtBuffer(r io.ReaderAt, blockSize int, blockNum int) io.ReaderAt {
	pool := pool.NewPoolCap(blockNum, func() []byte {
		return make([]byte, blockSize)
	})
	cache, err := lru.NewWithEvict(blockNum, func(key int, value []byte) {
		pool.Put(value)
	})
	if err != nil {
		panic(err)
	}

	return &readerAtBuffer{
		r:           r,
		pool:        pool,
		blockSize:   blockSize,
		cacheBlocks: cache,
	}
}

type readerAtBuffer struct {
	r           io.ReaderAt
	pool        *pool.PoolChan[[]byte]
	blockSize   int                     // 缓存块大小
	cacheBlocks *lru.Cache[int, []byte] // 块缓存
}

// 加载块到缓存
func (r *readerAtBuffer) loadBlock(index int) ([]byte, error) {
	if buf, ok := r.cacheBlocks.Get(index); ok {
		return buf, nil
	}

	buf := r.pool.Get()
	n, err := r.r.ReadAt(buf[:r.blockSize], int64(index)*int64(r.blockSize))
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	buf = buf[:n]
	r.cacheBlocks.Add(index, buf)
	return buf, nil
}

func (r *readerAtBuffer) ReadAt(p []byte, off int64) (rn int, err error) {
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
