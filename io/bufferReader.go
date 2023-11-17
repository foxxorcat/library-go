package ioutils

import (
	"errors"
	"io"
	"sync"
)

// NewBufferingReaderAt
// 将 io.Reader 读取到内存中以支持 io.ReaderAt & io.ReadSeeker
// 按需求读取到内存
func NewBufferReader(r io.Reader) *bufferReader {
	return &bufferReader{r: r}
}

type bufferReader struct {
	r    io.Reader
	buf  []byte
	lock sync.RWMutex

	offset int64
	eof    bool
}

// 读取指定长度的块到缓冲区
// 如果读取完毕，标记为eof并返回io.EOF
// io.ErrUnexpectedEOF 转化为 io.EOF
func (br *bufferReader) readToBuf(l int64) error {
	if !br.eof {
		buf := make([]byte, l)
		rn, err := io.ReadFull(br.r, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			br.eof = true
			err = io.EOF
		}
		br.buf = append(br.buf, buf[:rn]...)
		return err
	}
	return io.EOF
}

// 将剩余部分全部读取到缓存
// 仅返回非 io.EOF 错误
func (br *bufferReader) readAll() error {
	if !br.eof {
		buf, err := io.ReadAll(br.r)
		if err != nil {
			return err
		}
		br.eof = true
		br.buf = append(br.buf, buf...)
	}
	return nil
}

func (br *bufferReader) Read(p []byte) (n int, err error) {
	n, err = br.ReadAt(p, br.offset)
	br.offset += int64(n)
	return
}

func (br *bufferReader) Seek(offset int64, whence int) (int64, error) {
	br.lock.Lock()
	defer br.lock.Unlock()

	var off int64
	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off = offset + br.offset
		// 读取需要的部分到缓存
		if err := br.readToBuf(off - int64(len(br.buf))); err != nil && err != io.EOF {
			return br.offset, err
		}
	case io.SeekEnd:
		// 将所有读入缓存
		if err := br.readAll(); err != nil {
			return br.offset, err
		}
		off = int64(len(br.buf)) + offset
	}
	if off < 0 || off > int64(len(br.buf)) {
		return br.offset, errors.New("out of range")
	}

	br.offset = off
	return br.offset, nil
}

func (br *bufferReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ErrNegativeOffset
	}

	br.lock.RLock()

	// 需要读入缓存区的大小
	needSize := func() int64 { return off + int64(len(p)-len(br.buf)) }

	if need := needSize(); need > 0 {
		br.lock.RUnlock()
		br.lock.Lock()

		if need := needSize(); need > 0 {
			err = br.readToBuf(need)
		}

		br.lock.Unlock()
		br.lock.RLock()
	}
	if int64(len(br.buf)) >= off {
		n = copy(p, br.buf[off:])
	}
	br.lock.RUnlock()
	return
}

var _ ReadSeekReaderAt = (*bufferReader)(nil)
