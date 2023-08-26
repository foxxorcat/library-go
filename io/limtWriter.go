package ioutils

import "io"

type limitWriter struct {
	w     io.Writer
	limit int64
}

func (l *limitWriter) Write(p []byte) (n int, err error) {
	lp := len(p)
	if l.limit > 0 {
		if int64(lp) > l.limit {
			p = p[:l.limit]
		}
		l.limit -= int64(len(p))
		_, err = l.w.Write(p)
	}
	return lp, err
}

func LimitWriter(w io.Writer, limit int64) io.Writer {
	return &limitWriter{w: w, limit: limit}
}
