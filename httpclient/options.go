package httpclient

import "time"

type Option func(*CustomHttpClient)

// WithTimeout defines timeout when receiving response
func WithTimeout(timeout time.Duration) Option {
	return func(c *CustomHttpClient) {
		c.timeout = timeout
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier Retriable) Option {
	return func(c *CustomHttpClient) {
		c.retrier = retrier
	}
}

func WithRetryCount(retryCount int) Option {
	return func(c *CustomHttpClient) {
		c.retryCount = retryCount
	}
}

func WithHTTPClient(client DoReq) Option {
	return func(c *CustomHttpClient) {
		c.client = client
	}
}
