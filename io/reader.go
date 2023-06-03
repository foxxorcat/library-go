package ioutils

import "io"

type SizeReaderAt interface {
	io.ReaderAt
	Size() int64
}

type SizeReaderAtCloser interface {
	io.ReaderAt
	io.Closer
	Size() int64
}

type ReadSeekCloserAt interface {
	io.ReaderAt
	io.ReadSeekCloser
}

type ReadSeekerAt interface {
	io.ReaderAt
	io.ReadSeeker
}

type SizeReadSeekCloserAt interface {
	ReadSeekCloserAt
	Size() int64
}

type SizeReadSeekerAt interface {
	ReadSeekerAt
	Size() int64
}
