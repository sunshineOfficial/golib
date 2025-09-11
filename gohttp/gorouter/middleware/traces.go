package middleware

import (
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/internal/semconvutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	_userIdKey attribute.Key = "sunshine.auth.userId"
	_localeKey attribute.Key = "sunshine.locale"
	_errorKey  attribute.Key = "gorouter.error"
)

func Traces(service string, tracerProvider trace.TracerProvider) gorouter.Middleware {
	var (
		tracer      = tracerProvider.Tracer("github.com/sunshineOfficial/golib/gohttp/gorouter")
		propagators = otel.GetTextMapPropagator()
	)

	return func(_ gorouter.Context, nextHandler gorouter.Handler) gorouter.Handler {
		return func(c gorouter.Context) error {
			var (
				originCtx = c.Ctx()
				request   = c.Request()
				response  = c.Response()
			)

			traceCtx := propagators.Extract(originCtx, propagation.HeaderCarrier(request.Header))
			attrs := semconvutil.HTTPServerRequest(service, request)
			attrs = append(attrs,
				_userIdKey.Int(originCtx.Authorize.UserId),
				_localeKey.String(originCtx.Locale),
			)

			opts := []trace.SpanStartOption{
				trace.WithAttributes(attrs...),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			var (
				path     = request.URL.Path
				spanName = path
			)

			if len(path) > 0 {
				opts = append(opts, trace.WithAttributes(semconv.HTTPRoute(path)))
			} else {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			spanCtx, span := tracer.Start(traceCtx, spanName, opts...)
			defer span.End()

			handlerContext := gorouter.NewContext(c.Log(), c.Response(), request.WithContext(spanCtx), gorouter.WithTraces())
			err := nextHandler(handlerContext)
			if err != nil {
				span.SetAttributes(_errorKey.String(err.Error()))
				response.WriteHeader(http.StatusBadRequest)
			}

			status := response.Status()
			span.SetStatus(semconvutil.HTTPServerStatus(status))
			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}

			return err
		}
	}
}
