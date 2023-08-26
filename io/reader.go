package ioutils

import "io"

type Reader = io.Reader
type Closer = io.Closer
type Size interface{ Size() int64 }

/* ReaderAt */
type ReaderAt = io.ReaderAt
type ReadAtCloser interface {
	ReaderAt
	Closer
}
type SizeReaderAt interface {
	Size
	ReaderAt
}
type SizeReaderAtCloser interface {
	Size
	ReaderAt
	Closer
}

/* ReadSeeker */
type ReadSeeker = io.ReadSeeker
type ReadSeekCloser = io.ReadSeekCloser
type SizeReadSeeker interface {
	Size
	ReadSeeker
}
type SizeReadSeekCloser interface {
	Size
	ReadSeeker
	Closer
}

/* ReadSeekReaderAt */
type ReadSeekReaderAt interface {
	ReadSeeker
	ReaderAt
}
type ReadSeekReadAtCloser struct {
	ReadSeekReaderAt
	Closer
}
type SizeReadSeekReaderAt interface {
	Size
	ReadSeekReaderAt
}
type SizeReadSeekReadAtCloser interface {
	Size
	ReadSeekReaderAt
	Closer
}
