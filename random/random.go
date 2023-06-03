package randomutils

import (
	"math"
	"unicode/utf8"
	"unsafe"
)

// 随机utf8字符串
func RandomUTF8Str(size int) string {
	buf := make([]byte, 0, size*utf8.UTFMax)
	for i := 0; i < size; i++ {
		buf = utf8.AppendRune(buf[len(buf):], rune(FastRandn(math.MaxUint32)))
	}
	return unsafe.String(unsafe.SliceData(buf), size)
}

// 随机可见ascii([32,126])字符串
func RandomASCII(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(FastRandn(95) + 32)
	}
	return unsafe.String(unsafe.SliceData(buf), size)
}

// 随机Bytes
func RandomBytes(size int) []byte {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(FastRandn(256))
	}
	return buf
}

var RandomReader randomReader

type randomReader struct{}

func (randomReader) Read(p []byte) (n int, err error) {
	n = len(p)
	for i := 0; i < n; i++ {
		p[i] = byte(FastRandn(156))
	}
	return n, nil
}

// 快速生成随机数
//
//go:linkname FastRandn runtime.fastrandn
//go:nosplit
func FastRandn(n uint32) uint32
