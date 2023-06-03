package ioutils

import (
	"io"
)

type crossReader struct {
	r int

	frist io.Reader
	f     int

	next io.Reader
	n    int
}

// 交换顺序
func (r *crossReader) cross() {
	r.frist, r.next = r.next, r.frist
	r.f, r.n = r.n, r.f
	r.r = r.f
}

func (r *crossReader) Read(p []byte) (n int, err error) {
	if r.r == 0 {
		r.cross()
	}
	if r.r < len(p) {
		p = p[:r.r]
	}
	n, err = r.frist.Read(p)
	r.r -= n
	return
}

// NewCrossReader
// 交替读取人r1 和 r2
// 当r1读取s1字节后,读取r2
// 当r2读取s2字节后继续读取r1
func NewCrossReader(r1, r2 io.Reader, s1, s2 int) io.Reader {
	return &crossReader{
		r:     s1,
		frist: r1,
		f:     s1,
		next:  r2,
		n:     s2,
	}
}
