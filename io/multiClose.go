package ioutils

import (
	"errors"
	"io"
)

// 合并多个关闭接口
func NewMultiCloser(closes ...io.Closer) io.Closer {
	return &multiCloser{
		closes: closes,
	}
}

type multiCloser struct {
	closes []io.Closer
}

func (m *multiCloser) Close() error {
	errs := make([]error, 0, len(m.closes))
	for _, close := range m.closes {
		errs = append(errs, close.Close())
	}
	return errors.Join(errs...)
}
