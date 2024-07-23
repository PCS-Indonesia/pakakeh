package elasticapm

import (
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

// Option sets options for tracing.
type Option func(*apmMiddleware)

// WithTracer returns an Option as the tracer to use for tracing server requests.
func WithTracer(tracer *apm.Tracer) Option {
	if tracer == nil {
		return func(am *apmMiddleware) {
			am.tracer = apm.DefaultTracer
		}
	}

	return func(am *apmMiddleware) {
		am.tracer = tracer
	}
}

// WithRequestIgnorer returns a Option func(*apmMiddleware)
func WithRequestIgnorer(r apmhttp.RequestIgnorerFunc) Option {
	if r == nil {
		r = apmhttp.IgnoreNone // if r is nil, all requests will be reported.
	}

	return func(m *apmMiddleware) {
		m.requestIgnorer = r
	}
}
