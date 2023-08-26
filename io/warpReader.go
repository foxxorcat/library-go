package ioutils

import "io"

// 包装 ReadSeekReaderAt 接口为 SizeReadSeekReadAtCloser 接口
// 通过 io.Seek 接口获取大小
func WarpReadSeekReaderAtAddSizeCloser(r ReadSeekReaderAt) (*readSeekReaderAtAddSizeCloser, error) {
	size, err := StreamSizeBySeeking(r, true)
	if err != nil {
		return nil, err
	}

	nr := readSeekReaderAtAddSizeCloser{
		ReadSeekReaderAt: r,
		size:             size,
	}
	if c, ok := r.(Closer); ok {
		nr.close = c
	}

	return &nr, nil
}

type readSeekReaderAtAddSizeCloser struct {
	ReadSeekReaderAt
	size  int64
	close io.Closer
}

func (s *readSeekReaderAtAddSizeCloser) Size() int64 {
	return s.size
}

func (s *readSeekReaderAtAddSizeCloser) Close() error {
	if s.close != nil {
		return s.close.Close()
	}
	return nil
}

// 包装 ReadSeekReaderAt 接口为 SizeReadSeekReadAtCloser 接口
// 通过 io.Seek 接口获取大小
func WarpReadSeekerAddSizeCloser(r ReadSeeker) (*readSeekerAddSizeCloser, error) {
	size, err := StreamSizeBySeeking(r, true)
	if err != nil {
		return nil, err
	}

	nr := readSeekerAddSizeCloser{
		ReadSeeker: r,
		size:       size,
	}
	if c, ok := r.(Closer); ok {
		nr.close = c
	}

	return &nr, nil
}

type readSeekerAddSizeCloser struct {
	ReadSeeker
	size  int64
	close io.Closer
}

func (s *readSeekerAddSizeCloser) Size() int64 {
	return s.size
}

func (s *readSeekerAddSizeCloser) Close() error {
	if s.close != nil {
		return s.close.Close()
	}
	return nil
}
