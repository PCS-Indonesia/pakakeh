package opentelemetry

import (
	"context"
	"crypto/x509"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type OpenTelemetryClient struct {
	SchemaURL   string // format = url:port
	ServiceName string
	UseSecurity bool
	Creds       *TLSCredentials
}

type TLSCredentials struct {
	Cert               *x509.CertPool
	ServerNameOverride string
}

func NewOtelClient(url, svcName string, useSecurity bool, creds *TLSCredentials) OpenTelemetryClient {
	return OpenTelemetryClient{
		SchemaURL:   url,
		ServiceName: svcName,
		UseSecurity: useSecurity,
		Creds:       creds,
	}
}

func (ot OpenTelemetryClient) InitTracer() func(context.Context) error {
	var secureOption otlptracegrpc.Option

	if ot.UseSecurity {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(ot.Creds.Cert, ot.Creds.ServerNameOverride))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(ot.SchemaURL),
		),
	)

	if err != nil {
		log.Fatal()
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", ot.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Println("Could not set resources: ", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)

	return exporter.Shutdown
}
