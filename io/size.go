package ioutils

import (
	"io"

	"github.com/pkg/errors"
)

// GetStreamSize
// 使用 Size() 方法获取大小
func GetStreamSize(r any) (int64, error) {
	switch v := r.(type) {
	case interface{ Size() int64 }:
		return v.Size(), nil
	case interface{ Size() int32 }:
		return int64(v.Size()), nil
	case interface{ Size() int16 }:
		return int64(v.Size()), nil
	case interface{ Size() int8 }:
		return int64(v.Size()), nil
	case interface{ Size() int }:
		return int64(v.Size()), nil
	}
	return 0, errors.Errorf("input must be of Size() method")
}

// GetStreamLen
// 使用 Len() 方法获取可读取部分大小
func GetStreamLen(r any) (int64, error) {
	switch v := r.(type) {
	case interface{ Len() int64 }:
		return v.Len(), nil
	case interface{ Len() int32 }:
		return int64(v.Len()), nil
	case interface{ Len() int16 }:
		return int64(v.Len()), nil
	case interface{ Len() int8 }:
		return int64(v.Len()), nil
	case interface{ Len() int }:
		return int64(v.Len()), nil

	}
	return 0, errors.New("input func must be an Len()")
}

// StreamSizeBySeeking
// 通过 io.Seeker 方法获取文件大小
// 也可用于验证 io.Seeker 是否可用
// @param all 是否返回总大小，而不是可读大小
func StreamSizeBySeeking(s io.Reader, all bool) (int64, error) {
	v, ok := s.(io.Seeker)
	if !ok {
		return 0, errors.Errorf("input type must be an and io.Seeker")
	}

	// 获取当前位置
	currentPosition, err := v.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, errors.WithMessage(err, "getting current offset:")
	}
	// 获取文件长度
	maxPosition, err := v.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, errors.WithMessage(err, "fast-forwarding to end:")
	}
	// 回到原始位置
	_, err = v.Seek(currentPosition, io.SeekStart)
	if err != nil {
		return 0, errors.WithMessagef(err, "returning to prior offset %d:", currentPosition)
	}
	if all {
		return maxPosition, nil
	}
	return maxPosition - currentPosition, nil
}
