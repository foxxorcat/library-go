package ioutils

import (
	"io"

	"github.com/pkg/errors"
)

/* 限制 io.ReadSeek 读取范围 */
func LimitReadSeeker(r io.ReadSeeker, offset, size int64) SizeReadSeeker {
	lr := &limitReadSeeker{
		r:      r,
		offset: offset,
		size:   size,
	}
	return lr
}

type limitReadSeeker struct {
	r      io.ReadSeeker
	offset int64
	index  int64
	size   int64
}

func (lr *limitReadSeeker) Read(p []byte) (n int, err error) {
	if lr.index >= lr.size {
		return 0, io.EOF
	}
	if i := (lr.size - lr.index); i < int64(len(p)) {
		p = p[:i]
	}
	n, err = lr.r.Read(p)
	lr.index += int64(n)
	return
}

func (lr *limitReadSeeker) Seek(offset int64, whence int) (int64, error) {
	var off int64
	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off = lr.index + offset
	case io.SeekEnd:
		off = lr.size + offset
	}

	if off < 0 || off > lr.size {
		return lr.index, errors.Errorf("out of range off:%d",off)
	}

	if lr.index != off {
		n, err := lr.r.Seek(lr.offset+off, io.SeekStart)
		lr.index = n - lr.offset
		return lr.index, err
	}
	return lr.index, nil
}

func (lr *limitReadSeeker) Size() int64 {
	return lr.size
}
