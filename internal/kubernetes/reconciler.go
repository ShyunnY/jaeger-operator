package kubernetes

import (
	"context"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/set"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
)

type jaegerReconciler struct {
	name       string
	namespaces set.Set[string]

	logger       logging.Logger
	client       client.Client
	irMessage    *message.IRMessage
	statusIRMaps *message.StatusIRMaps
}

func NewJaegerController(cfg config.Server, mgr manager.Manager, jaegerIR *message.IRMessage, statusIRMaps *message.StatusIRMaps) error {

	r := &jaegerReconciler{
		logger:       cfg.Logger,
		namespaces:   cfg.NamespaceSet,
		client:       mgr.GetClient(),
		name:         cfg.JaegerOperatorName,
		irMessage:    jaegerIR,
		statusIRMaps: statusIRMaps,
	}

	c, err := controller.New("jaeger-reconciler", mgr, controller.Options{Reconciler: r})
	if err != nil {
		r.logger.Error(err, "failed to create controller")
		return err
	}

	// watch Jaeger CR object resources
	if err = r.watchResource(mgr, c); err != nil {
		r.logger.Error(err, "failed to watch Jaeger resource by controller")
		return err
	}

	return nil
}

// Reconcile
// TODO: 与Jaeger CRD有关的修改都放置在reconcile中
// TODO: 需要处理: 已经有一个jaeger实例, 又传入一个同名的不同配置jaeger, 该怎么处理?
func (r *jaegerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	r.logger.Info("reconcile jaeger object", "instance", req.Name)

	instance := jaegerv1a1.Jaeger{}
	if err := r.client.Get(ctx, req.NamespacedName, &instance); err != nil {
		if kerrors.IsNotFound(err) {
			r.logger.Info("specified jaeger instance was not found", "instance", req.String())

			return reconcile.Result{}, nil
		} else {
			r.logger.Error(err, "failed to get jaeger instance", "instance", req.String())

			return reconcile.Result{}, nil
		}
	}

	// Before reconciliation, we need to determine if the Jaeger CR Object is managed by the current Operator
	if operated, found := instance.Labels[jaegerv1a1.OperatorByLabelKey]; found {
		if operated != r.name {
			r.logger.Info("skipping Jaeger CR as we are not owners", "our-name", r.name, "owner-name", operated)

			return reconcile.Result{}, nil
		}
	} else {
		if instance.Labels == nil {
			instance.Labels = map[string]string{}
		}
		instance.Labels[jaegerv1a1.OperatorByLabelKey] = r.name
	}

	if err := r.client.Update(ctx, &instance); err != nil {
		r.logger.Error(err, "failed to update the Jaeger", "instance", req.String())
		return reconcile.Result{}, err
	}

	normalizeJaeger(&instance)

	// Let translate calculate the K8s resources required by Jaeger.
	r.irMessage.Store(req.String(), &instance)

	return reconcile.Result{}, nil
}

// normalizeJaeger Normalizes changes to Jaeger to detect incompatible patterns and apply default value
func normalizeJaeger(instance *jaegerv1a1.Jaeger) {

	if len(instance.Spec.Type) == 0 {
		instance.Spec.Type = jaegerv1a1.AllInOneType
	}

	if len(instance.Spec.Components.Storage.Type) == 0 {
		instance.Spec.Components.Storage.Type = jaegerv1a1.MemoryStorageType
	}

	if instance.Spec.Type != jaegerv1a1.AllInOneType &&
		instance.Spec.Components.Storage.Type == jaegerv1a1.MemoryStorageType {
		instance.Spec.Type = jaegerv1a1.AllInOneType
	}

	if len(instance.Name) == 0 {
		instance.Name = "simple-jaeger"
	}

}

// hasInWatchNamespace Determine if the resource is under the monitored namespace
func (r *jaegerReconciler) hasInWatchNamespace(object client.Object) bool {
	ns := object.GetNamespace()
	if r.namespaces.Len() != 0 && !r.namespaces.Has(ns) {
		return false
	}

	return true
}

func (r *jaegerReconciler) watchResource(mgr manager.Manager, c controller.Controller) error {

	predicates := []predicate.Predicate{
		predicate.GenerationChangedPredicate{},
	}

	if r.namespaces.Len() != 0 {
		predicates = append(predicates, predicate.NewPredicateFuncs(r.hasInWatchNamespace))
	}

	if err := c.Watch(
		source.Kind(mgr.GetCache(), &jaegerv1a1.Jaeger{}),
		&EnqueueHandler{},
		predicates...,
	); err != nil {
		return err
	}

	return nil
}

// EnqueueHandler With custom enqueueHandler, we don't handle delete cases, at least not for now
type EnqueueHandler struct {
	*handler.EnqueueRequestForObject
}

func (e *EnqueueHandler) Delete(ctx context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	return
}
