package message

import (
	"context"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type JaegerWithCtx struct {
	*jaegerv1a1.Jaeger
	Ctx context.Context
}

type JaegerStatusWithCtx struct {
	*jaegerv1a1.JaegerStatus
	ctx context.Context
}
