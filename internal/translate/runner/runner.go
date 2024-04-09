package runner

import (
	"context"

	"github.com/telepresenceio/watchable"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/tracing"
	"github.com/ShyunnY/jaeger-operator/internal/translate"
)

type Config struct {
	config.Server

	InfraIRMap  *message.InfraIRMaps
	StatusIRMap *message.StatusIRMaps
	IrMessage   *message.IRMessage
}

type Runner struct {
	Config
}

func New(cfg *Config) Runner {
	return Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return jaegerv1a1.TranslatorComponent
}

func (r *Runner) Start(ctx context.Context) error {
	r.Logger = r.Logger.WithName(r.Name())
	go r.translateResources(ctx)
	r.Logger.Info("translator started")

	return nil
}

func (r *Runner) translateResources(ctx context.Context) {
	translator := translate.Translator{
		StatusIRMap: r.Config.StatusIRMap,
	}

	message.SubscriptionIR(
		r.IrMessage.Subscribe(ctx),
		func(update watchable.Update[string, *message.JaegerWithCtx], errCh chan error) {
			r.Logger.Info("translator takes the ir instance and handler it", "instance", update.Key)

			tracer := otel.GetTracerProvider().Tracer(consts.ReconciliationTracer)
			ctx, span := tracer.Start(update.Value.Ctx, "translate")
			span.SetAttributes(
				attribute.String("runner", jaegerv1a1.TranslatorComponent),
				attribute.String("name", update.Value.Name),
				attribute.String("namespace", update.Value.Namespace),
			)
			defer span.End()

			if update.Delete {
				// TODO: handler delete
			}

			// Jaeger resources were translated into infra IR and statusIR
			infraIR, err := translator.Translate(update.Value.Jaeger)
			if err != nil {
				errCh <- tracing.HandleErr(span, err)
				return
			}

			// push ir
			infraIR.Ctx = ctx
			r.InfraIRMap.Store(update.Value.Jaeger.Name, infraIR)
		},
	)
}
