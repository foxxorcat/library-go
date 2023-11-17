package http_reader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	ioutils "github.com/foxxorcat/library-go/io"
	sutil "github.com/foxxorcat/library-go/system"
)

var ErrNotSupportRange = errors.New("not support range")
var ErrOutRange = errors.New("out of http range")

func NewHttpReader(method string, url string, opts ...Option) (*httpReader, error) {
	options := &HttpReaderOptions{
		Size: -1,
		RetryOption: []retry.Option{
			retry.Attempts(3),
			retry.Delay(time.Second),
			retry.DelayType(retry.BackOffDelay),
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	fetch := func(start, end int64) (*http.Response, error) {
		req, err := http.NewRequestWithContext(sutil.IFNULL(options.Ctx, context.Background()), method, url, nil)
		if err != nil {
			return nil, err
		}

		// 设置请求范围
		if start != -1 || end != -1 {
			rang := fmt.Sprintf("bytes=%s-%s",
				sutil.IFT(start != -1, strconv.FormatInt(start, 10)),
				sutil.IFT(end != -1, strconv.FormatInt(end, 10)),
			)
			req.Header.Set("Range", rang)
		}

		if options.SetRequest != nil {
			options.SetRequest(req)
		}

		// 发生请求
		resp, err := sutil.IFNULL(options.Client, http.DefaultClient).Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// 检测是否支持并获取内容大小
	if !options.SkipCheck || options.Size == -1 {
		resp, err := fetch(0, -1)
		if err != nil {
			return nil, err
		}
		// 判断是否支持
		if !options.SkipCheck && resp.Header.Get("Accept-Ranges") != "bytes" {
			return nil, ErrNotSupportRange
		}
		// 获取大小
		if options.Size == -1 {
			if _, size, ok := strings.Cut(resp.Header.Get("Content-Range"), "/"); ok {
				options.Size, _ = strconv.ParseInt(size, 10, 64)
			} else {
				options.Size = resp.ContentLength
			}
			if options.Size == -1 {
				return nil, ErrNotSupportRange
			}
		}
	}

	return &httpReader{
		size: options.Size,
		getReader: func(start, end int64) (r io.ReadCloser, err error) {
			err = retry.Do(func() error {
				resp, err := fetch(start, end)
				if err == nil {
					r = resp.Body
				}
				return err
			}, options.RetryOption...)
			return
		},
	}, nil
}

type httpReader struct {
	size      int64
	getReader func(start, end int64) (io.ReadCloser, error)

	r      io.ReadCloser
	offset int64
}

func (r *httpReader) Size() int64 {
	return r.size
}

func (r *httpReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		if r.r, err = r.getReader(r.offset, r.size); err != nil {
			return
		}
	}
	n, err = r.r.Read(p)
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
		return r.offset, ErrOutRange
	}

	if r.offset != off {
		nr, err := r.getReader(off, r.size)
		if err != nil {
			return r.offset, err
		}

		// 关闭旧的使用新的
		_ = r.Close()
		r.r = nr
		r.offset = off
	}
	return r.offset, nil
}

func (r *httpReader) Close() (err error) {
	if r.r != nil {
		err = r.r.Close()
	}
	return
}

func (r *httpReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ioutils.ErrNegativeOffset
	}

	if off >= r.size {
		return 0, io.EOF
	}

	end := off + int64(len(p))
	if end > r.size {
		end = r.size
	}

	nr, err := r.getReader(off, end)
	if err != nil {
		return
	}
	return io.ReadFull(nr, p[:end-off])
}

var _ ioutils.SizeReadSeekReadAtCloser = (*httpReader)(nil)
