package gotrace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/sunshineOfficial/golib/config"
	"go.opentelemetry.io/otel/attribute"

	"github.com/sunshineOfficial/golib/golog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	_sunshineEnvironmentKey attribute.Key = "sunshine.environment"
)

type TraceProvider interface {
	trace.TracerProvider
	Shutdown(ctx context.Context) error
}

var _ TraceProvider = (*sdktrace.TracerProvider)(nil)

func NewProvider(log golog.Logger, exporter sdktrace.SpanExporter, serviceName string) (*sdktrace.TracerProvider, error) {
	log = log.WithTags("traces")
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Debugf("Неизвестная ошибка с трейсами: %v", err)
	}))

	res, err := GetRequiredResource(serviceName)
	if err != nil {
		return nil, fmt.Errorf("can't create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter))

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

func GetRequiredResource(serviceName string) (*resource.Resource, error) {
	sdkInfo := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.TelemetrySDKName("opentelemetry"),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersion(sdk.Version()))

	serviceInfo := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		_sunshineEnvironmentKey.String(config.GetEnvironmentName()),
	)

	return resource.Merge(
		sdkInfo,
		serviceInfo,
	)
}
