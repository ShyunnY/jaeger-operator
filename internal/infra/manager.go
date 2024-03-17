package infra

import (
	"context"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager struct {
	logger      logging.Logger
	cli         client.Client
	infraClient Client
	StatusIRMap *message.StatusIRMaps
}

func NewManager(cli client.Client, logger logging.Logger, StatusIRMap *message.StatusIRMaps) *Manager {
	return &Manager{
		cli:         cli,
		logger:      logger,
		infraClient: Client{Client: cli},
		StatusIRMap: StatusIRMap,
	}
}

func (m *Manager) BuildInfraResources(ctx context.Context, infraIR *message.InfraIR) error {

	nsName := types.NamespacedName{
		Name:      infraIR.InstanceMedata.Name,
		Namespace: infraIR.InstanceMedata.Namespace,
	}

	ic := InventoryComputer{
		namespace:    infraIR.InstanceMedata.Name,
		instanceName: infraIR.InstanceMedata.Namespace,
		cli:          m.cli,
	}
	condJaeger := new(jaegerv1a1.Jaeger)
	condJaeger.Status.Phase = "Failed"
	condJaeger.ObjectMeta = infraIR.InstanceMedata

	// create service account
	if saObj, err := ic.ComputeServiceAccount(ctx, infraIR.ServiceAccount); err != nil {
		m.logger.Error(err, "failed to compute ServiceAccount")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return saObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update ServiceAccount")

			status.SetJaegerCondition(condJaeger, "Error", metav1.ConditionFalse, "Infra", "failed to create or update ServiceAccount")
		}
	}

	// create deployment
	if deployObj, err := ic.ComputeDeployment(ctx, infraIR.Deployment); err != nil {
		m.logger.Error(err, "failed to compute Deployment")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return deployObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update Deployment")

			status.SetJaegerCondition(condJaeger, "Error", metav1.ConditionFalse, "Infra", "failed to create or update Deployment")
		}
	}

	// create service
	if servicesObj, err := ic.ComputeService(ctx, infraIR.Service); err != nil {
		m.logger.Error(err, "failed to compute Service")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return servicesObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update Service")

			status.SetJaegerCondition(condJaeger, "Error", metav1.ConditionFalse, "Infra", "failed to create or update Services")
		}
	}

	// create httpRoute
	if httpRoutesObj, err := ic.ComputeHTTPRoutes(ctx, infraIR.HTTPRoutes); err != nil {
		m.logger.Error(err, "failed to compute HTTPRoute")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return httpRoutesObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update HTTPRoutes")

			status.SetJaegerCondition(condJaeger, "Error", metav1.ConditionFalse, "Infra", "failed to create or update HTTPRoutes")
		}
	}

	if condJaeger.Status.Conditions == nil || len(condJaeger.Status.Conditions) == 0 {
		condJaeger.Status.Phase = "Success"
		status.SetJaegerCondition(condJaeger, "Success", metav1.ConditionTrue, "Infra", "success to manager infra resource")
	}

	m.StatusIRMap.Store(nsName, &condJaeger.Status)
	return nil
}
