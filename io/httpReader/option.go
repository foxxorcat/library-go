package http_reader

import (
	"context"
	"net/http"

	"github.com/avast/retry-go"
)

type Option func(*HttpReaderOptions)

type HttpReaderOptions struct {
	Ctx       context.Context
	Size      int64 // 指定请求资源大小
	SkipCheck bool  // 跳过Range支持检测

	Client      *http.Client
	SetRequest  func(*http.Request)
	RetryOption []retry.Option
}

func SetSize(size int) Option {
	return func(hro *HttpReaderOptions) {
		hro.Size = int64(size)
	}
}

func SkipRangeCheck(skip bool) Option {
	return func(hro *HttpReaderOptions) {
		hro.SkipCheck = skip
	}
}

func SetContext(ctx context.Context) Option {
	return func(hro *HttpReaderOptions) {
		hro.Ctx = ctx
	}
}

func SetClient(client *http.Client) Option {
	return func(hro *HttpReaderOptions) {
		hro.Client = client
	}
}

func SetRequest(fn func(*http.Request)) Option {
	return func(hro *HttpReaderOptions) {
		hro.SetRequest = fn
	}
}

func SetRetryOption(ops ...retry.Option) Option {
	return func(hro *HttpReaderOptions) {
		hro.RetryOption = append(hro.RetryOption, ops...)
	}
}
