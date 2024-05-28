package httpclient

import "time"

// Retriable defines contract for retriers to implement
type Retriable interface {
	NextInterval(retry int) time.Duration
}

// RetriableFunc is an adapter to allow the use of ordinary functions as a Retryable
type RetriableFunc func(retry int) time.Duration

// NextInterval calls f(retry)
func (f RetriableFunc) NextInterval(retry int) time.Duration {
	return f(retry)
}

type retrier struct {
	backoff Backoff
}

// NewRetrier returns retrier instance with some backoff coef value
func NewRetrier(backoff Backoff) Retriable {
	return &retrier{
		backoff: backoff, // set backoff coeficient value
	}
}

// NewRetrierFunc returns a retrier instance with a retry function defined in Retriable interface
func NewRetrierFunc(f RetriableFunc) Retriable {
	return f
}

// NextInterval returns next retriable time
func (r *retrier) NextInterval(retry int) time.Duration {
	return r.backoff.Next(retry)
}

type noRetrier struct{}

// NewNoRetrier returns a null object for retriable
func NewNoRetrier() Retriable {
	return &noRetrier{}
}

// NextInterval returns next retriable time, always 0
func (r *noRetrier) NextInterval(retry int) time.Duration {
	return 0 * time.Second
}
