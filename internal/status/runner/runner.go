package runner

import (
	"context"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type Config struct {
}

type Runner struct {
	Config
}

func (r *Runner) WithName() string {
	return jaegerv1a1.StatusComponent
}

func (r *Runner) Start(ctx context.Context) error {

	return nil
}
