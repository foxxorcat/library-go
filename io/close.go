package ioutils

type sizeReadSeekCloserAt struct {
	SizeReadSeekerAt
	close func() error
}

func (s *sizeReadSeekCloserAt) Close() error {
	if s.close != nil {
		return s.close()
	}
	return nil
}
func NopCloserSizeReadSeekerAt(r SizeReadSeekerAt) SizeReadSeekCloserAt {
	return &sizeReadSeekCloserAt{
		SizeReadSeekerAt: r,
	}
}

func WarpCloserSizeReadSeekerAt(r SizeReadSeekerAt, close func() error) ReadSeekCloserAt {
	return &readSeekCloserAt{
		ReadSeekerAt: r,
		close:        close,
	}
}

type readSeekCloserAt struct {
	ReadSeekerAt
	close func() error
}

func (s *readSeekCloserAt) Close() error {
	if s.close != nil {
		return s.close()
	}
	return nil
}
func NopCloserReadSeekerAt(r ReadSeekerAt) ReadSeekCloserAt {
	return &readSeekCloserAt{
		ReadSeekerAt: r,
	}
}

func WarpCloserReadSeekerAt(r ReadSeekerAt, close func() error) ReadSeekCloserAt {
	return &readSeekCloserAt{
		ReadSeekerAt: r,
		close:        close,
	}
}

type sizeReaderAtCloser struct {
	SizeReaderAt
	close func() error
}

func (s *sizeReaderAtCloser) Close() error {
	if s.close != nil {
		return s.close()
	}
	return nil
}

func NopCloserSizeReaderAt(r SizeReaderAt) SizeReaderAtCloser {
	return &sizeReaderAtCloser{
		SizeReaderAt: r,
	}
}

func WarpCloserSizeReaderAt(r SizeReaderAt, close func() error) SizeReaderAtCloser {
	return &sizeReaderAtCloser{
		SizeReaderAt: r,
		close:        close,
	}
}
