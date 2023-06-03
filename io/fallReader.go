package ioutils

import "io"

type FallReader struct {
	r    io.Reader
	fall io.Reader
	eof  bool
}

func (fr *FallReader) Read(p []byte) (n int, err error) {
	if !fr.eof && fr.r != nil {
		n, err = fr.r.Read(p)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			fr.eof = true
			err = nil
		}
		return
	}
	if fr.fall != nil {
		return fr.fall.Read(p)
	}
	return len(p), nil
}

// FallReader
// 当r读取完毕后永远不会返回io.EOF
func NewFallReader(r io.Reader) io.Reader {
	return &FallReader{
		r: r,
	}
}

// FallReader
// 当r读取完毕后永远不会返回io.EOF,并读取fall填充
// fall必须支持无限读取
func NewFallReaderWithFall(r io.Reader, fall io.Reader) io.Reader {
	return &FallReader{
		r:    r,
		fall: fall,
	}
}
