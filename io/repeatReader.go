package ioutils

import "io"

type repeatReader struct {
	data []byte
	off  int
}

func (r *repeatReader) Reset() {
	r.off = 0
}

func (r *repeatReader) Read(p []byte) (n int, err error) {
	if r.off >= len(r.data) {
		r.Reset()
	}
	n = copy(p, r.data[r.off:])
	return
}

// RepeatReader
// 重复读取data数据，永远不会返回io.EOF
func NewRepeatReader(data ...byte) io.Reader {
	return &repeatReader{
		data: data,
	}
}
