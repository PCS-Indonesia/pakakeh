package httpclient

import (
	"io"
	"net/http"
)

type DoReq interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpMethods interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string, headers http.Header) (*http.Response, error)
	Post(url string, body io.Reader, headers http.Header) (*http.Response, error)
	Delete(url string, headers http.Header) (*http.Response, error)
	Put(url string, body io.Reader, headers http.Header) (*http.Response, error)
	Patch(url string, body io.Reader, headers http.Header) (*http.Response, error)
}
