package elasticapm

import (
	"net/url"

	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
)

// NewTracer returns tracer agent using the provided options.
// apmUrlStr apm url in string
// apiKey auth api key (if needed)
// svcName your application / service name
// svcVersion your application version, set "" if you never set it
func NewTracer(apmUrlStr, apiKey, svcName, svcVersion string) *apm.Tracer {
	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:    svcName,
		ServiceVersion: svcVersion,
	})

	if err != nil {
		return nil
	}

	// Set default HTTP Transport. If no URL is specified
	// then the transport will use the default URL "http://localhost:8200".
	transport, err := transport.NewHTTPTransport()
	if err != nil {
		return nil
	}
	transport.SetSecretToken(apiKey) // Set secret token or API Key if set

	var apmUrl *url.URL

	apmUrl, err = url.Parse(apmUrlStr) // convert url from string format to URL format
	if err != nil {
		return nil
	}
	transport.SetServerURL(apmUrl) // set server url. if empty it will get from default apmgin

	tracer.Transport = transport
	return tracer
}
