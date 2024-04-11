package runner

import (
	"context"
	"fmt"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/status"
	"github.com/telepresenceio/watchable"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type Config struct {
	config.Server
	StatusIRMap *message.StatusIRMaps
}

type Runner struct {
	Config
	*status.UpdateHandler
}

func (r *Runner) Name() string {
	return jaegerv1a1.StatusComponent
}

func New(cfg *Config) Runner {
	return Runner{
		Config: *cfg,
	}
}

func (r *Runner) Start(ctx context.Context) error {
	r.Logger = r.Logger.WithName(r.Name())
	restConfig := ctrl.GetConfigOrDie()

	cli, err := client.New(restConfig, client.Options{Scheme: jaeger.GetScheme()})
	if err != nil {
		return err
	}
	r.UpdateHandler = status.NewUpdateHandler(cli, r.Logger)
	go r.UpdateHandler.Start(ctx)
	go r.subscribeStatus(ctx)

	r.Logger.Info("status handler started")
	return nil
}

func (r *Runner) subscribeStatus(ctx context.Context) {
	message.SubscriptionIR(
		r.StatusIRMap.Map.Subscribe(ctx),
		func(update watchable.Update[types.NamespacedName, *jaegerv1a1.JaegerStatus], errCh chan error) {
			r.Logger.Info("status handler takes the ir instance and handler it", "instance", update.Key.String())

			r.UpdateHandler.Write(status.Update{
				NamespacedName: update.Key,
				Object:         new(jaegerv1a1.Jaeger),
				Mutator: func(oldObj client.Object) client.Object {
					obj, ok := oldObj.(*jaegerv1a1.Jaeger)
					if !ok {
						errCh <- fmt.Errorf("unsupported object type %T", obj)
					}
					dp := obj.DeepCopy()

					cond := sortConditions(status.MergeCondition(dp.Status.Conditions, update.Value.Conditions...))
					dp.Status.Conditions = cond
					dp.Status.Phase = update.Value.Phase
					return dp
				},
			})

			r.deleteStatusIR(update.Key)
		})
}

func (r *Runner) deleteStatusIR(key types.NamespacedName) {
	r.StatusIRMap.Delete(key)
}

func sortConditions(conditions []metav1.Condition) []metav1.Condition {

	retCond := make([]metav1.Condition, 0, len(conditions))
	for i := len(conditions) - 1; i >= 0; i-- {
		retCond = append(retCond, conditions[i])
	}

	return retCond
}
