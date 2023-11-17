package systemutil

import "unsafe"

func Str2Byte(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Byte2Str(b []byte) string {
	return unsafe.String(&b[0], len(b))
}
