package httpclient

import (
	"math"
	"math/rand"
	"time"
)

// Backoff interface defines contract for backoff interval startegies
type Backoff interface {
	Next(retry int) time.Duration
}

// exponentialBackoff defines backoff with exponential count
type exponentialBackoff struct {
	exponentFactor        float64
	initialTimeout        float64
	maxTimeout            float64
	maximumJitterInterval int64
}

// contantBackoff for constant Data
type constantBackoff struct {
	backoffInterval       int64
	maximumJitterInterval int64
}

// init will be executed when package is imported
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewConstantBackoff returns an instance of backoff constantly
func NewConstantBackoff(backoffInterval, maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}

	return &constantBackoff{
		backoffInterval:       int64(backoffInterval / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

// Next returns next time for retrying operation with constant strategy
func (cb *constantBackoff) Next(retry int) time.Duration {
	return (time.Duration(cb.backoffInterval) * time.Millisecond) + (time.Duration(rand.Int63n(cb.maximumJitterInterval+1)) * time.Millisecond)
}

// NewExponentialBackoff returns an instance of ExponentialBackoff
func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}

	return &exponentialBackoff{
		exponentFactor:        exponentFactor,
		initialTimeout:        float64(initialTimeout / time.Millisecond),
		maxTimeout:            float64(maxTimeout / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

// Next returns the next time for retrying operation with exponential strategy
func (eb *exponentialBackoff) Next(retry int) time.Duration {
	if retry < 0 {
		retry = 0
	}
	return time.Duration(math.Min(eb.initialTimeout*math.Pow(eb.exponentFactor, float64(retry)), eb.maxTimeout)+float64(rand.Int63n(eb.maximumJitterInterval+1))) * time.Millisecond
}
