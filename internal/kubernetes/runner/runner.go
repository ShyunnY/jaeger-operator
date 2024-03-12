package runner

import (
	"context"
	"github.com/ShyunnY/jaeger-operator/internal/kubernetes"
	"github.com/ShyunnY/jaeger-operator/internal/message"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/config"
)

type Config struct {
	config.Server
	IrMessage *message.IRMessage
}

type Runner struct {
	Config
}

func New(cfg *Config) Runner {
	return Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return jaegerv1a1.ReconcilerComponent
}

func (r *Runner) Start(ctx context.Context) error {
	r.Logger = r.Logger.WithValues("runner", r.Name())

	// create kubernetes manager and controller
	manager, err := kubernetes.New(r.Server, r.IrMessage)
	if err != nil {
		r.Logger.Error(err, "failed to create kuberntes controller")
		return err
	}

	go func() {
		if err := manager.Start(ctx); err != nil {
			r.Logger.Error(err, "wrong shutdown of manager")
		}
	}()

	return nil
}
