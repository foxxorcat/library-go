package ioutils

// 无限读取零值
var Zero zeroReader

type zeroReader struct{}

func (zeroReader) Read(p []byte) (n int, err error) {
	n = len(p)
	for i := 0; i < len(p); i++ {
		p[i] = 0
	}
	return n, nil
}
