package gotrace

import (
	"github.com/sunshineOfficial/golib/goctx"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	GetTracer() trace.Tracer
	NewSpan(ctx goctx.Context, name string) (goctx.Context, trace.Span)
}

type Impl struct {
	tracer trace.Tracer
}

func NewService(provider trace.TracerProvider) *Impl {
	return &Impl{
		tracer: provider.Tracer("github.com/sunshineOfficial/golib/gotrace"),
	}
}

func (i *Impl) GetTracer() trace.Tracer {
	return i.tracer
}

func (i *Impl) NewSpan(ctx goctx.Context, name string) (goctx.Context, trace.Span) {
	newCtx, span := i.tracer.Start(ctx, name)
	return ctx.CloneTo(newCtx), span
}
