package httputils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	ioutils "github.com/foxxorcat/library-go/io"
)

var ErrNotSupportRange = errors.New("not support range")

// NewMultiHttpReader
// 将多个url合并为一个 ioutils.SizeReaderAtCloser
// 默认不带缓存
func NewMultiHttpReader(ctx context.Context, client *http.Client, urls ...string) (ioutils.SizeReaderAtCloser, error) {
	rs := make([]ioutils.SizeReaderAt, 0, len(urls))
	closes := make([]io.Closer, 0, len(urls))
	for _, url := range urls {
		r, err := NewHttpReader(ctx, client, url)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
		closes = append(closes, r)
	}

	return struct {
		ioutils.SizeReaderAt
		io.Closer
	}{
		SizeReaderAt: ioutils.NewMultiReaderAt(rs...),
		Closer:       ioutils.NewMultiCloser(closes...),
	}, nil
}

// NewHttpReader
// 依靠 HTTP-Range 实现 io.ReadAt 和 io.ReadSeeker
// 大小来源于 resp.ContentLength
// 默认不带缓存
func NewHttpReader(ctx context.Context, client *http.Client, url string, ops ...httpRequestOption) (ioutils.SizeReadSeekCloserAt, error) {
	newReq := func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		for _, op := range ops {
			op(req)
		}
		return req, nil
	}

	req, err := newReq()
	if err != nil {
		return nil, err
	}

	// 使用Head请求判断是否支持
	req.Method = http.MethodHead
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 当http不支持分块
	if resp.Header.Get("Accept-Ranges") != "bytes" || resp.ContentLength == -1 {
		return nil, ErrNotSupportRange
	}

	return &httpReader{
		client: client,
		req:    newReq,
		size:   resp.ContentLength,
	}, nil
}

type httpReader struct {
	client *http.Client
	req    func() (*http.Request, error)
	size   int64

	resp   *http.Response
	offset int64
}

func (r *httpReader) Read(p []byte) (n int, err error) {
	if r.resp == nil || r.resp.Close {
		req, err := r.req()
		if err != nil {
			return 0, err
		}

		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", r.offset))
		r.resp, err = r.client.Do(req)
		if err != nil {
			return 0, err
		}
	}

	n, err = r.resp.Body.Read(p)
	r.offset += int64(n)
	return
}

func (r *httpReader) Seek(offset int64, whence int) (int64, error) {
	var off int64
	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off = r.offset + offset
	case io.SeekEnd:
		off = r.size + offset
	}

	if off < 0 || off > r.size {
		return r.offset, errors.New("out of http range")
	}

	// 偏移值改变，修改具体属性
	if r.offset != off {
		if r.resp != nil && !r.resp.Close {
			r.resp.Body.Close()
		}
		r.resp = nil
		r.offset = off
	}

	return r.offset, nil
}

func (r *httpReader) Close() (err error) {
	if r.resp != nil && !r.resp.Close {
		err = r.resp.Body.Close()
	}
	r.resp = nil
	r.req = nil
	r.client = nil
	return
}

func (r *httpReader) Size() int64 {
	return r.size
}

func (r *httpReader) ReadAt(p []byte, off int64) (n int, err error) {
	// 创建请求
	req, err := r.req()
	if err != nil {
		return 0, err
	}

	end := off + int64(len(p))
	if end > r.size {
		end = r.size
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, end-1))
	resp, err := r.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return io.ReadFull(resp.Body, p[:end-off])
}

type httpRequestOption func(req *http.Request)
