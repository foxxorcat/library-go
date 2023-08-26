package ioutils

import (
	"errors"
	"io"
	"sort"

	math_utils "github.com/foxxorcat/library-go/math"
)

// 合并多个关闭接口
func MultiCloser(closes ...io.Closer) *multiCloser {
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

// 合并多个 SizeReaderAt 接口
func MultiReaderAt(parts ...SizeReaderAt) SizeReaderAt {
	m := &multiReaderAt{parts: make([]offsetPart, 0, len(parts))}
	for _, p := range parts {
		m.parts = append(m.parts, offsetPart{m.size, p})
		m.size += p.Size()
	}
	return m
}

type offsetPart struct {
	off int64
	SizeReaderAt
}

type multiReaderAt struct {
	parts []offsetPart
	size  int64
}

func (m *multiReaderAt) ReadAt(p []byte, offset int64) (rn int, err error) {
	if offset < 0 {
		return 0, errors.New("negative offset")
	}

	// 超过文件可读范围
	if m.size <= offset {
		return 0, io.EOF
	}

	// 查找开始Part
	indexParts := sort.Search(len(m.parts), func(i int) bool {
		return m.parts[i].off+m.parts[i].Size() > offset
	})

	// 计算从该Part开始的偏移
	if indexParts < len(m.parts) {
		offset -= m.parts[indexParts].off
	}

	for len(p) != 0 {
		// TODO io.ReadAt规定 未读满p,必然返回错误。读满也可能返回io.EOF

		// 所有Part读取完毕
		if indexParts >= len(m.parts) {
			return rn, io.EOF
		}

		// 读取该Part
		part := m.parts[indexParts]
		n, err := part.ReadAt(p[:math_utils.Min(part.Size()-offset, len(p))], offset)

		rn += n
		p = p[n:]

		// 当错误不属于读取到文件末尾时
		if err != nil && err != io.EOF {
			return rn, err
		}

		// 下一部分
		indexParts++
		offset = 0
	}
	return rn, nil
}

func (m *multiReaderAt) Size() (s int64) { return m.size }
