package kubernetes

import (
	"fmt"
	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func New(cfg config.Server, irMessage *message.IRMessage, statusIRMaps *message.StatusIRMaps) (manager.Manager, error) {

	ctrl.SetLogger(cfg.Logger.Logger)

	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	mgrOpts := manager.Options{
		Scheme: jaeger.GetScheme(),
		Logger: cfg.Logger.Logger,

		// TODO: later we need to add leader-election
		LeaderElection:         false,
		HealthProbeBindAddress: ":8081",
		// TODO: later we need to add webhook server
		WebhookServer: nil,
	}

	mgr, err := ctrl.NewManager(restConfig, mgrOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Manager: %w", err)
	}

	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("failed to set ready probe to manager: %w", err)
	}

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("failed to set health probe to manager: %w", err)
	}

	// setup core reconcile controller
	if err = NewJaegerController(cfg, mgr, irMessage, statusIRMaps); err != nil {
		return nil, err
	}

	return mgr, nil
}
