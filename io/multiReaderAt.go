package ioutils

import (
	"io"
	"sort"
)

// NewMultiReaderAt return SizeReaderAt
// 合并多个 SizeReaderAt 接口
// io.ErrUnexpectedEOF 将被转换为 io.EOF
func NewMultiReaderAt(parts ...SizeReaderAt) SizeReaderAt {
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

func (m *multiReaderAt) ReadAt(p []byte, off int64) (rn int, err error) {
	// 查找开始部分
	indexParts := sort.Search(len(m.parts), func(i int) bool {
		return m.parts[i].off+m.parts[i].Size() > off
	})

	// 计算该部分偏移
	offset := off
	if indexParts < len(m.parts) {
		offset -= m.parts[indexParts].off
	}

	for len(p) != 0 {
		// 所有部分读取完毕
		if indexParts >= len(m.parts) {
			return rn, io.EOF
		}

		n, err := m.parts[indexParts].ReadAt(p, offset)
		// 当错误不属于读取到文件末尾时返回错误
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return rn, err
		}

		rn += n
		p = p[n:]

		// 下一部分
		indexParts++
		offset = 0
	}
	return rn, nil
}

func (m *multiReaderAt) Size() (s int64) { return m.size }
