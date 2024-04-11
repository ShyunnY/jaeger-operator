package runner

import (
	"context"

	"github.com/telepresenceio/watchable"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/infra"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/tracing"
)

type Config struct {
	config.Server
	InfraMap  *message.InfraIRMaps
	StatusMap *message.StatusIRMaps
}

type Runner struct {
	*Config
	*infra.Manager
}

func New(cfg *Config) Runner {
	return Runner{Config: cfg}
}

func (r *Runner) Name() string {
	return jaegerv1a1.InfrastructureComponent
}

func (r *Runner) Start(ctx context.Context) error {
	r.Logger = r.Logger.WithName(r.Name())

	cli, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: jaeger.GetScheme()})
	if err != nil {
		return err
	}
	r.Manager = infra.NewManager(cli, r.Logger, r.StatusMap)
	go r.subscriptionInfraResource(ctx)

	r.Logger.Info("infra manager started")
	return nil
}

func (r *Runner) subscriptionInfraResource(ctx context.Context) {

	message.SubscriptionIR(r.InfraMap.Subscribe(ctx), func(update watchable.Update[string, *message.InfraIR], errCh chan error) {
		r.Logger.Info("infra manager takes the ir instance and handler it", "instance", update.Key)

		tracer := otel.GetTracerProvider().Tracer(consts.ReconciliationTracer)
		ctx, span := tracer.Start(update.Value.Ctx, "infra")
		span.SetAttributes(
			attribute.String("runner", jaegerv1a1.InfrastructureComponent),
			attribute.String("name", update.Value.InstanceMedata.Name),
			attribute.String("namespace", update.Value.InstanceMedata.Namespace),
		)
		defer span.End()

		// According to infraIR, the corresponding k8s resource is created
		if err := r.Manager.BuildInfraResources(ctx, update.Value); err != nil {
			errCh <- tracing.HandleErr(span, err)
		}

		r.deleteInfraIR(update.Key)
	})

}

func (r *Runner) deleteInfraIR(key string) {
	r.InfraMap.Delete(key)
}
