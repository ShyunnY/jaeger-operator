package message

import (
	"context"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type JaegerWithCtx struct {
	*jaegerv1a1.Jaeger
	Ctx context.Context
}

func (j JaegerWithCtx) DeepCopy() *JaegerWithCtx {
	out := new(JaegerWithCtx)
	out.Jaeger = j.Jaeger.DeepCopy()
	out.Ctx = j.Ctx

	return out
}

type JaegerStatusWithCtx struct {
	*jaegerv1a1.JaegerStatus
	ctx context.Context
}

func (j *JaegerStatusWithCtx) DeepCopy() *JaegerStatusWithCtx {
	out := new(JaegerStatusWithCtx)
	out.JaegerStatus = j.JaegerStatus.DeepCopy()
	out.ctx = j.ctx

	return out
}
